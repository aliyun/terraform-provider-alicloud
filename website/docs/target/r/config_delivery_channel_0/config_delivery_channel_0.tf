provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_mns_topic" "example" {
  name = "test-topic"
}

# Example for create a MNS delivery channel
resource "alicloud_config_delivery_channel" "example" {
  description                      = "channel_description"
  delivery_channel_name            = "channel_name"
  delivery_channel_assume_role_arn = "acs:ram::11827252********:role/aliyunserviceroleforconfig"
  delivery_channel_type            = "MNS"
  delivery_channel_target_arn      = format("acs:oss:cn-shanghai:11827252********:/topics/%s", alicloud_mns_topic.example.name)
  delivery_channel_condition       = <<EOF
  [
      {
          "filterType":"ResourceType",
          "values":[
              "ACS::CEN::CenInstance",
              "ACS::CEN::Flowlog",
          ],
          "multiple":true
      }
  ]
    EOF
}
