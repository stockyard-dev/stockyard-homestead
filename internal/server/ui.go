package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Homestead</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1200px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.green{color:var(--green)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.table{background:var(--bg2);border:1px solid var(--bg3);overflow-x:auto}
.table table{width:100%;border-collapse:collapse;font-size:.7rem}
.table th{text-align:left;padding:.6rem .7rem;color:var(--cm);text-transform:uppercase;font-size:.55rem;letter-spacing:1px;border-bottom:1px solid var(--bg3);background:var(--bg)}
.table td{padding:.6rem .7rem;border-bottom:1px solid var(--bg3);color:var(--cream);vertical-align:top}
.table tr:hover td{background:var(--bg3);cursor:pointer}
.table tr.inactive td,.table tr.decommissioned td{opacity:.55}
.col-host{font-weight:700}
.col-ip{font-family:var(--mono);color:var(--cd);font-size:.65rem}
.col-os{color:var(--cd)}
.col-prov{color:var(--leather)}
.tag{display:inline-block;font-size:.5rem;padding:.05rem .3rem;background:var(--bg3);color:var(--cd);font-family:var(--mono);margin-right:.2rem}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.active{border-color:var(--green);color:var(--green)}
.badge.inactive{border-color:var(--cm);color:var(--cm)}
.badge.maintenance{border-color:var(--orange);color:var(--orange)}
.badge.decommissioned{border-color:var(--red);color:var(--red)}

.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}

.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:560px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> HOMESTEAD</h1>
<button class="btn btn-p" onclick="openNew()">+ Add Host</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search hostname, ip, tags, notes..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="active">Active</option>
<option value="inactive">Inactive</option>
<option value="maintenance">Maintenance</option>
<option value="decommissioned">Decommissioned</option>
</select>
<select class="filter-sel" id="provider-filter" onchange="render()">
<option value="">All Providers</option>
</select>
<select class="filter-sel" id="region-filter" onchange="render()">
<option value="">All Regions</option>
</select>
</div>
<div class="table" id="table-wrap"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='hosts';

var fields=[
{name:'hostname',label:'Hostname',type:'text',required:true},
{name:'ip',label:'IP Address',type:'text'},
{name:'os',label:'Operating System',type:'select_or_text',options:['Ubuntu 22.04','Ubuntu 24.04','Debian 12','RHEL 9','Alpine','Windows Server 2022']},
{name:'provider',label:'Provider',type:'select_or_text',options:['aws','gcp','azure','digitalocean','hetzner','linode','on-prem']},
{name:'region',label:'Region',type:'text'},
{name:'tags',label:'Tags',type:'text',placeholder:'comma separated'},
{name:'status',label:'Status',type:'select',options:['active','inactive','maintenance','decommissioned']},
{name:'notes',label:'Notes',type:'textarea'}
];

var hosts=[],hostExtras={},editId=null,searchTimer=null;

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
return d.toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'});
}catch(e){return s}
}

function fieldByName(n){for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];return null}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/hosts').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
hosts=resps[0].hosts||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
hostExtras=ex||{};
hosts.forEach(function(h){
var x=hostExtras[h.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(h[k]===undefined)h[k]=x[k]});
});
}catch(e){hostExtras={}}

populateFilters();
}catch(e){
console.error('load failed',e);
hosts=[];
}
render();
}

function populateFilters(){
function fillSelect(elId,key){
var sel=document.getElementById(elId);
if(!sel)return;
var current=sel.value;
var seen={};var items=[];
hosts.forEach(function(h){if(h[key]&&!seen[h[key]]){seen[h[key]]=true;items.push(h[key])}});
items.sort();
var label=sel.firstElementChild?sel.firstElementChild.textContent:'';
sel.innerHTML='<option value="">'+label+'</option>'+items.map(function(x){return'<option value="'+esc(x)+'"'+(x===current?' selected':'')+'>'+esc(x)+'</option>'}).join('');
}
fillSelect('provider-filter','provider');
fillSelect('region-filter','region');
}

