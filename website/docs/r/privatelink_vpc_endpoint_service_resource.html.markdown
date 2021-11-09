---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_resource"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-service-resource"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service Resource resource.
---

# alicloud\_privatelink\_vpc\_endpoint\_service\_resource

Provides a Private Link Vpc Endpoint Service Resource resource.

For information about Private Link Vpc Endpoint Service Resource and how to use it, see [What is Vpc Endpoint Service Resource](https://help.aliyun.com/document_detail/183548.html).

-> **NOTE:** Available in v1.110.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_privatelink_vpc_endpoint_service_resource" "example" {
  resource_id   = "lb-gw8nuym5xxxxx"
  resource_type = "slb"
  service_id    = "epsrv-gw8ii1xxxx"
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `resource_id` - (Required, ForceNew) The ID of Resource.
* `resource_type` - (Required, ForceNew) The Type of Resource.
* `service_id` - (Required, ForceNew) The ID of Vpc Endpoint Service.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Service Resource. The value is formatted `<service_id>:<resource_id>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 4 mins) Used when create the Vpc Endpoint Service Resource.

## Import

Private Link Vpc Endpoint Service Resource can be imported using the id, e.g.

```
$ terraform import alicloud_privatelink_vpc_endpoint_service_resource.example <service_id>:<resource_id>
```
