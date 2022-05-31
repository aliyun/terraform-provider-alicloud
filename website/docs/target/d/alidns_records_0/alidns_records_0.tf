data "alicloud_alidns_records" "records_ds" {
  domain_name = "xiaozhu.top"
  ids         = ["1978593525779****"]
  type        = "A"
  output_file = "records.txt"
}

output "first_record_id" {
  value = "${data.alicloud_alidns_records.records_ds.records.0.record_id}"
}
