#!/usr/bin/env bash
# bootstrap/board-html.sh — render board.html: light kanban dashboard
# (Harness-style). Sidebar nav + main pane w/ 5-column board. Data from board.sh.
# Cards link to Aone, priority chip + pool tag, data-pool for filter. No assets.
# 任务池 col = pool-state candidates (cap ~80 visible + "+N 更多"). 5 equal 300px
# cols overflow-x scroll. Pool filter = dropdown w/ checkboxes (names from pools.json),
# default all checked, hide/show by data-pool, live counts. data-pool=key, label=name.
# Writes both .my-day/board.html (ephemeral) and docs/board.html (repo-tracked).
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="${JARVIS_ROOT:-$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))}"
out="$root/.my-day/board.html"
docs_out="$root/docs/board.html"
mkdir -p "$root/.my-day" "$root/docs"
json_f="$(mktemp)"; trap 'rm -f "$json_f"' EXIT
"$script_dir/board.sh" > "$json_f"

python3 - "$out" "$docs_out" "$json_f" "$root/config/pools.json" <<PY
import json, sys, html, datetime
out=sys.argv[1]; docs_out=sys.argv[2]
data=json.load(open(sys.argv[3]))
def e(s): return html.escape(str(s or ""))
PRI={"紧急":("#d92d20","#fef3f2"),"高":("#d92d20","#fef3f2"),"中":("#b54708","#fffaeb"),"低":("#475467","#f2f4f7")}
COLS=[("任务池","pool"),("待开始","__pre"),("进行中","inflight"),("审核中","done"),("已完成","merged")]
POOLS=["tf_provider","tf_customer","mcp_server","api_toolkit","cloudspec"]
# key→hue: accent border + tag chip tint + filter swatch share one color per pool
COLOR={"tf_provider":"#3b82f6","tf_customer":"#f59e0b","mcp_server":"#a855f7","cloudspec":"#22c55e","api_toolkit":"#ec4899"}
cfg=json.load(open(sys.argv[4])); NAMES={k:v.get("name",k) for k,v in cfg.get("pools",{}).items()}
def pname(k): return NAMES.get(k,k)
def pcol(k): return COLOR.get(k,"#98a2b3")  # neutral for unknown/empty
def tint(hexc): return hexc+"22"  # light bg from hue (~13% alpha), dark text stays the hue
CAP=80
def card(r):
  fg,bg=PRI.get(r.get("priority"),("#475467","#f2f4f7")); p=e(r.get("priority"))
  pk=r.get("pool"); ac=pcol(pk)
  sm="" if r["summary"]==r["title"] else f'<div class="sum">{e(r["summary"])}</div>'
  tag=f'<span class="tag" style="color:{ac};background:{tint(ac)}">{e(pname(pk))}</span>' if pk else ""
  return f'''<a class="card" data-pool="{e(pk)}" style="border-left:4px solid {ac}" href="{e(r['url'])}" target="_blank">
<div class="cr"><span class="cid">#{e(r['id'])}</span><span class="pri" style="color:{fg};background:{bg}">{p or '·'}</span></div>
<div class="tt">{e(r['title'])}</div>{sm}<div class="cf">{tag}</div></a>'''
def col(label,st):
  rows=[x for x in data if x["state"]==st]
  vis=rows[:CAP]; more=len(rows)-len(vis)
  body="".join(card(r) for r in vis) if vis else '<div class="empty">No items</div>'
  foot=f'<div class="more">+{more} 更多</div>' if more>0 else ''
  cls=" col-empty" if not rows else ""
  return f'''<div class="col{cls}" data-col><div class="ch"><span>{label}</span><span class="badge">{len(vis)}</span></div><div class="cb">{body}{foot}</div></div>'''
