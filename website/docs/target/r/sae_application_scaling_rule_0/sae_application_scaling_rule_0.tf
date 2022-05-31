data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_sae_namespace" "default" {
  namespace_description = "example_value"
  namespace_id          = "example_value"
  namespace_name        = "example_value"
}
resource "alicloud_sae_application" "default" {
  app_description = "example_value"
  app_name        = "example_value"
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
resource "alicloud_sae_application_scaling_rule" "example" {
  app_id              = alicloud_sae_application.default.id
  scaling_rule_name   = "example-value"
  scaling_rule_enable = true
  scaling_rule_type   = "mix"
  scaling_rule_timer {
    begin_date = "2022-02-25"
    end_date   = "2022-03-25"
    period     = "* * *"
    schedules {
      at_time      = "08:00"
      max_replicas = 10
      min_replicas = 3
    }
    schedules {
      at_time      = "20:00"
      max_replicas = 50
      min_replicas = 3
    }
  }
  scaling_rule_metric {
    max_replicas = 50
    min_replicas = 3
    metrics {
      metric_type                       = "CPU"
      metric_target_average_utilization = 20
    }
    metrics {
      metric_type                       = "MEMORY"
      metric_target_average_utilization = 30
    }
    metrics {
      metric_type                       = "tcpActiveConn"
      metric_target_average_utilization = 20
    }
    scale_up_rules {
      step                         = 10
      disabled                     = false
      stabilization_window_seconds = 0
    }
    scale_down_rules {
      step                         = 10
      disabled                     = false
      stabilization_window_seconds = 10
    }
  }
}
