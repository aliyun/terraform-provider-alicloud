---
name: resource-onboarding-progress
description: >-
  Use when a customer Aone work item asks to 接入/支持 a specific alicloud Terraform resource
  (alicloud_xxx) and the user wants to assess progress / catch up / push it forward (NOT to start
  the resource implementation themselves). Triggers: 工单链接/ID + 进度怎么样 / 处理了吗 / 卡哪了 /
  覆盖度怎么样 / 催一下；标题含「接入 TF / Terraform 支持 alicloud_xxx」且 assignee 非空 / 已有
  POP 用例补充类评论的工单。Covers: 读单 → 查证「镇元 schema 是否建模 + 测试覆盖度是否达标」 →
  草拟回复（@产品侧负责人录入用例 / @谜拟 镇元侧监督跟进）→ 用户审核 → bookend。
  NOT for 资源代码开发（用 provider-resource-dev）/ 是否支持的纯查证（用 aone-triage）/ PR 评审
  （用 terraform-pr-review）。
---

# 资源接入 TF Provider 进度评估与催办

> 客户给 TF 客户需求池（tf_customer, project 1086837）提了「接入 alicloud_xxx」工单，已经在流转一段时间。本 skill **不是**从头开发资源——是回答「现在卡在哪、还差什么、催谁」。

## 卡点定位心智模型

资源接入有 4 个串联阶段，任一阶段卡住整个流程不动：

| 阶段 | 检查方式 | 卡住信号 |
|---|---|---|
| ① **镇元 schema 建模** | `pre-acube` get/list resourceTypeCode | get 无 data / list 无该 code |
| ② **测试用例录入并触发** | `APISpecInner / GetResourceModelTestCaseQualityByResource` | "未找到测试执行记录" / CoverageScore=null |
| ③ **覆盖度达标** | 同上接口 | CoverageScore < 1.0 |
| ④ **provider 代码合入** | `terraform-provider-alicloud` master 是否有 resource_xxx.go | 缺主体 / 缺文档 / 缺 provider 注册 |

本 skill 重心在 **② + ③**——这是经典卡点，一般产品侧补完 POP 用例就停，没人推动镇元侧用例录入。

## 开局：读单 + 抽资源名

```bash
# 1. 读工单（3h 缓存）
bash bootstrap/aone-get.sh <id>
# 2. 读评论 / 活动 — 看是不是已经有产品侧动作
bin/a1id -- project workitem comment list <id>
bin/a1id -- project workitem activity <id>
```

关键字段：标题（抽 `alicloud_xxx`）/ status / assignee（= 产品侧负责人，催办对象）/ 评论里有没有 "POP 补充测试用例" 这类信号 / 已抄送/参与人。

资源名解析规则同 [[provider-resource-dev]]：`alicloud_<product>_<resource>` → product/resourceCode PascalCase（边界不确定时按 Next API / OpenAPI product 收敛）。

## 查证一：镇元 schema 是否建模

直接复用 [[provider-resource-dev]] 的 acube resourceTypeCode get/list 段（同一套接口，无需重写）：

```bash
product=Apig
resourceCode=Domain

# 单资源详情 — env 参数:pre / online；prod 在线上 acube 通常返 hasData False，不要误判
for env in pre online; do
  curl -s "https://pre-acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/get?env=${env}&isShowChangeLog=false&product=${product}&resourceCode=${resourceCode}" \
    -H "accept: */*" | python3 -c '
import json,sys
d=json.load(sys.stdin); data=d.get("data") or {}
print("hasData:",bool(data),"version:",data.get("resourceTypeVersion"),"title:",data.get("title"))'
done

# 已发布列表里有没有
curl -s "https://acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/list?product=${product}&released=true" \
  -H "accept: */*" | python3 -c "import json,sys; d=json.load(sys.stdin); print('in released list:',('${resourceCode}' in (d.get('data') or [])))"
```

