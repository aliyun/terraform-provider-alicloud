resource "alicloud_slb" "instance" {
  name                 = "${var.slb_name}"
  internet_charge_type = "${var.internet_charge_type}"
  internet             = "${var.internet}"
}

resource "alicloud_slb_listener" "tcp" {
  load_balancer_id          = "${alicloud_slb.instance.id}"
  backend_port              = "22"
  frontend_port             = "22"
  protocol                  = "tcp"
  bandwidth                 = "10"
  health_check_type         = "tcp"
  persistence_timeout       = 3600
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx"
  health_check_timeout      = 8
  health_check_connect_port = 20
  health_check_uri          = "/console"
  acl_status                = "off"
  acl_type                  = "white"
  acl_id                    = "${alicloud_slb_acl.acl.id}"
  established_timeout       = 600
}

resource "alicloud_slb_listener" "udp" {
  load_balancer_id          = "${alicloud_slb.instance.id}"
  backend_port              = 2001
  frontend_port             = 2001
  protocol                  = "udp"
  bandwidth                 = 10
  persistence_timeout       = 3600
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 4
  health_check_timeout      = 8
  health_check_connect_port = 20
  acl_status                = "on"
  acl_type                  = "white"
  acl_id                    = "${alicloud_slb_acl.acl.id}"
}

resource "alicloud_slb_listener" "http" {
  load_balancer_id          = "${alicloud_slb.instance.id}"
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  sticky_session            = "on"
  sticky_session_type       = "insert"
  cookie                    = "testslblistenercookie"
  cookie_timeout            = 86400
  health_check              = "on"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx,http_3xx"
  bandwidth                 = 10
  acl_status                = "on"
  acl_type                  = "black"
  acl_id                    = "${alicloud_slb_acl.acl.id}"
}

resource "alicloud_slb_acl" "acl" {
  name       = "tf-testAccSlbAcl"
  ip_version = "ipv4"

  entry_list = [
    {
      entry   = "10.10.10.0/24"
      comment = "first"
    },
    {
      entry   = "168.10.10.0/24"
      comment = "second"
    },
    {
      entry   = "172.10.10.0/24"
      comment = "third"
    },
  ]
}

data "alicloud_slb_acls" "slb_acls" {
  ids         = ["${alicloud_slb_acl.acl.id}"]
  output_file = "${path.module}/acl.json"
}
