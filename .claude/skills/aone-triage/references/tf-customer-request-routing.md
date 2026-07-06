# TF 客户需求单诊断与路由(aone-triage reference)

> 触发条件:aone-triage 读单后,发现工单在 **tf_customer 池(1086837)** 或标题/涉及云产品含 alicloud_xxx / Terraform 关键字。用于诊断缺口在哪一层、路由到对的人,避免 jarvis / 谜拟 / 过载 / 临钧 / 新山 五方互踢皮球。

## 核心决策树

```
┌── 读单 + 抽真实诉求 (末段"限制/差异/仍需/不支持类似 X" 是真实诉求) ──┐
│                                                                    │
▼                                                                    │
诉求范围: 单一云产品/资源 vs Provider 侧全局改造?                     │
├─ Provider 全局改造 (不涉及单一 alicloud_xxx 资源:region 白名单 /   │
│   框架 utility / 公共 endpoint / provider.go 基础 / SDK bump 等)   │
│   → 关联单指派 新山(521957) [分支 G,end]                           │
└─ 单一产品/资源 ↓                                                    │
                                                                     │
产品 in 【专属维护名单】?                                              │
├─ YES → 直接指派对应负责人, status=问题解决中, @负责人 [end]         │
└─ NO ↓                                                              │
                                                                     │
诉求仅涉及文档改造(website/docs 变更,无 provider 代码改动)?           │
├─ YES → 关联单指派 过载(484483) [end]                                │
└─ NO ↓                                                              │
                                                                     │
诉求是"现有资源缺 X 属性/值/行为"(有 analog 产品/资源可类比)?         │
├─ YES → 查 analog 是"API 原生"还是"Provider 侧适配":                │
│         ├─ API 原生 + provider 透传 → 缺口在上游产品 API →         │
│         │   @提单人 + status=待上游排期 [end]                       │
│         └─ Provider 侧适配 → 我们团队可复制路径 → 走下方镇元分支    │
└─ NO (纯接入新资源 / 修 provider bug) ↓                             │
                                                                     │
镇元 OK? (三条件全满足才算 OK:                                        │
  ① API 在镇元已定义并发布                                            │
  ② 当前资源 schema 属性满足客户诉求(不缺字段)                        │
  ③ acube V2 覆盖度 CoverageScore == 1.0)                             │
├─ NO (三条件任一不满足) ↓                                            │
│     优先级=紧急 OR 距计划截止 < 14 天?                              │
│     ├─ YES → 关联单指派 新山(521957)                                │
│     └─ NO  → 关联单指派 谜拟(479782)                                │
│     原单同步指派 + status=问题解决中 + @指派人                       │
└─ YES → provider 代码类型:                                          │
        ├─ 自动生成 (`This file is generated automatically`) →       │
        │   关联单指派 临钧(429768)                                   │
        └─ 手写 → 关联单指派 过载(484483)                             │
        原单同步指派 + status=问题解决中 + @指派人                     │
```

关联单一律建在 **terraform-alicloud** 项目(528766, `pools.tf_provider`),类型 = 缺陷/需求(视诉求),双向关联到源客户单。**例外:临钧路由(生成器产出)不由 jarvis `a1 workitem create` 手动建单**——走 acube `createBuildTaskV2` 接口,acube 内部自动建单+指派临钧+触发生成/PR 工作流,jarvis 只查回 aoneId 做关联,详见 Step 3。


## 团队分工速查

详见 [team-roster.md](./team-roster.md)。

## Step 1 — 读单 + 抽真实诉求

aone-triage 主流程已跑 `aone-get.sh`;此处补充抽取关键字段:

- **description 全文,尤其末段** —— 小蜜/AI 预答工单末尾常有"注意/与 X 不同/仍需/只支持"类描述,那就是客户想改的真痛点(经典误诊源,见 memory `read-description-last-paragraph`)
- 标题抽 `alicloud_<product>_<resource>` + 关键动词(接入/支持/缺少/催办)
- `assignedTo`(通常挂在谁头上不代表最终该谁修)
- `priority`(紧急/高/中/低) + `计划截止日期(80)` —— 决定紧急路由分支
- `涉及云产品(140097)` —— 优先匹配【专属维护名单】
- `工单ID(104264)` —— 小蜜工单号,备用
- `space` / `workitemType` / `creator` —— 分类误建/承接单判定(见 Step 1.5)

