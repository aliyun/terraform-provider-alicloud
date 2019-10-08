package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)

type Principal struct {
	Service []string
	RAM     []string
}

type RolePolicyStatement struct {
	Effect    Effect
	Action    string
	Principal Principal
}

type RolePolicy struct {
	Statement []RolePolicyStatement
	Version   string
}

type PolicyStatement struct {
	Effect   Effect
	Action   interface{}
	Resource interface{}
}

type Policy struct {
	Statement []PolicyStatement
	Version   string
}

type RamService struct {
	client *connectivity.AliyunClient
}

func (s *RamService) ParseRolePolicyDocument(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, WrapError(err)
	}
	return policy, nil
}

func (s *RamService) ParsePolicyDocument(policyDocument string) (statement []map[string]interface{}, version string, err error) {
	policy := Policy{}
	err = json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		err = WrapError(err)
		return
	}

	version = policy.Version
	statement = make([]map[string]interface{}, 0, len(policy.Statement))
	for _, v := range policy.Statement {
		item := make(map[string]interface{})

		item["effect"] = v.Effect
		if val, ok := v.Action.([]interface{}); ok {
			item["action"] = val
		} else {
			item["action"] = []interface{}{v.Action}
		}

		if val, ok := v.Resource.([]interface{}); ok {
			item["resource"] = val
		} else {
			item["resource"] = []interface{}{v.Resource}
		}
		statement = append(statement, item)
	}
	return
}

func (s *RamService) AssembleRolePolicyDocument(ramUser, service []interface{}, version string) (string, error) {
	services := expandStringList(service)
	users := expandStringList(ramUser)

	statement := RolePolicyStatement{
		Effect: Allow,
		Action: "sts:AssumeRole",
		Principal: Principal{
			RAM:     users,
			Service: services,
		},
	}

	policy := RolePolicy{
		Version:   version,
		Statement: []RolePolicyStatement{statement},
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", WrapError(err)
	}
	return string(data), nil
}

func (s *RamService) AssemblePolicyDocument(document []interface{}, version string) (string, error) {
	var statements []PolicyStatement

	for _, v := range document {
		doc := v.(map[string]interface{})

		actions := expandStringList(doc["action"].([]interface{}))
		resources := expandStringList(doc["resource"].([]interface{}))

		statement := PolicyStatement{
			Effect:   Effect(doc["effect"].(string)),
			Action:   actions,
			Resource: resources,
		}
		statements = append(statements, statement)
	}

	policy := Policy{
		Version:   version,
		Statement: statements,
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", WrapError(err)
	}
	return string(data), nil
}

// Judge whether the role policy contains service "ecs.aliyuncs.com"
func (s *RamService) JudgeRolePolicyPrincipal(roleName string) error {
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = roleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, roleName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ram.GetRoleResponse)
	policy, err := s.ParseRolePolicyDocument(resp.Role.AssumeRolePolicyDocument)
	if err != nil {
		return WrapError(err)
	}
	for _, v := range policy.Statement {
		for _, val := range v.Principal.Service {
			if strings.Trim(val, " ") == "ecs.aliyuncs.com" {
				return nil
			}
		}
	}
	return WrapError(fmt.Errorf("Role policy services must contains 'ecs.aliyuncs.com', Now is \n%v.", resp.Role.AssumeRolePolicyDocument))
}

func (s *RamService) GetIntersection(dataMap []map[string]interface{}, allDataMap map[string]interface{}) (allData []interface{}) {
	for _, v := range dataMap {
		if len(v) > 0 {
			for key := range allDataMap {
				if _, ok := v[key]; !ok {
					allDataMap[key] = nil
				}
			}
		}
	}

	for _, v := range allDataMap {
		if v != nil {
			allData = append(allData, v)
		}
	}
	return
}

func (s *RamService) DescribeRamUser(id string) (*ram.User, error) {

	listUsersRequest := ram.CreateListUsersRequest()
	listUsersRequest.RegionId = s.client.RegionId
	var userName string
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListUsers(listUsersRequest)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, listUsersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(listUsersRequest.GetActionName(), raw, listUsersRequest.RegionId, listUsersRequest)
	users, _ := raw.(*ram.ListUsersResponse)
	for _, user := range users.Users.User {
		if user.UserId == id {
			userName = user.UserName
		}
	}
	if userName == "" {
		// the d.Id() has changed from userName to userId since v1.44.0, add the logic for backward compatibility.
		userName = id
	}
	getUserRequest := ram.CreateGetUserRequest()
	getUserRequest.RegionId = s.client.RegionId
	getUserRequest.UserName = userName
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetUser(getUserRequest)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, getUserRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(getUserRequest.GetActionName(), raw, getUserRequest.RpcRequest, getUserRequest)
	user, _ := raw.(*ram.GetUserResponse)

	return &user.User, nil
}

