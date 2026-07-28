# 镇元查证与路由分支

> 本 reference 是「镇元/Cloudspec 覆盖度查证 + provider 代码路由」的单点维护。跨 skill 复用:aone-triage tf-customer 路由 + provider-resource-dev 资源开发都读它。

## Step 2 — 定位缺口(按决策树)

### 前置 · upstream PR / commit 前扫(涉及 provider 源码的分支必跑)

Provider 源码查证**不能只 grep 本地 workspace**——workspace 可能停在旧 HEAD、本地磁盘可能滞后 upstream 数十小时。`sync-provider.sh` 会 `fetch + reset --hard FETCH_HEAD` 强对齐 upstream master,但即便对齐仍要额外扫 upstream PR,防止漏掉刚 merged / 正在 review 的同题改动。

```bash
# 1) 强制同步 workspace 到 upstream master 最新
bash .Codex/skills/aone-triage/scripts/sync-provider.sh

# 2) 扫 upstream open + recent merged PR,命中同题关键字 → 直接命中已有改动
bash bootstrap/github-identity.sh gh pr list \
  --repo aliyun/terraform-provider-alicloud \
  --search "<关键字,如 ap-southeast-8 / <alicloud_xxx> / <属性名>>" \
  --state all --limit 10 \
  --json number,title,state,mergedAt,url

# 3) 若关心某个具体文件的 upstream master 实际内容,gh api 直取(不受本地 workspace 影响):
bash bootstrap/github-identity.sh gh api \
  "repos/aliyun/terraform-provider-alicloud/contents/alicloud/<file>?ref=master" \
  -q '.content' | base64 -d | grep -n <关键字>
```

**命中处理**:

- **MERGED PR**:记 PR 链接,查 `mergedAt` vs 最新 release `publishedAt`:
  - `mergedAt < 最新 release publishedAt` → 已发版,客户升级到该 release 或更新即可
  - `mergedAt > 最新 release publishedAt` → 待下一版本发版,客户先按临时方案(如 `skip_region_validation`)或等 release
  - 回复直接引用 PR,**不必再建关联单**、不必查镇元;分支 A/B/C/D/E/G 建单流程直接跳过
- **OPEN PR**:贴 PR 链接 + 状态到源单评论,让客户/提单人可跟进;是否建关联单看 PR 进度(已 review 中可先等)
- **无命中**:才走下方决策树的分支 A/B/C/D/E/F/G 查证 + 建单流程

**为什么放前置**:避免"以为 provider 没支持,建了关联单指派团队做,结果同题 PR 昨天刚 merged"。先例:工单 83718139(同题 PR 9909 merged 21h 后处理仍未命中)。

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

# ② 资源质量覆盖度 - acube V2 GET,内网直连,无需 AK/SK
# 参数口径:POP 三元组(popCode/popVersion/resourceName),与 ① 的 product/resourceCode 不同
#   popCode      POP 产品代码,常 UPPER,如 APIG;与 cloudspec product(Apig)可能不同
#   popVersion   POP 接口版本,如 2024-03-27;从 aliyun CLI meta 或 next.api 取
#   resourceName POP 资源名,PascalCase,多数场景同 ① 的 resourceCode
popCode=<PopCode>
popVersion=<PopVersion>
resourceName="$resourceCode"

curl -sSG 'https://acube.aliyun-inc.com/api/v1/terraform/generator/getResourceQualityDetailCoverageScoreV2' \
  --data-urlencode "popCode=${popCode}" \
  --data-urlencode "popVersion=${popVersion}" \
  --data-urlencode "resourceName=${resourceName}" \
  -H "accept: application/json" | python3 -c '
import json,sys
d=json.load(sys.stdin)
if d.get("code") != "SUCCESS":
    print("acube:", d.get("code"), d.get("message")); sys.exit(0)
cov = ((d.get("data") or {}).get("data") or {}).get("CoverageDetail") or {}
print("CoverageScore:", cov.get("CoverageScore"))
print("  Property/Operation/PrimaryOperation:",
      cov.get("PropertyCoverageScore"), "/",
      cov.get("OperationCoverageScore"), "/",
      cov.get("PrimaryOperationCoverageScore"))