**缺陷类型优先级覆写**:`workitemType` 为缺陷(功能缺陷/线上问题/性能瓶颈)时,**优先级一律视为紧急**,无视原单 `priority` 字段值。此覆写影响决策树镇元 NOT OK 分支路由(紧急 → 新山 521957)及关联单 `--priority` 传值。

**真实诉求重述**:核心动作。用一句话把客户想要的能力/结果写下来,与 description 全文对齐。不对齐直接停,读客户原描述二次确认。

## Step 1.5 — canned 前置分诊(命中直接回,免走决策树)

以下 8 类高频场景在读单后**先于决策树判定**——命中就发对应 canned,通常只 comment、不建关联单、不改状态,等客户补料或云产品回帖后再进 Step 2 决策树。

**共通 gate — 承接单检查**:发追料/科普 canned 前,先做一次同期 duplicate 排查(池/类型/creator/关键词),避免与已有承接单重复打搅客户:

```bash
# 按 creator 近 2 天在 tf_customer 池(1086837)搜同一提单人的活跃单
bin/a1id -- project workitem list --project 1086837 \
  --creator <empId> --updated-since "$(date -v-2d +%Y-%m-%d)"
# 或按 title 关键词搜
bin/a1id -- project workitem list --project 1086837 --query "<关键字>"
```

- **命中承接单** → 本 canned **免发**,评论"重复单,跟进转 <承接单ID>"并关本单,追料由承接单负责
- **未命中** → 按下方对应场景发 canned

### 场景 1 — 工单信息不完整

**判定**:客户只贴报错标题或截图,没有完整 tf 代码 / 完整报错 / 期望结果;description 里无法抽出真实诉求。

**发前 gate**:先按段首"共通 gate"排承接单——同期已有 tf_customer 完整单则本场景免发,让承接单去追料,别重复打搅客户。**误落 tf_provider 池(528766)分类错误单也走这条**:关本单 + 转对应客户单跟进。

```
请与我们确认清楚您的需求是什么、期望是什么;如有报错,请提供完整的 tf
代码和完整报错信息。
```

status 不动,等客户回帖。

### 场景 2 — OpenAPI 报错含 RequestId

**判定**:工单里贴了 RequestId + 具体错误码(如 `InvalidParameter.XXX`)。这类是上游 API 返回的错误,Provider 只是透传,应由云产品定位。

```
OpenAPI 返回的报错,请拉云产品的研发看一下,为什么报 <错误码>,哪个属性的
值引起的,正确应该传什么。
RequestId: <XXXXXX>
```

后续路由:【专属维护名单】产品 → 分支 A 直接指派;非专属 → @提单人协助转对应云产品研发,不建关联单(等同分支 F)。

### 场景 3 — 询问「TF 是否支持某功能」

**判定**:诉求形如"能不能通过 TF 实现 <某功能>",且我们不确定上游 API 有无。走 Step 2 决策树分支 B 前的预筛。

```
请问一下云产品的研发,针对您描述的需求,当前是否有接口支持?如果有,通过
哪个接口 / 哪个属性实现的,请提供接口文档链接;如果没有接口支持,那 TF
也无法支持。
```

云产品回帖后:API 支持 → 走分支 C 我们接入;API 不支持 → 走分支 F 上游缺口。

### 场景 4 — 控制台改动导致 TF diff

**判定**:客户报 `terraform plan` 出现非预期 diff,排查发现是控制台上手动改过 TF 管理的资源。

```
Terraform 是全生命周期管理的,除非特殊情况,不建议在控制台上对 TF 创建的
资源进行修改操作;这样会造成本地状态与云上实际数据不一致,从而 plan 出
diff。建议要么统一在 TF 侧管理,要么将控制台的最新数据 import 回 TF state
或 apply 覆盖。
```

### 场景 5 — `hashicorp/alicloud` vs `aliyun/alicloud` source

**判定**:客户问两个 source 有什么区别、该用哪个。

