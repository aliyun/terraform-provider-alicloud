#!/usr/bin/env bash
# bootstrap/status.sh — terminal status board. Calls board.sh, groups by state.
# Sections: 跟进中 (inflight) / 已上报 (escalated) / 审核中=待合 (done) / 已完成=已合 (merged).
# Columns: ID | 标题(trunc) | 优先级 | 摘要. ANSI color, footer counts.
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

json="$("$script_dir/board.sh")"

python3 - <<PY
import json, sys
data = json.loads('''$json''')
def C(c): return "" if not sys.stdout.isatty() else c
RST,DIM=C("\033[0m"),C("\033[2m"); B=C("\033[1m")
RED,GRN,YEL,CYN=C("\033[31m"),C("\033[32m"),C("\033[33m"),C("\033[36m")
PRI={"高":RED,"中":YEL,"低":CYN}
def w(s,n):
    s=s or ""; o=0; out=""
    for ch in s:
        cw=2 if ord(ch)>0x2e80 else 1
        if o+cw>n: return out+"…"
        out+=ch; o+=cw
    return out+" "*(n-o)
secs=[("跟进中","inflight",CYN),("已上报","escalated",RED),("审核中(待合)","done",YEL),("已完成(已合)","merged",GRN)]
for label,st,col in secs:
    rows=[x for x in data if x["state"]==st]
    if not rows: continue
    print(f"{B}{col}━━ {label} ({len(rows)}) ━━{RST}")
    for r in rows:
        p=r.get("priority","")
        sm="" if r["summary"]==r["title"] else r["summary"]
        print(f"  {DIM}{r['id']:>9}{RST} {w(r['title'],40)} {PRI.get(p,'')}{p or '·':>2}{RST} {DIM}{w(sm,46)}{RST}")
    print()
n=lambda s:len([x for x in data if x["state"]==s])
print(f"{B}总计{RST} {len(data)}  {CYN}跟进中 {n('inflight')}{RST}  {RED}已上报 {n('escalated')}{RST}  {YEL}审核中 {n('done')}{RST}  {GRN}已完成 {n('merged')}{RST}")
PY
