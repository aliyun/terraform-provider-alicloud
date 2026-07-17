# 交付链路:cloudspec

> **红线:任何改动一律新分支 + CR/MR 评审,master 只接已评审合入。禁止直接 push master,禁止从 master 拉零 diff 空 CR 直发正式——没有可评审 diff 就不上线。**

分诊后**自己实现并交付**(不转上游)走这条。全程 `a1` CLI,从建需求到上正式。
应用 = `cloudspec`(阿里云 API MCP SERVER,owner 原根)。

## 路由(命中即走本链路,默认项目 2124589 + app 260634)
满足任一即落本链路,无须再问"哪个项目":
- 当前在 **cloudspec 仓库**目录(cwd 含该 repo),或
- 明确给 **cloudspec / OpenAPI MCP Server / API MCP**(及 ListApis、RunIaC 等 cloudspec 内置能力)提需求
→ 需求建到 https://project.aone.alibaba-inc.com/v2/project/2124589/req,基于 app https://cd.aone.alibaba-inc.com/unite/micro/cr/app/260634/list 开发,自指派工号 **320687**(辰羿)。
别的应用 / 别处的需求池 → 不走本文件。

固定坐标(本应用,权威值见 config/pools.json mcp_server.apps):
- 项目空间(建需求):**2124589**
- 应用:**260634**(repo cloudspec-mcp/cloudspec)
- 发布流水线:**65 日常 / 420 预发 / 67 正式**

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
a1 app link 260634
a1 app cr create "3-5行摘要" --branch <suffix> --workitem-ids <需求ID>
```
分支名 a1 自动加 `feature/<date>_<n>_<suffix>_1`。记下 CR id + 分支。

**注意**:open_jarvis 账号在 app 260634 可能无 committer 身份,自动建 CR 报 `600001 无效用户` 时需人工建 CR 或先加权限。

## 3. worktree 开发
```
git fetch origin <branch>
git worktree add .claude/worktrees/<name> <branch>
```
进 worktree 改代码、commit、`git push origin HEAD`。
若为修 bug:补/改一个会因该 bug 失败的用例锁定回归(无可测则在 CR 说明为何),用例随 CR 一并评审。

构建验证(无 mvnw,用系统 mvn):
```
mvn -q -DskipTests package              # 编译
mvn -q -Dtest=<TestClass> test           # 跑单个测试类
```

## 4. 预发
```
a1 app link 260634            # 坑②:换目录后必重绑
a1 app cr submit <cr> --pipeline-id 420
a1 app pipeline status --pipeline-id 420
```
监控部署:轮询 `pipeline status`,看「预发部署」阶段到 SUCCESS(后面「预发验证 CR_PUBLISH_RELEASE」是人工发布卡点,WAITING 不影响已部署)。
**预发后停下,等用户验证 + 明确反馈再进正式。** 别预发绿了就自动发线上——正式不可逆,放行权在用户。

## 5. 正式(用户预发验证通过 + 反馈 OK 后)
```
a1 app cr submit <cr> --pipeline-id 67
a1 app pipeline status --pipeline-id 67
```
**收口确认**:`pipeline status` 整体 SUCCESS → CR 自动置 FINISH。完成后按需建 MR 评审。

## 6. 清理(发布完成后)
```
git worktree remove .claude/worktrees/<name> --force
git branch -D <branch>
```
CR FINISH + 正式 SUCCESS 后才清。

## 坑(逐条踩过)
1. **必填字段**:建 req 前先 `field options` 查 version,带 `--version` + `--cfs 100340=日期`。
2. **worktree 丢绑定**:换工作目录后 `app cr submit` 报 `no app linked`,先 `a1 app link 260634` 重绑。
3. **预发 pipeline 是 420 不是 66**:66 是 aliyun-automation-agent 的,别混;`a1 app pipeline list --app 260634` 核实。
4. **open_jarvis 身份限制**:open_jarvis 在此 app 可能无 committer,自动建 CR 会报错,需人工建或加权限。
5. **无 mvnw**:本仓库没有 Maven Wrapper,用系统 `mvn` 命令。
6. **CI 全模块受内网镜像影响**:全模块 `mvn test` 可能因内网坏包失败,单类独立编译+跑通即可。
7. **共享文件撞兄弟 CR**:submit 报 `pipeline_submit_failed` 或代码冲突预检失败时，**永久禁止** `app pipeline exit-cr`、等价的 `app pipeline quit` 以及定向的 `app cr quit`，`bin/a1id` 与 PreToolUse 已在真实 a1 前硬拒绝。立即停止后续发布动作，反馈“部署失败”，附上当前 CR、流水线及冲突信息，等待人工给出解决方案；禁止 Jarvis 自行查找或 merge 兄弟分支，也禁止自行重 submit。任何 CR/分支退出动作都交给人工处理。

该护栏的威胁模型是防止 Jarvis 在正常工具入口误操作或绕过 wrapper；它不声称隔离同一 UID 下的恶意本地代码——后者可直接读取本机凭据并调用绝对路径 a1，不属于仓库内 guard 能提供的密码学边界。

## 交付链路目录
| 应用 | 项目 | app | 预发/正式 pipeline | 链路文件 |
|------|------|-----|--------------------|----------|
| cloudspec | 2124589 | 260634 | 420 / 67 | 本文件 |
| aliyun-automation-agent | 2124589 | 283346 | 66 / 67 | delivery-aliyun-automation-agent.md |
