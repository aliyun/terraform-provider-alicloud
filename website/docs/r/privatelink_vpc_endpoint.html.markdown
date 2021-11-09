---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint resource.
---

# alicloud\_privatelink\_vpc\_endpoint

Provides a Private Link Vpc Endpoint resource.

For information about Private Link Vpc Endpoint and how to use it, see [What is Vpc Endpoint](https://help.aliyun.com/document_detail/120479.html).

-> **NOTE:** Available in v1.109.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_privatelink_vpc_endpoint" "example" {
  service_id        = "YourServiceId"
  security_group_id = ["sg-ercx1234"]
  vpc_id            = "YourVpcId"
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run. Default to: `false`.
* `endpoint_description` - (Optional) The description of Vpc Endpoint. The length is 2~256 characters and cannot start with `http://` and `https://`.
* `security_group_ids` - (Required) The security group associated with the terminal node network card.
* `service_id` - (Optional, ForceNew) The terminal node service associated with the terminal node.
* `service_name` - (Optional, Computed, ForceNew) The name of the terminal node service associated with the terminal node.
* `vpc_endpoint_name` - (Optional) The name of Vpc Endpoint. The length is between 2 and 128 characters, starting with English letters or Chinese, and can include numbers, hyphens (-) and underscores (_).
* `vpc_id` - (Required, ForceNew) The private network to which the terminal node belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpc Endpoint. Value as `endpoint_id`.
* `bandwidth` - The Bandwidth.
* `connection_status` - The status of Connection.
* `endpoint_business_status` - The status of Endpoint Business.
* `endpoint_domain` - The Endpoint Domain.
* `status` - The status of Vpc Endpoint.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Vpc Endpoint.
* `update` - (Defaults to 4 mins) Used when update the Vpc Endpoint.

## Import

Private Link Vpc Endpoint can be imported using the id, e.g.

```
$ terraform import alicloud_privatelink_vpc_endpoint.example <endpoint_id>
```
