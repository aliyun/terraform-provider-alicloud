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
├─ YES → 仍走镇元查证(确保文档与 schema/API 一致),路由固定            │
│        关联单指派 过载(484483),jarvis claim 跟进 [end]               │
│        镇元查证发现 metadata 问题 → 额外建镇元 agent 侧单(2165097)   │
└─ NO ↓                                                              │
                                                                     │
诉求是"现有资源缺 X 属性/值/行为"(有 analog 产品/资源可类比)?         │
├─ YES → 查 analog 是"API 原生"还是"Provider 侧适配":                │
│         ├─ API 原生 + provider 透传 → 缺口在上游产品 API →         │
│         │   @提单人 + status=待上游排期 [end]                       │
│         └─ Provider 侧适配 → 我们团队可复制路径 → 走下方镇元分支    │
└─ NO (纯接入新资源 / 修 provider bug) ↓                             │
                                                                     │
纯 datasource 问题?(诉求只涉 data.alicloud_xxx 查询/过滤/输出字段,    │
  不涉资源 schema/生命周期;resource+datasource 混合不算"纯")          │
├─ YES → 【与镇元不相关】——datasource 是 provider 侧对查询 API 的     │
│   只读封装,镇元只管资源 schema,查镇元无意义;**跳过镇元查证**,       │
│   直接走下方【与镇元不相关分流】                                     │
└─ NO(资源类) ↓                                                       │
                                                                     │
镇元 OK? (三条件全满足才算 OK,**手写资源亦严格适用,不豁免**:            │
  ① API 在镇元已定义并发布                                            │
  ② 当前资源 schema 属性满足客户诉求(不缺字段)                        │
  ③ acube V2 覆盖度 CoverageScore == 1.0)                             │
├─ NO (任一不满足) = 与镇元相关且镇元 NOT OK [分支 E] ↓               │
│     关联单 → 镇元 agent(WORKER_1783326253279), 落 Terraform镇元    │
│     对接(2165097) 池 · body 必带 `## 机读信息` JSON(见 templates.md)│
│     ——镇元侧根因主责,agent 自动接单,**无论紧急与否都建**          │
│     再判紧急: 优先级=紧急 OR 距计划截止 < 14 天?                     │
│     ├─ YES → **同时再建一张**关联单 → 新山(521957), 落 528766      │
│     │        ——双单并行(agent 修镇元侧根因,新山紧急兜底 provider) │
│     └─ NO  → 仅 agent 单一张                                        │
│     原单指派 谜拟(479782,人类 owner) + status=问题解决中           │
│     + @谜拟(紧急时并 @新山)——钉钉私信只发谜拟/新山,不发 agent      │
└─ YES = 镇元侧无问题、provider 侧存在问题 → 【与镇元不相关】↓        │
                                                                     │
【与镇元不相关分流】(纯 datasource 从上方直接进入):                    │
├─ 资源类 且 provider 代码=自动生成                                    │
│   (`This file is generated automatically`) →                       │
│   关联单指派 临钧(429768) [D-临钧,acube 重跑生成器管道不变]         │
└─ 其余(纯 datasource / 手写代码):                                    │
    ├─ 紧急 → 关联单指派 新山(521957) [D-新山]                        │
    └─ 不紧急 → 关联单指派 过载(484483) [D-过载]                      │
    原单同步指派 + status=问题解决中 + @指派人                         │
    过载(484483)关联单 → jarvis claim 跟进 + 评论跟进度                 │
                                                                     │
以上所有分支均未匹配的特殊情况(NPE 兜底)                              │
→ 关联单指派 夏节(401498) + 打标签 jarvis-npe [分支 H,end]           │
```

**「与镇元相关性」定义**(2026-07-06 路由重整的核心概念):以下两类为**与镇元不相关的问题**——① **纯 datasource 问题**:诉求只涉 `data.alicloud_xxx`(查询/过滤/输出字段),不涉资源 schema/生命周期。datasource 是 provider 侧对 List/Describe 类查询 API 的只读封装,镇元(Cloudspec)只管资源 schema,查镇元无意义(跳过镇元查证);resource+datasource 混合诉求不算"纯",按资源主线走。② **镇元侧无问题、provider 侧存在问题**:镇元 OK 三条件全满足,缺口在 provider 实现(bug/适配缺失/文档行为不符)。与镇元不相关的问题**不进 2165097 池**——紧急转新山(521957),不紧急转过载(484483);生成器产出资源例外走临钧 acube 管道。反之,**2165097 池镇元 agent 只接「与镇元相关且镇元 NOT OK」**的单;该类单若紧急,则镇元 agent 关联单 + 新山关联单**两张都建**并行处理。

**过载 = jarvis 自动接手**(2026-07-09 新增):所有路由到**过载(484483)**的关联单(文档改造分支 + D-过载不紧急分支),jarvis **直接 claim 关联单干活**(worktree 开发 / PR / 合入),不等过载本人处理。评论区跟进进度(每步动作都发评论:查证结论 → PR 链接 → 合入状态)。**bookend 纪律不变**:claim → 干活 → wrap done → release/finish。过载是团队共享的"jarvis 工作账号"——路由到过载等于 jarvis 自己接单闭环。**钉钉私信仍发过载(484483)**(建关联单时通知,通用规则不变)。

**镇元 agent 接单机制**(2026-07-08 谜拟本人切换 · 关键变更):分支 E 的 Terraform镇元对接(2165097) 关联单不再指派谜拟本人,而是指派 **镇元 agent (`WORKER_1783326253279`)** 自动接单——谜拟自己不解单,由 agent 从关联单 description 里的机读 JSON 解析后驱动 spec 补齐 / 映射建立 / 生成器触发。**契约硬要求**:关联单 body 必须严格按 [templates.md 的 "Requirement skeleton (Cloudspec 关联单 · 镇元 agent 接单硬契约)"](./templates.md#requirement-skeleton-cloudspec-关联单--镇元-agent-接单硬契约) 骨架写(`## 背景` / `## 需求` / `## 机读信息` + `\`\`\`json` 代码块 + 7 字段全);少任何一项 agent 都无法接,单会沉底。谜拟(479782) **保留在源客户主单的 assignee + 参与者/@ 里** 做**人类兜底 owner**(客户可见),但**关联单 assignee 不再是她**;钉钉私信也只发谜拟/新山,不发 agent 工号(`WORKER_` 前缀无 IM 通道)。

**手写资源不豁免镇元 OK 三条件**(2026-07-08 84021197 事后补规):即使主资源是手写代码(非生成器产出),镇元 spec 仍是团队合同——面向长期"手写→生成器迁移"与文档一致性,条件②③(schema 属性覆盖 / acube 覆盖度)照样评估。**反模式**:以"手写不走生成链路,镇元 spec 缺字段非阻塞"为由把手写资源的镇元缺字段判为"与镇元不相关",直接转分支 D-过载/新山。**正确做法**:镇元 NOT OK 一律走分支 E 转镇元 agent (WORKER_1783326253279) 落 2165097 池主责镇元 spec 补齐;若同时 provider 手写代码也需补透传,分支 E 的紧急双单机制不适用于"不紧急"情形——不紧急时只 agent 一张,provider 手写补齐作为镇元 spec 落地后的自然跟进,或视场景再挂一张 tf_provider 关联单跟研发(不重复"紧急"标签)。

关联单默认建在 **terraform-alicloud** 项目(528766, `pools.tf_provider`),类型 = 缺陷/需求(视诉求),双向关联到源客户单。**四类例外**(不走 `a1 workitem create` 手动建 tf_provider 关联单):

1. **专属维护名单产品(分支 A)不建关联单**——ACK/SLS/OSS/RDS 等专人维护的云产品,provider 代码就是该负责人自己的活,源单直接 `--assignee <工号>` + `--status 问题解决中` + `@负责人`,追踪走源客户单本身;jarvis 无 tf_provider 关联单参与研发闭环。**若前次已按分支 A 转对(源单 assignee 已是名单人),不要"补建关联单"**——分支 A 本身就没关联单,不存在"缺"。
2. **分支 E 关联单落 Terraform镇元对接(2165097, `pools.upstream.cloudspec_gap`) 指派镇元 agent (`WORKER_1783326253279`)**——镇元侧根因由 agent 自动接手在自己池内跟进(谜拟做人类兜底 owner,不再自己解单),tf_provider(528766) 池只放 provider 侧修复(如分支 E 紧急双单里的新山单)。**分支 E 双单时两池并存**:agent 单 → 2165097,新山单 → 528766。**agent 接单硬契约**:关联单 body 必须按 [templates.md 的 Cloudspec 关联单 · 镇元 agent 接单硬契约](./templates.md) 骨架写足 7 字段机读 JSON,缺字段/marker/JSON 语法错 = agent 无法接单 = 单沉底。
3. **临钧路由(生成器产出)不由 jarvis `a1 workitem create` 手动建单**——走 acube `createBuildTaskV2` 接口,acube 内部自动建单+指派临钧+触发生成/PR 工作流,jarvis 只查回 aoneId 做关联,详见 Step 3。
4. **分支 F(上游 API 缺口)不建关联单**——只 @提单人 + status=待上游排期,详见 Step 3。


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

