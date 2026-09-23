resource "alicloud_slb_load_balancer" "instance" {
  load_balancer_name                 = var.slb_name
  internet_charge_type = var.internet_charge_type
  address_type         = var.address_type
  load_balancer_spec        = "slb.s2.small"

  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}

resource "alicloud_slb_listener" "tcp" {
  load_balancer_id          = alicloud_slb_load_balancer.instance.id
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
  health_check_connect_port = 20
  health_check_uri          = "/console"
  acl_status                = "off"
  acl_type                  = "white"
  acl_id                    = alicloud_slb_acl.acl.id
  established_timeout       = 600
}

resource "alicloud_slb_listener" "udp" {
  load_balancer_id          = alicloud_slb_load_balancer.instance.id
  backend_port              = 2001
  frontend_port             = 2001
  protocol                  = "udp"
  bandwidth                 = 10
  persistence_timeout       = 3600
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 4
  health_check_connect_port = 20
  acl_status                = "on"
  acl_type                  = "white"
  acl_id                    = alicloud_slb_acl.acl.id
}

resource "alicloud_slb_listener" "http" {
  load_balancer_id          = alicloud_slb_load_balancer.instance.id
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
  acl_id                    = alicloud_slb_acl.acl.id
  request_timeout           = 80
  idle_timeout              = 30
}

resource "alicloud_slb_acl" "acl" {
  name       = "tf-testAccSlbAcl"
  ip_version = "ipv4"

  entry_list {
    entry   = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
    entry   = "168.10.10.0/24"
    comment = "second"
  }
  entry_list {
    entry   = "172.10.10.0/24"
    comment = "third"
  }
}

data "alicloud_slb_acls" "slb_acls" {
  ids = [alicloud_slb_acl.acl.id]
}

resource "alicloud_slb_server_certificate" "foo-file" {
  name               = "tf-testAccSlbServerCertificate-file"
  server_certificate = file("${path.module}/server_certificate.pem")
  private_key        = file("${path.module}/private_key.pem")
}

resource "alicloud_slb_listener" "https-file" {
  load_balancer_id          = alicloud_slb_load_balancer.instance.id
  backend_port              = 80
  frontend_port             = 8443
  protocol                  = "https"
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
  server_certificate_id     = alicloud_slb_server_certificate.foo-file.id
  request_timeout           = 80
  idle_timeout              = 30
  enable_http2              = "on"
  tls_cipher_policy         = "tls_cipher_policy_1_2"
}

data "alicloud_slbs" "balancers" {
  tags = {
    tag_a = 1
  }
}

