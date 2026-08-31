---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_policy"
description: |-
  Provides an Alicloud APIG Policy and attachment resource.
---

# alicloud_apig_policy

Provides an APIG Policy together with its attachment.

For information about APIG Policies, see [CreateAndAttachPolicy](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateAndAttachPolicy).

-> **NOTE:** Available since v1.291.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_policy&exampleId=2fad0ea1-fd2f-eaf3-f556-18bcd6c49ab58a7e3390&activeTab=example&spm=docs.r.apig_policy.0.2fad0ea1fd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_apig_policy" "default" {
  policy_name = "terraform-example"
  class_name  = "RateLimit"
  config = jsonencode({
    enable             = true
    threshold          = 10
    responseStatusCode = 429
  })
  description          = "Example route policy"
  attach_resource_ids  = [alicloud_apig_route.default.route_id]
  attach_resource_type = "GatewayRoute"
  environment_id       = alicloud_apig_gateway.default.environments[0].environment_id
  gateway_id           = alicloud_apig_gateway.default.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_apig_policy&spm=docs.r.apig_policy.example&intl_lang=EN_US)


## Argument Reference

* `attach_resource_ids` - (Required, Set) IDs of the resources to which the policy is attached.
* `attach_resource_type` - (Optional, ForceNew) Attachment resource type. Defaults to `GatewayRoute`.
* `class_name` - (Required, ForceNew) The APIG policy class name.
* `config` - (Required, Sensitive) JSON policy configuration. Equivalent JSON formatting is ignored.
* `description` - (Optional) The policy description.
* `environment_id` - (Required) The environment ID for the attachment.
* `gateway_id` - (Required) The gateway ID for the attachment.
* `policy_name` - (Required) The policy name.

## Attributes Reference

* `id` - The policy ID.
* `policy_attachment_id` - The policy attachment ID.

## Timeouts

* `create` - (Defaults to 5 mins) Used when creating and attaching the policy.
* `delete` - (Defaults to 5 mins) Used when deleting the attachment and policy.
* `update` - (Defaults to 5 mins) Used when updating the policy and attachment.

## Import

Import requires both the policy and attachment IDs:

```shell
$ terraform import alicloud_apig_policy.example '<policy_id>:<policy_attachment_id>'
```