'
```

#### 分支 I — CloudSpec 文档文本 metadata

文档问题先比较 OpenAPI 长期语义、CloudSpec `resource/property/operation description` 与
Provider docs。变更只涉及 description、字段解释、NOTE 与枚举文案，且
**不改变字段集合、类型、约束或 CRUD** 时，才是 text-only I：

- 从 `config/pools.json` 的 `upstream.cloudspec_docs_quality` **创建或复用** 2169561，
  指派念依（373108），入口保持 `submit_only`；
- Provider 公开 docs 同时错误时，补一条**独立 528766 紧急兜底腿**，指派过载（484483）；
- 两池执行**分池防重**，一个池已有 relation 不能抑制另一个池的缺失补建；
- PD/QA 不外写，terraform-rd finalizer 作为 downstream `single-writer` 执行动作；
  executor 只负责原主单 bookend。

CloudSpec 文档源正确、仅 Provider 本地生成/展示偏差时走 D；文档改动同时涉及字段集合、类型、
约束、枚举集合、CRUD/operation 或映射时不再是 I，走结构分支 E。

**CloudSpec 结构 OK 三条件**(全满足才算 OK,任一不满足即结构 NOT OK):

1. **API 在镇元有对应资源**:`get` 返回 data 且 `released` list 命中(资源已定义并发布)
2. **当前资源属性满足客户诉求**:比对 Step 1 抽取的真实诉求字段,镇元资源 schema 的 properties **全覆盖** —— 缺字段即视为 NOT OK(即便覆盖度分再高也不算 OK,因为覆盖度只测已建 schema 的属性,客户想要新字段时属性不覆盖=缺口在镇元)
3. **测试覆盖度 100%**:`CoverageDetail.CoverageScore == 1.0`(V2 仅以覆盖度综合分判定,无 PASS/FAIL 用例计数)
**判定**(按「与镇元相关性」口径):
- 三条件全满足 → **CloudSpec 结构 OK** = 结构侧无问题、缺口在 provider 侧，走分支 D
- 任一不满足(资源未定义 / 属性缺客户要的字段 / 覆盖度 < 1.0) →
  **结构 NOT OK**，走分支 E 的 CloudSpec 结构 metadata 原主单自闭环

**前置短路 · 纯 datasource 问题不查镇元**:诉求只涉 `data.alicloud_xxx`(查询/过滤/输出字段)、不涉资源 schema/生命周期的,直接判**与镇元不相关**、进分支 D 分流(非临钧子分支)——datasource 是 provider 侧对 List/Describe 查询 API 的只读封装,镇元只管资源 schema,本节的 get/覆盖度查证全部跳过。resource+datasource 混合诉求不算"纯",仍按资源主线走本节查证。

**环境说明**:V2 走 acube 正式(`acube.aliyun-inc.com`),默认线上;需要预发数据把域名换成 `pre-acube.aliyun-inc.com` 即可(路径 / 参数 / 返回结构一致)。服务端实现见邻仓 `a-cube-aliyun-com` 里 `acube-auto-generator/.../ProductMetaUtil.java#getResourceQualityDetailCoverageScoreV2ByCommon`。

**sanity check(必跑)**:拿同产品已发布的其他资源(如查 ① released list 里另一项)以相同接口复查,能拿到 `CoverageDetail` 说明内网/接口正常;否则先排查内网访问(`pre-acube.aliyun-inc.com` 需办公网/VPN;`/api/v1/**` 免鉴权,但走内网 DNS)。

### G / 紧急普通 D 的双 owner 契约

本契约覆盖 Provider 全局改造 G，以及紧急普通 D（纯 datasource；或 CloudSpec 结构 OK +
Provider 手写实现）：

- **源客户主单 assignee 保持新山（521957）**，但
  **528766 研发关联单 assignee 固定过载（484483）**，由
  **Jarvis/TerraformRD claim 并尝试修复**。新山是客户主单映射 owner，不是等待接手研发的人。
- 写前 point-read relation、同题单与 lease；**healthy existing claim 不抢占**。没有健康
  claim 时，**同题 528766 已指派新山时原地复用**。同一 terraform-rd Task 先进入
  route-finalizer phase 并 fail-closed claim，成功后才幂等改派过载、补 relation 并切 dev；
  claim 失败立即停止，禁止 relation/update/wrap。不存在同题单才 create 一次并 claim。
- relation/assignee/status 只证明路由物化，**无 PR/CI/QA 完成信号必须继续 RD**。只有已有
  PR 等待人工合并，或明确外部依赖/人工决策，才允许 observe/release 或 blocked。
- **build/test/CI 失败**与 QA fail 留在 RD ↔ QA 修复闭环，**不得转交新山**。能力缺失或重试
  耗尽进入 **blocked / SUSPENDED**，保持源单新山、研发单过载并 release，绝不 finish。
- TerraformRD control plane 幂等执行物化、必要改派、relation 与 claim；PD/QA 不外写。
  源单由 executor、实际 claim 的 528766 由最终 finalizer 分别执行一次聚合 bookend；
  PR 未合并只 release。