**缺陷类型优先级覆写**:`workitemType` 为缺陷(功能缺陷/线上问题/性能瓶颈)时,**优先级一律视为紧急**,无视原单 `priority` 字段值。此覆写影响决策树两处紧急判定——分支 E(与镇元相关且 NOT OK)紧急时在镇元 agent 单之外**加建新山双单**;与镇元不相关分流紧急时路由新山(521957)——及关联单 `--priority` 传值。

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

### 钉钉私信 · 通用调用姿势(所有实质动作)

Aone 评论 + @花名(工号) 仍是**主通知渠道**;钉钉私信是**补充**,让承接方/提单人在 Aone 通知易被漏看时也能第一时间知道。**jarvis 做实质动作就私信相关方**:

| 动作 | 私信对象 |
|---|---|
| **转单**(分支 A/B/C/D-新山/D-过载/D-临钧/G/H) | 承接方(新指派人) |
| **转单**(分支 E 镇元 agent 关联单) | **谜拟(479782)——agent 归谜拟维护,`WORKER_` 前缀无 IM 通道,`notify-dingtalk.sh` 传 agent id 会 400/静默,一律改私信谜拟**;紧急双单场景另加私信新山(521957) |
| **分支 F 上游 API 缺口** | 提单人(阿里前线/PD) |
| **补建关联单**(前次漏建) | 承接方(前次没关联单不知道);**分支 E 私信谜拟**(同上,不发 agent) |
| **模板 D 进度跟进**(≥30 天无实质进展) | 承接方(Aone 明显没在看,私信更有效);**分支 E 承接方是 agent → 私信谜拟兜底催 agent** |
| **模板 F 补料提醒** | 提单人 + 承接方(球回客户手里,承接方知会可歇);**分支 E 承接方位替换成谜拟** |
| **模板 E 关单提示** | 提单人 + 最后处理人;**分支 E 最后处理人是 agent → 用谜拟顶替 agent 位** |

**不发**:观察等待(<30 天,承接方正在处理,不打搅)。

**通用规则 · agent 承接方一律用谜拟顶替**:凡承接方是镇元 agent(`WORKER_1783326253279`)的场景,私信对象一律替换为谜拟(479782)——agent 归谜拟维护、`WORKER_` 前缀无 IM 通道、`notify-dingtalk.sh` 传 agent id 会 400/静默。不要私信 agent 工号,永远走人类兜底 owner。

提单人一定是**阿里前线/PD**(不是外部客户),内部账号可以私信。

**调用**:

```bash
# 单行 body
bash bootstrap/notify-dingtalk.sh <staffId> "<title>" "<body>"

# 多行 body(推荐,避免 shell 转义)
bash bootstrap/notify-dingtalk.sh <staffId> "<title>" --body-file /tmp/notify-<id>.txt
# 或
cat <<'EOF' | bash bootstrap/notify-dingtalk.sh <staffId> "<title>" --body-stdin
第一行
第二行
EOF

# dry-run(不实际发,打印骨架用来 debug 消息内容)
bash bootstrap/notify-dingtalk.sh --dry-run <staffId> "<title>" "<body>"
```

**消息内容规范**——**只贴关键索引**,不复述 Aone 评论全文(双通道内容重复=噪声):
- **一句结论**(接了什么/关了什么)
- **分支说明**(为什么走到你名下,如"分支 D-过载·纯 datasource·非紧急")
- **Aone 主单链接** `https://project.aone.alibaba-inc.com/v2/project/1086837/req/<主单ID>`
- **关联单链接**(若有,如 tf_provider 关联单 528766 / 镇元 agent 关联单 2165097 / acube 自动建单)

**失败降级**(`notify-dingtalk.sh` 自身兜住,不阻断 bookend):
- 缺 `DINGTALK_APP_KEY/SECRET/TEMPLATE_ID` → stderr `[NOTIFY-SKIP: missing env]`,退 0
- `JARVIS_NOTIFY_DINGTALK=0` → 全局关闭,退 0
- staffId ∈ `config/dingtalk-optout.txt` → 个人 opt-out,退 0
- 网络/API 失败 → 落 `escalation/notify-fail-<ts>-<staffId>.md`,退 0

**staffId 假设**:阿里 empId(Aone 工号)= 钉钉 staffId,直接用 team-roster 里的工号即可。首次全量启用前先跑 dry-run 三人小样(仓库主人 + 过载 + 新山)确认能收到。

### 前置 Gate — 评论区/状态变化扫描

Step 1 读单 → Step 2 查证期间存在**时间窗口**：原指派人(常被前线随机指派)或团队成员可能已回帖修复/贴 PR/接手。执行任何路由写操作**前**,必须 point-read 一次:

```bash
# 1. 扫最新评论(限最近 10 条,看有没有 Fixed/PR 链接/@接手)
bin/a1id -- project workitem comment list <源工单ID> 2>&1 | tail -20

# 2. 看状态/指派人是否已变(如 New → Fixed,或已改指派给过载/谜拟/镇元 agent 等)
bin/a1id -- project workitem get <源工单ID> -f json 2>/dev/null | \
  python3 -c 'import json,sys
d=json.load(sys.stdin)
for f in d.get("fields",[]):
  if f.get("identifier") in ("status","assignedTo"):
    print(f"{f[\"label\"]}: {f.get(\"displayValue\",\"\")} ({f.get(\"value\",\"\")})")'
```

**短路条件**(必须**同时满足**才短路):
1. 评论含团队成员(新山/谜拟/过载/临钧/镇元 agent 等)贴 PR 链接或明确说"已修复"
2. **且**该 PR/修复命中主工单的客户问题根因(不是只修自己负责的那部分)

**不算短路的情形**(仍需建关联单):
- 状态非 New —— 主单可能需多节点串行(镇元 agent 修完 spec→转临钧生成),各人/agent 修自己的关联单,主单同步状态,全部 Fixed 才最终 Fixed
- 指派人已为路由目标人 —— 前线可能随机指对,仍需关联单追踪进度
- 评论只是 @接手/讨论/追问 —— 未真正修复主问题

**短路动作**:
1. 评论确认对方 PR 命中客户问题根因,贴 PR 链接
2. 不建关联单(避免重复)
3. 原单状态跟进 PR 进度(未合并 → 问题解决中;已合并 → 已发布待需求方验收)
4. wrap done + release

**未命中短路** → 继续下方对应分支执行正常路由。

### 分支 A(专属维护名单产品:只指派 + @,**不建关联单**)

专属维护名单(ACK/SLS/MNS/OSS/ESS/OTS/EMR/RDS/PolarDB/MSE/ClickHouse,见 [team-roster.md](./team-roster.md))的产品诉求由该云产品 provider 专人维护,**不接镇元、不建 tf_provider 关联单**,追踪走源客户单本身。落地只三步:

```bash
# 1. 源单指派专属维护人
bin/a1id -- project workitem update <源工单ID> --assignee <工号>
# 2. 源单状态改「问题解决中」
bin/a1id -- project workitem update <源工单ID> --status 问题解决中
# 3. 走 wrap.sh done 发模板 C 评论并 @负责人;bookend 内 claim → wrap done → release
#    评论正文参见 Step 4 模板 C(不提关联单)
# 4. 钉钉私信承接方(补充通知,失败不阻断)
bash bootstrap/notify-dingtalk.sh <工号> \
  "Aone 转单 · <产品中文名>" \
  "分支 A · 专属维护名单
主单: <一句诉求>
链接: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>"
```

**为什么不建关联单**:专属维护名单本身表示"该产品 provider 由这个人独立维护、不走 tf_provider 共享池",研发在其自己团队/仓库内闭环,建 tf_provider 关联单是重复档案。若该产品出现"专人已推 PR / 已修复"的短路情形,按 Step 3 前置 Gate 短路即可,一样不建单。

### 分支 D-新山 / D-过载 / E(jarvis 手动建关联单+指派)

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

