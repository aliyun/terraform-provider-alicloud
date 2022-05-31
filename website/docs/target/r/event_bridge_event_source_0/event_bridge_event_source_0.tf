resource "alicloud_event_bridge_event_source" "example" {
  event_bus_name         = "bus_name"
  event_source_name      = "tftest"
  description            = "tf-test"
  linked_external_source = true
  external_source_type   = "MNS"
  external_source_config = {
    QueueName = "mns_queuqe_name"
  }
}

