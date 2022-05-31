# Create a new action trail.
resource "alicloud_actiontrail" "foo" {
  name            = "action-trail"
  event_rw        = "Write-test"
  oss_bucket_name = alicloud_oss_bucket.bucket.id
  role_name       = alicloud_ram_role_policy_attachment.attach.role_name
  oss_key_prefix  = "at-product-account-audit-B"
}
