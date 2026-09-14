---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_digital_employee"
description: |-
  Provides a Alicloud Cms Digital Employee resource.
---

# alicloud_cms_digital_employee

Provides a Cms Digital Employee resource.

A digital employee is an AI agent managed by Cloud Monitor (Cms) that can be configured with a RAM role, knowledge bases, sandbox network policies, and tool security policies. It orchestrates LLM-based automation against your cloud resources.

For information about Cms Digital Employee and how to use it, see [What is Digital Employee](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateDigitalEmployee).

-> **NOTE:** Available since v1.282.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ram_role" "default" {
  name     = "tf-cms-de-role-${var.name}"
  document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["cms.aliyuncs.com"]
      }
    }
  ]
}
EOF
}

resource "alicloud_cms_digital_employee" "default" {
  digital_employee_name = var.name
  role_arn              = alicloud_ram_role.default.arn
  description           = "a digital employee managed by terraform"
  display_name          = "tf-digital-employee"
  default_rule          = "default-rule"
  attributes = {
    team = "platform"
  }
  knowledges {
    bailian {
      workspace_id = "ws-xxxxxxxx"
      index_id     = "idx-xxxxxxxx"
      region       = "cn-beijing"
      attributes   = "{\"source\":\"terraform\"}"
    }
  }
  sandbox_network_policy {
    allow_fqdns = ["api.example.com"]
    allow_cidrs = ["10.0.0.0/8"]
    enable_acl  = false
  }
  tool_policy {
    aliyun {
      enable           = true
      deny_policy      = ["ecs:Delete*"]
      auto_pass_policy = ["log:Get*", "log:List*"]
      statements {
        decision    = "user_ack"
        product     = "Sls"
        api_version = "2020-12-30"
        actions     = ["log:GetProject"]
      }
    }
  }
  tags {
    key   = "env"
    value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `digital_employee_name` - (Required, ForceNew) The name of the digital employee. It must be 1 to 64 characters in length.
* `role_arn` - (Required) The ARN of the RAM role that the digital employee assumes.
* `default_rule` - (Optional) The default rule of the digital employee.
* `description` - (Optional) The description of the digital employee.
* `display_name` - (Optional) The display name of the digital employee. It must be 0 to 128 characters in length.
* `resource_group_id` - (Optional) The ID of the resource group to which the digital employee belongs.
* `attributes` - (Optional) A map of custom attributes of the digital employee.
* `knowledges` - (Optional) The knowledge base configuration of the digital employee. See [`knowledges`](#knowledges) below.
* `sandbox_network_policy` - (Optional) The sandbox network ACL policy of the digital employee. See [`sandbox_network_policy`](#sandbox_network_policy) below.
* `tool_policy` - (Optional, Computed) The tool security policy of the digital employee. See [`tool_policy`](#tool_policy) below.
* `tags` - (Optional) A set of tags for the digital employee. See [`tags`](#tags) below.

### `knowledges`

The `knowledges` block supports:

* `bailian` - (Optional) A list of Bailian knowledge bases. See [`bailian`](#knowledges-bailian) below.

### `knowledges-bailian`

The `bailian` block supports:

* `workspace_id` - (Optional) The workspace ID of the Bailian knowledge base.
* `index_id` - (Optional) The index ID of the Bailian knowledge base.
* `region` - (Optional) The region of the Bailian knowledge base.
* `attributes` - (Optional) The attributes of the Bailian knowledge base, as a JSON string.

### `sandbox_network_policy`

The `sandbox_network_policy` block supports:

* `allow_fqdns` - (Optional) A list of allowed FQDNs, up to 50 entries.
* `allow_cidrs` - (Optional) A list of allowed CIDR blocks or IP addresses, up to 50 entries.
* `enable_acl` - (Optional) Whether to enable the sandbox network ACL.

### `tool_policy`

The `tool_policy` block supports:

* `aliyun` - (Optional, Computed) The Aliyun MCP tool security policy. See [`aliyun`](#tool_policy-aliyun) below.

### `tool_policy-aliyun`

The `aliyun` block supports:

* `enable` - (Optional) Whether to enable the Aliyun MCP tool policy. Defaults to enabled; set to `false` to disable.
* `deny_policy` - (Optional) A list of explicitly denied RAM actions, in the format `product:ApiName`, `product:Prefix*`, or `product:*`. Has the highest priority.
* `auto_pass_policy` - (Optional) A list of automatically allowed RAM actions. When empty, only read-only actions (`Get*`, `List*`, `Describe*`) are auto-passed; others require human confirmation.
* `statements` - (Optional, Deprecated) A list of legacy Aliyun OpenAPI tool policy statements. Use `deny_policy` and `auto_pass_policy` instead. See [`statements`](#tool_policy-aliyun-statements) below.

### `tool_policy-aliyun-statements`

The `statements` block supports:

* `decision` - (Optional) The execution policy after the API is matched, such as `user_ack`.
* `product` - (Optional) The Aliyun OpenAPI product name of the statement, such as `Sls`.
* `api_version` - (Optional, Deprecated) The Aliyun OpenAPI API version of the statement.
* `actions` - (Optional) A list of Aliyun OpenAPI actions, in the format `product:ApiName`, `product:Prefix*`, or `product:*`.

### `tags`

The `tags` block supports:

* `key` - (Optional) The key of the tag.
* `value` - (Optional) The value of the tag.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the digital employee. It is the value of `digital_employee_name`.
* `create_time` - The creation time of the digital employee.
* `update_time` - The update time of the digital employee.
* `employee_type` - The type of the digital employee.
* `region_id` - The region ID of the resource.
* `resource_type` - The resource type, fixed as `ALIYUN::CMS::DIGITALEMPLOYEE`.

## Timeouts

The `timeouts` block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the digital employee.
* `update` - (Defaults to 5 mins) Used when updating the digital employee.
* `delete` - (Defaults to 5 mins) Used when deleting the digital employee.

## Import

Cms Digital Employee can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_digital_employee.example <name>
```
