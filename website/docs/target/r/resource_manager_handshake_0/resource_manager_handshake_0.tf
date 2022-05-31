# Add a Resource Manager handshake.
resource "alicloud_resource_manager_handshake" "example" {
  target_entity = "1182775234******"
  target_type   = "Account"
  note          = "test resource manager handshake"
}
