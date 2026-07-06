#!/usr/bin/env bash
# bootstrap/board.sh — classify scan.json (single source = Aone tags/status) into one status array.
#
# Subcommands:
#   (no args)     Emits the board JSON ARRAY to stdout (unchanged contract; consumed as a
#                 bare array by bridge BoardScheduler → /api/board/sync and by board-html.sh).
#   probe [--text] Emits the F3 度量看板 "probe 飞轮健康度" metrics SECTION (JSON object, or
#                 human digest with --text). Aggregates runs/probe verdicts + probe-drafts
#                 frontmatter + a1 jarvis-probe workitems + playground scenarios + tier0
#                 rotate coverage. Kept as a sibling subcommand precisely so the bare-array
#                 stdout stays byte-compatible with the existing consumers above.
#
# Board array contract (no-arg): [{id,title,state,summary,priority,pool,project,url,ts}].
#   state: escalated > merged > done > idle > inflight > pool (precedence on id collision)
#   escalated ← escalation/<id>.md exists                 (first body line = reason)
#   merged    ← tag jarvis-done  + status 已发布/验收通过/已完成/已发布待需求方验收/已发布待需求排期
#   done      ← tag jarvis-done  + any other status        (审核中)
#   idle      ← tag jarvis-idle                            (jarvis 本轮释放,等待人或下一个 jarvis)
#   inflight  ← tag jarvis-claimed                         (进行中)
#   pool      ← untagged scan candidate (scan already applied exclude_status); cap 2000/pool by 紧急>高>中>低, pool_total=full count + req/bug/task split
# runs/<id>.md is NO LONGER source of done/merged — only enriches summary text;
# scan items with no runs record use title as summary. escalation/ still supplies reason.
# Cross-refs scan.json by id for title/priority/pool/project; escalated ids absent from
# scan still appear (title=reason, default project 528766).
# Pure python3; no external deps. JARVIS_ROOT overrides repo root.
set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
root="$(jarvis_root)"

# ════════════════════════════════════════════════════════════════════
# `probe` subcommand — 飞轮健康度指标段(F3 度量看板)
# ════════════════════════════════════════════════════════════════════
if [ "${1:-}" = "probe" ]; then
    shift 2>/dev/null || true
    text=0
    if [ "${1:-}" = "--text" ]; then text=1; fi

    probe_config="${PROBE_CONFIG:-$root/config/probe.json}"
    proj="$(jq -r '.ticket.project' "$probe_config" 2>/dev/null || true)"
    case "$proj" in ""|null) proj=528766 ;; esac
    tag="$(jq -r '.ticket.tag' "$probe_config" 2>/dev/null || true)"
    case "$tag" in ""|null) tag="jarvis-probe" ;; esac

    # a1 jarvis-probe 工单查询(唯一网络调用);失败/不可解析 → 降级为本地 drafts + WARN。
    # A1 与 scan.sh 同款:默认 bin/a1id(jarvis 身份),JARVIS_A1 供测试打桩。
    A1="${JARVIS_A1:-$root/bin/a1id --}"
    a1_file="$(mktemp)"; trap 'rm -f "$a1_file"' EXIT
    a1_ok=0
    tmo=""
    if command -v timeout >/dev/null 2>&1; then tmo="timeout 30"
    elif command -v gtimeout >/dev/null 2>&1; then tmo="gtimeout 30"; fi
    if raw="$($tmo $A1 project workitem list --project "$proj" --tag "$tag" -f json 2>/dev/null)" \
        && [ -n "$raw" ] && printf '%s' "$raw" | jq -e 'type=="array"' >/dev/null 2>&1; then
        printf '%s' "$raw" > "$a1_file"; a1_ok=1
    else
        echo '[]' > "$a1_file"
        echo "WARN board.sh probe: a1 workitem 查询失败/不可解析(project=$proj tag=$tag),ticket 指标降级为本地数据(source=local)" >&2
    fi

    days="${JARVIS_PROBE_WEEK_DAYS:-7}"
    python3 - "$root" "$probe_config" "$a1_file" "$a1_ok" "$days" "$text" <<'PY'
import sys, os, re, json, glob, datetime

root      = sys.argv[1]
cfgp      = sys.argv[2]
a1_file   = sys.argv[3]
a1_ok     = sys.argv[4] == "1"
days      = int(sys.argv[5] or "7")
text_mode = sys.argv[6] == "1"

def load(p, default=None):
    try:
        return json.load(open(p, encoding="utf-8"))
    except Exception:
        return default

cfg   = load(cfgp, {}) or {}
paths = cfg.get("paths") or {}

now    = datetime.datetime.utcnow()
cutoff = now - datetime.timedelta(days=days)

def parse_iso(s):
    if not s: return None
    s = str(s).strip().replace("Z", "")
    for fmt in ("%Y-%m-%dT%H:%M:%S", "%Y-%m-%dT%H:%M"):
        try: return datetime.datetime.strptime(s, fmt)
        except Exception: pass
    return None

