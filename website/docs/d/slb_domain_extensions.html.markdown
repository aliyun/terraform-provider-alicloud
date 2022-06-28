---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_domain_extensions"
sidebar_current: "docs-alicloud-resource-slb-domain-extensions"
description: |-
  Provides a Load Banlancer domain extension Resource and add it to one Listener.
---

# alicloud\_slb\_domain_extensions

This data source provides the domain extensions associated with a server load balancer listener.

-> **NOTE:** Available in 1.60.0+

## Example Usage
```
data "alicloud_slb_domain_extensions" "foo" {
  ids               = ["fake-de-id"]
  load_balancer_id  = "fake-lb-id"
  frontend_port     = "fake-port"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) IDs of the SLB domain extensions.
* `load_balancer_id` - (Required) The ID of the SLB instance.
* `frontend_port` - (Required) The frontend port used by the HTTPS listener of the SLB instance. Valid values: 1–65535.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `extensions` - A list of SLB domain extension. Each element contains the following attributes:
    * `id` - The ID of the domain extension.
    * `domain` - The domain name.
    * `server_certificate_id` - The ID of the certificate used by the domain name.
