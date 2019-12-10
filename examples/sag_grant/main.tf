provider "alicloud" {
  alias = "sag_account"
}
provider "alicloud" {
  region     = "cn-shanghai"
  access_key = "xxx"
  secret_key = "xxx"
  alias      = "ccn_account"
}
resource "alicloud_cloud_connect_network" "ccn" {
  provider   = "alicloud.ccn_account"
  name       = "tf-testAccCloudConnectNetwork-xxx"
  is_default = "true"
}
resource "alicloud_sag_grant" "default" {
  provider   = "alicloud.sag_account"
  sag_id     = "tf-testAccSagGrant-xxx"
  ccn_id     = "${alicloud_cloud_connect_network.ccn.id}"
  ccn_uid    = "xxx"
  depends_on = ["alicloud_cloud_connect_network.ccn"]
}