variable "name" {
  default = "example_value"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_slb" "default" {
  name          = var.name
  specification = "slb.s2.small"
  vswitch_id    = data.alicloud_vswitches.default.ids.0
}

variable "namespace_id" {
  default = "cn-hangzhou:yourname"
}

resource "alicloud_sae_namespace" "default" {
  namespace_id          = var.namespace_id
  namespace_name        = var.name
  namespace_description = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = "your_app_description"
  app_name        = "your_app_name"
  namespace_id    = "your_namespace_id"
  package_url     = "your_package_url"
  package_type    = "your_package_url"
  jdk             = "jdk_specifications"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  replicas        = "your_replicas"
  cpu             = "cpu_specifications"
  memory          = "memory_specifications"

}
resource "alicloud_sae_ingress" "default" {
  slb_id        = alicloud_slb.default.id
  namespace_id  = alicloud_sae_namespace.default.id
  listener_port = "your_listener_port"
  rules {
    app_id         = alicloud_sae_application.default.id
    container_port = "your_container_port"
    domain         = "your_domain"
    app_name       = "your_name"
    path           = "your_path"
  }
}