def verdict_dt(v, fname):
    dt = parse_iso(v.get("started_at"))
    if dt: return dt
    m = re.match(r"(\d{8})-", os.path.basename(fname))
    if m:
        try: return datetime.datetime.strptime(m.group(1), "%Y%m%d")
        except Exception: pass
    return None

# ── findings: runs/probe/*.json(本周窗口聚合)──────────────────────────
audit_dir = os.environ.get("PROBE_AUDIT_DIR") or os.path.join(root, paths.get("audit", "runs/probe"))
rounds = tier0 = tier1 = ftotal = 0
by_sev, api_gap, mech = {}, {}, {}
last_t1 = None  # (dt, iso) — 最近一轮 tier1(all-time,新鲜度信号)
for f in sorted(glob.glob(os.path.join(audit_dir, "*.json"))):
    v = load(f)
    if not isinstance(v, dict): continue
    mode = v.get("mode", "")
    if mode == "tier1":
        sdt = parse_iso(v.get("started_at"))
        if sdt and (last_t1 is None or sdt > last_t1[0]):
            last_t1 = (sdt, v.get("started_at"))
    dt = verdict_dt(v, f)
    if dt is None or dt < cutoff:   # 本周窗口外剔除
        continue
    rounds += 1
    if mode == "tier0":
        tier0 += 1
        m = v.get("mech") or "unknown"
        mech[m] = mech.get(m, 0) + 1
    elif mode == "tier1":
        tier1 += 1
    for fd in (v.get("findings") or []):
        ftotal += 1
        sev = fd.get("severity_hint") or "?"
        by_sev[sev] = by_sev.get(sev, 0) + 1
        code = fd.get("code") or ""
        if code.startswith("api_gap"):
            api_gap[code] = api_gap.get(code, 0) + 1

findings = {"rounds": rounds, "tier0_rounds": tier0, "tier1_rounds": tier1,
            "total": ftotal, "by_severity": by_sev, "api_gap": api_gap, "mech": mech}

# ── drafts: escalation/probe-drafts/*.md frontmatter status ───────────
# 采纳率 = filed / (filed + rejected)(只算已决断的;pending/未知不入分母)。
drafts_dir = os.path.join(root, paths.get("drafts", "escalation/probe-drafts"))
dc = {"filed": 0, "pending": 0, "rejected": 0, "other": 0}

def draft_status(path):
    try:
        lines = open(path, encoding="utf-8").read().splitlines()
    except Exception:
        return None
    if not lines or lines[0].strip() != "---":
        return None
    for l in lines[1:]:
        if l.strip() == "---": break
        m = re.match(r"\s*status\s*:\s*(.+)", l)
        if m: return m.group(1).strip().strip('"').strip("'")
    return None

dtotal = 0
for f in glob.glob(os.path.join(drafts_dir, "*.md")):
    if os.path.basename(f) == ".gitkeep": continue
    dtotal += 1
    st = (draft_status(f) or "").lower()
    if   st.startswith("filed"):  dc["filed"] += 1
    elif st.startswith("pending"): dc["pending"] += 1   # covers pending-review
    elif st.startswith("reject"): dc["rejected"] += 1
    else: dc["other"] += 1
decided  = dc["filed"] + dc["rejected"]
adoption = round(dc["filed"] / decided, 3) if decided > 0 else None
drafts = {"total": dtotal, "filed": dc["filed"], "pending": dc["pending"],
          "rejected": dc["rejected"], "other": dc["other"], "adoption_rate": adoption}

# ── tickets: a1 jarvis-probe 工单(降级安全)──────────────────────────
CLOSED_KEYS = ("已完成", "已关闭", "已解决", "已修复", "已发布", "验收通过", "已取消", "已交付", "关闭", "完成")
if a1_ok:
    arr = load(a1_file, []) or []
    by_status, closed = {}, 0
    for w in (arr if isinstance(arr, list) else []):
        s = (w.get("status") or "").strip() or "未知"
        by_status[s] = by_status.get(s, 0) + 1
        if any(k in s for k in CLOSED_KEYS): closed += 1
    tickets = {"source": "aone", "total": len(arr) if isinstance(arr, list) else 0,
               "closed": closed, "by_status": by_status}
else:
    tickets = {"source": "local", "total": None, "closed": None, "by_status": {}}

# ── scenarios: playground <dir>/<product>/<id>/scenario.yaml ──────────
def playground_dir():
    env = os.environ.get("JARVIS_TF_PLAYGROUND")
    if env and os.path.isdir(env): return env
    cd = paths.get("playground_dir")
    if cd and cd != "null" and os.path.isdir(cd): return cd
    return os.path.join(os.path.dirname(root.rstrip("/")), "terraform_playground")

