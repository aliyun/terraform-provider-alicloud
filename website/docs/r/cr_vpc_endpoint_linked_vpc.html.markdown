---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_vpc_endpoint_linked_vpc"
sidebar_current: "docs-alicloud-resource-cr-vpc-endpoint-linked-vpc"
description: |-
  Provides a Alicloud CR Vpc Endpoint Linked Vpc resource.
---

# alicloud\_cr\_vpc\_endpoint\_linked\_vpc

Provides a CR Vpc Endpoint Linked Vpc resource.

For information about CR Vpc Endpoint Linked Vpc and how to use it, see [What is Vpc Endpoint Linked Vpc](https://www.alibabacloud.com/help/en/container-registry/latest/api-doc-cr-2018-12-01-api-doc-createinstancevpcendpointlinkedvpc).

-> **NOTE:** Available in v1.199.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cr_vpc_endpoint_linked_vpc" "default" {
  instance_id                      = "your_cr_instance_id"
  vpc_id                           = "your_vpc_id"
  vswitch_id                       = "your_vswitch_id"
  module_name                      = "Registry"
  enable_create_dns_record_in_pvzt = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `module_name` - (Required, ForceNew) The name of the module that you want to access. Valid Values:
  - `Registry`: the image repository.
  - `Chart`: a Helm chart.
* `enable_create_dns_record_in_pvzt` - (Optional) Specifies whether to automatically create an Alibaba Cloud DNS PrivateZone record. Valid Values:
  - `true`: automatically creates a PrivateZone record.
  - `false`: does not automatically create a PrivateZone record.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpc Endpoint Linked Vpc. It formats as `<instance_id>:<vpc_id>:<vswitch_id>:<module_name>`.
* `status` - The status of the Vpc Endpoint Linked Vpc.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Vpc Endpoint Linked Vpc.
* `delete` - (Defaults to 3 mins) Used when delete the Vpc Endpoint Linked Vpc.

## Import

CR Vpc Endpoint Linked Vpc can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_vpc_endpoint_linked_vpc.example <instance_id>:<vpc_id>:<vswitch_id>:<module_name>
```
