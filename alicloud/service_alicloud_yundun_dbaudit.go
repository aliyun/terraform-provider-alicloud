package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DbauditService struct {
	client *connectivity.AliyunClient
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

func (s *DbauditService) DescribeYundunDbauditInstance(id string) (v yundun_dbaudit.Instance, err error) {
	request := yundun_dbaudit.CreateDescribeInstancesRequest()
	var instanceIds []string
	instanceIds = append(instanceIds, id)
	request.InstanceId = &instanceIds
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.DescribeInstances(request)
	})
	if err != nil {
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*yundun_dbaudit.DescribeInstancesResponse)

	if len(response.Instances) == 0 || response.Instances[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Yundun_dbaudit Instance", id)), NotFoundMsg, ProviderERROR)
	}
	v = response.Instances[0]
	return
}

func (s *DbauditService) DescribeDbauditInstanceAttribute(id string) (v yundun_dbaudit.InstanceAttribute, err error) {
	request := yundun_dbaudit.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id

	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.DescribeInstanceAttribute(request)
	})

	if err != nil {
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*yundun_dbaudit.DescribeInstanceAttributeResponse)
	if response.InstanceAttribute.InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Yundun_dbaudit Instance", id)), NotFoundMsg, ProviderERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	v = response.InstanceAttribute
	return v, WrapError(err)
}

func (s *DbauditService) StartDbauditInstance(instanceId string, vSwitchId string) error {
	request := yundun_dbaudit.CreateStartInstanceRequest()
	request.InstanceId = instanceId
	request.VswitchId = vSwitchId
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.StartInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DbauditService) UpdateDbauditInstanceDescription(instanceId string, description string) error {
	request := yundun_dbaudit.CreateModifyInstanceAttributeRequest()
	request.InstanceId = instanceId
	request.Description = description
	raw, err := s.client.WithDbauditClient(func(dbauditClient *yundun_dbaudit.Client) (interface{}, error) {
		return dbauditClient.ModifyInstanceAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DbauditService) UpdateInstanceSpec(schemaName string, specName string, d *schema.ResourceData, meta interface{}) error {
	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	request.ProductCode = "dbaudit"
	request.SubscriptionType = "Subscription"
	// only support upgrade
	request.ModifyType = "Upgrade"

	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  specName,
			Value: d.Get(schemaName).(string),
		},
	}

	raw, err := s.client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if response, _ := raw.(*bssopenapi.ModifyInstanceResponse); !response.Success {
		return WrapError(Error(response.Message))
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DbauditService) DbauditInstanceRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeYundunDbauditInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.InstanceStatus == failState {
				return object, object.InstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.InstanceStatus))
			}
		}

		return object, object.InstanceStatus, nil
	}
}

func (s *DbauditService) createRole() error {
	createRoleRequest := ram.CreateCreateRoleRequest()
	createRoleRequest.RoleName = RoleName
	createRoleRequest.Description = RoleDefaultDescription
	createRoleRequest.AssumeRolePolicyDocument = AssumeRolePolicy
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateRole(createRoleRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, RoleName, createRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(createRoleRequest.GetActionName(), raw, createRoleRequest.RpcRequest, createRoleRequest)
	return nil
}

func (s *DbauditService) attachPolicy(policyToBeAttached []PolicyRequired) error {
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
			return WrapErrorf(err, DefaultErrorMsg, RoleName, attachPolicyRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if response, err := raw.(*ram.AttachPolicyToRoleResponse); !err || !response.IsSuccess() {
			log.Printf("AttachPolicyToRoleError [%v]", response)
			return errors.New("attach policy to role failed")
		}
		addDebug(attachPolicyRequest.GetActionName(), raw, attachPolicyRequest.RpcRequest, attachPolicyRequest)

	}
	return nil
}

func (s *DbauditService) ProcessRolePolicy() error {
	getRoleRequest := ram.CreateGetRoleRequest()
	getRoleRequest.RoleName = RoleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(getRoleRequest)
	})
	log.Printf("DEBUG ProcessRolePolicy create role %v", raw.(*ram.GetRoleResponse))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, RoleName, getRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(getRoleRequest.GetActionName(), raw, getRoleRequest.RpcRequest, getRoleRequest)
	if response, _ := raw.(*ram.GetRoleResponse); response == nil || response.Role.RoleName != RoleName {
		if err := s.createRole(); err != nil {
			return WrapError(err)
		}
	}
	listPolicyForRoleRequest := ram.CreateListPoliciesForRoleRequest()
	listPolicyForRoleRequest.RoleName = RoleName
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(listPolicyForRoleRequest)
	})
	addDebug(listPolicyForRoleRequest.GetActionName(), raw, listPolicyForRoleRequest.RpcRequest, listPolicyForRoleRequest)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, RoleName, listPolicyForRoleRequest.GetActionName(), AlibabaCloudSdkGoERROR)
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

func (s *DbauditService) WaitForYundunDbauditInstance(instanceId string, status Status, timeoutSenconds time.Duration) error {
	deadline := time.Now().Add(timeoutSenconds * time.Second)
	for {
		_, err := s.DescribeYundunDbauditInstance(instanceId)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, instanceId, GetFunc(1), timeoutSenconds, "", "", ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
