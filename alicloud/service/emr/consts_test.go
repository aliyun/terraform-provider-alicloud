package emr_test

// EMR-specific acceptance-test fixtures relocated from the flat suite's
// alicloud/service_alicloud_common_test.go during the alicloud/service/emr
// product sub-package migration. The raw HCL is preserved verbatim.

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain delegates to the plugin-testing entrypoint so the package-level
// sweepers (alicloud_emr_cluster, alicloud_emrv2_cluster) registered via
// resource.AddTestSweepers are honored by `go test` and the acceptance suite.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

const EmrCommonTestCase = `
data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_emr_main_versions" "default" {
	cluster_type = ["DRUID"]
}

data "alicloud_emr_instance_types" "default" {
    destination_resource = "InstanceType"
    cluster_type = "DRUID"
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER", "CORE"]
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = "DRUID"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = "DRUID"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  zone_id = "${data.alicloud_emr_instance_types.default.types.0.zone_id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
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
    force = true
}
`

const EmrV2CommonTestCase = `
data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

data "alicloud_kms_keys" "default" {
	status = "Enabled"
}

data "alicloud_zones" "default" {
	available_instance_type = "ecs.g7.xlarge"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "${var.name}"
}

resource "alicloud_ecs_capacity_reservation" default {
  capacity_reservation_name = "${var.name}"
  platform                  = "linux"
  end_time_type             = "Unlimited"
  instance_amount           = 1
  instance_type             = "ecs.g7.xlarge"
  match_criteria            = "Target"
  resource_group_id         = "${data.alicloud_resource_manager_resource_groups.default.ids.0}"
  zone_ids                  = ["${data.alicloud_zones.default.zones.0.id}"]

  tags = {
    Created = "terraform-example"
  }
}

resource "alicloud_security_group" "default" {
    security_group_name = "${var.name}"
    vpc_id              = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
  role_name                   = var.name
  assume_role_policy_document = <<EOF
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
`

const EmrV2AckConfigTestCase = `
data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

data "alicloud_kms_keys" "default" {
	status = "Enabled"
}

data "alicloud_zones" "default" {
	available_instance_type = "ecs.g7.xlarge"
}

resource "alicloud_vpc" "default" {
    vpc_name   = "${var.name}"
    cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "vswitches" {
    count      = 2
    vpc_id     = join("", alicloud_vpc.default.*.id)
    cidr_block = ["10.1.0.0/16", "10.2.0.0/16"][count.index]
    zone_id    = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "terway_vswitches" {
    count      = 2
    vpc_id     = join("", alicloud_vpc.default.*.id)
    cidr_block = ["10.4.0.0/16", "10.5.0.0/16"][count.index]
    zone_id    = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "default" {
    vpc_id = "${alicloud_vpc.default.id}"
    cidr_block = "10.6.0.0/16"
    zone_id = "${data.alicloud_zones.default.zones.0.id}"
    vswitch_name = "${var.name}"
}

resource "alicloud_ecs_key_pair" "default" {
    key_pair_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    security_group_name = "${var.name}"
    vpc_id              = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
    role_name                   = var.name
    assume_role_policy_document = <<EOF
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

resource "alicloud_cs_managed_kubernetes" "k8s" {
	name               = "${var.name}"
	cluster_spec       = "ack.pro.small"
	worker_vswitch_ids = split(",", join(",", alicloud_vswitch.vswitches.*.id))
	pod_vswitch_ids    = split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
	new_nat_gateway    = true
	node_cidr_mask     = 24
	proxy_mode         = "ipvs"
	service_cidr       = "192.168.0.0/16"

    deletion_protection = "false"

    addons {
        name = "terway-eniip"
    }
    addons {
        name = "csi-plugin"
    }
    addons {
        name = "csi-provisioner"
    }
    addons {
        name = "logtail-ds"
        config = jsonencode({
            IngressDashboardEnabled = "true"
      })
    }
    addons {
        name = "nginx-ingress-controller"
        config = jsonencode({
            IngressSlbNetworkType = "internet"
        })
    }
    addons {
        name = "arms-prometheus"
    }
    addons {
        name = "ack-node-problem-detector"
        config = jsonencode({
        })
    }
}
`

