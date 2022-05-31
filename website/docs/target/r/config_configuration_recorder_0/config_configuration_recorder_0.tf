resource "alicloud_config_configuration_recorder" "example" {
  resource_types = [
    "ACS::ECS::Instance",
    "ACS::ECS::Disk"
    # other resource types ...
  ]
}
