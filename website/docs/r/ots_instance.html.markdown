---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_instance"
sidebar_current: "docs-alicloud-resource-ots-instance"
description: |-
  Provides an OTS (Open Table Service) instance resource.
---

# alicloud_ots_instance

This resource will help you to manager a [Table Store](https://www.alibabacloud.com/help/doc-detail/27280.htm) Instance.
It is foundation of creating data table.

-> **NOTE:** Available since v1.10.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ots_instance" "default" {
  name             = "${var.name}-${random_integer.default.result}"
  description      = var.name
  network_type_acl = ["VPC"]
  tags = {
    Created = "TF"
    For     = "Building table"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the instance.
* `network_type_acl` - (Optional, Available since v1.221.0) The set of network types that are allowed access. Valid optional values:
    * `CLASSIC` - Classic network.
    * `VPC` - VPC network.
    * `INTERNET` - Public internet.

    Default to ["VPC", "CLASSIC"].
* `network_source_acl` - (Optional, Available since v1.221.0) The set of request sources that are allowed access. Valid optional values:
  * `TRUST_PROXY` - Trusted proxy, usually the Alibaba Cloud console.

  Default to ["TRUST_PROXY"].
* `accessed_by` - (Optional, Deprecated since v1.221.0) The network limitation of accessing instance. Valid values:
  * `Any` - Allow all network to access the instance.
  * `Vpc` - Only can the attached VPC allow to access the instance.
  * `ConsoleOrVpc` - Allow web console or the attached VPC to access the instance.

  Default to "Any".
* `resource_group_id` - (Optional, Available since v1.221.0) The resource group the instance belongs to.
  Default to Alibaba Cloud default resource group.
* `instance_type` - (Optional, ForceNew) The type of instance. Valid values are "Capacity" and "HighPerformance". Default to "HighPerformance".
* `description` - (Optional, ForceNew) The description of the instance. Currently, it does not support modifying.
* `tags` - (Optional) A mapping of tags to assign to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as the "name".

## Import

OTS instance can be imported using instance id or name, e.g.

```shell
$ terraform import alicloud_ots_instance.foo "my-ots-instance"
```

