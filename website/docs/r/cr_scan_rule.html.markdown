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

-> **NOTE:** Available since v1.265.0.

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
  rule_name               = var.name
  namespaces              = ["aa"]
  repo_names              = ["bb", "cc", "dd", "ee"]
  instance_id             = alicloud_cr_ee_instance.default2Aqoce.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) Instance ID
* `namespaces` - (Optional, List) Set of namespaces:  
  - This parameter must not be empty when the scan scope is NAMESPACE.  
  - This parameter must contain exactly one namespace when the scan scope is REPO.
* `repo_names` - (Optional, List) Repository list:  
  - This parameter must be empty when the scan scope is NAMESPACE.  
  - This parameter must not be empty when the scan scope is REPO.
* `repo_tag_filter_pattern` - (Required) Regular expression for matching tags that trigger a scan
* `rule_name` - (Required) Event rule name  
* `scan_scope` - (Required) Scan scope
* `scan_type` - (Required, ForceNew) Scan type:  
  - `VUL`: Artifact vulnerability scan  
  - `SBOM`: Artifact content analysis  

The default value of this parameter is `VUL`.
* `trigger_type` - (Required) Trigger type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<scan_rule_id>`.
* `create_time` - Creation time
* `scan_rule_id` - Rule ID

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