- **D-临钧/A/F/H/非紧急 D 边界不变**；非紧急 D 与 I 的 Provider docs 紧急兜底腿保持
  既有内部 claim/bookend，I→念依、E→CloudSpec pre→D-临钧保持既有路径。

### 分支 D:与镇元不相关(镇元 OK、Provider 本地问题 / 纯 datasource),分流

资源类先判 provider 代码类型:

```bash
provider_repo="$(bash bootstrap/workspace.sh dir terraform_provider)"
head -3 "$provider_repo/alicloud/resource_alicloud_<product>_<resource>.go" 2>/dev/null
```

- 首行有类似 `// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!` → **生成器产出** → 走 acube V2 接口触发临钧工作流(见 Step 3 · 分支 D-临钧;生成器产出修复=重跑生成器,管道不变)
- 无该注释(手写)**或纯 datasource 问题**(datasource 不走临钧管道)→ 按紧急度分流:
  - 紧急(优先级=紧急 OR 距 DDL<14 天 OR 缺陷类型覆写)→ 源单保持新山(521957)，创建或复用
    528766 并固定指派过载(484483)，由内部 TerraformRD claim 后按 Provider 开发流程修复
  - 不紧急 → 保持既有 D-过载边界，源单与 528766 均指派过载(484483)

文档问题只有在证据证明 **CloudSpec 文档源正确，Provider 本地文档生成/展示偏差** 时才允许
进入本分支。CloudSpec text-only 文档源错误必须回到 I；文档与结构同时变化时进入 E。

若资源文件不存在:说明 provider 代码尚未合入(镇元 OK 但 provider 未生成/合入)——按"生成器产出待跑"处理,同样走 Step 3 · 分支 D-临钧 的 acube V2 接口(接口内部会跑生成器 + PR),不必 jarvis 手动 comment 提醒生成。

### 分支 E:CloudSpec 结构 NOT OK → CloudSpec 结构 metadata 原主单自闭环

本分支只承接字段集合、类型、约束、CRUD/operation 与映射缺口；text-only 文档 metadata
固定走 I：

1. PD 返回 `requested_external_actions: []` 与 `next=terraform-rd/dev`。不得提出 create_related、
   relation、assign、另建文档兜底单或切个人身份。
2. RD 加载 `terraform-provider-release/references/cloudspec-pre-resource-loop.md`，先运行
   `bash bootstrap/cloudspec-core.sh doctor`。
3. 调用 `cloudspec-amp-workflow` 创建/切换 task 专属 feature 分支，并使用 **AMP 返回的 SSH URL**
   clone 对应 cloudspec-model；禁止在 master/main 编辑。
4. 在已有 `main.cspec` 的模型目录调用 `cloudspec-idl-guide`，再按改动类型使用
   `cloudspec-resource-edit`、必要时 `cloudspec-operation-edit`；build 失败才调用
   `cloudspec-build-fix`，并用 `cloudspec-norm-check-fix` 收敛本次增量。
5. `aliyun cspec build` 与资源级 `aliyun cspec check` 全绿后，提交/推送 feature 分支，
   执行 `amp publish pre --dry-run`，通过后 `amp publish pre`，轮询 pre Meta 收敛。
6. **pre 未收敛不得触发 Acube**。pre Meta 收敛后必须 **E → D-临钧**：
   已有正确 relation/taskId/aoneId 时只查询/复用，否则由 single-writer 调用一次
   `createBuildTaskV2`，让 Acube 自动创建/复用 528766 并指派临钧（429768）。
7. **不得由 E 直接执行 Provider PR/CI/ACC**，不得在 E 完成后直接 release/idle；只有
   D-临钧 handoff 回执确认后才由 finalizer 聚合。Acube 不确定时保留 taskId 并 blocked，
   不得手工补建。

权限、AMP 登录、SSH、模型仓或 pre 能力失败时返回 `missing_capability` / `blocked` 并记录原主单，
不得回退其它承接人或身份。`amp publish prod` / prod/online、master/main merge/push 与正式发布
始终是人工硬门，不得 finish 或宣称正式发布。

**只允许分支 E 进入此转换**；不得泛化到 A/F/G/H/I、纯 datasource 或纯手写 Provider-only bug。

### 分支 F:上游 API 缺口

- 不建关联单(镇元/我们团队无事可做)
- 提取 `creator.displayName` 作为提单人,评论正文 @提单人 + 请其协助转对应云产品 API 团队评估
- status 改 `待上游排期`(若 CLI 报 "unsupported target status" 却能实际写入,重试;或落 `待排期` 作二级降级)
