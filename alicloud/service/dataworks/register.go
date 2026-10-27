package dataworks

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "dataworks", Resources: []conns.Resource{{TypeName: "alicloud_data_works_certificate", Factory: NewCertificateResource}}}
}
