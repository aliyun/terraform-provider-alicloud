---
subcategory: "Cloud Native API Gateway (APIG)"
layout: "alicloud"
page_title: "Alicloud: alicloud_apig_consumer"
description: |-
  Provides an Alicloud APIG Consumer resource.
---

# alicloud_apig_consumer

Provides an APIG Consumer resource with a service-generated AK/SK identity.

For information about APIG Consumers, see [CreateConsumer](https://next.api.alibabacloud.com/document/APIG/2024-03-27/CreateConsumer).

-> **NOTE:** Available since v1.291.0.

-> **NOTE:** The generated AK/SK is not returned by this resource and is never stored in Terraform state. Resource refresh uses `ListConsumers`, not `GetConsumer`, because `GetConsumer` includes credential fields in its response.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_apig_consumer&exampleId=b3a733f1-1640-cffe-1719-55524316ceb3d9182dcc&activeTab=example&spm=docs.r.apig_consumer.0.b3a733f116&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_apig_consumer" "default" {
  consumer_name = "terraform-example"
  description   = "Example APIG consumer"
  enable        = true
  gateway_type  = "API"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_apig_consumer&spm=docs.r.apig_consumer.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `consumer_name` - (Required, ForceNew) The consumer name.
* `credential_generate_mode` - (Optional, ForceNew) Credential generation mode. The supported value is `System`. Defaults to `System`.
* `description` - (Optional) The consumer description.
* `enable` - (Optional) Whether the consumer is enabled. Defaults to `true`.
* `gateway_type` - (Optional, ForceNew) Gateway type. Valid values are `API` and `AI`. Defaults to `API`.

## Attributes Reference

The following attributes are exported:

* `id` - The consumer ID.

## Timeouts

* `create` - (Defaults to 5 mins) Used when creating the consumer.
* `delete` - (Defaults to 5 mins) Used when deleting the consumer.
* `update` - (Defaults to 5 mins) Used when updating the consumer.

## Import

APIG Consumers can be imported using the consumer ID, e.g.

```shell
$ terraform import alicloud_apig_consumer.example <consumer_id>
```
