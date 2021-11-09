---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service_user"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-service-user"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service User resource.
---

# alicloud\_privatelink\_vpc\_endpoint\_service\_user

Provides a Private Link Vpc Endpoint Service User resource.

For information about Private Link Vpc Endpoint Service User and how to use it, see [What is Vpc Endpoint Service User](https://help.aliyun.com/document_detail/183545.html).

-> **NOTE:** Available in v1.110.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_privatelink_vpc_endpoint_service_user" "example" {
  service_id = "epsrv-gw81c6xxxxxx"
  user_id    = "YourRamUserId"
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `service_id` - (Required, ForceNew) The Id of Vpc Endpoint Service.
* `user_id` - (Required, ForceNew) The Id of Ram User.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Service User. The value is formatted `<service_id>:<user_id>`.

## Import

Private Link Vpc Endpoint Service User can be imported using the id, e.g.

```
$ terraform import alicloud_privatelink_vpc_endpoint_service_user.example <service_id>:<user_id>
```
