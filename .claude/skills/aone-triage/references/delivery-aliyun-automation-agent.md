# 交付链路:aliyun-automation-agent

> **红线:任何改动一律新分支 + CR/MR 评审,master 只接已评审合入。禁止直接 push master,禁止从 master 拉零 diff 空 CR 直发正式——没有可评审 diff 就不上线。**

分诊后**自己实现并交付**(不转上游)走这条。全程 `a1` CLI,从建需求到上正式。
应用 = `aliyun-automation-agent`(包 `com.aliyun.terraformai`,PlayGround 在 `templates/playground.html`)。

## 路由(命中即走本链路,默认项目 2124589 + app 283346)
满足任一即落本链路,无须再问"哪个项目":
- 当前在 **aliyun-automation-agent 仓库**目录(cwd 含该 repo),或
- 明确给 **Agent门户 / AgentRuntime / aliyun-automation-agent**(及 PlayGround、Agent应用)提需求
→ 需求建到 https://project.aone.alibaba-inc.com/v2/project/2124589/req,基于 app https://cd.aone.alibaba-inc.com/unite/micro/cr/app/283346/list 开发,自指派工号 **320687**(辰羿)。
别的应用 / 别处的需求池 → 不走本文件。

固定坐标(本应用,权威值见 config/pools.json agent_portal 节):
- 项目空间(建需求):**2124589**
- 应用:**283346**(repo opensource-tools/aliyun-automation-agent)
- 发布流水线:**65 日常 / 66 预发 / 67 正式**

每步写操作前先报当前用户:`a1 auth whoami`;CR/部署命令依赖绑定,见坑 ②。

## 1. 建需求
```
a1 project workitem create --project 2124589 --category req \
  --title "..." --assignee <empId> --version <id> --cfs 100340=YYYY-MM-DD --body "..."
```
- 本项目 **version 与期望日期(字段 100340)必填**,缺了报 `400 【版本】不能为空,【期望日期】不能为空`。
- 取枚举:`workitem field options version --project 2124589`(选最新「正式」版本);期望日期是普通 date,填到该版本末。
- 自指派直接给工号(320687=辰羿)。

## 2. 建变更(CR),自动建分支+关联需求
```
a1 app link 283346
a1 app cr create "3-5行摘要" --branch <suffix> --workitem-ids <需求ID>
```
分支名 a1 自动加 `feature/<date>_<n>_<suffix>_1`。记下 CR id + 分支。

## 3. worktree 开发
```
git fetch origin <branch>
git worktree add .claude/worktrees/<name> <branch>
```
进 worktree 改代码、commit(co-author 行)、`git push origin HEAD`。
若为修 bug:补/改一个会因该 bug 失败的用例锁定回归(无可测则在 CR 说明为何),用例随 CR 一并评审。

## 4. 预发
```
a1 app link 283346            # 坑②:换目录后必重绑
a1 app cr submit <cr> --pipeline-id 66
a1 app pipeline status --pipeline-id 66
```
监控部署:轮询 `pipeline status`,看「预发部署」阶段到 SUCCESS(后面「预发验证 CR_PUBLISH_RELEASE」是人工发布卡点,WAITING 不影响已部署)。
验证预发:URL `https://pre-agent.aliyun-inc.com/playground?dev=1`(`?dev=1` 跳 MCP 配置门)。匿名 curl 撞 SSO,前端验证需 BUC 登录;无登录时退而 `grep` 部署产物里的新标记(如 CSS class)确认已上线。
**预发后停下,等用户验证 + 明确反馈再进正式。** 别预发绿了就自动发线上——正式不可逆,放行权在用户。

## 5. 正式(用户预发验证通过 + 反馈 OK 后)
```
a1 app cr submit <cr> --pipeline-id 67
a1 app pipeline status --pipeline-id 67
```
比预发多 CM 卡点,见坑⑤;`submit_publish_component` 若 `canAutoSubmit:false`,CLI 推不动,要用户去发布页点「提交发布单」。
**收口确认**:`pipeline status` 整体 SUCCESS(正式部署 / 写基线 / 线上测试全绿)→ CR **自动置 FINISH**(发布即闭单,不用手填状态,坑⑦的「已发布」枚举不存在)。完成后按需建 MR 评审。

## 6. 清理(发布完成后)
```
# 会话在 worktree 内:先 ExitWorktree(action: keep)回主目录,再删
git worktree remove .claude/worktrees/<name> --force
git branch -D <branch>            # 已合主干,本地分支可删
```
CR FINISH + 正式 SUCCESS 后才清;worktree 干净(已 push)。手建的 worktree ExitWorktree 不代删,要手动 `git worktree remove`。

## 坑(逐条踩过)
1. **必填字段**:建 req 前先 `field options` 查 version,带 `--version` + `--cfs 100340=日期`。
2. **worktree 丢绑定**:换工作目录后 `app cr submit` 报 `no app linked`,先 `a1 app link 283346` 重绑。
3. **预发/正式 pipeline 固定** 66/67,别凭名字猜,`a1 app pipeline list --app 283346` 核。
4. **监控收口**:用 poll-until 看「X 部署 [SUCCESS]」收尾,空日志/RUNNING 不算结束;阻塞在人工卡点别死等。
5. **CM 管控人工卡点**:正式准入 CheckList 含安全扫描(自动)+ CM 管控(WAITING),CLI 推不动,要去 BG 变更审批页人工通过,过后自动续跑。
6. **共享文件撞兄弟 CR**:submit 报 `pipeline_submit_failed` + 代码冲突预检 → 多半是另一在飞 CR 改了同一文件(`playground.html` 高发,多需求并行)。`git merge origin/<兄弟分支>` 合入共存、解冲突、push、重 submit。先 `git log --oneline --all -- <文件>` 找兄弟分支。
7. **状态枚举无"已发布"**:`workitem update --status` 用 `验收通过`(或 `已发布待需求方验收`),写"已发布"会报 unsupported,枚举见报错列表。
8. **监控 grep 别太宽**:盯部署用紧 grep,`部署.*(SUCCESS|FAIL)` 会被 WAITING 行里的"预发部署"误命中假收口;锚到 `预发部署 \[SUCCESS\]` 整阶段。

## 交付链路目录
| 应用 | 项目 | app | 预发/正式 pipeline | 链路文件 |
|------|------|-----|--------------------|----------|
| aliyun-automation-agent | 2124589 | 283346 | 66 / 67 | 本文件 |
| cloudspec | 2124589 | 260634 | 420 / 67 | delivery-cloudspec.md |
