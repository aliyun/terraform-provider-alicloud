# Terraformer Resource Development Skill v0.1 Forward Evaluation

Date: 2026-07-16

Method: three fresh-context, read-only evaluators first loaded the canonical `terraformer-resource-dev` Skill and technical reference, then answered one bounded scenario. No evaluator modified Terraformer or accessed cloud resources. Repeat the evaluation by loading those two files and presenting the prompts below.

## Scenario A — PASS

Prompt: a parent-free global List API returns every segment of a two-part Provider Import ID.

Observed decision: select pattern B, confirm segment order and delimiter from Provider `d.SetId`, `ParseResourceId`, Import docs, or Import tests, paginate the direct List, and do not add parent traversal merely because the ID is multipart.

## Scenario B — PASS

Prompt: the child List API requires `workspace_id`; the Provider Import ID is `workspace_id:member_id`; the Data Source also requires `workspace_id`.

Observed decision: select pattern C because the API requires parent scope, enumerate parents, create a fresh child request for each parent, reset child pagination inside the parent loop, and construct the leaf ID from Provider evidence. The Data Source requirement is supporting List evidence, not a global Terraformer input rule.

## Scenario C — PASS

Prompt: Provider schema and API fields imply a relationship, but the unified relationship artifact has no declaration for the new resource.

Observed decision: leave the relationship consumer unchanged, record the producer gap, and do not infer or produce a relationship. Continue core discovery and Import ID support unless relationship delivery is an explicit acceptance requirement.

## 中文化后的复验

将 canonical Skill 和技术 reference 改成中文后，又用三个互不共享上下文的只读执行者重复了上述场景。执行者只获得中文 Skill 路径和用户场景，没有获得预期答案：

- **Scenario A — PASS：** 选择模式 B；没有因为两段式 Import ID 引入父资源遍历；正确区分 token 与页码分页。
- **Scenario B — PASS：** 选择模式 C；明确父 ID 来自父资源 List，而不是 Data Source Required 参数；child 分页状态位于父资源循环内部。
- **Scenario C — PASS：** 统一关系产物无声明时保持 consumer 不变；不从 Provider schema 或 API 字段推导关系；核心资源发现和 Import ID 接入可以继续。

三个执行者均使用中文输出，并在无法解析已登记的本地 checkout 时明确报告“仅静态分析”，没有伪造源码查证、测试或真实云验证结果。
