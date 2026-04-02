package web

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func LoginPage() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_ = ctx
		_, err := io.WriteString(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>ParkOps Login</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <style>
    .toast-container{position:fixed;top:1rem;right:1rem;z-index:9999;display:flex;flex-direction:column;gap:.5rem;pointer-events:none}
    .toast{pointer-events:auto;min-width:280px;max-width:380px;padding:.75rem 1rem;border-radius:.5rem;font-size:.875rem;font-weight:500;box-shadow:0 4px 12px rgba(0,0,0,.15);display:flex;align-items:center;gap:.5rem;animation:toast-in .3s ease-out,toast-out .3s ease-in 3.7s forwards}
    .toast-success{background:#065f46;color:#fff}
    .toast-error{background:#991b1b;color:#fff}
    .toast-info{background:#1e3a5f;color:#fff}
    @keyframes toast-in{from{opacity:0;transform:translateX(100%)}to{opacity:1;transform:translateX(0)}}
    @keyframes toast-out{from{opacity:1}to{opacity:0}}
  </style>
</head>
<body class="min-h-screen bg-slate-100 text-slate-900">
  <div id="toastContainer" class="toast-container"></div>
  <main class="mx-auto flex min-h-screen max-w-md items-center p-6">
    <section class="w-full rounded-xl bg-white p-8 shadow-sm ring-1 ring-slate-200">
      <p class="text-sm font-medium uppercase tracking-widest text-emerald-700">ParkOps</p>
      <h1 class="mt-2 text-3xl font-semibold">Sign in</h1>
      <p class="mt-1 text-sm text-slate-600">Local operations console</p>
      <form class="mt-8 space-y-4" method="post" action="/auth/login">
        <label class="block">
          <span class="mb-1 block text-sm font-medium">Username</span>
          <input class="w-full rounded-md border border-slate-300 px-3 py-2" type="text" name="username" autocomplete="username" required>
        </label>
        <label class="block">
          <span class="mb-1 block text-sm font-medium">Password</span>
          <input class="w-full rounded-md border border-slate-300 px-3 py-2" type="password" name="password" autocomplete="current-password" required>
        </label>
        <button class="w-full rounded-md bg-emerald-700 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-800" type="submit">Sign in</button>
      </form>
    </section>
  </main>
  <script>
  (function(){
    var msgs={login_success:['Signed in successfully','success'],login_error:['Invalid username or password','error'],logout_success:['You have been signed out','success'],session_expired:['Session expired, please sign in again','info'],password_changed:['Password changed successfully','success']};
    function showToast(text,type){
      var c=document.getElementById('toastContainer');if(!c)return;
      var el=document.createElement('div');el.className='toast toast-'+(type||'info');
      var icon=type==='success'?'\u2713':type==='error'?'\u2717':'\u2139';
      el.textContent=icon+' '+text;c.appendChild(el);
      setTimeout(function(){el.remove()},4200);
    }
    var p=new URLSearchParams(window.location.search);var t=p.get('toast');
    if(t&&msgs[t]){showToast(msgs[t][0],msgs[t][1]);var u=new URL(window.location);u.searchParams.delete('toast');window.history.replaceState({},'',u)}
  })();
  </script>
</body>
</html>`)
		return err
	})
}
