#!/usr/bin/env bash
# bootstrap/board-html.sh — render board.html: light kanban dashboard
# (Harness-style). Sidebar nav + main pane w/ 5-column board. Data from board.sh.
# Cards link to Aone, priority chip + pool tag. No external assets. JARVIS_ROOT ok.
# Writes both .my-day/board.html (ephemeral) and docs/board.html (repo-tracked).
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="${JARVIS_ROOT:-$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))}"
out="$root/.my-day/board.html"
docs_out="$root/docs/board.html"
mkdir -p "$root/.my-day" "$root/docs"
json="$("$script_dir/board.sh")"

python3 - "$out" "$docs_out" <<PY
import json, sys, html, datetime
out=sys.argv[1]; docs_out=sys.argv[2]
data=json.loads('''$json''')
def e(s): return html.escape(str(s or ""))
PRI={"高":("#d92d20","#fef3f2"),"中":("#b54708","#fffaeb"),"低":("#475467","#f2f4f7")}
COLS=[("任务池","escalated"),("待开始","__pre"),("进行中","inflight"),("审核中","done"),("已完成","merged")]
def card(r):
  fg,bg=PRI.get(r.get("priority"),("#475467","#f2f4f7")); p=e(r.get("priority"))
  sm="" if r["summary"]==r["title"] else f'<div class="sum">{e(r["summary"])}</div>'
  tag=f'<span class="tag">{e(r["pool"])}</span>' if r.get("pool") else ""
  return f'''<a class="card" href="{e(r['url'])}" target="_blank">
<div class="cr"><span class="cid">#{e(r['id'])}</span><span class="pri" style="color:{fg};background:{bg}">{p or '·'}</span></div>
<div class="tt">{e(r['title'])}</div>{sm}<div class="cf">{tag}</div></a>'''
def col(label,st):
  rows=[x for x in data if x["state"]==st]
  body="".join(card(r) for r in rows) if rows else '<div class="empty">No items</div>'
  cls=" col-empty" if not rows else ""
  return f'''<div class="col{cls}"><div class="ch"><span>{label}</span><span class="badge">{len(rows)}</span></div><div class="cb">{body}</div></div>'''
board="".join(col(l,s) for l,s in COLS)
arun=len([x for x in data if x["state"]=="inflight"])
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
.fl{{display:flex;align-items:center;gap:10px;padding:12px 22px}}.tog{{display:flex;border:1px solid #d0d5dd;border-radius:8px;overflow:hidden}}.tog span{{padding:5px 14px;font-size:13px;color:#475467}}.tog .on{{background:#f2f4f7;color:#1d2939;font-weight:600}}
.pill{{padding:5px 12px;border:1px solid #eaecf0;border-radius:16px;font-size:12.5px;color:#475467}}.pill.on{{background:#f2f4f7;color:#1d2939;border-color:#d0d5dd}}
.sr{{margin-left:auto;border:1px solid #d0d5dd;border-radius:8px;padding:6px 12px;font-size:13px;color:#98a2b3;min-width:200px}}
.bd{{display:grid;grid-template-columns:repeat(5,1fr);gap:14px;padding:8px 22px 26px;align-items:start}}
.col{{background:#f9fafb;border:1px solid #eaecf0;border-radius:10px;padding:8px}}.col-empty{{background:#fff;border-style:dashed}}
.ch{{display:flex;justify-content:space-between;padding:6px 6px 10px;font-weight:600;font-size:13px;color:#344054}}.badge{{background:#eaecf0;color:#475467;border-radius:10px;font-size:11px;padding:0 7px}}
.cb{{display:flex;flex-direction:column;gap:8px;min-height:60px}}.empty{{color:#98a2b3;text-align:center;padding:24px 0;font-size:13px}}
.card{{display:block;background:#fff;border:1px solid #eaecf0;border-radius:10px;padding:11px;text-decoration:none;color:inherit;box-shadow:0 1px 2px rgba(16,24,40,.04);transition:.15s}}
.card:hover{{box-shadow:0 4px 12px rgba(16,24,40,.1);transform:translateY(-2px)}}
.cr{{display:flex;justify-content:space-between;margin-bottom:5px}}.cid{{font:11px ui-monospace,monospace;color:#98a2b3}}.pri{{font-size:11px;font-weight:600;padding:0 8px;border-radius:10px}}
.tt{{font-weight:600;font-size:13px;color:#1d2939;margin-bottom:4px}}.sum{{color:#667085;font-size:12px}}.cf{{margin-top:8px}}.tag{{font-size:11px;color:#667085;background:#f2f4f7;border-radius:6px;padding:1px 7px}}
</style><div class="app"><aside class="sb"><div class="brand"><span class="dot"></span>Jarvis</div>
{nav("全部")}{nav("Manager")}{nav("收件箱")}{nav("自动化")}<div class="grp">Workspace</div>
{nav("工作板"," act")}{nav("Agents")}{nav("Skills")}{nav("知识·记忆")}<div class="grp">管理</div>
{nav("Workspace管理")}{nav("应用管理")}<div class="sf"><span class="av">辰</span>辰羿<span class="ico">⚙</span></div></aside>
<main class="main"><div class="tb"><div class="bc">Workspace › <b>工作板</b></div><div class="r"><button class="btn">刷新</button><button class="btn k">+ 新增任务</button></div></div>
<div class="fl"><div class="tog"><span class="on">Board</span><span>List</span></div><span class="pill on">All</span><span class="pill">Assigned to me</span><span class="pill">Agent runs · {arun}</span><input class="sr" placeholder="搜索…"></div>
<div class="bd">{board}</div></main></div>'''
open(out,"w",encoding="utf-8").write(doc); print(out)
open(docs_out,"w",encoding="utf-8").write(doc); print(docs_out)
PY