---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_image_component"
description: |-
  Provides a Alicloud ECS Image Component resource.
---

# alicloud_ecs_image_component

Provides a ECS Image Component resource.



For information about ECS Image Component and how to use it, see [What is Image Component](https://www.alibabacloud.com/help/en/doc-detail/200424.htm).

-> **NOTE:** Available since v1.159.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_image_component&exampleId=8f92edce-12ca-de82-52fd-cf441590ea9a6d95664b&activeTab=example&spm=docs.r.ecs_image_component.0.8f92edce12&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_ecs_image_component" "example" {
  component_type       = "Build"
  content              = "RUN yum update -y"
  description          = "example_value"
  image_component_name = "example_value"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  system_type          = "Linux"
  tags = {
    Created = "TF"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_image_component&spm=docs.r.ecs_image_component.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `component_type` - (Optional, ForceNew, Computed) The component type. Supports mirrored build components and test components.

  Value range:
  - Build
  - Test

  Default value: Build.

-> **NOTE:**  Build components can only be used in build templates and test components can only be used in test templates.

* `component_version` - (Optional, ForceNew, Available since v1.235.0) The component version number, which is used in conjunction with the component name, is in the format of major.minor.patch and is a non-negative integer.

  Default value:(x +1).0.0, where x is the maximum major version of the current component.
* `content` - (Required, ForceNew) Component content. Consists of multiple commands. The maximum number of commands cannot exceed 127. Details of supported commands and command formats,
* `description` - (Optional, ForceNew) Description information. It must be 2 to 256 characters in length and cannot start with http:// or https.
* `image_component_name` - (Optional, ForceNew, Computed) The component name. It must be 2 to 128 characters in length and start with an uppercase letter or a Chinese character. It cannot start with http:// or https. Can contain Chinese, English, numbers, half-length colons (:), underscores (_), half-length periods (.), or dashes (-).

-> **NOTE:**  When 'Name' is not set, the 'ImageComponentId' return value is used by default.

* `resource_group_id` - (Optional, Computed) The ID of the enterprise resource group to which the created image component belongs.
* `system_type` - (Optional, ForceNew, Computed) The operating system supported by the component.

  Value range:
  - Linux
  - Windows

  Default value: Linux.
* `tags` - (Optional, Map) List of label key-value pairs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Component creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image Component.
* `delete` - (Defaults to 5 mins) Used when delete the Image Component.
* `update` - (Defaults to 5 mins) Used when update the Image Component.

## Import

ECS Image Component can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_image_component.example <id>
```