pg = playground_dir()
by_product, stotal = {}, 0
if os.path.isdir(pg):
    for prod in sorted(os.listdir(pg)):
        pdir = os.path.join(pg, prod)
        if not os.path.isdir(pdir): continue
        for sid in sorted(os.listdir(pdir)):
            if os.path.isfile(os.path.join(pdir, sid, "scenario.yaml")):
                by_product[prod] = by_product.get(prod, 0) + 1
                stotal += 1
scenarios = {"total": stotal, "by_product": by_product}

# ── tier0 coverage: .my-day/probe/t0mech-scanned.json(res→last epoch)──
rotate = os.environ.get("PROBE_ROTATE_STATE")
if not rotate:
    wd = os.environ.get("PROBE_WORKDIR") or os.path.join(root, paths.get("workdir", ".my-day/probe"))
    rotate = os.path.join(wd, "t0mech-scanned.json")
scanned = load(rotate, {})
scanned = scanned if isinstance(scanned, dict) else {}
res_n = len(scanned)
last_scan = None
if scanned:
    try:
        mx = max(int(x) for x in scanned.values())
        last_scan = datetime.datetime.utcfromtimestamp(mx).strftime("%Y-%m-%dT%H:%M:%SZ")
    except Exception:
        pass
tier0_coverage = {"resources_scanned": res_n, "last_scanned_at": last_scan}

out = {
    "generated_at":    now.strftime("%Y-%m-%dT%H:%M:%SZ"),
    "window_days":     days,
    "since":           cutoff.strftime("%Y-%m-%dT%H:%M:%SZ"),
    "findings":        findings,
    "drafts":          drafts,
    "tickets":         tickets,
    "scenarios":       scenarios,
    "tier0_coverage":  tier0_coverage,
    "last_tier1_round": last_t1[1] if last_t1 else None,
}

if text_mode:
    ag = ", ".join("%s:%d" % (k.replace("api_gap_", ""), v) for k, v in sorted(api_gap.items())) or "—"
    bs = ", ".join("%s:%d" % (k, v) for k, v in sorted(by_sev.items())) or "—"
    mc = ", ".join("%s:%d" % (k, v) for k, v in sorted(mech.items())) or "—"
    bp = ", ".join("%s:%d" % (k, v) for k, v in sorted(by_product.items())) or "—"
    ts = (", ".join("%s:%d" % (k, v) for k, v in sorted(tickets["by_status"].items()))
          if tickets["source"] == "aone" else "(降级:本地)")
    ar = ("%.0f%%" % (adoption * 100)) if adoption is not None else "—"
    tt = tickets["total"] if tickets["total"] is not None else "—(降级)"
    tcl = tickets["closed"] if tickets["closed"] is not None else "—"
    print("\n".join([
        "== probe 飞轮健康度(近 %d 天,%s ~ %s)==" % (days, out["since"], out["generated_at"]),
        "本周发现: %d findings / %d 轮 (tier0 %d, tier1 %d)" % (ftotal, rounds, tier0, tier1),
        "  严重度: %s" % bs,
        "  api_gap 分布: %s" % ag,
        "  tier0 mech: %s" % mc,
        "草稿: %d (filed %d / pending %d / rejected %d) 采纳率 %s" % (
            dtotal, dc["filed"], dc["pending"], dc["rejected"], ar),
        "工单[%s]: 建单 %s / 关单 %s / 状态 %s" % (tickets["source"], tt, tcl, ts),
        "场景总数(按产品): %d (%s)" % (stotal, bp),
        "tier0 已巡检资源: %d" % res_n,
        "最近 tier1 轮: %s" % (out["last_tier1_round"] or "—"),
    ]))
else:
    print(json.dumps(out, ensure_ascii=False, indent=2))
PY
    exit 0
fi

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

items = {}  # id -> record ; precedence escalated>merged>done>idle>inflight>pool
RANK = {"pool": 0, "inflight": 1, "idle": 2, "done": 3, "merged": 4, "escalated": 5}
def put(rec):
    old = items.get(rec["id"])
    if old is None or RANK[rec["state"]] >= RANK[old["state"]]:
        items[rec["id"]] = rec

# Aone tags = single source for done/merged/idle/inflight.
# 优先级：done > idle > claimed（同时存在时取真完成状态）
MERGED_STATUS = {"已发布", "验收通过", "已完成", "已发布待需求方验收", "已发布待需求排期"}
for i, s in scan.items():
    tag = s.get("tag") or ""
    if "jarvis-done" in tag:
        st = "merged" if s.get("status") in MERGED_STATUS else "done"
        put(enrich({"id": i, "state": st, "summary": "", "ts": ""}))
    elif "jarvis-idle" in tag:
        put(enrich({"id": i, "state": "idle", "summary": s.get("status", ""), "ts": ""}))
    elif "jarvis-claimed" in tag:
        put(enrich({"id": i, "state": "inflight", "summary": s.get("status", ""), "ts": ""}))

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
