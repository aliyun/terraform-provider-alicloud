# Get RAM user ID 
data "alicloud_ram_users" "users_ds" {
  name_regex = "your ram user name"
}

# Create a new RAM Policy.
resource "alicloud_ram_policy" "policy" {
  policy_name     = "AckClusterReadOnlyAccess"
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
          "acs:cs:*:*:cluster/${target_cluster_ID}"
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
  user_name   = data.alicloud_ram_users.users_ds.users.0.name
}

# RBAC authorization for the cluster
resource "alicloud_cs_kubernetes_permissions" "default" {
  uid = data.alicloud_ram_users.users_ds.users.0.id
  permissions {
    cluster     = "target cluster id1"
    role_type   = "cluster"
    role_name   = "ops"
    is_custom   = false
    is_ram_role = false
    namespace   = ""
  }
  permissions {
    cluster     = "target cluster id2"
    role_type   = "cluster"
    role_name   = "ops"
    is_custom   = false
    is_ram_role = false
    namespace   = ""
  }
  depends_on = [
    alicloud_ram_user_policy_attachment.attach
  ]
}
