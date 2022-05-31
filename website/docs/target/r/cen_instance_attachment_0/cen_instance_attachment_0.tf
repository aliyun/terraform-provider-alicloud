# Create a new instance-attachment and use it to attach one child instance to a new CEN
variable "name" {
  default = "tf-testAccCenInstanceAttachmentBasic"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = "terraform01"
}

resource "alicloud_vpc" "vpc" {
  name       = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance_attachment" "foo" {
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc.id
  child_instance_type      = "VPC"
  child_instance_region_id = "cn-beijing"
}