**判定**：
- get 有 data + list released=true 命中 → 已建模发布（继续查证二）
- get 无 data + list 不命中 → 卡在镇元 schema 这层，催镇元侧建模发布；不进查证二

## 查证二：测试覆盖度查询（核心动作）

POP API `APISpecInner / GetResourceModelTestCaseQualityByResource`（v2021-07-13）查"该资源是否跑过测试用例 + 覆盖度多少"。这条接口集团内只通过 aliyun POP（AK/SK 签名）调；直接 curl amp-apispec REST 路径会被 SSO 拦（401 logout）。

```bash
product=Apig
resourceCode=Domain

# 跨 env 查（覆盖度按"online"是主战场；pre/daily 可旁证）
for env in online pre daily; do
  echo "=== Env=$env ==="
  aliyun --product APISpecInner --version 2021-07-13 \
         --endpoint apispec-share.cn-zhangjiakou.aliyuncs.com \
         --method POST --force \
         APISpecInner GetResourceModelTestCaseQualityByResource \
         --ServiceCode "$product" --ResourceCode "$resourceCode" --Env "$env" 2>&1 | python3 -c '
import json,sys
try:
    d=json.loads(sys.stdin.read())
    cov=d.get("CoverageDetail") or {}
    tc=d.get("TotalCases") or {}
    print("RequestId:",d.get("RequestId"))
    print("Message:",d.get("Message"))
    print("CoverageScore:",cov.get("CoverageScore"))
    print("OperationCoverageScore:",cov.get("OperationCoverageScore"))
    print("PropertyCoverageScore:",cov.get("PropertyCoverageScore"))
    print("PASS/FAIL cases:",len(tc.get("PASS") or []),"/",len(tc.get("FAIL") or []))
except Exception as e:
    print("parse error:",e)'
done
```

**响应语义对照**：

| 返回 | 含义 |
|---|---|
| `CoverageScore: <float>` + `TotalCases.PASS > 0` | 用例已跑，看具体数字判断是否 100% |
| `Message: "未找到当前资源[X]测试的执行记录..."` + `CoverageScore: null` | 镇元建模过但**没人录入并触发执行用例**——典型卡点 |
| `Message: "指定条件：环境[X]未找到对应的资源[ServiceCode::Product::Version::ResourceCode]"` | 该 env 镇元没注册此资源（不一定是"未发布"——`Env=prod` 是非法值也会这样，按定义 Env 只接 online/pre/daily/gray） |
| 复合错误 / 404 / signing 失败 | 检查 aliyun CLI AK/SK 是否 valid（`aliyun configure list`） |

**sanity check**（**必跑**，避免空查证误报）：拿**同产品**已经在跑的资源（list released=true 里随便挑一个，如 `Apig/Gateway`）以相同 env 查一次。如果它返回完整 `CoverageDetail` → 接口本身正常，目标资源"无记录"是真缺；如果它也无记录 → 是接口 / 凭据 / env 选错，先排查再下结论。

## 查证三（可选）：provider 仓代码合入状态

如果上两步都通过，覆盖度还不到 100%，看一下 provider 仓有没有合：

```bash
provider_repo="$(bash bootstrap/workspace.sh dir terraform_provider 2>/dev/null)"
[ -d "$provider_repo" ] || echo "provider repo not on disk; skip"
[ -d "$provider_repo" ] && {
  ls "$provider_repo/alicloud/resource_${product,,}_${resourceCode,,}.go" 2>/dev/null || echo "resource file 缺失"
  ls "$provider_repo/website/docs/r/${product,,}_${resourceCode,,}.html.markdown" 2>/dev/null || echo "website doc 缺失"
}
```

但工单本身的卡点判定主要看 ① + ② + ③，provider 仓状态是辅助信息。

## 同步上游：查关联工单是否已有「Terraform 镇元对接」单

