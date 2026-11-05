package alikafka

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// alikafkaInstanceConfigDiffSuppressFunc suppresses the diff on the config attribute
// of alicloud_alikafka_instance: the service reports the effective configuration as a
// JSON object, so a declared key matches whenever its value is unchanged regardless of
// extra keys the service reports.
func alikafkaInstanceConfigDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return true
	}
	if old == "" {
		return false
	}

	oldMap := make(map[string]string)
	err := json.Unmarshal([]byte(old), &oldMap)
	if err != nil {
		return false
	}

	newMap := make(map[string]string)
	err = json.Unmarshal([]byte(new), &newMap)
	if err != nil {
		return false
	}

	// key exist in oldMap && found new value item different with old item
	for k, newValueItem := range newMap {
		oldValueItem, ok := oldMap[k]
		if ok && newValueItem != oldValueItem {
			return false
		}
	}

	return true
}

// alikafkaTopicConfigsDiffSuppressFunc suppresses the diff on the configs attribute of
// alicloud_alikafka_topic as long as every key declared in the configuration matches the
// value reported by the service. The service returns the whole effective topic configuration,
// which on a serverless instance always contains keys the user never declared
// (cloud.native.topic.type and replication-factor), so a whole-document comparison can never
// converge.
//
// Removing a key that was declared before is not reported as a diff: the update merges the
// submitted keys and cannot reset a single key. A key the service does not report yet is still
// reported as a diff so that it gets applied.
//
// helper.CompareJsonHasSubset does not fit here: it compares the decoded values with
// reflect.DeepEqual, so a number or a boolean declared in the configuration would never match the
// string the service reports for the same value. The values are compared through their literal
// text instead.
func alikafkaTopicConfigsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return true
	}
	if old == "" {
		return false
	}

	oldMap, err := decodeJsonObjectKeepingNumbers(old)
	if err != nil {
		log.Printf("[DEBUG] %s: the reported configuration is not a JSON object, keeping the diff: %s", k, err)
		return false
	}
	newMap, err := decodeJsonObjectKeepingNumbers(new)
	if err != nil {
		log.Printf("[DEBUG] %s: the declared configuration is not a JSON object, keeping the diff: %s", k, err)
		return false
	}

	for key, newValue := range newMap {
		oldValue, ok := oldMap[key]
		if !ok {
			log.Printf("[DEBUG] %s: the reported configuration does not contain the declared key %q, keeping the diff", k, key)
			return false
		}
		if formatJsonConfigValue(oldValue) != formatJsonConfigValue(newValue) {
			log.Printf("[DEBUG] %s: the declared key %q is %q while the service reports %q, keeping the diff",
				k, key, formatJsonConfigValue(newValue), formatJsonConfigValue(oldValue))
			return false
		}
	}

	return true
}

// decodeJsonObjectKeepingNumbers decodes a JSON object without converting numbers to float64,
// so their literal text stays available for comparison.
func decodeJsonObjectKeepingNumbers(s string) (map[string]interface{}, error) {
	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()

	object := make(map[string]interface{})
	if err := decoder.Decode(&object); err != nil {
		return nil, err
	}
	if err := decoder.Decode(new(interface{})); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("unexpected trailing data after the JSON object")
		}
		return nil, err
	}
	return object, nil
}

// formatJsonConfigValue renders a decoded JSON value as the text the API carries it in, so that
// values that only differ in their JSON type compare equal.
func formatJsonConfigValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case json.Number:
		return v.String()
	case bool:
		return strconv.FormatBool(v)
	default:
		encoded, _ := json.Marshal(v)
		return string(encoded)
	}
}

// kmsDiffSuppressFunc mirrors alicloud.kmsDiffSuppressFunc
// (diff_suppress_funcs.go:659). It suppresses the diff on a kms-encrypted
// password field whenever a cleartext password (or account_password) is set,
// so the provider does not force a replace when both forms are present. The
// alikafka SASL user resource (resource_sasl_user.go) uses it for
// kms_encrypted_password.
func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}
