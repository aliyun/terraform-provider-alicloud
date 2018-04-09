// Provider specific configs
provider "alicloud" {
  access_key = "${var.alicloud_access_key}"
  secret_key = "${var.alicloud_secret_key}"
  region = "${var.region}"
}

// Instance_types data source for instance_type
data "alicloud_instance_types" "default" {
  cpu_core_count = "${var.cpu_core_count}"
  memory_size = "${var.memory_size}"
}

// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count = "${var.vpc_id == "" ? 1 : 0}"
  cidr_block = "${var.vpc_cidr}"
  name = "${var.vpc_name == "" ? var.example_name : var.vpc_name}"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count = "${length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)}"
  vpc_id = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
  cidr_block = "${element(var.vswitch_cidrs, count.index)}"
  availability_zone = "${lookup(data.alicloud_zones.default.zones[count.index%length(data.alicloud_zones.default.zones)], "id")}"
  name = "${var.vswitch_name_prefix == "" ? format("%s-%s", var.example_name, format(var.number_format, count.index+1)) : format("%s-%s", var.vswitch_name_prefix, format(var.number_format, count.index+1))}"
}

resource "alicloud_nat_gateway" "default" {
  count = "${var.new_nat_gateway == true ? 1 : 0}"
  vpc_id = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
  name = "${var.example_name}"
}

resource "alicloud_eip" "default" {
  count = "${var.new_nat_gateway == "true" ? 1 : 0}"
  bandwidth = 10
}

resource "alicloud_eip_association" "default" {
  count = "${var.new_nat_gateway == "true" ? 1 : 0}"
  allocation_id = "${alicloud_eip.default.id}"
  instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default"{
  count = "${var.new_nat_gateway == "false" ? 0 : length(var.vswitch_ids) > 0 ? length(var.vswitch_ids) : length(var.vswitch_cidrs)}"
  snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
  source_vswitch_id = "${length(var.vswitch_ids) > 0 ? element(split(",", join(",", var.vswitch_ids)), count.index%length(split(",", join(",", var.vswitch_ids)))) : length(var.vswitch_cidrs) < 1 ? "" : element(split(",", join(",", alicloud_vswitch.vswitches.*.id)), count.index%length(split(",", join(",", alicloud_vswitch.vswitches.*.id))))}"
  snat_ip = "${alicloud_eip.default.ip_address}"
}

resource "alicloud_cs_kubernetes" "k8s" {
  count = "${var.k8s_number}"
  name = "${var.k8s_name_prefix == "" ? format("%s-%s", var.example_name, format(var.number_format, count.index+1)) : format("%s-%s", var.k8s_name_prefix, format(var.number_format, count.index+1))}"
  vswitch_id = "${length(var.vswitch_ids) > 0 ? element(split(",", join(",", var.vswitch_ids)), count.index%length(split(",", join(",", var.vswitch_ids)))) : length(var.vswitch_cidrs) < 1 ? "" : element(split(",", join(",", alicloud_vswitch.vswitches.*.id)), count.index%length(split(",", join(",", alicloud_vswitch.vswitches.*.id))))}"
  new_nat_gateway = false
  master_instance_type = "${var.master_instance_type == "" ? data.alicloud_instance_types.default.instance_types.0.id : var.master_instance_type}"
  worker_instance_type = "${var.worker_instance_type == "" ? data.alicloud_instance_types.default.instance_types.0.id : var.worker_instance_type}"
  worker_number = "${var.k8s_worker_number}"
  master_disk_category = "${var.master_disk_category}"
  worker_disk_category = "${var.worker_disk_category}"
  master_disk_size = "${var.master_disk_size}"
  worker_disk_size = "${var.master_disk_size}"
  password = "${var.ecs_password}"
  pod_cidr = "${var.k8s_pod_cidr}"
  service_cidr = "${var.k8s_service_cidr}"
  enable_ssh = true
  install_cloud_monitor = true

  depends_on = ["alicloud_snat_entry.default"]
}
