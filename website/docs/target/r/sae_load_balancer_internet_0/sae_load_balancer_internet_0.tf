resource "alicloud_sae_load_balancer_internet" "example" {
  app_id          = "your_application_id"
  internet_slb_id = "your_internet_slb_id"
  internet {
    protocol    = "TCP"
    port        = 80
    target_port = 8080
  }
}

