# TF 客户需求单诊断与路由(aone-triage reference)

> 触发条件:aone-triage 读单后,发现工单在 **tf_customer 池(1086837)** 或标题/涉及云产品含 alicloud_xxx / Terraform 关键字。用于诊断缺口在哪一层、路由到对的人,避免 jarvis / 谜拟 / 过载 / 临钧 / 新山 五方互踢皮球。

## 核心决策树

```
┌── 读单 + 抽真实诉求 (末段"限制/差异/仍需/不支持类似 X" 是真实诉求) ──┐
│                                                                    │
▼                                                                    │
产品 in 【专属维护名单】?                                              │
├─ YES → 直接指派对应负责人, status=问题解决中, @负责人 [end]         │
└─ NO ↓                                                              │
                                                                     │
诉求是"现有资源缺 X 属性/值/行为"(有 analog 产品/资源可类比)?         │
├─ YES → 查 analog 是"API 原生"还是"Provider 侧适配":                │
│         ├─ API 原生 + provider 透传 → 缺口在上游产品 API →         │
│         │   @提单人 + status=待上游排期 [end]                       │
│         └─ Provider 侧适配 → 我们团队可复制路径 → 走下方镇元分支    │
└─ NO (纯接入新资源 / 修 provider bug) ↓                             │
                                                                     │
镇元 schema 建模 + 测试覆盖度 100%?                                   │
├─ NO ↓                                                              │
│     优先级=紧急 OR 距计划截止 < 14 天?                              │
│     ├─ YES → 关联单指派 新山(521957)                                │
│     └─ NO  → 关联单指派 谜拟(479782)                                │
│     原单同步指派 + status=问题处理中 + @指派人                       │
└─ YES → provider 代码类型:                                          │
        ├─ 自动生成 (`This file is generated automatically`) →       │
        │   关联单指派 临钧(429768)                                   │
        └─ 手写 → 关联单指派 过载(484483)                             │
        原单同步指派 + status=问题处理中 + @指派人                     │
```

关联单一律建在 **terraform-alicloud** 项目(528766, `pools.tf_provider`),类型 = 缺陷/需求(视诉求),双向关联到源客户单。

## 团队分工速查

### 专属维护名单(直接指派该负责人,不走镇元/生成器)

这批云产品的 provider 代码由专人维护,**不接镇元**;不用查覆盖度,直接指派 + @。

| 云产品 | 花名 | 工号 |
|---|---|---|
| 容器服务 Kubernetes (ACK) | 若即 | 377376 |
| 日志服务 SLS | 豁朗 | 269032 |
| 消息服务 MNS | 曼红 | 38570 |
| OSS | 凡修 | 71145 |
| 弹性伸缩 ESS | 扶柳 | WB530580 |
| 表格存储 OTS | 景哲 | 263417 |
| E-MapReduce (EMR) | 鱼戏 | 373227 |
| RDS | 柴天生 | WB01586841 |
| PolarDB | 米汐 | 527630 |
| MSE | 棠溪 | 401341 |
| ClickHouse | 逸颉 | 439859 |

### 通用路由角色(其他云产品走此表)

| 场景 | 花名 | 工号 |
|---|---|---|
| 镇元 schema/覆盖度未 OK,非紧急且距 DDL ≥14 天 | 谜拟 | 479782 |
| 镇元 schema/覆盖度未 OK,紧急 或 距 DDL <14 天 | 新山 | 521957 |
| 镇元 OK + provider 代码由生成器产出 | 临钧 | 429768 |
| 镇元 OK + provider 代码手写(默认兜底) | 过载 | 484483 |

### 上游 API 缺口(纯上游产品团队问题)

- 不走上述任一路由,**@提单人**(工单 creator)+ status=待上游排期
- 由提单人协助转对应云产品的 API/OpenAPI 团队评估;jarvis 不代跨团队协调

## Step 1 — 读单 + 抽真实诉求

aone-triage 主流程已跑 `aone-get.sh`;此处补充抽取关键字段:

