---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_application_scale"
sidebar_current: "docs-alicloud-resource-edas_application_scale"
description: |-
  Provides a Alicloud EDAS Application Scale resource.
---

# alicloud\_edas\_application\_scale

Provides an Alicloud EDAS Application Scale resource.

-> **NOTE:** Available in 1.81.0+

## Example Usage

Basic Usage

```
resource "alicloud_edas_application_scale" "default" {
  app_id = var.app_id
  deploy_group = var.deploy_group
  ecu_info = var.ecu_info
  force_status = var.force_status
}

```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application that you want to deploy.
* `deploy_group` - (Required, ForceNew) The ID of the instance group to which you want to add ECS instances to scale out the application.
* `ecu_info` - (Required, ForceNew) The IDs of the Elastic Compute Unit (ECU) where you want to deploy the application. Type: List.
* `force_status` - (Optional) This parameter specifies whether to forcibly remove an ECS instance where the application is deployed. It is set as true only after the ECS instance expires. In normal cases, this parameter do not need to be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<ecu1,ecu2>`.
* `ecc_info` - The ecc information of the resource supplied above. The value is formulated as `<ecc1,ecc2>`.

