---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_ai_model_providers"
sidebar_current: "docs-alicloud-datasource-apig-ai-model-providers"
description: |-
  Provides a list of Apig Ai Model Provider owned by an Alibaba Cloud account.
---

# alicloud_apig_ai_model_providers

This data source provides Apig Ai Model Provider available to the user.[What is Ai Model Provider](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateAiModelProvider)

-> **NOTE:** Available since v1.286.0.

## Example Usage

```terraform
variable "gateway_id" {
  description = "The ID of an existing AI-type APIG gateway"
  type        = string
}

data "alicloud_apig_ai_model_providers" "default" {
  gateway_id = var.gateway_id
}

output "first_provider_id" {
  value = data.alicloud_apig_ai_model_providers.default.providers[0].model_provider_id
}
```

## Argument Reference

The following arguments are supported:
* `gateway_id` - (Required, ForceNew) The ID of the AI gateway instance. The target instance must exist, belong to the current account, and be of the AI gateway type.
* `ids` - (Optional, Computed) A list of Ai Model Provider IDs. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ai Model Provider IDs.
* `providers` - A list of Ai Model Provider Entries. Each element contains the following attributes:
    * `display_name` - Model supplier presentation name.
    * `gateway_id` - The ID of the AI gateway instance.
    * `model_count` - The number of model cards currently associated with the model supplier.
    * `model_provider` - Stable model vendor ID, no more than 128 characters in length.
    * `model_provider_id` - The first ID of the resource.
    * `source` - Model supplier source.
    * `update_time` - The last update time of the model vendor, in the format of yyyy-MM-dd HH:mm:ss.
    * `id` - The ID of the resource supplied above.
    * `bound_services` - A list of AI service summaries currently bound to this model vendor. Each element contains the following attributes:
        * `service_id` - The ID of the AI service.
        * `name` - The name of the AI service.
        * `namespace` - The namespace of the AI service.
        * `source_type` - The source type of the AI service.
        * `group_name` - The group name of the AI service.
        * `qualifier` - The qualifier of the AI service.
        * `express_type` - The express type of the AI service.
        * `pai_workspace_id` - The PAI workspace ID.
        * `pai_workspace_name` - The PAI workspace name.
        * `status` - The status of the AI service.
    * `model_cards` - A list of model cards currently associated with the model supplier. Each element contains the following attributes:
        * `model_card_id` - The ID of the model card.
        * `gateway_id` - The gateway ID of the model card.
        * `model_provider` - The model provider identifier.
        * `model_name` - The model name.
        * `source` - The model source.
        * `update_time` - The last update time of the model card.
