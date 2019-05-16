---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_rules"
sidebar_current: "docs-alicloud-datasource-slb-rules"
description: |-
    Provides a list of server load balancer rules to the user.
---

# alicloud\_slb_rules

This data source provides the rules associated with a server load balancer listener.

## Example Usage

```
data "alicloud_slb_rules" "sample_ds" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  frontend_port = 80
}

output "first_slb_rule_id" {
  value = "${data.alicloud_slb_rules.sample_ds.slb_rules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - ID of the SLB with listener rules.
* `frontend_port` - SLB listener port.
* `ids` - (Optional) A list of rules IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB listener rules IDs.
* `names` - A list of SLB listener rules names.
* `slb_rules` - A list of SLB listener rules. Each element contains the following attributes:
  * `id` - Rule ID.
  * `name` - Rule name.
  * `domain` - Domain name in the HTTP request where the rule applies (e.g. "*.aliyun.com").
  * `url` - Path in the HTTP request where the rule applies (e.g. "/image").
  * `server_group_id` - ID of the linked VServer group.
