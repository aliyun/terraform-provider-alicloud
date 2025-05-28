---
subcategory: "Lindorm"
layout: "alicloud"
page_title: "Alicloud: alicloud_lindorm_public_network"
description: |-
  Provides a Alicloud Lindorm Public Network resource.
---

# alicloud_lindorm_public_network

Provides a Lindorm Public Network resource.

Public network connection of Lindorm instance.

For information about Lindorm Public Network and how to use it, see [What is Public Network](https://next.api.alibabacloud.com/document/hitsdb/2020-06-15/SwitchInstancePublicNetwork).

-> **NOTE:** Available since v1.250.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

variable "zone_id" {
  default = "cn-shanghai-f"
}

variable "region_id" {
  default = "cn-shanghai"
}

resource "alicloud_vpc" "defaultX7MgJO" {
  description = var.name
  cidr_block  = "10.0.0.0/8"
  vpc_name    = "amp-example-shanghai"
}

resource "alicloud_vswitch" "default45mCzM" {
  description = var.name
  vpc_id      = alicloud_vpc.defaultX7MgJO.id
  zone_id     = var.zone_id
  cidr_block  = "10.0.0.0/24"
}

resource "alicloud_lindorm_instance" "defaultQpsLKr" {
  payment_type               = "PayAsYouGo"
  table_engine_node_count    = "2"
  instance_storage           = "80"
  zone_id                    = var.zone_id
  vswitch_id                 = alicloud_vswitch.default45mCzM.id
  disk_category              = "cloud_efficiency"
  table_engine_specification = "lindorm.g.xlarge"
  instance_name              = "tf-example"
  vpc_id                     = alicloud_vpc.defaultX7MgJO.id
}


resource "alicloud_lindorm_public_network" "default" {
  instance_id           = alicloud_lindorm_instance.defaultQpsLKr.id
  enable_public_network = "1"
  engine_type           = "lindorm"
}
```

### Deleting `alicloud_lindorm_public_network` or removing it from your configuration

Terraform cannot destroy resource `alicloud_lindorm_public_network`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `enable_public_network` - (Optional, Int) Open or close the public connection. Value:
  - `0`: Closes the public network connection.
  - `1`: Enable the public network connection.
* `engine_type` - (Required) Engine type, value:
  - `lindorm`: Wide table engine.
  - `tsdb`: Time series engine.
  - `solr`: Search engine.
  - `blob`:S3 compatible.
* `instance_id` - (Required, ForceNew) Instance ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Instance status, returns:_EXPANDING`: Capacity-based cloud storage is being expanded._version_transing`: The minor version is being upgraded._CHANGING`: The specification is being upgraded or downgraded._SWITCHING`:SSL is being changed._OPENING`: The data subscription function is being activated._TRANSFER`: migrates data to the database._CREATING`: in the production disaster recovery instance._RECOVERING`: The backup is being restored._IMPORTING`: Data is being imported._MODIFYING`: The network is being changed._SWITCHING`: The internal network and the external network are being switched._CREATING`: creates a network link._DELETING`: deletes a network link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 11 mins) Used when create the Public Network.

## Import

Lindorm Public Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_lindorm_public_network.example <id>
```