---
subcategory: "E-MapReduce"
layout: "alicloud"
page_title: "Alicloud: alicloud_emr_cluster"
sidebar_current: "docs-alicloud-resource-emr-cluster"
description: |-
  Provides a EMR Cluster resource.
---

# alicloud\_emr\_cluster

Provides a EMR Cluster resource. With this you can create, read, and release  EMR Cluster. 

-> **NOTE:** Available in 1.57.0+.

## Example Usage

#### 1. Create A Cluster

```terraform
data "alicloud_emr_main_versions" "default" {}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["MASTER", "CORE", "TASK"]
}

data "alicloud_emr_disk_types" "data_disk" {
  destination_resource = "DataDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
  destination_resource = "SystemDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "vpc" {
  count = var.vpc_id == "" ? 1 : 0

  name       = var.vpc_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_security_group" "default" {
  count = var.security_group_id == "" ? 1 : 0

  name   = var.security_group_name
  vpc_id = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count = var.vswitch_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? data.alicloud_emr_instance_types.default.types.0.zone_id : var.availability_zone
  vswitch_name      = var.vswitch_name
  cidr_block        = var.vswitch_cidr
  vpc_id            = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// Ram role Resource for Module
resource "alicloud_ram_role" "default" {
  name        = var.ram_name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}

resource "alicloud_emr_cluster" "default" {
  name = "terraform-create-cluster-test"

  emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

  cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

  host_group {
    host_group_name   = "master_group"
    host_group_type   = "MASTER"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "1"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "core_group"
    host_group_type   = "CORE"
    node_count        = "3"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "task_group"
    host_group_type   = "TASK"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  high_availability_enable  = true
  zone_id                   = data.alicloud_emr_instance_types.default.types.0.zone_id
  security_group_id         = var.security_group_id == "" ? alicloud_security_group.default[0].id : var.security_group_id
  is_open_public_ip         = true
  charge_type               = "PostPaid"
  vswitch_id                = var.vswitch_id == "" ? alicloud_vswitch.vswitch[0].id : var.vswitch_id
  user_defined_emr_ecs_role = alicloud_ram_role.default.name
  ssh_enable                = true
  master_pwd                = "ABCtest1234!"
}
```

#### 2. Scale Up
The hosts of EMR Cluster are orginized as host group. Scaling up/down is operating host group. 

In the case of scaling up cluster, we should add the node_count of some host group. 

-> **NOTE:** Scaling up is only applicable to CORE and TASK group. Cost time of scaling up will vary with the number of scaling-up nodes. 
Scaling down is only applicable to TASK group. If you want to scale down CORE group, please submit tickets or contact EMR support team.

As the following case, we scale up the TASK group 2 nodes by increasing host_group.node_count by 2.

```terraform
data "alicloud_emr_main_versions" "default" {}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["MASTER", "CORE", "TASK"]
}

data "alicloud_emr_disk_types" "data_disk" {
  destination_resource = "DataDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
  destination_resource = "SystemDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "vpc" {
  count = var.vpc_id == "" ? 1 : 0

  name       = var.vpc_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_security_group" "default" {
  count = var.security_group_id == "" ? 1 : 0

  name   = var.security_group_name
  vpc_id = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count = var.vswitch_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? data.alicloud_emr_instance_types.default.types.0.zone_id : var.availability_zone
  vswitch_name      = var.vswitch_name
  cidr_block        = var.vswitch_cidr
  vpc_id            = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// Ram role Resource for Module
resource "alicloud_ram_role" "default" {
  name        = var.ram_name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}

resource "alicloud_emr_cluster" "default" {
  name = "terraform-resize-cluster-test"

  emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

  cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

  host_group {
    host_group_name   = "master_group"
    host_group_type   = "MASTER"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "1"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "core_group"
    host_group_type   = "CORE"
    node_count        = "3"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "task_group"
    host_group_type   = "TASK"
    node_count        = "4"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  high_availability_enable  = true
  zone_id                   = data.alicloud_emr_instance_types.default.types.0.zone_id
  security_group_id         = var.security_group_id == "" ? alicloud_security_group.default[0].id : var.security_group_id
  is_open_public_ip         = true
  charge_type               = "PostPaid"
  vswitch_id                = var.vswitch_id == "" ? alicloud_vswitch.vswitch[0].id : var.vswitch_id
  user_defined_emr_ecs_role = alicloud_ram_role.default.name
  ssh_enable                = true
  master_pwd                = "ABCtest1234!"
}
```

#### 3. Scale Down

In the case of scaling down a cluster, we need to specified the host group and the instance list. 

-> **NOTE:** Graceful decommission of hadoop cluster has been supported Available in 1.168.0+.

