# 批量 bookend 参考模板

> 场景:一次处理 N 张工单(转单 / 关单提示 / 追料 / 进度跟进等)。
> 目标:避开 `claim.sh lost race` 静默继续、Bash 工具 2min timeout 截断、
> `wrap.sh done` heredoc 反引号/`$var:字母` 展开等已知踩坑。
>
> 引用见 `.Codex/skills/aone-triage/references/tf-customer-request-routing.md`
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

## 骨架 C · pure datasource source-only 源单路由

本骨架只用于纯 datasource：仅涉及 `data.alicloud_xxx` 查询、过滤、分页、输出字段或 Read，
且不含 resource 变更。resource+datasource 混合、G 全局与手写 resource D 都不得进入。
RD route phase **只幂等同步源单 assignee + per-type progress_status**；
**bridge executor 独占源单 claim/唯一回复/tag/release/finish**。历史 relation 只读保留；
禁止 create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766。

```bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
SRC=78056841
URGENT="${URGENT:?1=紧急，0=非紧急}"
URGENT_ASSIGNEE=521957
NONURGENT_ASSIGNEE=484483
SOURCE_STATUS="${SOURCE_STATUS:?从 pools.json progress_status[workitemType] 解析}"
SOURCE_ROUTE_DRIFT="${SOURCE_ROUTE_DRIFT:-0}"

if [ "$URGENT" = 1 ]; then
  PURE_DATASOURCE_ASSIGNEE="$URGENT_ASSIGNEE"
else
  PURE_DATASOURCE_ASSIGNEE="$NONURGENT_ASSIGNEE"
fi

# point-read 仅用于比较源单 owner/status，并可只读引用历史 relation 上已有 PR 防重复。
# 不删除/迁移/关闭/改派历史 relation；它不是开发、完成或 blocker 门。
if [ "$SOURCE_ROUTE_DRIFT" = 1 ]; then
  bin/a1id as terraform-rd -- project workitem update "$SRC" \
    --assignee "$PURE_DATASOURCE_ASSIGNEE" --status "$SOURCE_STATUS"
fi

# 后续在源单上下文完成 RD→QA；finalizer 仅返回 AONE_RESULT。
# open PR + QA pass 时由 bridge executor release 源单，不 finish。
```

## 骨架 D · D/G source-only + D route DM

本骨架用于所有 D/G；pure datasource 继续走上方窄 source-only。控制面 executor 独占源工单
bookend；RD finalizer 只同步源单 route 字段，D 再通过持久化 ledger enqueue owner DM。
D/E/G 严禁对 528766 执行 create/reuse/reassign/relation/claim/wrap/release/finish。

```bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
SRC=78056841
SOURCE_PROJECT=1086837
ROUTE="${ROUTE:?d 或 g}"
SUBTYPE="${SUBTYPE:-}"   # D: handwritten-urgent|handwritten-normal|generated；G 留空
SOURCE_STATUS="${SOURCE_STATUS:?从 pools.json progress_status[workitemType] 解析}"
SOURCE_ROUTE_DRIFT="${SOURCE_ROUTE_DRIFT:-0}"

case "$ROUTE:$SUBTYPE" in
  d:handwritten-urgent) SOURCE_ASSIGNEE=521957 ;;
  d:handwritten-normal) SOURCE_ASSIGNEE=484483 ;;
  d:generated) SOURCE_ASSIGNEE=429768 ;;
  g:) SOURCE_ASSIGNEE=521957 ;;
  *) echo "非法 D/G subtype" >&2; exit 2 ;;
esac

# === 1. 先幂等同步源单 owner + per-type status ===
# point-read 可只读历史 relation 防重复代码工作，但不得更新/claim/bookend relation。
if [ "$SOURCE_ROUTE_DRIFT" = 1 ]; then
  bin/a1id as terraform-rd -- project workitem update "$SRC" \
    --assignee "$SOURCE_ASSIGNEE" --status "$SOURCE_STATUS"
fi

# === 2. D 才 enqueue 类型化 DM；G 不发新增 route DM ===
if [ "$ROUTE" = d ]; then
  notify_result="$(
    python3 -m bridge.terraform_route_notify \
      --ticket "$SRC" --project "$SOURCE_PROJECT" --subtype "$SUBTYPE"
  )" || notify_rc=$?
  notify_rc="${notify_rc:-0}"
  echo "$notify_result"
  # rc=0 包括 posted/suppressed/durable pending；均不阻断开发。
  # rc=1 表示 ledger 未持久化：继续开发，但 AONE_RESULT 必须写“通知未完成”，不得宣称完成。
  [ "$notify_rc" -le 1 ] || exit "$notify_rc"
fi

# === 3. 继续 RD→QA；最后交 AONE_RESULT ===
# 不因 owner/status、通知或历史 relation 观察等待。
# open PR + QA pass 时由 bridge executor release 源单，不 finish。
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
