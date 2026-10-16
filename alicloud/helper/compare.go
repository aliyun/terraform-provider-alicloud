package helper

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"

	"gopkg.in/yaml.v2"
)

// Template and value comparison helpers moved from the root package
// (alicloud/common.go); package alicloud re-exports them as var aliases
// so the flat layout keeps compiling unchanged.

func CompareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	obj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	obj2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

// CompareJsonHasSubset reports whether every key in sub exists in sup with an
// equal value. Used for ForceNew JSON fields the server enriches with extra
// keys after create: the user config is a subset of the server-returned state,
// so the diff should be suppressed instead of forcing replacement.
func CompareJsonHasSubset(sup, sub string) (bool, error) {
	supObj := make(map[string]interface{})
	if err := json.Unmarshal([]byte(sup), &supObj); err != nil {
		return false, err
	}
	subObj := make(map[string]interface{})
	if err := json.Unmarshal([]byte(sub), &subObj); err != nil {
		return false, err
	}
	for k, v := range subObj {
		sv, ok := supObj[k]
		if !ok || !reflect.DeepEqual(sv, v) {
			return false, nil
		}
	}
	return true, nil
}

func CompareArrayJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	obj1 := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	obj2 := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}

	return equal, nil
}

func CompareYamlTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := yaml.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalYaml1, _ := yaml.Marshal(obj1)

	var obj2 interface{}
	err = yaml.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalYaml2, _ := yaml.Marshal(obj2)

	equal := bytes.Compare(canonicalYaml1, canonicalYaml2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalYaml1, canonicalYaml2)
	}
	return equal, nil
}

func CompareMapWithIgnoreEquivalent(tem1, tem2 map[string]interface{}, ignore []string) bool {

	if len(tem1) != len(tem2) {
		return false
	}

OuterLoop:
	for key1, val1 := range tem1 {
		for _, item := range ignore {
			if key1 == item {
				continue OuterLoop
			}
		}

		val2 := tem2[key1]
		if val2 != val1 {
			return false
		}
	}

	return true
}

func IsSlicesValuesEqual(left, right []interface{}) bool {
	if len(left) != len(right) {
		return false
	}

	leftMap := make(map[interface{}]struct{}, len(left))
	for _, v := range left {
		leftMap[v] = struct{}{}
	}

	for _, v := range right {
		if _, ok := leftMap[v]; !ok {
			return false
		}
	}

	return true
}
