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
