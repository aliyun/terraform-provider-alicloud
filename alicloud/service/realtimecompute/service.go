package realtimecompute

import (
	"fmt"
	"strconv"
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

// Service holds the Realtime Compute (foasconsole 2019-06-01 / ververica
// 2022-07-18) client and the product-specific API helpers used by the SDKv2
// realtimecompute resources and data sources migrated into this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs the realtimecompute Service from a connectivity client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeRealtimeComputeVvpInstance <<< Encapsulated get interface for RealtimeCompute VvpInstance.

func (s *Service) DescribeRealtimeComputeVvpInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DescribeInstancesRequest.InstanceId"] = id
	query["DescribeInstancesRequest.Region"] = s.client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("foasconsole", "2019-06-01", action, query, request)

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

	v, err := jsonpath.Get("$.Instances[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instances[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VvpInstance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) RealtimeComputeVvpInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRealtimeComputeVvpInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeRealtimeComputeVvpInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for RealtimeCompute.
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
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("foasconsole", "2019-06-01", action, query, request, true)

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
				response, err = client.RpcPost("foasconsole", "2019-06-01", action, query, request, true)

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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeRealtimeComputeDeployment <<< Encapsulated get interface for RealtimeCompute Deployment.

func (s *Service) DescribeRealtimeComputeDeployment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	deploymentId := parts[2]
	namespace := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = helper.StringPointer(parts[0])

	action := fmt.Sprintf("/api/v2/namespaces/%s/deployments/%s", namespace, deploymentId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"990301"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Deployment", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) RealtimeComputeDeploymentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RealtimeComputeDeploymentStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeDeployment)
}

func (s *Service) RealtimeComputeDeploymentStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeRealtimeComputeDeployment >>> Encapsulated.

// DescribeRealtimeComputeJob <<< Encapsulated get interface for RealtimeCompute Job.

func (s *Service) DescribeRealtimeComputeJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	jobId := parts[2]
	namespace := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = helper.StringPointer(parts[0])

	action := fmt.Sprintf("/api/v2/namespaces/%s/jobs/%s", namespace, jobId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"990301"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Job", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) RealtimeComputeJobStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RealtimeComputeJobStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeJob)
}

func (s *Service) RealtimeComputeJobStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeRealtimeComputeJob >>> Encapsulated.

// DescribeRealtimeComputeMember <<< Encapsulated get interface for RealtimeCompute member.

func (s *Service) DescribeRealtimeComputeMember(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	member := parts[2]
	namespace := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = helper.StringPointer(parts[0])

	action := fmt.Sprintf("/gateway/v2/namespaces/%s/members/%s", namespace, member)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"990301"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("member", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) RealtimeComputeMemberStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RealtimeComputeMemberStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeMember)
}

func (s *Service) RealtimeComputeMemberStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeRealtimeComputeMember >>> Encapsulated.

// DescribeRealtimeComputeSqlFile <<< Encapsulated get interface for RealtimeCompute SqlFile.

func (s *Service) DescribeRealtimeComputeSqlFile(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	namespace := parts[1]
	sqlFileId := parts[2]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = helper.StringPointer(parts[0])

	action := fmt.Sprintf("/api/v2/namespaces/%s/sql-file/%s", namespace, sqlFileId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"990301"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SqlFile", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) RealtimeComputeSqlFileStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RealtimeComputeSqlFileStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeSqlFile)
}

func (s *Service) RealtimeComputeSqlFileStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// GetSqlFileRootFolderId returns the root folder ID of SQL scripts (rootType=SQL_FILE),
// used as the default parentId when creating a SqlFile without an explicit parent_id.
func (s *Service) GetSqlFileRootFolderId(workspace, namespace string) (string, error) {
	client := s.client
	var response map[string]interface{}
	action := fmt.Sprintf("/api/v2/namespaces/%s/folder", namespace)
	query := make(map[string]*string)
	query["rootType"] = helper.StringPointer("SQL_FILE")
	header := make(map[string]*string)
	header["workspace"] = helper.StringPointer(workspace)
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)
		if err != nil {
			if helper.NeedRetry(err) || strings.Contains(err.Error(), "uk_folder_id") {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, nil)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, namespace, action, errs.AlibabaCloudSdkGoERROR)
	}
	folderId, err := jsonpath.Get("$.data.folderId", response)
	if err != nil || folderId == nil || fmt.Sprint(folderId) == "" {
		return "", errs.WrapErrorf(fmt.Errorf("failed to get the SQL_FILE root folder id of namespace %s", namespace), errs.FailedGetAttributeMsg, namespace, "$.data.folderId", response)
	}
	return fmt.Sprint(folderId), nil
}

