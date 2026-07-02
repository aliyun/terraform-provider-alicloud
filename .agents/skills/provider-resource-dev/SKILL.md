---
name: provider-resource-dev
description: Use when DEVELOPING, DIAGNOSING, or FIXING an alicloud Terraform provider resource — either (a) NEW resource end-to-end (Terraform 资源名解析 → 镇元/Cloudspec resourceTypeCode 查证 → acube 映射/生成 → 生成代码 vs 手写代码 diff → 手改 → acc 验收 → PR), OR (b) **修/改一个非自动化生成（hand-written 或历史遗留 generated-but-mutated）的既有资源** — 补属性、修 bug、补重试、修 schema drift、加 import 支持等。触发场景：客户/Aone 需求要 接入/支持 一个资源(e.g. alicloud_oss_bucket_inventory)；生成资源空/缺；cloudspec terraform 无资源；要拿 镇元 spec 落 provider 代码；**或既有资源出现 bug（错误码未重试 / attribute 缺失 / CRUD 不对齐 API / import 断链 等）**。NOT for 现有 PR 评审(用 terraform-pr-review)或简单 是否支持 查询(用 aone-triage 查证)。
---

# Provider 资源开发全流程

**两种模式，共用后半段（手改 → 回归用例 → 验收 → PR）**：
- **新增资源模式**：走完全部 8 步，`getTerraformResourceSpec` / acube 生成 / 生成-手写 diff 是必经。
- **修复/补丁模式**（非自动化生成资源的 bug fix / 属性补齐）：跳过步骤 1–5（不走生成），直接 **步骤 6（手改）→ 6.5（补/改回归用例）→ 7（AccTest）→ 8（PR）**；且 **6.5 是硬门**：修 bug 必须补/改一个"未打 patch 前会 fail、打 patch 后 pass"的锁定回归用例，随 PR 一并交付；确实无法测（如概率事件、需实体网关、外部依赖）则在 PR body 明写原因 + 已做的替代验证（静态查证 / 逻辑推演 / 同类资源对比）。

从客户需求到 provider PR。真源=镇元(Zhenyuan)`ResourceTypeSchema`;生成器吃镇元料出一整套(resource.go + _test.go + service + provider.go 注册行 + website 文档 markdown),但它可能空生成/部分生成;**只信 acc 测**。

## 工具/路径
- 工作区一律通过 `bootstrap/workspace.sh dir <key>` 解析;provider key=`terraform_provider`;acube key=`acube`;生成器 key=`terraform_generator_v4`。
- cspec 仓 `cloudspec-model/<Product>_pop_*`;provider upstream=aliyun;Jarvis 提交 GitHub PR/评论/推分支必须使用 `JARVIS_GITHUB_TOKEN` 对应的 `api-tool-agent` 身份,head=`api-tool-agent:<branch>`;acube/terraform-generator-v4 见 `config/workspaces.json`。
- Acube 在线生成工具: `tools/acube_terraform_generate.py`。
- 生成差异/语义检查工具: `tools/terraform_generated_diff.py`。
- **错误码语义查证**(客户 acc/apply 报错、retry 白名单决策):读 `.claude/skills/aone-triage/references/aliyun-error-code-lookup.md`(跨 skill 复用,给定 product+code 出 HTTP/中英 message/官方 retry 建议/相邻错误码)。
- 编码交 developer 子代理,改文件先 worktree,acc 测过才交。

## Aone 分单与同步
非自动化生成链路、需要 Jarvis 内部研发处理的 Terraform Provider 资源,必须创建或复用 **terraform-alicloud** 内部研发单:
- 项目: `tf_provider` / `528766`;指派给 Jarvis 自己 `WORKER_1782379562571`。
- 与 Terraform-客户需求池的客户主单双向关联;拿到内部单 id 后按 bookend 先 claim,收尾 done+release。
- 主要研发进展、生成/手改差异、验证细节、PR/CI/验收信息优先同步到内部研发单。
- 客户主单只同步关键节点摘要:已转内部单、发现镇元/Cloudspec/API 卡点、资源模型问题、需要客户感知的决策或阻塞。
- 如客户主单还关联 `cloudspec_gap` 或云产品上游 Aone,依赖方协作的详细问题同步到对应依赖单,不要混写在客户主单里。

