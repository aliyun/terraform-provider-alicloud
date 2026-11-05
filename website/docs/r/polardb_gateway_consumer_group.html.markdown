---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_consumer_group"
sidebar_current: "docs-alicloud-resource-polardb-gateway-consumer-group"
description: |-
  Provides a PolarDB AI Gateway Consumer Group resource.
---

# alicloud_polardb_gateway_consumer_group

Provides a PolarDB AI Gateway Consumer Group resource.

-> **NOTE:** Available since v2.1.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_gateway_consumer_group&exampleId=d3a87da6-f93e-7248-bf1c-64dd5b571e2bcc141f35&activeTab=example&spm=docs.r.polardb_gateway_consumer_group.0.d3a87da6f9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_gateway_consumer_group" "example" {
  gateway_id = alicloud_polardb_gateway.example.id
  name       = "production"
  nickname   = "Production consumers"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_gateway_consumer_group&spm=docs.r.polardb_gateway_consumer_group.example&intl_lang=EN_US)


## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `name` - (Required, ForceNew, Available since v2.1.0) The name of the consumer group.
* `nickname` - (Optional, Available since v2.1.0) The nickname of the consumer group.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<consumer_group_id>`.
* `consumer_group_id` - (Available since v2.1.0) The consumer group ID returned by the service API.
* `is_default` - (Available since v2.1.0) Whether the consumer group is the default group.
* `allowed_models` - (Available since v2.1.0) The models that the consumer group is allowed to access.
* `create_time` - (Available since v2.1.0) The creation time.
* `modify_time` - (Available since v2.1.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_consumer_group.example pg-abc:cg-abc
```
