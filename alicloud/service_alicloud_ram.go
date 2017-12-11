package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/denverdino/aliyungo/ram"
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

func ParseRolePolicyDocument(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, err
	}
	return policy, nil
}

func ParsePolicyDocument(policyDocument string) (statement []map[string]interface{}, version string, err error) {
	policy := Policy{}
	err = json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
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

func AssembleRolePolicyDocument(ramUser, service []interface{}, version string) (string, error) {
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
		return "", err
	}
	return string(data), nil
}

func AssemblePolicyDocument(document []interface{}, version string) (string, error) {
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
		return "", err
	}
	return string(data), nil
}

// Judge whether the role policy contains service "ecs.aliyuncs.com"
func (client *AliyunClient) JudgeRolePolicyPrincipal(roleName string) error {
	conn := client.ramconn
	resp, err := conn.GetRole(ram.RoleQueryRequest{RoleName: roleName})
	if err != nil {
		return fmt.Errorf("GetRole %s got an error: %#v", roleName, err)
	}

	policy, err := ParseRolePolicyDocument(resp.Role.AssumeRolePolicyDocument)
	if err != nil {
		return err
	}
	for _, v := range policy.Statement {
		for _, val := range v.Principal.Service {
			if strings.Trim(val, " ") == "ecs.aliyuncs.com" {
				return nil
			}
		}
	}
	return fmt.Errorf("Role policy services must contains 'ecs.aliyuncs.com', Now is \n%v.", resp.Role.AssumeRolePolicyDocument)
}

func GetIntersection(dataMap []map[string]interface{}, allDataMap map[string]interface{}) (allData []interface{}) {
	if len(dataMap) == 1 {
		allDataMap = dataMap[0]
	} else {
		for _, v := range dataMap {
			if len(v) > 0 {
				for key := range allDataMap {
					if _, ok := v[key]; !ok {
						allDataMap[key] = nil
					}
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
