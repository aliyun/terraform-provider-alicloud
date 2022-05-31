resource "alicloud_resource_manager_control_policy" "example" {
  control_policy_name = "tf-testAccRDControlPolicy"
  description         = "tf-testAccRDControlPolicy"
  effect_scope        = "RAM"
  policy_document     = <<EOF
  {
    "Version": "1",
    "Statement": [
      {
        "Effect": "Deny",
        "Action": [
          "ram:UpdateRole",
          "ram:DeleteRole",
          "ram:AttachPolicyToRole",
          "ram:DetachPolicyFromRole"
        ],
        "Resource": "acs:ram:*:*:role/ResourceDirectoryAccountAccessRole"
      }
    ]
  }
  EOF
}

