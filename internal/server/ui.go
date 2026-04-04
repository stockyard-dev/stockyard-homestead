package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Homestead</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:1000px;margin:0 auto}
table{width:100%;border-collapse:collapse;font-size:.72rem}
th{text-align:left;padding:.5rem .6rem;border-bottom:2px solid var(--bg3);font-size:.6rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px}
td{padding:.4rem .6rem;border-bottom:1px solid var(--bg3)}
tr:hover{background:var(--bg2)}
.badge-active{color:var(--green)}.badge-inactive{color:var(--red)}.badge-maintenance{color:var(--gold)}
.tag{font-size:.5rem;padding:.05rem .25rem;background:var(--bg3);color:var(--cm);margin-right:.2rem}
.btn{font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>HOMESTEAD</h1><button class="btn btn-p" onclick="openForm()">+ Add Server</button></div>
<div class="main" id="main"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let servers=[];
async function load(){const r=await fetch(A+'/servers').then(r=>r.json());servers=r.servers||[];render();}
function render(){if(!servers.length){document.getElementById('main').innerHTML='<div class="empty">No servers registered. Add your first one.</div>';return;}
let h='<table><tr><th>Hostname</th><th>IP</th><th>OS</th><th>Provider</th><th>Region</th><th>Status</th><th></th></tr>';
servers.forEach(s=>{
h+='<tr><td style="color:var(--cream)">'+esc(s.hostname)+'</td><td>'+esc(s.ip)+'</td><td>'+esc(s.os)+'</td><td>'+esc(s.provider)+'</td><td>'+esc(s.region)+'</td><td><span class="badge-'+s.status+'">'+s.status+'</span></td><td><button class="btn" onclick="del(\''+s.id+'\')" style="font-size:.5rem;color:var(--cm)">✕</button></td></tr>';});
h+='</table>';document.getElementById('main').innerHTML=h;}
async function del(id){if(confirm('Remove?')){await fetch(A+'/servers/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>Add Server</h2><div class="fr"><label>Hostname</label><input id="f-h" placeholder="e.g. web-prod-01"></div><div class="fr"><label>IP Address</label><input id="f-ip" placeholder="10.0.1.5"></div><div class="fr"><label>OS</label><input id="f-os" placeholder="Ubuntu 24.04"></div><div class="fr"><label>Provider</label><input id="f-p" placeholder="AWS, Hetzner, Railway"></div><div class="fr"><label>Region</label><input id="f-r" placeholder="us-east-1"></div><div class="fr"><label>Tags</label><input id="f-t" placeholder="web, production"></div><div class="fr"><label>Status</label><select id="f-s"><option value="active">Active</option><option value="maintenance">Maintenance</option><option value="inactive">Inactive</option></select></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Add</button></div>';document.getElementById('mbg').classList.add('open');}
async function sub(){await fetch(A+'/servers',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({hostname:document.getElementById('f-h').value,ip:document.getElementById('f-ip').value,os:document.getElementById('f-os').value,provider:document.getElementById('f-p').value,region:document.getElementById('f-r').value,tags:document.getElementById('f-t').value,status:document.getElementById('f-s').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
