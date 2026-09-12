---
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_aigw_instance"
sidebar_current: "docs-alicloud-resource-esa-aigw-instance"
description: |-
  Provides an ESA AI Gateway Instance resource.
---


# alicloud_esa_aigw_instance

Provides an ESA AI Gateway Instance resource.

For information about ESA AI Gateway Instance and how to use it, refer to Alibaba Cloud ESA documentation.

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_esa_aigw_instance" "default" {
  aigw_instance_name = "tf-example-aigw"
  comment            = "terraform example"
}
```

## Argument Reference

The following arguments are supported:

* `aigw_instance_name` - (Required, ForceNew) The name of the AI Gateway instance. The name must start with a lowercase letter and end with a lowercase letter or digit, and can only contain lowercase letters, digits, and hyphens. 1-128 characters.
* `comment` - (Optional, Computed) The remark of the AI Gateway instance.
* `enable_auth` - (Optional, Computed) Whether to enable AuthKey authentication for the AI Gateway instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform, equals to the `aigw_instance_id`.
* `aigw_instance_id` - The ID of the AI Gateway instance.
* `auth_key` - The authentication key of the AI Gateway instance.
* `create_time` - The creation time of the AI Gateway instance.
* `record_count` - The number of domains bound to the AI Gateway instance.
* `status` - The status of the AI Gateway instance. Valid values: `Pending`, `Deploying`, `Running`, `Failed`, `Updating`, `Deleting`, `Deleted`.
* `update_time` - The update time of the AI Gateway instance.

## Import

ESA AI Gateway Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_aigw_instance.example <aigw_instance_id>
```
