# Grant users developer permissions for the cluster.
resource "alicloud_cs_kubernetes_permissions" "default" {
  # uid
  uid = alicloud_ram_user.user.id
  # permissions
  permissions {
    cluster     = alicloud_cs_managed_kubernetes.default.0.id
    role_type   = "cluster"
    role_name   = "dev"
    namespace   = ""
    is_custom   = false
    is_ram_role = false
  }
  # If you want to grant users multiple cluster permissions, you can define multiple sets of permissions 
  #  permissions {
  #    cluster     = "cluster_id_2"
  #    role_type   = "cluster"
  #    role_name   = "ops"
  #    namespace   =  ""
  #    is_custom   = false
  #    is_ram_role = false
  #  }
  depends_on = [
    alicloud_ram_user_policy_attachment.attach
  ]
}
