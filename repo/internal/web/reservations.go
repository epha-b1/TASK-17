package web

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func ReservationsPage(user CurrentUser) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		content := `<section class="space-y-6">
    <header>
      <h1 class="text-3xl font-semibold">Reservation Calendar</h1>
      <p class="text-slate-600 mt-2">Manage reservations, check availability, and create new holds.</p>
    </header>

    <!-- Date range filter for reservation list -->
    <section class="rounded-xl bg-white p-5 shadow-sm border border-slate-200 space-y-4">
      <h2 class="text-lg font-semibold">Reservations</h2>
      <div class="flex flex-wrap gap-3 items-end">
        <label class="text-sm">
          <span class="mb-1 block text-slate-600">Status</span>
          <select id="filterStatus" class="rounded border border-slate-300 px-3 py-2 text-sm">
            <option value="">All</option>
            <option value="hold">Hold</option>
            <option value="confirmed" selected>Confirmed</option>
            <option value="cancelled">Cancelled</option>
            <option value="expired">Expired</option>
          </select>
        </label>
        <button id="filterBtn" class="rounded bg-emerald-700 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-800">Load</button>
      </div>
      <div id="resState" class="rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600">Click Load to fetch reservations.</div>
      <div class="overflow-x-auto">
        <table class="min-w-full text-left text-sm">
          <thead class="border-b border-slate-200 text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="px-3 py-2">ID</th>
              <th class="px-3 py-2">Zone</th>
              <th class="px-3 py-2">Status</th>
              <th class="px-3 py-2">Stalls</th>
              <th class="px-3 py-2">Window Start</th>
              <th class="px-3 py-2">Window End</th>
              <th class="px-3 py-2">Hold Expires</th>
            </tr>
          </thead>
          <tbody id="resBody" class="divide-y divide-slate-100"></tbody>
        </table>
      </div>
    </section>

    <!-- Availability checker with date pickers -->
    <section class="rounded-xl bg-white p-5 shadow-sm border border-slate-200 space-y-4">
      <h2 class="text-lg font-semibold">Check Availability</h2>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <label class="text-sm">
          <span class="mb-1 block text-slate-600">Zone</span>
          <select id="zoneSelect" class="w-full rounded border border-slate-300 px-3 py-2">
            <option value="">Loading zones...</option>
          </select>
        </label>
        <label class="text-sm">
          <span class="mb-1 block text-slate-600">Start</span>
          <input id="startAt" type="datetime-local" class="w-full rounded border border-slate-300 px-3 py-2" />
        </label>
        <label class="text-sm">
          <span class="mb-1 block text-slate-600">End</span>
          <input id="endAt" type="datetime-local" class="w-full rounded border border-slate-300 px-3 py-2" />
        </label>
        <label class="text-sm">
          <span class="mb-1 block text-slate-600">Stalls needed</span>
          <input id="stallCount" type="number" min="1" value="1" class="w-full rounded border border-slate-300 px-3 py-2" />
        </label>
      </div>
      <button id="checkBtn" class="rounded bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800">Check Availability</button>
      <div id="result" class="rounded border border-slate-200 bg-slate-50 p-3 text-sm text-slate-700">Select a zone and time window, then click Check.</div>
      <div id="warning" class="hidden rounded border border-amber-300 bg-amber-50 p-3 text-sm font-medium text-amber-800">
        Conflict warning: not enough stalls available for this request.
      </div>
    </section>

  <script>
    // Set default dates
    const now = new Date();
    const later = new Date(now.getTime() + 4*3600*1000);
    document.getElementById('startAt').value = now.toISOString().slice(0,16);
    document.getElementById('endAt').value = later.toISOString().slice(0,16);

    // Load zones into dropdown
    (async function loadZones() {
      try {
        const res = await fetch('/api/capacity/dashboard', { credentials: 'same-origin' });
        const data = await res.json();
        const zones = Array.isArray(data.zones) ? data.zones : [];
        const sel = document.getElementById('zoneSelect');
        sel.innerHTML = zones.length === 0
          ? '<option value="">No zones</option>'
          : zones.map(z => '<option value="'+z.zone_id+'">'+z.zone_name+' ('+z.available_stalls+'/'+z.total_stalls+' avail)</option>').join('');
      } catch(e) { console.error(e); }
    })();

    // Load reservations
    document.getElementById('filterBtn').addEventListener('click', async () => {
      const status = document.getElementById('filterStatus').value;
      let url = '/api/reservations?limit=50';
      if (status) url += '&status=' + status;
      const state = document.getElementById('resState');
      const body = document.getElementById('resBody');
      try {
        const res = await fetch(url, { credentials: 'same-origin' });
        const payload = await res.json();
        if (!res.ok) { state.textContent = payload.message || 'Failed'; return; }
        const rows = Array.isArray(payload.items) ? payload.items : (Array.isArray(payload) ? payload : []);
        if (rows.length === 0) { state.textContent = 'No reservations found.'; body.innerHTML=''; return; }
        body.innerHTML = rows.map(r => {
          const statusClass = r.status==='confirmed'?'text-emerald-700':r.status==='hold'?'text-amber-700':r.status==='cancelled'?'text-rose-600':'text-slate-500';
          return '<tr>'
            +'<td class="px-3 py-2 font-mono text-xs">'+r.id.slice(0,8)+'</td>'
            +'<td class="px-3 py-2">'+(r.zone_id||'').slice(0,8)+'</td>'
            +'<td class="px-3 py-2 font-semibold '+statusClass+'">'+r.status+'</td>'
            +'<td class="px-3 py-2">'+r.stall_count+'</td>'
            +'<td class="px-3 py-2 text-xs">'+fmt(r.time_window_start)+'</td>'
            +'<td class="px-3 py-2 text-xs">'+fmt(r.time_window_end)+'</td>'
            +'<td class="px-3 py-2 text-xs">'+(r.hold_expires_at?fmt(r.hold_expires_at):'-')+'</td>'
            +'</tr>';
        }).join('');
        state.textContent = 'Showing '+rows.length+' reservations.';
        if(window.parkopsToast) window.parkopsToast('Loaded '+rows.length+' reservations','success');
      } catch(e) { state.textContent = 'Error loading reservations'; }
    });

    function fmt(iso) { if(!iso) return '-'; try { return new Date(iso).toLocaleString(); } catch(e) { return iso; } }

    // Availability check
    document.getElementById('checkBtn').addEventListener('click', async () => {
      const warning = document.getElementById('warning');
      const result = document.getElementById('result');
      warning.classList.add('hidden');
      const zoneId = document.getElementById('zoneSelect').value;
      const startVal = document.getElementById('startAt').value;
      const endVal = document.getElementById('endAt').value;
      const requested = parseInt(document.getElementById('stallCount').value || '1', 10);
      if (!zoneId || !startVal || !endVal) { result.textContent = 'Please fill in zone and dates.'; return; }
      const startISO = new Date(startVal).toISOString();
      const endISO = new Date(endVal).toISOString();
      const url = '/api/availability?zone_id='+encodeURIComponent(zoneId)+'&time_window_start='+encodeURIComponent(startISO)+'&time_window_end='+encodeURIComponent(endISO);
      try {
        const res = await fetch(url, { credentials: 'same-origin' });
        const payload = await res.json();
        if (!res.ok) { result.textContent = payload.message || 'Failed'; return; }
        result.textContent = 'Available: '+payload.available_stalls+' / '+payload.total_stalls+' stalls';
        if (Number(payload.available_stalls) < requested) { warning.classList.remove('hidden'); }
        else if(window.parkopsToast) window.parkopsToast('Availability check passed','success');
      } catch(e) { result.textContent = 'Error checking availability'; }
    });

    // Auto-load confirmed reservations on page load
    document.getElementById('filterBtn').click();
  </script>
</section>`
		return AppLayout(user, "Reservations", "/reservations", content).Render(ctx, w)
	})
}
