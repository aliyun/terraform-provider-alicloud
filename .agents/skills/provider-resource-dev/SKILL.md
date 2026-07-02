---
name: provider-resource-dev
description: Use when developing or diagnosing a NEW alicloud Terraform provider resource end-to-end (Terraform 资源名解析 → 镇元/Cloudspec resourceTypeCode 查证 → acube 映射/生成 → 生成代码 vs 手写代码 diff → 手改 → acc 验收 → PR). Trigger when a customer/Aone req asks to 接入/支持 a resource (e.g. alicloud_oss_bucket_inventory), a generated resource is empty/missing, cloudspec terraform reports no resource, or you must take 镇元 spec to provider code. NOT for reviewing an existing PR (use terraform-pr-review) or answering simple 是否支持 (use aone-triage 查证).
---

# Provider 资源开发全流程

从客户需求到 provider PR。真源=镇元(Zhenyuan)`ResourceTypeSchema`;生成器吃镇元料出一整套(resource.go + _test.go + service + provider.go 注册行 + website 文档 markdown),但它可能空生成/部分生成;**只信 acc 测**。

## 工具/路径
- 工作区一律通过 `bootstrap/workspace.sh dir <key>` 解析;provider key=`terraform_provider`;acube key=`acube`;生成器 key=`terraform_generator_v4`。
- cspec 仓 `cloudspec-model/<Product>_pop_*`;provider upstream=aliyun;Jarvis 提交 GitHub PR/评论/推分支必须使用 `JARVIS_GITHUB_TOKEN` 对应的 `api-tool-agent` 身份,head=`api-tool-agent:<branch>`;acube/terraform-generator-v4 见 `config/workspaces.json`。
- Acube 在线生成工具: `tools/acube_terraform_generate.py`。
- 生成差异/语义检查工具: `tools/terraform_generated_diff.py`。
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
6. **手改生成缺陷** — OSS 等 XML 产品:`client.Do("Oss",xmlParam(...))` 非 Roa(治 `<` 解析错);PUT body 按 schema **固定元素序**(治 MalformedXML);update 删-再-PUT 绕 AlreadyExists。
7. **验收** — 真实 AccTest 优先用 `invoke-terraform-acc-test-remote` 走 ACube/FC 远程执行,避免长时间占用本机;本地只跑 `go test ./alicloud -run '^$'`、小单测、lint、示例 `terraform validate` 等轻量检查。远程 AccTest 过 create+update+import 才算数;跨账号/企业账号资源要隔离 ambient `ALICLOUD_ACCESS_KEY`/`ALICLOUD_SECRET_KEY`,显式声明测试需要的多把 AK 环境变量,并用 STS/CLI 验证每把 AK 的 caller account,但任何文档/评论/示例都不能泄露真实 AK/SK。
8. **PR** — `bootstrap/github-identity.sh check` → `bootstrap/github-identity.sh push api-tool-agent/terraform-provider-alicloud HEAD <branch>` → `bootstrap/github-identity.sh gh pr create --repo aliyun/terraform-provider-alicloud --head api-tool-agent:<branch>`;带 resource+test+service+provider注册+website 文档;无 AI 署名。缺 `JARVIS_GITHUB_TOKEN` 或登录名不是 `api-tool-agent` 时阻断并升级,禁止回退个人账号或 ambient git 凭据。

## Terraform 资源名解析
规则: `alicloud_【产品名下划线】_【资源名下划线】` → product/resourceCode 各自转 PascalCase。

示例:
```text
alicloud_resource_manager_handshake_acceptance
=> product=ResourceManager
=> resourceCode=HandshakeAcceptance
```

如果边界不确定,结合客户描述/Next API/OpenAPI product 先确定 product;仍不确定时对 `alicloud_` 后缀做最长 product 前缀试探,逐个调用 resourceTypeCode get/list,不要凭命名猜死。

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
不碰 master;无可评审 diff 不空发;对外无 AI 署名;Aone 唯一真源(进展 sync/完工 done)。
