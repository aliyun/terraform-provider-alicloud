---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_app_group"
description: |-
  Provides a Alicloud Schedulerx App Group resource.
---

# alicloud_schedulerx_app_group

Provides a Schedulerx App Group resource.



For information about Schedulerx App Group and how to use it, see [What is App Group](https://www.alibabacloud.com/help/en/schedulerx/schedulerx-serverless/developer-reference/api-schedulerx2-2019-04-30-createappgroup).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_schedulerx_app_group&exampleId=0e641d36-d54e-e4f5-97e1-7dfcc47b79c3d98e3e34&activeTab=example&spm=docs.r.schedulerx_app_group.0.0e641d36d5&intl_lang=EN_US" target="_blank">
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

resource "alicloud_schedulerx_namespace" "CreateNameSpace" {
  namespace_name = var.name
  description    = var.name
}

resource "alicloud_schedulerx_app_group" "default" {
  max_jobs              = "100"
  monitor_contacts_json = jsonencode([{ "name" : "contact-group-1" }, { "name" : "contact-group-2" }])
  delete_jobs           = "false"
  app_type              = "1"
  namespace_source      = "schedulerx"
  group_id              = "example-appgroup-pop-autoexample"
  namespace_name        = "default"
  description           = var.name
  monitor_config_json   = jsonencode({ "sendChannel" : "sms,ding", "alarmType" : "Contacts", "webhookIsAtAll" : "false" })
  app_version           = "1"
  app_name              = "example-appgroup-pop-autoexample"
  namespace             = alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid
  enable_log            = "false"
  schedule_busy_workers = "false"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_schedulerx_app_group&spm=docs.r.schedulerx_app_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `app_name` - (Required, ForceNew) Application Name
* `app_type` - (Optional, Int) Application type.
  - 1, general application.
  - 2, k8s application.
* `app_version` - (Optional) Application Version, 1: Basic Edition, 2: Professional Edition
* `delete_jobs` - (Optional) Whether to delete the task in the application Group. The values are as follows:
  - `true`: Delete.
  - `false`: Do not delete.
* `description` - (Optional) Application Description
* `enable_log` - (Optional) Whether to enable the log.
  - true: On
  - false: Close
* `group_id` - (Required, ForceNew) Application ID
* `max_concurrency` - (Optional, Int) The maximum number of instances running at the same time. The default value is 1, that is, the last trigger is not completed, and the next trigger will not be performed even at the running time.
* `max_jobs` - (Optional, ForceNew, Int) Application Grouping Configurable Maximum Number of Tasks
* `monitor_config_json` - (Optional) Alarm configuration JSON field. Supported keys include `sendChannel` (alarm channels, e.g. `"sms,ding"`), `alarmType` (alarm type, e.g. `"Contacts"` or `"CustomContacts"`), and `webhookIsAtAll` (whether webhook @all). **Note:** When `monitor_contacts_json` is specified, `alarmType` must be explicitly included in `monitor_config_json` (typically `"CustomContacts"` for custom contacts or `"Contacts"` for contact groups); otherwise the API will automatically append `alarmType` which causes configuration drift on subsequent plans.
* `monitor_contacts_json` - (Optional) Alarm contact JSON format. **Note:** This field only takes effect when `monitor_config_json` contains an `alarmType` value (e.g. `"CustomContacts"` or `"Contacts"`). The format depends on `alarmType`: for `"CustomContacts"`, use `[{"userName":"name","userPhone":"phone","ding":"webhook_url"}]`; for `"Contacts"`, use `[{"name":"contact_group_name"}]`.
* `namespace` - (Required, ForceNew) The namespace ID, which is obtained on the namespace page of the console.
* `namespace_name` - (Required) The namespace name.
* `namespace_source` - (Optional) Not supported for the time being, no need to fill in.
* `schedule_busy_workers` - (Optional) Whether to schedule a busy machine.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<namespace>:<group_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the App Group.
* `delete` - (Defaults to 5 mins) Used when delete the App Group.
* `update` - (Defaults to 5 mins) Used when update the App Group.

## Import

Schedulerx App Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_schedulerx_app_group.example <namespace>:<group_id>
```