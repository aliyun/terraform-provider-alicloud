---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_cost_rule"
sidebar_current: "docs-alicloud-resource-polardb-gateway-cost-rule"
description: |-
  Provides a PolarDB AI Gateway Cost Rule resource.
---

# alicloud_polardb_gateway_cost_rule

Provides a PolarDB AI Gateway Cost Rule resource.

-> **NOTE:** Available since v2.1.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_gateway_cost_rule&exampleId=9a27e15a-33a6-61be-b6c8-ba915f017d245a6157d8&activeTab=example&spm=docs.r.polardb_gateway_cost_rule.0.9a27e15a33&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_gateway_cost_rule" "example" {
  gateway_id                     = alicloud_polardb_gateway.example.id
  model_name                     = "qwen-plus"
  model_service_id               = alicloud_polardb_gateway_model_service.example.model_service_id
  input_cost_points_per_million  = "10"
  output_cost_points_per_million = "20"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_gateway_cost_rule&spm=docs.r.polardb_gateway_cost_rule.example&intl_lang=EN_US)


## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `model_name` - (Required, ForceNew, Available since v2.1.0) The model name.
* `model_service_id` - (Required, ForceNew, Available since v2.1.0) The model service ID returned by the service API, without the gateway portion of the Terraform composite ID.
* `input_cost_points_per_million` - (Optional, Available since v2.1.0) Input-token cost points. Default: `0`.
* `output_cost_points_per_million` - (Optional, Available since v2.1.0) Output-token cost points. Default: `0`.
* `cache_cost_points_per_million` - (Optional, Available since v2.1.0) Cached-token cost points. Default: `0`.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<cost_rule_id>`.
* `cost_rule_id` - (Available since v2.1.0) The cost rule ID returned by the service API.
* `create_time` - (Available since v2.1.0) The creation time.
* `modify_time` - (Available since v2.1.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_cost_rule.example pg-abc:cost-rule-abc
```
