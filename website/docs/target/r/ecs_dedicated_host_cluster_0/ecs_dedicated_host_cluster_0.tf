data "alicloud_zones" example {}

resource "alicloud_ecs_dedicated_host_cluster" "example" {
  dedicated_host_cluster_name = "example_value"
  description                 = "example_value"
  zone_id                     = data.alicloud_zones.example.zones.0.id
  tags = {
    Create = "TF"
    For    = "DDH_Cluster_Test",
  }
}

