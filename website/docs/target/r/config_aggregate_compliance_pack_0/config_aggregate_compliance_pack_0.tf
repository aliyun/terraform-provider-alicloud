variable "name" {
  default = "example_name"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_instances" "default" {}

resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = "140278452670****"
    account_name = "test-2"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccaggregator"
  description     = "tf-testaccaggregator"
}


resource "alicloud_config_aggregate_config_rule" "default" {
  aggregator_id              = alicloud_config_aggregator.default.id
  aggregate_config_rule_name = var.name
  source_owner               = "ALIYUN"
  source_identifier          = "ecs-cpu-min-count-limit"
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  resource_types_scope       = ["ACS::ECS::Instance"]
  risk_level                 = 1
  description                = var.name
  exclude_resource_ids_scope = data.alicloud_instances.default.ids.0
  input_parameters = {
    cpuCount = "4",
  }
  region_ids_scope         = "cn-hangzhou"
  resource_group_ids_scope = data.alicloud_resource_manager_resource_groups.default.ids.0
  tag_key_scope            = "tFTest"
  tag_value_scope          = "forTF 123"
}

resource "alicloud_config_aggregate_compliance_pack" "default" {
  aggregate_compliance_pack_name = "tf-testaccConfig1234"
  aggregator_id                  = alicloud_config_aggregator.default.id
  description                    = "tf-testaccConfig1234"
  risk_level                     = 1
  config_rule_ids {
    config_rule_id = alicloud_config_aggregate_config_rule.default.config_rule_id
  }
}

