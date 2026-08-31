---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_db_clusters"
sidebar_current: "docs-alicloud-datasource-click-house-db-clusters"
description: |-
  Provides a list of Click House DBCluster to the user.
---

# alicloud\_click\_house\_db\_clusters

This data source provides the Click House DBCluster of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = "your_vswitch_id"
}

data "alicloud_click_house_db_clusters" "default" {
  ids = [alicloud_click_house_db_cluster.default.id]
}
output "db_cluster" {
  value = data.alicloud_click_house_db_clusters.default.ids.0
}

```

## Argument Reference

The following arguments are supported:

* `db_cluster_description` - (Optional, ForceNew) The DBCluster description.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of DBCluster IDs.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Running`,`Creating`,`Deleting`,`Restarting`,`Preparing`,.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `clusters` - A list of Click House DBClusters. Each element contains the following attributes:
    * `ali_uid` - Alibaba Cloud account Id.
    * `bid` - The ID of the business process flow.
    * `category` - Instance family values include: Basic: Basic edition; HighAvailability: high availability edition.
      * `commodity_code` - The Commodity Code of the DBCluster.
    * `connection_string` - Connection string.
    * `create_time` - The creation time of the resource.
    * `db_cluster_description` - The DBCluster description.
    * `id` - The DBCluster id.
    * `db_cluster_network_type` - The DBCluster network type.
    * `db_cluster_type` - The DBCluster type.
    * `db_node_class` - The node class of the DBCluster. 
    * `db_node_count` - The node count of the DBCluster.
    * `db_node_storage` - The node storage of the DBCluster.
    * `encryption_key` - Key management service KMS key ID.
    * `encryption_type` - Currently only supports ECS disk encryption, with a value of CloudDisk, not encrypted when empty.
    * `engine` - The Engine of the DBCluster.
    * `engine_version` - The engine version of the DBCluster.
    * `expire_time` - The expiration time of the DBCluster.
    * `id` - The ID of the DBCluster.
    * `is_expired` - If the instance has expired.
    * `lock_mode` - The lock mode of the DBCluster.
    * `lock_reason` - Lock reason of the DBCluster.
    * `maintain_time` - Examples of the maintenance window, in the format of hh:mmZ-hh:mm Z.
    * `payment_type` - The payment type of the resource. Valid values: `PayAsYouGo`,`Subscription`.
    * `port` - Connection port.
    * `public_connection_string` - A public IP address for the connection.
    * `public_port` - Public network port.
    * `scale_out_status` - Scale state.
        * `progress` - Process.
        * `ratio` - Efficiency.
    * `storage_type` - Storage type of DBCluster. Valid values: `cloud_essd`, `cloud_efficiency`, `cloud_essd_pl2`, `cloud_essd_pl3`.
    * `support_backup` - Support fallback scheme.
    * `support_https_port` - The system supports http port number.
    * `support_mysql_port` - Supports Mysql, and those of the ports.
    * `vpc_cloud_instance_id` - Virtual Private Cloud (VPC cloud instance ID.
    * `vpc_id` - The VPC ID of the DBCluster.
    * `vswitch_id` - The vswitch id of the DBCluster.
    * `zone_id` - The zone ID of the DBCluster.
    * `control_version` - The control version of the DBCluster.
    * `status` - The status of the DBCluster. Valid values: `Running`,`Creating`,`Deleting`,`Restarting`,`Preparing`.
    * `db_cluster_access_white_list` - The db cluster access white list.
        * `db_cluster_ip_array_attribute` - Field `db_cluster_ip_array_attribute` has been removed from provider.
        * `db_cluster_ip_array_name` - Whitelist group name.
        * `security_ip_list` - The IP address list under the whitelist group.