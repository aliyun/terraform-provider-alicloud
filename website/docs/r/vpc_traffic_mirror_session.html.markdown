---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_traffic_mirror_session"
sidebar_current: "docs-alicloud-resource-vpc-traffic-mirror-session"
description: |-
  Provides a Alicloud VPC Traffic Mirror Session resource.
---

# alicloud_vpc_traffic_mirror_session

Provides a VPC Traffic Mirror Session resource. Traffic mirroring session.

For information about VPC Traffic Mirror Session and how to use it, see [What is Traffic Mirror Session](https://www.alibabacloud.com/help/en/doc-detail/261364.htm).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_traffic_mirror_session&exampleId=3e93056a-2e0d-89da-24c4-74e9269890ad9a069d8f&activeTab=example&spm=docs.r.vpc_traffic_mirror_session.0.3e93056a2e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g7"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = alicloud_vpc.default.id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
  count                = 2
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_name        = var.name
  host_name            = var.name
  image_id             = data.alicloud_images.default.images.0.id
  instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  security_groups      = [alicloud_security_group.default.id]
  vswitch_id           = alicloud_vswitch.default.id
  system_disk_category = "cloud_essd"
}

resource "alicloud_ecs_network_interface" "default" {
  count                  = 2
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  count                = 2
  instance_id          = alicloud_instance.default[count.index].id
  network_interface_id = alicloud_ecs_network_interface.default[count.index].id
}

resource "alicloud_vpc_traffic_mirror_filter" "default" {
  traffic_mirror_filter_name        = var.name
  traffic_mirror_filter_description = var.name
}


resource "alicloud_vpc_traffic_mirror_session" "default" {
  priority                           = 1
  virtual_network_id                 = 10
  traffic_mirror_session_description = var.name
  traffic_mirror_session_name        = var.name
  traffic_mirror_target_id           = alicloud_ecs_network_interface_attachment.default[0].network_interface_id
  traffic_mirror_source_ids          = [alicloud_ecs_network_interface_attachment.default[1].network_interface_id]
  traffic_mirror_filter_id           = alicloud_vpc_traffic_mirror_filter.default.id
  traffic_mirror_target_type         = "NetworkInterface"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_traffic_mirror_session&spm=docs.r.vpc_traffic_mirror_session.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Whether to PreCheck only this request, value:
  - **true**: sends a check request and does not create a mirror session. Check items include whether required parameters are filled in, request format, and restrictions. If the check fails, the corresponding error is returned. If the check passes, the error code 'DryRunOperation' is returned '.
  - **false** (default): Sends a normal request and directly creates a mirror session after checking.
* `enabled` - (Optional) Specifies whether to enable traffic mirror sessions. default to `false`.
* `packet_length` - (Optional, ForceNew, Available since v1.206.0) Maximum Transmission Unit (MTU).
* `priority` - (Required) The priority of the traffic mirror session. Valid values: `1` to `32766`. A smaller value indicates a higher priority. You cannot specify the same priority for traffic mirror sessions that are created in the same region with the same Alibaba Cloud account.
* `resource_group_id` - (Optional, Available since v1.206.0) The ID of the resource group.
* `tags` - (Optional, Map, Available since v1.206.0) The tags of this resource.
* `traffic_mirror_filter_id` - (Required) The ID of the filter.
* `traffic_mirror_session_description` - (Optional) The description of the traffic mirror session. The description must be `2` to `256` characters in length and cannot start with `http://` or `https://`.
* `traffic_mirror_session_name` - (Optional)  The name of the traffic mirror session. The name must be `2` to `128` characters in length and can contain digits, underscores (_), and hyphens (-). It must start with a letter.
* `traffic_mirror_source_ids` - (Required) The ID of the image source instance. Currently, the Eni is supported as the image source. The default value of N is 1, that is, only one mirror source can be added to a mirror session.
* `traffic_mirror_target_id` - (Required) The ID of the mirror destination. You can specify only an ENI or a Server Load Balancer (SLB) instance as a mirror destination.
* `traffic_mirror_target_type` - (Required) The type of the mirror destination. Valid values: `NetworkInterface` or `SLB`. `NetworkInterface`: an ENI. `SLB`: an internal-facing SLB instance.
* `virtual_network_id` - (Optional) The VXLAN network identifier (VNI) that is used to distinguish different mirrored traffic. Valid values: `0` to `16777215`. You can specify VNIs for the traffic mirror destination to identify mirrored traffic from different sessions. If you do not specify a VNI, the system randomly allocates a VNI. If you want the system to randomly allocate a VNI, ignore this parameter.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Mirror Session.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Mirror Session.
* `update` - (Defaults to 5 mins) Used when update the Traffic Mirror Session.

## Import

VPC Traffic Mirror Session can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_traffic_mirror_session.example <id>
```