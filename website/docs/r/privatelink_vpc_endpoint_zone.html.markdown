---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_zone"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-zone"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Zone resource.
---

# alicloud\_privatelink\_vpc\_endpoint\_zone

Provides a Private Link Vpc Endpoint Zone resource.

For information about Private Link Vpc Endpoint Zone and how to use it, see [What is Vpc Endpoint Zone](https://help.aliyun.com/document_detail/183561.html).

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_privatelink_vpc_endpoint_zone" "example" {
  endpoint_id = "ep-gw8boxxxxx"
  vswitch_id  = "vsw-rtycxxxxx"
  zone_id     = "eu-central-1a"
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `endpoint_id` - (Required, ForceNew) The ID of the Vpc Endpoint.
* `vswitch_id` - (Required, ForceNew) The VSwitch id.
* `zone_id` - (Optional, Computed, ForceNew) The Zone Id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Zone. The value is formatted `<endpoint_id>:<zone_id>`.
* `status` - Status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 4 mins) Used when create the Vpc Endpoint Zone.
* `delete` - (Defaults to 4 mins) Used when delete the Vpc Endpoint Zone.

## Import

Private Link Vpc Endpoint Zone can be imported using the id, e.g.

```
$ terraform import alicloud_privatelink_vpc_endpoint_zone.example <endpoint_id>:<zone_id>
```
