data "alicloud_brain_industrial_pid_projects" "example" {
  ids        = ["3e74e684-cbb5-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_brain_industrial_pid_project_id" {
  value = data.alicloud_brain_industrial_pid_projects.example.projects.0.id
}
