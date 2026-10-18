package cms

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the CMS service.
//
// alicloud_cms_alert_robot and its data source are Framework-native: declared
// here, registered through the framework provider's Resources()/DataSources()
// aggregation, and never added to the SDK v2 provider.go map.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cms",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_cms_alert_robot",
				Factory:  NewAlertRobotResource,
			},
		},
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_cms_alert_robots",
				Factory:  NewAlertRobotsDataSource,
			},
		},
	}
}

// stringPointer returns a pointer to s. A thin local alias keeps the CMS
// package free of the SDK v2 helper import.
func stringPointer(s string) *string { return &s }
