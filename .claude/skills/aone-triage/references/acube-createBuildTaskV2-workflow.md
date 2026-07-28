# 临钧路由（生成器产出）· Acube V2 单写者工作流

> 本文件是 `tf-customer-request-routing.md` 分支 D-临钧的执行契约。Acube
> `createBuildTaskV2` 会创建真实的 528766 研发单并触发生成器，不是只读探测。

## 单写者与身份边界

- **executor 托管**的 headless run：PD、开发阶段 RD、QA 都只返回结构化结果，
  `requested_external_actions` 由 executor 执行；terraform-rd finalizer 只校验回执并把本轮
  唯一最终聚合写入 `AONE_RESULT`，不得自行调用 `wrap.sh` 或重复执行动作。
- 独立运行的 **terraform-rd finalizer** 才能直接执行下文的 Aone 善后与本轮唯一回复。先运行
  `bin/a1id ready terraform-rd`；未登录即返回 `missing_public_identity`，不得切换其他身份。
- **禁止回退 jarvis**，也不得用个人身份。PD/QA 永远不写 Aone、钉钉或 Acube。
- 下文出现的 `bin/a1id as terraform-rd -- ...` 仅是独立 finalizer 示例；executor 托管时
  绝不能运行这些命令。

## Existing-related 状态机（POST 前硬门）

先 point-read 源单最新 relation、评论/activity 中的 Acube 回执、assignee、workitemType 与
status，再按以下顺序决策：

1. 已有正确 528766 relation：复用关联单，禁止 POST。若源单
   `assignee=临钧（429768）` 和按类型映射的 status 有漂移，只幂等修差异字段；都一致则观察。
2. 评论/activity 已有 `taskId/aoneId`，但 relation 尚未回填：按 taskId
   `queryAoneByTaskId`，**只查询/复用**；拿到 aoneId 后再次 point-read，relation 只写一次。
3. 只有 taskId：只查询该 taskId，不能因暂时没有 aoneId 再 POST。
4. 没有正确 relation、taskId 或 aoneId，且三层证据确认是 D-临钧生成器路径：才允许调用一次
   `createBuildTaskV2`。
5. 无法判定已有回执是否属于当前诉求：返回 `blocked` 请求人工核对，禁止猜测后创建。

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
| `workId` / `workName` | executor 或 terraform-rd 的可审计调用身份，绝不能伪装成 jarvis |

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

独立 finalizer 的写操作形态如下；这段**不得在 executor 托管 run 中执行**：

```bash
# relation 的存在性必须已由 point-read 确认
bin/a1id as terraform-rd -- project workitem relation add <源单ID> relate:<aoneId>
bin/a1id as terraform-rd -- project workitem update <源单ID> --assignee 429768
bin/a1id as terraform-rd -- project workitem update <源单ID> --status "$progress_status"
```

实际执行时只选择尚未满足的行，不能把三行当成无条件脚本。`progress_status` 必须来自当前
源单类型的配置映射。

## 回执与最终聚合边界

动作执行者返回以下结构化回执给 terraform-rd finalizer：

- `taskId`、`aoneId`、`aoneUrl`，以及本轮是 created 还是 reused；
- relation 是 existed 还是 added；
- assignee/status 的 expected、before、after 与是否 changed；
- Acube 查询次数、最终状态和任何 blocker。

executor 托管时，finalizer 校验回执后只生成一份 `AONE_RESULT.reply_body`；独立 finalizer
把相同信息纳入唯一一次最终聚合。两种模式都禁止中途评论、`wrap.sh sync`、重复 relation
或第二次最终回复。
