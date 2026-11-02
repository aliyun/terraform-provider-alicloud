package ram

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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

type PolicyDocumentStatementCondition struct {
	Operator string
	Variable string
	Values   interface{}
}

type PolicyDocumentStatementPrincipal struct {
	Entity      string
	Identifiers interface{}
}

type PolicyDocumentStatement struct {
	Effect    Effect
	Action    interface{}
	Resource  interface{}                         `json:",omitempty"`
	Principal PolicyDocumentStatementPrincipalSet `json:",omitempty"`
	Condition PolicyDocumentStatementConditionSet `json:",omitempty"`
}

type PolicyDocument struct {
	Statement []PolicyDocumentStatement
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

type Status string

const (
	Active   Status = "Active"
	Inactive Status = "Inactive"
	Normal   Status = "Normal"
	Deleted  Status = "Deleted"
)

type RamService struct {
	client *connectivity.AliyunClient
}

func NewRamService(client *connectivity.AliyunClient) *RamService {
	return &RamService{client: client}
}

type PolicyDocumentStatementPrincipalSet []PolicyDocumentStatementPrincipal
type PolicyDocumentStatementConditionSet []PolicyDocumentStatementCondition

func (s *RamService) ParseRolePolicyDocument(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, errs.WrapError(err)
	}
	return policy, nil
}

func (s *RamService) ParsePolicyDocument(policyDocument string) (statement []map[string]interface{}, version string, err error) {
	policy := Policy{}
	err = json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		err = errs.WrapError(err)
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
	services := convert.ExpandStringList(service)
	users := convert.ExpandStringList(ramUser)

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
		return "", errs.WrapError(err)
	}
	return string(data), nil
}

func (s *RamService) AssemblePolicyDocument(document []interface{}, version string) (string, error) {
	var statements []PolicyStatement

	for _, v := range document {
		doc := v.(map[string]interface{})

		actions := convert.ExpandStringList(doc["action"].([]interface{}))
		resources := convert.ExpandStringList(doc["resource"].([]interface{}))

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
		return "", errs.WrapError(err)
	}
	return string(data), nil
}

// Judge whether the role policy contains service "ecs.aliyuncs.com"
func (s *RamService) JudgeRolePolicyPrincipal(roleName string) error {
	request := ramsdk.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = roleName
	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, roleName, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ramsdk.GetRoleResponse)
	policy, err := s.ParseRolePolicyDocument(resp.Role.AssumeRolePolicyDocument)
	if err != nil {
		return errs.WrapError(err)
	}
	for _, v := range policy.Statement {
		for _, val := range v.Principal.Service {
			if strings.Trim(val, " ") == "ecs.aliyuncs.com" {
				return nil
			}
		}
	}
	return errs.WrapError(fmt.Errorf("Role policy services must contains 'ecs.aliyuncs.com', Now is \n%v.", resp.Role.AssumeRolePolicyDocument))
}