```
两个 source 背后对应的是同一套 provider 代码:
- `hashicorp/alicloud` 是历史路径,当时 provider 在社区 GitHub org 维护;
- 后来 transfer 到我们自己维护,新增了 `aliyun/alicloud` 路径。
- 为保持兼容,`hashicorp/alicloud` 依然可用,功能与 `aliyun/alicloud` 完全一致。
- 更推荐使用 `hashicorp/alicloud`。
```

### 场景 6 — 本地无法复现客户问题

**判定**:客户报 error 但我们本地跑同 tf 不复现;需要客户抓完整 debug 日志。

```
麻烦配置一下环境变量,重新跑一次并把生成的 terraform.log 发上来:

export TF_ACC=1
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log
export DEBUG=terraform
```

### 场景 7 — `Post "https://xxx.aliyuncs.com/?..."` 连接超时

**判定**:报错形如 `Post "https://<product>.aliyuncs.com/?AccessKeyId=...": dial tcp: i/o timeout` 或 `context deadline exceeded`。网络问题,非 TF/API 本身。

```
网络问题,连接超时了,两种解决方案:
1. 切换网络环境 / 挂代理后重试;
2. 显式配置 endpoint。endpoint 值向对应云产品确认(region 为地域名如
   cn-hangzhou 的内网 endpoint):

provider "alicloud" {
  region     = var.region
  access_key = ""
  secret_key = ""
  endpoints {
    <product_key> = "<endpoint 值>"
  }
}
```

### 场景 8 — 「为什么 TF 不支持某功能」科普

**判定**:客户抱怨 TF 缺某功能,需要澄清职责边界(与场景 3 不同:这里更偏理念澄清,而不是我们代问)。

```
Terraform 只是一个客户端资源编排工具,所有功能特性都基于云产品的接口
实现。云产品接口如果支持,且 TF 评估之后认为可以接入,那 TF 就可以支持;
如果云产品本身接口不支持,那 TF 也无法支持。
```

### TF 文档三件套(引导客户自学)

对不熟 TF 语法/用法的客户,先引导读文档,别兜住所有基础问题:

- TF 官方文档: https://developer.hashicorp.com/terraform/language
- TF 中文文档: https://lonegunmanb.github.io/introduction-terraform/
- 阿里云 provider 资源文档: https://registry.terraform.io/providers/aliyun/alicloud/latest/docs


## Step 2 — 定位缺口(按决策树)

详见 [镇元查证与路由分支](../../provider-resource-dev/references/zhenyuan-verification.md)(跨 skill 单点维护)。

## Step 3 — 执行路由动作(写操作,先授权)

### 前置 Gate — 评论区/状态变化扫描

Step 1 读单 → Step 2 查证期间存在**时间窗口**：原指派人(常被前线随机指派)或团队成员可能已回帖修复/贴 PR/接手。执行任何路由写操作**前**,必须 point-read 一次:

```bash
# 1. 扫最新评论(限最近 10 条,看有没有 Fixed/PR 链接/@接手)
bin/a1id -- project workitem comment list <源工单ID> 2>&1 | tail -20

# 2. 看状态/指派人是否已变(如 New → Fixed,或已改指派给过载/谜拟等)
bin/a1id -- project workitem get <源工单ID> -f json 2>/dev/null | \
  python3 -c 'import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
  if f.get("identifier") in ("status","assignedTo"):
    print(f"{f[\"label\"]}: {f.get(\"displayValue\",\"\")} ({f.get(\"value\",\"\")})")'
```

**短路条件**(必须**同时满足**才短路):
1. 评论含团队成员(新山/谜拟/过载/临钧等)贴 PR 链接或明确说"已修复"
2. **且**该 PR/修复命中主工单的客户问题根因(不是只修自己负责的那部分)

**不算短路的情形**(仍需建关联单):
- 状态非 New —— 主单可能需多节点串行(谜拟修完→转临钧),各人修自己的关联单,主单同步状态,全部 Fixed 才最终 Fixed
- 指派人已为路由目标人 —— 前线可能随机指对,仍需关联单追踪进度
- 评论只是 @接手/讨论/追问 —— 未真正修复主问题

**短路动作**:
1. 评论确认对方 PR 命中客户问题根因,贴 PR 链接
2. 不建关联单(避免重复)
3. 原单状态跟进 PR 进度(未合并 → 问题解决中;已合并 → 已发布待需求方验收)
4. wrap done + release

