package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type GaService struct {
	client *connectivity.AliyunClient
}

func (s *GaService) DescribeGaAccelerator(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccelerator"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Accelerator", "UnknownError"}) {
			return object, WrapErrorf(NotFoundErr("Ga:Accelerator", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *GaService) GaAcceleratorStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAccelerator(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaListener(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeListener"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ListenerId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Listener", "UnknownError"}) {
			return object, WrapErrorf(NotFoundErr("Ga:GaListener", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *GaService) GaListenerStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaListener(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaBandwidthPackage(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBandwidthPackage"
	request := map[string]interface{}{
		"RegionId":           s.client.RegionId,
		"BandwidthPackageId": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.BandwidthPackage"}) {
			return object, WrapErrorf(NotFoundErr("Ga:BandwidthPackage", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *GaService) GaBandwidthPackageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBandwidthPackage(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaEndpointGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeEndpointGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"EndpointGroupId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return nil, WrapErrorf(NotFoundErr("Ga:EndpointGroup", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:EndpointGroup", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *GaService) DescribeGaForwardingRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListForwardingRules"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}

	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"ClientToken":      buildClientToken("ListForwardingRules"),
		"AcceleratorId":    parts[0],
		"ListenerId":       parts[1],
		"ForwardingRuleId": parts[2],
		"MaxResults":       PageSizeLarge,
	}

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.ForwardingRules", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ForwardingRules", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(NotFoundErr("Ga:ForwardingRule", id), NotFoundWithResponse, response)
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
		return object, WrapErrorf(NotFoundErr("Ga:ForwardingRule", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *GaService) DescribeGaIpSet(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"IpSetId":  id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeIpSet"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"UnknownError", "NotExist.IpSet"}) {
			return object, WrapErrorf(NotFoundErr("Ga:IpSet", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:IpSet", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *GaService) DescribeGaBandwidthPackageAttachment(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	object, err = s.DescribeGaBandwidthPackage(fmt.Sprint(parts[1]))
	if err != nil {
		return object, WrapError(err)
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
		return object, WrapErrorf(NotFoundErr("Ga:BandwidthPackageAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(object["RequestId"]))
	}

	return object, nil
}

func (s *GaService) GaBandwidthPackageAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBandwidthPackageAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) GaEndpointGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaEndpointGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["State"].(string) == failState {
				return object, object["State"].(string), WrapError(Error(FailedToReachTargetStatus, object["State"].(string)))
			}
		}
		return object, object["State"].(string), nil
	}
}

func (s *GaService) GaForwardingRuleStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaForwardingRule(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ForwardingRuleStatus"].(string) == failState {
				return object, object["ForwardingRuleStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ForwardingRuleStatus"].(string)))
			}
		}

		return object, object["ForwardingRuleStatus"].(string), nil
	}
}

func (s *GaService) GaIpSetStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaIpSet(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeAcceleratorAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAcceleratorAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *GaService) DescribeGaAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAcl"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return object, WrapErrorf(NotFoundErr("Ga:Acl", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *GaService) GaAclStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAcl(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AclStatus"]) == failState {
				return object, fmt.Sprint(object["AclStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AclStatus"])))
			}
		}

		return object, fmt.Sprint(object["AclStatus"]), nil
	}
}

func (s *GaService) GetAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAcl"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return object, WrapErrorf(NotFoundErr("Ga:Acl", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *GaService) DescribeGaAclAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeListener"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ListenerId": parts[0],
	}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Listener"}) {
			return object, WrapErrorf(NotFoundErr("Ga::AclAttachment", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
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
		return object, WrapErrorf(NotFoundErr("Ga::AclAttachment", id), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *GaService) GaAclAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAclAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["status"]) == failState {
				return object, fmt.Sprint(object["status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["RelatedAcls[*].Status"])))
			}
		}
		return object, fmt.Sprint(object["status"]), nil
	}
}

func (s *GaService) DescribeGaAdditionalCertificate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListListenerCertificates"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
		"Role":          "additional",
		"MaxResults":    PageSizeMedium,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Certificates", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Certificates", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("Ga", id), NotFoundWithResponse, response)
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
		return object, WrapErrorf(NotFoundErr("Ga", id), NotFoundWithResponse, response)
	}
	return
}

