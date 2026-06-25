---
name: provider-resource-dev
description: Use when developing a NEW alicloud Terraform provider resource end-to-end (镇元/Cloudspec 建模 → acube 映射 → 在线/本地生成 → 手改 → acc 验收 → PR). Trigger when a customer/Aone req asks to 接入/支持 a resource (e.g. alicloud_oss_bucket_inventory) and there is no provider resource yet, or 生成器/cloudspec terraform 出空/报错, or you must take 镇元 spec → 可交付 provider 代码. NOT for reviewing an existing PR (use terraform-pr-review) or answering 是否支持 (use aone-triage 查证).
---

# Provider 资源开发全流程

从客户需求到 provider PR。真源=镇元(Zhenyuan)`ResourceTypeSchema`;生成器吃镇元料出 .go,生成器对 XML 产品有缺陷需手改;**只信 acc 测**。

## 工具/路径
- cspec 仓 `cloudspec-model/<Product>_pop_*`;provider `~/go/src/.../terraform-provider-alicloud`(fork=ChenHanZhang,upstream=aliyun);acube `~/IdeaProjects/a-cube-aliyun-com`(详 [[acube-tf-generation-pipeline]])。
- 编码交 developer 子代理,改文件先 worktree,acc 测过才交。

## 步骤
1. **查证** — OpenAPI + provider 源码两层,确认缺;`getTerraformResourceSpec` 看映射(找不到返第一行 ClickHouse,别信)。
2. **镇元建模** cspec — resource+ops。容器(object/array)节点**只映叶子+`[*]`**,整块映 `$` 崩生成器;Id `@readonly`;对象嵌套 ≤3 层。`cloudspec check`/`build`/`terraform -r X` 本地验。
3. **在线生成** — 见下「接口示例」。未发线上用 `env=pre`。
4. **手改生成缺陷** — OSS 等 XML 产品:`client.Do("Oss",xmlParam(...))` 非 Roa(治 `<` 解析错);PUT body 按 schema **固定元素序**(治 MalformedXML);update 删-再-PUT 绕 AlreadyExists。
5. **验收** — `TF_ACC=1 go test ./alicloud/ -run TestAcc<R> -v -timeout 40m`;过 create+update+import 才算数。
6. **PR** — fork 分支 → `gh pr create --repo aliyun/terraform-provider-alicloud`;带 resource+test+service+provider注册+website 文档;无 AI 署名。

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

## 常见坑
- 生成空/`获取spec数据失败` → env 没到预发镇元(查 createLocalBuildTask Env 是否 online)。
- `MalformedXML` → body map 随机序;`invalid character '<'` → 用了 Roa(JSON)发 XML。
- 直接发生成版不验收 → 必踩 OSS XML;acc 不过别交付。

## 红线
不碰 master;无可评审 diff 不空发;对外无 AI 署名;Aone 唯一真源(进展 sync/完工 done)。
