---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerate_ip_endpoint_relation"
sidebar_current: "docs-alicloud-resource-ga-basic-accelerate-ip-endpoint-relation"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation resource.
---

# alicloud\_ga\_basic\_accelerate\_ip\_endpoint\_relation

Provides a Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation resource.

For information about Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation and how to use it, see [What is Basic Accelerate Ip Endpoint Relation](https://help.aliyun.com/document_detail/466842.html).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_basic_accelerate_ip_endpoint_relation" "default" {
  accelerator_id   = "your_accelerator_id"
  accelerate_ip_id = "your_accelerate_ip_id"
  endpoint_id      = "your_endpoint_id"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Basic GA instance.
* `accelerate_ip_id` - (Required, ForceNew) The ID of the Basic Accelerate IP.
* `endpoint_id` - (Required, ForceNew) The ID of the Basic Endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Accelerate Ip Endpoint Relation. It formats as `<accelerator_id>:<accelerate_ip_id>:<endpoint_id>`.
* `status` - The status of the Basic Accelerate Ip Endpoint Relation.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Basic Accelerate Ip Endpoint Relation.
* `delete` - (Defaults to 5 mins) Used when delete the Basic Accelerate Ip Endpoint Relation.

## Import

Global Accelerator (GA) Basic Accelerate Ip Endpoint Relation can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_accelerate_ip_endpoint_relation.example <accelerator_id>:<accelerate_ip_id>:<endpoint_id>
```
