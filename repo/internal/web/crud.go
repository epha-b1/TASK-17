package web

import (
	"context"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// CrudField describes a form field for create/edit modals.
type CrudField struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Type        string   `json:"type"` // text, number, select, hidden, textarea, uuid, lookup
	Required    bool     `json:"required"`
	Placeholder string   `json:"placeholder,omitempty"`
	Options     []Option `json:"options,omitempty"` // for select type
	Default     string   `json:"default,omitempty"`
	ReadOnly    bool     `json:"readOnly,omitempty"` // show in table but not editable
	LookupAPI   string   `json:"lookupAPI,omitempty"`   // e.g. "/api/facilities" — fetches options dynamically
	LookupLabel string   `json:"lookupLabel,omitempty"` // which field to show as label (e.g. "name")
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// CrudPageConfig holds everything needed to render a full CRUD page.
type CrudPageConfig struct {
	Title      string
	Path       string
	APIBase    string      // e.g. /api/facilities
	Fields     []CrudField // fields for create/edit
	CanCreate  bool
	CanEdit    bool
	CanDelete  bool
	IDField    string // which JSON field is the ID (default "id")
	ExtraJS    string // optional extra JS to inject
}

func CrudPage(user CurrentUser, cfg CrudPageConfig) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if cfg.IDField == "" {
			cfg.IDField = "id"
		}

		// Build fields JSON for JS
		fieldsJSON := "["
		for i, f := range cfg.Fields {
			if i > 0 {
				fieldsJSON += ","
			}
			optsJSON := "[]"
			if len(f.Options) > 0 {
				var opts []string
				for _, o := range f.Options {
					opts = append(opts, fmt.Sprintf(`{"value":%q,"label":%q}`, o.Value, o.Label))
				}
				optsJSON = "[" + strings.Join(opts, ",") + "]"
			}
			fieldsJSON += fmt.Sprintf(`{"key":%q,"label":%q,"type":%q,"required":%t,"placeholder":%q,"options":%s,"default":%q,"readOnly":%t,"lookupAPI":%q,"lookupLabel":%q}`,
				f.Key, f.Label, f.Type, f.Required, f.Placeholder, optsJSON, f.Default, f.ReadOnly, f.LookupAPI, f.LookupLabel)
		}
		fieldsJSON += "]"

		content := fmt.Sprintf(`<section class="space-y-4">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-semibold">%s</h1>
      <p class="mt-1 text-slate-600">Manage %s records.</p>
    </div>
    <div class="flex gap-2">
      <button id="refreshBtn" class="rounded border border-slate-300 bg-white px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 shadow-sm">Refresh</button>
      %s
    </div>
  </div>

  <section class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
    <div id="listState" class="mb-3 rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600">Loading...</div>
    <div class="overflow-x-auto">
      <table class="min-w-full text-left text-sm">
        <thead id="listHead" class="border-b border-slate-200 text-xs uppercase tracking-wider text-slate-500"></thead>
        <tbody id="listBody" class="divide-y divide-slate-100"></tbody>
      </table>
    </div>
  </section>

  <!-- Modal backdrop -->
  <div id="modalBg" class="fixed inset-0 z-40 hidden bg-black/40"></div>

  <!-- Create/Edit modal -->
  <div id="formModal" class="fixed inset-0 z-50 hidden items-center justify-center p-4" style="display:none">
    <div class="mx-auto w-full max-w-lg rounded-xl bg-white p-6 shadow-2xl ring-1 ring-slate-200">
      <div class="flex items-center justify-between mb-4">
        <h2 id="modalTitle" class="text-xl font-semibold">Create</h2>
        <button id="modalClose" class="text-slate-400 hover:text-slate-700 text-xl font-bold">&times;</button>
      </div>
      <form id="crudForm" class="space-y-4">
        <div id="formFields"></div>
        <div id="formError" class="hidden rounded border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700"></div>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" id="cancelBtn" class="rounded border border-slate-300 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50">Cancel</button>
          <button type="submit" id="submitBtn" class="rounded bg-emerald-700 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-800">Save</button>
        </div>
      </form>
    </div>
  </div>

  <!-- Delete confirm modal -->
  <div id="deleteModal" class="fixed inset-0 z-50 hidden items-center justify-center p-4" style="display:none">
    <div class="mx-auto w-full max-w-sm rounded-xl bg-white p-6 shadow-2xl ring-1 ring-slate-200">
      <h2 class="text-lg font-semibold text-slate-900">Confirm Delete</h2>
      <p class="mt-2 text-sm text-slate-600">Are you sure you want to delete this record? This cannot be undone.</p>
      <div class="mt-4 flex justify-end gap-2">
        <button id="delCancelBtn" class="rounded border border-slate-300 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50">Cancel</button>
        <button id="delConfirmBtn" class="rounded bg-rose-600 px-4 py-2 text-sm font-medium text-white hover:bg-rose-700">Delete</button>
      </div>
    </div>
  </div>
</section>

<script>
(function(){
  const API = %q;
  const ID_FIELD = %q;
  const fields = %s;
  const canCreate = %t;
  const canEdit = %t;
  const canDelete = %t;

  const listState = document.getElementById('listState');
  const listHead = document.getElementById('listHead');
  const listBody = document.getElementById('listBody');
  const modalBg = document.getElementById('modalBg');
  const formModal = document.getElementById('formModal');
  const modalTitle = document.getElementById('modalTitle');
  const formFields = document.getElementById('formFields');
  const formError = document.getElementById('formError');
  const crudForm = document.getElementById('crudForm');
  const deleteModal = document.getElementById('deleteModal');

  let editingId = null;
  let deletingId = null;

  // ── Modal helpers ──
  function openModal() { modalBg.classList.remove('hidden'); formModal.style.display='flex'; formModal.classList.remove('hidden'); formError.classList.add('hidden'); }
  function closeModal() { modalBg.classList.add('hidden'); formModal.style.display='none'; formModal.classList.add('hidden'); editingId=null; }
  function openDeleteModal(id) { deletingId=id; modalBg.classList.remove('hidden'); deleteModal.style.display='flex'; deleteModal.classList.remove('hidden'); }
  function closeDeleteModal() { modalBg.classList.add('hidden'); deleteModal.style.display='none'; deleteModal.classList.add('hidden'); deletingId=null; }

  document.getElementById('modalClose').onclick = closeModal;
  document.getElementById('cancelBtn').onclick = closeModal;
  document.getElementById('delCancelBtn').onclick = closeDeleteModal;
  modalBg.onclick = function(){ closeModal(); closeDeleteModal(); };

  // ── Build form fields ──
  function buildForm(data) {
    formFields.innerHTML = '';
    fields.forEach(function(f) {
      if (f.readOnly) return;
      var wrap = document.createElement('label');
      wrap.className = 'block text-sm';
      var span = document.createElement('span');
      span.className = 'mb-1 block font-medium text-slate-700';
      span.textContent = f.label + (f.required ? ' *' : '');
      wrap.appendChild(span);

      var input;
      if (f.type === 'select') {
        input = document.createElement('select');
        input.className = 'w-full rounded border border-slate-300 px-3 py-2';
        input.name = f.key;
        if (!f.required) {
          var empty = document.createElement('option');
          empty.value = ''; empty.textContent = '-- select --';
          input.appendChild(empty);
        }
        (f.options||[]).forEach(function(o) {
          var opt = document.createElement('option');
          opt.value = o.value; opt.textContent = o.label;
          input.appendChild(opt);
        });
      } else if (f.type === 'uuid') {
        // If lookupAPI is set, render as a searchable dropdown
        if (f.lookupAPI) {
          input = document.createElement('select');
          input.className = 'w-full rounded border border-slate-300 px-3 py-2';
          input.name = f.key;
          var loadingOpt = document.createElement('option');
          loadingOpt.value = ''; loadingOpt.textContent = 'Loading...';
          input.appendChild(loadingOpt);
          if (f.required) input.required = true;
          // Load options from API
          (function(sel, api, labelKey, currentVal) {
            fetch(api, {credentials:'same-origin'}).then(function(r){return r.json();}).then(function(payload){
              var rows = Array.isArray(payload) ? payload : (Array.isArray(payload.items) ? payload.items : []);
              sel.innerHTML = '';
              if (!f.required) { var e = document.createElement('option'); e.value=''; e.textContent='-- select --'; sel.appendChild(e); }
              rows.forEach(function(row) {
                var o = document.createElement('option');
                o.value = row.id || '';
                var lbl = row[labelKey] || row.name || row.display_name || row.title || row.plate_number || row.id || '';
                // Show extra info if available
                var extra = '';
                if (row.address) extra = ' — ' + row.address;
                else if (row.display_name && labelKey !== 'display_name') extra = ' — ' + row.display_name;
                o.textContent = lbl + extra + ' (' + (row.id||'').slice(0,8) + '...)';
                if (currentVal && row.id === currentVal) o.selected = true;
                sel.appendChild(o);
              });
              if (rows.length === 0) { var e = document.createElement('option'); e.value=''; e.textContent='No records found'; sel.appendChild(e); }
            }).catch(function(){
              sel.innerHTML = '<option value="">Failed to load</option>';
            });
          })(input, f.lookupAPI, f.lookupLabel || 'name', data ? data[f.key] : '');
          wrap.appendChild(span);
          wrap.appendChild(input);
          formFields.appendChild(wrap);
          return;
        }
        // Fallback: manual UUID input with auto-format
        input = document.createElement('input');
        input.type = 'text';
        input.className = 'w-full rounded border border-slate-300 px-3 py-2 font-mono text-sm';
        input.name = f.key;
        input.placeholder = f.placeholder || 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
        input.maxLength = 36;
        input.pattern = '[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}';
        input.title = 'UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx';
        input.addEventListener('input', function() {
          var raw = this.value.replace(/[^0-9a-fA-F]/g, '').toLowerCase().slice(0,32);
          var parts = [];
          if (raw.length > 0) parts.push(raw.slice(0,Math.min(8,raw.length)));
          if (raw.length > 8) parts.push(raw.slice(8,Math.min(12,raw.length)));
          if (raw.length > 12) parts.push(raw.slice(12,Math.min(16,raw.length)));
          if (raw.length > 16) parts.push(raw.slice(16,Math.min(20,raw.length)));
          if (raw.length > 20) parts.push(raw.slice(20,Math.min(32,raw.length)));
          this.value = parts.join('-');
          var valid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/.test(this.value);
          this.style.borderColor = this.value.length === 0 ? '' : valid ? '#059669' : '#dc2626';
        });
        var hint = document.createElement('p');
        hint.className = 'mt-1 text-xs text-slate-400';
        hint.textContent = 'Paste or type a UUID — hyphens added automatically';
        wrap.appendChild(span);
        wrap.appendChild(input);
        wrap.appendChild(hint);
        if (f.required) input.required = true;
        var val = data ? data[f.key] : (f['default'] || '');
        if (val !== undefined && val !== null) input.value = String(val);
        formFields.appendChild(wrap);
        return;
      } else if (f.type === 'textarea') {
        input = document.createElement('textarea');
        input.className = 'w-full rounded border border-slate-300 px-3 py-2';
        input.name = f.key;
        input.rows = 3;
        input.placeholder = f.placeholder || '';
      } else if (f.type === 'hidden') {
        input = document.createElement('input');
        input.type = 'hidden'; input.name = f.key;
      } else {
        input = document.createElement('input');
        input.type = f.type || 'text';
        input.className = 'w-full rounded border border-slate-300 px-3 py-2';
        input.name = f.key;
        input.placeholder = f.placeholder || '';
        if (f.type === 'number') { input.min = '0'; input.step = '1'; }
      }
      if (f.required) input.required = true;

      // Set value
      var val = data ? data[f.key] : (f['default'] || '');
      if (val !== undefined && val !== null) {
        if (typeof val === 'object') val = JSON.stringify(val);
        input.value = String(val);
      }

      wrap.appendChild(input);
      if (f.type !== 'hidden') formFields.appendChild(wrap);
      else formFields.appendChild(input);
    });
  }

  // ── Create button ──
  var createBtn = document.getElementById('createBtn');
  if (createBtn) {
    createBtn.onclick = function() {
      editingId = null;
      modalTitle.textContent = 'Create';
      buildForm(null);
      openModal();
    };
  }

  // ── Form submit ──
  crudForm.onsubmit = async function(e) {
    e.preventDefault();
    formError.classList.add('hidden');
    var uuidRe = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    var body = {};
    var validationErr = '';
    fields.forEach(function(f) {
      if (f.readOnly) return;
      var el = crudForm.querySelector('[name="'+f.key+'"]');
      if (!el) return;
      var v = el.value.trim();
      // UUID validation
      if (f.type === 'uuid' && v !== '') {
        if (!uuidRe.test(v)) { validationErr = f.label + ' must be a valid UUID (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)'; return; }
      }
      // Required check
      if (f.required && v === '') { validationErr = f.label + ' is required'; return; }
      if (f.type === 'number' && v !== '') body[f.key] = Number(v);
      else if (v !== '') body[f.key] = v;
    });

    if (validationErr) {
      formError.textContent = validationErr;
      formError.classList.remove('hidden');
      return;
    }

    var method = editingId ? 'PATCH' : 'POST';
    var url = editingId ? API + '/' + editingId : API;

    try {
      var res = await fetch(url, { method: method, credentials: 'same-origin', headers: {'Content-Type':'application/json'}, body: JSON.stringify(body) });
      if (res.status === 204 || res.status === 200 || res.status === 201) {
        closeModal();
        loadData();
        if(window.parkopsToast) window.parkopsToast(editingId?'Updated successfully':'Created successfully','success');
      } else {
        var err = await res.json().catch(function(){return{}});
        formError.textContent = err.message || 'Operation failed ('+res.status+')';
        formError.classList.remove('hidden');
      }
    } catch(ex) {
      formError.textContent = 'Network error'; formError.classList.remove('hidden');
    }
  };

  // ── Delete confirm ──
  document.getElementById('delConfirmBtn').onclick = async function() {
    if (!deletingId) return;
    try {
      var res = await fetch(API+'/'+deletingId, { method:'DELETE', credentials:'same-origin' });
      closeDeleteModal();
      if (res.ok || res.status === 204) {
        loadData();
        if(window.parkopsToast) window.parkopsToast('Deleted successfully','success');
      } else {
        if(window.parkopsToast) window.parkopsToast('Delete failed','error');
      }
    } catch(ex) {
      closeDeleteModal();
      if(window.parkopsToast) window.parkopsToast('Delete failed','error');
    }
  };

  // ── Load data ──
  async function loadData() {
    listState.textContent = 'Loading...';
    listState.className = 'mb-3 rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600';
    listHead.innerHTML = ''; listBody.innerHTML = '';
    try {
      var res = await fetch(API, { credentials:'same-origin' });
      var payload = await res.json();
      if (!res.ok) { listState.textContent = payload.message||'Failed'; listState.className='mb-3 rounded border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700'; return; }
      var rows = Array.isArray(payload) ? payload : (Array.isArray(payload.items) ? payload.items : []);
      if (rows.length === 0) { listState.textContent = 'No records found.'; return; }

      var columns = Object.keys(rows[0]);
      var hasActions = canEdit || canDelete;
      listHead.innerHTML = '<tr>' + columns.map(function(c){return '<th class="px-3 py-2 font-semibold">'+c.replaceAll('_',' ')+'</th>';}).join('') + (hasActions?'<th class="px-3 py-2 font-semibold">Actions</th>':'') + '</tr>';

      listBody.innerHTML = rows.map(function(row) {
        var rid = row[ID_FIELD] || '';
        var cells = columns.map(function(c) {
          var raw = row[c];
          var val = '';
          if (Array.isArray(raw)) val = raw.join(', ');
          else if (raw===null||raw===undefined) val='';
          else if (typeof raw==='object') val=JSON.stringify(raw);
          else val=String(raw);
          if ((c==='id'||c.endsWith('_id'))&&val.length>8) val=val.slice(0,8)+'...';
          if (c.endsWith('_at')&&val&&val.includes('T')) try{val=new Date(val).toLocaleString();}catch(e){}
          return '<td class="px-3 py-2 text-slate-700 max-w-xs truncate">'+val+'</td>';
        }).join('');

        var actions = '';
        if (hasActions) {
          actions = '<td class="px-3 py-2 flex gap-1">';
          if (canEdit) actions += '<button data-edit="'+rid+'" class="rounded border border-slate-300 px-2 py-1 text-xs text-slate-700 hover:bg-slate-100">Edit</button>';
          if (canDelete) actions += '<button data-del="'+rid+'" class="rounded border border-rose-300 px-2 py-1 text-xs text-rose-600 hover:bg-rose-50">Delete</button>';
          actions += '</td>';
        }
        return '<tr class="hover:bg-slate-50" data-row=\''+JSON.stringify(row).replace(/'/g,"&#39;")+'\'>'+cells+actions+'</tr>';
      }).join('');

      listState.textContent = 'Loaded '+rows.length+' records.';
      if(window.parkopsToast) window.parkopsToast('Loaded '+rows.length+' records','success');

      // Bind edit/delete buttons
      listBody.querySelectorAll('[data-edit]').forEach(function(btn){
        btn.onclick = function(){
          var id = btn.getAttribute('data-edit');
          var rowData = JSON.parse(btn.closest('tr').getAttribute('data-row'));
          editingId = id;
          modalTitle.textContent = 'Edit';
          buildForm(rowData);
          openModal();
        };
      });
      listBody.querySelectorAll('[data-del]').forEach(function(btn){
        btn.onclick = function(){ openDeleteModal(btn.getAttribute('data-del')); };
      });
    } catch(_e) {
      listState.textContent = 'Unable to load data';
      listState.className = 'mb-3 rounded border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700';
    }
  }

  document.getElementById('refreshBtn').onclick = loadData;
  loadData();
  %s
})();
</script>`,
			html.EscapeString(cfg.Title),
			strings.ToLower(html.EscapeString(cfg.Title)),
			func() string {
				if cfg.CanCreate {
					return `<button id="createBtn" class="rounded bg-emerald-700 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-800 shadow-sm">+ New</button>`
				}
				return ""
			}(),
			cfg.APIBase,
			cfg.IDField,
			fieldsJSON,
			cfg.CanCreate,
			cfg.CanEdit,
			cfg.CanDelete,
			cfg.ExtraJS,
		)

		return AppLayout(user, cfg.Title, cfg.Path, content).Render(ctx, w)
	})
}
