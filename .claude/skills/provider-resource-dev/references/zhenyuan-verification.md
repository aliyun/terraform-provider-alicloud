# 镇元查证与路由分支

> 本 reference 是「镇元/Cloudspec 覆盖度查证 + provider 代码路由」的单点维护(P3.a 从 aone-triage 抽出)。跨 skill 复用:aone-triage tf-customer 路由 + provider-resource-dev 资源开发都读它。

## Step 2 — 定位缺口(按决策树)

### 前置 · upstream PR / commit 前扫(涉及 provider 源码的分支必跑)

Provider 源码查证**不能只 grep 本地 workspace**——workspace 可能停在旧 HEAD、本地磁盘可能滞后 upstream 数十小时。`sync-provider.sh` 已 hardened 为 `fetch + reset --hard FETCH_HEAD` 强对齐 upstream master,但即便如此仍要额外扫 upstream PR,防止漏掉刚 merged / 正在 review 的同题改动。

```bash
# 1) 强制同步 workspace 到 upstream master 最新
bash .claude/skills/aone-triage/scripts/sync-provider.sh

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

**为什么放前置**:避免"以为 provider 没支持,建了关联单指派团队做,结果同题 PR 昨天刚 merged"。参见工单 83718139 教训——PR 9909 merged 21 小时后 jarvis 才处理仍未命中,是本 skill 的历史缺口。

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

**镇元 OK 三条件**(全满足才算 OK,任一不满足即 NOT OK):

1. **API 在镇元有对应资源**:`get` 返回 data 且 `released` list 命中(资源已定义并发布)
2. **当前资源属性满足客户诉求**:比对 Step 1 抽取的真实诉求字段,镇元资源 schema 的 properties **全覆盖** —— 缺字段即视为 NOT OK(即便覆盖度分再高也不算 OK,因为覆盖度只测已建 schema 的属性,客户想要新字段时属性不覆盖=缺口在镇元)
3. **测试覆盖度 100%**:`CoverageDetail.CoverageScore == 1.0`(V2 已不返回 PASS/FAIL 用例计数,仅以覆盖度综合分判定)

**判定**:
- 三条件全满足 → **镇元 OK**,走分支 D
- 任一不满足(资源未定义 / 属性缺客户要的字段 / 覆盖度 < 1.0) → **镇元 NOT OK**,走分支 E

**环境说明**:V2 已发布到 acube 正式(`acube.aliyun-inc.com`),默认走线上;需要预发数据把域名换成 `pre-acube.aliyun-inc.com` 即可(路径 / 参数 / 返回结构一致)。服务端实现见邻仓 `a-cube-aliyun-com` 里 `acube-auto-generator/.../ProductMetaUtil.java#getResourceQualityDetailCoverageScoreV2ByCommon`。

**sanity check(必跑)**:拿同产品已发布的其他资源(如查 ① released list 里另一项)以相同接口复查,能拿到 `CoverageDetail` 说明内网/接口正常;否则先排查内网访问(`pre-acube.aliyun-inc.com` 需办公网/VPN;`/api/v1/**` 免鉴权,但走内网 DNS)。

### 分支 D:镇元 OK,判定 provider 代码类型

```bash
provider_repo="$(bash bootstrap/workspace.sh dir terraform_provider)"
head -3 "$provider_repo/alicloud/resource_alicloud_<product>_<resource>.go" 2>/dev/null
```

- 首行有类似 `// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!` → **生成器产出** → 走 acube V2 接口触发临钧工作流(见 Step 3 · 分支 D-临钧)
- 无该注释(手写) → 指派 过载(484483)(见 Step 3 · 分支 A / D-过载 / E)

若文件不存在:说明 provider 代码尚未合入(镇元 OK 但 provider 未生成/合入)——按"生成器产出待跑"处理,同样走 Step 3 · 分支 D-临钧 的 acube V2 接口(接口内部会跑生成器 + PR),不必 jarvis 手动 comment 提醒生成。

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

