---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_instances"
sidebar_current: "docs-alicloud-datasource-waf-instances"
description: |-
  Provides a datasource to retrieve WAF instances.
---

# alicloud\_waf\_instances

Provides a WAF datasource to retrieve instances.

For information about WAF and how to use it, see [What is Alibaba Cloud WAF](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available since v1.90.0.

## Example Usage

```terraform
data "alicloud_waf_instances" "default" {
  ids               = ["waf-cn-09k********"]
  status            = "1"
  resource_group_id = "rg-acfmwvv********"
  instance_source   = "waf-cloud"
}

output "the_first_waf_instance_id" {
  value = data.alicloud_waf_instances.default.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of WAF instance IDs. 
* `status` - (Optional) The status of WAF instance to filter results. Optional value: `0`: The instance has expired, `1` : The instance has not expired and is working properly.
* `resource_group_id` - (Optional) The ID of resource group to which WAF instance belongs.
* `instance_source` - (Optional) The source of the WAF instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of WAF instance IDs. 
* `instances` - A list of WAF instances. Each element contains the following attributes:
  * `id` - The ID of the WAF instance.
  * `instance_id` - The ID of WAF the instance.
  * `end_date` - The timestamp (in seconds) indicating when the WAF instance expires.
  * `in_debt` - Indicates whether the WAF instance has overdue payments.
  * `remain_day` - The number of days before the trial period of the WAF instance expires.
  * `trial` - Indicates whether this is a trial instance.
  * `status` - Indicates whether the WAF instance has expired.
  * `subscription_typed` - The billing method of the WAF instance. 
