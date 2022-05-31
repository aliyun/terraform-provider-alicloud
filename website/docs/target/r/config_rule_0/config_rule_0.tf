# Audit ECS instances under VPC using preset rules
resource "alicloud_config_rule" "example" {
  rule_name            = "instances-in-vpc"
  source_identifier    = "ecs-instances-in-vpc"
  source_owner         = "ALIYUN"
  resource_types_scope = ["ACS::ECS::Instance"]
  description          = "ecs instances in vpc"
  input_parameters = {
    vpcIds = "vpc-uf6gksw4ctjd******"
  }
  risk_level                = 1
  config_rule_trigger_types = "ConfigurationItemChangeNotification"
}

