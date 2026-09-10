resource "alicloud_slb_acl" "foo" {
  name       = "tf-testAccSlbAcl"
  ip_version = "ipv4"

  entry_list {
    entry   = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
    entry   = "168.10.10.0/24"
    comment = "second"
  }
  entry_list {
    entry   = "172.10.10.0/24"
    comment = "third"
  }
}

data "alicloud_slb_acls" "slb_acls" {
  ids = [alicloud_slb_acl.foo.id]
}

