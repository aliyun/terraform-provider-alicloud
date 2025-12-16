---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_template"
sidebar_current: "docs-alicloud-resource-ros-template"
description: |-
  Provides a Alicloud ROS Template resource.
---

# alicloud\_ros\_template

Provides a ROS Template resource.

For information about ROS Template and how to use it, see [What is Template](https://www.alibabacloud.com/help/en/doc-detail/141851.htm).

-> **NOTE:** Available in v1.108.0+.

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

* `description` - (Optional) The description of the template. The description can be up to 256 characters in length.
* `template_body` - (Optional) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length.  If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.  You must specify one of the TemplateBody and TemplateURL parameters, but you cannot specify both of them.
* `template_name` - (Required) The name of the template. The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `template_url` - (Optional) The template url.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Template. Value as `template_id`.

## Import

ROS Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_ros_template.example <template_id>
```