# 1. 建关联单
#    正文用 --body 或 --body-file(a1 不吃 --description,会报 unknown flag)
#    池选择(见 Step 3 开头例外说明):
#      - 分支 E 镇元 agent 单 → --project 2165097 (Terraform镇元对接) --assignee WORKER_1783326253279
#        · 2165097 池当前无 528766 的「计划开始/截止日期 / 实际工时 / Terraform需求类型」
#          cfs 强校验,可省 cfs 直接建;若日后 400 按报错补,别默认沿用 528766 cfs 组合
#        · **body 硬契约**:必须严格按 templates.md 的 "Cloudspec 关联单 · 镇元 agent 接单硬契约"
#          骨架写(`## 背景` / `## 需求` / `## 机读信息` + ```json 代码块 + 7 字段全)
#          缺 marker / 缺字段 / JSON 语法错 = agent 无法接单 = 单沉底
#      - 其它路由(过载/新山/夏节/文档改造/G/H) → --project 528766 (terraform-alicloud)
#        必带下方三件套 cfs;bug 类还必带 --cfs "Terraform需求类型=..."
bin/a1id -- project workitem create \
  --project <528766 或 2165097> \
  --category <req|bug|task> \
  --title "<清晰标题:资源/属性/诉求>" \
  --assignee <工号 或 WORKER_1783326253279> \
  --priority "$src_prio" \
  --body-file <path-to-body.txt> \
  --cfs "计划开始日期=$(date +%Y-%m-%d)" \
  --cfs "计划截止日期=$new_ddl" \
  --cfs "实际工时=0" \
  --quiet
# 注:上述 --cfs 三件套仅 528766 池必需;分支 E 镇元 agent 单落 2165097 时可去掉整块 --cfs
# --quiet 输出只有 "<id>\t<title>\t<status>\t<assignee>",取第一列作 NEW_ID
NEW_ID=<新单 id>

# 分支 E body 生成参考(镇元 agent 硬契约,详见 templates.md):
# cat > /tmp/e-body-<id>.md <<'MD'
# ## 背景
# <原文【背景】段落 · 缺口一句话 + 来源诉求>
#
# ## 需求
# 1. 资源:<Namespace::Type>(alicloud_x) 缺<attr> 对应 <Product>::<Action> 预期透传+import
# 2. Cloudspec:补 <Namespace::Type> 资源 spec,建立 alicloud_x ↔ Cloudspec 映射
#
# ## 机读信息
# ```json
# {
#   "background": "<与上方【背景】段落一致>",
#   "requirement": "<与上方【需求】段落一致>",
#   "documentUrl": "https://api.aliyun.com/document/<Product>/<ver>/<Action>",
#   "mappingCheckUrl": "https://acube.aliyun-inc.com/api/v1/terraform/generator/getTerraformResourceSpec?terraformResourceType=alicloud_x",
#   "acceptance": "Cloudspec 资源 spec 补齐,alicloud_x ↔ Cloudspec 映射建立,getTerraformResourceSpec 查询返回有效映射",
#   "deadline": "$new_ddl",
#   "source": "工单 <源客户单ID>（Terraform - 客户问题）"
# }
# ```
# MD

# 2. 关联(aone 自动双向,单次 relation add 即建 A↔B;第二次会 400 已存在)
bin/a1id -- project workitem relation add <源工单ID> relate:$NEW_ID

# 3. 源工单指派 + 状态(原单优先级 / DDL 保持不动)
#    补建关联单场景(判定表"补建"行)跳过本步:前次已改 assignee/status,不重复改
bin/a1id -- project workitem update <源工单ID> --assignee <工号>
bin/a1id -- project workitem update <源工单ID> --status 问题解决中

# 4. 钉钉私信承接方(补充通知,失败不阻断;补建也发)
#    分支 E 承接方是镇元 agent(WORKER_ 无 IM 通道):
#      · 谜拟单一律私信谜拟(479782)——agent 归谜拟维护,人类兜底 owner
#      · 紧急双单场景另发一条给新山(521957)
bash bootstrap/notify-dingtalk.sh <工号> \
  "Aone 转单 · <一句诉求>" \
  "分支 <D-过载|D-新山|E-agent|E-新山> · <紧急度/纯 datasource/镇元 NOT OK 等>
关联单: #$NEW_ID (<528766 tf_provider|2165097 镇元对接>)
主单: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>
关联单: https://project.aone.alibaba-inc.com/v2/project/<528766|2165097>/req/\$NEW_ID"
```

**分支 E · 关联单(镇元 agent WORKER_1783326253279 + 紧急时并新山 521957)**:与镇元相关且镇元 NOT OK 的单,**镇元 agent 关联单无论紧急与否都建**(镇元侧根因主责,agent 自动接单);若紧急(优先级=紧急 OR 距 DDL<14 天 OR 缺陷类型覆写),**再按同一脚本建第二张**关联单指派新山,两单并行。要点:
- **池归属两分**:镇元 agent 单 `--project 2165097` (Terraform镇元对接) `--assignee WORKER_1783326253279`,新山单 `--project 528766` (terraform-alicloud) `--assignee 521957`。两池 cfs 校验不同——**2165097** 当前无「计划开始/截止/实际工时/Terraform需求类型」强校验,可省整块 `--cfs` 直接建;**528766** 必带三件套 cfs,bug 类还必带 `--cfs "Terraform需求类型=..."`。日后 2165097 若 400 按报错补,别默认沿用 528766 cfs 组合。
- **镇元 agent 单 body 必须机读**:严格按 [templates.md 的 "Requirement skeleton (Cloudspec 关联单 · 镇元 agent 接单硬契约)"](./templates.md) 骨架写(`## 背景` / `## 需求` / `## 机读信息` + ```json 代码块 + 7 字段全)。缺 marker/字段/JSON 语法错 = agent 无法接单 = 单沉底(jarvis 不能替 agent 补做,漏契约就是把单丢黑洞)。
- 两张关联单**各自**双向关联源客户单(`relation add <源> relate:<agent 单>` + `relation add <源> relate:<新山单>`)
- 分工写进各自 `--body`:agent 单注明"镇元侧根因修复(schema 定义/属性补齐/覆盖度)"(用机读契约表达);新山单注明"紧急兜底——provider 侧可先行绕过/加速方案,与 agent 单 <ID> 并行"
- **原单指派谜拟(479782,人类兜底 owner)** + status=问题解决中;评论 `@谜拟(479782)`(紧急时并 `@新山(521957)`),注明"关联单已交镇元 agent 自动处理,agent 断电/复杂决策时 @谜拟兜底"。**钉钉私信只发谜拟/新山,不发 agent 工号**(`WORKER_` 前缀无 IM 通道,notify-dingtalk.sh 传 WORKER_ id 会 400)
- 非紧急:仅镇元 agent 一张(落 2165097),原单指派谜拟,@谜拟

**分支 D-新山 / D-过载 · 与镇元不相关(纯 datasource / 镇元 OK 但 provider 侧问题,手写代码)**:落地脚本与分支 A 完全一致,按紧急度选指派人——紧急 `--assignee 521957`(新山),不紧急 `--assignee 484483`(过载);`--title` 注明问题面("纯 datasource"/"provider 侧实现"),纯 datasource 单在 `--body` 说明已按定义跳过镇元查证。**D-过载**:jarvis 直接 claim 关联单跟进(worktree → PR → 合入 → bookend),评论区逐步跟进度(见上方"过载 = jarvis 自动接手"段);**D-新山**:建单 + @对方。两者都须**钉钉私信承接方**(通用规则 line 434-443)。

**分支 G · Provider 全局改造(→ 新山 521957)**:适用于"诉求不涉及单一 alicloud_xxx 资源、而是 provider 侧全局改动"的场景(region 白名单/框架 utility/公共 endpoint/provider.go 基础/SDK bump 等)。**落地脚本与上面完全一致**,只需 3 处微调:
- `--assignee` 填 `521957`(新山)
- `--category` **跟原单类型一致**(`req`/`bug`/`task`,视原单 `workitemType`)——全局改造工单类型跨度大(需求/任务/缺陷都可能),类型跟原单走避免类型突变导致 Aone 报表统计错位
- `--title` 与 `--body-file` 正文注明"provider 全局改造"字样(例:"provider 支持 ap-southeast-8 region"),便于新山识别范围
- 源单 `--assignee 521957` + `--status 问题解决中`

分支 G **不走镇元查证**(镇元管资源 schema,不管 provider 基础),Step 2 分支 D/E 的 acube 覆盖度检查可跳过,直接进入 Step 3 建单流程。

