---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_aigw_record"
description: |-
  Provides a Alicloud ESA AIGW Record resource.
---

# alicloud_esa_aigw_record

Provides a ESA AIGW Record resource.

AI Gateway (AIGW) record is a domain name mapping entry bound to an ESA AI Gateway instance. It associates a domain name with a gateway instance to route AI API traffic through the gateway.

For information about ESA AIGW Record and how to use it, see [What is AIGW Record](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateAIGWRecord).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "aigw_instance_id" {
  default = "your-aigw-instance-id"
}

resource "alicloud_esa_aigw_record" "default" {
  site_id     = "your-site-id"
  instance_id = var.aigw_instance_id
  record_name = "aigw.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The ID of the ESA site.
* `instance_id` - (Required, ForceNew) The ID of the AI Gateway instance.
* `record_name` - (Required, ForceNew) The domain name of the AIGW record.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in format `instance_id:record_name`.
* `create_time` - The creation time of the resource.

## Import

ESA AIGW Record can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_aigw_record.example <instance_id>:<record_name>
```
