---
name: provider-resource-dev
description: Use when DEVELOPING, DIAGNOSING, or FIXING an alicloud Terraform provider resource — either (a) NEW resource end-to-end (Terraform 资源名解析 → 镇元/Cloudspec resourceTypeCode 查证 → acube 映射/生成 → 生成代码 vs 手写代码 diff → 手改 → acc 验收 → PR), OR (b) **修/改一个非自动化生成（hand-written 或历史遗留 generated-but-mutated）的既有资源** — 补属性、修 bug、补重试、修 schema drift、加 import 支持等。触发场景：客户/Aone 需求要 接入/支持 一个资源(e.g. alicloud_oss_bucket_inventory)；生成资源空/缺；cloudspec terraform 无资源；要拿 镇元 spec 落 provider 代码；**或既有资源出现 bug（错误码未重试 / attribute 缺失 / CRUD 不对齐 API / import 断链 等）**。NOT for 现有 PR 评审(用 terraform-pr-review)或简单 是否支持 查询(用 aone-triage 查证)。 NOT for structured release SOP with mandatory Aone gap-analysis and PR-merge governance — use terraform-provider-release.
---

# Provider 资源开发全流程

**两种模式，共用后半段（手改 → 回归用例 → 验收 → PR）**：
- **新增资源模式**：走完全部 8 步，`getTerraformResourceSpec` / acube 生成 / 生成-手写 diff 是必经。
- **修复/补丁模式**（非自动化生成资源的 bug fix / 属性补齐）：跳过步骤 1–5（不走生成），直接 **步骤 6（手改）→ 6.5（补/改回归用例）→ 7（AccTest）→ 8（PR）**；且 **6.5 是硬门**：修 bug 必须补/改一个"未打 patch 前会 fail、打 patch 后 pass"的锁定回归用例，随 PR 一并交付；确实无法测（如概率事件、需实体网关、外部依赖）则在 PR body 明写原因 + 已做的替代验证（静态查证 / 逻辑推演 / 同类资源对比）。

从客户需求到 provider PR。真源=镇元(Zhenyuan)`ResourceTypeSchema`;生成器吃镇元料出一整套(resource.go + _test.go + service + provider.go 注册行 + website 文档 markdown),但它可能空生成/部分生成;**只信 acc 测**。

## 工具/路径
- 工作区一律通过 `bootstrap/workspace.sh dir <key>` 解析;provider key=`terraform_provider`;acube key=`acube`;生成器 key=`terraform_generator_v4`。
- cspec 仓 `cloudspec-model/<Product>_pop_*`;provider upstream=aliyun;Jarvis 提交 GitHub PR/评论/推分支必须使用 `JARVIS_GITHUB_TOKEN` 对应的 `api-tool-agent` 身份,head=`api-tool-agent:<branch>`;acube/terraform-generator-v4 见 `config/workspaces.json`。
- Acube 在线生成工具: `tools/acube_terraform_generate.py`(tools/ 是 repo 顶层的 Python 工具目录,不在 skill scripts/ 内,便于跨会话共享 + 独立测试)。
- 生成差异/语义检查工具: `tools/terraform_generated_diff.py`(同上)。
- **错误码语义查证**(客户 acc/apply 报错、retry 白名单决策):读 `.claude/skills/aone-triage/references/aliyun-error-code-lookup.md`(跨 skill 复用,给定 product+code 出 HTTP/中英 message/官方 retry 建议/相邻错误码)。
- **镇元查证与路由分支**(诊断资源在哪一层缺、按决策树选执行分支):读 `references/zhenyuan-verification.md`(跨 skill 单点维护:aone-triage tf-customer 路由与本 skill 资源开发都读它)。
- **开发前先读 `<playground>/<product>/KNOWLEDGE.md`(存在即读)**:jarvis 蒸馏的产品级可执行知识(命名/参数 quirk/生命周期/API 行为/报错→原因→解法五节);playground 路径解析走 `bootstrap/workspace.sh dir tf_playground` 或 env `JARVIS_TF_PLAYGROUND`,文件不存在跳过。契约见 `.claude/skills/tf-customer-probe/references/knowledge-distillation.md`。
- 编码交 terraform-rd 子代理,改文件先 worktree,acc 由 terraform-qa 验过才交。
- **PR 提交后**按 `.claude/skills/tf-customer-probe/references/knowledge-distillation.md` 契约把本次开发学到的产品级事实(API 行为差异、schema 陷阱、必须的重试码等)蒸馏进 `<playground>/<product>/KNOWLEDGE.md`(触发点③provider-resource-dev 完成开发后),来源锚点写 upstream PR URL + provider 源码行号。