**未命中短路** → 继续下方对应分支执行正常路由。

### 分支 A / D-过载(手写) / E(jarvis 手动建关联单+指派)

```bash
# 0. 从原单读优先级 + DDL,关联单继承
#    - 缺陷类型(功能缺陷/线上问题/性能瓶颈)优先级一律覆写为"紧急",无视原单值
#    - 非缺陷类型直接复制原单优先级
#    - 截止日期:原单 DDL 提前 2 天(留余量给下一棒);原单无 DDL 时 today+3
src_json=$(bash bootstrap/aone-get.sh <源工单ID>)
src_type=$(echo "$src_json" | python3 -c 'import json,sys
d=json.load(sys.stdin)
print(d.get("workitemType","") or d.get("categoryIdentifier",""))')
src_prio=$(echo "$src_json" | python3 -c 'import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
  if f.get("identifier")=="priority": print(f.get("value") or "")')
# 缺陷类型强制紧急
case "$src_type" in *缺陷*|*Bug*|*bug*) src_prio="紧急" ;; esac
src_ddl=$(echo "$src_json" | python3 -c 'import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
  if f.get("identifier")=="80": print(f.get("value") or "")')
new_ddl=$(python3 -c "
from datetime import date,timedelta
src='$src_ddl'.strip()
if src:
  d=date.fromisoformat(src)-timedelta(days=2)   # 短于原单 2 天,给下一棒留余量
  today=date.today()
  if d<=today: d=today+timedelta(days=1)         # 防倒挂
else:
  d=date.today()+timedelta(days=3)               # 原单无 DDL,默认 today+3
print(d.isoformat())")

# 1. 建关联单在 terraform-alicloud (528766)
#    正文用 --body 或 --body-file(a1 不吃 --description,会报 unknown flag)
#    tf_provider(528766) 校验必填:计划开始/截止日期 + 实际工时,通过 --cfs 传
bin/a1id -- project workitem create \
  --project 528766 \
  --category <req|bug|task> \
  --title "<清晰标题:资源/属性/诉求>" \
  --assignee <工号> \
  --priority "$src_prio" \
  --body-file <path-to-body.txt> \
  --cfs "计划开始日期=$(date +%Y-%m-%d)" \
  --cfs "计划截止日期=$new_ddl" \
  --cfs "实际工时=0" \
  --quiet
# --quiet 输出只有 "<id>\t<title>\t<status>\t<assignee>",取第一列作 NEW_ID
NEW_ID=<新单 id>

# 2. 关联(aone 自动双向,单次 relation add 即建 A↔B;第二次会 400 已存在)
bin/a1id -- project workitem relation add <源工单ID> relate:$NEW_ID

# 3. 源工单指派 + 状态(原单优先级 / DDL 保持不动)
bin/a1id -- project workitem update <源工单ID> --assignee <工号>
bin/a1id -- project workitem update <源工单ID> --status 问题解决中
```

**分支 G · Provider 全局改造(→ 新山 521957)**:适用于"诉求不涉及单一 alicloud_xxx 资源、而是 provider 侧全局改动"的场景(region 白名单/框架 utility/公共 endpoint/provider.go 基础/SDK bump 等)。**落地脚本与上面完全一致**,只需 3 处微调:
- `--assignee` 填 `521957`(新山)
- `--category` 填 `task`(全局改造多为工程任务,而非缺陷/需求)
- `--title` 与 `--body-file` 正文注明"provider 全局改造"字样(例:"provider 支持 ap-southeast-8 region"),便于新山识别范围
- 源单 `--assignee 521957` + `--status 问题解决中`

分支 G **不走镇元查证**(镇元管资源 schema,不管 provider 基础),Step 2 分支 D/E 的 acube 覆盖度检查可跳过,直接进入 Step 3 建单流程。

**文档改造分支 · 仅 website/docs 变更(→ 过载 484483)**:适用于"诉求仅涉及文档改造(website/docs 变更,无 provider 代码改动)"的场景。**落地脚本与分支 A/D-过载/E 完全一致**,只需以下微调:
- `--assignee` 填 `484483`(过载)
- `--category` 视原单类型(需求/缺陷)
- `--title` 注明"文档改造"字样
- **不走镇元查证**,直接跳过 Step 2 进入 Step 3 建单
- 过载的关联单 jarvis 直接 claim 跟进(与 D-过载手写分支一致的 bookend 流程)

