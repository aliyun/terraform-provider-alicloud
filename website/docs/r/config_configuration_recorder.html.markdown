---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_configuration_recorder"
sidebar_current: "docs-alicloud-resource-config-configuration-recorder"
description: |-
  Provides a Alicloud Config Configuration Recorder resource.
---

# alicloud_config_configuration_recorder

Provides a Alicloud Config Configuration Recorder resource. Cloud Config is a specialized service for evaluating resources. Cloud Config tracks configuration changes of your resources and evaluates configuration compliance. Cloud Config can help you evaluate numerous resources and maintain the continuous compliance of your cloud infrastructure.
For information about Alicloud Config Configuration Recorder and how to use it, see [What is Configuration Recorder.](https://www.alibabacloud.com/help/en/cloud-config/latest/startconfigurationrecorder)

-> **NOTE:** Available since v1.99.0.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-southeast-1`.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_configuration_recorder&exampleId=5d26c3ee-261e-af62-66f0-00f852c8373743e4337a&activeTab=example&spm=docs.r.config_configuration_recorder.0.5d26c3ee26&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_config_configuration_recorder" "example" {
  resource_types = [
    "ACS::ECS::Instance",
    "ACS::ECS::Disk"
    # other resource types ...
  ]
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_config_configuration_recorder&spm=docs.r.config_configuration_recorder.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `enterprise_edition` - (Optional, ForceNew) - Whether to use the enterprise version configuration audit. Valid values: `true` and `false`. Default value `false`. For enterprise accounts, We recommend you to use the resource [alicloud_config_aggregator](https://www.terraform.io/docs/providers/alicloud/r/config_aggregator).
* `resource_types` - (Optional) A list of resource types to be monitored. [Resource types that support Cloud Config.](https://www.alibabacloud.com/help/en/doc-detail/127411.htm)
  * If you use an ordinary account, the `resource_types` supports the update operation after the process of creation is completed.
  * If you use an enterprise account, the `resource_types` does not support updating. 

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Config Configuration Recorder. Value as alicloud account ID.
* `status` - Status of resource monitoring. Values: `REGISTRABLE`: Not registered, `BUILDING`: Under construction, `REGISTERED`: Registered and `REBUILDING`: Rebuilding.
* `organization_enable_status` - Enterprise version configuration audit enabled status. Values: `REGISTRABLE`: Not enabled, `BUILDING`: Building and `REGISTERED`: Enabled.
* `organization_master_id` - The ID of the Enterprise management account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `update` - (Defaults to 11 mins) Used when updating the Config Configuration Recorder.

## Import

Alicloud Config Configuration Recorder can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_configuration_recorder.example 122378463********
```
