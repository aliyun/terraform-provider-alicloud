resource "alicloud_amqp_instance" "professional" {
  instance_type  = "professional"
  max_tps        = 1000
  queue_capacity = 50
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}

resource "alicloud_amqp_instance" "vip" {
  instance_type  = "vip"
  max_tps        = 5000
  queue_capacity = 50
  storage_size   = 700
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}
