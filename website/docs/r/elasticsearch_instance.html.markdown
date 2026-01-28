---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_instance"
description: |-
  Provides a Alicloud Elasticsearch Instance resource.
---

# alicloud_elasticsearch_instance

Provides a Elasticsearch Instance resource.



For information about Elasticsearch Instance and how to use it, see [What is Instance](https://next.api.alibabacloud.com/document/elasticsearch/2017-06-13/createInstance).

-> **NOTE:** Available since v1.30.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_elasticsearch_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_elasticsearch_zones.default.zones.0.id
}

resource "alicloud_elasticsearch_instance" "default" {
  description                      = var.name
  vswitch_id                       = alicloud_vswitch.default.id
  password                         = "Examplw1234"
  version                          = "7.10_with_X-Pack"
  instance_charge_type             = "PostPaid"
  data_node_amount                 = "2"
  data_node_spec                   = "elasticsearch.sn2ne.large"
  data_node_disk_size              = "20"
  data_node_disk_type              = "cloud_ssd"
  kibana_node_spec                 = "elasticsearch.sn2ne.large"
  data_node_disk_performance_level = "PL1"
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

## Argument Reference

The following arguments are supported:
* `vswitch_id` - (Required, ForceNew) The ID of VSwitch.
* `kms_encrypted_password` - (Optional, Available since v1.57.1) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored, but you have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available since v1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `period` - (Optional) The duration that you will buy Elasticsearch instance (in month). It is valid when PaymentType is `Subscription`. Valid values: [1~9], 12, 24, 36. Default to 1. From version 1.69.2, when to modify this value, the resource can renewal a `PrePaid` instance.
* `auto_renew_duration` - (Optional, Int) Renewal Period
* `client_node_configuration` - (Optional, Computed, Set, Available since v1.267.0) Elasticsearch cluster coordination node configuration See [`client_node_configuration`](#client_node_configuration) below.
* `data_node_configuration` - (Optional, Computed, Set, Available since v1.267.0) Elasticsearch data node information See [`data_node_configuration`](#data_node_configuration) below.
* `description` - (Optional, Computed) Instance name
* `enable_kibana_private_network` - (Optional, Computed, Available since v1.87.0) Whether to enable Kibana private network access.

The meaning of the value is as follows:
  - true: On.
  - false: does not open.
* `enable_kibana_public_network` - (Optional, Computed, Available since v1.87.0) Does Kibana enable public access
* `enable_public` - (Optional, Computed, Available since v1.87.0) Whether to enable Kibana public network access.

The meaning of the value is as follows:
  - true: On.
  - false: does not open.
* `force` - (Optional, Available since v1.267.0) Whether to force changes

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `instance_category` - (Optional, ForceNew, Computed, Available since v1.267.0) Version type.
* `kibana_configuration` - (Optional, Computed, Set, Available since v1.267.0) Elasticsearch Kibana node settings See [`kibana_configuration`](#kibana_configuration) below.
* `kibana_private_security_group_id` - (Optional) Kibana private network security group ID
* `kibana_private_whitelist` - (Optional, Computed, List, Available since v1.87.0) Cluster Kibana node private network access whitelist
* `kibana_whitelist` - (Optional, Computed, List) Kibana private network access whitelist
* `master_configuration` - (Optional, Computed, Set, Available since v1.267.0) Elasticsearch proprietary master node configuration information See [`master_configuration`](#master_configuration) below.
* `order_action_type` - (Optional, Available since v1.267.0) The instance changes the operation type. UPGRADE, UPGRADE. DOWNGRADE, DOWNGRADE.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `password` - (Optional) The access password of the instance.
* `payment_type` - (Optional, Computed, Available since v1.267.0) The payment method of the instance. Optional values: `prepaid` (subscription) and `postpaid` (pay-as-you-go)
* `private_whitelist` - (Optional, Computed, List) Elasticsearch private network whitelist. (Same as EsIpWhitelist)
* `protocol` - (Optional, Computed, Available since v1.101.0) Access protocol. Optional values: `HTTP` and **HTTPS * *.
* `public_whitelist` - (Optional, Computed, List) Elasticseach public network access whitelist IP list
* `renew_status` - (Optional, Computed) Renewal Status
* `renewal_duration_unit` - (Optional, Computed) Renewal Period Unit
* `resource_group_id` - (Optional, ForceNew, Computed, Available since v1.86.0) Resource group to which the instance belongs
* `setting_config` - (Optional, Computed, Map, Available since v1.125.0) Configuration information
* `tags` - (Optional, Computed, Map, Available since v1.73.0) Collection of tag key-value pairs
* `update_strategy` - (Optional, Available since v1.267.0) The change policy for Elasticsearch.

The values are as follows:
  - blue_green: blue-green change, which can realize seamless switching by running two identical environments (blue environment and green environment) in parallel.
  - normal: In-place changes, changes are made directly in the current environment (for example, upgrades, scaling) without additional resources.
  - intelligent: intelligent change, the system automatically analyzes the change type and environmental status, and dynamically selects the optimal change method (that is, blue-green change or in-situ change).

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `version` - (Required, ForceNew) Instance version
* `warm_node_configuration` - (Optional, Computed, Set, Available since v1.267.0) Elasticsearch cluster cold data node configuration See [`warm_node_configuration`](#warm_node_configuration) below.
* `zone_count` - (Optional, ForceNew, Computed, Int) The number of zones in the Elasticsearch instance.

The following arguments will be discarded. Please use new fields as soon as possible:

* `instance_charge_type` - (Optional, Deprecated since v1.261.0) Valid values are `PrePaid`, `PostPaid`. Default to `PostPaid`. From version 1.69.0, the Elasticsearch cluster allows you to update your instance_charge_ype from `PostPaid` to `PrePaid`, the following attributes are required: `period`.
* `client_node_spec` - (Optional, Deprecated since v1.261.0) The client node spec. If specified, client node will be created.
* `client_node_amount` - (Optional, Deprecated since v1.261.0) The Elasticsearch cluster's client node quantity, between 2 and 25.
* `master_node_spec` - (Optional, Deprecated since v1.261.0) The dedicated master node spec. If specified, dedicated master node will be created.
* `master_node_disk_type` - (Optional, Deprecated since v1.261.0) The single master node storage space. Valid values are `PrePaid`, `PostPaid`.
* `warm_node_spec` - (Optional, Deprecated since v.1.261.0) The warm node specifications of the Elasticsearch instance.
* `warm_node_amount` - (Optional, Deprecated since v1.261.0) The Elasticsearch cluster's warm node quantity, between 3 and 50.
* `warm_node_disk_encrypted` - (Optional, Deprecated since v1.261.0) If encrypt the warm node disk. Valid values are `true`, `false`. Default to `false`.
* `warm_node_disk_size` - (Optional, Deprecated since v.1.261.0) The single warm node storage space, should between 500 and 20480
* `warm_node_disk_type` - (Optional, Deprecated since v.1.261.0) The warm node disk type. Supported values:  cloud_efficiency.
* `kibana_node_spec` - (Optional, Deprecated since v1.261.0) The kibana node specifications of the Elasticsearch instance. Default is `elasticsearch.n4.small`.
* `data_node_spec` - (Optional, Deprecated since v1.261.0) The data node specifications of the Elasticsearch instance.
* `data_node_amount` - (Optional, Deprecated since v1.261.0) The Elasticsearch cluster's data node quantity, between 2 and 50.
* `data_node_disk_performance_level` - (Optional, Deprecated since v1.261.0) Cloud disk performance level. Valid values are `PL0`, `PL1`, `PL2`, `PL3`. The `data_node_disk_type` muse be `cloud_essd`.
* `data_node_disk_encrypted` - (Optional, Deprecated since v1.261.0) If encrypt the data node disk. Valid values are `true`, `false`. Default to `false`.
* `data_node_disk_type` - (Optional, ForceNew, Deprecated since v1.261.0) The data node disk type. Supported values: cloud_ssd, cloud_efficiency.
* `data_node_disk_size` - (Optional, Deprecated since v1.261.0) The single data node storage space.
  - `cloud_ssd`: An SSD disk, supports a maximum of 2048 GiB (2 TB).
  - `cloud_efficiency` An ultra disk, supports a maximum of 5120 GiB (5 TB). If the data to be stored is larger than 2048 GiB, an ultra disk can only support the following data sizes (GiB): [`2560`, `3072`, `3584`, `4096`, `4608`, `5120`].

### `client_node_configuration`

The client_node_configuration supports the following:
* `amount` - (Optional, Int, Available since v1.267.0) Number of disks in the Elasticsearch cluster coordination node
* `disk` - (Optional, ForceNew, Computed, Int, Available since v1.267.0) Elasticsearch cluster coordinates node disk size
* `disk_type` - (Optional, ForceNew, Available since v1.267.0) Elasticsearch cluster coordination node disk type
* `spec` - (Optional, Available since v1.267.0) Elasticsearch cluster coordination node specification

### `data_node_configuration`

The data_node_configuration supports the following:
* `amount` - (Optional, Computed, Int, Available since v1.267.0) Number of data nodes in the Elasticsearch cluster
* `disk` - (Optional, Int, Available since v1.267.0) Elasticsearch data node disk size
* `disk_encryption` - (Optional, ForceNew, Computed, Available since v1.267.0) Whether the Elasticsearch data node disk is encrypted
* `disk_type` - (Optional, ForceNew, Computed, Available since v1.267.0) Elasticsearch cluster data node disk type
* `performance_level` - (Optional, Computed, Available since v1.267.0) Elasticsearch cluster data node Essd disk level
* `spec` - (Required, Available since v1.267.0) Elasticsearch data node specification

### `kibana_configuration`

The kibana_configuration supports the following:
* `amount` - (Optional, ForceNew, Computed, Int, Available since v1.267.0) The number of disks of the Elasticsearch Kibana node. The default value is 1.
* `disk` - (Optional, ForceNew, Computed, Int, Available since v1.267.0) Elasticsearch Kibana node disk size
* `spec` - (Required, Available since v1.267.0) Elasticsearch Kibana node disk specifications

### `master_configuration`

The master_configuration supports the following:
* `amount` - (Optional, ForceNew, Int, Available since v1.267.0) Elasticsearch proprietary master node number of disks
* `disk` - (Optional, ForceNew, Int, Available since v1.267.0) Elasticsearch proprietary master node disk size
* `disk_type` - (Optional, ForceNew, Available since v1.267.0) Elasticsearch proprietary master node disk type
* `spec` - (Optional, Available since v1.267.0) Elasticsearch proprietary master node specifications

### `warm_node_configuration`

The warm_node_configuration supports the following:
* `amount` - (Optional, Int, Available since v1.267.0) Elasticsearch cluster cold data node disk number
* `disk` - (Optional, Int, Available since v1.267.0) Elasticsearch cluster cold data node disk size
* `disk_encryption` - (Optional, ForceNew, Available since v1.267.0) Elasticsearch cluster cold data node Disk encryption
* `disk_type` - (Optional, ForceNew, Available since v1.267.0) Elasticsearch cluster cold data node disk type
* `spec` - (Optional, Available since v1.267.0) Elasticsearch cluster cold data node Disk Specification

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `arch_type` - Schema Type:.
* `create_time` - Instance creation time.
* `domain` - Elasticsearch cluster private domain name.
* `kibana_domain` - Kibana address.
* `kibana_port` - The port assigned by the Kibana node.
* `public_domain` - The public network address of the current instance.
* `public_port` - Elasticsearch cluster public network access port
* `port` - Instance connection port.
* `status` - Instance change status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 61 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 360 mins) Used when update the Instance.

## Import

Elasticsearch Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_elasticsearch_instance.example <instance_id>
```