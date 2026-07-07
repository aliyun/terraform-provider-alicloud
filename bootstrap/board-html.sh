#!/usr/bin/env bash
# bootstrap/board-html.sh — render board.html: light kanban dashboard
# (Harness-style). Sidebar nav + main pane w/ 5-column board. Data from board.sh.
# Cards link to Aone, priority chip + pool tag, data-pool for filter. No assets.
# 任务池 col = pool-state candidates (ALL rendered; col body scrolls vertically). 5 equal
# 300px cols overflow-x scroll. Pool filter = dropdown w/ checkboxes (names from pools.json),
# default all checked, hide/show by data-pool, live counts. data-pool=key, label=name.
# Writes docs/board.html (gitignored build artifact; served by serve.sh).
# --refresh: force scan.sh --force then rebuild (原 refresh.sh 折入)。
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
root="$(jarvis_root)"
docs_out="$root/docs/board.html"

force_scan=0
if [ "${1:-}" = "--refresh" ]; then
    force_scan=1
    echo "refresh: rescanning Aone (--force)…" >&2
    bash "$script_dir/scan.sh" --force >/dev/null || { echo "refresh: scan failed" >&2; exit 1; }
    echo "refresh: rebuilding board…" >&2
fi

mkdir -p "$root/docs"
json_f="$(mktemp)"; probe_f="$(mktemp)"; trap 'rm -f "$json_f" "$probe_f"' EXIT
"$script_dir/board.sh" > "$json_f"
# probe 段(飞轮健康度):board.sh probe 恒退 0(a1 失败自动降级),WARN 落 stderr,html 构建静默
"$script_dir/board.sh" probe > "$probe_f" 2>/dev/null || echo '{}' > "$probe_f"

python3 - "$docs_out" "$json_f" "$root/config/pools.json" "$probe_f" <<PY
import json, sys, html, datetime
docs_out=sys.argv[1]
data=json.load(open(sys.argv[2]))
try: prb=json.load(open(sys.argv[4]))
except Exception: prb={}
def e(s): return html.escape(str(s or ""))
PRI={"紧急":("#d92d20","#fef3f2"),"高":("#d92d20","#fef3f2"),"中":("#b54708","#fffaeb"),"低":("#475467","#f2f4f7")}
COLS=[("任务池","pool"),("待开始","__pre"),("进行中","inflight"),("审核中","done"),("已完成","merged")]
POOLS=["tf_provider","tf_customer","mcp_server","api_toolkit"]
# key→hue: accent border + tag chip tint + filter swatch share one color per pool
COLOR={"tf_provider":"#3b82f6","tf_customer":"#f59e0b","mcp_server":"#a855f7","api_toolkit":"#ec4899"}
CAT={"req":("需求","#0e7090","#ecfdff"),"bug":("缺陷","#b42318","#fef3f2"),"task":("任务","#5925dc","#f4f3ff")}
cfg=json.load(open(sys.argv[3])); NAMES={k:v.get("name",k) for k,v in cfg.get("pools",{}).items()}
PCNT={}; PSPLIT={}  # pool key → candidate count + {req,bug,task} split (full, even if cap clipped DOM)
for x in data:
  if x["state"]=="pool":
    PCNT[x.get("pool")]=x.get("pool_total") or PCNT.get(x.get("pool"),0)
    PSPLIT[x.get("pool")]={"req":x.get("pool_req",0),"bug":x.get("pool_bug",0),"task":x.get("pool_task",0)}
