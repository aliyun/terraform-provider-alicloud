---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_custom_deployment_set"
description: |-
  Provides a Alicloud RDS Custom Deployment Set resource.
---

# alicloud_rds_custom_deployment_set

Provides a RDS Custom Deployment Set resource.

Custom Deployment set.

For information about RDS Custom Deployment Set and how to use it, see [What is Custom Deployment Set](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_custom_deployment_set&exampleId=b629c950-ef86-3b30-c341-a7bfda7554adb23efa76&activeTab=example&spm=docs.r.rds_custom_deployment_set.0.b629c950ef&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}


resource "alicloud_rds_custom_deployment_set" "default" {
  on_unable_to_redeploy_failed_instance = "CancelMembershipAndStart"
  custom_deployment_set_name            = var.name
  description                           = "2024:11:19 1010:1111:0808"
  group_count                           = "3"
  strategy                              = "Availability"
}
```

## Argument Reference

The following arguments are supported:
* `custom_deployment_set_name` - (Optional, ForceNew) The name of the resource
* `description` - (Optional, ForceNew) Deployment set description information. It must be 2 to 256 characters in length and cannot start with http:// or https.
* `group_count` - (Optional, Int) Set the number of groups for the deployment set group high availability policy. Value range: 1~7.

  Default value: 3

-> **NOTE:**  This parameter takes effect only when 'Strategy = AvailabilityGroup.

* `on_unable_to_redeploy_failed_instance` - (Optional) After the instance in the deployment set is down and migrated, there is no emergency solution for the scattered instance inventory. Value range:
  - `CancelMembershipAndStart`: removes the instance from the deployment set and starts the instance immediately after the instance is down and migrated.
  - `KeepStopped`: The deployment set of the instance is maintained. The instance remains in the stopped state.

  Default value: CancelMembershipAndStart.
* `strategy` - (Optional, ForceNew, Computed) Deployment strategy. Value range:
  - `Availability`: High Availability policy.
  - `AvailabilityGroup`: the high availability policy of the deployment set group.
  - `LowLatency`: Network low latency policy.

  Default value: Availability.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom Deployment Set.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Deployment Set.

## Import

RDS Custom Deployment Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_custom_deployment_set.example <id>
```