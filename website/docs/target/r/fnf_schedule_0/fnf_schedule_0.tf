resource "alicloud_fnf_flow" "example" {
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: pass
      name: helloworld
  EOF  
  description = "tf-testaccFnFFlow983041"
  name        = "tf-testAccSchedule"
  type        = "FDL"
}

resource "alicloud_fnf_schedule" "example" {
  cron_expression = "30 9 * * * *"
  description     = "tf-testaccFnFSchedule983041"
  enable          = "true"
  flow_name       = alicloud_fnf_flow.example.name
  payload         = "{\"tf-test\": \"test success\"}"
  schedule_name   = "tf-testaccFnFSchedule983041"
}