def pcnt(k): return PCNT.get(k,0)
def psplit(k): return PSPLIT.get(k,{"req":0,"bug":0,"task":0})
TOTAL=sum(PCNT.get(p,0) for p in POOLS)
def pname(k): return NAMES.get(k,k)
def pcol(k): return COLOR.get(k,"#98a2b3")  # neutral for unknown/empty
def tint(hexc): return hexc+"22"  # light bg from hue (~13% alpha), dark text stays the hue
def probe_strip(p):
  # 飞轮健康度 stat 带:board.sh probe 聚合结果 → 一排 tile。空/失败则不渲染。
  if not p: return ""
  f=p.get("findings") or {}; d=p.get("drafts") or {}; t=p.get("tickets") or {}
  sc=p.get("scenarios") or {}; cov=p.get("tier0_coverage") or {}
  def tile(k,v,s="",warn=False):
    cls="ptile warn" if warn else "ptile"
    sub=f'<div class="ps">{e(s)}</div>' if s else ""
    # str(v) first: 现有 e() 用 (s or "") 会把整数 0 渲染成空,数值 tile 须显式 str 保 "0" 可见
    return f'<div class="{cls}"><div class="pk">{e(k)}</div><div class="pv">{e(str(v))}</div>{sub}</div>'
  ar=d.get("adoption_rate"); ar_s=(str(round(ar*100))+"%") if ar is not None else "—"
  tot=t.get("total"); tot_s="—" if tot is None else str(tot); warn=(t.get("source")!="aone")
  bp=sc.get("by_product") or {}
  bptop=", ".join(f"{k} {v}" for k,v in sorted(bp.items(), key=lambda kv:-kv[1])[:3]) or "—"
  last=p.get("last_tier1_round") or ""; last=last[:10] if last else "—"
  ag=f.get("api_gap") or {}
  agtop=", ".join(f'{k.replace("api_gap_","")} {v}' for k,v in sorted(ag.items(), key=lambda kv:-kv[1])[:2]) or "无"
  tiles=[
    tile("本周发现", f.get("total",0), f'{f.get("rounds",0)} 轮 · tier0 {f.get("tier0_rounds",0)}/tier1 {f.get("tier1_rounds",0)}'),
    tile("采纳率", ar_s, f'filed {d.get("filed",0)}/pend {d.get("pending",0)}/rej {d.get("rejected",0)}'),
    tile("建单", tot_s, ("降级·本地" if warn else f'关单 {t.get("closed",0)}'), warn),
    tile("api_gap", sum(ag.values()) if ag else 0, agtop),
    tile("场景总数", sc.get("total",0), bptop),
    tile("tier0 覆盖", cov.get("resources_scanned",0), "已巡检资源"),
    tile("最近 tier1", last, f'近 {p.get("window_days",7)} 天窗口'),
  ]
  return '<div class="probe" title="probe 飞轮健康度 · board.sh probe">'+"".join(tiles)+'</div>'
def card(r):
  fg,bg=PRI.get(r.get("priority"),("#475467","#f2f4f7")); p=e(r.get("priority"))
  pk=r.get("pool"); ac=pcol(pk); cat=r.get("category") or ""
  sm="" if r["summary"]==r["title"] else f'<div class="sum">{e(r["summary"])}</div>'
  tag=f'<span class="tag" style="color:{ac};background:{tint(ac)}">{e(pname(pk))}</span>' if pk else ""
  cl,cfg2,cbg=CAT.get(cat,("","","")); badge=f'<span class="cat" style="color:{cfg2};background:{cbg}">{cl}</span>' if cl else ""
  return f'''<a class="card" data-pool="{e(pk)}" data-cat="{e(cat)}" data-title="{e(r['title']).lower()}" style="border-left:4px solid {ac}" href="{e(r['url'])}" target="_blank">
<div class="cr"><span class="cid">#{e(r['id'])}</span><span class="pri" style="color:{fg};background:{bg}">{p or '·'}</span></div>
<div class="tt">{e(r['title'])}</div>{sm}<div class="cf">{tag}{badge}</div></a>'''
def col(label,st):
  rows=[x for x in data if x["state"]==st]
  body="".join(card(r) for r in rows) if rows else '<div class="empty">No items</div>'
  cls=" col-empty" if not rows else ""
  return f'''<div class="col{cls}" data-col><div class="ch"><span>{label}</span><span class="badge">{len(rows)}</span></div><div class="cb">{body}</div></div>'''
board="".join(col(l,s) for l,s in COLS)
probe_html=probe_strip(prb)
arun=len([x for x in data if x["state"]=="inflight"])
def drow(p):
  s=psplit(p)
  return f'<label class="dr"><input type=checkbox class=pf data-pf="{p}" checked><span class="sw" style="background:{pcol(p)}"></span><span class="nm">{e(pname(p))}</span><span class="cnt">{pcnt(p)} (需{s["req"]}/缺{s["bug"]}/务{s["task"]})</span></label>'
