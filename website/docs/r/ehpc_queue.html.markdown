---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_queue"
description: |-
  Provides a Alicloud Ehpc Queue resource.
---

# alicloud_ehpc_queue

Provides a Ehpc Queue resource.

E-HPC the compute queue of the cluster.

For information about Ehpc Queue and how to use it, see [What is Queue](https://next.api.alibabacloud.com/document/EHPC/2024-07-30/CreateQueue).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "queue_minimal_example_vpc" {
  is_default = false
  cidr_block = "10.0.0.0/8"
  vpc_name   = "example-cluster-vpc"
}

resource "alicloud_vswitch" "queue_minimal_example_vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.queue_minimal_example_vpc.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "example-cluster-vsw"
}

resource "alicloud_nas_file_system" "queue_minimal_example_nas" {
  description  = "example-cluster-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_nas_access_group" "queue_minimal_example_access_group" {
  access_group_type = "Vpc"
  description       = var.name
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_security_group" "queue_minimal_example_security_group" {
  vpc_id              = alicloud_vpc.queue_minimal_example_vpc.id
  security_group_type = "normal"
}

resource "alicloud_nas_mount_target" "queue_minimal_example_mount_domain" {
  vpc_id            = alicloud_vpc.queue_minimal_example_vpc.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.queue_minimal_example_access_group.access_group_name
  vswitch_id        = alicloud_vswitch.queue_minimal_example_vswitch.id
  file_system_id    = alicloud_nas_file_system.queue_minimal_example_nas.id
}

resource "alicloud_nas_access_rule" "queue_minimal_example_access_rule" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.queue_minimal_example_access_group.access_group_name
  file_system_type  = alicloud_nas_file_system.queue_minimal_example_nas.file_system_type
  source_cidr_ip    = "10.0.0.0/8"
}

