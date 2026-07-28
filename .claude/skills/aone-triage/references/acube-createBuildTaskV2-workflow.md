# 临钧路由（普通 D 或 E 的 pre 交接）· Acube V2 单写者工作流

> 本文件是 `tf-customer-request-routing.md` 分支 D-临钧的执行契约，也接收
> **分支 E 的 pre Meta 已收敛**且 QA `verification_mode: cloudspec_pre` pass 的交接。Acube
> `createBuildTaskV2` 会创建真实的 528766 研发单并触发生成器，不是只读探测。

## 单写者与身份边界

- PD、开发阶段 RD、QA 都只返回结构化结果；**terraform-rd finalizer 是 downstream single-writer**，
  负责 Acube POST/query、下游单复用以及 relation/路由字段的幂等善后。
- **executor 托管**的 headless run 中，executor 只负责原主单 bookend（claim、唯一回复、
  outcome 状态、release/finish），不解析或重放 `requested_external_actions`。terraform-rd
  finalizer 先运行 `bin/a1id ready terraform-rd`，执行 downstream 动作并把回执写进
  `AONE_RESULT.reply_body`；不得调用 `wrap.sh`，executor 也不得重复 POST/relation。
- 独立运行的 **terraform-rd finalizer** 同样是唯一动作执行者，并额外负责本轮唯一回复。
  先运行 `bin/a1id ready terraform-rd`；未登录即返回 `missing_public_identity`，不得切换其他身份。
- **禁止回退 jarvis**，也不得用个人身份。PD/QA 永远不写 Aone、钉钉或 Acube。
- 下文 `bin/a1id as terraform-rd -- ...` 同时适用于 executor 托管 run 内的 RD finalizer
  与独立 finalizer；两者都由 RD finalizer 执行。executor 进程本身只做 bookend，不运行或
  重放 downstream 命令。

## Existing-related 状态机（POST 前硬门）

先 point-read 源单最新 relation、评论/activity 中的 Acube 回执、assignee、workitemType 与
status，再按以下顺序决策：

1. 已有正确 528766 relation：复用关联单，禁止 POST。若源单
   `assignee=临钧（429768）` 和按类型映射的 status 有漂移，只幂等修差异字段；都一致则观察。
2. 评论/activity 已有 `taskId/aoneId`，但 relation 尚未回填：按 taskId
   `queryAoneByTaskId`，**只查询/复用**；拿到 aoneId 后再次 point-read，relation 只写一次。
3. 只有 taskId：只查询该 taskId，不能因暂时没有 aoneId 再 POST。
4. 没有正确 relation、taskId 或 aoneId，且满足以下入口之一，才允许调用一次
   `createBuildTaskV2`：
   - 三层证据确认是普通 D-临钧生成器路径；
   - 分支 E 的 build/check/publish pre 与 pre Meta 已收敛，QA
     `verification_mode: cloudspec_pre` 已 pass 并返回 `pre_handoff`。
5. 无法判定已有回执是否属于当前诉求：返回 `blocked` 请求人工核对，禁止猜测后创建。

分支 E 在 pre 未收敛或 QA 未 pass 时禁止 POST。E 交接门不得泛化到 A/F/G/H/I、纯
datasource 或纯手写 Provider-only bug；这些路径仍按各自路由处理。

错误的历史 relation 不迁移、不关闭，也不能替代正确的 528766 relation。任何重试都从
point-read 开始；禁止重复创建、重复 relation、重复改派和重复阶段回复。

## Acube 接口

- `POST /api/v1/terraform_vendor_build/createBuildTaskV2`
- `GET /api/v1/terraform_vendor_build/queryAoneByTaskId?taskId={taskId}`

POST body 使用 `TerraformVendorBuildTaskDTO`：

| 字段 | 值 |
|---|---|
| `namespace` | 产品 namespace |
| `resourceTypeCode` | PascalCase 资源名 |
| `resourceTypeVersion` | 首版生成场景为 `0.0.0` |
| `osType` | `Linux` |
| `flowType` | `ACubeRelease` |
| `workId` / `workName` | terraform-rd finalizer 的可审计调用身份；禁止填 executor、jarvis 或个人身份 |

POST 同步返回 taskId；随后只用 `queryAoneByTaskId` 查询最多 60 秒，等待异步返回
`aoneId/aoneUrl`。60 秒没有 aoneId 时返回 `blocked`，回执必须保留 taskId 供下次继续查询；
禁止重新 POST，也禁止手工创建 528766 单。

正式域名为 `acube.aliyun-inc.com`，预发为 `pre-acube.aliyun-inc.com`。预发 POST 同样会写
真实 Aone；只做连通性检查时只能调用只读 query 接口。

## aoneId 返回后的幂等善后

先再次 point-read relation 和源单映射字段，然后：

1. relation 已存在则保持；不存在时执行一次
   `relation add <源单ID> relate:<aoneId>`。Aone 自动双向关联，禁止反向再写。
2. 源单 assignee 期望值固定为临钧（429768）；仅在不一致时更新。
3. 源单 status 必须按源单 workitemType 从 `config/pools.json` 的
   `.pools.tf_customer.progress_status[workitemType]` 精确解析。缺 mapping、值为空或不是当前
   合法枚举时返回 `blocked`；不得选择第一个状态或硬编码统一中文值。
4. assignee/status 已一致时不写；字段漂移时只写差异字段，不产生阶段评论。

RD finalizer 的写操作形态如下；executor 托管 run 与独立 run 都由 finalizer 以
terraform-rd 身份执行，区别只在于托管 run 的原主单 bookend 仍交 executor：

```bash
# relation 的存在性必须已由 point-read 确认
bin/a1id as terraform-rd -- project workitem relation add <源单ID> relate:<aoneId>
bin/a1id as terraform-rd -- project workitem update <源单ID> --assignee 429768
bin/a1id as terraform-rd -- project workitem update <源单ID> --status "$progress_status"
```

实际执行时只选择尚未满足的行，不能把三行当成无条件脚本。`progress_status` 必须来自当前
源单类型的配置映射。

## 回执与最终聚合边界

terraform-rd finalizer 保存以下结构化动作回执：

- `taskId`、`aoneId`、`aoneUrl`，以及本轮是 created 还是 reused；
- relation 是 existed 还是 added；
- assignee/status 的 expected、before、after 与是否 changed；
- Acube 查询次数、最终状态和任何 blocker。

executor 托管时，finalizer 执行动作、校验回执后只生成一份 `AONE_RESULT.reply_body`，
executor 仅做原主单 bookend；独立 finalizer 把相同信息纳入唯一一次最终聚合。两种模式都
禁止中途评论、`wrap.sh sync`、重复 POST/relation 或第二次最终回复。
