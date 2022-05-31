resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
  ratio          = 30
}

