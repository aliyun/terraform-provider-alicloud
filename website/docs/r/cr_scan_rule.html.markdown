---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_scan_rule"
description: |-
  Provides a Alicloud CR Scan Rule resource.
---

# alicloud_cr_scan_rule

Provides a CR Scan Rule resource.

Artifact Scan Rule.

For information about CR Scan Rule and how to use it, see [What is Scan Rule](https://next.api.alibabacloud.com/document/cr/2018-12-01/CreateScanRule).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cr_ee_instance" "default2Aqoce" {
  default_oss_bucket = "false"
  renewal_status     = "ManualRenewal"
  period             = "1"
  instance_name      = "pl-example-2"
  payment_type       = "Subscription"
  instance_type      = "Basic"
}


resource "alicloud_cr_scan_rule" "default" {
  repo_tag_filter_pattern = ".*"
  scan_scope              = "REPO"
  trigger_type            = "MANUAL"
  scan_type               = "VUL"
  rule_name               = "699"
  namespaces              = ["aa"]
  repo_names              = ["bb", "cc", "dd", "ee"]
  instance_id             = alicloud_cr_ee_instance.default2Aqoce.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) instance id
* `namespaces` - (Optional, List) Namespace scope of the scan
* `repo_names` - (Optional, List) Scope of warehouse
* `repo_tag_filter_pattern` - (Required) Product tag matching rules
* `rule_name` - (Optional) The name of the resource
* `scan_scope` - (Required) Scan Range
* `scan_type` - (Required, ForceNew) 规则类型
* `trigger_type` - (Required) Scan rule triggering method

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<scan_rule_id>`.
* `create_time` - Creation time
* `scan_rule_id` - The first ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Scan Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Scan Rule.
* `update` - (Defaults to 5 mins) Used when update the Scan Rule.

## Import

CR Scan Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_scan_rule.example <instance_id>:<scan_rule_id>
```