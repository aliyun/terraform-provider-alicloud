module "vpc" {
  availability_zones = "${var.availability_zones}"
  source             = "../vpc"
  short_name         = "${var.short_name}"
  region             = "${var.region}"
}

module "security-groups" {
  source     = "../vpc-cluster-sg"
  short_name = "${var.short_name}"
  vpc_id     = "${module.vpc.vpc_id}"
}

module "control-nodes" {
  source       = "../ecs-vpc"
  count_format = "${var.control_count}"
  role         = "control"
  ecs_type     = "${var.control_ecs_type}"
  disk_size    = "${var.control_disk_size}"
  ssh_username = "${var.ssh_username}"
  short_name   = "${var.short_name}"
  sg_id        = "${module.security-groups.control_security_group}"
  vswitch_id   = "${module.vpc.vswitch_ids}"
}

module "edge-nodes" {
  source       = "../ecs-vpc"
  count_format = "${var.edge_count}"
  role         = "edge"
  ecs_type     = "${var.edge_ecs_type}"
  ssh_username = "${var.ssh_username}"
  short_name   = "${var.short_name}"
  sg_id        = "${module.security-groups.edge_security_group}"
  vswitch_id   = "${element(split(",", module.vpc.vswitch_ids), 0)}"
}

module "worker-nodes" {
  source       = "../ecs-vpc"
  count_format = "${var.worker_count}"
  role         = "worker"
  ecs_type     = "${var.worker_ecs_type}"
  ssh_username = "${var.ssh_username}"
  short_name   = "${var.short_name}"
  sg_id        = "${module.security-groups.worker_security_group}"
  vswitch_id   = "${module.vpc.vswitch_ids}"
}