**文档改造分支 · 仅 website/docs 变更(→ 过载 484483)**:适用于"诉求仅涉及文档改造(website/docs 变更,无 provider 代码改动)"的场景。**落地脚本与分支 A/D-过载/E 完全一致**,只需以下微调:
- `--assignee` 填 `484483`(过载)
- `--category` 视原单类型(需求/缺陷)
- `--title` 注明"文档改造"字样
- **仍走镇元查证**(Step 2 不跳过):probe tier-0 的 `doc_gap_*` 发现依赖 TF 文档 ↔ OpenAPI 文档 ↔ provider 源码三方比对,文档改造同样需要确认镇元 schema 定义准确
- 路由固定过载,**不受镇元 OK/NOT OK 影响**(不像分支 D/E 会因镇元结果分到不同人);jarvis claim 关联单跟进(见"过载 = jarvis 自动接手"段)
- **镇元 metadata 问题 → 额外建镇元 agent 侧单**(2026-07-09 新增):若镇元查证发现 Cloudspec/镇元侧 metadata 也存在同类问题(如枚举值描述过时、属性定义与实际不符),**在主路由关联单之外**,额外建一张镇元 agent 关联单(`--project 2165097 --assignee WORKER_1783326253279`,body 按 templates.md 硬契约)。两张单各管各的:过载单修 TF 文档,jarvis 跟进;镇元 agent 单修 Cloudspec metadata,agent 自动接。**不冲突、不替代**。钉钉私信发谜拟(479782,agent 兜底 owner)

**分支 H · NPE 兜底(→ 夏节 401498 + 标签 `jarvis-npe`)**:适用于"以上所有分支均未命中"的兜底场景。典型触发:
- 需求跨多个云产品无法拆解到单一负责人(如 EMR SparkServerless / StarRocksServerless / DLF 三块独立产品线,本团队无路径拆分)
- 客户诉求超出 tf-alicloud 团队职责边界但不明确该转给谁
- 分类/归属模糊,决策树各分支判定结果均为 N(不适用)
- 决策依据不足以走前述任一分支,又不能直接分支 F 甩上游

**落地脚本与分支 A/D-过载/E 一致**,只需以下微调:
- `--assignee` 填 `401498`(夏节)
- `--category` 一般填 `task`(视原单类型)
- `--title` 注明"NPE 兜底/需二次分诊"字样,便于夏节识别
- 源单 `--assignee 401498` + `--status 问题解决中` + **打标签 `jarvis-npe`**(而非 jarvis-idle/jarvis-claimed):
  ```bash
  bin/a1id -- project workitem tag add <源工单ID> jarvis-npe
  ```
- **不走镇元查证**(NPE 意味着诉求本身还没定位到具体资源/API 层),Step 2 可跳过,直接建单

**为什么走夏节兜底**:夏节是团队分诊 owner,能横向拉云产品/评估拆解路径,能处理"路由决策树没覆盖到"的边缘场景。jarvis-npe 标签便于后续统计分诊决策树盲区,补齐路由规则。

### 分支 D-临钧(生成器产出):走 acube V2 接口,jarvis 不手动建单

生成器产出资源交给 acube 的 `TerraformVendorBuildTaskOpenapiController#createBuildTaskV2` 接口——接口内部**自动**在 terraform-alicloud (528766) 建关联单、指派临钧(429768)、触发生成/PR 工作流。jarvis 只负责:(1) 触发 build 任务,(2) 轮询 `queryAoneByTaskId` 拿 aoneId(60s 内),(3) 关联源单 + 指派临钧 + status 改「问题解决中」。**严禁**同时走 `a1 workitem create` 手动建单流程,否则双单污染临钧队列。60s 内查不到 aoneId → 升级 escalation,**禁**回退手动建单。

**完整 bash 脚本(含 workId/workName 抓取、轮询、关联)、接口 URL/body 详情、环境说明** 见 [acube-createBuildTaskV2-workflow.md](./acube-createBuildTaskV2-workflow.md)。

**钉钉私信临钧**(拿到 aoneId + 关联完成后追加,与其他分支一致):

```bash
bash bootstrap/notify-dingtalk.sh 429768 \
  "Aone 转单 · <一句诉求,如 alicloud_xxx 新资源>" \
  "分支 D-临钧 · 生成器产出
关联单: #\$AONE_ID (528766 tf_provider,acube V2 自动建)
taskId: <acube taskId>
主单: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>
关联单: https://project.aone.alibaba-inc.com/v2/project/528766/req/\$AONE_ID"
```

### 分支 F(上游 API 缺口,只 @提单人)

无需建关联单,只发评论 + 改状态 + 钉钉私信提单人:

```bash
bin/a1id -- project workitem update <源工单ID> --status 待上游排期
# (CLI 报 unsupported 可实际写入,若拒绝改用 "待排期")

# 钉钉私信提单人(补充通知,失败不阻断)
bash bootstrap/notify-dingtalk.sh <提单人工号> \
  "Aone 需转上游 · <一句诉求>" \
  "分支 F · 上游 <产品> API 缺口,Terraform Provider 侧无法闭环
状态: 待上游排期
链接: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>"
```

## Step 4 — 回复模板

### 模板选择判定表(先横向判定动作,再选具体模板)

分诊后要发评论前,先按下表定位**该走哪种动作**;每种动作对应下方一个具体模板:

| 触发条件(必须都命中) | 动作 | 模板 | 写操作 |
|---|---|---|---|
| 新工单 or **前次路由错**(如 78365505 从 NPE 撤回改路由到豁朗) | **转单** | A/B/C(按分支 A-H) | **分支 A(专属维护名单)**:仅源单 assignee 改 + 状态改「问题解决中」+ wrap done + release,**不建关联单**;**分支 B/C(其它)**:建关联单 + 源单 assignee 改 + 状态改「问题解决中」+ wrap done + release。**所有转单动作 + 钉钉私信承接方**(A/D/E/G/H/临钧 每个分支落地脚本尾部各带一步,分支 F 只 @提单人不建单也不私信) |
| **前次路由对 + 缺关联单**(源单 assignee 已改到本团队分工里正确的人,但没有对应关联单——常见于前一轮 jarvis 只改了 assignee/状态却漏建关联单,或客户单历史长且前线随机指派;对应池按分工走:过载/新山/夏节/G/H → 528766,分支 E 谜拟(人类兜底 owner) → 2165097 镇元 agent 关联单) | **补建关联单** | B/C(按分支,不重复改 assignee/状态) | 建关联单 + relation add(双向)+ comment 通知承接方新单号;**源单 assignee/状态保持不动**(前次已对);**钉钉私信承接方**(前次没关联单不知道,现在必须通知;agent 无 IM,分支 E 私信谜拟)。**分支 A(专属维护名单)不适用本行**——A 本身无关联单,前次已按 A 转对即为终态,不再"补建" |
| **前次路由对 + 关联单齐 + 距上次实质进展 <30 天** | **观察等待** | 无 | **不发评论、不改状态、不建单、不私信**;jarvis 内部记录本轮观察时间(bridge revisit 日轮会重扫,新进展进来自然触发);避免频繁打搅承接方 |
| 前次路由对 + 关联单齐 + **距上次实质进展 ≥30 天**(不算 canned/@ 追问,以承接方给出的技术信息或时间承诺为准) | **进度跟进** | D | 只发 comment,免 bookend;不改状态、不建关联单;**钉钉私信承接方**(Aone 明显没在看,私信更有效) |
| 承接方已给结论(拒接/根因/待客户验证)但客户或云产品未回,且已有 ≥1 次追问未响应 | **追料/补料提醒** | F | 只发 comment,免 bookend;提醒里给出建议补齐时间(默认 today+14 天)与到期处理方式(默认"按无法复现/信息不足关单");评论**同时 @ 客户/提单人 + 承接方**(提单人负责补料/确认,承接方获知补料进度);**钉钉私信提单人 + 承接方**(球回客户,承接方可歇) |
| 承接方已闭环(PR merged / 已发布 / 明确拒接结论),但主单状态仍是 New/评估中/问题解决中 | **关单提示** | E | wrap done + 状态改「已发布待需求方验收」/「已拒绝」/「方案功能已存在」+ finish 或 release;评论**同时 @ 提单人 + 最后处理人**;**钉钉私信提单人 + 最后处理人**(见模板 E 段落) |

### assignee 转向规则(补齐,防止隐式推断)

判定表列的"写操作"里 assignee 处置常被忽略,汇总如下(每种动作对应的 assignee 处置):

| 动作 | 源单 assignee 处置 | 依据 |
|---|---|---|
| **转单** | 改**承接方**(过载/新山/谜拟/临钧/夏节/云产品专属人 —— 按分支 A-H 决策树落到谁;**分支 E 主单 assignee 仍改谜拟人类兜底,不改成镇元 agent 工号**——客户主单指派机器人客户看不懂,agent 只承接关联单) | 建关联单同时把主单 assignee 同步到承接方,便于承接方看到自己单子;**分支 A(专属维护名单)例外**——只改 assignee + 状态,**不建关联单**(见 Step 3 分支 A 段落) |
| **补建关联单** | **保持不动**(前次已对) | 前次的 assignee 已经对,只是漏建关联单;不能重复改 |
| **观察等待**(<30 天) | **保持不动** | 承接方正在处理,不打搅 |
| **进度跟进**(D) | **保持不动** | 承接方仍是主责,只是催进度 |
| **追料/补料提醒**(F) | **保持不动** | 承接方已给结论等客户回料,不改 assignee |
| **关单提示**(E) | **保持原承接方**(即最后处理人 —— PR 合入者、云产品拒接结论作者、根因定位人) | 承接方看到自己已完成任务归档;客户回帖时承接方能收到 notification 而非 jarvis 收到再转;若原承接方就是过载(见文档改造分支),关单时 assignee **本来就是过载**,不需要"转"——是自然结果 |

