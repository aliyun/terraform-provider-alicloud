resource "alicloud_havip" "foo" {
  vswitch_id  = "vsw-fakeid"
  description = "test_havip"
}
