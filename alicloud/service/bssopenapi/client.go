package bssopenapi

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func NewBssOpenApiService(client *connectivity.AliyunClient) *BssOpenApiService {
	return &BssOpenApiService{client: client}
}
