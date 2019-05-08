package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
	request.RoleName = roleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, roleName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
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

func (s *RamService) GetSpecifiedUser(id string) (*ram.User, error) {
	request := ram.CreateGetUserRequest()
	request.UserName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetUser(request)
	})
	if err == nil {
		users, _ := raw.(*ram.GetUserResponse)
		return &users.User, nil
	}
	request1 := ram.CreateListUsersRequest()
	raw, err = s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListUsers(request1)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request1.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	users, _ := raw.(*ram.ListUsersResponse)
	for _, user := range users.Users.User {
		if user.UserId == id {
			return &user, nil
		}
	}
	return nil, WrapErrorf(Error(GetNotFoundMessage("ram_user", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
}