## Aone 分单与同步
非自动化生成链路、需要 Jarvis 内部研发处理的 Terraform Provider 资源,必须创建或复用 **terraform-alicloud** 内部研发单:
- 项目: `tf_provider` / `528766`;指派按 aone-triage skill
  `references/tf-customer-request-routing.md` 的分支契约执行，不把源客户主单 owner 和研发单
  owner 混为一人，也不把 assignee 写成数字 worker `WORKER_1782379562571`。
- 与 Terraform-客户需求池的客户主单双向关联；无客户主单的 adhoc 场景按
  `loops/adhoc-intake.md` 走。
- 主要研发进展、生成/手改差异、验证细节、PR/CI/验收信息优先同步到内部研发单。
- 客户主单只同步关键节点摘要:已转内部单、发现镇元/Cloudspec/API 卡点、资源模型问题、需要客户感知的决策或阻塞。

### G / 紧急普通 D 的双 owner 契约

本契约覆盖 G Provider 全局改造，以及紧急普通 D（纯 datasource；或 CloudSpec 结构 OK +
手写 Provider）：

- **源客户主单 assignee 保持新山（521957）**；**528766 研发关联单 assignee 固定过载（484483）**，
  由 **Jarvis/TerraformRD claim 并尝试修复**，不等待新山接手研发。
- 写前 point-read relation、同题 528766、assignee 与 lease。**healthy existing claim 不抢占**；
  无健康 claim 时对**同题单原地复用并幂等改派过载**，但严格顺序是同一 terraform-rd Task
  先进入 route-finalizer phase 并 fail-closed claim，成功后才改派、补缺失 relation 并切 dev；
  claim 失败立即停止，禁止 relation/update/wrap。不存在同题单才 create 一次并 claim，禁止
  重复建单。
- **relation/assignee/status 不是完成信号**；**无 PR/CI/QA 完成信号必须继续开发**，按本
  skill 完成 TDD、Provider 修改、CI 与 QA。
- **build/test/CI 失败只走 RD ↔ QA 修复闭环**，QA fail 也回 RD 修复重验，**不得转交新山**。
- `missing_capability / retry exhausted` 进入 **blocked / SUSPENDED**，源单仍是新山、研发单
  仍是过载，保留失败证据并 release，**不得 finish**。
- route materialization、必要改派、relation 与 claim 只由 TerraformRD 控制面幂等执行；
  PD/QA 不外写。源单由 executor、实际 claim 的 528766 由最终 RD finalizer 分别完成，
  **源单与实际 claim 的 528766 各自最多一次聚合 bookend**，开发阶段不发阶段评论。
  **PR 未合并只 release**，不 finish。
- 源工单禁令不约束按既有契约由内部链承接的 528766；非紧急 D 与 I 的 Provider docs
  紧急兜底腿保持既有内部 claim/bookend。G/紧急普通 D hard gate 只新增双 owner、先 claim
  后 dev 与不可观察语义；D-临钧/A/F/H 和 E 路径边界不变。

- **分支 I — CloudSpec 文档文本 metadata**：resource/property/operation description、字段解释、
  NOTE 与枚举文案，且不改变字段集合、类型、约束或 CRUD。I 不进入本 skill 的 CloudSpec
  开发路径；finalizer 创建或复用 `upstream.cloudspec_docs_quality`（2169561，念依 373108，
  `submit_only`）。若公开 Provider docs 也错误，保留独立 528766 紧急兜底腿并分池防重；
  一个池已有 relation 不能抑制另一个池的缺失补建。
