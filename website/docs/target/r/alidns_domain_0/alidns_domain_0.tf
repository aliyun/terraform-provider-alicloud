# Add a new Domain.
resource "alicloud_alidns_domain" "dns" {
  domain_name = "starmove.com"
  group_id    = "85ab8713-4a30-4de4-9d20-155ff830****"
  tags = {
    Created     = "Terraform"
    Environment = "test"
  }
}
