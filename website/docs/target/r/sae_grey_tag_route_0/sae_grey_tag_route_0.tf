variable "name" {
  default = "tf-testacc"
}

variable "region" {
  default = "cn-hangzhou"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = join(":", [var.region, var.name])
  namespace_name        = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.namespace_id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}

resource "alicloud_sae_grey_tag_route" "default" {
  grey_tag_route_name = var.name
  description         = var.name
  app_id              = alicloud_sae_application.default.id
  sc_rules {
    items {
      type     = "param"
      name     = "tftest"
      operator = "rawvalue"
      value    = "test"
      cond     = "=="
    }
    path      = "/tf/test"
    condition = "AND"
  }

  dubbo_rules {
    items {
      cond     = "=="
      expr     = ".key1"
      index    = "1"
      operator = "rawvalue"
      value    = "value1"
    }
    condition    = "OR"
    group        = "DUBBO"
    method_name  = "test"
    service_name = "com.test.service"
    version      = "1.0.0"
  }
}
