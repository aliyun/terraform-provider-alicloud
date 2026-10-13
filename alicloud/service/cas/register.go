package cas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the CAS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cas",
		Actions: []conns.Action{
			{
				TypeName: "alicloud_ssl_certificates_service_pca_cert_sync",
				Factory:  NewPcaCertSyncAction,
			},
		},
	}
}