- **分支 E — CloudSpec 结构 metadata 原主单自闭环**：只处理字段集合、类型、约束、CRUD、
  operationMapping 与生命周期等结构合同。PD 返回 `requested_external_actions: []`、
  `next=terraform-rd/dev`；RD 按
  `terraform-provider-release/references/cloudspec-pre-resource-loop.md` 使用 CloudSpec skills + AMP
  修到 build/check/pre Meta 收敛，不创建 2165097。pre 未收敛不得触发 Acube。
- E 收敛后必须 **E → D-临钧**：已有正确 relation/taskId/aoneId 时只查询/复用，否则由
  finalizer 的 `single-writer` 通过 Acube `createBuildTaskV2` 自动创建或复用 528766，指派
  临钧（429768）。PD/QA 不外写；不得由 E 直接执行 Provider PR/CI/ACC，不得在 E 完成后直接
  release/idle。
- 只允许分支 E 进入此转换；不得泛化到 A/F/G/H/I、纯 datasource 或纯手写 Provider-only bug。
  普通分支 D 仍按本 skill 的 Provider 开发、PR CI 和远程 ACC 流程执行。CloudSpec
  prod/online、master/main merge/push 与正式发布始终是人工硬门，不得 finish。
- AMP 登录、SSH、模型仓权限、pre 发布或 Acube 能力缺失时返回
  `missing_capability` / `blocked`，不得回退外部承接人或个人身份。

## 步骤
1. **查证 Terraform ↔ Cloudspec 身份** — OpenAPI + provider 源码确认缺;`getTerraformResourceSpec` 只看映射,不代表实现,且找不到可能返无关资源,别信。
2. **查 acube resourceTypeCode** — 通过 Terraform 资源名推 product/resourceCode 后,先 get,再 list 降级;get 只有 `SUCCESS` 无 `data` 就按未命中处理。
3. **镇元建模/发布判断** — 如果 get/list 都找不到,说明资源未进预发/线上 resourceTypeCode;先回复 Aone 推动镇元发布或确认别名,除非本地已有 cspec 分支可继续验证。
4. **生成** — 标准入口走 Acube `createLocalBuildTask`:先 `resourceTypeCode/get` 取料,再 `createMapping`,再 `createLocalBuildTask`,用 `tools/acube_terraform_generate.py` 落盘 raw JSON/logs/files/generated/summary。只有 Acube 不可用或需验证本地 cspec 分支时,才 fallback 到本地 `cloudspec terraform -r <terraform_resource> -e pre -o <dir>`;报 `no that resource` 时用 `<CloudspecResourceCode>` 重试并记录 partial output。
5. **生成 vs 手写 diff** — 用 `tools/terraform_generated_diff.py` 看 Acube generated 与手写分支差异,先判断缺主体/缺测试/缺文档/缺 provider 注册/缺 service,再看 `resourceNotExistCondition` 等语义风险,最后手改。
6. **手改生成缺陷 / bug fix** — OSS 等 XML 产品:`client.Do("Oss",xmlParam(...))` 非 Roa(治 `<` 解析错);PUT body 按 schema **固定元素序**(治 MalformedXML);update 删-再-PUT 绕 AlreadyExists。**修复模式**下,查证根因(OpenAPI 错误码/schema/CRUD 语义)后照同 package 已有正确写法照抄,不 refactor 无关代码。
6.5. **补/改回归用例(修复模式硬门)** — 修 bug 时必须在同 PR 补/改一个"未打 patch 前会 fail、打 patch 后 pass"的用例(unit 或 acc,能锁定回归即可),随 CR/PR 一并评审。选型:
   - 若可用小并发/构造错误码/mock RPC 稳定复现,写 unit 或 acc 用例。
   - 若是概率事件(如服务端锁冲突、限流),可用同一 Config 里 declare N 个 resource 让 Terraform 并发 apply 触发。参照 PR 9916 `TestAccAliCloudESARoutineRoute_lockRetry` 模式:同一 SiteId+RoutineName 下起 N 个 route,未打 patch 会 LockFailed 挂;打 patch 后 apply/destroy 全过。
   - 真无法测(需实体网关/外部依赖/无稳定复现路径)才允许豁免:在 PR body 明写"无法用例复现:<原因>",并列出替代验证(同类资源对比 / 静态查证 / manual repro log)。审阅人可拒。
