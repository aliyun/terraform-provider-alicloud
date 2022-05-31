variable "name" {
  default = "tftest"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "Iscsi"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}


resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}

resource "alicloud_cloud_storage_gateway_gateway_block_volume" "default" {
  cache_mode                = "Cache"
  chap_enabled              = true
  chap_in_user              = var.name
  chap_in_password          = var.name
  chunk_size                = "8192"
  gateway_block_volume_name = var.name
  gateway_id                = alicloud_cloud_storage_gateway_gateway.default.id
  local_path                = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_path
  oss_bucket_name           = alicloud_oss_bucket.default.bucket
  oss_bucket_ssl            = true
  oss_endpoint              = alicloud_oss_bucket.default.extranet_endpoint
  protocol                  = "iSCSI"
  size                      = 100
}
