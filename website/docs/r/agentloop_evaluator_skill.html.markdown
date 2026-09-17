---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_evaluator_skill"
description: |-
  Provides a Alicloud AgentLoop Evaluator Skill resource.
---

# alicloud_agentloop_evaluator_skill

Provides a AgentLoop Evaluator Skill resource.

Evaluator Skill is the skill file package bound to an Evaluator, which defines the evaluation logic of the evaluator.

For information about AgentLoop Evaluator Skill and how to use it, see [What is Evaluator Skill](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateEvaluatorSkill).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_evaluator" "default" {
  agent_space = alicloud_agentloop_agent_space.default.agent_space
  name        = "terraform_example_evaluator"
  type        = "AGENT"
  metric_name = "example_metric"
  version     = "v1"
}

resource "alicloud_agentloop_evaluator_skill" "default" {
  agent_space  = alicloud_agentloop_evaluator.default.agent_space
  name         = alicloud_agentloop_evaluator.default.name
  skill_name   = "terraform_example"
  display_name = var.name
  description  = "terraform-example"
  enable       = true
  files {
    name = "SKILL.md"
    # The backend requires the SKILL.md frontmatter 'name' to equal
    # skill_name, so keep them in sync.
    content = <<-EOT
---
name: terraform_example
description: Example evaluator skill
---
# Test Skill
EOT
    remark  = "initial version"
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Evaluator belongs.
* `description` - (Optional) The description of the Evaluator Skill. It can be updated in place.
* `display_name` - (Optional) The display name of the Evaluator Skill. It can be updated in place.
* `enable` - (Optional) Whether to enable the Evaluator Skill. It can be updated in place.
* `files` - (Required, Set) The skill files of the Evaluator Skill. Changing `files` updates the skill in place and replaces the whole file list. **NOTE:** Updates send the complete current metadata and file list because the API uses full-replacement semantics. See [`files`](#files) below.
* `name` - (Required, ForceNew) The name of the Evaluator to which the skill belongs.
* `skill_name` - (Required, ForceNew) The name of the Evaluator Skill.

### `files`

The files supports the following:
* `content` - (Required) The content of the skill file.
* `name` - (Required) The name of the skill file.
* `remark` - (Optional) The remark of the skill file. When the Get API returns no remark (null or an empty string) for a file, the configured value is kept in the state; therefore clearing a remark outside Terraform cannot be detected as drift, and the remark cannot be recovered by `terraform import`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Evaluator Skill. It formats as `<agent_space>:<name>:<skill_name>`.
* `created_at` - The creation time of the Evaluator Skill, as a Unix timestamp in seconds.
* `updated_at` - The last update time of the Evaluator Skill, as a Unix timestamp in seconds.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Evaluator Skill.
* `delete` - (Defaults to 5 mins) Used when delete the Evaluator Skill.
* `update` - (Defaults to 5 mins) Used when update the Evaluator Skill.

## Import

AgentLoop Evaluator Skill can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_evaluator_skill.example <agent_space>:<name>:<skill_name>
```
