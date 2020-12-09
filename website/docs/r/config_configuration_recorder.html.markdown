---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_configuration_recorder"
sidebar_current: "docs-alicloud-resource-config-configuration-recorder"
description: |-
  Provides a Alicloud Config Configuration Recorder resource.
---

# alicloud\_config\_configuration\_recorder

Provides a Alicloud Config Configuration Recorder resource. Cloud Config is a specialized service for evaluating resources. Cloud Config tracks configuration changes of your resources and evaluates configuration compliance. Cloud Config can help you evaluate numerous resources and maintain the continuous compliance of your cloud infrastructure.
For information about Alicloud Config Configuration Recorder and how to use it, see [What is Cloud Config.](https://www.alibabacloud.com/help/en/doc-detail/127388.htm)

-> **NOTE:** Available in v1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-northeast-1`.

## Example Usage

```terraform
resource "alicloud_config_configuration_recorder" "example" {
  resource_types = [
    "ACS::ECS::Instance",
    "ACS::ECS::Disk"
    # other resource types ...
  ]
}
```
## Argument Reference

The following arguments are supported:

* `enterprise_edition` - (Optional, ForceNew) - Whether to use the enterprise version configuration audit. Valid values: `true` and `fales`. Default value `false`.
* `resource_types` - (Optional) A list of resource types to be monitored. [Resource types that support Cloud Config.](https://www.alibabacloud.com/help/en/doc-detail/127411.htm)

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Config Configuration Recorder. Value as alicloud account ID.
* `status` - Enterprise version configuration audit enabled status. Values: `REGISTRABLE`: Not registered, `BUILDING`: Under construction, `REGISTERED`: Registered and `REBUILDING`: Rebuilding.
* `organization_enable_status` - Status of resource monitoring. Values: `REGISTRABLE`: Not enabled, `BUILDING`: Building and `REGISTERED`: Enabled.
* `organization_master_id` - The ID of the Enterprise management account.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 11 mins) Used when updating the Config Configuration Recorder.

## Import

Alicloud Config Configuration Recorder can be imported using the id, e.g.

```
$ terraform import alicloud_config_configuration_recorder.example 122378463********
```