**反常识警示 · 关单提示不转过载**:
- **不要**在关单提示阶段把 assignee 从原承接方(如许章/齐澄)改到过载(484483) —— 这样承接方后续收不到客户回帖 notification,jarvis 反而背了不属于自己的追踪责任
- 例外:若承接方本身就是**过载**(纯文档改造分支的路径就是 jarvis 接手,assignee 从建关联单时就是过载),关单时 assignee 保持过载即可,是**结果**不是"转向"

**如果诉求是纯文档问题应该做什么**:
- 按分支 G 前的"文档改造分支"路由:**先建 tf_provider 关联单指派过载(484483)** + 源单 assignee 改过载 → jarvis 直接 claim 关联单干活(worktree/PR/合入) → 关联单 done + 主单关单提示
- **不要跳过"建关联单"直接开 provider worktree 提 PR**——虽然 jarvis 就是过载能干活,但没建关联单就丢了研发档案,主单也难留下"研发详情 vs 关键节点"的分层;2026-07-06 78504233/78523353/78554774 三单就走了这个隐式捷径(见反模式)

**辅助判定原则**:
- 「实质进展」= 承接方(不是提单人也不是 jarvis 自己)在评论里给出**技术信息、排期时间、决策结论**其中一样;单纯 "@某人跟进"、"辛苦看下" 不算实质进展
- 「30 天」按 UTC 计算,以最后一条实质进展评论的时间为基准
- 「路由对不对」的判定:承接方是否在**本团队分工表**内且与诉求领域匹配(专属维护名单/新山/过载/谜拟(分支 E 人类兜底 owner)/镇元 agent WORKER_1783326253279 (分支 E 关联单)/临钧/夏节 NPE,详见 [team-roster.md](./team-roster.md));不在本团队分工的(如"刘源"外包工号、"皓桦"云产品同学、客户 TAM)一律视为"前次路由错",走转单
- 「缺关联单」的判定:源单看 relation list,若没有指向承接方对应池的关联单,即算缺——承接方 → 池映射:过载/新山/夏节/G/H → tf_provider(528766),**分支 E 谜拟 → 2165097 池镇元 agent 关联单**(谜拟保留在主单 assignee,关联单 assignee=镇元 agent WORKER_1783326253279),临钧 → 528766(由 acube 自动建);单纯与其他客户单的 relate 不算"承接关联单"。**分支 A(专属维护名单)不适用此判定**——A 本身不建关联单,前次已将 assignee 改到名单人即为终态,不存在"缺"
- 一张工单可能**同时**满足"关单提示 + 追料"(承接方已给拒接结论 + 客户未确认),此时优先走关单提示(状态改「已拒绝」直接闭环),不做追料
- 一张工单可能**同时**满足"缺关联单 + 距上次进展 ≥30 天"(前次只改 assignee 没建单),此时先补建关联单(带上进度跟进语气的评论),不再另发 D
- 若判定犹豫,默认往"观察等待"或"进度跟进"倾斜(comment 免 bookend/不发,轻量;不改状态,不建重复关联单)

### 通用骨架 · 6 类模板都用「结论 → 查证资料 → 建议行动」3 段

所有模板(A/B/C/D/E/F)都遵循**结论 → 查证资料 → 建议行动**骨架,差异在特化字段和 @ 语法。

**查证资料清单(通用)**——涉及具体 tf 资源/API 时按需选贴,不必全列:

- 【provider 源码】(涉 provider 改动): `alicloud/resource_alicloud_xxx.go:LINE`,url `https://github.com/aliyun/terraform-provider-alicloud/blob/master/alicloud/<file>.go#L<n>`
- 【上游 API】(涉 API): `<Product> <Action>`,url `https://help.aliyun.com/document_detail/<xxx>.html` 或 `https://api.aliyun.com/product/<Product>`
- 【镇元/Cloudspec】(涉资源接入): resourceTypeCode = `<>`,CoverageScore = `<>`
- 【关联工单】: `<ID> <title>`,url `https://project.aone.alibaba-inc.com/v2/project/<pool>/req/<ID>`
- 【相关 PR】(涉 provider 已 PR/已发布): `#NNNN`,url `https://github.com/aliyun/terraform-provider-alicloud/pull/NNNN`(合入版本 `vX.Y.Z`)
- 【前次进展评论】`<日期> <评论人>: "<摘要>" (comment ID: <ID>)`

**建议行动段的双 @ 规范**(模板 E/F 强制,其他视场景):

- **提单人段**(客户/补料侧)——负责验收 / 补料 / 后续复现
- **承接方段**(最后处理人知会)——获知本单归档 or 补料状态,避免再来追

**下方 6 个模板** 列出各分支/场景的**特化字段**(结论用语、查证依据的场景特化项、行动请求要点);查证资料通用清单不再逐条重复,按需从上面选贴。

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
客户诉求「<一句话诉求重述>」<与镇元相关且镇元 NOT OK / 与镇元不相关(纯
datasource / 镇元 OK 但 provider 侧问题)>,应由 <指派人所在层> 侧承接。
已建关联单 <NEW_ID> 到 <terraform-alicloud (528766) / Terraform镇元对接 (2165097)>
项目,指派 @<花名>(<工号>) 跟进,源工单同步指派并改状态为「问题解决中」。
[分支 E 紧急双单场景改写为:已建两张关联单——<agent 单ID> 到 Terraform镇元对接
(2165097) 指派 镇元 agent (WORKER_1783326253279) 自动接单,负责镇元侧根因修复
(spec/映射/覆盖度,body 已按机读契约写清);<新山单ID> 到 terraform-alicloud
(528766) 指派 @新山(521957) 紧急兜底 provider 侧,双单并行;源工单指派谜拟
(479782,人类兜底 owner)并改状态「问题解决中」,@谜拟 + @新山。]

### 查证依据
1. **镇元相关性**:<纯 datasource(跳过镇元查证)/ 资源类镇元 OK(与镇元
   不相关)/ 资源类镇元 NOT OK(与镇元相关)>
   - 镇元:<get: 有/无 data | list released 是否命中 | CoverageScore=<x>>
   - acube V2 覆盖度(线上):CoverageScore=<x>
   - Property/Operation/PrimaryOperation=<>/<>/<>
2. **provider 代码类型**:<自动生成 / 手写>(依据 `alicloud/resource_...go:1-3`
   <是否有 generated automatically 注释>;纯 datasource 看 data_source_...go)
3. **紧急度**:priority=<>, 计划截止=<>, 剩余=<>天<,缺陷类型覆写=紧急>,
   故路由到 <人 / 分支 E 镇元 agent 单(+紧急时并新山单)>

### 关联单
- <NEW_ID>: <标题>,项目 <528766 terraform-alicloud / 2165097 Terraform镇元对接>
  (**分支 E 镇元 agent 单** → 2165097 Terraform镇元对接,assignee=WORKER_1783326253279;**其它人**(过载/新山/夏节/文档改造/G/H) → 528766 terraform-alicloud)
  (临钧场景:aoneId 由 acube V2 createBuildTaskV2 接口异步创建;taskId=<>)
  (分支 E 紧急场景列两张:agent 单 <ID> @2165097 镇元侧根因 / 新山单 <ID> @528766 紧急兜底)
- 双向关联已加(双单场景两张各自与源单双向关联)

@<花名>(<工号>) 烦请跟进上述查证结论,进度请在两侧工单同步回帖。
[双单场景:@谜拟(479782)(人类兜底 owner) @新山(521957) 烦请按上述分工并行跟进——
agent 单已进入镇元 agent 自动处理队列(spec/映射),新山单请紧急兜底 provider 侧;
进度请在各自关联单与本单同步回帖。]
```

### 模板 C:专属名单产品(分支 A)

```
### 结论
「alicloud_<product>_<resource> <诉求>」属 <产品中文名> 云产品域,由该云产品
provider 专人维护(不接镇元)。已指派 @<花名>(<工号>) 跟进,状态改为
「问题解决中」。

