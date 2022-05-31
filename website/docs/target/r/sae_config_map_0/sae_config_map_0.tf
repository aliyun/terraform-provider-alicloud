variable "ConfigMapName" {
  default = "examplename"
}
resource "alicloud_sae_config_map" "example" {
  data         = jsonencode({ "env.home" : "/root", "env.shell" : "/bin/sh" })
  name         = var.ConfigMapName
  namespace_id = alicloud_sae_namespace.example.namespace_id
}
resource "alicloud_sae_namespace" "example" {
  namespace_id          = "cn-hangzhou:yourname"
  namespace_name        = "example_value"
  namespace_description = "your_description"
}

