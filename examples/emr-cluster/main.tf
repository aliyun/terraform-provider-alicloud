data "alicloud_emr_main_versions" "default" {
}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["MASTER", "CORE", "TASK", "GATEWAY"]
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

  vpc_name   = var.vpc_name
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

  zone_id      = var.availability_zone == "" ? data.alicloud_emr_instance_types.default.types.0.zone_id : var.availability_zone
  vswitch_name = var.vswitch_name
  cidr_block   = var.vswitch_cidr
  vpc_id       = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
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
  name = "terraform-resize-test-0926"

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
    charge_type       = "PostPaid"
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
    charge_type       = "PostPaid"
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
    charge_type       = "PostPaid"
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

resource "alicloud_emr_cluster" "gateway" {
  name = "terraform-gateway-test-1101"

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
  related_cluster_id        = alicloud_emr_cluster.default.id
}
