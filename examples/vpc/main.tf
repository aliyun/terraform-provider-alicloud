provider "alicloud" {
  configuration_source = "terraform-provider-alicloud/examples/vpc"
}

resource "alicloud_vpc" "main" {
  vpc_name   = var.long_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "main" {
  vpc_id     = alicloud_vpc.main.id
  count      = length(var.cidr_blocks)
  cidr_block = var.cidr_blocks["az${count.index}"]
  zone_id    = var.availability_zones

  depends_on = [alicloud_vpc.main]
}

resource "alicloud_nat_gateway" "main" {
  vpc_id        = alicloud_vpc.main.id
  specification = "Small"
  name          = "from-tf-example"
}

resource "alicloud_eip_address" "foo" {
}

resource "alicloud_eip_association" "foo" {
  allocation_id = alicloud_eip_address.foo.id
  instance_id   = alicloud_nat_gateway.main.id
}