7. **验收** — 真实 AccTest 优先用 `invoke-terraform-acc-test-remote` 走 ACube/FC 远程执行,避免长时间占用本机;本地只跑 `go test ./alicloud -run '^$'`、小单测、lint、示例 `terraform validate` 等轻量检查。远程 AccTest 过 create+update+import 才算数(**修复模式**至少要过 6.5 补的回归用例 + 原有主用例);跨账号/企业账号资源要隔离 ambient `ALICLOUD_ACCESS_KEY`/`ALICLOUD_SECRET_KEY`,显式声明测试需要的多把 AK 环境变量,并用 STS/CLI 验证每把 AK 的 caller account,但任何文档/评论/示例都不能泄露真实 AK/SK。
7.5. **TestingCoverageRate CI 门(改既有资源必过)** — PR 动了 `alicloud/resource_alicloud_<name>.go` 或其 `_test.go`,CI 就按**该资源全量 schema** 跑 `scripts/testing/testing_coverage_rate_check.go`;本地先复现:`go run scripts/testing/testing_coverage_rate_check.go -resource=alicloud_<name>`。三层检查全过才绿:
   - ① **must-set**:每个 active `Optional/Required` 属性(仅豁免 `dry_run`/`Deprecated`/`Removed`)至少在一个 test 的 config map 出现;
   - ② **ignore 合法**:`ForceNew` 属性禁入 `ImportStateVerifyIgnore`;数组禁含非 active-schema 项(拼错遗留常见)。**数组通常在文件内所有带 import 的 test 重复出现,改必须全改**(工具取全文件并集);
   - ③ **modify**:非 ForceNew 可修改属性需在 step 间取过不同值;已在 ignore 数组或文档标注 immutable 的豁免。
   缺口处理次序(**禁占位值/空串凑覆盖**,对齐 CI 失败禁绕过纪律):
   - 能真实测的补真实 test;下发型/ForceNew/无回填属性集中放一个 **create-only、无 import step** 的属性覆盖 test(化解「ForceNew 禁入 ignore」vs「set 后 import verify 必 diff」的死锁);
   - 无回填的非 ForceNew 属性(如公网 `port`)加 ignore 数组——import 本就无法 verify,语义正当,顺带豁免 modify;
   - 常用真实化手段:`kms_encrypted_password` → `alicloud_kms_ciphertext` 真密文(先例 `resource_alicloud_rds_account_test.go`);`coupon_no` → 官方"不用券"占位值 `youhuiquan_promotion_option_id_for_blank`;`private_ip` → dependence 自建 vpc/vswitch 固定 cidr;`capacity` → 取与 `instance_class` 规格一致值;企业版特性(如 TDE)→ 扩展既有 amber/企业版规格 test;
   - **明确补不了的**(需外部前置:真实备份/专属集群/既有全球实例等)**不硬凑**,在关联工单逐属性列明原因(区分「不可测」与「可补但需扩 scope,建议 follow-up」),PR 里给 maintainer 一段简短英文说明。
