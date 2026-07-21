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

## 骨架 C · 建关联单 + wrap 主单(转单场景)

```bash
cd "$(git rev-parse --show-toplevel)"   # jarvis 仓根(勿硬编码单机路径)
SRC=78056841                            # 客户主单 id
POOL=1086837                            # 主单所在池
NEW_PROJECT=528766                      # tf_provider 池
ASSIGNEE=521957                         # 承接方工号
CATEGORY=task                           # req/bug/task,按分支决定
PRIORITY="紧急"                         # 复制原单 or 缺陷覆写

# === 1. 建关联单 ===
# --quiet 输出是空格分隔不是 tab; 用 awk '{print $1}' 抓第一列(不带 -F)
NEW_ID=$(bin/a1id -- project workitem create \
  --project "$NEW_PROJECT" \
  --category "$CATEGORY" \
  --title "分支 X · <具体标题>" \
  --assignee "$ASSIGNEE" \
  --priority "$PRIORITY" \
  --body-file /tmp/body-${SRC}.txt \
  --cfs "计划开始日期=$(date +%Y-%m-%d)" \
  --cfs "计划截止日期=$(date -v+3d +%Y-%m-%d)" \
  --cfs "实际工时=0" \
  --quiet 2>&1 | head -1 | awk '{print $1}')

# 显式 [方括号] 打印,便于肉眼验证边界(空格 / 脏字符触发问题)
echo "NEW_ID=[$NEW_ID]"
[ -z "$NEW_ID" ] && { echo "创建失败"; exit 1; }

# === 2. 双向关联(aone 自动双向,单次即可; 第二次调 relation add 会 400) ===
bin/a1id -- project workitem relation add "$SRC" "relate:$NEW_ID" 2>&1 | tail -1

# === 3. 主单 assignee 改到承接方 + status 改「问题解决中」 + wrap 关键节点评论 ===
bin/a1id -- project workitem update "$SRC" --assignee "$ASSIGNEE" 2>&1 | tail -1

# 用 sed 预替换 PLACEHOLDER_NEW_ID(避开 heredoc 展开坑),写到 /tmp/wrap-<SRC>.txt
sed "s/PLACEHOLDER_NEW_ID/${NEW_ID}/g" > /tmp/wrap-${SRC}.txt <<'TXT'
### 结论
<一句话根因 + 分支路由>。已建关联单 PLACEHOLDER_NEW_ID,指派 @<承接方>。

### 关联单
关联单 ID = PLACEHOLDER_NEW_ID
标题 = <具体标题>
项目 <NEW_PROJECT>;双向关联已加。

@<承接方>(<工号>) 详情见关联单,进度请在两侧同步回帖。
TXT

if bash bootstrap/claim.sh claim "$SRC" "$POOL" 2>&1 | tail -1 | grep -q claimed; then
  bash bootstrap/wrap.sh done "$SRC" --summary-file /tmp/wrap-${SRC}.txt 问题解决中 2>&1 \
    | grep -E "工作项|Error" | head -2
  bash bootstrap/claim.sh release "$SRC" "$POOL" 2>&1 | tail -1
else
  echo "[$SRC] claim 失败,主单 wrap 未发;NEW_ID $NEW_ID 已建,记得手动同步"
fi
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
