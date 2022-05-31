data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_bastionhost_instance" "default" {
  description        = "Terraform-test"
  license_code       = "bhah_ent_50_asset"
  period             = 1
  security_group_ids = [alicloud_security_group.default.0.id, alicloud_security_group.default.1.id]
  vswitch_id         = "v-testVswitch"
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  ad_auth_server {
    server         = "192.168.1.1"
    standby_server = "192.168.1.3"
    port           = "80"
    domain         = "domain"
    account        = "cn=Manager,dc=test,dc=com"
    password       = "YouPassword123"
    filter         = "objectClass=person"
    name_mapping   = "nameAttr"
    email_mapping  = "emailAttr"
    mobile_mapping = "mobileAttr"
    is_ssl         = false
    base_dn        = "dc=test,dc=com"
  }
  ldap_auth_server {
    server             = "192.168.1.1"
    standby_server     = "192.168.1.3"
    port               = "80"
    login_name_mapping = "uid"
    account            = "cn=Manager,dc=test,dc=com"
    password           = "YouPassword123"
    filter             = "objectClass=person"
    name_mapping       = "nameAttr"
    email_mapping      = "emailAttr"
    mobile_mapping     = "mobileAttr"
    is_ssl             = false
    base_dn            = "dc=test,dc=com"
  }

  lifecycle {
    ignore_changes = [ldap_auth_server.0.password, ad_auth_server.0.password]
  }
}