The following is an example. We scale down the cluster by decreasing the node count by 2, and specifying the scale-down instance list.

```terraform
data "alicloud_emr_main_versions" "default" {}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["MASTER", "CORE", "TASK"]
}

data "alicloud_emr_disk_types" "data_disk" {
  destination_resource = "DataDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
  destination_resource = "SystemDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "vpc" {
  count = var.vpc_id == "" ? 1 : 0

  name       = var.vpc_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_security_group" "default" {
  count = var.security_group_id == "" ? 1 : 0

  name   = var.security_group_name
  vpc_id = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count = var.vswitch_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? data.alicloud_emr_instance_types.default.types.0.zone_id : var.availability_zone
  vswitch_name      = var.vswitch_name
  cidr_block        = var.vswitch_cidr
  vpc_id            = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// Ram role Resource for Module
resource "alicloud_ram_role" "default" {
  name        = var.ram_name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}

resource "alicloud_emr_cluster" "default" {
  name = "terraform-resize-cluster-test"

  emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

  cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

  host_group {
    host_group_name   = "master_group"
    host_group_type   = "MASTER"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "1"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "core_group"
    host_group_type   = "CORE"
    node_count        = "3"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "task_group"
    host_group_type   = "TASK"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  high_availability_enable  = true
  zone_id                   = data.alicloud_emr_instance_types.default.types.0.zone_id
  security_group_id         = var.security_group_id == "" ? alicloud_security_group.default[0].id : var.security_group_id
  is_open_public_ip         = true
  charge_type               = "PostPaid"
  vswitch_id                = var.vswitch_id == "" ? alicloud_vswitch.vswitch[0].id : var.vswitch_id
  user_defined_emr_ecs_role = alicloud_ram_role.default.name
  ssh_enable                = true
  master_pwd                = "ABCtest1234!"
}
```

#### 4. Create a emr gateway cluster

