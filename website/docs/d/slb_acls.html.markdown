---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_acls"
sidebar_current: "docs-alicloud-datasource-slb-acls"
description: |-
    Provides a list of server load balancer acls (access control lists) to the user.
---

# alicloud\_slb_acls

This data source provides the acls in the region.

## Example Usage

```
data "alicloud_slb_acls" "sample_ds" {
}

output "first_slb_acl_id" {
  value = "${data.alicloud_slb_acls.sample_ds.acls.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of acls IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by acl name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew, Available in 1.60.0+) The Id of resource group which acl belongs.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB acls IDs.
* `names` - A list of SLB acls names.
* `acls` - A list of SLB  acls. Each element contains the following attributes:
  * `id` - Acl ID.
  * `name` - Acl name.
  * `entry_list` - A list of entry (IP addresses or CIDR blocks).  Each entry contains two sub-fields as `Entry Block` follows.
  * `related_listeners` - A list of listener are attached by the acl.  Each listener contains four sub-fields as `Listener Block` follows.
  * `tags` - A mapping of tags to assign to the resource.
  * `resource_group_id` - Resource group ID.
## Entry Block

The entry mapping supports the following:

* `entry`   - An IP addresses or CIDR blocks.
* `comment` - the comment of the entry.

## Listener Block

The Listener mapping supports the following:

* `load_balancer_id` - the id of load balancer instance, the listener belongs to.
* `frontend_port` - the listener port.
* `protocol`      - the listener protocol (such as tcp/udp/http/https, etc).
* `acl_type`      - the type of acl (such as white/black).
