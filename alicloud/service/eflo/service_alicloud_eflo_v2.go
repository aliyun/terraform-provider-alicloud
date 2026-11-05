package eflo

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

type EfloServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewEfloServiceV2 constructs an EfloServiceV2 bound to client. The client field
// is unexported, so split-package tests build the value through this helper.
func NewEfloServiceV2(client *connectivity.AliyunClient) *EfloServiceV2 {
	return &EfloServiceV2{client: client}
}

// DescribeEcsDataDisksByNodeId lists the ECS data disks attached to the given
// node (instance) id. The ECS service is not yet migrated to a subpackage, so
// the eflo package carries this self-contained read for the node resource's
// data_disk flattening; the request/response contract is identical to the
// root EcsServiceV2 method of the same name.
func (s *EfloServiceV2) DescribeEcsDataDisksByNodeId(nodeId string) (objects []map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDisks"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = nodeId
	request["DiskType"] = "data"
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, nodeId, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Disks.Disk[*]", response)
	if err != nil {
		return nil, nil
	}

	for _, disk := range v.([]interface{}) {
		objects = append(objects, disk.(map[string]interface{}))
	}

	return objects, nil
}

// DescribeEfloNode <<< Encapsulated get interface for Eflo Node.

func (s *EfloServiceV2) DescribeEfloNode(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NodeId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeNode"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalDependencyError.RequestLimitExceeded"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"RESOURCE_NOT_FOUND"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Node", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *EfloServiceV2) DescribeNodeListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = id
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "Node"
	action := "ListTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"SystemError"}) || helper.NeedRetry(err) {
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

func (s *EfloServiceV2) DescribeNodeQueryAvailableInstances(d *schema.ResourceData) (object map[string]interface{}, err error) {

	installPai := false
	if v, ok := d.GetOk("install_pai"); ok && v.(bool) {
		installPai = true
	}
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
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if installPai {
		request["ProductCode"] = "learn"
		request["ProductType"] = "learn_eflocomputing_public_cn"
	}
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
		if installPai {
			return nil, errs.WrapError(errs.Error("InstallPai currently does not support pay-as-you-go products."))
		}
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if installPai {
			request["ProductCode"] = "learn"
			request["ProductType"] = "learn_eflocomputing_public_intl"
		}
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
			if installPai {
				return nil, errs.WrapError(errs.Error("InstallPai currently does not support pay-as-you-go products."))
			}
		}
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
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				if installPai {
					request["ProductCode"] = "learn"
					request["ProductType"] = "learn_eflocomputing_public_intl"
				}
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_computinginstance_public_intl"
					if installPai {
						return retry.RetryableError(err)
					}
				}
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
		return object, errs.WrapErrorf(errs.NotFoundErr("Node", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *EfloServiceV2) DescribeNodeListClusterNodes(d *schema.ResourceData) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClusterId"] = d.Get("cluster_id")
	request["NodeGroupId"] = d.Get("node_group_id")
	request["MaxResults"] = 100
	action := "ListClusterNodes"
	id := d.Id()

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

		v, err := jsonpath.Get("$.Nodes[*]", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Nodes[*]", response)
		}

		if len(v.([]interface{})) == 0 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Node", id), errs.NotFoundMsg, response)
		}
		for _, vv := range v.([]interface{}) {
			if fmt.Sprint(vv.(map[string]interface{})["NodeId"]) == id {
				return vv.(map[string]interface{}), nil
			}
		}
		if v, ok := response["NextToken"].(string); ok && v != "" {
			request["NextToken"] = v
		} else {
			break
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("Node", id), errs.NotFoundMsg, response)
}

func (s *EfloServiceV2) EfloNodeStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EfloNodeStateRefreshFuncWithApi(id, field, failStates, s.DescribeEfloNode)
}

func (s *EfloServiceV2) EfloNodeStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeEfloNode >>> Encapsulated.

// DescribeEfloCluster <<< Encapsulated get interface for Eflo Cluster.