resource "alicloud_ehpc_cluster_v2" "queue_minimal_cluster_example" {
  depends_on = [alicloud_nas_access_rule.queue_minimal_example_access_rule]
  cluster_credentials {
    password = "aliHPC123"
  }
  cluster_vpc_id    = alicloud_vpc.queue_minimal_example_vpc.id
  cluster_category  = "Standard"
  cluster_mode      = "Integrated"
  security_group_id = alicloud_security_group.queue_minimal_example_security_group.id
  addons {
    version        = "1.0"
    services_spec  = <<EOF
[
        {
          "ServiceName": "SSH",
          "NetworkACL": [
            {
              "Port": 22,
              "SourceCidrIp": "0.0.0.0/0",
              "IpProtocol": "TCP"
            }
          ]
        },
        {
          "ServiceName": "VNC",
          "NetworkACL": [
            {
              "Port": 12016,
              "SourceCidrIp": "0.0.0.0/0",
              "IpProtocol": "TCP"
            }
          ]
        },
        {
          "ServiceName": "CLIENT",
          "ServiceAccessType": "URL",
          "ServiceAccessUrl": "https://ehpc-app.oss-cn-hangzhou.aliyuncs.com/ClientRelease/E-HPC-Client-Mac-zh-cn.zip",
          "NetworkACL": [
            {
              "Port": 12011,
              "SourceCidrIp": "0.0.0.0/0",
              "IpProtocol": "TCP"
            }
          ]
        }
      ]
  EOF
    resources_spec = <<EOF
{
        "EipResource": {
          "AutoCreate": true
        },
        "EcsResources": [
          {
            "ImageId": "centos_7_6_x64_20G_alibase_20211130.vhd",
            "EnableHT": true,
            "InstanceChargeType": "PostPaid",
            "InstanceType": "ecs.c7.xlarge",
            "SpotStrategy": "NoSpot",
            "SystemDisk": {
              "Category": "cloud_essd",
              "Size": 40,
              "Level": "PL0"
            },
            "DataDisks": [
              {
                "Category": "cloud_essd",
                "Size": 40,
                "Level": "PL0"
              }
            ]
          }
        ]
      }
  EOF
    name           = "Login"
  }
  cluster_name        = "minimal-example-cluster"
  deletion_protection = false
  shared_storages {
    mount_directory     = "/home"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_minimal_example_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_minimal_example_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/opt"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_minimal_example_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_minimal_example_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.queue_minimal_example_mount_domain.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.queue_minimal_example_nas.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  cluster_vswitch_id = alicloud_vswitch.queue_minimal_example_vswitch.id
  manager {
    manager_node {
      system_disk {
        category = "cloud_essd"
        size     = "40"
        level    = "PL0"
      }
      enable_ht            = true
      instance_charge_type = "PostPaid"
      image_id             = "centos_7_6_x64_20G_alibase_20211130.vhd"
      instance_type        = "ecs.c6.xlarge"
      spot_strategy        = "NoSpot"
    }
    scheduler {
      type    = "SLURM"
      version = "22.05.8"
    }
    dns {
      type    = "nis"
      version = "1.0"
    }
    directory_service {
      type    = "nis"
      version = "1.0"
    }
  }
}


resource "alicloud_ehpc_queue" "default" {
  cluster_id = alicloud_ehpc_cluster_v2.queue_minimal_cluster_example.id
  queue_name = "autoque1"
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Optional, ForceNew, Computed) The cluster ID.
You can call the [ListClusters](~~87116~~) operation to query the cluster ID.
* `compute_nodes` - (Optional, Computed, List) The hardware configurations of the compute nodes in the queue. Valid values of N: 1 to 10. See [`compute_nodes`](#compute_nodes) below.
* `enable_scale_in` - (Optional, Computed) Specifies whether to enable auto scale-in for the queue. Valid values:

  - true
  - false
* `enable_scale_out` - (Optional, Computed) Specifies whether to enable auto scale-out for the queue. Valid values:

  - true
  - false
* `hostname_prefix` - (Optional, Computed) The hostname prefix of the added compute nodes.
* `hostname_suffix` - (Optional, Computed) The hostname suffix of the compute nodes in the queue.
* `initial_count` - (Optional, ForceNew, Computed, Int) The initial number of compute nodes in the queue.
* `inter_connect` - (Optional, Computed) The type of the network for interconnecting compute nodes in the queue.
* `max_count` - (Optional, Computed, Int) The maximum number of compute nodes that the queue can contain.
* `min_count` - (Optional, Computed, Int) The minimum number of compute nodes that the queue must contain.
* `queue_name` - (Optional, ForceNew, Computed) The queue name.
* `vswitch_ids` - (Optional, List) The vSwitches available for use by compute nodes in the queue.

### `compute_nodes`

The compute_nodes supports the following:
* `auto_renew` - (Optional) AutoRenew
* `auto_renew_period` - (Optional, Int) AutoRenewPeriod
* `duration` - (Optional, Int) Duration
* `enable_ht` - (Optional, Computed) Whether HT is enabled for the computing node.
* `image_id` - (Optional) ImageId
* `instance_charge_type` - (Optional) InstanceChargeType
* `instance_type` - (Optional) InstanceTypes
* `period` - (Optional, Int) Period
* `period_unit` - (Optional) PeriodUnit
* `spot_price_limit` - (Optional, Float) SpotPriceLimit
* `spot_strategy` - (Optional) SpotStrategy
* `system_disk` - (Optional, List) SystemDisk See [`system_disk`](#compute_nodes-system_disk) below.

### `compute_nodes-system_disk`

The compute_nodes-system_disk supports the following:
* `category` - (Optional) Category
* `level` - (Optional) Level
* `size` - (Optional, Int) Size

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<queue_name>`.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Queue.
* `delete` - (Defaults to 5 mins) Used when delete the Queue.
* `update` - (Defaults to 5 mins) Used when update the Queue.

## Import

Ehpc Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_ehpc_queue.example <cluster_id>:<queue_name>
```