variable "name" {
  default = "tf-testacc"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_testacc"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = "cn-hangzhou:tfacctest"
  namespace_name        = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = "tf-testaccDescription"
  app_name        = "tf-testaccAppName"
  namespace_id    = alicloud_sae_namespace.default.id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  vswitch_id      = alicloud_vswitch.vsw.id
  timezone        = "Asia/Beijing"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}