8. **PR** — 提交走 `bootstrap/github-identity.sh commit -m "..."`(而非裸 `git commit`)→ `bootstrap/github-identity.sh check` → `bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud HEAD <branch>` → `bootstrap/github-identity.sh gh pr create --repo aliyun/terraform-provider-alicloud --head api-tool-agent:<branch>`;带 resource+test+service+provider注册+website 文档;无 AI 署名。缺 `JARVIS_GITHUB_TOKEN` 或登录名不是 `api-tool-agent` 时阻断并升级,禁止回退个人账号或 ambient git 凭据。
   - **commit 作者硬门(CLA)**:CLA-assistant 按 **commit 作者邮箱**核验,不是 push token 也不是 PR opener。裸 `git commit` 会用本地默认身份(如 `jarvis@jarvis.local`)→ `license/cla` 必挂。子代理提交一律走 `github-identity.sh commit`(自动署名 `api-tool-agent <cloudspec_bot@alibaba-inc.com>`);已用裸 commit 的用 `bootstrap/github-identity.sh commit --amend --no-edit` 重署名后再 force-push。`push` 会对 tip 作者不符时 WARN。
   - **force-push 自有 fork PR-head 是预授权动作**:为满足单提交 CI 门禁而 squash/rebase/重署名后 force-push 到 `api-tool-agent:<branch>`（`... push api-tool-agent/terraform-provider-alicloud +<ref> <branch>`）属 `autonomy.md` 的 `fork_push`,headless 下**直接执行,不 SUSPEND/escalate/等工单放行**。仅限自有 fork PR-head;force-push 上游 `aliyun/…` 或任何 master = `release_prod` 人工硬门。

## Terraform 资源名解析
优先查 Acube Terraform 映射接口,不要维护固定映射表:

```bash
curl -s "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=<terraform_resource>" \
  -H "accept: */*"
```

读取 `data.terraformResourceSpecModel.namespace` / `data.terraformResourceSpecModel.resourceTypeCode` 作为 product/resourceCode。接口查不到时,再把 `alicloud_【产品名下划线】_【资源名下划线】` → product/resourceCode 各自转 PascalCase 作为候选,并进入下方 resourceTypeCode get/list 查证;边界不确定时结合客户描述/Next API/OpenAPI product 先确定 product,不要凭命名猜死。

## acube resourceTypeCode 查证

完整配方(schema get/list、released 复核、覆盖度 V2、sanity check)**单点维护在
[references/zhenyuan-verification.md](references/zhenyuan-verification.md) 的镇元 schema 查证节**,照它执行。本节只留守则:

- get 返回 `SUCCESS` 但无 `data` = 未命中;降级查 `resourceTypeCode/list`(released=false),列表也没有才判定未发布/未同步。
- **env 合法枚举只有 `pre` / `online`**——`env=prod` 接口不报错但**静默返回空 data**(实测),会误判「镇元无此资源」,勿用。
- sanity check:用同产品已存在资源(如 `Handshake`)反查一次;它有 data 说明接口链路正常,目标 resourceCode 无 data 才是真缺失。

## 接口示例(可复制)
```bash
# 预发镇元是否有模型(developing 仅 pre)
curl -s "https://pre-amp-apispec.aliyun-inc.com/api/v1/open_api/api_groups/describe_namespace_resource_type_list?product=OSS" \
 | python3 -c 'import json,sys;[print(i["id"],i["name"],i["status"]) for i in json.load(sys.stdin)["data"] if "Invent" in i["name"]]'
# ① 注册映射(GET;prod=acube.aliyun-inc.com,pre=pre-acube;未发线上加 ignoreVerify)
curl -s "https://pre-acube.aliyun-inc.com/api/v1/terraform/resource/mapping/createMapping?name=alicloud_oss_bucket_inventory&namespace=OSS&resourceCode=BucketInventory&isSpec=true&ignoreVerify=true"
# ② 在线打包出制品(POST form,非 JSON;env=pre 需 acube env 透传补丁已合)
curl -s -X POST "https://pre-acube.aliyun-inc.com/api/v1/terraform_vendor_build/createLocalBuildTask" \
 --data-urlencode namespace=OSS --data-urlencode resourceTypeCode=BucketInventory --data-urlencode env=pre
# ③ 本地等价(要本地 cspec 输入)
cloudspec terraform -r alicloud_oss_bucket_inventory -e pre -o ./gen_out
```
namespace/code 用 PascalCase;比对现成:`getTerraformResourceSpec?terraformResourceType=alicloud_oss_bucket` → OSS/Bucket。

