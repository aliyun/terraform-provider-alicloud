package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// slsIndexApiDefaults defines the known default values that SLS API
// automatically adds to each key's config when they are not specified by the user.
var slsIndexApiDefaults = map[string]interface{}{
	"alias":         "",
	"caseSensitive": false,
	"chn":           false,
}

// isApiDefaultValue checks if the given field value is a known API-added default.
func isApiDefaultValue(fieldName string, value interface{}) bool {
	defaultValue, exists := slsIndexApiDefaults[fieldName]
	if !exists {
		return false
	}
	return areFieldValuesEquivalent(value, defaultValue)
}

// areKeysJsonEquivalent compares two keys JSON strings semantically,
// ignoring API-added default values. Returns true if user config matches state.
func areKeysJsonEquivalent(old, new string) (bool, error) {
	// Handle empty string edge cases
	if old == "" && new == "" {
		return true, nil
	}
	if old == "" || new == "" {
		return false, nil
	}

	oldKeys := make(map[string]interface{})
	if err := json.Unmarshal([]byte(old), &oldKeys); err != nil {
		// Fallback to original comparison on parse error
		return compareJsonTemplateAreEquivalent(old, new)
	}

	newKeys := make(map[string]interface{})
	if err := json.Unmarshal([]byte(new), &newKeys); err != nil {
		return compareJsonTemplateAreEquivalent(old, new)
	}

	// Check if user added a new key
	for keyName, newValue := range newKeys {
		oldValue, exists := oldKeys[keyName]
		if !exists {
			return false, nil
		}

		newConfig, ok1 := newValue.(map[string]interface{})
		oldConfig, ok2 := oldValue.(map[string]interface{})
		if !ok1 || !ok2 {
			// Fallback to original comparison if type mismatch
			return compareJsonTemplateAreEquivalent(old, new)
		}

		// Check if user modified any field in the new config
		for fieldName, fieldValue := range newConfig {
			oldFieldValue, hasInOld := oldConfig[fieldName]
			if !hasInOld {
				return false, nil
			}
			if !areFieldValuesEquivalent(oldFieldValue, fieldValue) {
				return false, nil
			}
		}

		// Check old-only fields: suppress diff ONLY if they are known API defaults.
		// If a field exists in state but not in user config and is NOT a known API
		// default, it means the user deleted it or there is external drift — surface the diff.
		for oldFieldName, oldFieldValue := range oldConfig {
			if _, hasInNew := newConfig[oldFieldName]; !hasInNew {
				if !isApiDefaultValue(oldFieldName, oldFieldValue) {
					return false, nil
				}
			}
		}
	}

	// Check if user deleted a key
	for keyName := range oldKeys {
		if _, exists := newKeys[keyName]; !exists {
			return false, nil
		}
	}

	return true, nil
}

// areFieldValuesEquivalent compares two field values for equivalence.
// Numeric types (int/int64/float64) are compared loosely to handle JSON parsing.
func areFieldValuesEquivalent(old, new interface{}) bool {
	// Handle numeric types loosely - JSON unmarshals all numbers to float64
	switch old.(type) {
	case int, int64, float64:
		switch new.(type) {
		case int, int64, float64:
			return toFloat64(old) == toFloat64(new)
		}
	}

	// Strict type matching for other types
	if reflect.TypeOf(old) != reflect.TypeOf(new) {
		return false
	}

	switch oldVal := old.(type) {
	case string:
		return oldVal == new.(string)
	case bool:
		return oldVal == new.(bool)
	case []interface{}:
		newArr := new.([]interface{})
		if len(oldVal) != len(newArr) {
			return false
		}
		// Order-insensitive array comparison using element counting
		oldCount := make(map[string]int)
		for _, v := range oldVal {
			oldCount[fmt.Sprintf("%v", v)]++
		}
		newCount := make(map[string]int)
		for _, v := range newArr {
			newCount[fmt.Sprintf("%v", v)]++
		}
		if len(oldCount) != len(newCount) {
			return false
		}
		for k, count := range oldCount {
			if newCount[k] != count {
				return false
			}
		}
		return true
	case map[string]interface{}:
		newMap := new.(map[string]interface{})
		if len(oldVal) != len(newMap) {
			return false
		}
		for k, v := range oldVal {
			if !areFieldValuesEquivalent(v, newMap[k]) {
				return false
			}
		}
		return true
	default:
		return old == new
	}
}

