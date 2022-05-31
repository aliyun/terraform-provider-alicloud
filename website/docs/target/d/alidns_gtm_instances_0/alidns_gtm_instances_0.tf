data "alicloud_alidns_gtm_instances" "ids" {}
output "alidns_gtm_instance_id_1" {
  value = data.alicloud_alidns_gtm_instances.ids.instances.0.id
}
