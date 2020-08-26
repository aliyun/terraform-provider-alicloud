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

```
resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type = "Eureka"
  cluster_version = "EUREKA_1_9_3"
  instance_count = 1
  net_type = "privatenet"
  vswitch_id = "vsw-123456"
  pub_network_flow = "1"
  acl_entry_list= ["127.0.0.1/32"]
  cluster_alias_name= "tf-testAccMseCluster"
}
```

## Argument Reference

The following arguments are supported:

* `acl_entry_list` - (Optional) The whitelist.
* `cluster_alias_name` - (Optional) The alias of MSE Cluster.
* `cluster_specification` - (Required, ForceNew) The engine specification of MSE Cluster. Valid values: `MSE_SC_1_2_200_c`, `MSE_SC_2`, `MSE_SC_4_8_200_c_4_200_c`, `MSE_SC_8_16_200_c`.
* `cluster_type` - (Required, ForceNew) The type of MSE Cluster.
* `cluster_version` - (Required, ForceNew) The version of MSE Cluster.
* `disk_type` - (Optional, ForceNew) The type of Disk.
* `instance_count` - (Optional, ForceNew) The count of instance.
* `net_type` - (Required, ForceNew) The type of network. Range limit: 1~5.
* `private_slb_specification` - (Optional, ForceNew) The specification of private network SLB.
* `pub_network_flow` - (Optional, ForceNew) The public network bandwidth. `0` means no access to the public network.
* `pub_slb_specification` - (Optional, ForceNew) The specification of public network SLB.
* `vswitch_id` - (Optional, ForceNew) The id of VSwitch.
                    
## Attributes Reference

The following attributes are exported:

* `id` - The id of MSE Cluster.
* `status` - The status of MSE Cluster.

## Import

MSE Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_mse_cluster.example mse-cn-0d9xxxx
```
