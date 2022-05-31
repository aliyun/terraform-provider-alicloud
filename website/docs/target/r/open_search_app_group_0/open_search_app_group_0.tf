variable "name" {
  default = "name"
}
resource "alicloud_open_search_app_group" "default" {
  app_group_name = var.name
  payment_type   = "PayAsYouGo"
  type           = "standard"
  quota {
    doc_size         = 1
    compute_resource = 20
    spec             = "opensearch.share.common"
  }
}

