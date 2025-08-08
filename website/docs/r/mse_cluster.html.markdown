---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_cluster"
sidebar_current: "docs-alicloud-resource-mse-cluster"
description: |-
  Provides an Alicloud MSE Cluster resource.
---

# alicloud_mse_cluster

Provides a MSE Cluster resource. It is a one-stop microservice platform for the industry's mainstream open source microservice frameworks Spring Cloud and Dubbo, providing governance center, managed registry and managed configuration center.

-> **NOTE:** Available since v1.94.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mse_cluster&exampleId=5d68da6c-7d9a-38e1-a002-d2e37218505fa3c72b8a&activeTab=example&spm=docs.r.mse_cluster.0.5d68da6c7d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create resource
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_2_0_0"
  version_code          = "NACOS_2_3_2_1"
  instance_count        = 3
  net_type              = "privatenet"
  pub_network_flow      = "1"
  connection_type       = "slb"
  cluster_alias_name    = "terraform-example"
  mse_version           = "mse_pro"
  vswitch_id            = alicloud_vswitch.example.id
  vpc_id                = alicloud_vpc.example.id
}
```

### Deleting `alicloud_mse_cluster` or removing it from your configuration

The `alicloud_mse_cluster` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.


## Argument Reference

The following arguments are supported:

* `acl_entry_list` - (Optional) The whitelist. **NOTE:** This attribute is invalid when the value of `pub_network_flow` is `0` and the value of `net_type` is `privatenet`.
* `cluster_alias_name` - (Optional, Computed) The alias of MSE Cluster.
* `cluster_specification` - (Required) The engine specification of MSE Cluster. **NOTE:** From version 1.188.0, `cluster_specification` can be modified. If you were an international user, please use the specification version ending with `_200_c`.Valid values:
  - Professional Edition
    - `MSE_SC_1_2_60_c`: 1C2G
    - `MSE_SC_2_4_60_c`: 2C4G
    - `MSE_SC_4_8_60_c`: 4C8G
    - `MSE_SC_8_16_60_c`: 8C16G
    - `MSE_SC_16_32_60_c`:16C32G
    - `MSE_SC_1_2_200_c`: 1C2G
    - `MSE_SC_2_4_200_c`: 2C4G
    - `MSE_SC_4_8_200_c`: 4C8G
    - `MSE_SC_8_16_200_c`: 8C16G
    - `MSE_SC_16_32_200_c`:16C32G
  - Developer Edition
    - `MSE_SC_1_2_60_c`: 1C2G
    - `MSE_SC_2_4_60_c`: 2C4G
    - `MSE_SC_1_2_200_c`: 1C2G
    - `MSE_SC_2_4_200_c`: 2C4G
  - Serverless Edition
    - `MSE_SC_SERVERLESS`: Available since v1.232.0
* `cluster_type` - (Required, ForceNew) The type of MSE Cluster.
* `cluster_version` - (Required, ForceNew) The version of MSE Cluster. See [details](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createcluster)
* `version_code` - (Optional) The version code of MSE Cluster. You can keep the instance version up to date by setting the value to `LATEST` (Available since v1.257.0).
* `disk_type` - (Optional) The type of Disk.
* `instance_count` - (Required) The count of instance. **NOTE:** From version 1.188.0, `instance_count` can be modified.
* `net_type` - (Required, ForceNew) The type of network. Valid values: `privatenet` and `pubnet` and `both`(Available since v1.232.0).
* `payment_type` - (Optional, ForceNew, Computed, Available since v1.220.0) Payment type: Subscription (prepaid), PayAsYouGo (postpaid). Default PayAsYouGo.
* `tags` - (Optional, Map, Computed, Available since v1.220.0) The tag of the resource.
* `resource_group_id` - (Optional, Computed, Available since v1.220.0) The resource group of the resource.
* `private_slb_specification` - (Optional) The specification of private network SLB.
* `pub_network_flow` - (Required from 1.173.0) The public network bandwidth.
* `pub_slb_specification` - (Optional) The specification of public network SLB. Serverless Instance could ignore this parameter.
* `vswitch_id` - (Optional) The id of VSwitch.
* `mse_version` - (Optional, ForceNew, Computed, Available since v1.177.0) The version of MSE. Valid values: `mse_dev` or `mse_pro` or `mse_serverless`(Available since v1.232.0).
* `connection_type` - (Optional, ForceNew, Available since v1.183.0) The connection type. Valid values: `slb`,`single_eni`(Available since v1.232.0). If your region is one of `ap-southeast-6、us-west-1、eu-central-1、us-east-1、ap-southeast-1`,and your cluster's mse_version is `mse_dev`,please use `single_eni`.
* `request_pars` - (Optional, Available since v1.183.0) The extended request parameters in the JSON format.
* `vpc_id` - (Optional, Available since v1.185.0) The id of the VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the resource.The instance id of cluster.
* `cluster_id` - (Available since v1.162.0) The cluster id of Cluster.
* `app_version` - (Available since v1.205.0) The application version.
* `status` - The status of MSE Cluster.

## Timeouts

-> **NOTE:** Available since v1.188.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the MSE Cluster.
* `update` - (Defaults to 15 mins) Used when update the MSE Cluster.
* `delete` - (Defaults to 15 mins) Used when delete the MSE Cluster.

## Import

MSE Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_mse_cluster.example mse-cn-0d9xxxx
```
