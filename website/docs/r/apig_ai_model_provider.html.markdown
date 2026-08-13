---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_ai_model_provider"
description: |-
  Provides a Alicloud APIG Ai Model Provider resource.
---

# alicloud_apig_ai_model_provider

Provides a APIG Ai Model Provider resource.

AI Model Provider.

For information about APIG Ai Model Provider and how to use it, see [What is Ai Model Provider](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateAiModelProvider).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_ai_model_provider&exampleId=b6a4f2a3-0cc8-b2a6-6048-61b657bbc8610e95e3b8&activeTab=example&spm=docs.r.apig_ai_model_provider.0.b6a4f2a30c&intl_lang=EN_US" target="_blank">
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

# AiModelProvider must be bound to an existing AI-type APIG gateway.
# Set the gateway_id variable to your own AI gateway.
variable "gateway_id" {
  description = "The ID of an existing AI-type APIG gateway"
  type        = string
}

resource "alicloud_apig_ai_model_provider" "default" {
  gateway_id     = var.gateway_id
  model_provider = "openai"
  display_name   = var.name
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_apig_ai_model_provider&spm=docs.r.apig_ai_model_provider.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `display_name` - (Required) Model supplier presentation name. Required, no more than 128 characters in length.

* `gateway_id` - (Required, ForceNew) The ID of the AI gateway instance. The target instance must exist, belong to the current account, and be of the AI gateway type.

* `model_provider` - (Required, ForceNew) Stable model vendor ID, no more than 128 characters in length.

* `service_ids` - (Optional, List) The full collection of AI service IDs to bind to the model vendor. If not passed, existing bindings are retained; if empty array, all bindings are cleared.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `model_count` - The number of model cards currently associated with the model supplier.
* `source` - Model supplier source.
* `update_time` - The last update time of the model vendor, in the format of yyyy-MM-dd HH:mm:ss.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ai Model Provider.
* `delete` - (Defaults to 5 mins) Used when delete the Ai Model Provider.
* `update` - (Defaults to 5 mins) Used when update the Ai Model Provider.

## Import

APIG Ai Model Provider can be imported using the id, e.g.

```shell
$ terraform import alicloud_apig_ai_model_provider.example <model_provider_id>
```