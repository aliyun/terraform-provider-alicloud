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

-> **NOTE:** Available in 1.94.0+.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
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
  instance_count        = 1
  net_type              = "privatenet"
  pub_network_flow      = "1"
  connection_type       = "slb"
  cluster_alias_name    = "terraform-example"
  mse_version           = "mse_dev"
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
* `cluster_alias_name` - (Optional) The alias of MSE Cluster.
* `cluster_specification` - (Required) The engine specification of MSE Cluster. **NOTE:** From version 1.188.0, `cluster_specification` can be modified. Valid values:
  - `MSE_SC_1_2_60_c`: 1C2G
  - `MSE_SC_2_4_60_c`: 2C4G
  - `MSE_SC_4_8_60_c`: 4C8G
  - `MSE_SC_8_16_60_c`: 8C16G
* `cluster_type` - (Required, ForceNew) The type of MSE Cluster.
* `cluster_version` - (Required, ForceNew) The version of MSE Cluster. See [details](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createcluster)
* `disk_type` - (Optional, ForceNew) The type of Disk.
* `instance_count` - (Required) The count of instance. **NOTE:** From version 1.188.0, `instance_count` can be modified.
* `net_type` - (Required, ForceNew) The type of network. Valid values: "privatenet" and "pubnet".
* `payment_type` - (Optional, ForceNew, Computed, Available since v1.220.0) Payment type: Subscription (prepaid), PayAsYouGo (postpaid). Default PayAsYouGo.
* `tags` - (Optional, Map, Computed, Available since v1.220.0) The tag of the resource.
* `resource_group_id` - (Optional, Computed, Available since v1.220.0) The resource group of the resource.
* `private_slb_specification` - (Optional, ForceNew) The specification of private network SLB.
* `pub_network_flow` - (Required from 1.173.0, ForceNew) The public network bandwidth. `0` means no access to the public network.
* `pub_slb_specification` - (Optional, ForceNew) The specification of public network SLB.
* `vswitch_id` - (Optional, ForceNew) The id of VSwitch.
* `mse_version` - (Optional, ForceNew, Computed, Available in v1.177.0+) The version of MSE. Valid values: `mse_dev` or `mse_pro`.
* `connection_type` - (Optional, ForceNew, Available in v1.183.0+) The connection type. Valid values: `slb`.
* `request_pars` - (Optional, Available in v1.183.0+) The extended request parameters in the JSON format.
* `vpc_id` - (Optional, ForceNew, Available in v1.185.0+) The id of the VPC.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the resource.
* `cluster_id` - (Available in v1.162.0+) The id of Cluster.
* `app_version` - (Available in v1.205.0+) The application version.
* `status` - The status of MSE Cluster.

## Timeouts

-> **NOTE:** Available in 1.188.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the MSE Cluster.
* `update` - (Defaults to 15 mins) Used when update the MSE Cluster.
* `delete` - (Defaults to 15 mins) Used when delete the MSE Cluster.

## Import

MSE Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_mse_cluster.example mse-cn-0d9xxxx
```
