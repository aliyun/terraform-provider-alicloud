resource "alicloud_alidns_instance" "example" {
  dns_security   = "no"
  domain_numbers = "2"
  period         = 1
  renew_period   = 1
  renewal_status = "ManualRenewal"
  version_code   = "version_personal"
}

