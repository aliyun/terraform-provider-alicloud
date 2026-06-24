#!/usr/bin/env bash
# bootstrap/board.sh — classify scan.json (single source = Aone tags/status) into one status array.
# Emits JSON to stdout: [{id,title,state,summary,priority,pool,project,url,ts}].
#   state: escalated > merged > done > inflight > pool (precedence on id collision)
#   escalated ← escalation/<id>.md exists                 (first body line = reason)
#   merged    ← tag jarvis-done  + status 已发布/验收通过/已完成/已发布待需求方验收
#   done      ← tag jarvis-done  + any other status        (审核中)
#   inflight  ← tag jarvis-claimed                         (进行中)
#   pool      ← untagged scan candidate (scan already applied exclude_status); cap 2000/pool by 紧急>高>中>低, pool_total=full count + req/bug/task split
# runs/<id>.md is NO LONGER source of done/merged — only enriches summary text;
# scan items with no runs record use title as summary. escalation/ still supplies reason.
# Cross-refs scan.json by id for title/priority/pool/project; escalated ids absent from
# scan still appear (title=reason, default project 528766).
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

# runs/<id>.md → summary only (no longer drives state). Map id → summary/ts.
runs = {}
for f in glob.glob(os.path.join(root, "runs", "*.md")):
    b = os.path.basename(f)
    if b.startswith("plan-"): continue
    m = re.match(r"\d{4}-\d{2}-\d{2}-(\d+)\.md$", b)
    if not m: continue
    txt = open(f, encoding="utf-8").read()
    g = lambda k: (re.search(r"\*\*%s:\*\*\s*(.+)" % k, txt) or [None, ""])[1].strip()
    runs[g("id") or m.group(1)] = {"summary": g("summary"), "ts": g("timestamp")}

def enrich(item):
    s = scan.get(item["id"], {})
    item["title"] = s.get("title") or "Run %s" % item["id"]
    item["priority"] = s.get("priority", "")
    # default uncategorized non-pool items to "req" so the badge + type filter
    # keep them visible instead of vanishing as uncategorized.
    item["category"] = s.get("category") or ("req" if item["state"] in ("done", "merged", "escalated") else "")
    item["pool"] = pool_for(s)
    item["project"] = project_for(item["pool"])
    item["url"] = URL.format(p=item["project"], i=item["id"])
    # runs/ enrich: prefer its summary text + ts when present, else fall back to title
    r = runs.get(item["id"], {})
    item["summary"] = item.get("summary") or r.get("summary") or item["title"]
    item["ts"] = item.get("ts") or r.get("ts") or ""
    return item

items = {}  # id -> record ; precedence escalated>merged>done>inflight>pool
RANK = {"pool": 0, "inflight": 1, "done": 2, "merged": 3, "escalated": 4}
def put(rec):
    old = items.get(rec["id"])
    if old is None or RANK[rec["state"]] >= RANK[old["state"]]:
        items[rec["id"]] = rec

# Aone tags = single source for done/merged/inflight.
MERGED_STATUS = {"已发布", "验收通过", "已完成", "已发布待需求方验收"}
for i, s in scan.items():
    tag = s.get("tag") or ""
    if "jarvis-claimed" in tag:
        put(enrich({"id": i, "state": "inflight", "summary": s.get("status", ""), "ts": ""}))
    elif "jarvis-done" in tag:
        st = "merged" if s.get("status") in MERGED_STATUS else "done"
        put(enrich({"id": i, "state": st, "summary": "", "ts": ""}))

# escalated from escalation/ (overrides tag-derived state)
for f in glob.glob(os.path.join(root, "escalation", "*.md")):
    b = os.path.basename(f)
    if not b.endswith(".md") or b == ".gitkeep": continue
    rid = b[:-3]
    reason = next((l.strip() for l in open(f, encoding="utf-8") if l.strip()), "")
    put(enrich({"id": rid, "state": "escalated", "summary": reason, "ts": ""}))

# pool ← scan candidates: any scan item not already tracked + no jarvis tag.
# scan already applied exclude_status, so trust it — no active-status whitelist; status passes through.
PRANK = {"紧急": 0, "高": 1, "中": 2, "低": 3}
cands, totals, catc = {}, {}, {}  # pool → ids ; pool_total ; pool → {req,bug,task counts}
for i, s in scan.items():
    if i in items: continue
    if "jarvis" in (s.get("tag") or ""): continue
    pk = pool_for(s)  # backfilled key so empty-pool items group + cap correctly
    totals[pk] = totals.get(pk, 0) + 1
    c = catc.setdefault(pk, {"req": 0, "bug": 0, "task": 0})
    cat = s.get("category")
    if cat in c: c[cat] += 1
    cands.setdefault(pk, []).append(i)
for pool, ids in cands.items():
    ids.sort(key=lambda i: PRANK.get(scan[i].get("priority"), 9))
    cc = catc[pool]
    for i in ids[:2000]:
        put(enrich({"id": i, "state": "pool", "summary": scan[i].get("status", ""), "ts": "",
                    "pool_total": totals[pool], "pool_req": cc["req"], "pool_bug": cc["bug"], "pool_task": cc["task"]}))

print(json.dumps(sorted(items.values(), key=lambda x: x["ts"], reverse=True), ensure_ascii=False, indent=2))
PY
