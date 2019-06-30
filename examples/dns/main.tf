data "alicloud_dns_domains" "domain" {
  domain_name_regex = "^hegu"
}

data "alicloud_dns_groups" "group" {
  name_regex = "^y[A-Za-z]+"
}

data "alicloud_dns_records" "record" {
  domain_name       = "${data.alicloud_dns_domains.domain.domains.0.domain_name}"
  is_locked         = false
  type              = "A"
  host_record_regex = "^@"
  output_file       = "records.txt"
}

resource "alicloud_dns_group" "group" {
  name  = "${var.group_name}"
  count = "${var.number}"
}

resource "alicloud_dns" "dns" {
  name     = "${var.domain_name}"
  group_id = "${element(alicloud_dns_group.group.*.id, count.index)}"
}

resource "alicloud_dns_record" "record" {
  name        = "${alicloud_dns.dns.name}"
  host_record = "alimailskajdh"
  type        = "CNAME"
  value       = "mail.mxhichind.com"
}
