package dmsenterprise

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func NewDMSEnterpriseServiceV2(client *connectivity.AliyunClient) *DMSEnterpriseServiceV2 {
	return &DMSEnterpriseServiceV2{client: client}
}
func NewDmsEnterpriseService(client *connectivity.AliyunClient) *DmsEnterpriseService {
	return &DmsEnterpriseService{client: client}
}