```terraform
data "alicloud_emr_main_versions" "default" {}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["GATEWAY"]
}

data "alicloud_emr_disk_types" "data_disk" {
  destination_resource = "DataDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
  destination_resource = "SystemDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "vpc" {
  count = var.vpc_id == "" ? 1 : 0

  name       = var.vpc_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_security_group" "default" {
  count = var.security_group_id == "" ? 1 : 0

  name   = var.security_group_name
  vpc_id = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count = var.vswitch_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? data.alicloud_emr_instance_types.default.types.0.zone_id : var.availability_zone
  vswitch_name      = var.vswitch_name
  cidr_block        = var.vswitch_cidr
  vpc_id            = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

// Ram role Resource for Module
resource "alicloud_ram_role" "default" {
  name        = var.ram_name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}

resource "alicloud_emr_cluster" "gateway" {
  name = "terraform-gateway-cluster-test"

  emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

  # supported 'GATEWAY' available in 1.61.0+.
  cluster_type = "GATEWAY"

  host_group {
    host_group_name   = "master_group"
    host_group_type   = "GATEWAY"
    node_count        = "1"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "1"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  high_availability_enable  = true
  zone_id                   = data.alicloud_emr_instance_types.default.types.0.zone_id
  security_group_id         = var.security_group_id == "" ? alicloud_security_group.default[0].id : var.security_group_id
  is_open_public_ip         = true
  charge_type               = "PostPaid"
  vswitch_id                = var.vswitch_id == "" ? alicloud_vswitch.vswitch[0].id : var.vswitch_id
  user_defined_emr_ecs_role = alicloud_ram_role.default.name
  ssh_enable                = true
  master_pwd                = "ABCtest1234!"

  // the gateway related emr cluster id, if you want to bind a cluster dynamically, please refer to the examples.
  related_cluster_id = related_cluster_id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of emr cluster. The name length must be less than 64. Supported characters: chinese character, english character, number, "-", "_".
* `emr_ver` - (Required, ForceNew) EMR Version, e.g. EMR-3.22.0. You can find the all valid EMR Version in emr web console.
* `cluster_type` - (Required, ForceNew) EMR Cluster Type, e.g. HADOOP, KAFKA, DRUID, GATEWAY etc. You can find all valid EMR cluster type in emr web console. Supported 'GATEWAY' available in 1.61.0+.
* `charge_type` - (Required, ForceNew) Charge Type for this cluster. Supported value: PostPaid or PrePaid. Default value: PostPaid.
* `zone_id` - (Required, ForceNew) Zone ID, e.g. cn-huhehaote-a
* `security_group_id` (Optional, ForceNew) Security Group ID for Cluster, you can also specify this key for each host group.
* `vswitch_id` (Optional, ForceNew) Global vswitch id, you can also specify it in host group.
* `option_software_list` (Optional, ForceNew) Optional software list.
* `high_availability_enable` (Optional, ForceNew) High Available for HDFS and YARN. If this is set true, MASTER group must have two nodes.
* `use_local_metadb` (Optional, ForceNew) Use local metadb. Default is false.
* `ssh_enable` (Optional, ForceNew) If this is set true, we can ssh into cluster. Default value is false.
* `master_pwd` (Optional, ForceNew) Master ssh password.
* `eas_enable` (Optional, ForceNew) High security cluster (true) or not. Default value is false.
* `user_defined_emr_ecs_role` (Optional, ForceNew) Alicloud EMR uses roles to perform actions on your behalf when provisioning cluster resources, running applications, dynamically scaling resources. EMR uses the following roles when interacting with other Alicloud services. Default value is AliyunEmrEcsDefaultRole.
* `key_pair_name` (Optional, ForceNew) Ssh key pair.
* `deposit_type` (Optional, ForceNew) Cluster deposit type, HALF_MANAGED or FULL_MANAGED.
* `related_cluster_id` (Optional, ForceNew) This specify the related cluster id, if this cluster is a Gateway.
* `host_group` - (Optional) Groups of Host, You can specify MASTER as a group, CORE as a group (just like the above example).
* `tags` - (Optional, Available in v1.67.0+) A mapping of tags to assign to the resource.
* `period` - (Optional, Available in v1.101.0+) If charge type is PrePaid, this should be specified, unit is month. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36.
* `is_open_public_ip` (Optional, ForceNew) Whether the MASTER node has a public IP address enabled. Default value is false.
* `bootstrap_action` (Optional) Boot action parameters.

#### Block host_group

The host_group mapping supports the following: 

* `host_group_name` - (Required, ForceNew) host group name.
* `host_group_type` - (Required) host group type, supported value: MASTER, CORE or TASK, supported 'GATEWAY' available in 1.61.0+.
* `charge_type` - (Optional) Charge Type for this group of hosts: PostPaid or PrePaid. If this is not specified, charge type will follow global charge_type value.
* `period` - (Optional) If charge type is PrePaid, this should be specified, unit is month. Supported value: 1、2、3、4、5、6、7、8、9、12、24、36.
* `node_count` - (Required) Host number in this group.
* `instance_type` - (Required) Host Ecs instance type.
* `disk_type` - (Required) Data disk type. Supported value: cloud,cloud_efficiency,cloud_ssd,local_disk,cloud_essd.
* `disk_capacity` - (Required) Data disk capacity.
* `disk_count` - (Required) Data disk count.
* `sys_disk_type` - (Required) System disk type. Supported value: cloud,cloud_efficiency,cloud_ssd,cloud_essd.
* `sys_disk_capacity` - (Required) System disk capacity.
* `auto_renew` - (Optional) Auto renew for prepaid, ’true’ or ‘false’ . Default value: false.
* `instance_list` - (Optional) Instance list for cluster scale down. This value follows the json format, e.g. ["instance_id1","instance_id2"]. escape character for " is \".
* `enable_graceful_decommission` - (Optional, Available in 1.168.0+) Enable hadoop cluster of task node graceful decommission, ’true’ or ‘false’ . Default value: false.
* `decommission_timeout` - (Optional, Available in 1.168.0+) Graceful decommission timeout, unit: seconds.

#### Block bootstrap_action

The bootstrap_action mapping supports the following:

* `name` - (Optional, Available in 1.71.2+) bootstrap action name.
* `path` - (Optional, Available in 1.71.2+) bootstrap action path, e.g. "oss://bucket/path".
* `arg` - (Optional, Available in 1.71.2+) bootstrap action args, e.g. "--a=b".
* `execution_target` - (Optional, Available in 1.173.0+) bootstrap action execution target, you can specify the host group name, e.g. "core_group". If this is not specified, the bootstrap action execution target is whole cluster.
* `execution_moment` - (Optional, Available in 1.173.0+) bootstrap action execution moment, ’BEFORE_INSTALL’ or ‘AFTER_STARTED’ . Default value: "BEFORE_INSTALL".
* `execution_fail_strategy` - (Optional, Available in 1.173.0+) bootstrap action execution fail strategy, ’FAILED_BLOCKED’ or ‘FAILED_CONTINUE’ . Default value: "FAILED_BLOCKED

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the cluster (until it reaches the initial `IDLE` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the instance.

## Attribute Reference

The following attributes are exported:

* `id` - The cluster ID.

## Import

Aliclioud E-MapReduce cluster can be imported using the id e.g.

```
$ terraform import alicloud_emr_cluster.default C-B47FB8FE96C67XXXX
```

