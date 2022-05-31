# Create a RAM user
resource "alicloud_ram_user" "example" {
  name = "tf-testaccramuser"
}

# Create a Resource Manager Policy
resource "alicloud_resource_manager_policy" "example" {
  policy_name     = "tf-testaccrdpolicy"
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

# Create a Resource Group
resource "alicloud_resource_manager_resource_group" "example" {
  display_name = "tf_test"
  name         = "tf_test"
}

# Get Alicloud Account Id
data "alicloud_account" "example" {}

# Attach the custom policy to resource group
resource "alicloud_resource_manager_policy_attachment" "example" {
  policy_name       = alicloud_resource_manager_policy.example.policy_name
  policy_type       = "Custom"
  principal_name    = format("%s@%s.onaliyun.com", alicloud_ram_user.example.name, data.alicloud_account.example.id)
  principal_type    = "IMSUser"
  resource_group_id = alicloud_resource_manager_resource_group.example.id
}
