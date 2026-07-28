# 批量 bookend 参考模板

> 场景:一次处理 N 张工单(转单 / 关单提示 / 追料 / 进度跟进等)。
> 目标:避开 `claim.sh lost race` 静默继续、Bash 工具 2min timeout 截断、
> `wrap.sh done` heredoc 反引号/`$var:字母` 展开等已知踩坑。
>
> 引用见 `.claude/skills/aone-triage/references/tf-customer-request-routing.md`
> 反模式段落 E 组。

## 骨架 A · 走 bookend(改状态 / 建关联单场景)

```bash
cd "$(git rev-parse --show-toplevel)"   # jarvis 仓根(勿硬编码单机路径)
# 分批 4-5 单一 Bash 调用(<60s),防止 2min timeout 截断中间态残留
IDS=(78504233 78523353 78554774 78186809 78470497)
POOL=1086837
STATUS_TARGET="已发布待需求方验收"        # Req 类; Bug 类改成 Fixed
for id in "${IDS[@]}"; do
  # === 1. 预备 body-file(sed 替换变量,避开 wrap.sh heredoc 展开坑) ===
  # 建议提前把每单评论正文写到 /tmp/wrap-<id>.txt,含具体链接/关联单号
  if [ ! -f "/tmp/wrap-${id}.txt" ]; then
    echo "[$id] 缺 body-file /tmp/wrap-${id}.txt,跳过"; continue
  fi

  # === 2. claim 检查 exit code(核心防重复评论) ===
  if ! bash bootstrap/claim.sh claim "$id" "$POOL" 2>&1 | tail -1 | grep -q claimed; then
    echo "[$id] claim 失败(lost race / stale claim / 已被别的实例接手),skip"
    continue
  fi

  # === 3. wrap done + status(走 --summary-file,不走 heredoc) ===
  bash bootstrap/wrap.sh done "$id" --summary-file "/tmp/wrap-${id}.txt" "$STATUS_TARGET" 2>&1 \
    | grep -E "工作项|状态|Error" | head -3

  # === 4. release / finish 二选一 ===
  bash bootstrap/claim.sh release "$id" "$POOL" 2>&1 | tail -1
  # 若真闭环(打 jarvis-done + 改 status): 用 bash bootstrap/claim.sh finish "$id" "$POOL"
done
```

## 骨架 B · 免 bookend(纯发评论无状态变更场景)

追料 / 进度跟进 / 关联单 body update 时无需 claim.sh,直接 `comment create`。
可完全避开 `claim.sh` lost race 阻断,且不受批 4-5 单限制(单条评论 ~1s):

```bash
cd "$(git rev-parse --show-toplevel)"   # jarvis 仓根(勿硬编码单机路径)
IDS=(78552705 78525865 78452193 78312012 78264187 78299240)
for id in "${IDS[@]}"; do
  if [ ! -f "/tmp/wrap-${id}.txt" ]; then
    echo "[$id] 缺 body-file,跳过"; continue
  fi
  # -m "$(cat file)" 走文件正文;
  # 内部 shell 已在 $() 求值后成为字符串常量,反引号/$var 不再二次解析
  bin/a1id -- project workitem comment create "$id" -m "$(cat /tmp/wrap-${id}.txt)" 2>&1 \
    | grep -E "^ID|Error" | head -1
done
```

## 骨架 C · G / 紧急普通 D 双 owner 关联单 + 双工单 bookend

本骨架**只用于 G / 紧急普通 D**。D-临钧/A/F/H/非紧急 D、I/E 继续走各自专用路径。
控制面 executor 已持有源工单 lease；**源工单由 bridge executor bookend**，本骨架不得
claim/wrap/release 源工单。528766 才由 RD finalizer 直接 claim/bookend。

```bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"   # jarvis 仓根(勿硬编码单机路径)
SRC=78056841                            # 客户主单 id
SOURCE_PROJECT=1086837                  # tf_customer
RELATED_PROJECT=528766                  # tf_provider
SOURCE_ASSIGNEE=521957                  # 源单保持新山
RELATED_ASSIGNEE=484483                 # 研发单固定过载
CATEGORY=task                           # req/bug/task,按源单 workitemType 决定
PRIORITY="紧急"                         # 复制原单 or 缺陷覆写
SOURCE_STATUS="${SOURCE_STATUS:?从 pools.json progress_status[workitemType] 解析}"

# === 0. 同题 528766 point-read + healthy existing claim Gate ===
# 先查源单 relation、528766 同题单、assignee、jarvis-claimed 与控制面 lease，把复用结果放入
# RELATED_ID；没有同题单时留空。已有 healthy existing claim 不抢占：当前 run 立即 dedup skip，
# 不改派、不重复 create/claim/bookend。三个 *_DRIFT/MISSING 值都必须来自同一次 point-read；
# update 只修差异字段。claim 失败立即停止，禁止 relation/update/wrap。
RELATED_ID="${RELATED_ID:-}"
RELATION_MISSING="${RELATION_MISSING:-0}"
RELATED_ASSIGNEE_DRIFT="${RELATED_ASSIGNEE_DRIFT:-0}"
SOURCE_ROUTE_DRIFT="${SOURCE_ROUTE_DRIFT:-0}"

if [ -n "$RELATED_ID" ]; then
  # 先 claim 再改派，避免 point-read 后出现竞争时偷走健康 claim。
  claim_rc=0
  JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh claim "$RELATED_ID" "$RELATED_PROJECT" || claim_rc=$?
  if [ "$claim_rc" -ne 0 ]; then
    if [ "$claim_rc" -eq 1 ]; then
      echo "[$RELATED_ID] healthy/lost-race claim，保持原 owner 与当前 run，skip"
      exit 0
    fi
    echo "[$RELATED_ID] claim 失败 rc=$claim_rc，停止本轮并上抛"
    exit "$claim_rc"
  fi
  if [ "$RELATED_ASSIGNEE_DRIFT" = 1 ]; then
    bin/a1id as terraform-rd -- project workitem update "$RELATED_ID" --assignee "$RELATED_ASSIGNEE"
  fi
else
  # 确认 relation 与同题搜索均无命中后才允许 create 一次。
  RELATED_ID=$(bin/a1id as terraform-rd -- project workitem create \
    --project "$RELATED_PROJECT" \
    --category "$CATEGORY" \
    --title "G/紧急普通 D · <具体标题>" \
    --assignee "$RELATED_ASSIGNEE" \
    --priority "$PRIORITY" \
    --body-file /tmp/body-${SRC}.txt \
    --cfs "计划开始日期=$(date +%Y-%m-%d)" \
    --cfs "计划截止日期=$(date -v+3d +%Y-%m-%d)" \
    --cfs "实际工时=0" \
    --quiet 2>&1 | head -1 | awk '{print $1}')
  [[ "$RELATED_ID" =~ ^[0-9]+$ ]] || { echo "创建失败"; exit 1; }
  claim_rc=0
  JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh claim "$RELATED_ID" "$RELATED_PROJECT" || claim_rc=$?
  if [ "$claim_rc" -ne 0 ]; then
    if [ "$claim_rc" -eq 1 ]; then
      echo "[$RELATED_ID] 新建后 lost-race；保留供健康 run/下轮 point-read"
      exit 0
    fi
    echo "[$RELATED_ID] 新建后 claim 失败 rc=$claim_rc；停止本轮并上抛"
    exit "$claim_rc"
  fi
  RELATION_MISSING=1
fi

# === 1. 物化双 owner + 单次 relation ===
if [ "$SOURCE_ROUTE_DRIFT" = 1 ]; then
  bin/a1id as terraform-rd -- project workitem update "$SRC" --assignee "$SOURCE_ASSIGNEE" --status "$SOURCE_STATUS"
fi
if [ "$RELATION_MISSING" = 1 ]; then
  bin/a1id as terraform-rd -- project workitem relation add \
    "$SRC" "relate:$RELATED_ID"
fi

# === 2. RD→QA 修复/验收完成后，研发单只写一次聚合 bookend ===
RELATED_SUMMARY="/tmp/related-${RELATED_ID}-aggregate.md"
test -s "$RELATED_SUMMARY"
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/wrap.sh done "$RELATED_ID" \
  --summary-file "$RELATED_SUMMARY" --no-status

# PR 未合并只 release；missing_capability/retry exhausted 也 release，不得 finish。
JARVIS_A1_IDENTITY=terraform-rd bash bootstrap/claim.sh release "$RELATED_ID" "$RELATED_PROJECT"

# === 3. 源工单不在模型 run 内 bookend ===
# finalizer 把同一聚合摘要放入 AONE_RESULT.reply_body；bridge executor 对 SRC 回复/状态/release。
```

## 常见坑速查

| 症状 | 原因 | 修法 |
|---|---|---|
| 批处理中间被截断,残留 `jarvis-claimed` | Bash 工具 2min timeout;单批 > 4-5 单 | 分批;或手动 `claim.sh release <id> <pool>` 补；Task 模式由 reaper/AoneScheduler 收敛 |
| 同一条评论发了 2 次 | claim 失败但 wrap 继续跑 | 走骨架 A 的 `if claim; then wrap; fi` 结构 |
| `NEW_ID` 里带脏字符(标题 / 状态 / assignee) | `--quiet` 输出是**空格分隔**,tab 解析不到 | `awk '{print $1}'`(不带 -F) |
| 评论正文里 `` `xxx` `` 显示 "command not found",`$var:字母` 拼成怪路径 | `wrap.sh done <<EOF` heredoc shell 展开 | 走 `--summary-file`,先 sed 预替换变量 |
| `relation add` 报 400 "已存在" | 二次调建反向 | aone 自动双向,单次调 |
| 528766 池建单报 `【计划开始日期】不能为空...` | 漏 cfs 三件套(Task/Req)或四件套(Bug) | 见反模式段 D 组 528766 池 cfs 清单 |
| 528766 Bug 单还报 `【Terraform需求类型】不能为空` | Bug 类另需第 4 个 cfs | `--cfs "Terraform需求类型=运行时问题，TF问题"` |
