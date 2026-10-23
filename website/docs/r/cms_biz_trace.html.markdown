---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_biz_trace"
description: |-
  Provides a Alicloud Cms Biz Trace resource.
---

# alicloud_cms_biz_trace

Provides a Cms Biz Trace resource.

business trace analysis Service-related configuration.

For information about Cms Biz Trace and how to use it, see [What is Biz Trace](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateBizTrace).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cms_biz_trace" "default" {
  biz_trace_code  = "vip_order"
  biz_trace_name  = "vip下单"
  rule_config     = "[{\"entrancePid\":\"ggxw4lnjuz@fef9c08e66b5a4c\",\"characteristics\":{\"rules\":[{\"id\":\"l3wwnk6o\",\"matcher\":{\"matchType\":\"COLLECTION_CONTAINS\",\"pattern\":[\"GC5\",\"GC7\"]},\"key\":\"biz.vip.level\",\"target\":\"CUSTOM_EXTRACT\"}],\"operation\":\"AND\"},\"rpcMatcher\":{\"matchType\":\"EQUALS\",\"pattern\":\"/api/order\"}}]"
  advanced_config = "{\"bodyDetail\":[\"REQUEST\",\"RESPONSE\"],\"sample\":{\"strategy\":\"RATE\",\"detail\":\"100\"},\"dyingScope\":[\"TRACE\",\"METRIC\"]}"
  workspace       = "default"
}
```

## Argument Reference

The following arguments are supported:
* `advanced_config` - (Optional) Advanced Configuration
* `biz_trace_code` - (Required, ForceNew) The service link code, which is used for the value of the service link label.
* `biz_trace_name` - (Optional) The name of the resource
* `rule_config` - (Optional) Business Trace Rule Configuration
* `workspace` - (Optional) Work space

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Biz Trace.
* `delete` - (Defaults to 5 mins) Used when delete the Biz Trace.
* `update` - (Defaults to 5 mins) Used when update the Biz Trace.

## Import

Cms Biz Trace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_biz_trace.example <biz_trace_id>
```