---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_consumer"
sidebar_current: "docs-alicloud-resource-polardb-gateway-consumer"
description: |-
  Provides a PolarDB AI Gateway Consumer resource.
---

# alicloud_polardb_gateway_consumer

Provides a PolarDB AI Gateway Consumer resource.

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** `api_key` is sensitive but is stored in Terraform state. Use an encrypted remote state backend and restrict access to the state.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_gateway_consumer&exampleId=8c8056ee-63c2-8dde-451d-ce28f4ac4f7263dc9df4&activeTab=example&spm=docs.r.polardb_gateway_consumer.0.8c8056ee63&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_gateway_consumer" "example" {
  gateway_id        = alicloud_polardb_gateway.example.id
  name              = "application-a"
  consumer_group_id = alicloud_polardb_gateway_consumer_group.example.consumer_group_id
  nickname          = "Application A"
  key_type          = "ApiKey"
  is_default        = "0"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_gateway_consumer&spm=docs.r.polardb_gateway_consumer.example&intl_lang=EN_US)


## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `name` - (Required, ForceNew, Available since v2.1.0) The name of the consumer.
* `consumer_group_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the consumer group to which the consumer belongs.
* `nickname` - (Optional, ForceNew, Available since v2.1.0) The nickname of the consumer.
* `key_type` - (Optional, ForceNew, Available since v2.1.0) The key type. Currently, only `ApiKey` is supported. Default: `ApiKey`.
* `is_default` - (Optional, ForceNew, Available since v2.1.0) Whether the consumer belongs to the default group. Valid values: `0`, `1`. Default: `0`.
* `api_key_reset_token` - (Optional, Available since v2.1.0) An arbitrary value that, when changed to a non-empty value, resets the API key. Store this trigger separately from the generated key.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<consumer_id>`.
* `consumer_id` - (Available since v2.1.0) The consumer ID returned by the service API.
* `consumer_group_name` - (Available since v2.1.0) The consumer group name.
* `api_key` - (Sensitive, Available since v2.1.0) The full API key returned only during creation or reset.
* `allowed_models` - (Available since v2.1.0) The models that the consumer is allowed to access.
* `month_to_date_cost_count` - (Available since v2.1.0) The month-to-date cost count.
* `lifetime_cost_count` - (Available since v2.1.0) The lifetime cost count.
* `month_to_date_token_count` - (Available since v2.1.0) The month-to-date token count.
* `lifetime_token_count` - (Available since v2.1.0) The lifetime token count.
* `create_time` - (Available since v2.1.0) The creation time.
* `modify_time` - (Available since v2.1.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_consumer.example pg-abc:c-abc
```

-> **NOTE:** The API does not return an existing consumer's full API key, key type, or default-consumer flag. After import, change `api_key_reset_token` to generate a new key if required. The imported state cannot reconstruct the original `key_type` or `is_default` values.
