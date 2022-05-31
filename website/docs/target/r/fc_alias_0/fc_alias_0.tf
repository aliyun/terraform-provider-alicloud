resource "alicloud_fc_alias" "example" {
  alias_name      = "my_alias"
  description     = "a sample description"
  service_name    = "my_service_name"
  service_version = "1"

  routing_config {
    additional_version_weights = {
      "2" = 0.5
    }
  }
}