- **description 全文,尤其末段** —— 小蜜/AI 预答工单末尾常有"注意/与 X 不同/仍需/只支持"类描述,那就是客户想改的真痛点(经典误诊源,见 memory `read-description-last-paragraph`)
- 标题抽 `alicloud_<product>_<resource>` + 关键动词(接入/支持/缺少/催办)
- `assignedTo`(通常挂在谁头上不代表最终该谁修)
- `priority`(紧急/高/中/低) + `计划截止日期(80)` —— 决定紧急路由分支
- `涉及云产品(140097)` —— 优先匹配【专属维护名单】
- `工单ID(104264)` —— 小蜜工单号,备用

**真实诉求重述**:核心动作。用一句话把客户想要的能力/结果写下来,与 description 全文对齐。不对齐直接停,读客户原描述二次确认。

## Step 2 — 定位缺口(按决策树)

### 分支 A:产品在专属维护名单

直接跳 Step 3,指派对应负责人。**不查镇元 / 不做类比分析**——这批产品的 schema/代码流程与镇元无关。

### 分支 B:类比查证(诉求 = 现有资源缺某属性/值/行为)

诉求形如"能不能支持 X 属性 / X 状态值 / 行为像 Y 那样",且能找到 analog(通常客户已在描述里明说,如"与 PolarDB 不同")。

**查证套路(以 SSL Update 案为参考,ADB vs PolarDB)**:

```bash
# 1. 类比产品的 API meta - 用 aliyun CLI 官方 meta(不是 next.api 网页/预览):
aliyun <analog_product> <相关 Action> --help 2>&1 | grep -A 8 "<关键 Param>"

# 2. 目标产品的 API meta - 同样方式:
aliyun <target_product> <相关 Action> --help 2>&1 | grep -A 8 "<关键 Param>"

# 3. 类比产品的 provider 处理 - grep 相关字段:
grep -n "<字段>\|<API Param>" \
  "$(bash bootstrap/workspace.sh dir terraform_provider)/alicloud/resource_alicloud_<analog>_*.go"
```

**判定**:

| 类比产品的支持方式 | 结论 | 后续 |
|---|---|---|
| API 原生(合法值直接列出目标值) + Provider 直接透传 无转换 | 上游 API 缺口(目标产品 API 未暴露该值/属性) | 分支 F |
| API 侧同样只有基础值 + Provider 做了额外适配(如二次调用/参数拼装) | 我们团队可复制该适配路径 | 走分支 C |
| API 原生 + Provider 做了转换(如映射 Enable↔Enabled) | 视 API 本身支不支持 —— 支持则上游缺口(分支 F),否则我们团队做(分支 C) | 视情况 |

判定原则:**只要目标产品 API 层没有该能力,我们无法在 provider 侧凭空造出**;此时 100% 是上游产品团队职责,走 F。

### 分支 C:镇元 schema + 覆盖度查证(非专属名单产品,需要我们团队做)

```bash
product=<Product>          # PascalCase, e.g. Apig
resourceCode=<Resource>    # PascalCase, e.g. Domain

# ① 镇元 schema get/list
for env in pre online; do
  curl -s "https://pre-acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/get?env=${env}&isShowChangeLog=false&product=${product}&resourceCode=${resourceCode}" \
    -H "accept: */*" | python3 -c '
import json,sys
d=json.load(sys.stdin); data=d.get("data") or {}
print("hasData:",bool(data),"version:",data.get("resourceTypeVersion"),"title:",data.get("title"))'
done
curl -s "https://acube.aliyun-inc.com/api/v1/terraform/generator/cloudspec/resourceTypeCode/list?product=${product}&released=true" \
  -H "accept: */*" | python3 -c "import json,sys; d=json.load(sys.stdin); print('in released list:',('${resourceCode}' in (d.get('data') or [])))"

# ② 测试覆盖度 - POP API GetResourceModelTestCaseQualityByResource,须 --force + AK/SK
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
    cov=d.get("CoverageDetail") or {}; tc=d.get("TotalCases") or {}
    print("RequestId:",d.get("RequestId"),"| Message:",d.get("Message"))
    print("CoverageScore:",cov.get("CoverageScore"),"| PASS/FAIL:",len(tc.get("PASS") or []),"/",len(tc.get("FAIL") or []))
except Exception as e: print("parse error:",e)'
done
```

