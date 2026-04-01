package server

import "net/http"

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Homestead — Stockyard</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>*{margin:0;padding:0;box-sizing:border-box}body{background:#1a1410;color:#f0e6d3;font-family:'JetBrains Mono',monospace;min-height:100vh;padding:2rem}
.hdr{font-size:.7rem;color:#a0845c;letter-spacing:3px;text-transform:uppercase;margin-bottom:2rem;border-bottom:2px solid #8b3d1a;padding-bottom:.8rem}
.grid{display:grid;grid-template-columns:1fr 1fr;gap:2rem}
.section{margin-bottom:2rem}.section h2{font-size:.65rem;letter-spacing:3px;text-transform:uppercase;color:#e8753a;margin-bottom:.8rem;border-bottom:1px solid #2e261e;padding-bottom:.4rem}
.item{background:#241e18;padding:.6rem .8rem;margin-bottom:.4rem;border:1px solid #2e261e;font-size:.72rem}
.item a{color:#e8753a;text-decoration:none}.item a:hover{color:#d4a843}
.item .cat{font-size:.55rem;color:#7a7060;margin-top:.2rem}
.note{background:#241e18;padding:.8rem;margin-bottom:.6rem;border:1px solid #2e261e;font-size:.72rem;color:#bfb5a3;white-space:pre-wrap;line-height:1.5}
.note-title{color:#f0e6d3;font-weight:600;margin-bottom:.3rem}
.empty{color:#7a7060;font-style:italic;padding:1rem;text-align:center}
</style></head><body>
<div class="hdr">Stockyard · Homestead</div>
<div class="grid">
<div>
<div class="section"><h2>Bookmarks</h2><div id="bm-list"></div></div>
<div class="section"><h2>Feeds</h2><div id="feed-list"></div></div>
</div>
<div>
<div class="section"><h2>Notes</h2><div id="note-list"></div></div>
</div>
</div>
<script>
async function refresh(){
  try{const d=await(await fetch('/api/bookmarks')).json();const bms=d.bookmarks||[];
  document.getElementById('bm-list').innerHTML=bms.length?bms.map(b=>'<div class="item"><a href="'+esc(b.url)+'" target="_blank">'+esc(b.title)+'</a>'+(b.category?'<div class="cat">'+esc(b.category)+'</div>':'')+'</div>').join(''):'<div class="empty">No bookmarks</div>';}catch(e){}
  try{const d=await(await fetch('/api/notes')).json();const ns=d.notes||[];
  document.getElementById('note-list').innerHTML=ns.length?ns.map(n=>'<div class="note">'+(n.title?'<div class="note-title">'+esc(n.title)+'</div>':'')+esc(n.content)+'</div>').join(''):'<div class="empty">No notes</div>';}catch(e){}
  try{const d=await(await fetch('/api/feeds')).json();const fs=d.feeds||[];
  document.getElementById('feed-list').innerHTML=fs.length?fs.map(f=>'<div class="item"><a href="'+esc(f.url)+'" target="_blank">'+esc(f.title)+'</a></div>').join(''):'<div class="empty">No feeds</div>';}catch(e){}
}
function esc(s){const d=document.createElement('div');d.textContent=s||'';return d.innerHTML;}
refresh();setInterval(refresh,10000);
</script></body></html>`))
}
