---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router resource.
---

# alicloud_cen_transit_router

Provides a Cloud Enterprise Network (CEN) Transit Router resource.



For information about Cloud Enterprise Network (CEN) Transit Router and how to use it, see [What is Transit Router](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateTransitRouter).

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router&exampleId=06f18f6e-9f4a-8bb8-b381-5abdc62dd60113b15574&activeTab=example&spm=docs.r.cen_transit_router.0.06f18f6e9f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alicloud_cen_transit_router" "example" {
  cen_id              = alicloud_cen_instance.example.id
  transit_router_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cen_transit_router&spm=docs.r.cen_transit_router.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `support_multicast` - (Optional, ForceNew, Bool, Available since v1.195.0) Specifies whether to enable the multicast feature for the Enterprise Edition transit router. Valid values:
  - `false` (default): no
  - `true`: yes
The multicast feature is supported only in specific regions. You can call [ListTransitRouterAvailableResource](https://www.alibabacloud.com/help/en/doc-detail/261356.html) to query the regions that support multicast.
* `tags` - (Optional, Map, Available since v1.193.1) The tag of the resource
* `transit_router_description` - (Optional) The description of the Enterprise Edition transit router instance.
The description must be 1 to 256 characters in length, and cannot start with http:// or https://. You can also leave this parameter empty.
* `transit_router_name` - (Optional) The name of the Enterprise Edition transit router.
The name must be 1 to 128 characters in length, and cannot start with http:// or https://. You can also leave this parameter empty.
* `dry_run` - (Optional, Bool) The dry run.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It formats as `<cen_id>:<transit_router_id>`.
* `create_time` - (Available since v1.247.0) The time when the transit router was created.
* `region_id` - (Available since v1.247.0) The ID of the region where the transit router is deployed.
* `status` - The status of the transit router.
* `transit_router_id` - The ID of the transit router.
* `type` - The edition of the transit router.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router.
* `update` - (Defaults to 5 mins) Used when update the Transit Router.

## Import

Cloud Enterprise Network (CEN) Transit Router can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router.example <cen_id>:<transit_router_id>
```
