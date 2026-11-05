package ga

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the flat
// layout's acceptance tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeGaAccelerator(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccelerator"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Accelerator", "UnknownError"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:Accelerator", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) GaAcceleratorStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAccelerator(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaListener(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeListener"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ListenerId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Listener", "UnknownError"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:GaListener", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) GaListenerStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaListener(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, object["State"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaBandwidthPackage(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBandwidthPackage"
	request := map[string]interface{}{
		"RegionId":           s.client.RegionId,
		"BandwidthPackageId": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.BandwidthPackage"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BandwidthPackage", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) GaBandwidthPackageStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBandwidthPackage(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaEndpointGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeEndpointGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"EndpointGroupId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("Ga:EndpointGroup", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:EndpointGroup", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DescribeGaForwardingRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListForwardingRules"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		err = errs.WrapError(err)
		return
	}

	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"ClientToken":      helper.BuildClientToken("ListForwardingRules"),
		"AcceleratorId":    parts[0],
		"ListenerId":       parts[1],
		"ForwardingRuleId": parts[2],
		"MaxResults":       helper.PageSizeLarge,
	}

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"StateError.Accelerator"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.ForwardingRules", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ForwardingRules", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:ForwardingRule", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ListenerId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["ForwardingRuleId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:ForwardingRule", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DescribeGaIpSet(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"IpSetId":  id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeIpSet"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"UnknownError", "NotExist.IpSet"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:IpSet", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:IpSet", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DescribeGaBandwidthPackageAttachment(id string) (object map[string]interface{}, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	object, err = s.DescribeGaBandwidthPackage(fmt.Sprint(parts[1]))
	if err != nil {
		return object, errs.WrapError(err)
	}

	idExist := false
	if accelerators, ok := object["Accelerators"]; ok {
		acceleratorList := accelerators.([]interface{})
		for _, accelerator := range acceleratorList {
			if accelerator == fmt.Sprint(parts[0]) {
				idExist = true
				return object, nil
			}
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BandwidthPackageAttachment", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(object["RequestId"]))
	}

	return object, nil
}

func (s *Service) GaBandwidthPackageAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBandwidthPackageAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) GaEndpointGroupStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaEndpointGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *Service) GaForwardingRuleStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaForwardingRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["ForwardingRuleStatus"].(string) == failState {
				return object, object["ForwardingRuleStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["ForwardingRuleStatus"].(string)))
			}
		}

		return object, object["ForwardingRuleStatus"].(string), nil
	}
}

func (s *Service) GaIpSetStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaIpSet(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeAcceleratorAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAcceleratorAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeGaAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAcl"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:Acl", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) GaAclStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAcl(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AclStatus"]) == failState {
				return object, fmt.Sprint(object["AclStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["AclStatus"])))
			}
		}

		return object, fmt.Sprint(object["AclStatus"]), nil
	}
}

func (s *Service) GetAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAcl"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:Acl", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeGaAclAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeListener"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ListenerId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Listener"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga::AclAttachment", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	resp, _ := v.(map[string]interface{})
	if v, ok := resp["RelatedAcls"].([]interface{}); ok && len(v) > 0 {
		for _, aclArgs := range v {
			aclArg := aclArgs.(map[string]interface{})
			if fmt.Sprint(aclArg["AclId"]) == parts[1] {
				idExist = true
				resp["status"] = aclArg["Status"]
				return resp, nil
			}
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga::AclAttachment", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *Service) GaAclAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAclAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["status"]) == failState {
				return object, fmt.Sprint(object["status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["RelatedAcls[*].Status"])))
			}
		}
		return object, fmt.Sprint(object["status"]), nil
	}
}

func (s *Service) DescribeGaAdditionalCertificate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListListenerCertificates"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
		"Role":          "additional",
		"MaxResults":    helper.PageSizeMedium,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Certificates", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Certificates", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Domain"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *Service) DescribeGaAcceleratorSpareIpAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetSpareIp"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": parts[0],
		"SpareIp":       parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if _, ok := object["State"]; !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:SpareIp", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}
	return object, nil
}

func (s *Service) GaAcceleratorSpareIpAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAcceleratorSpareIpAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeBandwidthPackageAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBandwidthPackageAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeGaAccessLog(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeLogStoreOfEndpointGroup"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"AcceleratorId":   parts[0],
		"ListenerId":      parts[1],
		"EndpointGroupId": parts[2],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["Status"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("GaAccessLog", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaAccessLogStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAccessLog(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *Service) DescribeGaAclEntryAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAcl"

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    parts[0],
	}

	if strings.Contains(fmt.Sprint(parts[1]), "_") {
		parts[1] = strings.Replace(fmt.Sprint(parts[1]), "_", ":", -1)
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:AclEntryAttachment", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.AclEntries", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AclEntries", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:AclEntryAttachment", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Entry"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:AclEntryAttachment", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DescribeGaBasicAccelerator(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicAccelerator"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BasicAccelerator", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaBasicAcceleratorStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicAccelerator(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaBasicEndpointGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicEndpointGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClientToken":     helper.BuildClientToken("GetBasicEndpointGroup"),
		"EndpointGroupId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BasicEndpointGroup", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) GaBasicEndpointGroupStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicEndpointGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaBasicIpSet(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicIpSet"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ClientToken": helper.BuildClientToken("GetBasicIpSet"),
		"IpSetId":     id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.IpSet"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BasicIpSet", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) GaBasicIpSetStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicIpSet(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaBasicAccelerateIp(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetBasicAccelerateIp"
	client := s.client
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ClientToken":    helper.BuildClientToken("GetBasicAccelerateIp"),
		"AccelerateIpId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.AccelerateIpId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BasicAccelerateIp", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) GaBasicAccelerateIpStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicAccelerateIp(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaBasicEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetBasicEndpoint"
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ClientToken": helper.BuildClientToken("GetBasicEndpoint"),
		"EndpointId":  parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.EndPoints"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BasicEndpoint", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) GaBasicEndpointStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicEndpoint(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaBasicAccelerateIpEndpointRelation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicAccelerateIpEndpointRelation"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ClientToken":    helper.BuildClientToken("GetBasicAccelerateIpEndpointRelation"),
		"AcceleratorId":  parts[0],
		"AccelerateIpId": parts[1],
		"EndpointId":     parts[2],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:BasicAccelerateIpEndpointRelation", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaBasicAccelerateIpEndpointRelationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicAccelerateIpEndpointRelation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"AcceleratorId": parts[0],
		"Domain":        parts[1],
		"RegionId":      s.client.RegionId,
	}

	var response map[string]interface{}
	action := "ListDomains"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Domains", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Domains", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga", id), errs.NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DescribeGaCustomRoutingEndpointGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndpointGroup"

	client := s.client

	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"EndpointGroupId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpointGroup", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpointGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaCustomRoutingEndpointGroupStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpointGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaCustomRoutingEndpointGroupDestination(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndpointGroupDestinations"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"DestinationId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Destination"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpointGroupDestination", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpointGroupDestination", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaCustomRoutingEndpointGroupDestinationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpointGroupDestination(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaCustomRoutingEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndpoint"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"EndpointGroup": parts[0],
		"EndpointId":    parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpoint", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpoint", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaCustomRoutingEndpointStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpoint(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) DescribeGaCustomRoutingEndpointTrafficPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndPointTrafficPolicy"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": parts[0],
		"PolicyId":   parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"NotExist.Policy"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpointTrafficPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ga:CustomRoutingEndpointTrafficPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *Service) GaCustomRoutingEndpointTrafficPolicyStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpointTrafficPolicy(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *Service) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}

	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}

			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}

			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil || response["NextToken"].(string) == "" {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	client := s.client
	resourceIdNum := strings.Count(d.Id(), ":")

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		var response map[string]interface{}
		var err error
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}

		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := helper.ParseResourceId(d.Id(), 2)
				if err != nil {
					return errs.WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := helper.ParseResourceId(d.Id(), 2)
				if err != nil {
					return errs.WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}
