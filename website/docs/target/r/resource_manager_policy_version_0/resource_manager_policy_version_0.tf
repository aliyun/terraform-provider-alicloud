resource "alicloud_resource_manager_policy" "example" {
  policy_name     = "tftest"
  policy_document = <<EOF
		{
			"Statement": [{
				"Action": ["oss:*"],
				"Effect": "Allow",
				"Resource": ["acs:oss:*:*:*"]
			}],
			"Version": "1"
		}
    EOF
}

resource "alicloud_resource_manager_policy_version" "example" {
  policy_name     = alicloud_resource_manager_policy.example.policy_name
  policy_document = <<EOF
		{
			"Statement": [{
				"Action": ["oss:*"],
				"Effect": "Allow",
				"Resource": ["acs:oss:*:*:myphotos"]
			}],
			"Version": "1"
		}
    EOF
}
