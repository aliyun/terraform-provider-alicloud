provider "alicloud" {
  alias = "ccn_account"
}

provider "alicloud" {
  region     = "cn-hangzhou"
  access_key = "xxxxxx"
  secret_key = "xxxxxx"
  alias      = "cen_account"
}

resource "alicloud_cen_instance" "cen" {
  provider = "alicloud.cen_account"
  name     = "tf-testAccCenInstance-xxx"
}

resource "alicloud_cloud_connect_network" "ccn" {
  provider   = "alicloud.ccn_account"
  name       = "tf-testAccCloudConnectNetwork-xxx"
  is_default = "true"
}

resource "alicloud_cloud_connect_network_grant" "default" {
  ccn_id  = alicloud_cloud_connect_network.ccn.id
  cen_id  = alicloud_cen_instance.cen.id
  cen_uid = "xxxxxx"
  depends_on = [
    "alicloud_cloud_connect_network.ccn",
  "alicloud_cen_instance.cen"]
}