---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_rate_limit_policy"
sidebar_current: "docs-alicloud-resource-polardb-gateway-rate-limit-policy"
description: |-
  Provides a PolarDB AI Gateway Rate Limit Policy resource.
---

# alicloud_polardb_gateway_rate_limit_policy

Provides a PolarDB AI Gateway Rate Limit Policy resource.

-> **NOTE:** Available since v2.1.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_gateway_rate_limit_policy&exampleId=a09ed95a-a7fb-b5af-f0af-6d59a74243d9408ee08a&activeTab=example&spm=docs.r.polardb_gateway_rate_limit_policy.0.a09ed95aa7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_gateway_rate_limit_policy" "example" {
  gateway_id     = alicloud_polardb_gateway.example.id
  scope_type     = "Consumer"
  scope_ref_id   = alicloud_polardb_gateway_consumer.example.consumer_id
  rate_limit_rpm = "100"
  rate_limit_tpm = "10000"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_gateway_rate_limit_policy&spm=docs.r.polardb_gateway_rate_limit_policy.example&intl_lang=EN_US)


## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `scope_type` - (Required, ForceNew, Available since v2.1.0) The rate limit dimension. Valid values: `ConsumerGroup`, `Consumer`.
* `scope_ref_id` - (Required, ForceNew, Available since v2.1.0) The consumer group or consumer ID.
* `rate_limit_rpm` - (Required, Available since v2.1.0) The maximum number of requests per minute.
* `rate_limit_tpm` - (Required, Available since v2.1.0) The maximum number of tokens per minute.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<policy_id>`.
* `policy_id` - (Available since v2.1.0) The rate limit policy ID.
* `policy_type` - (Available since v2.1.0) The policy type.
* `status` - (Available since v2.1.0) The policy status.
* `create_time` - (Available since v2.1.0) The creation time.
* `modify_time` - (Available since v2.1.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_rate_limit_policy.example pg-abc:policy-abc
```