const EmrHadoopClusterTestCase = `
data "alicloud_emr_main_versions" "default" {
	cluster_type = ["DRUID"]
}

data "alicloud_db_zones" "default" {
	engine = "MySQL"
	engine_version = "8.0"
	category = "Basic"
	instance_charge_type = "PostPaid"
	db_instance_storage_type = "cloud_essd"
}

data "alicloud_emr_instance_types" "default" {
	destination_resource = "InstanceType"
	cluster_type = "DRUID"
	zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
	support_local_storage = false
	instance_charge_type = "PostPaid"
	support_node_type = ["MASTER", "CORE"]
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = "DRUID"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = "DRUID"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_db_instance_classes" "default" {
	zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
	engine = "MySQL"
	engine_version = "8.0"
	category = "Basic"
	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
	vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "8.0"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
	instance_charge_type = "Postpaid"
	db_instance_storage_type = "cloud_essd"
	vswitch_id = "${alicloud_vswitch.default.id}"
	instance_name = "${var.name}"
	security_ips = ["${alicloud_vswitch.default.cidr_block}"]
}

resource "alicloud_rds_account" "default" {
	db_instance_id = "${alicloud_db_instance.default.id}"
	account_type = "Normal"
	account_name = "taihao"
	account_password = "EMRtest1234!"
	account_description = "tf-test"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
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
    force = true
}
`

const EmrGatewayTestCase = `
data "alicloud_emr_main_versions" "default" {
	cluster_type = ["DRUID"]
}

data "alicloud_emr_instance_types" "default" {
    destination_resource = "InstanceType"
    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER","CORE"]
}

data "alicloud_emr_instance_types" "gateway" {
    destination_resource = "InstanceType"
    cluster_type = "GATEWAY"
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["GATEWAY"]
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "gateway_data_disk" {
	destination_resource = "DataDisk"
	cluster_type = "GATEWAY"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.gateway.types.0.id
	zone_id = data.alicloud_emr_instance_types.gateway.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "gateway_system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = "GATEWAY"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.gateway.types.0.id
	zone_id = data.alicloud_emr_instance_types.gateway.types.0.zone_id
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  zone_id = "${data.alicloud_emr_instance_types.default.types.0.zone_id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
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
    force = true
}

resource "alicloud_emr_cluster" "default" {
    name = "${var.name}"

    emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

    host_group {
        host_group_name = "master_group"
        host_group_type = "MASTER"
        node_count = "2"
        instance_type = data.alicloud_emr_instance_types.default.types.0.id
        disk_type = data.alicloud_emr_disk_types.data_disk.types.0.value
        disk_capacity = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
        disk_count = "1"
        sys_disk_type = data.alicloud_emr_disk_types.system_disk.types.0.value
		sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
    }

	host_group {
        host_group_name = "core_group"
        host_group_type = "CORE"
        node_count = "2"
        instance_type = data.alicloud_emr_instance_types.default.types.0.id
        disk_type = data.alicloud_emr_disk_types.data_disk.types.0.value
        disk_capacity = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
        disk_count = "4"
        sys_disk_type = data.alicloud_emr_disk_types.system_disk.types.0.value
        sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
    }

    high_availability_enable = true
    meta_store_type = "local"
    zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
    security_group_id = alicloud_security_group.default.id
    is_open_public_ip = true
    charge_type = "PostPaid"
    vswitch_id = alicloud_vswitch.default.id
    user_defined_emr_ecs_role = alicloud_ram_role.default.name
    ssh_enable = true
    master_pwd = "ABCtest1234!"
}
`
const EmrLocalStorageTestCase = `
data "alicloud_emr_main_versions" "default" {
	cluster_type = ["DRUID"]
}

data "alicloud_emr_instance_types" "local_disk" {
    destination_resource = "InstanceType"
    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = true
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER","CORE"]
}

data "alicloud_emr_instance_types" "cloud_disk" {
    destination_resource = "InstanceType"
    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER"]
    zone_id = data.alicloud_emr_instance_types.local_disk.types.0.zone_id
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.cloud_disk.types.0.id
	zone_id = data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.cloud_disk.types.0.id
	zone_id = data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  zone_id = "${data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
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
    force = true
}
`
