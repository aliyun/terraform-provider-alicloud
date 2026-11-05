package file

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage declares the resources and data sources owned by this package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "file",
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_file_crc64_checksum", Factory: dataSourceAlicloudFileCRC64Checksum},
		},
	}
}