### 分支 D-临钧(生成器产出):走 acube V2 接口,jarvis 不手动建单

生成器产出资源交给 acube 的 `TerraformVendorBuildTaskOpenapiController#createBuildTaskV2` 接口——接口内部**自动**在 terraform-alicloud (528766) 建关联单、指派临钧(429768)、触发生成/PR 工作流,jarvis 只负责查回 aoneId 并做源单关联+指派。**严禁**同时走上面 `a1 workitem create` 手动建单流程,否则双单污染临钧队列。

服务端实现见邻仓 `a-cube-aliyun-com`:
- `POST /api/v1/terraform_vendor_build/createBuildTaskV2` — body `TerraformVendorBuildTaskDTO`,返回 `ResultDTO<Long>` (taskId,同步返回)
- `GET  /api/v1/terraform_vendor_build/queryAoneByTaskId?taskId={taskId}` — 返回 `{taskId, aoneId, aoneUrl}`,aoneId 异步产生(acube 内部建单完成后回写),需轮询

```bash
# 0. 拿 jarvis 工号(acube 侧 workId/workName 用当前 a1 身份,便于事后追溯)
jarvis_empid=$(bin/a1id -- auth whoami 2>/dev/null | python3 -c '
import sys,re,json
raw=sys.stdin.read()
# 兼容 whoami 输出的多种格式,取第一个 5-11 位数字/WB 前缀作为工号
m=re.search(r"\b(WB\d+|\d{5,11})\b", raw)
print(m.group(1) if m else "")')
[ -z "$jarvis_empid" ] && echo "jarvis 未登录 a1(bin/a1id login jarvis),阻断" && exit 1

# 1. 触发 build 任务(acube 自动建单+指派临钧+跑生成器)
#    必填字段: namespace / resourceTypeCode / resourceTypeVersion / osType / flowType / workId / workName
#    resourceTypeVersion 走"生成器产出待跑"场景填 0.0.0(acube 会跑首版生成)
task_id=$(curl -s -X POST "https://acube.aliyun-inc.com/api/v1/terraform_vendor_build/createBuildTaskV2" \
  -H "Content-Type: application/json" -H "accept: */*" \
  -d "{
    \"namespace\":\"<product>\",
    \"resourceTypeCode\":\"<PascalCase Resource>\",
    \"resourceTypeVersion\":\"0.0.0\",
    \"osType\":\"Linux\",
    \"flowType\":\"ACubeRelease\",
    \"workId\":\"$jarvis_empid\",
    \"workName\":\"jarvis\"
  }" | python3 -c '
import json,sys
d=json.load(sys.stdin)
if d.get("code")!="SUCCESS":
  sys.stderr.write(f"createBuildTaskV2 failed: {d.get(\"code\")} {d.get(\"message\")}\n"); sys.exit(1)
print(d.get("data"))')
[ -z "$task_id" ] && echo "acube createBuildTaskV2 未返回 taskId,阻断" && exit 1
echo "taskId=$task_id"

# 2. 轮询查 aoneId(taskId 立返,aoneId 需等 acube 异步建单完成;60s 内应有值)
NEW_ID=""
for i in 1 2 3 4 5 6; do
  NEW_ID=$(curl -s "https://acube.aliyun-inc.com/api/v1/terraform_vendor_build/queryAoneByTaskId?taskId=${task_id}" \
    -H "accept: */*" | python3 -c '
import json,sys
d=json.load(sys.stdin); data=(d.get("data") or {})
print(data.get("aoneId") or "")')
  [ -n "$NEW_ID" ] && break
  sleep 10
done
if [ -z "$NEW_ID" ]; then
  echo "acube 60s 内未返回 aoneId,升级 escalation/(不要回退到手动 a1 workitem create,可能双建)"
  bootstrap/log.sh escalate <源工单ID> "acube build task $task_id 60s 内未返回 aoneId,人工排查"
  exit 1
fi
echo "临钧关联单 aoneId=$NEW_ID"

# 3. 关联到源客户单(aone 自动双向,单次 relation add 即建 A↔B)
bin/a1id -- project workitem relation add <源工单ID> relate:$NEW_ID

# 4. 源工单同步指派临钧 + 状态(评论 @临钧 走 Step 4 模板 B)
bin/a1id -- project workitem update <源工单ID> --assignee 429768
bin/a1id -- project workitem update <源工单ID> --status 问题解决中
```

