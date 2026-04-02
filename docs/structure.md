# ParkOps — Submission Folder Structure

Task ID: 17
Project Type: fullstack
Stack: Go (Gin) + Templ + PostgreSQL

---

## ZIP Root Layout

```
17/
├── docs/
│   ├── design.md
│   ├── api-spec.md
│   ├── questions.md
│   ├── action-plan.md
│   ├── features.md
│   ├── requirements.md
│   ├── testing-plan.md
│   ├── structure.md
│   ├── AI-self-test.md
│   ├── aesthetics-assessment.md
│   ├── ai-self-test-completion-report.md
│   ├── build-order.md
│   ├── delivery-completeness-report.md
│   ├── engineering-architecture-report.md
│   ├── engineering-details-professionalism.md
│   ├── project-self-test-report.md
│   ├── prompt-requirements-understanding.md
│   └── ui-crud-enhancement-prompt.md
├── repo/                             # project code lives directly here
├── sessions/
│   ├── develop-1.json                # primary development session
│   └── bugfix-1.json                 # remediation session (if needed)
├── metadata.json
├── prompt.md
└── trajectory.json
```

### metadata.json

```json
{
  "prompt": "...",
  "project_type": "fullstack",
  "frontend_language": "go",
  "backend_language": "go",
  "frontend_framework": "templ",
  "backend_framework": "gin",
  "database": "postgresql"
}
```

---

## repo/ — Full Project Structure

```
repo/
├── cmd/
│   ├── hashgen/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   ├── model.go
│   │   ├── password.go
│   │   ├── lockout.go
│   │   └── session.go
│   │
│   ├── campaigns/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   └── model.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── db/
│   │   ├── postgres.go
│   │   ├── tx.go
│   │   └── migrate.go
│   │
│   ├── devices/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   ├── model.go
│   │   ├── ingest.go
│   │   └── dedupe.go
│   │
│   ├── exceptions/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   └── model.go
│   │
│   ├── jobs/
│   │   ├── worker.go
│   │   ├── scheduler.go
│   │   └── registry.go
│   │
│   ├── notifications/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   ├── model.go
│   │   ├── dispatcher.go
│   │   └── rules.go
│   │
│   ├── platform/
│   │   ├── logger/
│   │   ├── clock/
│   │   ├── security/
│   │   ├── pagination/
│   │   └── validator/
│   │
│   ├── reconciliation/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   ├── model.go
│   │   └── reconciliation.go
│   │
│   ├── segments/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   ├── model.go
│   │   └── runner.go
│   │
│   ├── server/
│   │   ├── app.go
│   │   ├── config.go
│   │   └── router.go
│   │
│   ├── tracking/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repo.go
│   │   ├── model.go
│   │   ├── smoother.go
│   │   └── stop_detector.go
│   │
│   └── web/
│       ├── handlers/
│       ├── middleware/
│       ├── templates/
│       │   ├── layouts/
│       │   │   ├── base.templ
│       │   │   └── auth.templ
│       │   ├── pages/
│       │   │   ├── login.templ
│       │   │   ├── dashboard.templ
│       │   │   ├── reservations.templ
│       │   │   ├── capacity.templ
│       │   │   ├── notifications.templ
│       │   │   ├── campaigns.templ
│       │   │   ├── segments.templ
│       │   │   ├── analytics.templ
│       │   │   ├── devices.templ
│       │   │   ├── audit.templ
│       │   │   └── admin/
│       │   │       ├── users.templ
│       │   │       └── content-rules.templ
│       │   ├── partials/
│       │   │   ├── activity-feed.templ
│       │   │   ├── conflict-warning.templ
│       │   │   ├── zone-card.templ
│       │   │   └── exception-list.templ
│       │   └── components/
│       │       ├── button.templ
│       │       ├── modal.templ
│       │       ├── table.templ
│       │       └── alert.templ
│       └── static/
│           ├── css/
│           │   └── app.css
│           ├── js/
│           │   └── poll.js
│           └── img/
│
├── migrations/
│   ├── 000001_initial_schema.up.sql
│   ├── 000002_seed_admin.up.sql
│   ├── 000003_master_data.up.sql
│   ├── 000004_reservations_capacity.up.sql
│   ├── 000005_device_integration.up.sql
│   ├── 000006_device_applied_sequence.up.sql
│   ├── 000007_exceptions.up.sql
│   ├── 000008_tracking.up.sql
│   ├── 000009_reconciliation.up.sql
│   ├── 000010_notifications.up.sql
│   ├── 000011_reconciliation_compat.up.sql
│   ├── 000012_campaigns_tasks.up.sql
│   ├── 000013_tagging_segmentation.up.sql
│   ├── 000014_analytics_exports.up.sql
│   └── 000015_seed_demo_data.up.sql
│
├── unit_tests/
│   ├── auth_test.go
│   ├── capacity_test.go
│   ├── device_test.go
│   ├── exception_test.go
│   ├── notifications_test.go
│   ├── rbac_test.go
│   ├── reconciliation_test.go
│   ├── security_test.go
│   └── tracking_test.go
│
├── API_tests/
│   ├── analytics_api_test.go
│   ├── auth_api_test.go
│   ├── campaigns_api_test.go
│   ├── devices_api_test.go
│   ├── exceptions_api_test.go
│   ├── master_data_api_test.go
│   ├── notifications_api_test.go
│   ├── rbac_api_test.go
│   ├── reconciliation_api_test.go
│   ├── reservations_api_test.go
│   ├── router_api_test.go
│   ├── segments_api_test.go
│   └── tracking_api_test.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── third_party/
│   └── templ/
│
├── run_tests.sh
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.test
├── .env.example
├── .dockerignore
├── go.mod
├── go.sum
└── README.md
```

---

## What Must NOT Be in the ZIP

- no `vendor/` directory
- no compiled binaries
- no `.env` with real credentials (only `.env.example`)
- no temp or scratch files

---

## Sessions Naming Rules

- primary development session → `sessions/develop-1.json`
- remediation session → `sessions/bugfix-1.json`
- additional sessions → `develop-2.json`, `bugfix-2.json`, etc.

---

## Submission Checklist

- [ ] `docker compose up` completes without errors
- [ ] Cold start tested in clean environment
- [ ] README URLs, ports, and credentials match running app
- [ ] `docs/design.md` and `docs/api-spec.md` present
- [ ] `docs/questions.md` has question + assumption + solution for each item
- [ ] `unit_tests/` and `API_tests/` exist in `repo/`, `run_tests.sh` passes
- [ ] No `vendor/`, cache, or compiled output in ZIP
- [ ] No real credentials in any config file
- [ ] All prompt requirements implemented — no silent substitutions
- [ ] `sessions/develop-1.json` trajectory file present
- [ ] `metadata.json` at root with all required fields
- [ ] `prompt.md` at root, unmodified
- [ ] Running application screenshots captured
- [ ] Self-test report generated and attached
