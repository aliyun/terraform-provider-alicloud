// register.go declares the Framework-native RocketMQ service package. The
// legacy Plugin SDK RocketMQ types remain in package alicloud; new types are
// added here on Framework.
package rocketmq

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the RocketMQ
// (2022-08-01) service for the Framework provider server.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Rocketmq",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_rocketmq_disaster_recovery_plan",
				Factory:  NewDisasterRecoveryPlanResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_rocketmq_disaster_recovery_plans",
				Factory:  NewDisasterRecoveryPlansDataSource,
			},
		},
	}
}
