---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_load_balancer_intranet"
sidebar_current: "docs-alicloud-resource-sae-load-balancer-intranet"
description: |-
  Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.
---

# alicloud\_sae\_load\_balancer\_intranet

Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.

For information about Serverless App Engine (SAE) Load Balancer Intranet Attachment and how to use it, see [alicloud_sae_load_balancer_intranet](https://help.aliyun.com/document_detail/126360.html).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sae_load_balancer_intranet" "example" {
  app_id          = "your_application_id"
  intranet_slb_id = "intranet_slb_id"
  intranet {
    protocol    = "TCP"
    port        = 80
    target_port = 8080
  }
}

```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The target application ID that needs to be bound to the SLB.
* `intranet_slb_id` - (Optional) The intranet SLB ID.
* `intranet` - (Required) The bound private network SLB. See the following `Block intranet`.

### Block intranet

The intranet supports the following:

* `protocol` - (Optional) The Network protocol. Valid values: `TCP` ,`HTTP`,`HTTPS`.
* `https_cert_id` - (Optional) The SSL certificate. `https_cert_id` is required when HTTPS is selected
* `target_port` - (Optional) The Container port.
* `port` - (Optional) The SLB Port.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the same as the application ID.
* `intranet_ip` - Use designated private network SLBs that have been purchased to support non-shared instances.

## Import

The resource can be imported using the id, e.g.

```
$ terraform import alicloud_sae_load_balancer_intranet.example <id>
```
