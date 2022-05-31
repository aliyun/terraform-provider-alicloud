resource "alicloud_config_aggregator" "example" {
  aggregator_accounts {
    account_id   = "123968452689****"
    account_name = "tf-testacc1234"
    account_type = "ResourceDirectory"
  }
  aggregator_name = "tf-testaccConfigAggregator1234"
  description     = "tf-testaccConfigAggregator1234"
}

