---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_endpoint_acl_policy"
sidebar_current: "docs-alicloud-resource-cr-endpoint-acl-policy"
description: |-
  Provides a Alicloud CR Endpoint Acl Policy resource.
---

# alicloud\_cr\_endpoint\_acl\_policy

Provides a CR Endpoint Acl Policy resource.

For information about CR Endpoint Acl Policy and how to use it, see [What is Endpoint Acl Policy](https://www.alibabacloud.com/help/doc-detail/145275.htm).

-> **NOTE:** Available in v1.139.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
}
data "alicloud_cr_ee_instances" "default" {}
resource "alicloud_cr_ee_instance" "default" {
  count          = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
  module_name   = "Registry"
}
resource "alicloud_cr_endpoint_acl_policy" "default" {
  instance_id   = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
  entry         = "192.168.1.0/24"
  description   = var.name
  module_name   = "Registry"
  endpoint_type = "internet"
  depends_on    = [data.alicloud_cr_endpoint_acl_service.default]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the entry.
* `endpoint_type` - (Required, ForceNew) The type of endpoint. Valid values: `internet`.
* `entry` - (Required, ForceNew) The IP segment that allowed to access.
* `instance_id` - (Required, ForceNew) The ID of the CR Instance.
* `module_name` - (Optional, ForceNew) The module that needs to set the access policy. Valid values: `Registry`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Endpoint Acl Policy. The value formats as `<instance_id>:<endpoint_type>:<entry>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Acl Policy.
* `delete` - (Defaults to 10 mins) Used when delete the Acl Policy.

## Import

CR Endpoint Acl Policy can be imported using the id, e.g.

```
$ terraform import alicloud_cr_endpoint_acl_policy.example <instance_id>:<endpoint_type>:<entry>
```
