resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name  = alicloud_fc_service.example.name
  function_name = alicloud_fc_function.example.name

  destination_config {
    on_failure {
      destination = the_example_mns_queue_arn
    }

    on_success {
      destination = the_example_mns_topic_arn
    }
  }
}
