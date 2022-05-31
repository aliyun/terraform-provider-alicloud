resource "alicloud_bastionhost_instance" "default" {
  description        = "Terraform-test"
  license_code       = "bhah_ent_50_asset"
  period             = "1"
  vswitch_id         = "v-testVswitch"
  security_group_ids = ["sg-test", "sg-12345"]
}
