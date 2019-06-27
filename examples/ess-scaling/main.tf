data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex  = "^centos_6\\w{1,5}[64].*"
}

// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_instance_type     = "${var.ecs_instance_type}"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count       = "${var.vswitch_id == "" ? (var.vpc_id == "" ? 1 : 0) : 0}"
  name        = "${var.vpc_name}"
  cidr_block  = "${var.vpc_cidr}"
  description = "${var.vpc_description}"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitch" {
  count             = "${var.vswitch_id == "" ? 1 : 0}"
  vpc_id            = "${var.vpc_id != "" ? var.vpc_id : alicloud_vpc.vpc.id}"
  cidr_block        = "${var.vswitch_cidr}"
  availability_zone = "${var.availability_zone == "" ? data.alicloud_zones.default.zones.0.id : var.availability_zone}"
  name              = "${var.vswitch_name}"
  description       = "${var.vswitch_description}"
}

resource "alicloud_security_group" "sg" {
  name        = "${var.security_group_name}"
  vpc_id      = "${var.vpc_id == "" ? alicloud_vpc.vpc.id : var.vpc_id}"
  description = "tf-sg"
}

resource "alicloud_security_group_rule" "ssh-in" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alicloud_security_group.sg.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_ess_scaling_group" "scaling" {
  min_size           = "${var.scaling_min_size}"
  max_size           = "${var.scaling_max_size}"
  scaling_group_name = "tf-example-scaling"
  removal_policies   = "${var.removal_policies}"
  vswitch_ids        = ["${var.vswitch_id == "" ? alicloud_vswitch.vswitch.id : var.vswitch_id}"]
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id  = "${alicloud_ess_scaling_group.scaling.id}"
  active            = true
  enable            = "${var.enable}"
  image_id          = "${data.alicloud_images.ecs_image.images.0.id}"
  instance_type     = "${var.ecs_instance_type}"
  security_group_id = "${alicloud_security_group.sg.id}"
  key_name          = "${alicloud_key_pair.key.id}"
  role_name         = "${alicloud_ram_role.role.id}"
  force_delete      = "true"
}

resource "alicloud_key_pair" "key" {
  key_name = "my-key-pair-for-ess"
}

resource "alicloud_ram_role" "role" {
  name        = "EcsRamRoleTest"
  services    = ["ecs.aliyuncs.com"]
  description = "New role for ECS."
  force       = true
}

resource "alicloud_ram_policy" "policy" {
  name = "EcsRamRolePolicyTest"

  document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "ecs:*"
        ],
        "Effect": "Allow",
        "Resource": [
          "*"
        ]
      }
    ],
      "Version": "1"
  }
  EOF

  description = "New role policy for ECS."
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "role-policy" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name   = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}
