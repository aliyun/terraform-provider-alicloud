// Enable the control policy
resource "alicloud_resource_manager_resource_directory" "example" {
  status = "Enabled"
}

resource "alicloud_resource_manager_control_policy" "example" {
  control_policy_name = "tf-testAccName"
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

resource "alicloud_resource_manager_folder" "example" {
  folder_name = "tf-testAccName"
}

resource "alicloud_resource_manager_control_policy_attachment" "example" {
  policy_id  = alicloud_resource_manager_control_policy.example.id
  target_id  = alicloud_resource_manager_folder.example.id
  depends_on = [alicloud_resource_manager_resource_directory.example]
}