@<花名>(<工号>) 请协助评估 <诉求>,进度请回帖。
```

### 模板 D:进度跟进(承接方 ≥30 天无实质进展)

用于**已转出且距上次实质进展 ≥30 天**的工单(见判定表)。不叫"追料"(语气偏催促),用"进度跟进"(语气偏专业协同)。**必须贴出至少 2 条真实链接**(provider 源码 + API 文档 + 关联工单 + 前次进展评论 comment ID 任选),让承接人有具体上下文,不空催。

```
### 进度跟进 · <一句话诉求>

**查证资料**(参见通用骨架清单,本场景常用【provider 源码】+【前次进展评论 comment ID】+【相关 PR】)

**进度跟进请求**:
@<承接人花名>(<工号>) 距上次进展 <X> 天,请更新:
1. <具体问题 1,基于查证资料给出可执行提问,不空问>
2. <具体问题 2>
3. <可选:兼顾双方的过渡方案 / 排期建议>

若已闭环烦请回帖同步(附 PR/资源链接);若延期请给出新时间线;若长期无进展,建议 <升级到 XX 或转 XX>。
```

**发送方式**:进度跟进属于**只发评论、不改状态、不建关联单**,可直接 `bin/a1id -- project workitem comment create <id> -m "$(cat body-file)"` 免 bookend。避免走 `wrap.sh done` 的 heredoc,规避反引号/`$var:字母` 展开风险(见反模式)。发前 point-read 一次最新评论,避免与承接人刚回帖的进展撞车。**评论发完追加一条钉钉私信承接方**(失败不阻断):

```bash
bash bootstrap/notify-dingtalk.sh <承接人工号> \
  "Aone 进度跟进 · <一句诉求>" \
  "距上次进展 <X> 天,请更新:<关键问题一句话>
链接: https://project.aone.alibaba-inc.com/v2/project/<pool>/req/<单号>"
```

### 模板 E:关单提示(承接方已闭环 / 云产品明确拒接 / 已修复但主单未同步)

用于**承接方给出终结论后**主单状态仍未收敛的场景(见判定表)。走 bookend + 状态改「已发布待需求方验收」/「已拒绝」/「方案功能已存在」;**评论必须同时 @ 提单人和最后处理人**——提单人才知道诉求闭环并可后续验收,最后处理人被通知本单归档避免后续被再次追问。

```
### 结论(<PR/已发布/已拒接/根因结论>)
<一句话诉求 + 一句话闭环状态>

**查证资料**(参见通用骨架清单,本场景常用【provider 源码】+【相关 PR / 合入版本】+【前次进展评论 comment ID】)

**建议行动**:
- 客户侧:@<提单人花名>(<工号>) 请升级到 vX.Y.Z 后验收;若验收通过烦请回帖同步,本单归档。
- 承接方:@<最后处理人花名>(<工号>) 感谢跟进,本单据此关闭,后续若客户复现请另开新单直连你。

若已闭环无需回帖,状态已同步为「<已发布待需求方验收 / 已拒绝 / 方案功能已存在>」。
```

**发送方式**:
- **走 bookend**(claim → wrap done + status → release/finish),避免关单动作与 jarvis-idle 循环冲突。
- **状态选择**:PR 已合待客户验收 → 「已发布待需求方验收」;云产品/API 明确拒接 → 「已拒绝」;资源已支持只是客户版本旧 → 「方案功能已存在」;客户配置根因非 provider bug → 走模板 F 追料确认后关(而非本模板)。
- **双 @ 语法**:`@<提单人花名>(<提单人工号>)` 和 `@<最后处理人花名>(<最后处理人工号>)` 各带工号,提单人放 "客户侧" 段落,承接方放 "承接方" 段落——语义清晰,notification 两侧都收到。
- **最后处理人识别**:从 activity 或 comment list 找最后一条**给出实质进展**的评论(非 canned/@ 追问),其作者即最后处理人;若最后处理人 = 提单人自己(自派单),仅 @ 一次即可。
- **钉钉私信提单人 + 最后处理人**(bookend 完成后追加两条,失败不阻断):

  ```bash
  # 提单人一定是阿里前线/PD(不是外部客户),可以私信
  bash bootstrap/notify-dingtalk.sh <提单人工号> \
    "Aone 关单 · <一句结论>" \
    "本单已归档:<PR 已发布 vX.Y.Z / 云产品已拒接 / 方案已存在>
  状态:<已发布待需求方验收 / 已拒绝 / 方案功能已存在>
  链接: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>"

  # 若最后处理人 ≠ 提单人,再发一条
  bash bootstrap/notify-dingtalk.sh <最后处理人工号> \
    "Aone 关单 · <一句结论>" \
    "感谢跟进,本单已归档,后续同题会另开新单直连你。
  链接: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>"
  ```

### 模板 F:追料/补料提醒(承接方已给结论,客户/云产品未回)

用于承接方已给根因/结论(如"你的 lifecycle 配置引起 diff""API 层限制单值")后,**客户/云产品 ≥1 次追问仍未回**的场景。**提醒里要给出建议补齐时间与到期处理方式**,避免无限期挂着。**评论必须同时 @ 提单人和承接方**——提单人负责补料/确认,承接方获知补料进度(避免工单到期关闭后承接方以为球还在他手里再来追)。

```
### 补料提醒 · <一句话根因或阻塞点>

**查证资料**:
- 【前次结论评论】<日期> <承接人>:"<根因/结论摘要>" (comment ID: <ID>)
- 【provider 源码 / API 文档】(定位根因用): <链接>
- 【已请求补料内容】: <具体要什么,如"完整 main.tf + terraform.log + Code 字段的错误信息">

**待补齐材料**(<次数>次追问未回):
1. <材料 1,如 完整 tf 测试模板>
2. <材料 2,如 报错时 ECS 返回的 Code 字段(非仅 RequestId)>
3. <材料 3,可选>

**补料请求与时限**:
- 客户/补料侧:@<提单人花名>(<工号>) 麻烦在两周内(**建议 YYYY-MM-DD 前**)补齐上述材料,以便我们继续定位;
- 承接方知会:@<承接方最后处理人花名>(<工号>) 本单已请客户补料,若到期仍未补齐,将按「<无法复现/客户未回料/云产品无排期>」暂时关闭;后续同题客户会开新单直连你,本单**无需你再来跟进**。

到期处理:
- 若两周内补齐 → 本团队继续定位;
- 若两周内仍未补齐 → 本单先按「<无法复现 / 客户未回料 / 云产品无排期>」关闭,后续再遇报错欢迎**开新单**并直接附上材料(减少往复)。
```

**发送方式**:
- **免 bookend**,走 `bin/a1id -- project workitem comment create <id> -m "$(cat body-file)"`,只发评论;状态**不立即改**,保持等回料。
- **双 @ 语法**:`@<提单人花名>(<提单人工号>)` 和 `@<承接方花名>(<承接方工号>)` 各带工号;提单人放"客户/补料侧",承接方放"承接方知会"段落,语义清晰,notification 两侧都收到,承接方知道本轮球在客户手里、不用他再追。
- **承接方识别**:activity/comment list 里给出实质进展/结论的最后作者(非 canned/@ 追问);若承接方 = 提单人自己(自派单/前线随机指派回自己),仅 @ 一次即可。
- **到期后**:若客户/云产品补齐 → 走正常处理;若仍未回 → 走模板 E 关单(状态改「客户未响应」或对应闭环状态)。
- **默认建议补齐时间**:today+14 天;若情况紧急(如 P1 缺陷)可缩短到 7 天,但需在评论中说明加急理由。
- **不要**在同一单反复发补料提醒——若前一轮提醒到期仍未回已按流程关单,再遇同题另开新单跟踪。
- **钉钉私信提单人 + 承接方**(评论发完后追加,失败不阻断):

  ```bash
  bash bootstrap/notify-dingtalk.sh <提单人工号> \
    "Aone 待补料 · <一句诉求或阻塞点>" \
    "麻烦两周内(YYYY-MM-DD 前)补齐:<关键材料一句话>
  链接: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>"

  # 若承接方 ≠ 提单人,再发一条知会
  bash bootstrap/notify-dingtalk.sh <承接方工号> \
    "Aone 待补料 · <一句诉求>" \
    "本轮球在客户手里,已请其补料,到期未回将按'无法复现'关单,无需你再追。
  链接: https://project.aone.alibaba-inc.com/v2/project/1086837/req/<源工单ID>"
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

按主题分 5 组。每条给出「表现」+「正确写法/规则」;新踩坑时按组回查更快。

### A. 分诊误诊(判断错 / 漏读关键信息)

