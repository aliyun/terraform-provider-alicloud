package cr

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tfTags "github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CrServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewCrServiceV2 constructs a CrServiceV2 from an AliyunClient. It is the
// supported entry point for callers outside the cr package (notably external tests).
func NewCrServiceV2(client *connectivity.AliyunClient) *CrServiceV2 {
	return &CrServiceV2{client: client}
}

// DescribeCrInstance <<< Encapsulated get interface for Cr Instance.

func (s *CrServiceV2) DescribeCrInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId
	action := "GetInstance"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"InvalidAction.NotFound"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"INSTANCE_NOT_EXIST"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return response, nil
}
func (s *CrServiceV2) DescribeInstanceQueryAvailableInstances(d *schema.ResourceData) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = id
	request["Region"] = client.RegionId
	var endpoint string
	request["ProductCode"] = "acr"
	request["ProductType"] = "acr_ee_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "acr_ee_public_intl"
	}
	action := "QueryAvailableInstances"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "acr"
				request["ProductType"] = "acr_ee_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *CrServiceV2) DescribeInstanceListInstanceEndpoint(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId
	action := "ListInstanceEndpoint"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

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

	return response, nil
}

func (s *CrServiceV2) CrInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CrInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeCrInstance)
}

func (s *CrServiceV2) CrInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

func (s *CrServiceV2) DescribeAsyncCrInstanceStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncGetInstance(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeCrInstance >>> Encapsulated.
// DescribeCrRepoSyncRule <<< Encapsulated get interface for Cr RepoSyncRule.

func (s *CrServiceV2) DescribeCrRepoSyncRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	action := "ListRepoSyncRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["NamespaceName"] = parts[1]
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNo"] = 1

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

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

		if fmt.Sprint(response["IsSuccess"]) == "false" {
			return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		resp, err := jsonpath.Get("$.SyncRules", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SyncRules", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cr:RepoSyncRule", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["LocalInstanceId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["LocalNamespaceName"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["SyncRuleId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNo"] = request["PageNo"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cr:RepoSyncRule", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

// DescribeAsyncGetInstance <<< Encapsulated for Cr.
func (s *CrServiceV2) DescribeAsyncGetInstance(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	action := "GetInstance"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncGetInstance >>> Encapsulated.

// DescribeCrScanRule <<< Encapsulated get interface for Cr ScanRule.

func (s *CrServiceV2) DescribeCrScanRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["ScanRuleId"] = parts[1]
	request["RegionId"] = client.RegionId
	action := "GetScanRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"SCAN_RULE_NOT_EXIST", "INSTANCE_NOT_EXIST"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ScanRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ScanRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScanRule", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *CrServiceV2) CrScanRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CrScanRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeCrScanRule)
}

func (s *CrServiceV2) CrScanRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeCrScanRule >>> Encapsulated.

// DescribeCrStorageDomainRoutingRule <<< Encapsulated get interface for Cr StorageDomainRoutingRule.

func (s *CrServiceV2) DescribeCrStorageDomainRoutingRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleId"] = parts[1]
	request["RegionId"] = client.RegionId
	action := "GetStorageDomainRoutingRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST", "InvalidInstanceStorageDomainRoutingRule.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("StorageDomainRoutingRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	currentStatus := response["RuleId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("StorageDomainRoutingRule", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *CrServiceV2) CrStorageDomainRoutingRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CrStorageDomainRoutingRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeCrStorageDomainRoutingRule)
}

func (s *CrServiceV2) CrStorageDomainRoutingRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeCrStorageDomainRoutingRule >>> Encapsulated.

// DescribeCrEERepo <<< Encapsulated get interface for Cr EE Repo.

func (s *CrServiceV2) DescribeCrEERepo(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"InstanceId":        parts[0],
		"RepoNamespaceName": parts[1],
		"RepoName":          parts[2],
	}
	query := make(map[string]interface{})
	action := "GetRepository"
	var response map[string]interface{}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CrEE:Repo", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) == "REPO_NOT_EXIST" {
		return object, errs.WrapErrorf(errs.NotFoundErr("CrEE:Repo", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	if isSuccess, ok := response["IsSuccess"].(bool); ok && !isSuccess {
		return object, errs.WrapErrorf(fmt.Errorf("%v", response), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeCrEERepo >>> Encapsulated.

// DescribeCrArtifactLifecycleRule <<< Encapsulated get interface for Cr ArtifactLifecycleRule.

func (s *CrServiceV2) DescribeCrArtifactLifecycleRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RuleId"] = parts[1]
	query["InstanceId"] = parts[0]
	query["RegionId"] = client.RegionId
	action := "GetArtifactLifecycleRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("cr", "2018-12-01", action, query, request)

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
		if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ArtifactLifecycleRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CrServiceV2) CrArtifactLifecycleRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CrArtifactLifecycleRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeCrArtifactLifecycleRule)
}

func (s *CrServiceV2) CrArtifactLifecycleRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeCrArtifactLifecycleRule >>> Encapsulated.

// ListCrInstanceTags <<< Encapsulated list tags interface for Cr Instance.

func (s *CrServiceV2) ListCrInstanceTags(id string) (tags map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "ListTagResources"
	tags = make(map[string]interface{})
	nextToken := ""

	for {
		request := make(map[string]interface{})
		query := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["ResourceType"] = "Instance"
		request["ResourceId.1"] = id
		if nextToken != "" {
			request["NextToken"] = nextToken
		}

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
			return tags, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		tagResourceRaw, _ := jsonpath.Get("$.TagResources.TagResource", response)
		if tagResourceRaw != nil {
			for key, value := range tfTags.TagsToMap(tagResourceRaw) {
				tags[key] = value
			}
		}

		if v, ok := response["NextToken"]; ok && fmt.Sprint(v) != "" {
			nextToken = fmt.Sprint(v)
		} else {
			break
		}
	}

	return tags, nil
}

// ListCrInstanceTags >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Cr.

func (s *CrServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tfTags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tfTags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
// DescribeCrArtifactSubscriptionRule <<< Encapsulated get interface for Cr ArtifactSubscriptionRule.

func (s *CrServiceV2) DescribeCrArtifactSubscriptionRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["RuleId"] = parts[1]
	query["InstanceId"] = parts[0]
	query["RegionId"] = client.RegionId
	action := "GetArtifactSubscriptionRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("cr", "2018-12-01", action, query, request)

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
		if errs.IsExpectedErrors(err, []string{"INSTANCE_NOT_EXIST"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ArtifactSubscriptionRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CrServiceV2) CrArtifactSubscriptionRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CrArtifactSubscriptionRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeCrArtifactSubscriptionRule)
}

func (s *CrServiceV2) CrArtifactSubscriptionRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeCrArtifactSubscriptionRule >>> Encapsulated.

// DescribeCrInstanceCustomizedDomain <<< Encapsulated get interface for Cr InstanceCustomizedDomain.

func (s *CrServiceV2) DescribeCrInstanceCustomizedDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": parts[0],
		"ModuleName": parts[1],
		"Domain":     parts[2],
	}
	query := make(map[string]interface{})
	action := "GetInstanceCustomizedDomain"
	var response map[string]interface{}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"DOMAIN_NOT_EXIST", "INSTANCE_NOT_EXIST"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("InstanceCustomizedDomain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CrServiceV2) CrInstanceCustomizedDomainStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.CrInstanceCustomizedDomainStateRefreshFuncWithApi(id, field, failStates, s.DescribeCrInstanceCustomizedDomain)
}

func (s *CrServiceV2) CrInstanceCustomizedDomainStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeCrInstanceCustomizedDomain >>> Encapsulated.
