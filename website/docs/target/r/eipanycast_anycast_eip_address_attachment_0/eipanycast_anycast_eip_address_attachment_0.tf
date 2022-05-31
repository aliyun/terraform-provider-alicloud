resource "alicloud_eipanycast_anycast_eip_address" "example" {
  service_location = "international"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "example" {
  anycast_id              = alicloud_eipanycast_anycast_eip_address.example.id
  bind_instance_id        = "lb-j6chlcr8lffy7********"
  bind_instance_region_id = "cn-hongkong"
  bind_instance_type      = "SlbInstance"
}