func (s *EfloServiceV2) DescribeEfloCluster(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeCluster"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

	currentStatus := response["ClusterId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *EfloServiceV2) DescribeClusterListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = id
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "Cluster"
	action := "ListTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

func (s *EfloServiceV2) EfloClusterStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloCluster(id)
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

// DescribeEfloCluster >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Eflo.
func (s *EfloServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
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
				response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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
				response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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

// DescribeEfloNodeGroup <<< Encapsulated get interface for Eflo NodeGroup.

func (s *EfloServiceV2) DescribeEfloNodeGroup(id string) (object map[string]interface{}, err error) {
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
	request["ClusterId"] = parts[0]
	request["NodeGroupId"] = parts[1]
	request["RegionId"] = client.RegionId
	action := "ListNodeGroups"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

	v, err := jsonpath.Get("$.Groups[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("NodeGroup", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("NodeGroup", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["GroupId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("NodeGroup", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *EfloServiceV2) DescribeNodeGroupListClusterNodes(id string) (object map[string]interface{}, err error) {
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
	request["ClusterId"] = parts[0]
	request["NodeGroupId"] = parts[1]
	request["RegionId"] = client.RegionId
	action := "ListClusterNodes"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

func (s *EfloServiceV2) EfloNodeGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloNodeGroup(id)
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

func (s *EfloServiceV2) DescribeAsyncEfloNodeGroupStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeTask(d, res)
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

// DescribeEfloNodeGroup >>> Encapsulated.

// DescribeAsyncDescribeTask <<< Encapsulated for Eflo.
func (s *EfloServiceV2) DescribeAsyncDescribeTask(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TaskId"], err = jsonpath.Get("$.TaskId", res)
	request["RegionId"] = client.RegionId
	action := "DescribeTask"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

// DescribeAsyncDescribeTask >>> Encapsulated.

// DescribeNodeIdByTaskId polls DescribeTask until NodeIds is available and returns the first NodeId.
func (s *EfloServiceV2) DescribeNodeIdByTaskId(taskId string, timeout time.Duration) (nodeId string, err error) {
	client := s.client
	action := "DescribeTask"
	var response map[string]interface{}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		request := map[string]interface{}{
			"RegionId": client.RegionId,
		}
		query := map[string]interface{}{
			"TaskId": taskId,
		}
		var e error
		response, e = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
		if e != nil {
			if errs.IsExpectedErrors(e, []string{"InternalDependencyError.RequestLimitExceeded"}) || helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		nodeIdsRaw, _ := jsonpath.Get("$.NodeIds", response)
		if nodeIds, ok := nodeIdsRaw.([]interface{}); ok && len(nodeIds) > 0 {
			return nil
		}
		wait()
		return retry.RetryableError(fmt.Errorf("task %s: waiting for NodeIds", taskId))
	})
	logs.AddDebug(action, response, nil)
	if err != nil {
		return "", errs.WrapError(err)
	}
	nodeIdsRaw, e := jsonpath.Get("$.NodeIds", response)
	if e != nil || nodeIdsRaw == nil {
		return "", errs.WrapError(fmt.Errorf("DescribeTask %s: NodeIds not found in response", taskId))
	}
	nodeIds, ok := nodeIdsRaw.([]interface{})
	if !ok || len(nodeIds) == 0 {
		return "", errs.WrapError(fmt.Errorf("DescribeTask %s: NodeIds is empty", taskId))
	}
	return fmt.Sprint(nodeIds[0]), nil
}

// DescribeEfloInvocation <<< Encapsulated get interface for Eflo Invocation.

func (s *EfloServiceV2) DescribeEfloInvocation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InvokeId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeInvocations"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

	v, err := jsonpath.Get("$.Invocations.Invocation[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Invocations.Invocation[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Invocation", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloInvocationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloInvocation(id)
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

// DescribeEfloInvocation >>> Encapsulated.

// DescribeEfloExperimentPlanTemplate <<< Encapsulated get interface for Eflo ExperimentPlanTemplate.

func (s *EfloServiceV2) DescribeEfloExperimentPlanTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = id
	request["RegionId"] = client.RegionId
	action := "GetExperimentPlanTemplate"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-cnp", "2023-08-28", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ExperimentPlanTemplate", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloExperimentPlanTemplateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloExperimentPlanTemplate(id)
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

// DescribeEfloExperimentPlanTemplate >>> Encapsulated.

// DescribeEfloResource <<< Encapsulated get interface for Eflo Resource.

func (s *EfloServiceV2) DescribeEfloResource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = id
	request["RegionId"] = client.RegionId
	action := "GetResource"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-cnp", "2023-08-28", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"A1002"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Resource", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloResourceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloResource(id)
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

// DescribeEfloResource >>> Encapsulated.

// DescribeEfloExperimentPlan <<< Encapsulated get interface for Eflo ExperimentPlan.

func (s *EfloServiceV2) DescribeEfloExperimentPlan(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PlanId"] = id
	request["RegionId"] = client.RegionId
	action := "GetExperimentPlan"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-cnp", "2023-08-28", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ExperimentPlan", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloExperimentPlanStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloExperimentPlan(id)
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

// DescribeEfloExperimentPlan >>> Encapsulated.

// DescribeEfloVsc <<< Encapsulated get interface for Eflo Vsc.

func (s *EfloServiceV2) DescribeEfloVsc(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VscId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeVsc"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Resource.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Vsc", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	currentStatus := response["VscId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Vsc", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *EfloServiceV2) DescribeVscListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = id
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "Vsc"
	action := "ListTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"SystemError"}) || helper.NeedRetry(err) {
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

func (s *EfloServiceV2) EfloVscStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloVsc(id)
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

// DescribeEfloVsc >>> Encapsulated.

// DescribeEfloNodeGroupAttachment <<< Encapsulated get interface for Eflo NodeGroupAttachment.

func (s *EfloServiceV2) DescribeEfloNodeGroupAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = parts[0]
	request["NodeGroupId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["MaxResults"] = helper.PageSizeLarge
	action := "ListClusterNodes"

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

		v, _ := jsonpath.Get("$.Nodes", response)
		if nodes, ok := v.([]interface{}); ok {
			for _, n := range nodes {
				item := n.(map[string]interface{})
				if fmt.Sprint(item["NodeId"]) == parts[2] {
					return item, nil
				}
			}
		}

		// robust nextToken handling: accept any type and stringify
		if v, exists := response["NextToken"]; !exists || v == nil || fmt.Sprint(v) == "" {
			break
		} else {
			request["NextToken"] = fmt.Sprint(v)
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("NodeGroupAttachment", id), errs.NotFoundMsg, response)
}

func (s *EfloServiceV2) EfloNodeGroupAttachmentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloNodeGroupAttachment(id)
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

func (s *EfloServiceV2) DescribeAsyncEfloNodeGroupAttachmentStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeTask(d, res)
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

// DescribeEfloNodeGroupAttachment >>> Encapsulated.
// DescribeEfloEr <<< Encapsulated get interface for Eflo Er.

func (s *EfloServiceV2) DescribeEfloEr(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ErId"] = id
	request["RegionId"] = client.RegionId
	action := "GetEr"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"1003"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Er", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Content", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Content", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloErStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEfloEr(id)
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

// DescribeEfloEr >>> Encapsulated.

// DescribeEfloVpd <<< Encapsulated get interface for Eflo Vpd.

func (s *EfloServiceV2) DescribeEfloVpd(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpdId"] = id
	request["RegionId"] = client.RegionId
	action := "GetVpd"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"1003"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Vpd", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Content", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Content", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloVpdStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EfloVpdStateRefreshFuncWithApi(id, field, failStates, s.DescribeEfloVpd)
}

func (s *EfloServiceV2) EfloVpdStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeEfloVpd >>> Encapsulated.

// DescribeEfloVpdGrantRule <<< Encapsulated get interface for Eflo VpdGrantRule.

func (s *EfloServiceV2) DescribeEfloVpdGrantRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GrantRuleId"] = id
	request["RegionId"] = client.RegionId
	action := "GetVpdGrantRule"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"1003"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("VpdGrantRule", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Content", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Content", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *EfloServiceV2) EfloVpdGrantRuleStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EfloVpdGrantRuleStateRefreshFuncWithApi(id, field, failStates, s.DescribeEfloVpdGrantRule)
}

func (s *EfloServiceV2) EfloVpdGrantRuleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeEfloVpdGrantRule >>> Encapsulated.

// DescribeEfloHyperNode <<< Encapsulated get interface for Eflo HyperNode.

func (s *EfloServiceV2) DescribeEfloHyperNode(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["HyperNodeId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeHyperNode"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidHyperNodeId.NotFound", "RESOURCE_NOT_FOUND"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("HyperNode", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *EfloServiceV2) DescribeHyperNodeQueryAvailableInstances(d *schema.ResourceData) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = id

	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
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
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_computinginstance_public_intl"
				}
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
		return object, errs.WrapErrorf(errs.NotFoundErr("HyperNode", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *EfloServiceV2) DescribeHyperNodeListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = id
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "HyperNode"
	action := "ListTagResources"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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
func (s *EfloServiceV2) DescribeHyperNodeListClusterHyperNodes(clusterId, id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClusterId"] = clusterId
	action := "ListClusterHyperNodes"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

	v, err := jsonpath.Get("$.HyperNodes[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.HyperNodes[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("HyperNode", id), errs.NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["HyperNodeId"]) != id {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("HyperNode", id), errs.NotFoundMsg, response)
}

func (s *EfloServiceV2) EfloHyperNodeStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.EfloHyperNodeStateRefreshFuncWithApi(id, field, failStates, s.DescribeEfloHyperNode)
}

func (s *EfloServiceV2) EfloHyperNodeStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeEfloHyperNode >>> Encapsulated.
