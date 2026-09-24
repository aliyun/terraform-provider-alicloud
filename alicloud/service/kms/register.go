package kms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the KMS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Kms",
		EphemeralResources: []conns.EphemeralResource{
			{
				TypeName: "alicloud_kms_secret",
				Factory:  NewSecretEphemeralResource,
			},
			{
				TypeName: "alicloud_kms_plaintext",
				Factory:  NewPlaintextEphemeralResource,
			},
		},
	}
}
