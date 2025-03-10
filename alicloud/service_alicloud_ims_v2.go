package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ImsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeImsOidcProvider <<< Encapsulated get interface for Ims OidcProvider.

func (s *ImsServiceV2) DescribeImsOidcProvider(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["OIDCProviderName"] = id

	action := "GetOIDCProvider"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.OIDCProvider"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("OidcProvider", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.OIDCProvider", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.OIDCProvider", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ImsServiceV2) ImsOidcProviderStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeImsOidcProvider(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeImsOidcProvider >>> Encapsulated.
