resource "alicloud_ga_accelerator" "example" {
  duration            = 3
  spec                = "2"
  accelerator_name    = "ga-tf"
  auto_use_coupon     = false
  description         = "ga-tf description"
  auto_renew_duration = "2"
  renewal_status      = "AutoRenewal"
}

resource "alicloud_ga_bandwidth_package" "example" {
  type           = "Basic"
  bandwidth      = 20
  bandwidth_type = "Basic"
  duration       = 1
  timeouts {
    create = "5m"
  }
  auto_pay               = true
  payment_type           = "Subscription"
  billing_type           = "PayByTraffic"
  ratio                  = 40
  auto_use_coupon        = false
  bandwidth_package_name = "bandwidth_package_name_tf"
  description            = "bandwidth_package_name_tf_description"

}

resource "alicloud_ga_bandwidth_package_attachment" "example" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.example.id
}


resource "alicloud_ga_listener" "example" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.example]
  accelerator_id = alicloud_ga_accelerator.example.id
  port_ranges {
    from_port = 60
    to_port   = 60
  }
  client_affinity = "SOURCE_IP"
  description     = "alicloud_ga_listener_description"
  name            = "alicloud_ga_listener_tf"
  protocol        = "HTTP"
  proxy_protocol  = true
}


resource "alicloud_ga_ip_set" "example" {
  depends_on           = [alicloud_ga_bandwidth_package_attachment.example]
  accelerate_region_id = "cn-shanghai"
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth            = "20"
}

resource "alicloud_eip_address" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id = alicloud_ga_accelerator.example.id
  endpoint_configurations {
    endpoint                     = alicloud_eip_address.example.ip_address
    type                         = "PublicIp"
    weight                       = "20"
    enable_clientip_preservation = true
  }
  endpoint_group_region         = "cn-shanghai"
  listener_id                   = alicloud_ga_listener.example.id
  description                   = "alicloud_ga_endpoint_group_description"
  endpoint_group_type           = "default"
  endpoint_request_protocol     = "HTTPS"
  health_check_interval_seconds = 4
  health_check_path             = "/path"
  name                          = "alicloud_ga_endpoint_group_name"
  threshold_count               = 4
  traffic_percentage            = 20
  port_overrides {
    endpoint_port = 80
    listener_port = 60
  }
}

resource "alicloud_ga_endpoint_group" "virtual" {
  accelerator_id = alicloud_ga_accelerator.example.id
  endpoint_configurations {
    endpoint                     = alicloud_eip_address.example.ip_address
    type                         = "PublicIp"
    weight                       = "20"
    enable_clientip_preservation = true
  }
  endpoint_group_region = "cn-shanghai"
  listener_id           = alicloud_ga_listener.example.id

  description                   = "alicloud_ga_endpoint_group_description"
  endpoint_group_type           = "virtual"
  endpoint_request_protocol     = "HTTPS"
  health_check_interval_seconds = 4
  health_check_path             = "/path"
  name                          = "alicloud_ga_endpoint_group_name"
  threshold_count               = 4
  traffic_percentage            = 20
  port_overrides {
    endpoint_port = 80
    listener_port = 60
  }
}

resource "alicloud_ga_forwarding_rule" "example" {
  accelerator_id = alicloud_ga_accelerator.example.id
  listener_id    = alicloud_ga_listener.example.id
  rule_conditions {
    rule_condition_type = "Path"
    path_config {
      values = ["/testpathconfig"]
    }
  }
  rule_conditions {
    rule_condition_type = "Host"
    host_config {
      values = ["www.test.com"]
    }
  }
  rule_actions {
    order            = "40"
    rule_action_type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        endpoint_group_id = alicloud_ga_endpoint_group.default.id
      }
    }
  }
  priority             = 2
  forwarding_rule_name = "forwarding_rule_name"
}
