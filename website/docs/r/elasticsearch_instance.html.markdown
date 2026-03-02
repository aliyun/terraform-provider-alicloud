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
* `auto_renew_duration` - (Optional, Int) Number of auto-renewal periods.  
* `client_node_configuration` - (Optional, Computed, Set, Available since v1.267.0) Configuration of dedicated coordinating nodes in the Elasticsearch cluster.   See [`client_node_configuration`](#client_node_configuration) below.
* `data_node_configuration` - (Optional, Computed, Set, Available since v1.267.0) Elasticsearch data node information. See [`data_node_configuration`](#data_node_configuration) below.
* `description` - (Optional, Computed) Instance name, which supports fuzzy search. For example, searching for all instances containing `abc` may return instances named `abc`, `abcde`, `xyabc`, or `xabcy`.
* `enable_kibana_private_network` - (Optional, Computed, Available since v1.87.0) Indicates whether private network access to Kibana is enabled. Valid values:  
  - true: Enabled  
  - false: Disabled  
* `enable_kibana_public_network` - (Optional, Computed, Available since v1.87.0) Specifies whether to enable public access to Kibana. Valid values:  
  - true: Enables public access.  
  - false: Disables public access.  
* `enable_public` - (Optional, Computed, Available since v1.87.0) Specifies whether to enable a public endpoint for the instance. Valid values:
  - true: Enables the public endpoint.
  - false: Disables the public endpoint.
* `force` - (Optional, Available since v1.267.0) Whether to force a restart:
  - true: Yes  
  - false (default): No.

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `instance_category` - (Optional, ForceNew, Computed, Available since v1.267.0) Edition type:  
  - x-pack: Creates a commercial edition instance, or a kernel-enhanced edition instance without Indexing Service or OpenStore enabled.  
  - IS: Creates a kernel-enhanced edition instance with Indexing Service or OpenStore enabled.  
* `kibana_configuration` - (Optional, Computed, Set, Available since v1.267.0) The configuration of Elasticsearch Kibana nodes. See [`kibana_configuration`](#kibana_configuration) below.
* `kibana_private_security_group_id` - (Optional) List of security groups.
* `kibana_private_whitelist` - (Optional, Computed, List, Available since v1.87.0) List of IP addresses in the whitelist. This parameter is available when whiteIpGroup is empty and is used to modify the default group's whitelist.  
* `kibana_whitelist` - (Optional, Computed, List) The list of IP addresses in the whitelist. This parameter is available when whiteIpGroup is empty and modifies the default group's whitelist.
* `master_configuration` - (Optional, Computed, Set, Available since v1.267.0) Configuration information for Elasticsearch dedicated master nodes. See [`master_configuration`](#master_configuration) below.
* `order_action_type` - (Optional, Available since v1.267.0) Configuration change type. Valid values:
  - upgrade (default): Upgrade configuration
  - downgrade: Downgrade configuration.

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `password` - (Optional) The access password for the instance. It must be 8 to 32 characters in length and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters (!@#$%^&*()_+-=).  
* `payment_type` - (Optional, Computed, Available since v1.267.0) The billing method of the instance. Supported values:
  - `prepaid`: Subscription
  - `postpaid`: Pay-as-you-go
* `private_whitelist` - (Optional, Computed, List) The list of IP addresses in the whitelist. This parameter is available when whiteIpGroup is empty and modifies the default group's whitelist.
* `protocol` - (Optional, Computed, Available since v1.101.0) The access protocol. Supported protocols: HTTP and HTTPS.  
* `public_whitelist` - (Optional, Computed, List) The IP address whitelist. This parameter is available when whiteIpGroup is empty and is used to modify the default group's whitelist.
* `renew_status` - (Optional, Computed) The renewal status. Valid values:
  - AutoRenewal: Auto-renewal.
  - ManualRenewal: Manual renewal.
  - NotRenewal: No renewal.
* `renewal_duration_unit` - (Optional, Computed) The unit of the auto-renewal period. Valid values:  
  - M: Month.  
  - Y: Year.  

-> **NOTE:**  This parameter is required when RenewalStatus is set to AutoRenewal.  

* `resource_group_id` - (Optional, ForceNew, Computed, Available since v1.86.0) The ID of the resource group to which the instance belongs.
* `setting_config` - (Optional, Computed, Map, Available since v1.125.0) YML configuration file settings for the instance.
* `tags` - (Optional, Computed, Map, Available since v1.73.0) Instance tag group.
* `update_strategy` - (Optional, Available since v1.267.0) Elasticsearch update strategy (for example, index updates, cluster upgrades, or service deployments). Valid values:
  - blue_green: Blue-green deployment, which enables seamless switching by running two identical environments (blue and green) in parallel.
  - normal: In-place update, which applies changes directly in the current environment (for example, upgrades or scaling) without requiring additional resources.
  - intelligent: Intelligent update, where the system automatically analyzes the update type and environment status to dynamically select the optimal strategy (either blue-green or in-place).

-> **NOTE:** This parameter only takes effect when other resource properties are also modified. Changing this parameter alone will not trigger a resource update.

* `version` - (Required, ForceNew) The instance version. Valid values:
  - 8.5.1_with_X-Pack
  - 7.10_with_X-Pack
  - 6.7_with_X-Pack
  - 7.7_with_X-Pack
  - 6.8_with_X-Pack
  - 6.3_with_X-Pack
  - 5.6_with_X-Pack
  - 5.5.3_with_X-Pack

-> **NOTE:**  The versions listed above might not include all versions supported by Elasticsearch instances. You can call the [GetRegionConfiguration](https://help.aliyun.com/document_detail/254099.html) operation to view the actual supported versions.

* `warm_node_configuration` - (Optional, Computed, Set, Available since v1.267.0) Cold data node configuration for the Elasticsearch cluster. See [`warm_node_configuration`](#warm_node_configuration) below.
* `zone_count` - (Optional, ForceNew, Computed, Int) The number of zones for the instance. Valid values: 1, 2, and 3. Default value: 1.  

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
* `amount` - (Optional, Int, Available since v1.267.0) Number of nodes.  
* `disk` - (Optional, ForceNew, Computed, Int, Available since v1.267.0) Node storage capacity, in GB.
* `disk_type` - (Optional, ForceNew, Available since v1.267.0) Storage type of the node. Only ultra disk (cloud_efficiency) is supported.  
* `spec` - (Optional) Node specification. You can view specification details in [Product Specifications](https://help.aliyun.com/document_detail/271718.html).

### `data_node_configuration`

The data_node_configuration supports the following:
* `amount` - (Optional, Computed, Int, Available since v1.267.0) Number of data nodes. Valid values: 2 to 50.
* `disk` - (Optional, Int, Available since v1.267.0) Storage capacity per node, in GB.
* `disk_encryption` - (Optional, ForceNew, Computed, Available since v1.267.0) Whether to enable cloud disk encryption:
  - true: Enabled
  - false: Disabled.
* `disk_type` - (Optional, ForceNew, Computed, Available since v1.267.0) Node disk type. Supported types:
  - cloud_ssd: SSD cloud disk
  - cloud_efficiency: Ultra cloud disk.
* `performance_level` - (Optional, Computed, Available since v1.267.0) Performance level of ESSD cloud disks. This parameter is required when diskType is set to cloud_essd. Supported values: PL1, PL2, PL3.
* `spec` - (Required, Available since v1.267.0) Node specification. For more information about specifications, see [Product Specifications](https://help.aliyun.com/document_detail/271718.html).

### `kibana_configuration`

The kibana_configuration supports the following:
* `amount` - (Optional, ForceNew, Computed, Int, Available since v1.267.0) The number of nodes.
* `disk` - (Optional, ForceNew, Computed, Int, Available since v1.267.0) Storage capacity per node, in GB.
* `spec` - (Required, Available since v1.267.0) Node specification. For specification details, see [Product Specifications](https://help.aliyun.com/document_detail/271718.html).

### `master_configuration`

The master_configuration supports the following:
* `amount` - (Optional, ForceNew, Int, Available since v1.267.0) Number of nodes.
* `disk` - (Optional, ForceNew, Int, Available since v1.267.0) Node storage capacity, in GB.
* `disk_type` - (Optional, ForceNew, Available since v1.267.0) Node storage type. Only cloud_ssd (SSD cloud disk) is supported.
* `spec` - (Optional, Available since v1.267.0) Node specification. For specifications, see [Product Specifications](https://help.aliyun.com/document_detail/271718.html).

### `warm_node_configuration`

The warm_node_configuration supports the following:
* `amount` - (Optional, Int, Available since v1.267.0) Number of nodes.
* `disk` - (Optional, Int, Available since v1.267.0) Storage capacity per node, in GB.
* `disk_encryption` - (Optional, ForceNew, Available since v1.267.0) Whether to enable disk encryption. The values are as follows:
  - true: Enabled.
  - false: Disabled.
* `disk_type` - (Optional, ForceNew, Available since v1.267.0) Storage type for the node. Only `cloud_efficiency` (ultra disk) is supported.
* `spec` - (Optional, Available since v1.267.0) Node specification. For specifications, see [Product Specifications](https://help.aliyun.com/document_detail/271718.html).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `arch_type` - The deployment mode or architecture type:.
* `create_time` - The time when the instance was created.
* `domain` - The internal network address of the instance.
* `kibana_domain` - Kibana endpoint.
* `kibana_port` - The access port for Kibana.
* `kibana_private_domain` - The private endpoint of Kibana.
* `public_domain` - The public endpoint of the instance.
* `public_port` - The public access port of the instance.
* `port` - Instance connection port.
* `status` - The status of the instance.

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