**环境**:
- 正式走 `acube.aliyun-inc.com`,预发把域名换成 `pre-acube.aliyun-inc.com`(路径/参数/返回结构一致)
- `/api/v1/**` 免鉴权,内网 DNS(需办公网/VPN)

**关键纪律**:
- acube 自动建单+指派+触发工作流是**原子动作**,jarvis 只做"查 aoneId + 关联源单"善后
- 60s 内没查到 aoneId → 直接升级 escalation,**禁**回退手动 `a1 workitem create`,双建会污染临钧研发队列
- workId/workName 填当前 jarvis 身份工号,acube 侧任务日志能追到调用方

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
跟进,源工单同步指派并改状态为「问题解决中」。

### 查证依据
1. **镇元**:<get: 有/无 data | list released 是否命中 | CoverageScore=<x>>
   - acube V2 覆盖度(线上):CoverageScore=<x>
   - Property/Operation/PrimaryOperation=<>/<>/<>
2. **provider 代码类型**:<自动生成 / 手写>(依据 `alicloud/resource_...go:1-3`
   <是否有 generated automatically 注释>)
3. **紧急度**:priority=<>, 计划截止=<>, 剩余=<>天,故路由到 <人>

### 关联单
- <NEW_ID>: <标题>,项目 528766 terraform-alicloud
  (临钧场景:aoneId 由 acube V2 createBuildTaskV2 接口异步创建;taskId=<>)
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

- 单行可用位置参数:`wrap.sh done <id> "<summary>" <status|--no-status>`
- 多行必须用 heredoc/stdin 或文件:`wrap.sh done <id> --summary-stdin <status|--no-status>` / `wrap.sh done <id> --summary-file <path> <status|--no-status>`
- 仍不支持 `--status` 命名参数;status 放在最后一个位置参数
- 字面量 `\n` 默认会被 wrap.sh 拦截;确认要发送反斜杠+n 文本时才临时设 `JARVIS_ALLOW_LITERAL_NEWLINE=1`
- **status 枚举随工单类型走两套流**:需求/任务类是池自定义状态(tf_customer 合法值见下方反模式清单);bug 类(功能缺陷/线上问题/性能瓶颈)是 Aone 缺陷独立枚举 `Open/Fixed/Won'tfix/Later/Worksforme/Duplicate/Invalid/External/ByDesign`——给 bug 单传「问题解决中」会报 `unsupported target status`,承接中传 `Open`,修复合入后传 `Fixed`(工单 83679740 实测)
- 关联单(528766)claim 规则:**指派给过载(484483)的,jarvis 直接 claim 跟进解决,bookend 同时处理客户主单与关联单**(研发细节 wrap 关联单,客户主单只 wrap 关键节点,收尾两边各自 done+release);指派其他人的不 claim,建单 + @对方即可
- 分支 F 的评论 + 状态修改不需要建关联单,直接 wrap.sh done + 单独一次 status update

## 反模式

