resource "alicloud_cms_alarm" "basic" {
  name    = "tf-testAccCmsAlarm_basic"
  project = "acs_ecs_dashboard"
  metric  = "disk_writebytes"
  dimensions = {
    instanceId = "i-bp1247,i-bp11gd"
    device     = "/dev/vda1,/dev/vdb1"
  }
  escalations_critical {
    statistics          = "Average"
    comparison_operator = "<="
    threshold           = 35
    times               = 2
  }
  period             = 900
  contact_groups     = ["test-group"]
  effective_interval = "0:00-2:00"
  webhook            = "https://${data.alicloud_account.current.id}.eu-central-1.fc.aliyuncs.com/2016-08-15/proxy/Terraform/AlarmEndpointMock/"
}