function renderStats(s){
var total=s.total||0;
var active=s.active||0;
var byProvider=s.by_provider||{};
var byRegion=s.by_region||{};
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Hosts</div></div>'+
'<div class="st"><div class="st-v green">'+active+'</div><div class="st-l">Active</div></div>'+
'<div class="st"><div class="st-v">'+Object.keys(byProvider).length+'</div><div class="st-l">Providers</div></div>'+
'<div class="st"><div class="st-v">'+Object.keys(byRegion).length+'</div><div class="st-l">Regions</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var pf=document.getElementById('provider-filter').value;
var rf=document.getElementById('region-filter').value;

var f=hosts.slice();
if(q)f=f.filter(function(h){
return(h.hostname||'').toLowerCase().includes(q)||
(h.ip||'').toLowerCase().includes(q)||
(h.tags||'').toLowerCase().includes(q)||
(h.notes||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(h){return h.status===sf});
if(pf)f=f.filter(function(h){return h.provider===pf});
if(rf)f=f.filter(function(h){return h.region===rf});

if(!f.length){
var msg=window._emptyMsg||'No hosts in the inventory yet.';
document.getElementById('table-wrap').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var customCols=fields.filter(function(fd){return fd.isCustom});

var h='<table><thead><tr>';
h+='<th>Hostname</th><th>IP</th><th>OS</th><th>Provider</th><th>Region</th><th>Tags</th><th>Status</th>';
customCols.forEach(function(fd){h+='<th>'+esc(fd.label)+'</th>'});
h+='</tr></thead><tbody>';

f.forEach(function(host){
var cls=host.status||'active';
h+='<tr class="'+esc(cls)+'" onclick="openEdit(\''+esc(host.id)+'\')">';
h+='<td class="col-host">'+esc(host.hostname)+'</td>';
h+='<td class="col-ip">'+esc(host.ip||'-')+'</td>';
h+='<td class="col-os">'+esc(host.os||'-')+'</td>';
h+='<td class="col-prov">'+esc(host.provider||'-')+'</td>';
h+='<td>'+esc(host.region||'-')+'</td>';
h+='<td>';
if(host.tags){
String(host.tags).split(',').forEach(function(t){
t=t.trim();
if(t)h+='<span class="tag">'+esc(t)+'</span>';
});
}else h+='-';
h+='</td>';
h+='<td><span class="badge '+esc(host.status||'active')+'">'+esc(host.status||'active')+'</span></td>';
customCols.forEach(function(fd){
var v=host[fd.name];
h+='<td>'+(v===undefined||v===null||v===''?'-':esc(String(v)))+'</td>';
});
h+='</tr>';
});
h+='</tbody></table>';

document.getElementById('table-wrap').innerHTML=h;
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
(f.options||[]).forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="3"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
h+='<input type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(host){
var h0=host||{};
var isEdit=!!host;
var h='<h2>'+(isEdit?'EDIT HOST':'NEW HOST')+'</h2>';

h+=fieldHTML(fieldByName('hostname'),h0.hostname);
h+='<div class="row2">'+fieldHTML(fieldByName('ip'),h0.ip)+fieldHTML(fieldByName('os'),h0.os)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('provider'),h0.provider)+fieldHTML(fieldByName('region'),h0.region)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('tags'),h0.tags)+fieldHTML(fieldByName('status'),h0.status||'active')+'</div>';
h+=fieldHTML(fieldByName('notes'),h0.notes);

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,h0[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit)h+='<button class="btn btn-del" onclick="delItem()">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button>';
h+='</div>';
return h;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-hostname');if(n)n.focus();
}

function openEdit(id){
var h=null;
for(var i=0;i<hosts.length;i++){if(hosts[i].id===id){h=hosts[i];break}}
if(!h)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(h);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var nameEl=document.getElementById('f-hostname');
if(!nameEl||!nameEl.value.trim()){alert('Hostname is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='number')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/hosts/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/hosts',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Add failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){alert('Network error: '+e.message);return}
closeModal();
load();
}

async function delItem(){
if(!editId)return;
if(!confirm('Delete this host?'))return;
await fetch(A+'/hosts/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.providers)){
var pf=fieldByName('provider');
if(pf)pf.options=cfg.providers;
}
if(Array.isArray(cfg.os_options)){
var of=fieldByName('os');
if(of)of.options=cfg.os_options;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
