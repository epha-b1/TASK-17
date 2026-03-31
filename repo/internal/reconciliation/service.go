package reconciliation

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditWriter interface {
	WriteAuditLog(ctx context.Context, actorID *string, action, resourceType string, resourceID *string, detail map[string]any) error
}

type Service struct {
	pool  *pgxpool.Pool
	audit AuditWriter
}

type RunResult struct {
	RunID              string
	ZonesChecked       int
	DiscrepanciesFound int
}

func NewService(pool *pgxpool.Pool, audit AuditWriter) *Service {
	return &Service{pool: pool, audit: audit}
}

func (s *Service) RunOnce(ctx context.Context, at time.Time) (RunResult, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return RunResult{}, err
	}
	defer tx.Rollback(ctx)

	var runID string
	err = tx.QueryRow(ctx, `
		INSERT INTO reconciliation_runs(run_at, zones_checked, discrepancies_found)
		VALUES ($1, 0, 0)
		RETURNING id::text
	`, at.UTC()).Scan(&runID)
	if err != nil {
		return RunResult{}, err
	}

	rows, err := tx.Query(ctx, `SELECT id::text FROM zones ORDER BY id`)
	if err != nil {
		return RunResult{}, err
	}
	zones := make([]string, 0)
	for rows.Next() {
		var zoneID string
		if err := rows.Scan(&zoneID); err != nil {
			rows.Close()
			return RunResult{}, err
		}
		zones = append(zones, zoneID)
	}
	rows.Close()

	zonesChecked := 0
	discrepancies := 0
	for _, zoneID := range zones {
		zonesChecked++

		var snapshotStalls int
		err = tx.QueryRow(ctx, `
			SELECT authoritative_stalls
			FROM capacity_snapshots
			WHERE zone_id=$1
			ORDER BY snapshot_at DESC
			LIMIT 1
		`, zoneID).Scan(&snapshotStalls)
		if err == pgx.ErrNoRows {
			continue
		}
		if err != nil {
			return RunResult{}, err
		}

		var eventDerived int
		err = tx.QueryRow(ctx, `
			SELECT COALESCE(SUM(
				CASE
					WHEN de.event_type IN ('stall_freed','stall_release','gate_opened') THEN 1
					WHEN de.event_type IN ('stall_occupied','stall_consumed','gate_closed') THEN -1
					ELSE 0
				END
			), 0)
			FROM device_events de
			JOIN devices d ON d.id = de.device_id
			WHERE d.zone_id = $1
		`, zoneID).Scan(&eventDerived)
		if err != nil {
			return RunResult{}, err
		}

		eventType, stallDelta, needed := DecideCompensatingEvent(snapshotStalls, eventDerived)
		if !needed {
			continue
		}
		discrepancies++

		detailJSON, _ := json.Marshal(map[string]any{
			"snapshot_stalls":      snapshotStalls,
			"event_derived_stalls": eventDerived,
			"delta":                stallDelta,
		})

		_, err = tx.Exec(ctx, `
			INSERT INTO compensating_events(
				reconciliation_run_id,
				zone_id,
				event_type,
				stall_count,
				detail
			)
			VALUES ($1, $2, $3, $4, $5::jsonb)
		`, runID, zoneID, eventType, stallDelta, string(detailJSON))
		if err != nil {
			return RunResult{}, err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE reconciliation_runs
		SET zones_checked=$2, discrepancies_found=$3
		WHERE id=$1
	`, runID, zonesChecked, discrepancies)
	if err != nil {
		return RunResult{}, err
	}

	if s.audit != nil {
		_ = s.audit.WriteAuditLog(ctx, nil, "reconciliation.run", "reconciliation_run", &runID, map[string]any{
			"zones_checked":       zonesChecked,
			"discrepancies_found": discrepancies,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return RunResult{}, err
	}

	return RunResult{RunID: runID, ZonesChecked: zonesChecked, DiscrepanciesFound: discrepancies}, nil
}

func StartScheduler(ctx context.Context, logger *slog.Logger, service *Service) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			if _, err := service.RunOnce(ctx, now.UTC()); err != nil {
				logger.Error("reconciliation run failed", "error", err)
			}
		}
	}
}