- ❌ 只按标题查证,不读 description 末段"限制/差异/仍需" —— 经典误诊(SSL Update 案就栽这)。参见 memory `read-description-last-paragraph`
- ❌ 直接把"类比产品支持 X"当"我们也应该做" —— 必须先分清是**API 原生**还是**Provider 侧适配**,前者 100% 是上游产品缺口
- ❌ 专属名单产品还去查镇元覆盖度 —— 不接镇元,直接指派
- ❌ 跳过 Step 1.5 共通 gate 直接发 canned —— 分类误建 / 重复单情形下会与承接单重复打搅客户
- ❌ Provider 全局改造(region 白名单/框架 utility/公共 endpoint/SDK bump)走"镇元 OK/NOT OK"判定 —— 镇元管资源 schema,不管 provider 基础;直接分支 G → 新山(521957),不必查 acube 覆盖度
- ❌ 只按 acube V2 `CoverageScore==1.0` 判"镇元 OK",忽略"当前 schema 属性是否覆盖客户诉求字段" —— 覆盖度分只反映已建 schema 的属性测试完备度,客户想要新字段而 schema 未建时,覆盖度分再高也是 NOT OK(缺口在镇元)
- ❌ 仅文档改造的工单因镇元 NOT OK 就把**主路由**转给镇元 agent/新山 —— TF 文档改造路由固定过载(484483) + jarvis 跟进,镇元查证仍要走(确保文档与 schema/API 一致),但结果不影响主指派对象;**但**(2026-07-09 补规):若查证发现镇元 metadata 也有同类问题,**必须额外建一张镇元 agent 侧单**(2165097),两张单各管各的
- ❌ 文档改造查证发现镇元 metadata 有问题却不建镇元 agent 侧单 —— 主路由走 TF 文档修复(过载),镇元 metadata 问题是独立线,不建单 = 镇元侧问题无人跟进、下次同类文档改动又撞同一个坑。正确:额外建镇元 agent 关联单(2165097,body 按硬契约),私信谜拟(479782)
- ❌ 缺陷类型工单沿用原单 `priority` 字段值 —— 缺陷(功能缺陷/线上问题/性能瓶颈)优先级一律覆写为"紧急",无视原单标记;覆写后影响两处紧急判定(分支 E 紧急→镇元 agent 单+新山双单;与镇元不相关分流紧急→新山)及关联单 `--priority` 传值
- ❌ 镇元 OK 但 provider 侧有问题的单仍派镇元 agent —— 2165097 池的镇元 agent 只接「与镇元相关且镇元 NOT OK」;镇元侧无问题即「与镇元不相关」,紧急→新山(521957),不紧急→过载(484483),生成器产出走临钧管道
- ❌ **分支 E 关联单 body 不含 `## 机读信息` + JSON 段** —— 镇元 agent 靠机读 JSON 驱动 spec/映射/覆盖度动作,没有机读段 agent 完全接不了,单沉底(2165097 池又不在 jarvis 视检范围,单会长期烂在里面);正确姿势:严格按 [templates.md 硬契约](./templates.md) 骨架写 `## 背景` / `## 需求` / `## 机读信息` 三段 + 7 字段 JSON,少一样就是漏契约
- ❌ **分支 E 关联单 assignee 写谜拟 479782** —— 谜拟已不解单(2026-07-08 切换),关联单硬指派镇元 agent (`WORKER_1783326253279`) 自动接单;写成 479782 = 落到人手不再走 agent 自动化,单会滞留
- ❌ **分支 E 主单 assignee 改成镇元 agent 工号** —— 客户主单在 tf_customer(1086837) 池,客户可见;指派机器人工号客户看不懂;正确姿势是主单 assignee 保留谜拟(479782,人类兜底 owner) + 关联单 assignee 才是镇元 agent
- ❌ **转单/建关联单不发钉钉私信** —— 通用规则(line 434-443)要求所有实质动作(转单/补建关联单/进度跟进/关单提示)后私信承接方;Aone @ 是主通道但易被漏看,钉钉是补充通知。**每个分支落地脚本尾部都带 `notify-dingtalk.sh` 步骤**,漏跑 = 承接方不知道有新单,工单响应延迟
- ❌ **钉钉私信直接发镇元 agent 工号** —— `WORKER_` 前缀是 agent 身份,无 IM 通道,notify-dingtalk.sh 传该 id 会 400/静默;正确姿势是私信谜拟(479782)兜底,紧急场景并私信新山(521957)
- ❌ 纯 datasource 工单去查镇元覆盖度 —— datasource 是 provider 侧对查询 API 的只读封装,镇元只管资源 schema;纯 datasource(不涉资源 schema/生命周期)直接判「与镇元不相关」分流,查镇元浪费一轮且可能误路由。resource+datasource 混合诉求不算"纯",仍按资源主线查镇元
- ❌ **未复现客户报错就直接路由** —— 分诊模糊时直接甩 NPE 兜底/上游是**懒惰误诊**。**正确顺序**:先读 description 全文(尤其客户 tf 代码 + 完整报错栈的堆栈行号),按报错文件路径 grep provider 源码定位根因,再决定路由。**规则**:客户 tf 代码 + 完整报错栈齐备时,先做静态复现(source 定位根因)再路由;仅"canned 类咨询/无报错/跨多产品"才走 NPE
- ❌ **「控制台 vs Terraform data source 结果不一致」类工单默认走分支 F 上游缺口** —— 或直接甩 NPE 兜底。大多数情况实际是 **provider 侧客户端多层过滤 + 客户配置差异**导致的表面不一致,**不是**上游 API 缺口。**判定路径 3 步**:① provider 调什么 API(grep `alicloud/data_source_alicloud_xxx.go` 定位实际 SDK Action,可能多个:主 API + 二次过滤 API + image 交集 API);② 控制台调什么 API(前端如 `ecs-buy.aliyun.com/api/...describeXxx.json` 通常是**公开 API 的内部聚合封装**,公开等价物是同族 OpenAPI,如 ecs-buy `describeAvailableInstanceTypes.json` ↔ 公开 `DescribeAvailableResource`);③ 过滤字段差异比对(常见 provider 侧额外客户端过滤:`image_id` 参数触发 `DescribeImageSupportInstanceTypes` 交集 / `spot_strategy` 传入即过滤 / `IoOptimized` 硬编码 `"optimized"` / `SoldOut` 强过滤,控制台前端通常不做这些)。**结论**:同族 API + provider 客户端过滤差异 = 纯 datasource 问题 → **与镇元不相关分流**(跳过镇元查证;不紧急 → **D-过载**,该类透明度改进多为非紧急:doc NOTE 补参数副作用 / 空集错误提示 / IoOptimized 类硬编码参数暴露;紧急 → D-新山);真正上游 API 缺口(公开 API 缺参数无法替代内部聚合)才走分支 F;绝不走 H NPE(场景明确、单一 data source、单一云产品域,不属"跨多产品无单一负责人/分类模糊")

### B. 路由动作(建关联单 / 优先级 DDL / upstream 前扫 / 分支 E 双单)

- ❌ 上游 API 缺口场景建了关联单 —— 应该只 @提单人 + 待上游排期,别拖镇元 agent/过载/临钧下水
- ❌ **专属维护名单产品(ACK/SLS/MNS/OSS/ESS/OTS/EMR/RDS/PolarDB/MSE/ClickHouse)建 tf_provider(528766) 关联单** —— 违反分支 A 定义:该批产品 provider 由专人独立维护、不接 tf_provider 共享池,追踪走源客户单本身;建关联单等于污染专人队列 + 双档案。**正确**:源单 `--assignee <名单工号>` + `--status 问题解决中` + wrap done + release,评论走模板 C(不提关联单)。若该产品已有专人推的 PR/修复,走 Step 3 前置 Gate 短路,同样不建单
- ❌ 强行 `--assignee` 指派专属名单以外的产品到过载/新山/镇元 agent —— 违反分工表,让本团队背不该背的锅
- ❌ 生成器产出(临钧)场景还手动 `a1 workitem create` 建关联单 —— acube V2 createBuildTaskV2 已自动建单+指派临钧,重复建会双单,污染临钧研发队列
- ❌ acube 60s 未返回 aoneId 就"降级"回手动 `a1 workitem create` —— 可能 acube 已建成功只是查询未及时,回退会双建;正确做法是升级 escalation 由人排查
- ❌ 转单不复制原单优先级 / 不设短于原单 DDL 的截止日期 —— 关联单接手方无优先级参考,DDL 与原单齐会让下一棒无余量;规则:`--priority` 复制原单,`--cfs 计划截止日期` = 原单 DDL - 2 天(至少 today+1);原单无 DDL 时默认 today+3
- ❌ Provider 源码查证跳过 Step 2 前置的 upstream PR 前扫,只 grep 本地 workspace 磁盘 —— workspace 的本地 branch 可能滞后 upstream 数十小时,或 sync-provider.sh 未 hardened 时只 fetch 不 reset;必先跑 `gh pr list --search` + `gh api contents?ref=master`,同题 recently-merged PR 直接引用避免重复建单
- ❌ Step 3 执行路由前不扫评论区/状态 —— 查证+决策期间有时间窗口,原指派人(常被前线随机指派)可能已回帖修复/贴 PR/接手;不扫就转单会重复建关联单、让接手方收多余通知
- ❌ 分支 E 紧急单只建镇元 agent 一张关联单 —— 紧急(优先级=紧急/距 DDL<14 天/缺陷覆写)时必须镇元 agent+新山**两张**关联单并行(agent 修镇元侧根因、新山紧急兜底 provider 侧),漏建新山单=紧急件失去快速通道;评论 @谜拟+@新山(agent 不接 IM 不用 @)并注明分工
- ❌ 纯文档诉求跳过"文档改造分支"转单直接开 provider worktree 提 PR —— 隐式捷径:jarvis 觉得"改一行 markdown 何必建关联单",直接就干了 + 主单发关单提示。**问题**:(1) 丢了 tf_provider(528766) 关联单的研发档案,(2) 主单研发详情与关键节点混在一起没分层,(3) 与其他分支的动作模式不一致——文档改造分支设计的正确路径是"建关联单指派过载 → 主单 assignee 改过载 → jarvis 直接 claim 关联单干活 → 关联单和主单各自 done + release"(见"文档改造分支 · 仅 website/docs 变更"段落)