// toFloat64 converts numeric types to float64 for unified comparison
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	default:
		return 0
	}
}

func resourceAliCloudSlsIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsIndexCreate,
		Read:   resourceAliCloudSlsIndexRead,
		Update: resourceAliCloudSlsIndexUpdate,
		Delete: resourceAliCloudSlsIndexDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"keys": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := areKeysJsonEquivalent(old, new)
					return equal
				},
			},
			"line": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"chn": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"case_sensitive": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"exclude_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"token": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"log_reduce": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_reduce_black_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_reduce_white_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_text_len": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudSlsIndexCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	logstore := d.Get("logstore_name")
	action := fmt.Sprintf("/logstores/%s/index", logstore)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project_name").(string))

	dataList := make(map[string]interface{})

	if v := d.Get("line"); !IsNil(v) {
		chn1, _ := jsonpath.Get("$[0].chn", v)
		if chn1 != nil && chn1 != "" {
			dataList["chn"] = chn1
		}
		caseSensitive1, _ := jsonpath.Get("$[0].case_sensitive", v)
		if caseSensitive1 != nil && caseSensitive1 != "" {
			dataList["caseSensitive"] = caseSensitive1
		}
		token1, _ := jsonpath.Get("$[0].token", v)
		if token1 != nil && token1 != "" {
			dataList["token"] = token1
		}
		excludeKeys, _ := jsonpath.Get("$[0].exclude_keys", v)
		if excludeKeys != nil && excludeKeys != "" && len(excludeKeys.([]interface{})) > 0 {
			dataList["exclude_keys"] = excludeKeys
		}
		includeKeys, _ := jsonpath.Get("$[0].include_keys", v)
		if includeKeys != nil && includeKeys != "" && len(includeKeys.([]interface{})) > 0 {
			dataList["include_keys"] = includeKeys
		}

		request["line"] = dataList
	}

	if v, ok := d.GetOkExists("log_reduce"); ok {
		request["log_reduce"] = v
	}
	if v, ok := d.GetOkExists("max_text_len"); ok {
		request["max_text_len"] = v
	}
	if v, ok := d.GetOk("log_reduce_white_list"); ok {
		log_reduce_white_listMapsArray := v.([]interface{})
		request["log_reduce_white_list"] = log_reduce_white_listMapsArray
	}

	if v, ok := d.GetOk("log_reduce_black_list"); ok {
		log_reduce_black_listMapsArray := v.([]interface{})
		request["log_reduce_black_list"] = log_reduce_black_listMapsArray
	}

	if v, ok := d.GetOk("keys"); ok {
		request["keys"] = NormalizeMap(convertJsonStringToObject(v))
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateIndex", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_index", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], logstore))

	return resourceAliCloudSlsIndexUpdate(d, meta)
}

func resourceAliCloudSlsIndexRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsIndex(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_index DescribeSlsIndex Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["keys"] != nil {
		keys, err := convertToJsonWithoutEscapeHTML(objectRaw["keys"].(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		d.Set("keys", keys)
	}
	d.Set("log_reduce", objectRaw["log_reduce"])
	d.Set("max_text_len", objectRaw["max_text_len"])

	lineMaps := make([]map[string]interface{}, 0)
	lineMap := make(map[string]interface{})
	lineRaw := make(map[string]interface{})
	if objectRaw["line"] != nil {
		lineRaw = objectRaw["line"].(map[string]interface{})
	}
	if len(lineRaw) > 0 {
		lineMap["case_sensitive"] = lineRaw["caseSensitive"]
		lineMap["chn"] = lineRaw["chn"]

		exclude_keysRaw := make([]interface{}, 0)
		if lineRaw["exclude_keys"] != nil {
			exclude_keysRaw = lineRaw["exclude_keys"].([]interface{})
		}

		lineMap["exclude_keys"] = exclude_keysRaw
		include_keysRaw := make([]interface{}, 0)
		if lineRaw["include_keys"] != nil {
			include_keysRaw = lineRaw["include_keys"].([]interface{})
		}

		lineMap["include_keys"] = include_keysRaw
		tokenRaw := make([]interface{}, 0)
		if lineRaw["token"] != nil {
			tokenRaw = lineRaw["token"].([]interface{})
		}

		lineMap["token"] = tokenRaw
		lineMaps = append(lineMaps, lineMap)
	}
	if err := d.Set("line", lineMaps); err != nil {
		return err
	}
	log_reduce_black_listRaw := make([]interface{}, 0)
	if objectRaw["log_reduce_black_list"] != nil {
		log_reduce_black_listRaw = objectRaw["log_reduce_black_list"].([]interface{})
	}

	d.Set("log_reduce_black_list", log_reduce_black_listRaw)
	log_reduce_white_listRaw := make([]interface{}, 0)
	if objectRaw["log_reduce_white_list"] != nil {
		log_reduce_white_listRaw = objectRaw["log_reduce_white_list"].([]interface{})
	}

	d.Set("log_reduce_white_list", log_reduce_white_listRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])
	d.Set("logstore_name", parts[1])

	return nil
}

func resourceAliCloudSlsIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	logstore := parts[1]
	action := fmt.Sprintf("/logstores/%s/index", logstore)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if !d.IsNewResource() && d.HasChange("line") {
		update = true
	}
	dataList := make(map[string]interface{})

	if v := d.Get("line"); v != nil {
		chn1, _ := jsonpath.Get("$[0].chn", v)
		if chn1 != nil && (d.HasChange("line.0.chn") || chn1 != "") {
			dataList["chn"] = chn1
		}
		caseSensitive1, _ := jsonpath.Get("$[0].case_sensitive", v)
		if caseSensitive1 != nil && (d.HasChange("line.0.case_sensitive") || caseSensitive1 != "") {
			dataList["caseSensitive"] = caseSensitive1
		}
		token1, _ := jsonpath.Get("$[0].token", d.Get("line"))
		if token1 != nil && (d.HasChange("line.0.token") || token1 != "") {
			dataList["token"] = token1
		}
		excludeKeys, _ := jsonpath.Get("$[0].exclude_keys", d.Get("line"))
		if excludeKeys != nil && excludeKeys != "" && len(excludeKeys.([]interface{})) > 0 {
			dataList["exclude_keys"] = excludeKeys
		}
		includeKeys, _ := jsonpath.Get("$[0].include_keys", d.Get("line"))
		if includeKeys != nil && includeKeys != "" && len(includeKeys.([]interface{})) > 0 {
			dataList["include_keys"] = includeKeys
		}

		request["line"] = dataList
	}

	if !d.IsNewResource() && d.HasChange("max_text_len") {
		update = true
	}
	if v, ok := d.GetOk("max_text_len"); ok || d.HasChange("max_text_len") {
		request["max_text_len"] = v
	}
	if !d.IsNewResource() && d.HasChange("log_reduce") {
		update = true
	}
	if v, ok := d.GetOk("log_reduce"); ok || d.HasChange("log_reduce") {
		request["log_reduce"] = v
	}
	if !d.IsNewResource() && d.HasChange("log_reduce_white_list") {
		update = true
	}
	if v, ok := d.GetOk("log_reduce_white_list"); ok || d.HasChange("log_reduce_white_list") {
		log_reduce_white_listMapsArray := v.([]interface{})
		request["log_reduce_white_list"] = log_reduce_white_listMapsArray
	}

	if !d.IsNewResource() && d.HasChange("log_reduce_black_list") {
		update = true
	}
	if v, ok := d.GetOk("log_reduce_black_list"); ok || d.HasChange("log_reduce_black_list") {
		log_reduce_black_listMapsArray := v.([]interface{})
		request["log_reduce_black_list"] = log_reduce_black_listMapsArray
	}

	if !d.IsNewResource() && d.HasChange("keys") {
		update = true
	}
	if v, ok := d.GetOk("keys"); ok || d.HasChange("keys") {
		request["keys"] = NormalizeMap(convertJsonStringToObject(v))
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateIndex", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudSlsIndexRead(d, meta)
}

func resourceAliCloudSlsIndexDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	logstore := parts[1]
	action := fmt.Sprintf("/logstores/%s/index", logstore)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteIndex", action), query, nil, nil, hostMap, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
