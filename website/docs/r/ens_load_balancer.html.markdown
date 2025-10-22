---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_load_balancer"
description: |-
  Provides a Alicloud ENS Load Balancer resource.
---

# alicloud_ens_load_balancer

Provides a ENS Load Balancer resource.

Load balancing.

For information about ENS Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createloadbalancer).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_load_balancer&exampleId=cd617ca6-e9d9-748a-b19c-a142ef0dd8cb59a8a80d&activeTab=example&spm=docs.r.ens_load_balancer.0.cd617ca6e9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_network" "network" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = var.name
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}


resource "alicloud_ens_load_balancer" "default" {
  load_balancer_name = var.name
  payment_type       = "PayAsYouGo"
  ens_region_id      = "cn-chenzhou-telecom_unicom_cmcc"
  load_balancer_spec = "elb.s1.small"
  vswitch_id         = alicloud_ens_vswitch.switch.id
  network_id         = alicloud_ens_network.network.id
}
```

## Argument Reference

The following arguments are supported:
* `backend_servers` - (Optional, Set, Available since v1.230.0) The list of backend servers. See [`backend_servers`](#backend_servers) below.
* `ens_region_id` - (Required, ForceNew) The ID of the ENS node.
* `load_balancer_name` - (Optional) Name of the Server Load Balancer instance. The length is 1~80 English or Chinese characters. When this parameter is not specified, the system randomly assigns an instance name. Cannot start with http:// and https.
* `load_balancer_spec` - (Required, ForceNew) Specifications of the Server Load Balancer instance. Optional values: elb.s1.small,elb.s3.medium,elb.s2.small,elb.s2.medium,elb.s3.small.
* `network_id` - (Required, ForceNew) The network ID of the created edge load balancing (ELB) instance.
* `payment_type` - (Required, ForceNew) Server Load Balancer Instance Payment Type. Value:PayAsYouGo
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch to which the VPC instance belongs.

### `backend_servers`

The backend_servers supports the following:
* `ip` - (Optional, Available since v1.230.0) IP address of the backend server  Example value: 192.168.0.5.
* `port` - (Optional, Computed, Int, Available since v1.230.0) Port used by the backend server.
* `server_id` - (Required, Available since v1.230.0) Backend server instance ID  Example value: i-5vb5h5njxiuhn48a * * * *.
* `type` - (Optional, Available since v1.230.0) Backend server type  Example value: ens.
* `weight` - (Optional, Int, Available since v1.230.0) Weight of the backend server  Example value: 100.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation Time (UTC) of the load balancing instance.
* `status` - The status of the SLB instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

ENS Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_load_balancer.example <id>
```