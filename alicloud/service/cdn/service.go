// Package cdn. The service helpers in this file were migrated verbatim from the
// former root-package helpers service_alicloud_cdn.go and service_alicloud_cdn_v2.go.
// Both the typed CdnService (v1, using the cdn SDK models) and the map-based
// CdnServiceV2 (v2, using RpcPost/RpcGet with map[string]interface{}) are kept as
// distinct types because different resources depend on either layer; their method
// bodies, retry, timeout and error semantics are unchanged. Root-package helper
// symbols are referenced through their leaf packages (errs, helper, convert, logs,
// tags) so this package does not reverse-import the root alicloud package.
package cdn

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CdnService wraps the typed cdn SDK calls used by the SDKv2 resources of this package.
type CdnService struct {
	client *connectivity.AliyunClient
}

// CdnServiceV2 wraps the map-based RpcPost/RpcGet calls used by the SDKv2 resources
// of this package that operate on raw response maps.
type CdnServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewCdnService is the exported constructor for CdnService: the client field stays
// unexported, so callers outside the package (including the acceptance tests) build
// values with it. Mirrors the ecs.NewService / oss.NewService convention.
func NewCdnService(client *connectivity.AliyunClient) *CdnService {
	return &CdnService{client: client}
}

// NewCdnServiceV2 is the exported constructor for CdnServiceV2.
func NewCdnServiceV2(client *connectivity.AliyunClient) *CdnServiceV2 {
	return &CdnServiceV2{client: client}
}

func (c *CdnService) convertCdnSourcesToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})

		priorityVal, err := convert.FormatInt(item["priority"])
		if err != nil {
			return "", err
		}

		weightVal, err := convert.FormatInt(item["weight"])
		if err != nil {
			return "", err
		}
		arrayMaps[i] = map[string]interface{}{
			"content":  item["content"],
			"port":     item["port"],
			"priority": priorityVal,
			"type":     item["type"],
			"weight":   weightVal,
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", errs.WrapError(err)
	}
	return string(maps), nil
}

