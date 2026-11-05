---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_rules"
sidebar_current: "docs-alicloud-datasource-pvtz-rules"
description: |-
  Provides a list of PrivateZone Rules to the user.
---

# alicloud_pvtz_rules

This data source provides the PrivateZone Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_pvtz_rules" "ids" {}

output "pvtz_rule_id_1" {
  value = data.alicloud_pvtz_rules.ids.rules.0.id
}

data "alicloud_pvtz_rules" "nameRegex" {
  name_regex = "^my-Rule"
}

output "pvtz_rule_id_2" {
  value = data.alicloud_pvtz_rules.nameRegex.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Optional, ForceNew) The ID of the Endpoint.
* `ids` - (Optional, ForceNew, Computed)  A list of Rule IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Rule names.
* `rules` - A list of PrivateZone Rules. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `endpoint_id` - The ID of the Endpoint.
  * `endpoint_name` - The Name of the Endpoint.
  * `forward_ips` - Forwarding targets. Each element contains the following attributes:
    * `ip` - The ip of the forwarding destination.
    * `port` - The port of the forwarding destination.
  * `id` - The ID of the Rule.
  * `rule_id` - The first ID of the resource.
  * `rule_name` - The name of the resource.
  * `type` - The type of the rule.
  * `zone_name` - The name of the forwarding zone.
  * `bind_vpcs` - The List of the VPC. Each element contains the following attributes:
    * `vpc_id` - The ID of the VPC.
    * `region_id` - The region ID of the vpc.
    * `vpc_name` - The Name of the VPC.
    * `region_name` - The Region Name of the vpc.
