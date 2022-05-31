data "alicloud_amqp_bindings" "examples" {
  instance_id       = "amqp-cn-xxxxx"
  virtual_host_name = "my-vh"
}
