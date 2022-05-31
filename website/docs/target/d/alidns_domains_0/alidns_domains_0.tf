data "alicloud_alidns_domains" "domains_ds" {
  domain_name_regex = "^hegu"
  output_file       = "domains.txt"
}

output "first_domain_id" {
  value = "${data.alicloud_alidns_domains.domains_ds.domains.0.domain_id}"
}
