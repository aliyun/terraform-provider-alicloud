---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_elasticity_assurance"
sidebar_current: "docs-alicloud-resource-ecs-elasticity-assurance"
description: |-
  Provides a Alicloud Ecs Elasticity Assurance resource.
---

# alicloud_ecs_elasticity_assurance

Provides a Ecs Elasticity Assurance resource.

For information about Ecs Elasticity Assurance and how to use it, see [What is Elasticity Assurance](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/createelasticityassurance).

-> **NOTE:** Available in v1.196.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_elasticity_assurance&exampleId=b14da29c-d75f-4580-6ea0-dffe7d201e0f620a1ac2&activeTab=example&spm=docs.r.ecs_elasticity_assurance.0.b14da29cd7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.c6"
}
resource "alicloud_ecs_elasticity_assurance" "default" {
  instance_amount                     = 1
  description                         = "before"
  zone_ids                            = [data.alicloud_zones.default.zones[2].id]
  private_pool_options_name           = "test_before"
  period                              = 1
  private_pool_options_match_criteria = "Open"
  instance_type                       = [data.alicloud_instance_types.default.instance_types.0.id]
  period_unit                         = "Month"
  assurance_times                     = "Unlimited"
  resource_group_id                   = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `assurance_times` - (ForceNew,Optional) The total number of times that the elasticity assurance can be applied. Set the value to Unlimited. This value indicates that the elasticity assurance can be applied an unlimited number of times within its effective duration. Default value: Unlimited.
* `description` - (Optional) Description of flexible guarantee service.
* `instance_amount` - (Required,ForceNew) The total number of instances for which to reserve the capacity of an instance type. Valid values: 1 to 1000.
* `instance_type` - (ForceNew,Required) Instance type. Currently, only one instance type is supported.
* `period` - (Optional) Length of purchase. The unit of duration is determined by the 'period_unit' parameter. Default value: 1.
  - When the `period_unit` parameter is set to Month, the valid values are 1, 2, 3, 4, 5, 6, 7, 8, and 9.
  - When the `period_unit` parameter is set to Year, the valid values are 1, 2, 3, 4, and 5.
* `period_unit` - (Optional) Duration unit. Value range:-Month: Month-Year: YearDefault value: Year
* `private_pool_options_match_criteria` - (ForceNew,Optional,Computed) The matching mode of flexible guarantee service. Possible values:-Open: flexible guarantee service for Open mode.-Target: specifies the flexible guarantee service of the mode.
* `private_pool_options_name` - (Optional,Computed) The name of the flexible protection service.
* `resource_group_id` - (ForceNew,Optional) The ID of the resource group.
* `start_time` - (ForceNew,Optional) Flexible guarantee service effective time.
* `tags` - (Optional) The tag key-value pair information bound by the elastic guarantee service.
* `zone_ids` - (ForceNew,Required) The zone ID of the region to which the elastic Protection Service belongs. Currently, only the creation of flexible protection services in one available area is supported.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `elasticity_assurance_id` - The first ID of the resource
* `end_time` - Flexible guarantee service failure time.
* `instance_charge_type` - The billing method of the instance. Possible value: PostPaid. Currently, only pay-as-you-go is supported.
* `start_time_type` - Flexible guarantee effective way. Possible values:-Now: Effective immediately.-Later: the specified time takes effect.
* `status` - The status of flexible guarantee services. Possible values:-Preparing: in preparation.-Prepared: to take effect.-Active: in effect.-Released: Released.
* `used_assurance_times` - This parameter is not yet available.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Elasticity Assurance.
* `update` - (Defaults to 5 mins) Used when update the Elasticity Assurance.

## Import

Ecs Elasticity Assurance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_elasticity_assurance.example <id>
```