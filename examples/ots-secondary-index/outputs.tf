output "table_name" {
  value = alicloud_ots_table.table.table_name
}

output "index_name" {
  value = alicloud_ots_secondary_index.index1.index_name
}

output "index_type" {
  value = alicloud_ots_secondary_index.index1.index_type
}

output "index_primary_keys" {
  value = alicloud_ots_secondary_index.index1.primary_keys
}

output "index_defined_columns" {
  value = alicloud_ots_secondary_index.index1.defined_columns
}



