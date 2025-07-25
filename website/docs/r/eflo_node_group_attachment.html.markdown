---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_node_group_attachment"
description: |-
  Provides a Alicloud Eflo Node Group Attachment resource.
---

# alicloud_eflo_node_group_attachment

Provides a Eflo Node Group Attachment resource.

Node Association Node Group Resources.

For information about Eflo Node Group Attachment and how to use it, see [What is Node Group Attachment](https://next.api.alibabacloud.com/document/eflo-controller/2022-12-15/ExtendCluster).

-> **NOTE:** Available since v1.255.0.

## Example Usage

Basic Usage

```terraform
# Before executing this example, you need to confirm with the product team whether the resources are sufficient or you will get an error message with "Failure to check order before create instance"
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_eflo_node_group_attachment" "default" {
  vswitch_id     = "vsw-xxx"
  hostname       = "attachment-example"
  login_password = "G7f$2kL9@vQx3Zp5*"
  cluster_id     = "i118976xxxx"
  node_group_id  = "i127582xxxx"
  node_id        = "e01-cn-xxxx"
  vpc_id         = "vpc-xxx"
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Optional, ForceNew, Computed) Cluster ID
* `data_disk` - (Optional, List) The data disk of the cloud disk to be attached to the node. See [`data_disk`](#data_disk) below.
* `hostname` - (Required, ForceNew) Node hostname
* `login_password` - (Required) Node login password
* `node_group_id` - (Optional, ForceNew, Computed) Node group ID
* `node_id` - (Optional, ForceNew, Computed) Node ID
* `user_data` - (Optional) User-defined data
* `vswitch_id` - (Required, ForceNew) vswitch id
* `vpc_id` - (Required, ForceNew) Vpc id

### `data_disk`

The data_disk supports the following:
* `category` - (Optional) Type
* `delete_with_node` - (Optional) Indicate whether the data disk is released with the node. true indicates that the data disk will be released together when the node unsubscribes.
* `performance_level` - (Optional) Performance level
* `size` - (Optional, Int) Data disk size

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<node_group_id>:<node_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 3605 mins) Used when create the Node Group Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Node Group Attachment.

## Import

Eflo Node Group Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_node_group_attachment.example <cluster_id>:<node_group_id>:<node_id>
```