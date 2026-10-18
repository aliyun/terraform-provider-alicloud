// Package threatdetection holds the Framework-native declarations of the Threat
// Detection (SOAR) service. It is registered through the provider registry so the
// muxed server serves it from the Framework half; existing SDK v2 Threat Detection
// resources and data sources stay in package alicloud.
package threatdetection

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Threat Detection
// service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ThreatDetection",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_threat_detection_components",
				Factory:  NewComponentsDataSource,
			},
		},
	}
}
