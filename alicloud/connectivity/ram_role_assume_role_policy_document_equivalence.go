package connectivity

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

type ConditionMap map[string]map[string]interface{}
type PrincipalMap map[string]interface{}

var validPrincipalKeys = [...]string{"RAM", "SERVICE", "FEDERATED"}

type intermediatePolicy struct {
	Version   string
	Statement interface{}
}

type policyStatement struct {
	Effect    string
	Action    interface{}
	NotAction interface{}
	Principal PrincipalMap
	Condition ConditionMap
}

type policy struct {
	Version   string
	Statement []*policyStatement
}

func AssumeRolePolicyDocumentAreEquivalentV2(policy1Str, policy2Str string) (bool, error) {
	var policy1intermediate, policy2intermediate intermediatePolicy

	err := json.Unmarshal([]byte(policy1Str), &policy1intermediate)
	if err != nil {
		return false, fmt.Errorf("Unmarshal policy1 failed: %w", err)
	}

	err = json.Unmarshal([]byte(policy2Str), &policy2intermediate)
	if err != nil {
		return false, fmt.Errorf("Unmarshal policy2 failed: %w", err)
	}

	if reflect.DeepEqual(policy1intermediate, policy2intermediate) {
		return true, nil
	}

	policy1, err := policy1intermediate.document()
	if err != nil {
		return false, fmt.Errorf("Parsing policy1: %s", err)
	}

	policy2, err := policy2intermediate.document()
	if err != nil {
		return false, fmt.Errorf("Parsing policy2: %s", err)
	}

	for i := range policy1.Statement {
		if err := normalizePrincipal(&policy1.Statement[i].Principal); err != nil {
			return false, err
		}
	}

	for i := range policy2.Statement {
		if err := normalizePrincipal(&policy2.Statement[i].Principal); err != nil {
			return false, err
		}
	}

	return policy1.Equals(policy2), nil
}

func (intermediate *intermediatePolicy) document() (*policy, error) {
	var statement []*policyStatement
	if intermediate.Statement != nil {
		switch s := intermediate.Statement.(type) {
		case []interface{}:
			config := &mapstructure.DecoderConfig{
				Result:      &statement,
				ErrorUnused: true,
			}
			decoder, err := mapstructure.NewDecoder(config)
			if err != nil {
				return nil, err
			}
			err = decoder.Decode(s)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("Policy syntax error in 'Statement' field")
		}
	}

	policy := &policy{
		Version:   intermediate.Version,
		Statement: statement,
	}

	return policy, nil
}

func normalizePrincipal(principal *PrincipalMap) error {
	newPrincipal := make(PrincipalMap)
	for k, v := range *principal {
		upperK := strings.ToUpper(k)
		if !isValidPrincipalKey(upperK) {
			return errors.New(fmt.Sprintf("Invalid principal %s", upperK))
		}
		switch v := v.(type) {
		case string:
			newPrincipal[upperK] = []string{v}
		case []interface{}:
			strArray := make([]string, len(v))
			for i, val := range v {
				strArray[i] = val.(string)
			}
			newPrincipal[upperK] = strArray
		default:
			return errors.New(fmt.Sprintf("Invalid principal %s", k))
		}
	}
	*principal = newPrincipal
	return nil
}

func isValidPrincipalKey(key string) bool {
	for _, validKey := range validPrincipalKeys {
		if key == validKey {
			return true
		}
	}
	return false
}

func (p *policy) Equals(other *policy) bool {
	if p.Version != other.Version || len(p.Statement) != len(other.Statement) {
		return false
	}

	for i := range p.Statement {
		if !p.Statement[i].Equals(other.Statement[i]) {
			return false
		}
	}

	return true
}

func (s *policyStatement) Equals(other *policyStatement) bool {
	if !equalActionOrNotAction(s.Action, other.Action) || !equalActionOrNotAction(s.NotAction, other.NotAction) || s.Effect != other.Effect || !s.Principal.Equals(other.Principal) || !s.Condition.Equals(other.Condition) {
		return false
	}
	return true
}

func equalActionOrNotAction(a, b interface{}) bool {
	if a == nil {
		return b == nil
	}
	switch aTyped := a.(type) {
	case string:
		if bTyped, ok := b.(string); ok && aTyped == bTyped {
			return true
		}
	case []string:
		if bTyped, ok := b.([]string); ok && len(aTyped) == len(bTyped) {
			for i, v := range aTyped {
				if v != bTyped[i] {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (pm PrincipalMap) Equals(other PrincipalMap) bool {
	if len(pm) != len(other) {
		return false
	}

	for key, value := range pm {
		if !compareValues(value, other[key]) {
			return false
		}
	}

	return true
}

func (cm ConditionMap) Equals(other ConditionMap) bool {
	if cm == nil {
		return other == nil
	}

	if len(cm) != len(other) {
		return false
	}

	for key, value := range cm {
		if !compareValues(value, other[key]) {
			return false
		}
	}

	return true
}

func compareValues(a, b interface{}) bool {
	switch a := a.(type) {
	case string:
		bStr, ok := b.(string)
		return ok && a == bStr
	case []string:
		bArr, ok := b.([]string)
		return ok && len(a) == len(bArr) && equalStrings(a, bArr)
	case []interface{}:
		bArr, ok := b.([]interface{})
		return ok && len(a) == len(bArr) && equalInterfaces(a, bArr)
	case map[string]interface{}:
		bMap, ok := b.(map[string]interface{})
		if !ok || len(a) != len(bMap) {
			return false
		}
		for key, value := range a {
			if !compareValues(value, bMap[key]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func equalStrings(a, b []string) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalInterfaces(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !compareValues(v, b[i]) {
			return false
		}
	}
	return true
}
