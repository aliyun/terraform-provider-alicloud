# Add a Resource Manager role.
resource "alicloud_resource_manager_role" "example" {
  role_name                   = "testrd"
  assume_role_policy_document = <<EOF
     {
          "Statement": [
               {
                    "Action": "sts:AssumeRole",
                    "Effect": "Allow",
                    "Principal": {
                        "RAM":[
                                "acs:ram::103755469187****:root"，
                                "acs:ram::104408977069****:root"
                        ]
                    }
                }
          ],
          "Version": "1"
     }
	 EOF
}
