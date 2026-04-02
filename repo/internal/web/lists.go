package web

import (
	"context"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/a-h/templ"
)

func ListPage(user CurrentUser, title, currentPath, endpoint string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		clean := strings.ReplaceAll(currentPath, "/", "")
		clean = strings.ReplaceAll(clean, "-", "")
		if clean == "" {
			clean = "list"
		}

		content := fmt.Sprintf(`<section class="space-y-4">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-semibold">%s</h1>
      <p class="mt-1 text-slate-600">Operational data view.</p>
    </div>
    <button id="%sRefresh" class="rounded border border-slate-300 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 shadow-sm">Refresh</button>
  </div>
  <section class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
    <div id="%sState" class="mb-3 rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600">Loading...</div>
    <div class="overflow-x-auto">
      <table class="min-w-full text-left text-sm" id="%sTable">
        <thead id="%sHead" class="border-b border-slate-200 text-xs uppercase tracking-wider text-slate-500"></thead>
        <tbody id="%sBody" class="divide-y divide-slate-100"></tbody>
      </table>
    </div>
  </section>
</section>

<script>
(function() {
  const state = document.getElementById('%sState');
  const head = document.getElementById('%sHead');
  const body = document.getElementById('%sBody');
  const refreshBtn = document.getElementById('%sRefresh');

  async function loadData() {
    state.textContent = 'Loading...';
    state.className = 'mb-3 rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600';
    head.innerHTML = '';
    body.innerHTML = '';
    try {
      const res = await fetch('%s', { credentials: 'same-origin' });
      const payload = await res.json();
      if (!res.ok) {
        state.textContent = payload.message || 'Failed to load data';
        state.className = 'mb-3 rounded border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700';
        if(window.parkopsToast) window.parkopsToast('Failed to load data','error');
        return;
      }

      const rows = Array.isArray(payload) ? payload : (Array.isArray(payload.items) ? payload.items : []);
      if (rows.length === 0) {
        state.textContent = 'No records found.';
        return;
      }

      const columns = Object.keys(rows[0]);
      head.innerHTML = '<tr>' + columns.map(function(c) {
        return '<th class="px-3 py-2 font-semibold">' + c.replaceAll('_', ' ') + '</th>';
      }).join('') + '</tr>';

      body.innerHTML = rows.map(function(row) {
        const cells = columns.map(function(c) {
          const raw = row[c];
          let val = '';
          if (Array.isArray(raw)) val = raw.join(', ');
          else if (raw === null || raw === undefined) val = '';
          else if (typeof raw === 'object') val = JSON.stringify(raw);
          else val = String(raw);

          // Truncate long UUIDs for readability
          if (c === 'id' || c.endsWith('_id')) {
            val = val.length > 8 ? val.slice(0,8) + '...' : val;
          }
          // Format dates
          if (c.endsWith('_at') && val && val.includes('T')) {
            try { val = new Date(val).toLocaleString(); } catch(e) {}
          }
          return '<td class="px-3 py-2 text-slate-700 max-w-xs truncate">' + val + '</td>';
        }).join('');
        return '<tr class="hover:bg-slate-50">' + cells + '</tr>';
      }).join('');

      state.textContent = 'Loaded ' + rows.length + ' records.';
      if(window.parkopsToast) window.parkopsToast('Loaded ' + rows.length + ' records','success');
    } catch (_err) {
      state.textContent = 'Unable to load data';
      state.className = 'mb-3 rounded border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700';
      if(window.parkopsToast) window.parkopsToast('Failed to load data','error');
    }
  }

  refreshBtn.addEventListener('click', loadData);
  loadData();
})();
</script>`,
			html.EscapeString(title),
			html.EscapeString(clean),
			html.EscapeString(clean), html.EscapeString(clean), html.EscapeString(clean), html.EscapeString(clean),
			html.EscapeString(clean), html.EscapeString(clean), html.EscapeString(clean), html.EscapeString(clean),
			html.EscapeString(endpoint),
		)
		return AppLayout(user, title, currentPath, content).Render(ctx, w)
	})
}