func (s *GaService) DescribeGaAcceleratorSpareIpAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetSpareIp"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": parts[0],
		"SpareIp":       parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if _, ok := object["State"]; !ok {
		return object, WrapErrorf(NotFoundErr("Ga:SpareIp", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}
	return object, nil
}

func (s *GaService) GaAcceleratorSpareIpAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAcceleratorSpareIpAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeBandwidthPackageAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBandwidthPackageAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *GaService) DescribeGaAccessLog(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeLogStoreOfEndpointGroup"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"AcceleratorId":   parts[0],
		"ListenerId":      parts[1],
		"EndpointGroupId": parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["Status"] == nil {
		return object, WrapErrorf(NotFoundErr("GaAccessLog", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaAccessLogStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaAccessLog(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *GaService) DescribeGaAclEntryAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAcl"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    parts[0],
	}

	if strings.Contains(fmt.Sprint(parts[1]), "_") {
		parts[1] = strings.Replace(fmt.Sprint(parts[1]), "_", ":", -1)
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return object, WrapErrorf(NotFoundErr("Ga:AclEntryAttachment", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.AclEntries", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AclEntries", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("Ga:AclEntryAttachment", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Entry"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("Ga:AclEntryAttachment", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *GaService) DescribeGaBasicAccelerator(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicAccelerator"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"AcceleratorId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:BasicAccelerator", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaBasicAcceleratorStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicAccelerator(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaBasicEndpointGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicEndpointGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ClientToken":     buildClientToken("GetBasicEndpointGroup"),
		"EndpointGroupId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return object, WrapErrorf(NotFoundErr("Ga:BasicEndpointGroup", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *GaService) GaBasicEndpointGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicEndpointGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaBasicIpSet(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicIpSet"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ClientToken": buildClientToken("GetBasicIpSet"),
		"IpSetId":     id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.IpSet"}) {
			return object, WrapErrorf(NotFoundErr("Ga:BasicIpSet", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *GaService) GaBasicIpSetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicIpSet(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaBasicAccelerateIp(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetBasicAccelerateIp"
	client := s.client
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ClientToken":    buildClientToken("GetBasicAccelerateIp"),
		"AccelerateIpId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.AccelerateIpId"}) {
			return object, WrapErrorf(NotFoundErr("Ga:BasicAccelerateIp", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *GaService) GaBasicAccelerateIpStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicAccelerateIp(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaBasicEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetBasicEndpoint"
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ClientToken": buildClientToken("GetBasicEndpoint"),
		"EndpointId":  parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.EndPoints"}) {
			return object, WrapErrorf(NotFoundErr("Ga:BasicEndpoint", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *GaService) GaBasicEndpointStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicEndpoint(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaBasicAccelerateIpEndpointRelation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetBasicAccelerateIpEndpointRelation"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"ClientToken":    buildClientToken("GetBasicAccelerateIpEndpointRelation"),
		"AcceleratorId":  parts[0],
		"AccelerateIpId": parts[1],
		"EndpointId":     parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:BasicAccelerateIpEndpointRelation", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaBasicAccelerateIpEndpointRelationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaBasicAccelerateIpEndpointRelation(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}
		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"AcceleratorId": parts[0],
		"Domain":        parts[1],
		"RegionId":      s.client.RegionId,
	}

	var response map[string]interface{}
	action := "ListDomains"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Domains", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Domains", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Ga", id), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *GaService) DescribeGaCustomRoutingEndpointGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndpointGroup"

	client := s.client

	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"EndpointGroupId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpointGroup", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpointGroup", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaCustomRoutingEndpointGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpointGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaCustomRoutingEndpointGroupDestination(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndpointGroupDestinations"

	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"DestinationId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Destination"}) {
			return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpointGroupDestination", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpointGroupDestination", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaCustomRoutingEndpointGroupDestinationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpointGroupDestination(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaCustomRoutingEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndpoint"

	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"EndpointGroup": parts[0],
		"EndpointId":    parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) {
			return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpoint", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpoint", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaCustomRoutingEndpointStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpoint(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) DescribeGaCustomRoutingEndpointTrafficPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCustomRoutingEndPointTrafficPolicy"

	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"EndpointId": parts[0],
		"PolicyId":   parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NotExist.Policy"}) {
			return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpointTrafficPolicy", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if object["State"] == nil {
		return object, WrapErrorf(NotFoundErr("Ga:CustomRoutingEndpointTrafficPolicy", id), NotFoundMsg, ProviderERROR)
	}

	return object, nil
}

func (s *GaService) GaCustomRoutingEndpointTrafficPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGaCustomRoutingEndpointTrafficPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["State"]) == failState {
				return object, fmt.Sprint(object["State"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["State"])))
			}
		}

		return object, fmt.Sprint(object["State"]), nil
	}
}

func (s *GaService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}

			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}

			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil || response["NextToken"].(string) == "" {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *GaService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	client := s.client
	resourceIdNum := strings.Count(d.Id(), ":")

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		var response map[string]interface{}
		var err error
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, false)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, false)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}
