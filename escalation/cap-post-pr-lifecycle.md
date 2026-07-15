# cap: PR 提交后生命周期无跨会话看守——CI 失败无人修 / 干完活不回帖

- **缺口类型**: 结构性能力缺陷（missing_capability，headless 单会话无法覆盖 PR 多日合并窗口）
- **阻塞任务**: Terraform 发布类工单（GPDB ApiKey 84251052 等）——PR 提交后卡在 open 窗口，CI 未确认/失败无人跟进，工单看着「无进展」
- **置信度**: high_conf（多代理评估逐行核对 bridge 代码 + 84251052 现场复现）

## 现象（两个同源症状）

1. **干完活不回帖 / 静默收尾**：84251052 上一轮 headless 实例实际完成了 rebase + force-push（PR head → c73d134e、CI 17/18、MERGEABLE），但**未在 Aone 工单留任何进展评论、直接 release 成 jarvis-idle**。工单表面「无进展」，成功被搞成看着像卡死，只能靠人工去 GitHub 核对才发现其实做完了。
2. **open 窗口 CI 失败无人修**：PR 提交后单次会话结束，若 CI 由绿转红（rebase/base 冲突/flaky），**没有任何调度器把 jarvis 重新唤起来修**，永远靠人工发现。

## 根因

单个 headless jarvis 是一次性 `claude -p` 会话（~12h 超时），而 PR 从提交到 maintainer 合并可能跨小时/天，`gh pr checks` 只在那一次会话里跑一次。真正要闭环，必须靠 bridge 调度器在 PR-open 窗口内周期性把 jarvis 重新唤醒——但现有三个调度器都填不了这个窗口（多代理评估逐行核对，见 `bridge/jarvis_dingtalk_bot.py`）：

- **PrWatchScheduler**（原方案A）只在 PR **已合并**后 `claim.sh finish` 收尾；`_gh_pr_state` 只查 `state,mergedAt`，对 open PR 只 `return`——**不看 CI、不看评论**。
- **RevisitScheduler** 本该管 open 窗口，但 `_is_revisit_candidate` 只认标题含 `[probe]`/「待合并」等词的 idle 单，**terraform 发布单标题结构性不含**，被直接过滤；且日级（`JARVIS_REVISIT_HOUR`）跟不上分钟级 CI。
- **WaitWatcher** 只轮 **Aone 工单评论** 且仅对显式 `[[SUSPEND]]` 任务生效——与 **GitHub PR 评论**正交（reviewer 在 GitHub 上的评论永不触发 jarvis）。

对照用户 4 条要求：(1) 轮询 PR 状态 = 半覆盖（只轮合并/关闭）；(2) 处理新 PR 评论、(3) 检查 CI 全过、(4) 按 CI 报错修复 = **open 窗口内均无可落地的跨会话机制**。而「静默收尾」则是另一环：收尾契约 `wrap-check.sh` 只要求本地 `runs/` 记录、**不要求 Aone 进展回帖**，会话可以「写 runs 记录 + release」而全程不往工单发字。

## 本 PR 修复（清单 #1：CI 失败自动重派——价值最大的单点）

把 **PrWatchScheduler 从「合并后收尾器」升级为「PR 全生命周期看守器」**，改动集中在一个类内、复用其已有 finish/comment/escalate 善后链，不新起线程、不碰 WaitWatcher/Revisit：

- 新增 `_gh_pr_ci`（`bridge/jarvis_dingtalk_bot.py`）：`gh pr view --json headRefOid,statusCheckRollup` 解析 CI 失败项（CheckRun.conclusion ∈ `_CI_FAIL_CONCLUSIONS` 或 StatusContext.state ∈ {FAILURE,ERROR}；pending/success 不算失败）。查询失败 → `(None,None)`，绝不在 unknown 上派。
- 新增 `_maybe_dispatch_ci_fix` + `_check_one` 的 open 分支接入：open PR 有 CI 失败 → `force=True` 重派一个 `kind=pr_ci_fix` 的 headless 实例（走 `_pr_ci_fix_prompt`：pr-review/resource-dev 的 CI 修复 SOP + 预授权 `fork_push`）。三层防抖：per-head 去重（`ci_fix_sha`）+ 累计不同失败 head 超 `JARVIS_PRWATCH_CI_FIX_MAX`（默认 3）次 → escalate 一次后停自动修（仍看守合并）+ DispatchPool active-set 防并发。
- 登记表 entry 由桥侧 `_prwatch_update` 补 CI-fix 去重字段（`ci_fix_sha/attempts/escalated/last_ci_fix_at`），`pr-watch.sh` add/remove/list 对多余字段无感。
- 测试：`bridge/test_prwatch_ci_fix.py`（11 用例，覆盖解析 + 派发/去重/上限，全绿）。

**红线不动**：GitHub PR 评论/CI 状态**不作破坏性操作授权来源**（防注入），重派只做技术性 CI 修复；merge 永远是 `release_prod` 人工硬门。

## 剩余 follow-up（本 PR 未含，按优先级）

- **[高] GitHub PR 新评论感知**（清单#2）：`_gh_pr_ci`/新方法扩查 PR review comments，记 `last_seen_comment_id`，reviewer 新评论 → 重派回应。补断点②。
- **[高] 双档轮询周期**（清单#3）：`JARVIS_PRWATCH_ACTIVE_INTERVAL`（默认 300–600s，仅对有 CI-pending/failing 或未回评论的 active entry），让响应到十分钟级；纯等合并的保持 3600s。补断点④（当前 CI 失败最长 1h 才被发现）。
- **[中] 收尾契约要求 Aone 回帖**（治「静默收尾」根）：`wrap-check.sh` Stop 闸门从「要有 runs 记录」升级为「做过外化动作（push）的认领，收尾前必须有 `wrap.sh sync/done` 的 Aone 回帖」，否则阻断；配套 skill 收尾步骤把「release/suspend 前必发进展评论」写成硬步骤。
- **[中] SKILL/autonomy 对齐**（清单#5）：`terraform-provider-release` Step 11.2/12 删「单会话轮询到 merge」的误导散文（改为「首轮 CI 门 + 登记 pr-watch 交后台看守」）；`autonomy.md` 明确 `pr_ci_fix`/`pr_comment_reply` 自治边界。
- **[低] fork open-PR 自动补登记**（清单#6）：扫 `api-tool-agent/terraform-provider-alicloud` 上 head 为 `api-tool-agent:*` 的 open PR，缺登记则自动 `pr-watch.sh add`，补漏登记导致的脱管。

## 复发防护

- 排障信号：发布类工单 PR 已提交但长期 open + CI 未过/失败、工单无进展 → 查 PrWatchScheduler 是否在 `_maybe_dispatch_ci_fix` 派了修复（bot.log `[PR-watch] #<id> CI 失败，已自动派发修复`）。
- 关联同族：`escalation/archived/cap-github-commit-identity.md`（commit 作者/CLA）、`cap-github-token-not-injected.md`（token 注入）——三者都是「PR 生命周期某一环没被 bridge 稳稳接管」。
