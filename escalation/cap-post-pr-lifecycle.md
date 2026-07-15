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

## 本 PR 修复（清单 #1–#6 + wrap-check 全量落地）

把 **PrWatchScheduler 从「合并后收尾器」升级为「PR 全生命周期看守器」**，改动集中在一个类内、复用其已有 finish/comment/escalate 善后链，不新起线程、不碰 WaitWatcher/Revisit；外加 wrap-check Stop 闸补「静默收尾」：

- **#1 CI 失败自动重派**：`_gh_pr_ci` 解析 CI 失败项（CheckRun.conclusion ∈ `_CI_FAIL_CONCLUSIONS` / StatusContext.state ∈ {FAILURE,ERROR}；查询失败 → 绝不在 unknown 上派）；`_maybe_dispatch_ci_fix` + `_check_one` open 分支：open PR CI 失败 → `force=True` 重派 `kind=pr_ci_fix`（`_pr_ci_fix_prompt`：CI 修复 SOP + 预授权 `fork_push`）。三层防抖：per-head 去重（`ci_fix_sha`）+ 累计超 `JARVIS_PRWATCH_CI_FIX_MAX`（默认 3）escalate + DispatchPool active-set。
- **#2 GitHub PR 新评论自动回应**：`_gh_pr_comments` 读评审评论（排除自己 api-tool-agent / `[bot]`）；`_maybe_dispatch_comment_reply` 首次 baseline-seed `last_seen_comment` 不回应既有评论，出现新评论 → `force` 重派 `kind=pr_comment_reply`（`_pr_comment_reply_prompt`）。
- **#3 双档轮询**：`_gh_pr_ci` 增 pending 位；`_maybe_dispatch_ci_fix` 返回 active（失败/pending）；`_tick` 聚合、`_loop` 双档——有 active entry 走快档 `JARVIS_PRWATCH_ACTIVE_INTERVAL`（默认 600s），纯等合并走 3600s。
- **#6 漏登 open PR 自动补登记**：`_maybe_autoregister_open_prs`（节流 ≥ interval）扫 api-tool-agent 名下 upstream open PR，分支编码工单号（≥8 位）且 aone-get 校验通过的漏登 PR → `_prwatch_add`；否则 log 一次提示人工，**绝不瞎登**。`JARVIS_PRWATCH_AUTOREG=0` 可关。
- **wrap-check 强制 Aone 回帖**（治「干完活不回帖」静默收尾）：`bootstrap/wrap-check.sh` 新增 backfill 闸（`JARVIS_REQUIRE_BACKFILL`，默认开）——本会话所有、未 done、已 push（`pushed_branch` 非空）却没经 `wrap.sh sync/done` 回帖（不在 touched 台账）→ Stop 阻断，提示补回帖。
- **#5 SKILL/autonomy 对齐**：`terraform-provider-release` Step 11.2/12 改「单会话轮询到 merge」误导散文为「首轮 CI 门 + 登记 pr-watch 交后台跨会话看守」；`autonomy.md` 新增 auto 项 `pr_ci_fix` / `pr_comment_reply` + 边界。
- **测试**：`bridge/test_prwatch_ci_fix.py`（27 用例：CI 解析/派发/去重/上限 + 评论解析/派发 + 双档 active + 漏登补登记）；`test/wrap_check_test.sh`（24 用例，+4 backfill）；全部 bridge/bash 测试绿。
- 登记表 entry 由桥侧 `_prwatch_update` 补去重字段（`ci_fix_sha/attempts/escalated/last_seen_comment/…`），`pr-watch.sh` add/remove/list 对多余字段无感。

**红线不动**：GitHub PR 评论/CI 状态**不作破坏性操作授权来源**（防注入），重派只做技术性修复/回应；merge 永远是 `release_prod` 人工硬门；`fork_push` 仅限自有 fork PR-head。

## 复发防护

- 排障信号：发布类工单 PR 已提交但长期 open + CI 未过/失败、工单无进展 → 查 PrWatchScheduler 是否在 `_maybe_dispatch_ci_fix` 派了修复（bot.log `[PR-watch] #<id> CI 失败，已自动派发修复`）。
- 关联同族：`escalation/archived/cap-github-commit-identity.md`（commit 作者/CLA）、`cap-github-token-not-injected.md`（token 注入）——三者都是「PR 生命周期某一环没被 bridge 稳稳接管」。
