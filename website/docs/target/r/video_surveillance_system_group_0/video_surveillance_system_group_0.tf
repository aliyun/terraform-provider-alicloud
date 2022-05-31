resource "alicloud_video_surveillance_system_group" "default" {
  group_name   = "your_group_name"
  in_protocol  = "rtmp"
  out_protocol = "flv"
  play_domain  = "your_plan_domain"
  push_domain  = "your_push_domain"
}
