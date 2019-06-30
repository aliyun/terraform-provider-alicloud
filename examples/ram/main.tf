resource "alicloud_ram_user" "user" {
  name         = "${var.user_name}"
  display_name = "${var.display_name}"
  mobile       = "${var.mobile}"
  email        = "${var.email}"
  comments     = "yoyoyo"
  force        = true
}

resource "alicloud_ram_login_profile" "profile" {
  user_name = "${alicloud_ram_user.user.name}"
  password  = "${var.password}"
}

resource "alicloud_ram_access_key" "ak" {
  user_name   = "${alicloud_ram_user.user.name}"
  status      = "Active"
  secret_file = "/Users/yu/accesskey.txt"
}

resource "alicloud_ram_group" "group" {
  name     = "${var.group_name}"
  comments = "this is a group comments."
  force    = true
}

resource "alicloud_ram_group_membership" "membership" {
  group_name = "${alicloud_ram_group.group.name}"

  user_names = [
    "${alicloud_ram_user.user.name}",
  ]
}

resource "alicloud_ram_role" "role" {
  name = "${var.role_name}"

  services = [
    "apigateway.aliyuncs.com",
    "ecs.aliyuncs.com",
  ]

  ram_users = [
    "acs:ram::${var.account_id}:root",
  ]

  //    "acs:ram::${var.member_account_id}:user/username"]
  description = "this is a role test."
  force       = true
}

resource "alicloud_ram_policy" "policy" {
  name = "${var.policy_name}"

  statement = [
    {
      effect = "Deny"

      action = [
        "oss:ListObjects",
        "oss:GetObject",
      ]

      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*",
      ]
    },
  ]

  description = "this is a policy test"
  force       = true
}

resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  user_name   = "${alicloud_ram_user.user.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  group_name  = "${alicloud_ram_group.group.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name   = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}

resource "alicloud_ram_account_alias" "alias" {
  account_alias = "hallo"
}
