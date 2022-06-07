---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_cluster"
sidebar_current: "docs-alicloud-resource-mse-cluster"
description: |-
    Provides an Alicloud MSE Cluster resource.
---

# alicloud\_mse\_cluster

Provides a MSE Cluster resource. It is a one-stop microservice platform for the industry's mainstream open source microservice frameworks Spring Cloud and Dubbo, providing governance center, managed registry and managed configuration center.

-> **NOTE:** Available in 1.94.0+.

## Example Usage

```terraform
resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_ANS_1_2_1"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = "vsw-123456"
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = "tf-testAccMseCluster"
}
```

## Argument Reference

The following arguments are supported:

* `acl_entry_list` - (Optional) The whitelist. **NOTE:** This attribute is invalid when the value of `pub_network_flow` is `0` and the value of `net_type` is `privatenet`.
* `cluster_alias_name` - (Optional) The alias of MSE Cluster.
* `cluster_specification` - (Required, ForceNew) The engine specification of MSE Cluster. Valid values:
  `MSE_SC_1_2_200_c`：1C2G
  `MSE_SC_2_4_200_c`：2C4G
  `MSE_SC_4_8_200_c`：4C8G
  `MSE_SC_8_16_200_c`：8C16G
  
* `cluster_type` - (Required, ForceNew) The type of MSE Cluster.
* `cluster_version` - (Required, ForceNew) The version of MSE Cluster.
* `disk_type` - (Optional, ForceNew) The type of Disk.
* `instance_count` - (Required, ForceNew) The count of instance.
* `net_type` - (Required, ForceNew) The type of network. Valid values: "privatenet" and "pubnet".
* `private_slb_specification` - (Optional, ForceNew) The specification of private network SLB.
* `pub_network_flow` - (Required from 1.173.0, ForceNew) The public network bandwidth. `0` means no access to the public network.
* `pub_slb_specification` - (Optional, ForceNew) The specification of public network SLB.
* `vswitch_id` - (Optional, ForceNew) The id of VSwitch.
* `mse_version` - (Optional, ForceNew, Computed, Available in v1.177.0+) The version of MSE. Valid values: `mse_basic` or `mse_pro`.
                    
## Attributes Reference

The following attributes are exported:

* `id` - The id of the resource.
* `status` - The status of MSE Cluster.
* `cluster_id` - (Available in v1.162.0+)  The id of Cluster.

## Import

MSE Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_mse_cluster.example mse-cn-0d9xxxx
```