**判定**:
- 镇元 get 有 data + released list 命中 + CoverageScore == 1.0 + PASS>0 → **镇元 OK**,走分支 D
- 否则 → **镇元 NOT OK**,走分支 E

**sanity check(必跑)**:拿同产品已发布的其他资源(如查同一 product 的 list,取 released=true 里另一个)以相同接口/env 复查,拿到完整 CoverageDetail 说明凭据/接口正常;否则先排查凭据(`aliyun configure list` default 是否 Valid)。

### 分支 D:镇元 OK,判定 provider 代码类型

```bash
provider_repo="$(bash bootstrap/workspace.sh dir terraform_provider)"
head -3 "$provider_repo/alicloud/resource_alicloud_<product>_<resource>.go" 2>/dev/null
```

- 首行有类似 `// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!` → **生成器产出** → 指派 临钧(429768)
- 无该注释(手写) → 指派 过载(484483)

若文件不存在:说明 provider 代码尚未合入(镇元 OK 但 provider 未生成/合入)——按"生成器产出待跑"处理 → 指派 临钧(429768),comment 里注明"资源代码尚未合入,请触发生成 + PR"。

### 分支 E:镇元 NOT OK,判定紧急度

```bash
priority=$(bash bootstrap/aone-get.sh <id> | python3 -c '
import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
  if f.get("identifier")=="priority": print(f.get("displayValue"))')
ddl=$(bash bootstrap/aone-get.sh <id> | python3 -c '
import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
  if f.get("identifier")=="80": print(f.get("value"))')
echo "priority=$priority ddl=$ddl"

# 距今天算 gap
days_left=$(python3 -c "
from datetime import date
d='$ddl'.strip()
if d:
  print((date.fromisoformat(d) - date.today()).days)
else:
  print('N/A')")
echo "days_left=$days_left"
```

- `priority == '紧急'` 或 `days_left < 14` → 指派 新山(521957)
- 否则 → 指派 谜拟(479782)

### 分支 F:上游 API 缺口

- 不建关联单(镇元/我们团队无事可做)
- 提取 `creator.displayName` 作为提单人,评论正文 @提单人 + 请其协助转对应云产品 API 团队评估
- status 改 `待上游排期`(若 CLI 报 "unsupported target status" 却能实际写入,重试;或落 `待排期` 作二级降级)

## Step 3 — 执行路由动作(写操作,先授权)

### 分支 A / D / E(需要指派我方或产品专属人员)

```bash
# 1. 建关联单在 terraform-alicloud (528766),category 视诉求
bin/a1id -- project workitem create \
  --project 528766 \
  --category <req|bug|task> \
  --title "<清晰标题:资源/属性/诉求>" \
  --assignee <工号> \
  --description "@<花名>(<工号>) 客户 <URL/ID> 请求 <诉求>;详见源单描述与本单描述..."
# 记下返回的新单 id
NEW_ID=<新单 id>

# 2. 双向关联
bin/a1id -- project workitem relation add <源工单ID> relate:$NEW_ID
bin/a1id -- project workitem relation add $NEW_ID relate:<源工单ID>

# 3. 源工单指派 + 状态
bin/a1id -- project workitem update <源工单ID> --assignee <工号>
bin/a1id -- project workitem update <源工单ID> --status 问题处理中
# (专属名单产品用 "问题解决中"——名字接近,别写混)
```

### 分支 F(上游 API 缺口,只 @提单人)

无需建关联单,只发评论 + 改状态:

```bash
bin/a1id -- project workitem update <源工单ID> --status 待上游排期
# (CLI 报 unsupported 可实际写入,若拒绝改用 "待排期")
```

## Step 4 — 回复模板

### 模板 A:类比查证 → 上游 API 缺口(分支 F)

