---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_session_manager_status"
description: |-
  Provides a Alicloud ECS Session Manager Status resource.
---

# alicloud_ecs_session_manager_status

Provides a ECS Session Manager Status resource.

For information about ECS Session Manager Status and how to use it, see [What is Session Manager Status](https://www.alibabacloud.com/help/zh/doc-detail/337915.html).

-> **NOTE:** Available since v1.148.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_session_manager_status&exampleId=d114a94f-1771-6fdc-574a-7f00d8d3d4abdc52fd36&activeTab=example&spm=docs.r.ecs_session_manager_status.0.d114a94f17&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecs_session_manager_status" "default" {
  session_manager_status_name = "sessionManagerStatus"
  status                      = "Disabled"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_session_manager_status&spm=docs.r.ecs_session_manager_status.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `session_manager_status_name` - (Required, ForceNew) The name of the Session Manager Status. Valid values: `sessionManagerStatus`.
* `status` - (Required) The status of the Session Manager Status. Valid values: `Enabled`, `Disabled`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Session Manager Status.

## Import

ECS Session Manager Status can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_session_manager_status.example <id>
```
