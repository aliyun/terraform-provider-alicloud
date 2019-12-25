package alicloud

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type bastionhostService struct {
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

func (s *bastionhostService) DescribeYundunBastionhostInstance(id string) (v yundun_bastionhost.Instance, err error) {
	request := yundun_bastionhost.CreateDescribeInstanceBastionhostRequest()
	var instanceIds []string
	instanceIds = append(instanceIds, id)
	request.InstanceId = &instanceIds
	request.PageSize = requests.NewInteger(PageSizeSmall)
	request.CurrentPage = requests.NewInteger(1)
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.DescribeInstanceBastionhost(request)
	})
	if err != nil {
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*yundun_bastionhost.DescribeInstanceBastionhostResponse)

	if len(response.Instances) == 0 || response.Instances[0].InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Yundun_bastionhost Instance", id)), NotFoundMsg, ProviderERROR)
	}
	v = response.Instances[0]
	return
}

func (s *bastionhostService) DescribeBastionhostInstanceAttribute(id string) (v yundun_bastionhost.InstanceAttribute, err error) {
	request := yundun_bastionhost.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id

	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.DescribeInstanceAttribute(request)
	})

	if err != nil {
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*yundun_bastionhost.DescribeInstanceAttributeResponse)
	if response == nil || response.InstanceAttribute.InstanceId != id {
		return v, WrapErrorf(Error(GetNotFoundMessage("Yundun_bastionhost Instance", id)), NotFoundMsg, ProviderERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	v = response.InstanceAttribute
	return v, WrapError(err)
}

func (s *bastionhostService) StartBastionhostInstance(instanceId string, vSwitchId string, securityGroupIds []string) error {
	request := yundun_bastionhost.CreateStartInstanceRequest()
	request.InstanceId = instanceId
	request.VswitchId = vSwitchId
	request.SecurityGroupIds = &securityGroupIds
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.StartInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *bastionhostService) UpdateBastionhostInstanceDescription(instanceId string, description string) error {
	request := yundun_bastionhost.CreateModifyInstanceAttributeRequest()
	request.InstanceId = instanceId
	request.Description = description
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.ModifyInstanceAttribute(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *bastionhostService) UpdateBastionhostSecurityGroups(instanceId string, securityGroups []string) error {
	request := yundun_bastionhost.CreateConfigInstanceSecurityGroupsRequest()
	request.InstanceId = instanceId
	request.SecurityGroupIds = &securityGroups
	raw, err := s.client.WithBastionhostClient(func(BastionhostClient *yundun_bastionhost.Client) (interface{}, error) {
		return BastionhostClient.ConfigInstanceSecurityGroups(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *bastionhostService) UpdateInstanceSpec(schemaSpecMap map[string]string, d *schema.ResourceData, meta interface{}) error {
	request := bssopenapi.CreateModifyInstanceRequest()
	request.InstanceId = d.Id()

	request.ProductCode = "bastionhost"
	request.SubscriptionType = "Subscription"
	// only support upgrade
	request.ModifyType = "Upgrade"

	params := make([]bssopenapi.ModifyInstanceParameter, 0, len(schemaSpecMap))
	for schemaName, spec := range schemaSpecMap {
		params = append(params, bssopenapi.ModifyInstanceParameter{schemaName, d.Get(spec).(string)})
	}

	request.Parameter = &params
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

func (s *bastionhostService) BastionhostInstanceRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeYundunBastionhostInstance(id)
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

func (s *bastionhostService) createRole() error {
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

func (s *bastionhostService) attachPolicy(policyToBeAttached []BastionhostPolicyRequired) error {
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

func (s *bastionhostService) ProcessRolePolicy() error {
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

func (s *bastionhostService) WaitForYundunBastionhostInstance(instanceId string, status Status, timeoutSenconds time.Duration) error {
	deadline := time.Now().Add(timeoutSenconds * time.Second)
	for {
		_, err := s.DescribeYundunBastionhostInstance(instanceId)

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

func (s *bastionhostService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []yundun_bastionhost.TagResource, err error) {
	request := yundun_bastionhost.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = strings.ToUpper(string(resourceType))
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []yundun_bastionhost.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, yundun_bastionhost.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	var raw interface{}

	raw, err = s.client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*yundun_bastionhost.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *bastionhostService) tagsToMap(tags []yundun_bastionhost.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *bastionhostService) ignoreTag(t yundun_bastionhost.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *bastionhostService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := yundun_bastionhost.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = strings.ToUpper(string(resourceType))
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := yundun_bastionhost.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = strings.ToUpper(string(resourceType))
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithBastionhostClient(func(client *yundun_bastionhost.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

	}

	return nil
}

func (s *bastionhostService) diffTags(oldTags, newTags []yundun_bastionhost.TagResourcesTag) ([]yundun_bastionhost.TagResourcesTag, []yundun_bastionhost.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []yundun_bastionhost.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *bastionhostService) tagsFromMap(m map[string]interface{}) []yundun_bastionhost.TagResourcesTag {
	result := make([]yundun_bastionhost.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, yundun_bastionhost.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}
