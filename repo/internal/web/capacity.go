package web

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func CapacityPage(user CurrentUser) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		content := `<section class="space-y-6">
  <div>
    <h1 class="text-3xl font-semibold">Capacity</h1>
    <p class="mt-1 text-slate-600">Zone-level capacity overview with real-time occupancy bars.</p>
  </div>

  <!-- Zone cards with progress bars -->
  <section id="zoneCards" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
    <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
      <p class="text-sm text-slate-500">Loading capacity...</p>
    </article>
  </section>

  <!-- Detailed table -->
  <section class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
    <h2 class="text-lg font-semibold mb-3">Detailed Breakdown</h2>
    <div class="overflow-x-auto">
      <table class="min-w-full text-left text-sm">
        <thead class="border-b border-slate-200 text-xs uppercase tracking-wider text-slate-500">
          <tr>
            <th class="px-3 py-2">Zone</th>
            <th class="px-3 py-2">Total</th>
            <th class="px-3 py-2">Held</th>
            <th class="px-3 py-2">Confirmed</th>
            <th class="px-3 py-2">Available</th>
            <th class="px-3 py-2">Occupancy</th>
          </tr>
        </thead>
        <tbody id="capacityBody" class="divide-y divide-slate-100"></tbody>
      </table>
    </div>
    <p id="capacityState" class="mt-3 text-xs text-slate-500"></p>
  </section>

  <!-- Snapshots history -->
  <section class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
    <h2 class="text-lg font-semibold mb-3">Recent Capacity Snapshots</h2>
    <div id="snapshotState" class="text-sm text-slate-600 mb-2">Loading...</div>
    <div class="overflow-x-auto">
      <table class="min-w-full text-left text-sm">
        <thead class="border-b border-slate-200 text-xs uppercase tracking-wider text-slate-500">
          <tr><th class="px-3 py-2">Zone</th><th class="px-3 py-2">Snapshot At</th><th class="px-3 py-2">Authoritative Stalls</th></tr>
        </thead>
        <tbody id="snapshotBody" class="divide-y divide-slate-100"></tbody>
      </table>
    </div>
  </section>
</section>

<script>
(async function() {
  const cards = document.getElementById('zoneCards');
  const body = document.getElementById('capacityBody');
  const state = document.getElementById('capacityState');
  try {
    const res = await fetch('/api/capacity/dashboard', { credentials: 'same-origin' });
    const payload = await res.json();
    if (!res.ok) { state.textContent = payload.message || 'Failed'; return; }
    const zones = Array.isArray(payload.zones) ? payload.zones : [];
    if (zones.length === 0) { cards.innerHTML='<p class="text-slate-500 p-4">No zones configured.</p>'; return; }

    // Cards with progress bars
    cards.innerHTML = zones.map(z => {
      const used = z.total_stalls - z.available_stalls;
      const pct = z.total_stalls > 0 ? Math.round(used/z.total_stalls*100) : 0;
      const barColor = pct >= 90 ? 'bg-rose-500' : pct >= 70 ? 'bg-amber-500' : 'bg-emerald-500';
      const borderColor = pct >= 90 ? 'border-rose-200' : pct >= 70 ? 'border-amber-200' : 'border-slate-200';
      return '<article class="rounded-xl border '+borderColor+' bg-white p-5 shadow-sm">'
        +'<div class="flex items-center justify-between">'
        +'<p class="font-semibold text-slate-900">'+z.zone_name+'</p>'
        +'<span class="text-xs font-medium px-2 py-0.5 rounded-full '+(pct>=90?'bg-rose-100 text-rose-700':pct>=70?'bg-amber-100 text-amber-700':'bg-emerald-100 text-emerald-700')+'">'+pct+'%</span>'
        +'</div>'
        +'<p class="mt-2 text-2xl font-bold text-slate-900">'+z.available_stalls+'<span class="text-sm font-normal text-slate-500"> / '+z.total_stalls+'</span></p>'
        +'<p class="text-xs text-slate-500">available stalls</p>'
        +'<div class="mt-3 h-2 w-full rounded-full bg-slate-100">'
        +'<div class="h-2 rounded-full '+barColor+'" style="width:'+pct+'%"></div>'
        +'</div>'
        +'<div class="mt-2 flex justify-between text-xs text-slate-500">'
        +'<span>Held: '+z.held_stalls+'</span>'
        +'<span>Confirmed: '+z.confirmed_stalls+'</span>'
        +'</div>'
        +'</article>';
    }).join('');

    // Table
    body.innerHTML = zones.map(z => {
      const used = z.total_stalls - z.available_stalls;
      const pct = z.total_stalls > 0 ? Math.round(used/z.total_stalls*100) : 0;
      return '<tr>'
        +'<td class="px-3 py-2 font-medium">'+z.zone_name+'</td>'
        +'<td class="px-3 py-2">'+z.total_stalls+'</td>'
        +'<td class="px-3 py-2">'+z.held_stalls+'</td>'
        +'<td class="px-3 py-2">'+z.confirmed_stalls+'</td>'
        +'<td class="px-3 py-2 font-semibold">'+z.available_stalls+'</td>'
        +'<td class="px-3 py-2"><span class="font-semibold '+(pct>=90?'text-rose-600':pct>=70?'text-amber-600':'text-emerald-600')+'">'+pct+'%</span></td>'
        +'</tr>';
    }).join('');
    state.textContent = 'Loaded '+zones.length+' zones.';
    if(window.parkopsToast) window.parkopsToast('Capacity loaded','success');
  } catch(_e) { state.textContent='Error loading capacity'; }
})();

// Snapshots
(async function() {
  const sState = document.getElementById('snapshotState');
  const sBody = document.getElementById('snapshotBody');
  try {
    const res = await fetch('/api/capacity/snapshots?limit=20', { credentials: 'same-origin' });
    const data = await res.json();
    if (!res.ok) { sState.textContent='Failed to load snapshots'; return; }
    const rows = Array.isArray(data) ? data : (Array.isArray(data.items) ? data.items : []);
    if (rows.length===0) { sState.textContent='No snapshots yet.'; return; }
    sBody.innerHTML = rows.map(s => '<tr>'
      +'<td class="px-3 py-2 font-mono text-xs">'+(s.zone_id||'').slice(0,8)+'</td>'
      +'<td class="px-3 py-2 text-xs">'+new Date(s.snapshot_at).toLocaleString()+'</td>'
      +'<td class="px-3 py-2">'+s.authoritative_stalls+'</td>'
      +'</tr>').join('');
    sState.textContent = rows.length+' snapshots loaded.';
  } catch(_e) { sState.textContent='Error loading snapshots'; }
})();
</script>`
		return AppLayout(user, "Capacity", "/capacity", content).Render(ctx, w)
	})
}
