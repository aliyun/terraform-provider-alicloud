package autoscaling

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "autoscaling", Resources: []conns.Resource{{TypeName: "alicloud_auto_scaling_alarm_task", Factory: NewAlarmTaskResource}}}
}
