---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_load_balancer_internet"
sidebar_current: "docs-alicloud-resource-sae-load-balancer-internet"
description: |-
  Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.
---

# alicloud\_sae\_load\_balancer\_internet

Provides an Alicloud Serverless App Engine (SAE) Application Load Balancer Attachment resource.

For information about Serverless App Engine (SAE) Load Balancer Internet Attachment and how to use it, see [alicloud_sae_load_balancer_internet](https://help.aliyun.com/document_detail/126360.html).

-> **NOTE:** Available in v1.164.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sae_load_balancer_internet" "example" {
  app_id          = "your_application_id"
  internet_slb_id = "your_internet_slb_id"
  internet {
    protocol    = "TCP"
    port        = 80
    target_port = 8080
  }
}

```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The target application ID that needs to be bound to the SLB.
* `internet_slb_id` - (Optional) The internet SLB ID.
* `internet` - (Required) The bound private network SLB. See the following `Block internet`.

### Block internet

The internet supports the following:

* `protocol` - (Optional) The Network protocol. Valid values: `TCP` ,`HTTP`,`HTTPS`.
* `https_cert_id` - (Optional) The SSL certificate. `https_cert_id` is required when HTTPS is selected
* `target_port` - (Optional) The Container port.
* `port` - (Optional) The SLB Port.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is the same as the application ID.
* `internet_ip` - Use designated public network SLBs that have been purchased to support non-shared instances.

## Import

The resource can be imported using the id, e.g.

```
$ terraform import alicloud_sae_load_balancer_internet.example <id>
```
