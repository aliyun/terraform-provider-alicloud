# Jarvis 舰队协调：崩溃续跑 + 跨仓可见 + 死活检测

> 状态：设计 v0 · 2026-06-25 · 待用户评审

## 一句话目标

一台 7×24 mac mini 上同时跑 0~N 个 Jarvis（Claude 会话），实例随时可能崩。让任务**不丢、不卡、不被错领**：崩了能接力续跑，跨仓 worktree 互相可见，死实例能被发现并回收——且不打扰正在被 Tata 定向委派的实例。

## 背景与三个缺口

现状：`claim.sh`（Aone 标签认领）+ `sweep.sh`（45min 超时升级）+ `reconcile.sh`（对账）+ `runs/`/`.my-day` 本地账本。三个真缺口：

1. **崩溃丢进度**：实例挂在任务中途，worktree + 半成品分支无人接管，重进只能从头重跑。
2. **跨仓不可见**：worktree 常建在其它代码仓库目录下，任务 → worktree 路径/分支无共享登记，实例间互不可知。
3. **无死活检测**：只有粗超时，无心跳，慢任务与死任务分不清，任务可挂任意久。

关键判断：**单机多进程,非真分布式**——所有实例共享文件系统,`git worktree list` 可枚举跨仓 worktree,本地共享目录即可协调,无需靠 Aone 标签当跨机协调器。

## 五个已锁定决策

1. **多机就绪,本地先行**：协调器抽象成接口，本地 FS 实现先上，多机时换 Aone/Redis 后端，签名不变（保 `claim.sh` 写 hostname 的多机口子）。
2. **接力续跑(非断点续传)**：死实例 worktree+分支保留，新实例从阶段 checkpoint 续，已 commit 不重做；不复活 Claude 会话上下文。
3. **角色隔离**：dispatch 实例（Tata 定向）只干本单，**心跳/登记照打但永不 adopt**；triage 实例才扫孤儿接管。
4. **watchdog 兜底**：常驻 launchd 进程只判死活、标 orphaned，**自己不跑 triage**；adopt 由真 triage 会话做。
5. **心跳进程级 0 上下文**：心跳是启动器的活,Claude 全程不碰；checkpoint 才由 Claude 调,仅阶段切换 1 次。

## 架构：协调接口 + 四件

### ① 协调接口 `bootstrap/coord.sh`
所有 jarvis 脚本只调它,不直接碰协调文件。命令：

| 命令 | 作用 |
|---|---|
| `register <role>` | 启动器登记本实例（role∈dispatch\|triage），写 `instances/<id>.json` |
| `heartbeat` | 刷新 `instances/<id>.hb` mtime（启动器后台调，0 上下文） |
| `checkpoint <aone_id> <stage>` | 更新 `tasks/<aone_id>.json` 阶段（Claude 阶段切换调） |
| `list-orphans` | 列出 owner 已死/超时的 orphaned 任务 |
| `adopt <aone_id>` | triage 实例接管孤儿任务 owner |

本地实现读写 `.my-day/`；后端可换。

### ② 登记结构（共享 `.my-day/`，接口屏蔽）
- `instances/<id>.json`：`{id, role, pid, host, started, last_hb, task}`——谁活着、在干啥。
- `instances/<id>.hb`：心跳文件，mtime 即最后存活时刻。
- `tasks/<aone_id>.json`：`{aone_id, owner_instance, stage, worktree, branch, repo, updated}`，stage ∈ `claimed→verifying→coding→prestage→done`。**worktree 路径在此登记**,跨仓可枚举。
- Aone `jarvis-claimed` 标签仍是认领真源；本地注册表是补充，本地崩不影响他人认领。

### ③ watchdog（launchd，~90s）
仅判死活：单机 `kill -0 <pid>` 直接判，或 `last_hb` 过 TTL（默认 3min）→ 把该实例 task 标 `orphaned`，worktree 留着不删。无 triage 实例时落 `escalation/` 兜底。**绝不跑 triage**。

### ④ adopt（仅 triage 开局）
扫 `orphaned`：`git worktree list` 找回分支，有未完阶段则接管 owner、从 stage 续（已 commit 不重做），回填 Aone。**dispatch 实例跳过**，只干本单。`sweep.sh`（超时）+`reconcile.sh`（对账）作二次兜底。

## 心跳机制（0 上下文）
启动器拉起会话时后台起循环，与会话同生共死：
```bash
while kill -0 $CLAUDE_PID 2>/dev/null; do touch .my-day/instances/$ID.hb; sleep 60; done
```
会话崩 → PID 没 → touch 停 → 90s 后 watchdog 判死。Claude 全程不参与。单机可仅靠 `kill -0`,`.hb` 是为多机后端留的抽象。

## 进度写哪
- **本地 `tasks/<id>.json`(机读)**：`coord.sh checkpoint` 高频细粒度,给接力实例;纯本地不打网络。
- **Aone(人读,唯一真源)**：`wrap.sh sync/done` 在里程碑回填;低频粗粒度。
- 两线解耦,默认大节点不自动捎 Aone;checkpoint 整单仅 4-5 次,上下文可忽略。

## 改动面
新增 `coord.sh` + 一个 watchdog plist；现有脚本插 `register/heartbeat/checkpoint` 调用。dispatch/triage 角色由启动器传入。

## 验收
- 实例崩 → ≤3min watchdog 标 orphaned → triage 接管续跑、不重复查证。
- dispatch 实例全程不 adopt 他单。
- 跨仓 worktree 经 `tasks/*.json` 可被任意实例枚举。
- checkpoint 0 网络、心跳 0 上下文。
