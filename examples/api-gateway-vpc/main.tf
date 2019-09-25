data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccInstanceConfigVPC"
}

resource "alicloud_vpc" "foo" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = alicloud_vpc.foo.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_security_group" "tf_test_foo" {
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.foo.id
}

resource "alicloud_instance" "foo" {
  vswitch_id = alicloud_vswitch.foo.id
  image_id   = data.alicloud_images.default.images[0].id

  # series III
  instance_type        = data.alicloud_instance_types.default.instance_types[0].id
  system_disk_category = "cloud_efficiency"

  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.tf_test_foo.id]
  instance_name              = var.name
}

resource "alicloud_api_gateway_vpc_access" "foo" {
  name        = "tf-testAccApiGatewayVpc"
  vpc_id      = alicloud_vpc.foo.id
  instance_id = alicloud_instance.foo.id
  port        = 8080
}