### C. assignee / status 转向

- ❌ 状态改成"问题处理中" —— tf_customer(1086837) 池合法枚举没这个值,合法名是"问题解决中";写错 a1 会 `unsupported target status` 阻断。合法枚举:需求待补充/待处理/评估中/待上游排期/问题讨论/长期跟进/待排期/已排期/问题解决中/已发布待需求方验收/验收中/验收通过/验收不通过/客户未响应/方案功能已存在/需求撤回/已拒绝
- ❌ 状态用 `--status 已完成` / `方案功能已存在` 兜底 —— 前者不在合法值,后者语义错(客户真诉求还没解决)
- ❌ 关单提示阶段把源单 assignee 从原承接方改到过载 —— 承接方(如原 PR 作者、云产品拒接结论作者)后续收不到客户回帖 notification,jarvis 反而背了不属于自己的追踪责任。**正确**:关单提示 assignee 保持原承接方,让承接方看到"自己已归档";只有原承接方本来就是过载(文档改造分支 / 分支 D-过载)时,关单时 assignee 保持过载才是**结果**(不是"转向")

### D. a1 CLI 陷阱

- ❌ 花名 @ 不带工号(`@谜拟`) —— a1 有时能补,有时不能,显式 `@谜拟(479782)` 保险
- ❌ 建完关联单调两次 `a1 relation add`(A→B + B→A) —— aone 自动双向,第二次会 400 "关联失败该条记录已存在";单次 `add <源> relate:<新>` 即建 A↔B
- ❌ 建关联单用 `--description` —— a1 CLI 不吃(报 `unknown flag: --description`),正文用 `--body` 或 `--body-file`
- ❌ **528766 池必填 cfs 清单**(按 category 分):
  - **Task / Req 类**:`计划开始日期` / `计划截止日期` / `实际工时` 三件套(漏 → 400 `【计划开始日期】不能为空...` 阻断)
  - **Bug 类**(功能缺陷/线上问题/性能瓶颈):在 Task/Req 三件套之外**追加** `Terraform需求类型`(默认传 `--cfs "Terraform需求类型=运行时问题，TF问题"`;查全部合法枚举 `a1 project workitem field options "Terraform需求类型" --project 528766 --type 36` [36=功能缺陷])
  - **限定**:此清单仅 528766 池必填;**分支 E 镇元 agent 单落 Terraform镇元对接(2165097) 无此约束**,可省 cfs 整块;实测遇到 2165097 池 cfs 校验(如日后加字段)按 400 报错补
- ❌ 用 `head -1 | cut -f1` 或 `awk -F'\t' '{print $1}'` 解析 `a1 workitem create --quiet` 输出抓 NEW_ID —— `--quiet` 输出是**空格分隔**不是 tab(实测 `<id><空格><title><空格><status><空格><assignee>`),tab 分隔符抓到的是整行;NEW_ID 会带脏字符,后续 `relation add` / heredoc 里 `$NEW_ID:xxx` 全部污染。**正确写法**:`... --quiet | head -1 | awk '{print $1}'`(不带 -F,按空白拆);抓完打印 `echo "NEW_ID=[$NEW_ID]"` 用方括号包裹肉眼验证边界
- ❌ `wrap.sh done <id> --summary-stdin <status> <<EOF ... EOF`(不带引号的 heredoc)里正文含反引号 code 块或 `$var:字母` —— shell 会展开 heredoc:反引号 `` `xxx` `` 被当命令替换执行(报 `command not found` 且被替换成空),`$NEW_ID:alicloud_xxx` 被 wrap.sh 前缀 pwd 拼成 `/path/to/jarvis/<id>licloud_xxx` 怪路径。**正确写法**二选一:(a) 用 `<<'EOF'`(带引号)阻止 shell 展开,但 `$NEW_ID` 也不展开,先 `envsubst` 或 sed 预替换;(b) 保留 unquoted heredoc 但正文里去掉反引号(用引号 `"xxx"` 代替 code)、`$NEW_ID` 后面接空格不接字母冒号(如写作 `关联单 ID = ${NEW_ID}`)。**批量场景推荐**:先把评论正文 sed 替换 NEW_ID 后写到 `/tmp/wrap-<id>.txt`,再走 `wrap.sh done <id> --summary-file /tmp/wrap-<id>.txt <status>`,彻底跳过 heredoc 展开风险(见 `./batch-bookend-template.md` 参考模板)
- ❌ 用 `a1 project workitem tag add <id> <tag>` 加标签 —— 该子命令**不存在**(`a1 project workitem` 只有 activity/attachment/comment/create/delete/field/get/list/relation/type/update)。标签统一走 `update --tag "a,b,c"` **覆盖式**写入;`claim.sh` 已实现 point-read 现有 tag + 合并再写,但外挂标签(如 `jarvis-npe`)必须在 `claim.sh release` 之后单独补一次:`bin/a1id -- project workitem update <id> --tag "jarvis-idle,jarvis-npe"`(release 后 tag=jarvis-idle,覆盖式写完整列表保留)
- ❌ `a1 update --tag "<name>"` 传 tag **name** 保留 existing tag —— 当 existing 里含跨项目/白名单外 tag(含括号 / 中文空格 / 不在 `field options tag --project X --type Y` 白名单)时,name-based `--tag` 校验会报 `not found` 或**静默截断**只写入白名单内的部分,业务 tag 丢失。**正确写法**:point-read `workitem get -f json` 拿 `.fields[].tag.value`(**数字 ID 列表**),`--tag` 参数用 ID 传或 ID+name 混合传(a1 CLI `--tag` 明说支持 name 或 ID);`claim.sh _get_tag_pairs` + `_update_tags_merged` 已实现 ID-aware 保留,外部调用侧手动 update 时也要遵循同规则

### E. 批量 / 流程

- ❌ 批量 bookend 里 `bash claim.sh claim ... ; wrap.sh ... ; claim.sh release` **不检查** claim exit code —— `claim.sh` lost race 时**退出码 1** 且 tag 已回滚(源单未被认领),但脚本继续跑 `wrap.sh done` 会把评论发出去,后续再回退用 `comment create` 补发,**同一条评论发出两次**污染工单历史。**正确写法**:`if bash bootstrap/claim.sh claim $id $pool; then bash bootstrap/wrap.sh done ... && bash bootstrap/claim.sh release $id $pool; else echo "$id lost race, skip"; fi`(完整可复制骨架见 `./batch-bookend-template.md`)
- ❌ 批量 bookend / 长循环脚本单次超过 Bash 工具 2min timeout —— 13 单 × 3 步 × 每步 ~5-10s ≈ 130-200s 会被截断,中断时最后一单常处于 `jarvis-claimed` 状态没走 `release`,遗留僵尸 claim 阻塞下一轮 lost race。**规则**:批处理**单 Bash 调用限 4-5 单**(<60s),分批;`preflight.sh` 已在开局主动跑 `reconcile.sh stale` 清理僵尸,但会话中的实时截断仍需手动 `claim.sh release <id> <pool>` 补救
- ❌ 「追料」canned 语气偏催 —— 长期未响应的工单发"追料"评论对承接人是**噪声**(没上下文只喊时间);进度跟进应走**模板 D**,贴 provider 源码 / API 文档 / 关联工单 / 前次评论 comment ID 等真实链接,让承接人有具体上下文;**规则**:批量进度跟进走 `comment create -m "$(cat body-file)"` **免 bookend**(纯发评论无状态变更),避免走 `wrap.sh done` 的 heredoc 展开风险和 `claim.sh` lost race 阻断