```
### 结论
客户诉求「<一句话诉求重述>」的缺口在**上游 <产品> 产品侧 OpenAPI**,非
Terraform Provider / 镇元(Cloudspec)侧可闭环。

### 双层核实
1. **<类比产品> `<字段>` 在 OpenAPI 原生**(aliyun CLI 官方 meta):
   - `<analog_product>/<Action>/<Param>` = <Type>,合法值 <值域>
   - Terraform Provider 侧仅透传(`resource_alicloud_<...>.go:<L>` `<字段赋值>`),无适配
2. **<目标产品> `<Action>` 当前不支持 <诉求>**:
   - `<target_product>/<Action>/<Param>` = <Type>(仅 <当前值域>),`aliyun <target_product> <Action> --help` 官方 meta 可复核
   - <目标产品> 无独立 <相关能力> action

### 建议行动
@<提单人花名>(<工号>) 请协助将此需求转 <目标产品> 产品团队评估;工单状态已改
为「待上游排期」,<目标产品> 侧 API 上线支持 <诉求> 后 Terraform Provider 侧
可同步跟进 schema 改造。
```

### 模板 B:镇元/provider 路由(分支 D / E)

```
### 结论
客户诉求「<一句话诉求重述>」应由 <指派人所在层:镇元 schema/provider 代码>
侧承接。已建关联单 <NEW_ID> 到 terraform-alicloud 项目,指派 @<花名>(<工号>)
跟进,源工单同步指派并改状态为「问题处理中」。

### 查证依据
1. **镇元**:<get: 有/无 data | list released 是否命中 | CoverageScore=<x>>
   - Env=online: <Message>(RequestId: <...>)
   - Env=pre:    <Message>(RequestId: <...>)
2. **provider 代码类型**:<自动生成 / 手写>(依据 `alicloud/resource_...go:1-3`
   <是否有 generated automatically 注释>)
3. **紧急度**:priority=<>, 计划截止=<>, 剩余=<>天,故路由到 <人>

### 关联单
- <NEW_ID>: <标题>,项目 528766 terraform-alicloud
- 双向关联已加

@<花名>(<工号>) 烦请跟进上述查证结论,进度请在两侧工单同步回帖。
```

### 模板 C:专属名单产品(分支 A)

```
### 结论
「alicloud_<product>_<resource> <诉求>」属 <产品中文名> 云产品域,由该云产品
provider 专人维护(不接镇元)。已指派 @<花名>(<工号>) 跟进,状态改为
「问题解决中」。

@<花名>(<工号>) 请协助评估 <诉求>,进度请回帖。
```

## bookend 与 wrap.sh 参数

评论发布走 aone-triage 主流程的 bookend(claim → wrap.sh done → release)。关键细节:

- `wrap.sh done <id> "<summary>" <status|--no-status>` —— **位置参数**,不吃 `--summary-file` / `--status` 命名参数
- 关联单(528766)一律**不 claim** —— jarvis 无 tf_provider 池管理权,建单 + @对方即可
- 分支 F 的评论 + 状态修改不需要建关联单,直接 wrap.sh done + 单独一次 status update

## 反模式

- ❌ 只按标题查证,不读 description 末段"限制/差异/仍需" —— 经典误诊(SSL Update 案就栽这)。参见 memory `read-description-last-paragraph`
- ❌ 直接把"类比产品支持 X"当"我们也应该做" —— 必须先分清是**API 原生**还是**Provider 侧适配**,前者 100% 是上游产品缺口
- ❌ 上游 API 缺口场景建了关联单 —— 应该只 @提单人 + 待上游排期,别拖谜拟/过载/临钧下水
- ❌ 专属名单产品还去查镇元覆盖度 —— 不接镇元,直接指派
- ❌ 花名 @ 不带工号(`@谜拟`) —— a1 有时能补,有时不能,显式 `@谜拟(479782)` 保险
- ❌ 关联单不双向 —— `a1 relation add` 必须调两次(A→B, B→A),否则一侧看不到对方
- ❌ 状态用 `--status 已完成` / `方案功能已存在` 兜底 —— 前者不在合法值,后者语义错(客户真诉求还没解决)
- ❌ 强行 `--assignee` 指派专属名单以外的产品到过载/谜拟 —— 违反分工表,让本团队背不该背的锅
