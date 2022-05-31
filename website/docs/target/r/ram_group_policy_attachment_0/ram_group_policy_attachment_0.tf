# Create a RAM Group Policy attachment.
resource "alicloud_ram_group" "group" {
  name     = "groupName"
  comments = "this is a group comments."
  force    = true
}

resource "alicloud_ram_policy" "policy" {
  name        = "policyName"
  document    = <<EOF
    {
      "Statement": [
        {
          "Action": [
            "oss:ListObjects",
            "oss:GetObject"
          ],
          "Effect": "Allow",
          "Resource": [
            "acs:oss:*:*:mybucket",
            "acs:oss:*:*:mybucket/*"
          ]
        }
      ],
        "Version": "1"
    }
  EOF
  description = "this is a policy test"
  force       = true
}

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  group_name  = alicloud_ram_group.group.name
}
