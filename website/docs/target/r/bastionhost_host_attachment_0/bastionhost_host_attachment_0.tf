resource "alicloud_bastionhost_host_attachment" "example" {
  host_group_id = "6"
  host_id       = "15"
  instance_id   = "bastionhost-cn-tl32bh0no30"
}

