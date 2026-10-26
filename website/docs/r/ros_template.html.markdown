---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_template"
sidebar_current: "docs-alicloud-resource-ros-template"
description: |-
  Provides a Alicloud ROS Template resource.
---

# alicloud_ros_template

Provides a ROS Template resource.

For information about ROS Template and how to use it, see [What is Template](https://www.alibabacloud.com/help/en/doc-detail/141851.htm).

-> **NOTE:** Available since v1.108.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ros_template&exampleId=f43f08c5-570b-91ee-2cd6-2465cfeddbca67cc9b02&activeTab=example&spm=docs.r.ros_template.0.f43f08c557&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ros_template" "example" {
  template_name = "example_value"
  template_body = <<EOF
    {
        "ROSTemplateFormatVersion": "2015-09-01"
    }
    EOF
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ros_template&spm=docs.r.ros_template.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `description` - (Optional, String) The description of the template. It can be up to 256 characters in length.
* `resource_group_id` - (Optional, Computed, String, Available since v2.2.0) The ID of the resource group.
* `rotate_strategy` - (Optional, String, Available since v2.2.0) The rotation policy for template versions. Valid values: `None` (default): no rotation, and an error is reported when the number of template versions reaches the upper limit. `DeleteOldestNonSharedVersionWhenLimitExceeded`: rotation deletes non-shared template versions. Defaults to `None`.

  -> **NOTE:** `rotate_strategy` is not returned by the `GetTemplate` operation. Re-apply it after importing the resource if it is configured.

* `tags` - (Optional, Map) A mapping of tags to assign to the resource.
* `template_body` - (Optional) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length.  If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.  You must specify one of the TemplateBody and TemplateURL parameters, but you cannot specify both of them.
* `template_name` - (Required, String) The name of the template. The name can be up to 255 characters in length (UTF-8 encoding). It must start with a digit, letter, or Chinese character, and can contain digits, letters, Chinese characters, hyphens (-), and underscores (_).
* `template_url` - (Optional, String) The URL of the file that contains the template body. The URL must point to a template located in an HTTP or HTTPS web server or an Alibaba Cloud OSS bucket. Examples: oss://ros/template/demo and oss://ros/template/demo?RegionId=cn-hangzhou. The template can be up to 524,288 bytes in length, and the URL can be up to 1,024 bytes in length.

  -> **NOTE:** If the region of the OSS bucket is not specified, the RegionId value is used. You can specify only one of the `template_body` and `template_url` parameters.

  -> **NOTE:** The content behind `template_url` is not returned by the `GetTemplate` operation. To remove `template_url`, set `template_body` in the same change. After importing a template created from a URL, re-apply the configuration if the content or the URL differs.

* `validation_options` - (Optional, List, Available since v2.2.0) The validation options. Default: no option is enabled and strict validation is performed. Valid values: `BasicValidation`, which performs only basic structural validation on the template content. A maximum of 1 item is supported.

  -> **NOTE:** `validation_options` only applies to the template content that is carried by the same apply, and it is not returned by the `GetTemplate` operation. After the import, confirm it in the configuration and apply again if the value differs from the default.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Template. Value as `template_id`.

## Import

ROS Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_ros_template.example <template_id>
```