如果查证结论涉及 @谜拟（覆盖度不足 / schema 未建模 / 用例未录入），**必须**额外检查本工单有没有指派给谜拟、归属「Terraform 镇元对接」(project **2165097**, `pools.json` upstream.cloudspec_gap 已登记)项目的关联工单——多数情况下产品/客户侧已经为对接侧建过这条上游单子，jarvis 应当**同步在那条上游单下也 @ 一次**，避免谜拟只在源工单收到通知。

```bash
# 1. 列关联工作项
bin/a1id -- project workitem relation list <主工单ID> --category workitem -f json \
  > /tmp/rel.json

# 2. 筛 Owner=谜拟 的候选（a1 JSON 不含 project，需二次 aone-get 确认 space）
candidates=$(python3 -c '
import json
d=json.load(open("/tmp/rel.json"))
for it in d.get("workitem",[]):
    if it.get("Owner")=="谜拟":
        print(it["Identifier"])')

# 3. 对每条候选 aone-get 查 space 字段值是否 == 2165097
#    注:aone-get.sh 输出在 description 字段可能含未转义控制字符(Aone 原文带的);
#    jq 严格模式会拒绝,所以这里用 python3 json(更宽松)解析
for rid in $candidates; do
  pid=$(bash bootstrap/aone-get.sh "$rid" 2>/dev/null | python3 -c '
import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
    if f.get("identifier")=="space":
        print(f.get("value")); break')
  echo "$rid -> project=$pid"
  [ "$pid" = "2165097" ] && echo "  ↑ Terraform 镇元对接,需 cross-link @"
done
```

**找到 → 在该上游工单下也发评论**：
- 内容是**短 ping**，不要把完整查证报告复制过去；直接指回源工单链接 + 简述（"该资源测试覆盖度查证后发现用例未录入，源工单 <URL> 已 @；烦请协助跟进"）+ `@谜拟(479782)`。
- 这是 **cross-link 通知**，**不**对上游工单 claim / release / 算 bookend（jarvis 无该池管理权，submit_only；强行 claim 会污染对方池标签）。
- 直接 `a1 project workitem comment create <上游ID> -m "..."`，不走 wrap.sh（wrap.sh 会走 touch_ledger + 评论流，不适合 cross-link）。

**找不到 → 仍只在源工单 @谜拟**，不要主动建上游单（上游池 submit_only + jarvis 无权臆造单子；如果谜拟真的需要建上游跟踪单，由她自己或产品方建）。

## 草拟回复（结构固定，先给用户审核）

四块：**结论 → 镇元 schema 证据 → 测试覆盖度证据（带 RequestId）→ sanity → @ 两人**。

不带 AI 署名（AGENTS.md 工作纪律 #7）。

```
查证结果（alicloud_<product>_<resource> 接入 TF 进度）：

1. 镇元 schema：<已建模 ACS::Product::ResourceCode vX.Y / 未建模>
   <如果未建模，证据：get 无 data / released list 不命中>

2. 测试覆盖度（APISpecInner / GetResourceModelTestCaseQualityByResource）：
   - Env=online：<Message>（RequestId: <...>）
   - Env=pre：<Message>（RequestId: <...>）
   - Env=daily：<Message>（RequestId: <...>）
   结论：<CoverageScore=null + PASS/FAIL=0 → 用例未录入/执行 / CoverageScore=<x> 未达 100% / 已达 100%>

3. sanity：<同产品已发布资源 X 在 pre 拿到完整 CoverageDetail，证明接口与凭据正常>

@<产品侧负责人花名>(<工号>) 作为 <产品> 云产品侧负责人，烦请<具体动作:在镇元平台为
ACS::Product::ResourceCode 录入测试用例并触发执行(online/pre/daily) / 推动 schema
发布到镇元 / ...>，跑通后在评论区同步进度，便于继续推动 alicloud_xxx 接入 TF Provider。

@谜拟(479782) 镇元端测试用例监督跟进，烦请协助跟进本资源的<用例录入与执行 / schema 发布>
进度，确保覆盖度达到 100%。
```

