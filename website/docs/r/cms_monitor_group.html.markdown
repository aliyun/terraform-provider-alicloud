---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_monitor_group"
sidebar_current: "docs-alicloud-resource-cms-monitor-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Monitor Group resource.
---

# alicloud_cms_monitor_group

Provides a Cloud Monitor Service Monitor Group resource.

For information about Cloud Monitor Service Monitor Group and how to use it, see [What is Monitor Group](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createmonitorgroup).

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_monitor_group&exampleId=6fbe8f1f-da7d-050f-d49c-b036c8707a64af17a6e3&activeTab=example&spm=docs.r.cms_monitor_group.0.6fbe8f1fda&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cms_monitor_group" "example" {
  monitor_group_name = "tf-example-accmonitorgroup"
}

resource "alicloud_cms_monitor_group" "default2" {
  contact_groups      = ["your_contact_groups"]
  resource_group_id   = "your_resource_group_id"
  resource_group_name = "resource_group_name"
  tags = {
    Created = "TF"
    For     = "Acceptance-example"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_monitor_group&spm=docs.r.cms_monitor_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `contact_groups` - (Optional) The alert group to which alert notifications will be sent.
* `monitor_group_name` - (Optional) The name of the application group.
* `resource_group_id` - (Optional, Available since v1.141.0) The ID of the resource group.
* `resource_group_name` - (Optional, Available since v1.141.0) The name of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Monitor Group.

## Import

Cloud Monitor Service Monitor Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_monitor_group.example <id>
```
