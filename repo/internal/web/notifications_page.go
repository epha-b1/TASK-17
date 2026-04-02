package web

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func NotificationsPage(user CurrentUser) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		content := `<section class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-semibold">Notifications</h1>
      <p class="mt-1 text-slate-600">View and manage your notifications.</p>
    </div>
    <button id="refreshNotifs" class="rounded border border-slate-300 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 shadow-sm">Refresh</button>
  </div>

  <div class="flex gap-2 text-sm">
    <button data-filter="all" class="nf-tab rounded-full border border-emerald-600 bg-emerald-700 px-3 py-1 text-white">All</button>
    <button data-filter="unread" class="nf-tab rounded-full border border-slate-300 px-3 py-1 text-slate-600 hover:bg-slate-100">Unread</button>
    <button data-filter="read" class="nf-tab rounded-full border border-slate-300 px-3 py-1 text-slate-600 hover:bg-slate-100">Read</button>
  </div>

  <div id="notifState" class="rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600">Loading...</div>
  <div id="notifList" class="space-y-2"></div>
</section>

<script>
(function(){
  var currentFilter = 'all';
  var allNotifs = [];

  async function loadNotifs() {
    document.getElementById('notifState').textContent = 'Loading...';
    document.getElementById('notifList').innerHTML = '';
    try {
      var res = await fetch('/api/notifications', {credentials:'same-origin'});
      var data = await res.json();
      if (!res.ok) { document.getElementById('notifState').textContent = data.message||'Failed'; return; }
      allNotifs = Array.isArray(data) ? data : (Array.isArray(data.items) ? data.items : []);
      renderNotifs();
    } catch(e) { document.getElementById('notifState').textContent = 'Error loading notifications'; }
  }

  function renderNotifs() {
    var filtered = allNotifs;
    if (currentFilter === 'unread') filtered = allNotifs.filter(function(n){return !n.read;});
    if (currentFilter === 'read') filtered = allNotifs.filter(function(n){return n.read;});

    if (filtered.length === 0) {
      document.getElementById('notifState').textContent = 'No notifications.';
      document.getElementById('notifList').innerHTML = '';
      return;
    }
    document.getElementById('notifState').textContent = filtered.length + ' notification(s)';

    document.getElementById('notifList').innerHTML = filtered.map(function(n) {
      var readClass = n.read ? 'border-slate-200 bg-white' : 'border-emerald-200 bg-emerald-50';
      var badge = n.read ? '<span class="text-xs text-slate-400">Read</span>' : '<span class="text-xs font-medium text-emerald-700 bg-emerald-100 px-2 py-0.5 rounded-full">New</span>';
      var dismissedBadge = n.dismissed ? ' <span class="text-xs text-slate-400">(Dismissed)</span>' : '';
      var date = '';
      try { date = new Date(n.created_at).toLocaleString(); } catch(e) { date = n.created_at||''; }

      var actions = '';
      if (!n.read) actions += '<button data-mark-read="'+n.id+'" class="rounded border border-slate-300 px-2 py-1 text-xs text-slate-700 hover:bg-slate-100">Mark Read</button> ';
      if (!n.dismissed) actions += '<button data-dismiss="'+n.id+'" class="rounded border border-rose-300 px-2 py-1 text-xs text-rose-600 hover:bg-rose-50">Dismiss</button>';

      return '<div class="rounded-lg border '+readClass+' p-4 shadow-sm">'
        +'<div class="flex items-start justify-between">'
        +'<div><p class="font-medium text-slate-900">'+esc(n.title||'Notification')+'</p>'
        +'<p class="mt-1 text-sm text-slate-600">'+esc(n.body||'')+'</p>'
        +'<p class="mt-1 text-xs text-slate-400">'+date+dismissedBadge+'</p></div>'
        +'<div class="flex items-center gap-2">'+badge+'</div>'
        +'</div>'
        +'<div class="mt-2 flex gap-2">'+actions+'</div>'
        +'</div>';
    }).join('');

    // Bind actions
    document.querySelectorAll('[data-mark-read]').forEach(function(btn){
      btn.onclick = async function(){
        await fetch('/api/notifications/'+btn.getAttribute('data-mark-read')+'/read',{method:'PATCH',credentials:'same-origin'});
        if(window.parkopsToast) window.parkopsToast('Marked as read','success');
        loadNotifs();
      };
    });
    document.querySelectorAll('[data-dismiss]').forEach(function(btn){
      btn.onclick = async function(){
        await fetch('/api/notifications/'+btn.getAttribute('data-dismiss')+'/dismiss',{method:'POST',credentials:'same-origin'});
        if(window.parkopsToast) window.parkopsToast('Dismissed','success');
        loadNotifs();
      };
    });
  }

  function esc(s) { var d=document.createElement('div'); d.textContent=s; return d.innerHTML; }

  // Tab switching
  document.querySelectorAll('.nf-tab').forEach(function(tab){
    tab.onclick = function(){
      currentFilter = tab.getAttribute('data-filter');
      document.querySelectorAll('.nf-tab').forEach(function(t){
        t.className = 'nf-tab rounded-full border border-slate-300 px-3 py-1 text-slate-600 hover:bg-slate-100';
      });
      tab.className = 'nf-tab rounded-full border border-emerald-600 bg-emerald-700 px-3 py-1 text-white';
      renderNotifs();
    };
  });

  document.getElementById('refreshNotifs').onclick = loadNotifs;
  loadNotifs();
})();
</script>`
		return AppLayout(user, "Notifications", "/notifications", content).Render(ctx, w)
	})
}
