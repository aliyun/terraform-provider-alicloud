# Create a cen vbr HealrhCheck resource and use it.
resource "alicloud_cen_instance" "default" {
  cen_instance_name = "test_name"
}

resource "alicloud_cen_instance_attachment" "default" {
  instance_id              = alicloud_cen_instance.default.id
  child_instance_id        = "vbr-xxxxx"
  child_instance_type      = "VBR"
  child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_cen_vbr_health_check" "default" {
  cen_id                 = alicloud_cen_instance.default.id
  health_check_source_ip = "192.168.1.2"
  health_check_target_ip = "10.0.0.2"
  vbr_instance_id        = "vbr-xxxxx"
  vbr_instance_region_id = "cn-hangzhou"
  health_check_interval  = 2
  healthy_threshold      = 8
  depends_on             = [alicloud_cen_instance_attachment.default]
}
