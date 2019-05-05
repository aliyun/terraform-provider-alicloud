variable "name" {
  default = "tf-testAccRamRolesDataSourceForPolicy-%d"
}
resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects",
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"]
    }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_role" "role" {
  name = "${var.name}"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}

data "alicloud_ram_roles" "role" {
  policy_name = "${alicloud_ram_role_policy_attachment.attach.policy_name}"
  policy_type = "Custom"
}