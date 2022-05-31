resource "alicloud_ehpc_job_template" "default" {
  job_template_name = "example_value"
  command_line      = "./LammpsTest/lammps.pbs"
}

