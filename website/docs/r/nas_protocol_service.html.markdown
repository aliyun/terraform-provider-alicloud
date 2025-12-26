---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_protocol_service"
description: |-
  Provides a Alicloud File Storage (NAS) Protocol Service resource.
---

# alicloud_nas_protocol_service

Provides a File Storage (NAS) Protocol Service resource.



For information about File Storage (NAS) Protocol Service and how to use it, see [What is Protocol Service](https://next.api.alibabacloud.com/document/NAS/2017-06-26/CreateProtocolService).

-> **NOTE:** Available since v1.267.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "example" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-examplee1031-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "example" {
  is_default   = false
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = "nas-examplee1031-vsw1sdw-F"
}

resource "alicloud_nas_file_system" "example" {
  description      = var.name
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  encrypt_type     = "0"
  vpc_id           = alicloud_vpc.example.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.example.id
  file_system_type = "cpfs"
}


resource "alicloud_nas_protocol_service" "default" {
  vpc_id         = alicloud_vpc.example.id
  protocol_type  = "NFS"
  protocol_spec  = "General"
  vswitch_id     = alicloud_vswitch.example.id
  dry_run        = false
  file_system_id = alicloud_nas_file_system.example.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description of the agreement service.

Limitations:
  - Length is 2~128 English or Chinese characters.
  - It must start with an uppercase or lowercase letter or Chinese, and cannot start with `http://` and `https://`.
  - Can contain numbers, colons (:), underscores (_), or dashes (-).
* `dry_run` - (Optional) Whether to PreCheck the creation request.

The pre-check operation helps you check the validity of parameters and dependency conditions, and does not actually create an instance, nor does it incur costs.

Value:
  - true: The check request is sent and the protocol service is not created. The check items include whether the required parameters, request format, and business restriction dependency conditions are filled in. If the check does not pass, the corresponding error is returned. If the check passes, the 200 HttpCode is returned, but the ProtocolServiceId is empty.
  - false (default): Send a normal request and directly create an instance after passing the check.

-> **NOTE:** This parameter only applies during resource creation, update or deletion. If modified in isolation without other property changes, Terraform will not trigger any action.

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `protocol_spec` - (Required, ForceNew) The specification of the protocol machine cluster.
  - Value range: General、CL1、CL2
  - Default value: General
* `protocol_throughput` - (Optional, ForceNew, Computed, Int) The throughput of the protocol service. Unit: MB/s.
* `protocol_type` - (Required, ForceNew) The protocol type supported by the protocol service.

Value range:
  - NFS: Protocol Service supports NFS protocol access.
* `vswitch_id` - (Optional, ForceNew) The VSwitchId of the protocol service.
* `vpc_id` - (Optional, ForceNew) The VpcId of the protocol service, which must be consistent with the VPC of the file system.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<protocol_service_id>`.
* `create_time` - The time when the protocol server service was created. The UTC time.
* `protocol_service_id` - Protocol Service ID
* `status` - Agreement service status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 20 mins) Used when create the Protocol Service.
* `delete` - (Defaults to 20 mins) Used when delete the Protocol Service.
* `update` - (Defaults to 10 mins) Used when update the Protocol Service.

## Import

File Storage (NAS) Protocol Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_protocol_service.example <file_system_id>:<protocol_service_id>
```