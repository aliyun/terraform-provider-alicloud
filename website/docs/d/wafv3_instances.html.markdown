---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_instances"
sidebar_current: "docs-alicloud-datasource-wafv3-instances"
description: |-
  Provides a list of Wafv3 Instance owned by an Alibaba Cloud account.
---

# alicloud_wafv3_instances

This data source provides Wafv3 Instance available to the user.[What is Instance](https://www.alibabacloud.com/help/en/web-application-firewall/latest/what-is-waf)

-> **NOTE:** Available since v1.200.0.

## Example Usage

```terraform
data "alicloud_wafv3_instances" "default" {
}

output "alicloud_wafv3_instance_example_id" {
  value = data.alicloud_wafv3_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, Available since 1.239.0) A list of WAF v3 instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance IDs.
* `instances` - A list of Instance Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `instance_id` - The first ID of the resource.
  * `id` - The ID of the resource.
  * `status` - The status of the resource.
