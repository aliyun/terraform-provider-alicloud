data "alicloud_brain_industrial_pid_organizations" "example" {
  ids        = ["3e74e684-cbb5-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_brain_industrial_pid_organization_id" {
  value = data.alicloud_brain_industrial_pid_organizations.example.organizations.0.id
}
