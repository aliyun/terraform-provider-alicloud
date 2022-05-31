resource "alicloud_ecs_key_pair_attachment" "example" {
  key_pair_name = "key_pair_name"
  instance_ids  = [i-gw80pxxxxxxxxxx]
}

