resource "alicloud_pvtz_zone" "main" {
  name = var.long_name
}

resource "alicloud_pvtz_zone_record" "main" {
  zone_id         = alicloud_pvtz_zone.main.id
  resource_record = var.resource_record
  type            = var.type
  value           = var.value
  priority        = var.priority
}

resource "alicloud_vpc" "main" {
  name       = var.long_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_pvtz_zone_attachment" "main" {
  zone_id = alicloud_pvtz_zone.main.id
  vpc_ids = [alicloud_vpc.main.id]
}