- ❌ 只按标题查证,不读 description 末段"限制/差异/仍需" —— 经典误诊(SSL Update 案就栽这)。参见 memory `read-description-last-paragraph`
- ❌ 直接把"类比产品支持 X"当"我们也应该做" —— 必须先分清是**API 原生**还是**Provider 侧适配**,前者 100% 是上游产品缺口
- ❌ 上游 API 缺口场景建了关联单 —— 应该只 @提单人 + 待上游排期,别拖谜拟/过载/临钧下水
- ❌ 专属名单产品还去查镇元覆盖度 —— 不接镇元,直接指派
- ❌ 花名 @ 不带工号(`@谜拟`) —— a1 有时能补,有时不能,显式 `@谜拟(479782)` 保险
- ❌ 建完关联单调两次 `a1 relation add`(A→B + B→A) —— aone 自动双向,第二次会 400 "关联失败该条记录已存在";单次 `add <源> relate:<新>` 即建 A↔B
- ❌ 状态用 `--status 已完成` / `方案功能已存在` 兜底 —— 前者不在合法值,后者语义错(客户真诉求还没解决)
- ❌ 强行 `--assignee` 指派专属名单以外的产品到过载/谜拟 —— 违反分工表,让本团队背不该背的锅
- ❌ 跳过 Step 1.5 共通 gate 直接发 canned —— 分类误建 / 重复单情形下会与承接单重复打搅客户
- ❌ 生成器产出(临钧)场景还手动 `a1 workitem create` 建关联单 —— acube V2 createBuildTaskV2 已自动建单+指派临钧,重复建会双单,污染临钧研发队列
- ❌ acube 60s 未返回 aoneId 就"降级"回手动 `a1 workitem create` —— 可能 acube 已建成功只是查询未及时,回退会双建;正确做法是升级 escalation 由人排查
- ❌ Provider 全局改造(region 白名单/框架 utility/公共 endpoint/SDK bump)走"镇元 OK/NOT OK"判定 —— 镇元管资源 schema,不管 provider 基础;直接分支 G → 新山(521957),不必查 acube 覆盖度
- ❌ 只按 acube V2 `CoverageScore==1.0` 判"镇元 OK",忽略"当前 schema 属性是否覆盖客户诉求字段" —— 覆盖度分只反映已建 schema 的属性测试完备度,客户想要新字段而 schema 未建时,覆盖度分再高也是 NOT OK(缺口在镇元)
- ❌ 状态改成"问题处理中" —— tf_customer(1086837) 池合法枚举没这个值,合法名是"问题解决中";写错 a1 会 `unsupported target status` 阻断。合法枚举:需求待补充/待处理/评估中/待上游排期/问题讨论/长期跟进/待排期/已排期/问题解决中/已发布待需求方验收/验收中/验收通过/验收不通过/客户未响应/方案功能已存在/需求撤回/已拒绝
- ❌ 转单不复制原单优先级 / 不设短于原单 DDL 的截止日期 —— 关联单接手方无优先级参考,DDL 与原单齐会让下一棒无余量;规则:`--priority` 复制原单,`--cfs 计划截止日期` = 原单 DDL - 2 天(至少 today+1);原单无 DDL 时默认 today+3
- ❌ 建关联单用 `--description` —— a1 CLI 不吃(报 `unknown flag: --description`),正文用 `--body` 或 `--body-file`
- ❌ 在 tf_provider(528766)建单不传"计划开始日期 / 计划截止日期 / 实际工时" cfs —— 池校验必填,漏传会 400 `【计划开始日期】不能为空...`;用 `--cfs "计划开始日期=YYYY-MM-DD"` 等传入
- ❌ Provider 源码查证跳过 Step 2 前置的 upstream PR 前扫,只 grep 本地 workspace 磁盘 —— workspace 的本地 branch 可能滞后 upstream 数十小时,或 sync-provider.sh 未 hardened 时只 fetch 不 reset;必先跑 `gh pr list --search` + `gh api contents?ref=master`,同题 recently-merged PR 直接引用避免重复建单;参见工单 83718139 教训(PR 9909 merged 21h 后 jarvis 才处理仍未命中)
- ❌ Step 3 执行路由前不扫评论区/状态 —— 查证+决策期间有时间窗口,原指派人(常被前线随机指派)可能已回帖修复/贴 PR/接手;不扫就转单会重复建关联单、让接手方收多余通知;参见工单 83861367 教训(新山在 jarvis 转单前 20 分钟已修并评论,jarvis 仍建单 @过载)
- ❌ 仅文档改造的工单还走镇元查证 —— 文档改造不涉及 provider 代码/schema,镇元验证无意义;直接跳过 Step 2,关联单指派过载(484483)
- ❌ 缺陷类型工单沿用原单 `priority` 字段值 —— 缺陷(功能缺陷/线上问题/性能瓶颈)优先级一律覆写为"紧急",无视原单标记;覆写后影响镇元 NOT OK 分支路由(紧急→新山)及关联单 `--priority` 传值