board="".join(col(l,s) for l,s in COLS)
arun=len([x for x in data if x["state"]=="inflight"])
rows="".join(f'<label class="dr"><input type=checkbox class=pf data-pf="{p}" checked><span class="sw" style="background:{pcol(p)}"></span><span>{e(pname(p))}</span></label>' for p in POOLS)
def nav(t,a=""): return f'<div class="nv{a}">{t}</div>'
gen=datetime.datetime.now().strftime("%Y-%m-%d %H:%M")
doc=f'''<!doctype html><meta charset=utf-8><title>Jarvis 工作板</title><style>
*{{box-sizing:border-box}}body{{margin:0;font:14px/1.5 -apple-system,Segoe UI,Roboto,sans-serif;color:#1d2939;background:#fff}}
.app{{display:flex;min-height:100vh}}
.sb{{width:210px;flex:none;background:#fafafa;border-right:1px solid #eaecf0;display:flex;flex-direction:column;padding:14px 10px}}
.brand{{display:flex;align-items:center;gap:8px;font-weight:700;font-size:16px;padding:6px 8px 14px}}
.dot{{width:10px;height:10px;border-radius:50%;background:#7f56d9}}
.grp{{font-size:11px;color:#98a2b3;text-transform:uppercase;letter-spacing:.4px;margin:14px 8px 4px}}
.nv{{padding:6px 10px;border-radius:7px;color:#475467;cursor:default;font-size:13.5px}}.nv:hover{{background:#f2f4f7}}.act{{background:#eef0f3;color:#1d2939;font-weight:600}}
.sf{{margin-top:auto;display:flex;align-items:center;gap:8px;padding:8px;border-top:1px solid #eaecf0}}
.av{{width:24px;height:24px;border-radius:50%;background:#7f56d9;color:#fff;font-size:11px;display:flex;align-items:center;justify-content:center}}
.ico{{margin-left:auto;color:#98a2b3}}.main{{flex:1;display:flex;flex-direction:column;min-width:0}}
.tb{{display:flex;align-items:center;padding:14px 22px;border-bottom:1px solid #eaecf0}}.bc{{color:#667085;font-size:13.5px}}.bc b{{color:#1d2939}}
.tb .r{{margin-left:auto;display:flex;gap:8px}}.btn{{border:1px solid #d0d5dd;background:#fff;border-radius:8px;padding:6px 12px;font-size:13px;color:#344054;cursor:pointer}}
.btn.k{{background:#1d2939;color:#fff;border-color:#1d2939}}
.fl{{display:flex;align-items:center;gap:8px;padding:12px 22px}}
.dd{{position:relative}}.ddb{{border:1px solid #d0d5dd;background:#fff;border-radius:8px;padding:6px 12px;font-size:13px;color:#344054;cursor:pointer}}
.ddp{{position:absolute;top:36px;left:0;z-index:5;background:#fff;border:1px solid #eaecf0;border-radius:10px;box-shadow:0 4px 16px rgba(16,24,40,.12);padding:6px;min-width:200px;display:none}}.ddp.open{{display:block}}
.dr{{display:flex;align-items:center;gap:8px;padding:6px 8px;border-radius:7px;font-size:13px;cursor:pointer}}.dr:hover{{background:#f2f4f7}}.dr input{{cursor:pointer}}.ddh{{border-bottom:1px solid #eaecf0;margin:2px 0 4px;font-weight:600}}.sw{{width:10px;height:10px;border-radius:3px;flex:none}}
.gen{{margin-left:auto;color:#98a2b3;font-size:12px}}
.bd{{display:flex;gap:14px;padding:8px 22px 26px;align-items:start;overflow-x:auto}}
.col{{flex:0 0 300px;background:#f9fafb;border:1px solid #eaecf0;border-radius:10px;padding:8px}}.col-empty{{background:#fff;border-style:dashed}}
.ch{{display:flex;justify-content:space-between;padding:6px 6px 10px;font-weight:600;font-size:13px;color:#344054}}.badge{{background:#eaecf0;color:#475467;border-radius:10px;font-size:11px;padding:0 7px}}
.cb{{display:flex;flex-direction:column;gap:8px;min-height:60px}}.empty{{color:#98a2b3;text-align:center;padding:24px 0;font-size:13px}}.more{{color:#98a2b3;text-align:center;font-size:12px;padding:6px 0}}
.card{{display:block;background:#fff;border:1px solid #eaecf0;border-radius:10px;padding:11px;text-decoration:none;color:inherit;box-shadow:0 1px 2px rgba(16,24,40,.04);transition:.15s}}
.card:hover{{box-shadow:0 4px 12px rgba(16,24,40,.1);transform:translateY(-2px)}}.card.hide{{display:none}}
.cr{{display:flex;justify-content:space-between;margin-bottom:5px}}.cid{{font:11px ui-monospace,monospace;color:#98a2b3}}.pri{{font-size:11px;font-weight:600;padding:0 8px;border-radius:10px}}
.tt{{font-weight:600;font-size:13px;color:#1d2939;margin-bottom:4px}}.sum{{color:#667085;font-size:12px}}.cf{{margin-top:8px}}.tag{{font-size:11px;color:#667085;background:#f2f4f7;border-radius:6px;padding:1px 7px}}
</style><div class="app"><aside class="sb"><div class="brand"><span class="dot"></span>Jarvis</div>
{nav("全部")}{nav("Manager")}{nav("收件箱")}{nav("自动化")}<div class="grp">Workspace</div>
{nav("工作板"," act")}{nav("Agents")}{nav("Skills")}{nav("知识·记忆")}<div class="grp">管理</div>
{nav("Workspace管理")}{nav("应用管理")}<div class="sf"><span class="av">辰</span>辰羿<span class="ico">⚙</span></div></aside>
<main class="main"><div class="tb"><div class="bc">Workspace › <b>工作板</b></div><div class="r"><button class="btn">刷新</button><button class="btn k">+ 新增任务</button></div></div>
<div class="fl"><div class="dd"><button class="ddb" id=ddb>工作池 ▾</button><div class="ddp" id=ddp><label class="dr ddh"><input type=checkbox id=pfall checked><span>全选</span></label>{rows}</div></div><span class="gen">{gen} · agent runs {arun}</span></div>
<div class="bd">{board}</div></main></div>
<script>
var B=document.querySelectorAll('.pf'),A=document.getElementById('pfall');
function sync(){{var on=new Set();B.forEach(c=>{{if(c.checked)on.add(c.dataset.pf)}});
document.querySelectorAll('.card').forEach(c=>{{c.classList.toggle('hide',c.dataset.pool!==''&&!on.has(c.dataset.pool))}});
document.querySelectorAll('[data-col]').forEach(col=>{{col.querySelector('.badge').textContent=col.querySelectorAll('.card:not(.hide)').length}});
A.checked=[...B].every(x=>x.checked);}}
B.forEach(c=>c.onchange=sync);A.onchange=function(){{B.forEach(x=>x.checked=A.checked);sync();}};
document.getElementById('ddb').onclick=function(e){{e.stopPropagation();document.getElementById('ddp').classList.toggle('open');}};
document.onclick=function(){{document.getElementById('ddp').classList.remove('open');}};
document.getElementById('ddp').onclick=function(e){{e.stopPropagation();}};
sync();
</script>'''
open(out,"w",encoding="utf-8").write(doc); print(out)
open(docs_out,"w",encoding="utf-8").write(doc); print(docs_out)
PY
