---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_resource_attachment"
sidebar_current: "docs-alicloud-resource-ga-resource-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Resource Attachment resource.
---

# alicloud_ga_resource_attachment

Provides a Global Accelerator (GA) Resource Attachment resource. It associates a cloud resource with a GA instance.

For information about Global Accelerator (GA) Resource Attachment and how to use it, see [What is Resource Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-associateresources).

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** The association API accepts only `cn-hangzhou` as the GA instance region, so configure the provider with `region = "cn-hangzhou"`.

-> **NOTE:** The GA WAF integration uses a WAF 3.0 instance. The example below reads an existing instance through the `alicloud_wafv3_instances` data source; create one with `alicloud_wafv3_instance` first if the account does not have a WAF instance.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "example" {}

resource "alicloud_ga_accelerator" "example" {
  accelerator_name       = "terraform-example"
  description            = "terraform-example"
  bandwidth_billing_type = "CDT"
  payment_type           = "PayAsYouGo"
  bandwidth              = 200
}

resource "alicloud_ga_resource_attachment" "example" {
  accelerator_id                = alicloud_ga_accelerator.example.id
  associated_resource_type      = "WAF"
  associated_resource_region_id = "cn-hangzhou"
  associated_resource_id        = data.alicloud_wafv3_instances.example.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the GA instance. Changing this forces a new resource.
* `associated_resource_type` - (Required, ForceNew) The type of the associated resource. Valid values: `WAF`. Changing this forces a new resource.
* `associated_resource_region_id` - (Required, ForceNew) The region ID of the associated resource. Changing this forces a new resource.
* `associated_resource_id` - (Required, ForceNew) The ID of the associated resource. For `WAF`, it is the ID of the Web Application Firewall instance. Changing this forces a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource attachment, in the format `<accelerator_id>:<associated_resource_type>:<associated_resource_region_id>:<associated_resource_id>`.
* `associated_mode` - The association mode. Valid values: `Associated` and `Managed`. The association API currently uses `Associated` (loosely coupled) for WAF.
* `state` - The association state. Valid values: `active`, `binding`, `unBinding` and `wafRelease`. `wafRelease` means the associated WAF instance has been released.

## Import

Global Accelerator Resource Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_resource_attachment.example <accelerator_id>:<associated_resource_type>:<associated_resource_region_id>:<associated_resource_id>
```
