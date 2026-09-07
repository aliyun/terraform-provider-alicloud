resource "alicloud_mns_topic" "topic" {
  name                 = var.name
  maximum_message_size = var.maximum_message_size
  logging_enabled      = var.loggin_enabled
}

resource "alicloud_mns_topic_subscription" "subscription" {
  topic_name            = alicloud_mns_topic.topic.name
  name                  = var.subscription_name
  endpoint              = var.endpoint
  filter_tag            = var.filter_tag
  notify_strategy       = var.notify_strategy
  notify_content_format = var.notify_content_format
}

