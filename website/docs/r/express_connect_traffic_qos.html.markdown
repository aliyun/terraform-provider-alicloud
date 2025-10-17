---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_traffic_qos"
description: |-
  Provides a Alicloud Express Connect Traffic Qos resource.
---

# alicloud_express_connect_traffic_qos

Provides a Express Connect Traffic Qos resource.

Express Connect Traffic QoS Policy.

For information about Express Connect Traffic Qos and how to use it, see [What is Traffic Qos](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateExpressConnectTrafficQos).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_traffic_qos&exampleId=cdfe0cdd-5031-c21b-d79f-a22265481bb05c38f7c1&activeTab=example&spm=docs.r.express_connect_traffic_qos.0.cdfe0cdd50&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_express_connect_traffic_qos" "createQos" {
  qos_name        = var.name
  qos_description = var.name
}
```

## Argument Reference

The following arguments are supported:
* `qos_description` - (Optional) The description of the QoS policy.
  The length is `0` to `256` characters and cannot start with 'http:// 'or 'https.
* `qos_name` - (Optional) The name of the QoS policy.
  The length is `0` to `128` characters and cannot start with 'http:// 'or 'https.
* `resource_group_id` - (Optional, Available since v1.261.0) The ID of the resource group.
* `tags` - (Optional, Map, Available since v1.261.0) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the QoS policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Qos.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Qos.
* `update` - (Defaults to 8 mins) Used when update the Traffic Qos.

## Import

Express Connect Traffic Qos can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_traffic_qos.example <id>
```