rows="".join(drow(p) for p in POOLS)
CTOT={c:sum(psplit(p)[c] for p in POOLS) for c in ("req","bug","task")}
pills="".join(f'<button class="cp" data-cf="{c}" aria-pressed=true><span class="cd" style="background:{CAT[c][1]}"></span>{CAT[c][0]}<span class="ccnt">{CTOT[c]}</span></button>' for c in ("req","bug","task"))
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
.srch{{border:1px solid #d0d5dd;background:#fff;border-radius:8px;padding:6px 12px;font-size:13px;color:#344054;width:220px;outline:none}}.srch:focus{{border-color:#7f56d9;box-shadow:0 0 0 3px rgba(127,86,217,.15)}}
.fl{{display:flex;align-items:center;gap:8px;padding:12px 22px}}
.dd{{position:relative}}.ddb{{border:1px solid #d0d5dd;background:#fff;border-radius:8px;padding:6px 12px;font-size:13px;color:#344054;cursor:pointer}}
.ddp{{position:absolute;top:36px;left:0;z-index:5;background:#fff;border:1px solid #eaecf0;border-radius:10px;box-shadow:0 4px 16px rgba(16,24,40,.12);padding:6px;width:280px;display:none}}.ddp.open{{display:block}}
.dr{{display:flex;align-items:center;gap:8px;padding:6px 8px;border-radius:7px;font-size:13px;cursor:pointer;white-space:nowrap}}.dr:hover{{background:#f2f4f7}}.dr input{{cursor:pointer}}.dr .nm{{overflow:hidden;text-overflow:ellipsis}}.cnt{{margin-left:auto;color:#98a2b3;font:11px ui-monospace,monospace;flex:none}}.ddh{{border-bottom:1px solid #eaecf0;margin:2px 0 4px;font-weight:600}}.sw{{width:10px;height:10px;border-radius:3px;flex:none}}
.cp{{display:inline-flex;align-items:center;gap:5px;border:1px solid #d0d5dd;background:#fff;border-radius:16px;padding:4px 10px;font-size:12.5px;color:#344054;cursor:pointer}}.cp .ccnt{{color:#98a2b3;font:10px ui-monospace,monospace}}.cp[aria-pressed=false]{{opacity:.4;background:#f2f4f7}}.cd{{width:9px;height:9px;border-radius:50%;flex:none}}
.gen{{margin-left:auto;color:#98a2b3;font-size:12px}}
.probe{{display:flex;gap:10px;padding:2px 22px 12px;flex-wrap:wrap;align-items:stretch}}
.ptile{{flex:0 0 auto;min-width:118px;background:#f9fafb;border:1px solid #eaecf0;border-radius:10px;padding:9px 13px}}
.ptile .pk{{font-size:11px;color:#98a2b3;text-transform:uppercase;letter-spacing:.3px}}
.ptile .pv{{font-size:19px;font-weight:700;color:#1d2939;margin-top:2px}}
.ptile .ps{{font-size:11px;color:#667085;margin-top:2px}}
.ptile.warn{{border-color:#fec84b;background:#fffaeb}}
.bd{{display:flex;gap:14px;padding:8px 22px 26px;align-items:start;overflow-x:auto}}
.col{{flex:0 0 300px;background:#f9fafb;border:1px solid #eaecf0;border-radius:10px;padding:8px}}.col-empty{{background:#fff;border-style:dashed}}
.ch{{display:flex;justify-content:space-between;padding:6px 6px 10px;font-weight:600;font-size:13px;color:#344054}}.badge{{background:#eaecf0;color:#475467;border-radius:10px;font-size:11px;padding:0 7px}}
.cb{{display:flex;flex-direction:column;gap:8px;min-height:60px;max-height:calc(100vh - 220px);overflow-y:auto}}.empty{{color:#98a2b3;text-align:center;padding:24px 0;font-size:13px}}
.card{{display:block;background:#fff;border:1px solid #eaecf0;border-radius:10px;padding:11px;text-decoration:none;color:inherit;box-shadow:0 1px 2px rgba(16,24,40,.04);transition:.15s}}
.card:hover{{box-shadow:0 4px 12px rgba(16,24,40,.1);transform:translateY(-2px)}}.card.hide{{display:none}}
.cr{{display:flex;justify-content:space-between;margin-bottom:5px}}.cid{{font:11px ui-monospace,monospace;color:#98a2b3}}.pri{{font-size:11px;font-weight:600;padding:0 8px;border-radius:10px}}
.tt{{font-weight:600;font-size:13px;color:#1d2939;margin-bottom:4px}}.sum{{color:#667085;font-size:12px}}.cf{{margin-top:8px;display:flex;gap:6px;flex-wrap:wrap}}.tag{{font-size:11px;color:#667085;background:#f2f4f7;border-radius:6px;padding:1px 7px}}.cat{{font-size:11px;font-weight:600;border-radius:6px;padding:1px 7px}}
</style><div class="app"><aside class="sb"><div class="brand"><span class="dot"></span>Jarvis</div>
{nav("全部")}{nav("Manager")}{nav("收件箱")}{nav("自动化")}<div class="grp">Workspace</div>
{nav("工作板"," act")}{nav("Agents")}{nav("Skills")}{nav("知识·记忆")}<div class="grp">管理</div>
{nav("Workspace管理")}{nav("应用管理")}<div class="sf"><span class="av">辰</span>辰羿<span class="ico">⚙</span></div></aside>
<main class="main"><div class="tb"><div class="bc">Workspace › <b>工作板</b></div><div class="r"><input class="srch" id=srch type=search placeholder="搜索工作项…" autocomplete=off><button class="btn" id=refresh title="运行 bash bootstrap/board-html.sh --refresh 重扫 Aone 并重建">刷新</button><button class="btn k">+ 新增任务</button></div></div>
<div class="fl"><div class="dd"><button class="ddb" id=ddb>工作池 ▾</button><div class="ddp" id=ddp><label class="dr ddh"><input type=checkbox id=pfall checked><span class="nm">全选</span><span class="cnt">{TOTAL}</span></label>{rows}</div></div>{pills}<span class="gen">{gen} · agent runs {arun}</span></div>
{probe_html}<div class="bd">{board}</div></main></div>
<script>
var B=document.querySelectorAll('.pf'),A=document.getElementById('pfall'),P=document.querySelectorAll('.cp'),S=document.getElementById('srch');
function sync(){{var on=new Set();B.forEach(c=>{{if(c.checked)on.add(c.dataset.pf)}});
var ct=new Set();P.forEach(c=>{{if(c.getAttribute('aria-pressed')==='true')ct.add(c.dataset.cf)}});
var q=S.value.trim().toLowerCase();
document.querySelectorAll('.card').forEach(c=>{{var okP=c.dataset.pool===''||on.has(c.dataset.pool);var okC=c.dataset.cat===''||ct.has(c.dataset.cat);var okT=q===''||(c.dataset.title||'').indexOf(q)>=0;c.classList.toggle('hide',!(okP&&okC&&okT))}});
document.querySelectorAll('[data-col]').forEach(col=>{{col.querySelector('.badge').textContent=col.querySelectorAll('.card:not(.hide)').length}});
A.checked=[...B].every(x=>x.checked);}}
B.forEach(c=>c.onchange=sync);A.onchange=function(){{B.forEach(x=>x.checked=A.checked);sync();}};
P.forEach(c=>c.onclick=function(){{c.setAttribute('aria-pressed',c.getAttribute('aria-pressed')==='true'?'false':'true');sync();}});
var sd;S.oninput=function(){{clearTimeout(sd);sd=setTimeout(sync,150);}};
document.getElementById('ddb').onclick=function(e){{e.stopPropagation();document.getElementById('ddp').classList.toggle('open');}};
document.onclick=function(){{document.getElementById('ddp').classList.remove('open');}};
document.getElementById('ddp').onclick=function(e){{e.stopPropagation();}};
var R=document.getElementById('refresh');R.onclick=function(){{
if(location.protocol==='file:'){{alert('未通过服务运行,请跑 bootstrap/serve.sh 后从 http://localhost:8787 打开');return;}}
R.disabled=true;R.textContent='刷新中…';
fetch('/refresh',{{method:'POST'}}).then(r=>{{if(r.ok)location.reload();else throw 0;}})
.catch(()=>{{alert('刷新失败,请手动运行 bash bootstrap/board-html.sh --refresh');R.disabled=false;R.textContent='刷新';}});}};
sync();
</script>'''
open(docs_out,"w",encoding="utf-8").write(doc); print(docs_out)
PY

if [ "$force_scan" = "1" ]; then
    scan_f="$root/.my-day/scan.json"
    total=$(jq 'length' "$scan_f" 2>/dev/null || echo 0)
    dist=$(jq -r 'group_by(.category)|map("\(.[0].category // "—"):\(length)")|join(" ")' "$scan_f" 2>/dev/null)
    echo "refresh: done — $total items ($dist)" >&2
fi
