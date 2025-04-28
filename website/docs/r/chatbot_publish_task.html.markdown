---
subcategory: "Chatbot"
layout: "alicloud"
page_title: "Alicloud: alicloud_chatbot_publish_task"
sidebar_current: "docs-alicloud-resource-chatbot-publish-task"
description: |-
  Provides a Alicloud Chatbot Publish Task resource.
---

# alicloud_chatbot_publish_task

Provides a Chatbot Publish Task resource.

For information about Chatbot Publish Task and how to use it, see [What is Publish Task](https://help.aliyun.com/document_detail/433996.html).

-> **NOTE:** Available since v1.203.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_chatbot_publish_task&exampleId=8bca0b7b-d5e1-b85b-748d-2ab7a5291c42c7b34ef6&activeTab=example&spm=docs.r.chatbot_publish_task.0.8bca0b7bd5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}
data "alicloud_chatbot_agents" "default" {}
resource "alicloud_chatbot_publish_task" "default" {
  biz_type  = "faq"
  agent_key = data.alicloud_chatbot_agents.default.agents.0.agent_key
}
```

## Argument Reference

The following arguments are supported:
* `agent_key` - (Optional) The business space key. If you do not set it, the default business space is accessed. The key value is obtained on the business management page of the primary account.
* `biz_type` - (Required, ForceNew) The type of the publishing unit. Please use the CreateInstancePublishTask API to publish the robot.
* `data_id_list` - (Optional) Additional release information. Currently supported: If the BizType is faq, enter the category Id in this field to indicate that only the knowledge under these categories is published.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - UTC time of task creation
* `modify_time` - UTC time for task modification
* `status` - The status of the task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Publish Task.

## Import

Chatbot Publish Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_chatbot_publish_task.example <id>
```