func (s *RamService) GetIntersection(dataMap []map[string]interface{}, allDataMap map[string]interface{}, groupOrUserNameOk, policyNameOk bool) (allData []interface{}) {
	if (groupOrUserNameOk || policyNameOk) && len(dataMap) == 0 {
		return
	}

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

func (s *RamService) DescribeRamUser(id string) (*ramsdk.User, error) {
	user := &ramsdk.User{}
	listUsersRequest := ramsdk.CreateListUsersRequest()
	listUsersRequest.RegionId = s.client.RegionId
	listUsersRequest.MaxItems = requests.NewInteger(100)
	var userName string

	for {
		raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListUsers(listUsersRequest)
		})
		if err != nil {
			return user, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, listUsersRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(listUsersRequest.GetActionName(), raw, listUsersRequest.RegionId, listUsersRequest)
		response, _ := raw.(*ramsdk.ListUsersResponse)
		for _, user := range response.Users.User {
			// the d.Id() has changed from userName to userId since v1.44.0, add the logic for backward compatibility.
			if user.UserId == id || user.UserName == id {
				userName = user.UserName
				break
			}
		}
		if userName != "" || !response.IsTruncated {
			break
		}
		listUsersRequest.Marker = response.Marker
	}

	if userName == "" {
		return user, errs.WrapErrorf(fmt.Errorf("there is no ram user with id or name is %s", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	getUserRequest := ramsdk.CreateGetUserRequest()
	getUserRequest.RegionId = s.client.RegionId
	getUserRequest.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.GetUser(getUserRequest)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return user, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return user, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, getUserRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(getUserRequest.GetActionName(), raw, getUserRequest.RpcRequest, getUserRequest)
	response, _ := raw.(*ramsdk.GetUserResponse)

	return &response.User, nil
}

func (s *RamService) WaitForRamUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamUser(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.UserId == id {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.DefaultTimeoutMsg, id, helper.GetFunc(1), errs.ProviderERROR)
		}
	}
	return nil
}
func (s *RamService) DescribeRamGroupMembership(id string) (*ramsdk.ListUsersForGroupResponse, error) {
	response := &ramsdk.ListUsersForGroupResponse{}
	request := ramsdk.CreateListUsersForGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupName = id
	users := []ramsdk.User{}

	var raw interface{}
	var err error
	for {
		raw, err = s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
				return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response = raw.(*ramsdk.ListUsersForGroupResponse)
		users = append(users, response.Users.User...)
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}
	response.Users.User = users
	if len(users) > 0 {
		return response, nil
	}
	return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *RamService) WaitForRamGroupMembership(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroupMembership(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, strconv.Itoa(len(object.Users.User)), status, errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamLoginProfile(id string) (*ramsdk.GetLoginProfileResponse, error) {
	response := &ramsdk.GetLoginProfileResponse{}
	request := ramsdk.CreateGetLoginProfileRequest()
	request.RegionId = s.client.RegionId
	request.UserName = id

	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.GetLoginProfile(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.User.LoginProfile", "EntityNotExist.User"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("Ram:LoginProfile", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response.RequestId))
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response = raw.(*ramsdk.GetLoginProfileResponse)

	return response, nil
}

func (s *RamService) WaitForRamLoginProfile(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamLoginProfile(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.LoginProfile.UserName == id {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.LoginProfile.UserName, id, errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamGroupPolicyAttachment(id string) (*ramsdk.Policy, error) {
	response := &ramsdk.Policy{}
	request := ramsdk.CreateListPoliciesForGroupRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return response, errs.WrapError(err)
	}
	request.GroupName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.ListPoliciesForGroup(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	listPoliciesForGroupResponse, _ := raw.(*ramsdk.ListPoliciesForGroupResponse)
	if len(listPoliciesForGroupResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForGroupResponse.Policies.Policy {
			if v.PolicyType == parts[2] && (v.PolicyName == parts[1] || strings.ToLower(v.PolicyName) == strings.ToLower(parts[1])) {
				return &v, nil
			}
		}
	}
	return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *RamService) WaitForRamGroupPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeRamGroupPolicyAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.PolicyName, parts[1], errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountAlias(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	action := "GetAccountAlias"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(5*time.Minute), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *RamService) DescribeRamAccessKey(id, userName string) (*ramsdk.AccessKey, error) {
	key := &ramsdk.AccessKey{}
	request := ramsdk.CreateListAccessKeysRequest()
	request.RegionId = s.client.RegionId
	request.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.ListAccessKeys(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return key, errs.WrapErrorf(errs.NotFoundErr("RamAccessKey", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return key, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ramsdk.ListAccessKeysResponse)
	for _, accessKey := range response.AccessKeys.AccessKey {
		if accessKey.AccessKeyId == id {
			return &accessKey, nil
		}
	}
	return key, errs.WrapErrorf(errs.NotFoundErr("RamAccessKey", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
}

func (s *RamService) WaitForRamAccessKey(id, useName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamAccessKey(id, useName)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if string(status) == object.Status {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Status, status, errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "GetPolicy"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": id,
		"PolicyType": "Custom",
	}
	response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("RamPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RamService) DescribeRamRoleAttachment(id string) (*ecs.DescribeInstanceRamRoleResponse, error) {
	response := &ecs.DescribeInstanceRamRoleResponse{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return response, errs.WrapError(err)
	}
	request := ecs.CreateDescribeInstanceRamRoleRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = parts[1]
	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return retry.RetryableError(errs.WrapError(err))
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidRamRole.NotFound"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

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
			return response, nil
		}
	}
	return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *RamService) WaitForRamRoleAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRoleAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, strconv.Itoa(object.TotalCount), status, errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRole(id string) (*ramsdk.GetRoleResponse, error) {
	response := &ramsdk.GetRoleResponse{}
	request := ramsdk.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = id
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ramsdk.GetRoleResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *RamService) WaitForRamRole(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRole(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Role.RoleName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Role.RoleName, id, errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamUserPolicyAttachment(id string) (*ramsdk.Policy, error) {
	response := &ramsdk.Policy{}
	request := ramsdk.CreateListPoliciesForUserRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return response, errs.WrapError(err)
	}
	request.UserName = parts[3]
	var listPoliciesForUserResponse *ramsdk.ListPoliciesForUserResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listPoliciesForUserResponse, _ = raw.(*ramsdk.ListPoliciesForUserResponse)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForUserResponse.Policies.Policy {
			if v.PolicyType == parts[2] && (v.PolicyName == parts[1] || strings.ToLower(v.PolicyName) == strings.ToLower(parts[1])) {
				return &v, nil
			}
		}
	}
	return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *RamService) WaitForRamUserPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeRamUserPolicyAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.PolicyName, parts[1], errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRolePolicyAttachment(id string) (*ramsdk.Policy, error) {
	response := &ramsdk.Policy{}
	request := ramsdk.CreateListPoliciesForRoleRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return response, errs.WrapError(err)
	}
	request.RoleName = parts[3]
	var listPoliciesForRoleResponse *ramsdk.ListPoliciesForRoleResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListPoliciesForRole(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)

		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		listPoliciesForRoleResponse, _ = raw.(*ramsdk.ListPoliciesForRoleResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(listPoliciesForRoleResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForRoleResponse.Policies.Policy {
			if v.PolicyType == parts[2] && (v.PolicyName == parts[1] || strings.ToLower(v.PolicyName) == strings.ToLower(parts[1])) {
				return &v, nil
			}
		}
	}
	return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *RamService) WaitForRamRolePolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeRamRolePolicyAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {

				}
				return nil
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.PolicyName, parts[1], errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamGroup(id string) (*ramsdk.GetGroupResponse, error) {
	response := &ramsdk.GetGroupResponse{}
	request := ramsdk.CreateGetGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.GetGroup(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ramsdk.GetGroupResponse)

	if response.Group.GroupName != id {
		return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *RamService) WaitForRamGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Group.GroupName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Group.GroupName, id, errs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountPasswordPolicy(id string) (*ramsdk.GetPasswordPolicyResponse, error) {
	response := &ramsdk.GetPasswordPolicyResponse{}
	request := ramsdk.CreateGetPasswordPolicyRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
		return ramClient.GetPasswordPolicy(request)
	})
	if err != nil {
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*ramsdk.GetPasswordPolicyResponse)

	return response, nil
}

func (s *RamService) DescribeRamSecurityPreference(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetSecurityPreference"
	request := map[string]interface{}{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
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
	v, err := jsonpath.Get("$.SecurityPreference", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SecurityPreference", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RamService) DescribeRamServiceLinkedRole(id string) (*ramsdk.GetRoleResponse, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapErrorf(err, "invalid service linked role id: %s", id)
	}
	id = parts[1]

	response := &ramsdk.GetRoleResponse{}
	request := ramsdk.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = id
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ramsdk.GetRoleResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}