优先用封装工具,避免漏保存 raw 证据:
```bash
python3 tools/acube_terraform_generate.py \
  --resource alicloud_resource_manager_handshake_acceptance \
  --env pre \
  --out-dir /tmp/acube-handshake-acceptance
```
若 Python urllib 报内部证书链校验失败,在确认目标是内部 Acube 预发/线上域名后加 `--insecure`;默认仍校验证书。

## 生成差异工具
对比生成目录和 provider 手写分支:
```bash
provider_repo="$(bootstrap/workspace.sh dir terraform_provider)"
python3 tools/terraform_generated_diff.py \
  --resource alicloud_resource_manager_handshake_acceptance \
  --acube-dir /tmp/acube-handshake-acceptance \
  --provider-repo "$provider_repo" \
  --handwritten-ref origin/f/resource_manager_handshake_acceptance
```

默认比较:
- `alicloud/resource_<resource>.go`
- `alicloud/resource_<resource>_test.go`
- `website/docs/r/<resource-without-alicloud>.html.markdown`
- Acube 生成目录中的 `alicloud/service_*.go`

输出看三块:
- `Structured summary`: 主体/测试/文档/provider/service 的 generated/handwritten 覆盖情况。
- `Semantic checks`: 从 `resourceTypeCode/get` raw JSON 检查生成代码语义;例如 `resourceNotExistCondition.notExistCheckTargetValueType=assertNotEqual` 时,生成代码应在 `!= <target>` 时返回 NotFound,若生成成 `== <target>` 要转 generator 问题。
- `Only in handwritten`: 生成器缺主体/测试/文档,常见于 resourceTypeCode 未发布或 partial generation。
- `Only in generated`: 生成器多出的文件,确认是否需要保留。
- `Diff`: 共同文件的 unified diff,左边 generated,右边 handwritten。

需要看 provider 注册时加 `--include-provider`;大 provider 仓可能很吵,默认不打开。

## Aone 回复模板(资源未发布)
```text
已按 Terraform 资源名 <terraform_resource> 推导 Cloudspec 查询参数:
- product=<Product>
- resourceCode=<ResourceCode>

验证结果:
1. resourceTypeCode/get 在预发检索 <ResourceCode>,接口返回 SUCCESS,但无 data。
2. resourceTypeCode/get 在线上域名检索 <ResourceCode>,接口返回 SUCCESS,但无 data。
3. resourceTypeCode/list 查询 <Product> 产品资源列表,相关命中仅有 <近似资源>,未命中 <ResourceCode>。
4. 用已存在资源 <近似资源> 反查,get 接口可返回完整 data,说明接口链路正常。

结论:<terraform_resource> 对应的 Cloudspec 资源 <ResourceCode> 尚未发布到 acube 预发/线上可检索的 resourceTypeCode 中,生成器无法按该资源完成正常检索/生成。建议先推动镇元定义发布到预发/线上,或确认是否需要配置 resourceCode/映射别名;待 acube 可检索到 <ResourceCode> 后再继续 Terraform 生成与差异对比。
```

