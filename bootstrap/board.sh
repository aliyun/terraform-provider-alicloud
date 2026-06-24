#!/usr/bin/env bash
# bootstrap/board.sh — merge runs/ + escalation/ + scan.json into one status array.
# Emits JSON to stdout: [{id,title,state,summary,priority,pool,project,url,ts}].
#   state: escalated > inflight > merged > done (precedence on id collision)
#   done      ← runs/<id>.md state=pending/absent (审核中, 待合)
#   merged    ← runs/<id>.md state=merged          (已完成, 已合)
#   inflight  ← scan.json tag contains "jarvis-claimed"
#   escalated ← escalation/<id>.md          (first body line = reason)
#   pool      ← scan candidates (待处理/New/Open, no jarvis tag, not tracked); cap 2000/pool by 紧急>高>中>低, pool_total=full count
# Cross-refs scan.json by id for title/priority/pool/project; ids absent from
# scan still appear (title=summary/reason, default project 528766).
# Pure python3; no external deps. JARVIS_ROOT overrides repo root.
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="${JARVIS_ROOT:-$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))}"

python3 - "$root" <<'PY'
import sys, os, re, json, glob
root = sys.argv[1]
DEF_PROJ = 528766
URL = "https://project.aone.alibaba-inc.com/v2/project/{p}/req/{i}"

scan_path = os.path.join(root, ".my-day", "scan.json")
scan = {}
try:
    for it in json.load(open(scan_path)):
        if it.get("id"): scan[str(it["id"])] = it
except Exception:
    pass

pools = {}
proj2pool = {}
try:
    cfg = json.load(open(os.path.join(root, "config", "pools.json")))
    pools = {k: v.get("project") for k, v in cfg.get("pools", {}).items()}
    proj2pool = {v.get("project"): k for k, v in cfg.get("pools", {}).items()}
except Exception:
    pass

def project_for(pool):
    return pools.get(pool) or DEF_PROJ

def pool_for(s):
    # backfill empty pool: infer key from project, else default tf_provider
    p = s.get("pool")
    if p:
        return p
    return proj2pool.get(s.get("project")) or "tf_provider"

def enrich(item):
    s = scan.get(item["id"], {})
    item["title"] = s.get("title") or "Run %s" % item["id"]
    item["priority"] = s.get("priority", "")
    item["pool"] = pool_for(s)
    item["project"] = project_for(item["pool"])
    item["url"] = URL.format(p=item["project"], i=item["id"])
    return item

items = {}  # id -> record ; precedence escalated>inflight>merged>done
RANK = {"done": 0, "merged": 1, "inflight": 2, "escalated": 3}
def put(rec):
    old = items.get(rec["id"])
    if old is None or RANK[rec["state"]] >= RANK[old["state"]]:
        items[rec["id"]] = rec

# done from runs/
for f in glob.glob(os.path.join(root, "runs", "*.md")):
    b = os.path.basename(f)
    if b.startswith("plan-"): continue
    m = re.match(r"\d{4}-\d{2}-\d{2}-(\d+)\.md$", b)
    if not m: continue
    txt = open(f, encoding="utf-8").read()
    g = lambda k: (re.search(r"\*\*%s:\*\*\s*(.+)" % k, txt) or [None, ""])[1].strip()
    rid = g("id") or m.group(1)
    st = "merged" if g("state") == "merged" else "done"
    put(enrich({"id": rid, "state": st, "summary": g("summary"), "ts": g("timestamp")}))

# inflight from scan tags
for i, s in scan.items():
    if "jarvis-claimed" in (s.get("tag") or ""):
        put(enrich({"id": i, "state": "inflight", "summary": s.get("status", ""), "ts": ""}))

# escalated from escalation/
for f in glob.glob(os.path.join(root, "escalation", "*.md")):
    b = os.path.basename(f)
    if not b.endswith(".md") or b == ".gitkeep": continue
    rid = b[:-3]
    reason = next((l.strip() for l in open(f, encoding="utf-8") if l.strip()), "")
    put(enrich({"id": rid, "state": "escalated", "summary": reason, "ts": ""}))

# pool ← scan candidates: status in (待处理,New,Open), no jarvis tag, not already tracked.
PRANK = {"紧急": 0, "高": 1, "中": 2, "低": 3}
OPEN = {"待处理", "New", "Open"}
cands, totals = {}, {}
for i, s in scan.items():
    if i in items: continue
    if s.get("status") not in OPEN: continue
    if "jarvis" in (s.get("tag") or ""): continue
    pk = pool_for(s)  # backfilled key so empty-pool items group + cap correctly
    totals[pk] = totals.get(pk, 0) + 1
    cands.setdefault(pk, []).append(i)
for pool, ids in cands.items():
    ids.sort(key=lambda i: PRANK.get(scan[i].get("priority"), 9))
    for i in ids[:2000]:
        put(enrich({"id": i, "state": "pool", "summary": scan[i].get("status", ""), "ts": "", "pool_total": totals[pool]}))

print(json.dumps(sorted(items.values(), key=lambda x: x["ts"], reverse=True), ensure_ascii=False, indent=2))
PY
