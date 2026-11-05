package yundundbaudit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

type PolicyRequired struct {
	PolicyName string
	PolicyType string
}

const (
	RoleDefaultDescription = "DbAudit perform the default role to access your other cloud resources"
	RoleName               = "AliyunDbAuditDefaultRole"
	AssumeRolePolicy       = `{
		"Statement": [
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": [
						"dbaudit.aliyuncs.com"
					]
				}
			}
		],
		"Version": "1"
	}`
)

var policyRequired = []PolicyRequired{
	{
		PolicyName: "AliyunDbAuditRolePolicy",
		PolicyType: "System",
	},
	{
		PolicyName: "AliyunLogFullAccess",
		PolicyType: "System",
	},
}

func (s *Service) DescribeYundunDbauditInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstanceAttribute"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
		request["ClientToken"] = helper.BuildClientToken(action)
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
	return response["InstanceAttribute"].(map[string]interface{}), nil
}

func (s *Service) StartDbauditInstance(instanceId string, vSwitchId string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "StartInstance"
	request := map[string]interface{}{
		"InstanceId": instanceId,
		"VswitchId":  vSwitchId,
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
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
	return errs.WrapError(err)
}

func (s *Service) UpdateDbauditInstanceDescription(instanceId string, description string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "ModifyInstanceAttribute"
	request := map[string]interface{}{
		"InstanceId":  instanceId,
		"Description": description,
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
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
	return errs.WrapError(err)
}

func (s *Service) UpdateInstanceSpec(schemaName string, specName string, d *schema.ResourceData) error {
	client := s.client
	action := "ModifyInstance"
	var response map[string]interface{}
	var err error
	var endpoint string
	query := make(map[string]interface{})
	request := map[string]interface{}{
		"InstanceId":       d.Id(),
		"ClientToken":      helper.BuildClientToken(action),
		"ProductCode":      "dbaudit",
		"SubscriptionType": "Subscription",
		"ModifyType":       "Upgrade",
	}
	parameterMapList := []map[string]interface{}{
		{
			"Code":  specName,
			"Value": d.Get(schemaName),
		},
	}
	request["Parameter"] = parameterMapList
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = helper.BuildClientToken(action)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	return errs.WrapError(err)
}

func (s *Service) DbauditInstanceRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeYundunDbauditInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		currentStatus := object["InstanceStatus"].(string)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}

		return object, currentStatus, nil
	}
}

func (s *Service) createRole() error {
	createRoleRequest := ram.CreateCreateRoleRequest()
	createRoleRequest.RoleName = RoleName
	createRoleRequest.Description = RoleDefaultDescription
	createRoleRequest.AssumeRolePolicyDocument = AssumeRolePolicy
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateRole(createRoleRequest)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, RoleName, createRoleRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(createRoleRequest.GetActionName(), raw, createRoleRequest.RpcRequest, createRoleRequest)
	return nil
}

func (s *Service) attachPolicy(policyToBeAttached []PolicyRequired) error {
	log.Printf("DEBUG attachPolicy policyRequred %v", policyToBeAttached)
	attachPolicyRequest := ram.CreateAttachPolicyToRoleRequest()
	for _, policy := range policyToBeAttached {
		log.Printf("DEBUG attach Policy in policyRequred %v", policy)
		attachPolicyRequest.RoleName = RoleName
		attachPolicyRequest.PolicyName = policy.PolicyName
		attachPolicyRequest.PolicyType = policy.PolicyType
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.AttachPolicyToRole(attachPolicyRequest)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, RoleName, attachPolicyRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		if response, err := raw.(*ram.AttachPolicyToRoleResponse); !err || !response.IsSuccess() {
			log.Printf("AttachPolicyToRoleError [%v]", response)
			return errors.New("attach policy to role failed")
		}
		logs.AddDebug(attachPolicyRequest.GetActionName(), raw, attachPolicyRequest.RpcRequest, attachPolicyRequest)

	}
	return nil
}

func (s *Service) ProcessRolePolicy() error {
	getRoleRequest := ram.CreateGetRoleRequest()
	getRoleRequest.RoleName = RoleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(getRoleRequest)
	})
	log.Printf("DEBUG ProcessRolePolicy create role %v", raw.(*ram.GetRoleResponse))
	logs.AddDebug(getRoleRequest.GetActionName(), raw, getRoleRequest.RpcRequest, getRoleRequest)
	response, _ := raw.(*ram.GetRoleResponse)
	if err != nil || response == nil || response.Role.RoleName != RoleName {
		if err := s.createRole(); err != nil {
			return errs.WrapError(err)
		}
	}
	listPolicyForRoleRequest := ram.CreateListPoliciesForRoleRequest()
	listPolicyForRoleRequest.RoleName = RoleName
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(listPolicyForRoleRequest)
	})
	logs.AddDebug(listPolicyForRoleRequest.GetActionName(), raw, listPolicyForRoleRequest.RpcRequest, listPolicyForRoleRequest)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, RoleName, listPolicyForRoleRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	var policyToAttach []PolicyRequired
	if response, _ := raw.(*ram.ListPoliciesForRoleResponse); response != nil && response.IsSuccess() {
		for _, required := range policyRequired {
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

func (s *Service) DescribeTags(id string, resourceType string) (object interface{}, err error) {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources", response))
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
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}
	return tags, nil
}

func (s *Service) setInstanceTags(d *schema.ResourceData, resourceType string) (err error) {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UnTagResources"
			request = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
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
				response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
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

func (s *Service) UpdateResourceGroup(resourceId, resourceGroupId string) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "MoveResourceGroup"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"ResourceId":      resourceId,
		"ResourceGroupId": resourceGroupId,
		"ResourceType":    "ALIYUN::DBAUDIT::INSTANCE",
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Yundun-dbaudit", "2018-10-29", action, nil, request, false)
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
	return err
}
