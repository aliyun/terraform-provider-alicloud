resource "alicloud_sae_load_balancer_intranet" "example" {
  app_id          = "your_application_id"
  intranet_slb_id = "intranet_slb_id"
  intranet {
    protocol    = "TCP"
    port        = 80
    target_port = 8080
  }
}

