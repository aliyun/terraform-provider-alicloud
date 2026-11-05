package rkvstore

// Local copies of cross-product service methods and shared enums/constants that the
// R-kvstore sub-package needs. The root package service types (KmsService, VpcService,
// BssOpenApiService in alicloud/service_alicloud_*.go) and root enums (PayType in
// alicloud/common.go, BACKUP_TIME in alicloud/extension_rds.go, ModulesSizeLimit in
// alicloud/service_alicloud_bss_open_api.go) are retained unchanged for their other
// consumers; this sub-package cannot import the root package, so the necessary methods
// are reproduced 1:1 here (qualified leaf-package calls), per stock-method-disposition.

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// ModulesSizeLimit mirrors alicloud/service_alicloud_bss_open_api.go.
const ModulesSizeLimit = 50

// PayType and its enum values mirror alicloud/common.go.
type PayType string

const (
	PrePaid    = PayType("PrePaid")
	PostPaid   = PayType("PostPaid")
	Prepaid    = PayType("Prepaid")
	Postpaid   = PayType("Postpaid")
	Serverless = PayType("Serverless")
)

// BACKUP_TIME mirrors alicloud/extension_rds.go (cross-product shared backup window list).
var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}

// KVStoreInstanceType and its enum values mirror alicloud/extension_rds.go.
type KVStoreInstanceType string

const (
	KVStoreRedis    = KVStoreInstanceType("Redis")
	KVStoreMemcache = KVStoreInstanceType("Memcache")
)

// KmsService is a local copy of the root KmsService providing the Decrypt method.
type KmsService struct {
	client *connectivity.AliyunClient
}

// Decrypt mirrors alicloud/service_alicloud_kms.go KmsService.Decrypt.
func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, errs.WrapError(err)
	}

	var response map[string]interface{}
	client := s.client
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return plaintext, errs.WrapErrorf(err, errs.DefaultErrorMsg, context, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

// VpcService is a local copy of the root VpcService providing the DescribeVSwitch method.
type VpcService struct {
	client *connectivity.AliyunClient
}

// DescribeVSwitch mirrors alicloud/service_alicloud_vpc.go VpcService.DescribeVSwitch.
func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

// BssOpenApiService is a local copy of the root BssOpenApiService providing GetInstanceTypePrice.
type BssOpenApiService struct {
	client *connectivity.AliyunClient
}

// GetInstanceTypePrice mirrors alicloud/service_alicloud_bss_open_api.go BssOpenApiService.GetInstanceTypePrice.
func (s *BssOpenApiService) GetInstanceTypePrice(productCode, productType, paymentType string, modules []map[string]interface{}) (object []float64, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	var priceList []float64
	action := "GetSubscriptionPrice"
	request := map[string]interface{}{
		"OrderType":        "NewOrder",
		"SubscriptionType": "Subscription",
		"Region":           client.RegionId,
		"ProductCode":      productCode,
		"ProductType":      productType,
	}

	if paymentType == "Subscription" {
		request["ServicePeriodQuantity"] = 1
		request["ServicePeriodUnit"] = "Month"
		request["Quantity"] = 1
	}

	if paymentType == "PayAsYouGo" {
		action = "GetPayAsYouGoPrice"
		request["SubscriptionType"] = "PayAsYouGo"
	}

	moduleLength := len(modules)

	for {
		if len(modules) < ModulesSizeLimit {
			moduleLength = len(modules)

			request["ModuleList"] = modules
		} else {
			tmp := modules[:ModulesSizeLimit]
			modules = modules[ModulesSizeLimit:]
			moduleLength = len(tmp)

			request["ModuleList"] = &tmp
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return object, errs.WrapError(err)
		}

		v, err := jsonpath.Get("$.Data.ModuleDetails.ModuleDetail", response)
		if err != nil {
			return object, errs.WrapError(err)
		}

		for _, originalCosts := range v.([]interface{}) {
			originalCostArg := originalCosts.(map[string]interface{})

			if originalCost, ok := originalCostArg["OriginalCost"]; ok {
				priceListItem, err := convert.FormatFloat64(originalCost)
				if err != nil {
					return nil, err
				}
				priceList = append(priceList, priceListItem)
			}
		}

		if moduleLength < ModulesSizeLimit {
			break
		}
	}

	return priceList, nil
}