## 步骤
1. **查证 Terraform ↔ Cloudspec 身份** — OpenAPI + provider 源码确认缺;`getTerraformResourceSpec` 只看映射,不代表实现,且找不到可能返无关资源,别信。
2. **查 acube resourceTypeCode** — 通过 Terraform 资源名推 product/resourceCode 后,先 get,再 list 降级;get 只有 `SUCCESS` 无 `data` 就按未命中处理。
3. **镇元建模/发布判断** — 如果 get/list 都找不到,说明资源未进预发/线上 resourceTypeCode;先回复 Aone 推动镇元发布或确认别名,除非本地已有 cspec 分支可继续验证。
4. **生成** — 标准入口走 Acube `createLocalBuildTask`:先 `resourceTypeCode/get` 取料,再 `createMapping`,再 `createLocalBuildTask`,用 `tools/acube_terraform_generate.py` 落盘 raw JSON/logs/files/generated/summary。只有 Acube 不可用或需验证本地 cspec 分支时,才 fallback 到本地 `cloudspec terraform -r <terraform_resource> -e pre -o <dir>`;报 `no that resource` 时用 `<CloudspecResourceCode>` 重试并记录 partial output。
5. **生成 vs 手写 diff** — 用 `tools/terraform_generated_diff.py` 看 Acube generated 与手写分支差异,先判断缺主体/缺测试/缺文档/缺 provider 注册/缺 service,再看 `resourceNotExistCondition` 等语义风险,最后手改。
6. **手改生成缺陷 / bug fix** — OSS 等 XML 产品:`client.Do("Oss",xmlParam(...))` 非 Roa(治 `<` 解析错);PUT body 按 schema **固定元素序**(治 MalformedXML);update 删-再-PUT 绕 AlreadyExists。**修复模式**下,查证根因(OpenAPI 错误码/schema/CRUD 语义)后照同 package 已有正确写法照抄,不 refactor 无关代码。
6.5. **补/改回归用例(修复模式硬门)** — 修 bug 时必须在同 PR 补/改一个"未打 patch 前会 fail、打 patch 后 pass"的用例(unit 或 acc,能锁定回归即可),随 CR/PR 一并评审。选型:
   - 若可用小并发/构造错误码/mock RPC 稳定复现,写 unit 或 acc 用例。
   - 若是概率事件(如服务端锁冲突、限流),可用同一 Config 里 declare N 个 resource 让 Terraform 并发 apply 触发。参照本 PR 9916 `TestAccAliCloudESARoutineRoute_lockRetry` 模式:同一 SiteId+RoutineName 下起 N 个 route,未打 patch 会 LockFailed 挂;打 patch 后 apply/destroy 全过。
   - 真无法测(需实体网关/外部依赖/无稳定复现路径)才允许豁免:在 PR body 明写"无法用例复现:<原因>",并列出替代验证(同类资源对比 / 静态查证 / manual repro log)。审阅人可拒。
7. **验收** — 真实 AccTest 优先用 `invoke-terraform-acc-test-remote` 走 ACube/FC 远程执行,避免长时间占用本机;本地只跑 `go test ./alicloud -run '^$'`、小单测、lint、示例 `terraform validate` 等轻量检查。远程 AccTest 过 create+update+import 才算数(**修复模式**至少要过 6.5 补的回归用例 + 原有主用例);跨账号/企业账号资源要隔离 ambient `ALICLOUD_ACCESS_KEY`/`ALICLOUD_SECRET_KEY`,显式声明测试需要的多把 AK 环境变量,并用 STS/CLI 验证每把 AK 的 caller account,但任何文档/评论/示例都不能泄露真实 AK/SK。
8. **PR** — `bootstrap/github-identity.sh check` → `bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud HEAD <branch>` → `bootstrap/github-identity.sh gh pr create --repo aliyun/terraform-provider-alicloud --head api-tool-agent:<branch>`;带 resource+test+service+provider注册+website 文档;无 AI 署名。缺 `JARVIS_GITHUB_TOKEN` 或登录名不是 `api-tool-agent` 时阻断并升级,禁止回退个人账号或 ambient git 凭据。

## Terraform 资源名解析
优先查 Acube Terraform 映射接口,不要维护固定映射表:

```bash
curl -s "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=<terraform_resource>" \
  -H "accept: */*"
```

读取 `data.terraformResourceSpecModel.namespace` / `data.terraformResourceSpecModel.resourceTypeCode` 作为 product/resourceCode。接口查不到时,再把 `alicloud_【产品名下划线】_【资源名下划线】` → product/resourceCode 各自转 PascalCase 作为候选,并进入下方 resourceTypeCode get/list 查证;边界不确定时结合客户描述/Next API/OpenAPI product 先确定 product,不要凭命名猜死。

## acube resourceTypeCode 查证
```bash
product=ResourceManager
resourceCode=HandshakeAcceptance

# 预发单资源定义;SUCCESS 但无 data = 未命中
curl -s "https://pre-acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/get?env=pre&isShowChangeLog=false&product=${product}&resourceCode=${resourceCode}" \
  -H "accept: */*"

# 线上域名也查一遍;必要时分别试 env=pre/prod,以是否有 data 为准
curl -s "https://acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/get?env=pre&isShowChangeLog=false&product=${product}&resourceCode=${resourceCode}" \
  -H "accept: */*"

# get 无 data 时降级查产品资源列表;列表里没有 resourceCode 才判定未发布/未同步
curl -s "https://pre-acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/list?product=${product}&released=false" \
  -H "accept: */*"
curl -s "https://acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/list?product=${product}&released=false" \
  -H "accept: */*"
```

sanity check: 用同产品已存在资源(如 `Handshake`) 反查一次;若它有 data,说明接口链路正常,目标 resourceCode 无 data 就是真缺失。

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

## 红线
不碰 master;无可评审 diff 不空发;对外无 AI 署名;Aone 唯一真源(进展 sync/完工 done);**修 bug 无回归用例又不写豁免原因 = 空发**。