// DescribeRealtimeComputeSqlFile >>> Encapsulated.

// DescribeRealtimeComputeVariable <<< Encapsulated get interface for RealtimeCompute Variable.

func (s *Service) DescribeRealtimeComputeVariable(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	namespace := parts[1]
	request := make(map[string]*string)
	query = make(map[string]*string)
	header := make(map[string]*string)
	request["pageSize"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
	request["pageIndex"] = helper.StringPointer("1")
	header["workspace"] = helper.StringPointer(parts[0])

	action := fmt.Sprintf("/api/v2/namespaces/%s/variables", namespace)

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action+"read", response, request)

		if err != nil {
			if fmt.Sprint(response["success"]) == "false" {
				if fmt.Sprint(response["errorCode"]) == "990301" {
					return object, errs.WrapErrorf(errs.NotFoundErr("Variable", id), errs.NotFoundMsg, response)
				}
				return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.data", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Variable", id), errs.NotFoundMsg, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["name"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < helper.PageSizeLarge {
			break
		}

		pageIndex, err := strconv.Atoi(*request["pageIndex"])
		if err != nil {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		request["pageIndex"] = helper.StringPointer(strconv.Itoa(pageIndex + 1))
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Variable", id), errs.NotFoundMsg, response)
	}

	return object, nil
}

func (s *Service) RealtimeComputeVariableStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RealtimeComputeVariableStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeVariable)
}

func (s *Service) RealtimeComputeVariableStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeRealtimeComputeVariable >>> Encapsulated.

// QueryAvailableInstances queries bssopenapi BssOpenApi 2017-12-14 for a
// subscription instance's renewal state. It is inlined here (mirroring the
// root-package alicloud.BssOpenApiService.QueryAvailableInstances) so this
// package does not need to import the root package; the root BssOpenApiService
// remains for still-flat consumers. Behaviour is preserved 1:1.
func (s *Service) QueryAvailableInstances(id, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"ProductCode": productCode,
		"ProductType": productType,
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = productCodeIntl
		request["ProductType"] = productTypeIntl
	}
	if id != "" {
		request["InstanceIDs"] = id
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = productCodeIntl
				request["ProductType"] = productTypeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		resp, _ := jsonpath.Get("$.Data.InstanceList", response)
		if len(resp.([]interface{})) < 1 {
			request["ProductCode"] = productCodeIntl
			if productTypeIntl != "" {
				request["ProductType"] = productTypeIntl
			}
			endpoint = connectivity.BssOpenAPIEndpointInternational
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
	} else if id != "" {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceID"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr(productCode+"Instance", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeRealtimeComputeNamespace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}

	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = parts[0]
	query["Namespace"] = parts[1]
	query["Region"] = client.RegionId
	action := "DescribeNamespaces"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("foasconsole", "2021-10-28", action, query, request)
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

	resp, err := jsonpath.Get("$.Namespaces", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Namespaces", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Namespace", id), errs.NotFoundMsg, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Namespace"]) == parts[1] {
			return v.(map[string]interface{}), nil
		}
	}

	return nil, errs.WrapErrorf(errs.NotFoundErr("Namespace", id), errs.NotFoundWithResponse, response)
}

func (s *Service) RealtimeComputeNamespaceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RealtimeComputeNamespaceStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeNamespace)
}

func (s *Service) RealtimeComputeNamespaceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
