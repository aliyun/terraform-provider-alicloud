resource "alicloud_event_bridge_event_bus" "example" {
  event_bus_name = "example_value"
}

resource "alicloud_event_bridge_rule" "example" {
  event_bus_name = alicloud_event_bridge_event_bus.example.id
  rule_name      = var.name
  description    = "test"
  filter_pattern = "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}"
  targets {
    target_id = "tf-test"
    endpoint  = "acs:mns:cn-hangzhou:118938335****:queues/tf-test"
    type      = "acs.mns.queue"
    param_list {
      resource_key = "queue"
      form         = "CONSTANT"
      value        = "tf-testaccEbRule"
    }
    param_list {
      resource_key = "Body"
      form         = "ORIGINAL"
    }
  }
}