func (c *CdnService) DescribeCdnDomainNew(id string) (*cdn.GetDomainDetailModel, error) {
	model := &cdn.GetDomainDetailModel{}
	request := cdn.CreateDescribeCdnDomainDetailRequest()
	request.RegionId = c.client.RegionId
	request.DomainName = id

	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeCdnDomainDetail(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound", "ConfigNotFound"}) {
			return model, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return model, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	domain, _ := raw.(*cdn.DescribeCdnDomainDetailResponse)
	if domain.GetDomainDetailModel.DomainName != id {
		return model, errs.WrapErrorf(errs.NotFoundErr("cdn_domain", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return &domain.GetDomainDetailModel, nil
}

func (s *CdnService) DescribeCdnDomainConfig(id string) (object interface{}, err error) {

	var response map[string]interface{}
	client := s.client
	action := "DescribeCdnDomainConfigs"

	parts := strings.Split(id, ":")
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"DomainName":    parts[0],
		"FunctionNames": parts[1],
	}

	if len(parts) > 2 {
		request["ConfigId"] = parts[2]
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"FlowControlError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	v, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("cdn_domain_config", id), errs.DefaultErrorMsg, err)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("cdn_domain_config", id), errs.ResourceNotfound, response)
	}

	val := v.([]interface{})[0].(map[string]interface{})
	return val, nil
}

func (c *CdnService) WaitForCdnDomain(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	time.Sleep(helper.DefaultIntervalShort * time.Second)

	for {
		domain, err := c.DescribeCdnDomainNew(id)
		if err != nil {
			if errs.NotFoundError(err) && status == Deleted {
				break
			}
			return errs.WrapError(err)
		}
		if domain.DomainStatus == string(status) {
			break
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, domain.DomainStatus, status, errs.ProviderERROR)
		}
	}
	return nil
}

func (c *CdnService) DescribeDomainCertificateInfo(id string) (certInfo cdn.CertInfo, err error) {
	request := cdn.CreateDescribeDomainCertificateInfoRequest()
	request.RegionId = c.client.RegionId
	request.DomainName = id
	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeDomainCertificateInfo(request)
	})
	if err != nil {
		return certInfo, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cdn.DescribeDomainCertificateInfoResponse)
	if len(response.CertInfos.CertInfo) <= 0 {
		return certInfo, errs.WrapErrorf(errs.NotFoundErr("DomainCertificateInfo", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	certInfo = response.CertInfos.CertInfo[0]
	return
}

func (c *CdnService) WaitForServerCertificateNew(id string, serverCertificate string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		certInfo, err := c.DescribeDomainCertificateInfo(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if strings.TrimSpace(certInfo.ServerCertificate) == strings.TrimSpace(serverCertificate) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, strings.TrimSpace(certInfo.ServerCertificate), strings.TrimSpace(serverCertificate), errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

// DescribeTags keeps the named return tagItems (rather than tags) so it does not
// shadow the imported tags package inside this method; the return name is not part
// of the call signature, so callers are unaffected.
func (c *CdnService) DescribeTags(resourceId string, resourceType TagResourceType) (tagItems []cdn.TagItem, err error) {
	request := cdn.CreateDescribeTagResourcesRequest()
	request.RegionId = c.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := c.client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.DescribeTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cdn.DescribeTagResourcesResponse)
	if len(response.TagResources) < 1 {
		return
	}
	for _, t := range response.TagResources {
		tagItems = append(tagItems, t.Tag...)
	}
	return
}

func (c *CdnService) CdnDomainConfigRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := c.DescribeCdnDomainConfig(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		val := v.(map[string]interface{})
		for _, failState := range failStates {
			if fmt.Sprint(val["Status"]) == failState {
				return val, fmt.Sprint(val["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(val["Status"])))
			}
		}
		return val, fmt.Sprint(val["Status"]), nil
	}
}

func (s *CdnService) DescribeCdnRealTimeLogDelivery(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDomainRealtimeLogDelivery"
	request := map[string]interface{}{
		"Domain": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Cdn", "2018-05-10", action, request, nil)
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
		if errs.IsExpectedErrors(err, []string{"Domain.NotFound", "InternalError"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CDN:RealTimeLogDelivery", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *CdnService) CdnRealTimeLogDeliveryStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCdnRealTimeLogDelivery(id)
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

func (s *CdnService) DescribeCdnFcTrigger(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeFCTrigger"
	request := map[string]interface{}{
		"TriggerARN": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Cdn", "2018-05-10", action, request, nil)
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
	v, err := jsonpath.Get("$.FCTrigger", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.FCTrigger", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

// --- CdnServiceV2: map-based helpers migrated from service_alicloud_cdn_v2.go ---

// DescribeCdnDomain <<< Encapsulated get interface for Cdn Domain.

func (s *CdnServiceV2) DescribeCdnDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeCdnDomainDetail"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DomainName"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.GetDomainDetailModel", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GetDomainDetailModel", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CdnServiceV2) DescribeDomainDescribeDomainCertificateInfo(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDomainCertificateInfo"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DomainName"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.CertInfos.CertInfo[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CertInfos.CertInfo[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *CdnServiceV2) DescribeDomainListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListTagResources"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = id
	request["RegionId"] = client.RegionId

	request["ResourceType"] = "DOMAIN"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CdnServiceV2) CdnDomainStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCdnDomain(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

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

// DescribeCdnDomain >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Cdn.
func (s *CdnServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, false)
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
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, false)
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

// SetResourceTags >>> tag function encapsulated.

// DescribeCdnRealTimeLogDelivery <<< Encapsulated get interface for Cdn RealTimeLogDelivery.

func (s *CdnServiceV2) DescribeCdnRealTimeLogDelivery(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDomainRealtimeLogDelivery"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Domain"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Cdn", "2018-05-10", action, query, nil)

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
		if errs.IsExpectedErrors(err, []string{"Domain.NotFound", "InternalError"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("RealTimeLogDelivery", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CdnServiceV2) CdnRealTimeLogDeliveryStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCdnRealTimeLogDelivery(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

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

// DescribeCdnRealTimeLogDelivery >>> Encapsulated.
