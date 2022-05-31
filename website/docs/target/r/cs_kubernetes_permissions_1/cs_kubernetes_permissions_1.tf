# Create a new RAM user.
resource "alicloud_ram_user" "user" {
  name         = var.name
  display_name = var.name
  mobile       = "86-18688888888"    # replace to your tel
  email        = "hello.uuu@aaa.com" # replace to your email
  comments     = "yoyoyo"
  force        = true
}

# Create a new RAM Policy.
resource "alicloud_ram_policy" "policy" {
  policy_name     = var.name
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "cs:Get*",
          "cs:List*",
          "cs:Describe*"
        ],
        "Effect": "Allow",
        "Resource": [
          "acs:cs:*:*:cluster/${alicloud_cs_managed_kubernetes.default.0.id}"
        ]
      }
    ],
    "Version": "1"
  }
  EOF
  description     = "this is a policy test by tf"
  force           = true
}

# Authorize the RAM user
resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  user_name   = alicloud_ram_user.user.name
}
