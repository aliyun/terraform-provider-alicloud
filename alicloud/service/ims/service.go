package ims

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service holds the IMS (2019-08-15) client and the product-specific API
// helpers used by the SDKv2 IMS resources and data sources migrated into this
// package. Framework resources in this package call the client directly via
// fwadapt; this Service is for the SDKv2 factories only.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs the IMS Service from a connectivity client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeImsOidcProvider <<< Encapsulated get interface for Ims OidcProvider.

func (s *Service) DescribeImsOidcProvider(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["OIDCProviderName"] = id

	action := "GetOIDCProvider"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.OIDCProvider"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("OidcProvider", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.OIDCProvider", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.OIDCProvider", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ImsOidcProviderStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImsOidcProvider(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.ClientIds" {
			e := jsonata.MustCompile("$split($.ClientIds, \",\")")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.Fingerprints" {
			e := jsonata.MustCompile("$split($.Fingerprints, \",\")")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeImsOidcProvider >>> Encapsulated.