**@ 写法**：a1 自动把 `@花名` 解析为 `@花名(工号)`；显式带工号防误命中。产品侧 @ 谁——通常就是工单 `assignee`（前期已经在做事），偶尔抄送/参与人里有更明确的产品测试负责人则按那个。

## bookend（写动作前）

不写工单仅查证可免 claim。**一旦要发评论就走完整 bookend**（AGENTS.md 工作纪律 #5）：

```bash
# 1. claim（仅打标签 + 冻结 jarvis-claim 痕迹到 .my-day/claim-prefix-<id>.txt）
bash bootstrap/claim.sh claim <id> 1086837   # tf_customer pool

# 2. wrap.sh done 发评论（自动 prefix claim 痕迹，不要单独 a1 comment create）
#    summary 就是上面审核过的完整回复；不传 status — jarvis 不擅自流转 Aone 状态
bash bootstrap/wrap.sh done <id> "<完整查证回复>"

# 3. 释放（本轮处理完，等产品/镇元侧跟进，未真完成 → release 打 jarvis-idle）
bash bootstrap/claim.sh release <id> 1086837
```

**真完成 vs 本轮释放**——本 skill 默认走 release（产品/镇元侧没真动起来前不算完）。仅当查证发现覆盖度已达 100% + provider 仓代码已合 + 工单产品方已确认验收 → 才走 `claim.sh finish`（打 jarvis-done 标签 + 改 status 为 `已发布待需求排期`）。

## 重要约定

- **aliyun CLI 走默认 profile** —— POP API 用 AK/SK 签名，不是 a1 token；前置确认 `aliyun configure list` default 是 Valid。
- **`--force` 不能省** —— aliyun CLI 不认 `APISpecInner` 是合法 product，没 `--force` 会 `ERROR: unchecked version 2021-07-13`。
- **`--Env` 合法值** —— online / pre / daily / gray；`prod` 不是合法值（API 定义里没有），传了会被服务端兜底报"未找到对应资源"，会误读为"prod 未发布"。
- **跨 env 同时查** —— 单 env 容易把"该 env 未发布"和"全无用例执行"搞混；至少 online + pre 两条都查清楚。
- **sanity check 不省** —— 同产品已发布资源对照，证明不是接口 / 凭据 / 环境问题；早期最容易出的就是"凭据过期返 401 误以为资源未跑"。
- **POP 接口定义/参数源头**：上游 model-design 文档（`api.aliyun-inc.com/v2?app=model-design ... APISpecInner ... GetResourceModelTestCaseQualityByResource`），需 SSO 登录；a-cube-aliyun-com 仓 `acube-service/src/main/java/com/aliyun/openplatform/helper/OpenApiHelper.java:382` 有 Java 等价实现。

## 反模式

- ❌ 直接 curl `pre-amp-apispec.aliyun-inc.com/api/v1/.../get_resource_model_test_case_quality_by_resource` —— 401 logout，要 SSO，jarvis 无能力。走 aliyun POP 才对。
- ❌ 只查 prod env / 只查 online，不做 sanity check —— "无记录" 可能是凭据问题。
- ❌ 用 wrap.sh done 之前先手动 `a1 comment create` 发查证 —— wrap.sh done 内部又会发一条，工单出现两条且 a1 无 comment delete，删不掉（参见 memory `feedback-wrap-done-single-comment`）。
- ❌ wrap.sh done 强传 status 改"长期跟进"/"问题解决中" —— jarvis 不擅自流转 status，由产品方/状态机推进；只在 finish 路径自动改为 `已发布待需求排期`。
- ❌ 没拿到用户对回复草稿的"发"就调 wrap.sh done —— supervised 模式 + 写动作必先授权。
- ❌ @ 写成 `@产品侧负责人` 这种伪 @ —— a1 解析不到 user，不触发通知；要 `@实际花名(工号)`。
