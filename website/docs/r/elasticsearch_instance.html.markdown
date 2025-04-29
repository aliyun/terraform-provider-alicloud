---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_instance"
sidebar_current: "docs-alicloud-resource-elasticsearch-instance"
description: |-
  Provides a Alicloud Elasticsearch instance resource.
---

# alicloud_elasticsearch_instance

Provides an Elasticsearch instance resource. It contains data nodes, dedicated master node(optional) and etc. It can be associated with private IP whitelists and kibana IP whitelist. see [What is Elasticsearch Instance](https://www.alibabacloud.com/help/en/es/developer-reference/api-createinstance).

-> **NOTE:** Only one operation is supported in a request. So if `data_node_spec` and `data_node_disk_size` are both changed, system will respond error.

-> **NOTE:** At present, `version` can not be modified once instance has been created.

-> **NOTE:** Available since v1.30.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_elasticsearch_instance&exampleId=73805ca0-8968-cf9a-4e8e-205d4afdd97b26ab1e4d&activeTab=example&spm=docs.r.elasticsearch_instance.0.73805ca089&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
### Deleting `alicloud_elasticsearch_instance` or removing it from your configuration

The `alicloud_elasticsearch_instance` resource allows you to manage `instance_charge_type = "Prepaid"` Elasticsearch instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Elasticsearch Instance.
You can resume managing the subscription Elasticsearch instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `description` - (Optional, Computed) The description of instance. It a string of 0 to 30 characters.
* `instance_charge_type` - (Optional) Valid values are `PrePaid`, `PostPaid`. Default to `PostPaid`. From version 1.69.0, the Elasticsearch cluster allows you to update your instance_charge_ype from `PostPaid` to `PrePaid`, the following attributes are required: `period`. But, updating from `PostPaid` to `PrePaid` is not supported.
* `period` - (Optional) The duration that you will buy Elasticsearch instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1. From version 1.69.2, when to modify this value, the resource can renewal a `PrePaid` instance.
* `data_node_amount` - (Required) The Elasticsearch cluster's data node quantity, between 2 and 50.
* `data_node_spec` - (Required) The data node specifications of the Elasticsearch instance.
* `data_node_disk_size` - (Required) The single data node storage space.
  - `cloud_ssd`: An SSD disk, supports a maximum of 2048 GiB (2 TB).
  - `cloud_efficiency` An ultra disk, supports a maximum of 5120 GiB (5 TB). If the data to be stored is larger than 2048 GiB, an ultra disk can only support the following data sizes (GiB): [`2560`, `3072`, `3584`, `4096`, `4608`, `5120`].
* `data_node_disk_type` - (Required, ForceNew) The data node disk type. Supported values: cloud_ssd, cloud_efficiency.
* `data_node_disk_encrypted` - (Optional, ForceNew, Available since 1.86.0) If encrypt the data node disk. Valid values are `true`, `false`. Default to `false`.
* `data_node_disk_performance_level` - (Optional, Available since 1.208.1) Cloud disk performance level. Valid values are `PL0`, `PL1`, `PL2`, `PL3`. The `data_node_disk_type` muse be `cloud_essd`.
* `vswitch_id` - (Required, ForceNew) The ID of VSwitch.
* `password` - (Optional, Sensitive) The password of the instance. The password can be 8 to 30 characters in length and must contain three of the following conditions: uppercase letters, lowercase letters, numbers, and special characters (`!@#$%^&*()_+-=`).
* `kms_encrypted_password` - (Optional, Available since 1.57.1) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored, but you have to specify one of `password` and `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, MapString, Available since 1.57.1) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `version` - (Required, ForceNew) Elasticsearch version. Supported values: `5.5.3_with_X-Pack`, `6.3_with_X-Pack`, `6.7_with_X-Pack`, `6.8_with_X-Pack`, `7.4_with_X-Pack` , `7.7_with_X-Pack`, `7.10_with_X-Pack`, `7.16_with_X-Pack`, `8.5_with_X-Pack`, `8.9_with_X-Pack`, `8.13_with_X-Pack`.
* `private_whitelist` - (Optional) Set the instance's IP whitelist in VPC network.
* `public_whitelist` - (Optional) Set the instance's IP whitelist in internet network.
* `enable_public` - (Optional, Available since v1.87.0) Bool, default to false. When it set to true, the instance can enable public network access。
* `kibana_whitelist` - (Optional) Set the Kibana's IP whitelist in internet network.
* `enable_kibana_public_network` - (Optional, Available since v1.87.0) Bool, default to true. When it set to false, the instance can enable kibana public network access。
* `kibana_private_whitelist` - (Optional, Available since v1.87.0) Set the Kibana's IP whitelist in private network, This option has been abandoned on newly created instance, please use `kibana_private_security_group_id` instead
* `enable_kibana_private_network` - (Optional, Available since v1.87.0) Bool, default to false. When it set to true, the instance can close kibana private network access。
* `master_node_spec` - (Optional) The dedicated master node spec. If specified, dedicated master node will be created.
* `master_node_disk_type` - (Optional, Available since 1.208.1) The single master node storage space. Valid values are `PrePaid`, `PostPaid`.
* `client_node_amount` - (Optional, Available since v1.101.0) The Elasticsearch cluster's client node quantity, between 2 and 25.
* `client_node_spec` - (Optional, Available since v1.101.0) The client node spec. If specified, client node will be created.
* `kibana_node_spec` - (Optional, Available since v1.163.0) The kibana node specifications of the Elasticsearch instance. Default is `elasticsearch.n4.small`.
* `protocol` - (Optional, Available since v1.101.0) Elasticsearch protocol. Supported values: `HTTP`, `HTTPS`.default is `HTTP`.
* `zone_count` - (Optional, ForceNew, Available since 1.44.0) The Multi-AZ supported for Elasticsearch, between 1 and 3. The `data_node_amount` value must be an integral multiple of the `zone_count` value.
* `tags` - (Optional, Available since v1.73.0) A mapping of tags to assign to the resource. 
  - `key`: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:". It cannot contain "http://" and "https://". It cannot be a null string.
  - `value`: It can be up to 128 characters in length. It cannot contain "http://" and "https://". It can be a null string.
* `resource_group_id` - (Optional, ForceNew, Computed, Available since v1.86.0) The ID of resource group which the Elasticsearch instance belongs.
* `setting_config` - (Optional, Computed, Available since v1.125.0) The YML configuration of the instance.[Detailed introduction](https://www.alibabacloud.com/help/doc-detail/61336.html).
* `renew_status` - (Optional, Available since 1.208.1) The renewal status of the specified instance. Valid values: `AutoRenewal`, `ManualRenewal`, `NotRenewal`.The `instance_charge_type` must be `PrePaid`.
* `auto_renew_duration` - (Optional, Available since 1.208.1) Auto-renewal period of an Elasticsearch Instance, in the unit of the month. It is valid when `instance_charge_type` is `PrePaid` and `renew_status` is `AutoRenewal`.
* `renewal_duration_unit` - (Optional, Available since 1.208.1) Auto-Renewal Cycle Unit Values Include: Month: Month. Year: Years. Valid values: `M`, `Y`.
* `domain` - (Computed, Available since 1.197.0) Instance connection domain (only VPC network access supported).
* `port` - (Computed, Available since 1.197.0) Instance connection port.
* `public_domain` - (Computed, Available since 1.197.0) Instance connection public domain.
* `public_port` - (Computed, Available since 1.197.0) Instance connection public port.
* `kibana_domain` - (Computed, Available since 1.197.0) Kibana console domain (Internet access supported).
* `kibana_port` - (Computed, Available since 1.197.0) Kibana console port.
* `status` - (Computed, Available since 1.197.0) The Elasticsearch instance status. Includes `active`, `activating`, `inactive`. Some operations are denied when status is not `active`.
* `warm_node_spec` - (Optional, Available since v.1.229.0) The warm node specifications of the Elasticsearch instance.
* `warm_node_amount` - (Optional, Available since v.1.229.0) The Elasticsearch cluster's warm node quantity, between 3 and 50.
* `warm_node_disk_size` - (Optional, Available since v.1.229.0) The single warm node storage space, should between 500 and 20480
* `warm_node_disk_type` - (Optional, Available since v.1.229.0) The warm node disk type. Supported values:  cloud_efficiency.
* `warm_node_disk_encrypted` - (Optional, ForceNew, Available since v.1.229.0) If encrypt the warm node disk. Valid values are `true`, `false`. Default to `false`.
* `kibana_private_security_group_id` - (Optional, Available since v.1.229.0) the security group id associated with Kibana private network, this param is required when `enable_kibana_private_network` set true, and the security group id should in the same VPC as `vswitch_id`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Elasticsearch instance.
* `domain` - Instance connection domain (only VPC network access supported).
* `port` - Instance connection port.
* `public_domain` - (Available since 1.197.0) Instance connection public domain.
* `public_port` - (Available since 1.197.0) Instance connection public port.
* `kibana_domain` - Kibana console domain (Internet access supported).
* `kibana_port` - Kibana console port.
* `status` - The Elasticsearch instance status. Includes `active`, `activating`, `inactive`. Some operations are denied when status is not `active`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 120 mins) Used when creating the elasticsearch instance (until it reaches the initial `active` status).
* `update` - (Defaults to 120 mins) Used when activating the elasticsearch instance when necessary during update - e.g. when changing elasticsearch instance description, whitelist, data node settings, master node spec and password.
* `delete` - (Defaults to 120 mins) Used when terminating the elasticsearch instance. `Note`: There are 5 minutes to sleep to eusure the instance is deleted. It is not in the timeouts configure.

## Import

Elasticsearch can be imported using the id, e.g.

```shell
$ terraform import alicloud_elasticsearch_instance.example es-cn-abcde123456
```
