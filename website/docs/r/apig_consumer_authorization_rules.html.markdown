---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_consumer_authorization_rules"
description: |-
  Provides a batch of Alicloud APIG Consumer Authorization Rules.
---

# alicloud_apig_consumer_authorization_rules

Provides an exact batch of long-term APIG Consumer authorization rules for HTTP API routes.

For information about APIG Consumer authorization, see [CreateConsumerAuthorizationRules](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateConsumerAuthorizationRules).

-> **NOTE:** Available since v1.291.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_consumer_authorization_rules&exampleId=56ee8d30-6420-bf75-de29-617ff4fce0497278c98f&activeTab=example&spm=docs.r.apig_consumer_authorization_rules.0.56ee8d3064&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_apig_consumer_authorization_rules" "default" {
  consumer_id        = alicloud_apig_consumer.default.id
  environment_id     = alicloud_apig_gateway.default.environments[0].environment_id
  parent_resource_id = alicloud_apig_http_api.default.id
  resource_ids       = [alicloud_apig_route.default.route_id]
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_apig_consumer_authorization_rules&spm=docs.r.apig_consumer_authorization_rules.example&intl_lang=EN_US)


## Argument Reference

* `consumer_id` - (Required, ForceNew) The consumer ID.
* `environment_id` - (Required, ForceNew) The environment ID.
* `expire_mode` - (Optional, ForceNew) Expiration mode. The supported value is `LongTerm`, which is the default.
* `parent_resource_id` - (Required, ForceNew) The parent HTTP API ID.
* `principal_type` - (Optional, ForceNew) Principal type. The supported value is `Consumer`, which is the default.
* `resource_ids` - (Required, ForceNew, Set) The exact route ID set. All rules are created in one batch request.
* `resource_type` - (Optional, ForceNew) Resource type. The supported value is `HttpApiRoute`, which is the default.

## Attributes Reference

* `id` - A composite ID identifying the consumer, environment, parent resource, and resource type.
* `authorization_rule_ids` - A map from route ID to consumer authorization rule ID. Response ordering is not used.

## Timeouts

* `create` - (Defaults to 5 mins) Used when creating and reading back the rule batch.
* `delete` - (Defaults to 5 mins) Used when deleting the rule batch.

## Import

Import resolves rule IDs through `QueryConsumerAuthorizationRules`; individual rule IDs are not required. Supply the exact route set as a comma-separated final component:

```shell
$ terraform import alicloud_apig_consumer_authorization_rules.example '<consumer_id>:<environment_id>:<http_api_id>:HttpApiRoute:<route_id_1>,<route_id_2>'
```
