variable "name" {
  default = "example_name"
}

data "alicloud_instances" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_config_rule" "default" {
  rule_name                  = var.name
  description                = var.name
  source_identifier          = "ecs-instances-in-vpc"
  source_owner               = "ALIYUN"
  resource_types_scope       = ["ACS::ECS::Instance"]
  risk_level                 = 1
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  tag_key_scope              = "tfTest"
  tag_value_scope            = "tfTest 123"
  resource_group_ids_scope   = data.alicloud_resource_manager_resource_groups.default.ids.0
  exclude_resource_ids_scope = data.alicloud_instances.default.instances[0].id
  region_ids_scope           = "cn-hangzhou"
  input_parameters = {
    vpcIds = data.alicloud_instances.default.instances[0].vpc_id
  }
}

resource "alicloud_config_compliance_pack" "default" {
  compliance_pack_name = "tf-testaccConfig1234"
  description          = "tf-testaccConfig1234"
  risk_level           = "1"
  config_rule_ids {
    config_rule_id = alicloud_config_rule.default.id
  }
}

