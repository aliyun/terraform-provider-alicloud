resource "alicloud_config_aggregator" "example" {
  aggregator_accounts {
    account_id   = "140278452670****"
    account_name = "test-2"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccaggregator"
  description     = "tf-testaccaggregator"
}

resource "alicloud_config_aggregate_config_rule" "example" {
  aggregate_config_rule_name = "tf-testaccconfig1234"
  aggregator_id              = alicloud_config_aggregator.example.id
  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  source_owner               = "ALIYUN"
  source_identifier          = "ecs-cpu-min-count-limit"
  risk_level                 = 1
  resource_types_scope       = ["ACS::ECS::Instance"]
  input_parameters = {
    cpuCount = "4",
  }
}

