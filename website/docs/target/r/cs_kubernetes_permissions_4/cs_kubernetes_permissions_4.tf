# remove the permissions on the "cluster_id_01", "cluster_id_02".
resource "alicloud_cs_kubernetes_permissions" "default" {
  uid = data.alicloud_ram_users.users_ds.users.0.id
}