## 常见坑
- 生成空/`获取spec数据失败` → env 没到预发镇元(查 createLocalBuildTask Env 是否 online)。
- Acube 生成的 NotFound 条件与 Cloudspec `resourceNotExistCondition` 相反 → 优先查 `resourceTypeCode/get` raw JSON;若 Cloudspec 为 `assertNotEqual` 而生成 Go 为 `==`,按 generator bug 转 `api_toolkit` 池并关联源 Aone。
- `cloudspec terraform -r alicloud_x` 报 `no that resource` → 先查 resourceTypeCode;本地 cspec 分支存在时可尝试 `-r <ResourceCode>`。
- get 接口 `SUCCESS` 但无 `data` → 未命中,不是成功获取定义。
- 产品列表只命中近似资源(如 `Handshake`),没命中目标(如 `HandshakeAcceptance`) → 目标未发布/未同步。
- `MalformedXML` → body map 随机序;`invalid character '<'` → 用了 Roa(JSON)发 XML。
- 直接发生成版不验收 → 必踩 OSS XML;acc 不过别交付。
- 长耗时 AccTest 在本地裸跑 → 占机器且易被 ambient AK/Profile 干扰;优先用 `invoke-terraform-acc-test-remote` 打包本地代码交给 ACube/FC 跑,再读 `run.log` / `tf-debug.log` 判定。
- ResourceManager handshake/成员账号类资源:若存在删除/移除关系 API,provider Delete 必须实现真实删除并校验幂等;AccTest 前置清理只用于清历史脏关系,不能替代资源 Delete。清理或删除后轮询 NotFound/终态并等待一致性后再 invite。
- 多 AK 测试脚本不要假设 `aliyun` 环境变量一定覆盖默认 profile;优先显式传 `--mode AK --access-key-id ... --access-key-secret ... --region ... --endpoint ...`,避免把 AK1/AK2 都打到同一个默认账号。
- 可复制 Example 必须能 `terraform init/validate`:包含 `required_providers`,跨账号资源用 provider alias 区分管理账号/受邀账号;真实 AK/SK 不写进文档、评论或仓库,只用 sensitive variables/env。
- **Optional + Computed 字段的 `HasChange` drift 陷阱**:当 schema 里字段是 `Optional: true, Computed: true`,`Read` 又调 `d.Set(<字段>, <API 值>)` 回填,那么二次 apply 时即使客户 tf config 里没设该字段,state 也会带 API 回填的值——`d.HasChange(<字段>)` 会返回 `true`(config 空 vs state 非空),Update 逻辑若无 guard 直接把 `d.Get(<字段>)` 空字符串塞给下游 API 就会踩 Duplicate/InvalidParameter 类错误。**修法**:Update 里所有 `HasChange(<Optional+Computed 字段>)` 分支必须叠加 `if v := d.Get(<字段>).(string); v != ""` 过滤,只有客户显式设了非空值才调下游 API;若下游 API 对同值请求返回 `Duplicate` 类错误,还要在错误处理里 `IsExpectedErrors(err, []string{"...Duplicate"})` tolerate 兜底(视为 no-op)。先例:kvstore_instance Modify `private_connection_port` 因 Computed drift 二次 apply 报 Duplicate。**规则**:审自己写的 Update 逻辑时,凡 schema 是 Optional+Computed 的字段,check 三处:① Read 是否有 `d.Set` 回填;② Update 是否只用 `HasChange` 不 filter 空;③ 下游 API 对"相同值"的错误码是否被 tolerate。

## PR 侧规则(对外)

- **PR title 必须符合 upstream CI 检查**:aliyun/terraform-provider-alicloud 的 CI 检查 PR title 前缀模式为 `resource/alicloud_<resource>:` / `data-source/alicloud_<resource>:` / `docs:` / `testcase:` 等,**不接受 conventional commit 格式**(`fix(xxx): ...` / `feat(xxx): ...`)。提 PR 前 title 命中格式错会被 CI 阻断;如需改 title,`gh pr edit --title` 后可能需 close/reopen PR 触发 CI 重跑(admin 权限才能 `gh run rerun`)。先例:PR 9937 因 `fix(resource_...)` 前缀被 CI 卡住。
- **对外产物 sanitize**:GitHub 公开仓 PR title/body、`git commit` message、code comments 严禁内部信息(Aone URL、客户名、账号 UID、实例 ID、RequestId、内部人员花名工号引用等)。完整清单见 CLAUDE.md 工作纪律 #5 + `terraform-provider-release` SKILL Step 11.1;**push 前自查(两步)**:先 `bash <jarvis仓>/bootstrap/pre-push-sanitize.sh`(禁品正则真源),再 `git log -p origin/master..HEAD` 通读 diff 与 commit message 兜底客户名类脚本无法穷举项,命中禁品即卡住修。

## 红线
不碰 master;无可评审 diff 不空发;对外无 AI 署名;Aone 唯一真源(进展 sync/完工 done);**修 bug 无回归用例又不写豁免原因 = 空发**。
