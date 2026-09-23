resource "alicloud_cloud_connect_network" "ccn" {
  name       = "tf-testAccCloudConnectNetworkAttachment-xxx"
  is_default = "true"
}

resource "alicloud_cloud_connect_network_attachment" "default" {
  ccn_id     = alicloud_cloud_connect_network.ccn.id
  sag_id     = "sag-xxxxx"
  depends_on = ["alicloud_cloud_connect_network.ccn"]
}