func (s *RamService) WaitForRamUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamUser(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.UserId == id {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, DefaultTimeoutMsg, id, GetFunc(1), ProviderERROR)
		}
	}
	return nil
}
func (s *RamService) DescribeRamGroupMembership(id string) (response *ram.ListUsersForGroupResponse, err error) {
	request := ram.CreateListUsersForGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListUsersForGroup(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*ram.ListUsersForGroupResponse)
	if len(response.Users.User) > 0 {
		return
	}
	return nil, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamGroupMembership(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroupMembership(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(len(object.Users.User)), status, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamLoginProfile(id string) (response *ram.GetLoginProfileResponse, err error) {
	request := ram.CreateGetLoginProfileRequest()
	request.RegionId = s.client.RegionId
	request.UserName = id

	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetLoginProfile(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*ram.GetLoginProfileResponse)
	return

}

func (s *RamService) WaitForRamLoginProfile(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamLoginProfile(id)
		if err != nil {
			if RamEntityNotExist(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.LoginProfile.UserName == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.LoginProfile.UserName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamGroupPolicyAttachment(id string) (response *ram.Policy, err error) {
	request := ram.CreateListPoliciesForGroupRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return nil, WrapError(err)
	}
	request.GroupName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForGroup(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	listPoliciesForGroupResponse, _ := raw.(*ram.ListPoliciesForGroupResponse)
	if len(listPoliciesForGroupResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForGroupResponse.Policies.Policy {
			if v.PolicyName == parts[1] && v.PolicyType == parts[2] {
				return &v, nil
			}
		}
	}
	return nil, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamGroupPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeRamGroupPolicyAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountAlias(id string) (*ram.GetAccountAliasResponse, error) {
	request := ram.CreateGetAccountAliasRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ram.GetAccountAliasResponse)

	return response, nil
}

func (s *RamService) DescribeRamAccessKey(id, userName string) (*ram.AccessKey, error) {
	request := ram.CreateListAccessKeysRequest()
	request.RegionId = s.client.RegionId
	request.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListAccessKeys(request)
	})

	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("RamAccessKey", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ram.ListAccessKeysResponse)
	var object *ram.AccessKey
	for _, accessKey := range response.AccessKeys.AccessKey {
		if accessKey.AccessKeyId == id {
			object = &accessKey
			return object, nil
		}
	}
	return nil, WrapErrorf(Error(GetNotFoundMessage("RamAccessKey", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
}

func (s *RamService) WaitForRamAccessKey(id, useName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamAccessKey(id, useName)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if string(status) == object.Status {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamPolicy(id string) (response *ram.GetPolicyResponse, err error) {
	request := ram.CreateGetPolicyRequest()
	request.RegionId = s.client.RegionId
	request.PolicyName = id
	request.PolicyType = "Custom"

	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetPolicy(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response = raw.(*ram.GetPolicyResponse)
	return
}

func (s *RamService) WaitForRamPolicy(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Policy.PolicyName == id {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Policy.PolicyName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRoleAttachment(id string) (response *ecs.DescribeInstanceRamRoleResponse, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	request := ecs.CreateDescribeInstanceRamRoleRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = parts[1]
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{RoleAttachmentUnExpectedJson}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidRamRoleNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response = raw.(*ecs.DescribeInstanceRamRoleResponse)
	instRoleSets := response.InstanceRamRoleSets.InstanceRamRoleSet
	if len(instRoleSets) > 0 {
		var instIds []string
		for _, item := range instRoleSets {
			if item.RamRoleName == parts[0] {
				instIds = append(instIds, item.InstanceId)
			}
		}
		ids := strings.Split(strings.TrimRight(strings.TrimLeft(strings.Replace(strings.Split(id, ":")[1], "\"", "", -1), "["), "]"), ",")
		sort.Strings(instIds)
		sort.Strings(ids)
		if reflect.DeepEqual(instIds, ids) {
			return
		}
	}
	return nil, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamRoleAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRoleAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(object.TotalCount), status, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRole(id string) (response *ram.GetRoleResponse, err error) {
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ram.GetRoleResponse)
	return
}

func (s *RamService) WaitForRamRole(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRole(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Role.RoleName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Role.RoleName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamUserPolicyAttachment(id string) (response *ram.Policy, err error) {
	request := ram.CreateListPoliciesForUserRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return nil, WrapError(err)
	}
	request.UserName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForUser(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	listPoliciesForUserResponse, _ := raw.(*ram.ListPoliciesForUserResponse)
	if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForUserResponse.Policies.Policy {
			if v.PolicyName == parts[1] && v.PolicyType == parts[2] {
				return &v, nil
			}
		}
	}
	return nil, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamUserPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeRamUserPolicyAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRolePolicyAttachment(id string) (response *ram.Policy, err error) {
	request := ram.CreateListPoliciesForRoleRequest()
	request.RegionId = s.client.RegionId
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return nil, WrapError(err)
	}
	request.RoleName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	listPoliciesForRoleResponse, _ := raw.(*ram.ListPoliciesForRoleResponse)
	if len(listPoliciesForRoleResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForRoleResponse.Policies.Policy {
			if v.PolicyName == parts[1] && v.PolicyType == parts[2] {
				return &v, nil
			}
		}
	}
	return nil, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamRolePolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeRamRolePolicyAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {

				}
				return nil
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamGroup(id string) (response *ram.GetGroupResponse, err error) {
	request := ram.CreateGetGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetGroup(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ram.GetGroupResponse)

	if response.Group.GroupName != id {
		return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return
}

func (s *RamService) WaitForRamGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Group.GroupName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Group.GroupName, id, ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountPasswordPolicy(id string) (response *ram.GetPasswordPolicyResponse, err error) {
	request := ram.CreateGetPasswordPolicyRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetPasswordPolicy(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ram.GetPasswordPolicyResponse)

	return
}
