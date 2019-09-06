resource "alicloud_nas_access_group" "main" {
  name        = "tf-testAccNasConfigName"
  type        = "Classic"
  description = "tf-testAccNasConfigDescription"
}

