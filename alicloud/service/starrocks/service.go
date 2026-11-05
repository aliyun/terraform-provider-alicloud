// Package starrocks. This file is generated automatically. Please do not modify it manually, thank you!
package starrocks

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

// Service wraps the StarRocks API calls used by the resources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeInstance <<< Encapsulated get interface for StarRocks Instance.

func (s *Service) DescribeInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["InstanceId"] = helper.StringPointer(id)
	query["RegionId"] = helper.StringPointer(client.RegionId)
	action := fmt.Sprintf("/webapi/starrocks/describeInstances")

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)

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

	v, err := jsonpath.Get("$.Data[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["InstanceStatus"]
	if fmt.Sprint(currentStatus) == "deleting" {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) InstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstance(id)
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

// DescribeInstance >>> Encapsulated.
// SetResourceTags <<< Encapsulated tag function for StarRocks.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/webapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceId"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ExpandSingletonToList(d.Id())))
			query["RegionId"] = helper.StringPointer(client.RegionId)
			query["ResourceType"] = helper.StringPointer(resourceType)

			query["TagKey"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(removedTagKeys)))
			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("starrocks", "2022-10-19", action, query, nil, nil, true)
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
			action = fmt.Sprintf("/webapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			request["ResourceType"] = resourceType
			request["RegionId"] = client.RegionId
			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["Key"] = key
				tagsMap["Value"] = value
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			request["Tag"] = tagsMaps

			jsonString := convert.ConvertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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

// DescribeNodeGroup <<< Encapsulated get interface for StarRocks NodeGroup.

func (s *Service) DescribeNodeGroup(id string) (map[string]interface{}, error) {
	return s.DescribeNodeGroupWithTimeout(id, time.Minute)
}

// DescribeNodeGroupWithTimeout reads one node group within the caller's budget.
func (s *Service) DescribeNodeGroupWithTimeout(id string, timeout time.Duration) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["ClusterId"] = helper.StringPointer(parts[0])
	query["RegionId"] = helper.StringPointer(client.RegionId)
	jsonString := convert.ConvertObjectToJsonString(query)
	jsonString, _ = sjson.Set(jsonString, "nodeGroupIds.0", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)
	body := request

	action := fmt.Sprintf("/webapi/nodegroup/describeNodeGroups")

	budget := timeout
	if budget <= 0 {
		budget = time.Minute
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(budget, func() *retry.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)

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

	v, err := jsonpath.Get("$.Data[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("NodeGroup", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["Status"]
	if fmt.Sprint(currentStatus) == "DELETED" {
		return object, errs.WrapErrorf(errs.NotFoundErr("NodeGroup", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) NodeGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.NodeGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeNodeGroup)
}

func (s *Service) NodeGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeNodeGroup >>> Encapsulated.

// DescribeBackendNodeGroups lists the backend compute groups of an instance. A BE
// or a CN group is a backend group depending on the run mode, and the instance
// resource does not persist node group ids, so both component types are read and
// returned together for attribute based matching.
func (s *Service) DescribeBackendNodeGroups(instanceId string, timeout time.Duration) ([]map[string]interface{}, error) {
	groups := make([]map[string]interface{}, 0)
	for _, componentType := range []string{"BE", "CN"} {
		items, err := s.DescribeNodeGroupsByComponentType(instanceId, componentType, timeout)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			if fmt.Sprint(item["Status"]) == "DELETED" {
				continue
			}
			groups = append(groups, item)
		}
	}
	return groups, nil
}

// DescribeNodeGroupsByComponentType lists every compute group of one component type.
// The API paginates and a group on a later page is as real as one on the first page,
// so the complete result set is read before a caller decides which group to use.
func (s *Service) DescribeNodeGroupsByComponentType(clusterId, componentType string, timeout time.Duration) ([]map[string]interface{}, error) {
	client := s.client
	action := fmt.Sprintf("/webapi/nodegroup/describeNodeGroups")
	const pageSize = 100
	const maxPages = 100
	items := make([]map[string]interface{}, 0)
	deadline := time.Now().Add(timeout)
	if timeout <= 0 {
		deadline = time.Now().Add(time.Minute)
	}
	for pageNumber := 1; ; pageNumber++ {
		if pageNumber > maxPages {
			return nil, errs.WrapErrorf(fmt.Errorf("more than %d pages of node groups", maxPages), errs.DefaultErrorMsg, clusterId, action, errs.AlibabaCloudSdkGoERROR)
		}
		query := make(map[string]*string)
		query["ClusterId"] = helper.StringPointer(clusterId)
		query["RegionId"] = helper.StringPointer(client.RegionId)
		query["PageNumber"] = helper.StringPointer(strconv.Itoa(pageNumber))
		query["PageSize"] = helper.StringPointer(strconv.Itoa(pageSize))
		request := map[string]interface{}{"componentType": componentType}

		// Every page consumes the same budget: a slow first page must not hand a full
		// fresh timeout to the next one.
		budget := time.Until(deadline)
		if budget <= 0 {
			return nil, errs.WrapErrorf(fmt.Errorf("listing node groups did not finish within the configured timeout"), errs.DefaultErrorMsg, clusterId, action, errs.AlibabaCloudSdkGoERROR)
		}
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(budget, func() *retry.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, request, true)
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
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, action, errs.AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.Data[*]", response)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, clusterId, "$.Data[*]", response)
		}
		batch := make([]map[string]interface{}, 0)
		for _, raw := range v.([]interface{}) {
			if item, ok := raw.(map[string]interface{}); ok {
				batch = append(batch, item)
			}
		}
		items = append(items, batch...)

		if len(batch) == 0 {
			break
		}
		if len(batch) == pageSize {
			// A full page always means more rows may follow. Stopping here on Total alone
			// would be unsafe: Total may or may not count the deleted groups that this
			// listing still returns, and a missed group can look unique to the matcher.
			continue
		}
		total, totalErr := nodeGroupsTotal(response)
		if totalErr != nil || total <= 0 || len(items) >= total {
			break
		}
		// The page is shorter than requested but the reported total says rows remain
		// (server side page size cap): keep paging.
	}
	return items, nil
}

func nodeGroupsTotal(response map[string]interface{}) (int, error) {
	v, err := jsonpath.Get("$.Total", response)
	if err != nil {
		return 0, err
	}
	total, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil || total < 0 {
		return 0, fmt.Errorf("invalid node group total %v", v)
	}
	return total, nil
}
