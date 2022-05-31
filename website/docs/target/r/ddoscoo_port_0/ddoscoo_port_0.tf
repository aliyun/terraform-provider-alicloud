resource "alicloud_ddoscoo_instance" "example" {
  name              = "yourDdoscooInstanceName"
  bandwidth         = "30"
  base_bandwidth    = "30"
  service_bandwidth = "100"
  port_count        = "50"
  domain_count      = "50"
}

resource "alicloud_ddoscoo_port" "example" {
  instance_id       = alicloud_ddoscoo_instance.example.id
  frontend_port     = "7001"
  frontend_protocol = "tcp"
  real_servers      = ["1.1.1.1", "2.2.2.2"]
}

