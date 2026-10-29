package dataworks

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
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

// NewService is the exported constructor for the DataWorks service.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeDataWorksProject <<< Encapsulated get interface for DataWorks Project.

func (s *Service) DescribeDataWorksProject(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetProject"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"9990020002", "9990040003"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"1101080008", "9990040003"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Project", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) DataWorksProjectStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksProject(id)
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

// DescribeDataWorksProject >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for DataWorks.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
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
				response, err = client.RpcPost("dataworks-public", "2020-05-18", action, query, request, false)
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
				response, err = client.RpcPost("dataworks-public", "2020-05-18", action, query, request, false)
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

// DescribeDataWorksProjectMember <<< Encapsulated get interface for DataWorks ProjectMember.

func (s *Service) DescribeDataWorksProjectMember(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetProjectMember"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ProjectId"] = parts[0]
	request["UserId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"9990020002"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// 400010 comes back with HTTP 403 once the queried UserId no longer belongs to the
		// tenant of the workspace, which means the project member does not exist any more.
		if errs.IsExpectedErrors(err, []string{"1101080166", "400010"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ProjectMember", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ProjectMember", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ProjectMember", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DataWorksProjectMemberStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksProjectMember(id)
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

// DescribeDataWorksProjectMember >>> Encapsulated.

// DescribeDataWorksDataSource <<< Encapsulated get interface for DataWorks DataSource.

func (s *Service) DescribeDataWorksDataSource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetDataSource"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = parts[1]
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

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
		if errs.IsExpectedErrors(err, []string{"1000010005"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DataSource", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DataSource", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DataSource", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DataWorksDataSourceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksDataSource(id)
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

// DescribeDataWorksDataSource >>> Encapsulated.

// DescribeDataWorksDataSourceSharedRule <<< Encapsulated get interface for DataWorks DataSourceSharedRule.

func (s *Service) DescribeDataWorksDataSourceSharedRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "ListDataSourceSharedRules"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DataSourceId"] = parts[0]
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

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
		if errs.IsExpectedErrors(err, []string{"1000011038", "1000010005"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DataSourceSharedRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DataSourceSharedRules[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DataSourceSharedRules[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DataSourceSharedRule", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["DataSourceId"]) != parts[0] {
			continue
		}
		if fmt.Sprint(item["Id"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("DataSourceSharedRule", id), errs.NotFoundMsg, response)
}

func (s *Service) DataWorksDataSourceSharedRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksDataSourceSharedRule(id)
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

// DescribeDataWorksDataSourceSharedRule >>> Encapsulated.

// DescribeDataWorksDiAlarmRule <<< Encapsulated get interface for DataWorks DiAlarmRule.

func (s *Service) DescribeDataWorksDiAlarmRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "ListDIAlarmRules"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DIAlarmRuleId"] = parts[1]
	query["JobId"] = parts[0]
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

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
		if errs.IsExpectedErrors(err, []string{"500130"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DiAlarmRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.PagingInfo.DIJobAlarmRules[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PagingInfo.DIJobAlarmRules[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DiAlarmRule", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DataWorksDiAlarmRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksDiAlarmRule(id)
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

// DescribeDataWorksDiAlarmRule >>> Encapsulated.

// DescribeDataWorksDiJob <<< Encapsulated get interface for DataWorks DiJob.

func (s *Service) DescribeDataWorksDiJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetDIJob"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DIJobId"] = parts[1]
	query["ProjectId"] = parts[0]
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

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
		if errs.IsExpectedErrors(err, []string{"500130"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DiJob", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.PagingInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PagingInfo", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DataWorksDiJobStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksDiJob(id)
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

// DescribeDataWorksDiJob >>> Encapsulated.

// DescribeDataWorksNetwork <<< Encapsulated get interface for DataWorks Network.

func (s *Service) DescribeDataWorksNetwork(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetNetwork"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

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

	v, err := jsonpath.Get("$.Network", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Network", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["Id"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Network", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DataWorksNetworkStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksNetwork(id)
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

// DescribeDataWorksNetwork >>> Encapsulated.

// DescribeDataWorksDwResourceGroup <<< Encapsulated get interface for DataWorks DwResourceGroup.

func (s *Service) DescribeDataWorksDwResourceGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("dataworks-public", "2024-05-18", action, query, nil)

		if err != nil || errs.IsExpectedErrors(err, []string{"9990040003"}) {
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

	v, err := jsonpath.Get("$.ResourceGroup", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DwResourceGroup", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["Id"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("DwResourceGroup", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) DataWorksDwResourceGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDataWorksDwResourceGroup(id)
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

// DescribeDataWorksDwResourceGroup >>> Encapsulated.
