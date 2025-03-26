package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

type YundunBastionhostService struct {
	client *connectivity.AliyunClient
}

type BastionhostPolicyRequired struct {
	PolicyName string
	PolicyType string
}

const (
	BastionhostRoleDefaultDescription = "Bastionhost will access other cloud resources by playing this role by default"
	BastionhostRoleName               = "AliyunBastionHostDefaultRole"
	BastionhostAssumeRolePolicy       = `{
		"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"bastionhost.aliyuncs.com"
					]
				}
			}
		],
		"Version": "1"
	}`
)

var bastionhostpolicyRequired = []BastionhostPolicyRequired{
	{
		PolicyName: "AliyunBastionHostRolePolicy",
		PolicyType: "System",
	},
}

func (s *YundunBastionhostService) DescribeBastionhostInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidApi"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus"}) {
			return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.InstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.InstanceAttribute", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostInstances(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeInstances"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": []string{id},
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidApi"}) {
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
	v, err := jsonpath.Get("$.Instances", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("BastionhostInstance", id), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) StartBastionhostInstance(instanceId string, vSwitchId string, securityGroupIds []string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "StartInstance"
	request := map[string]interface{}{
		"VswitchId":        vSwitchId,
		"InstanceId":       instanceId,
		"SecurityGroupIds": securityGroupIds,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *YundunBastionhostService) UpdateBastionhostInstanceDescription(instanceId string, description string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "ModifyInstanceAttribute"
	request := map[string]interface{}{
		"Description": description,
		"InstanceId":  instanceId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *YundunBastionhostService) UpdateBastionhostSecurityGroups(instanceId string, securityGroups []string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "ConfigInstanceSecurityGroups"
	request := map[string]interface{}{
		"AuthorizedSecurityGroups": securityGroups,
		"InstanceId":               instanceId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *YundunBastionhostService) UpdateInstanceSpec(schemaSpecMap map[string]string, d *schema.ResourceData, meta interface{}) (err error) {
	client := s.client
	var endpoint string
	action := "ModifyInstance"
	request := map[string]interface{}{
		"RegionId":         connectivity.Hangzhou,
		"InstanceId":       d.Id(),
		"ProductCode":      "bastionhost",
		"ProductType":      "bastionhost",
		"SubscriptionType": "Subscription",
		"ModifyType":       "Upgrade",
	}
	if client.IsInternationalAccount() {
		request["RegionId"] = connectivity.APSouthEast1
		request["ProductType"] = "bastionhost_std_public_intl"
	}
	parameterMapList := make([]map[string]interface{}, 0)
	for schemaName, spec := range schemaSpecMap {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  schemaName,
			"Value": d.Get(spec),
		})
	}
	request["Parameter"] = parameterMapList
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["RegionId"] = connectivity.APSouthEast1
				request["ProductType"] = "bastionhost_std_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func (s *YundunBastionhostService) BastionhostInstanceRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeBastionhostInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["InstanceStatus"]) == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

func (s *YundunBastionhostService) createRole() error {
	createRoleRequest := ram.CreateCreateRoleRequest()
	createRoleRequest.RoleName = BastionhostRoleName
	createRoleRequest.Description = BastionhostRoleDefaultDescription
	createRoleRequest.AssumeRolePolicyDocument = BastionhostAssumeRolePolicy
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateRole(createRoleRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, BastionhostRoleName, createRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(createRoleRequest.GetActionName(), raw, createRoleRequest.RpcRequest, createRoleRequest)
	return nil
}

func (s *YundunBastionhostService) attachPolicy(policyToBeAttached []BastionhostPolicyRequired) error {
	attachPolicyRequest := ram.CreateAttachPolicyToRoleRequest()
	for _, policy := range policyToBeAttached {
		attachPolicyRequest.RoleName = BastionhostRoleName
		attachPolicyRequest.PolicyName = policy.PolicyName
		attachPolicyRequest.PolicyType = policy.PolicyType
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.AttachPolicyToRole(attachPolicyRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, BastionhostRoleName, attachPolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if response, err := raw.(*ram.AttachPolicyToRoleResponse); !err || !response.IsSuccess() {
			return WrapError(errors.New("attach policy to role failed"))
		}
		addDebug(attachPolicyRequest.GetActionName(), raw, attachPolicyRequest.RpcRequest, attachPolicyRequest)

	}
	return nil
}

func (s *YundunBastionhostService) ProcessRolePolicy() error {
	getRoleRequest := ram.CreateGetRoleRequest()
	getRoleRequest.RoleName = BastionhostRoleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(getRoleRequest)
	})
	response, _ := raw.(*ram.GetRoleResponse)
	if err != nil || response == nil || response.Role.RoleName != BastionhostRoleName {
		if err := s.createRole(); err != nil {
			return err
		}
	}
	addDebug(getRoleRequest.GetActionName(), raw, getRoleRequest.RpcRequest, getRoleRequest)
	listPolicyForRoleRequest := ram.CreateListPoliciesForRoleRequest()
	listPolicyForRoleRequest.RoleName = BastionhostRoleName
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(listPolicyForRoleRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, BastionhostRoleName, listPolicyForRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(listPolicyForRoleRequest.GetActionName(), raw, listPolicyForRoleRequest.RpcRequest, listPolicyForRoleRequest)
	var policyToAttach []BastionhostPolicyRequired
	if response, _ := raw.(*ram.ListPoliciesForRoleResponse); response != nil && response.IsSuccess() {
		for _, required := range bastionhostpolicyRequired {
			contains := false
			for _, policy := range response.Policies.Policy {
				if required.PolicyName == policy.PolicyName {
					contains = true
				}
			}
			if !contains {
				policyToAttach = append(policyToAttach, required)
			}
		}
	}

	if policyToAttach != nil && len(policyToAttach) > 0 {
		return s.attachPolicy(policyToAttach)
	}

	return nil
}

func (s *YundunBastionhostService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
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
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *YundunBastionhostService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) (err error) {
	client := s.client
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		if len(removed) > 0 {
			var response map[string]interface{}
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removed {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
			return nil
		}

		if len(added) > 0 {
			var response map[string]interface{}
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
			return nil
		}
	}

	return nil
}

func (s *YundunBastionhostService) UpdateResourceGroup(resourceId, resourceType, resourceGroupId string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "MoveResourceGroup"
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ResourceType":    resourceType,
		"ResourceId":      resourceId,
		"ResourceGroupId": resourceGroupId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, resourceId, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *YundunBastionhostService) DescribeBastionhostUserGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUserGroup"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:UserGroup", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.UserGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.UserGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUser"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:User", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) DescribeBastionhostHostGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetHostGroup"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[1],
		"InstanceId":  parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:HostGroup", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostUserAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListUsers"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
		"PageNumber":  1,
		"PageSize":    50,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus"}) {
				return object, WrapErrorf(NotFoundErr("Bastionhost:UserAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Users", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["UserId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return
}

func (s *YundunBastionhostService) DescribeBastionhostHost(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetHost"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"HostId":     parts[1],
		"InstanceId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND", "HostNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:Host", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Host", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Host", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetHostAccount"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"HostAccountId": parts[1],
		"InstanceId":    parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND", "HostAccountNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:HostAccount", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccount", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccount", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) DescribeBastionhostHostAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListHosts"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[1],
		"InstanceId":  parts[0],
		"PageNumber":  1,
		"PageSize":    0,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
			if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus"}) {
				return object, WrapErrorf(NotFoundErr("Bastionhost:HostAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Hosts", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Hosts", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["HostId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return
}
func (s *YundunBastionhostService) DescribeBastionhostHostAccountUserAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListHostAccountsForUser"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"HostId":     parts[2],
		"InstanceId": parts[0],
		"UserId":     parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:HostAccountUserAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccounts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccounts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostAccountUserGroupAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListHostAccountsForUserGroup"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostId":      parts[2],
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:HostAccountUserGroupAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccounts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccounts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostGroupAccountUserAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListHostGroupAccountNamesForUser"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[2],
		"InstanceId":  parts[0],
		"UserId":      parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:HostGroupAccountUserAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccountNames", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccountNames", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) DescribeBastionhostHostGroupAccountUserGroupAttachment(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListHostGroupAccountNamesForUserGroup"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HostGroupId": parts[2],
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(10*time.Second, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Commodity.BizError.InvalidStatus", "OBJECT_NOT_FOUND"}) {
			return object, WrapErrorf(NotFoundErr("Bastionhost:HostGroupAccountUserGroupAttachment", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HostAccountNames", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccountNames", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return v.([]interface{}), nil
}

func (s *YundunBastionhostService) EnableInstancePublicAccess(id string) (err error) {
	var response map[string]interface{}
	client := s.client
	action := "EnableInstancePublicAccess"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func (s *YundunBastionhostService) DisableInstancePublicAccess(id string) (err error) {
	var response map[string]interface{}
	client := s.client
	action := "DisableInstancePublicAccess"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func (s *YundunBastionhostService) DescribeBastionhostHostShareKey(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetHostShareKey"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
	v, err := jsonpath.Get("$.HostShareKey", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostShareKey", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) GetHostShareKey(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetHostShareKey"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
	v, err := jsonpath.Get("$.HostShareKey", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostShareKey", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *YundunBastionhostService) DescribeBastionhostHostAccountShareKeyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListHostAccountsForHostShareKey"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
		"PageSize":       PageSizeMedium,
		"PageNumber":     1,
	}
	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
		v, err := jsonpath.Get("$.HostAccounts", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HostAccounts", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["HostsAccountId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("Bastionhost", id), NotFoundWithResponse, response)
	}
	return
}

func (s *YundunBastionhostService) DescribeBastionhostAdAuthServer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetInstanceADAuthServer"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
	v, err := jsonpath.Get("$.AD", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AD", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *YundunBastionhostService) DescribeBastionhostLdapAuthServer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetInstanceLDAPAuthServer"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Yundun-bastionhost", "2019-12-09", action, nil, request, true)
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
	v, err := jsonpath.Get("$.LDAP", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.LDAP", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
