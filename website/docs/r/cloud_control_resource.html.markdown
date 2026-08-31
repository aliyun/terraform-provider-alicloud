---
subcategory: "Cloud Control"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_control_resource"
description: |-
  Provides a Alicloud Cloud Control Resource resource.
---

# alicloud_cloud_control_resource

Provides a Cloud Control Resource resource.



For information about Cloud Control Resource and how to use it, see [What is Resource](https://next.api.aliyun.com/document/cloudcontrol/2022-08-30/GetResourceType).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_control_resource&exampleId=5d11b39c-6938-1d3d-e711-ea7d2daad8738ad49fd3&activeTab=example&spm=docs.r.cloud_control_resource.0.5d11b39c69&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cloud_control_resource" "mq_instance" {
  desire_attributes = jsonencode({ "InstanceName" : "terraform-example-ons-instance" })
  product           = "Ons"
  resource_code     = "Instance"
}

resource "alicloud_cloud_control_resource" "default" {
  product           = "Ons"
  resource_code     = "Instance::Topic"
  resource_id       = alicloud_cloud_control_resource.mq_instance.resource_id
  desire_attributes = jsonencode({ "InstanceId" : "${alicloud_cloud_control_resource.mq_instance.resource_id}", "TopicName" : "terraform-example-ons-topic", "MessageType" : "1" })
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_control_resource&spm=docs.r.cloud_control_resource.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `desire_attributes` - (Optional, JsonString) Resource attributes specified when a user creates or updates a resource.
* `product` - (Required, ForceNew) The product Code represents the product to be operated. Currently supported products and resources can be queried at the following link: [supported-services-and-resource-types](https://help.aliyun.com/zh/cloud-control-api/product-overview/supported-services-and-resource-types).
* `resource_code` - (Required, ForceNew) Resource Code, if there is a parent resource, split with `::`, such as VPC::VSwitch. The supported resource Code can be obtained from the following link: [supported-services-and-resource-types](https://help.aliyun.com/zh/cloud-control-api/product-overview/supported-services-and-resource-types).
* `resource_id` (Optional, ForceNew) - If there is a parent resource, you need to enter the id of the parent resource, for example, in the VPC::VSwtich resource, you need to enter the id of the VPC: vpc-dexadfe3r4ad. If there are more than one level of parent resources, you need to use `:` to split.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<product>:<resource_code>:<resource_id>`.
* `resource_attributes` - The collection of properties for the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource.
* `delete` - (Defaults to 5 mins) Used when delete the Resource.
* `update` - (Defaults to 5 mins) Used when update the Resource.

## Import

Cloud Control Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_control_resource.example <provider>:<product>:<resource_code>:<resource_id>
```