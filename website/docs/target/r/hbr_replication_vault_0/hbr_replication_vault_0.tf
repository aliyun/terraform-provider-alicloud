variable "name" {
  default = "tf-testAccReplicationVault"
}

variable "region_source" {
  default = "you Replication value source region"
}

provider "alicloud" {
  alias  = "source"
  region = var.region_source
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  provider   = alicloud.source
}

data "alicloud_hbr_replication_vault_regions" "default" {}

locals {
  region_replication = data.alicloud_hbr_replication_vault_regions.default.regions.0.replication_region_id
}

provider "alicloud" {
  alias  = "replication"
  region = local.region_replication
}

resource "alicloud_hbr_replication_vault" "default" {
  replication_source_region_id = local.region_replication
  replication_source_vault_id  = alicloud_hbr_vault.default.id
  vault_name                   = var.name
  vault_storage_class          = "STANDARD"
  description                  = var.name
  provider                     = alicloud.replication
}
