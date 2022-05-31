resource "alicloud_ecs_prefix_list" "default" {
  address_family   = "IPv4"
  max_entries      = 2
  prefix_list_name = "tftest"
  description      = "description"
  entry {
    cidr        = "192.168.0.0/24"
    description = "description"
  }
}
