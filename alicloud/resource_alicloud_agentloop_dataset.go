// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAgentloopDataset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopDatasetCreate,
		Read:   resourceAliCloudAgentloopDatasetRead,
		Update: resourceAliCloudAgentloopDatasetUpdate,
		Delete: resourceAliCloudAgentloopDatasetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9]*(_[a-z0-9]+)*$`), "must start with a lowercase letter and contain only lowercase letters, digits, and non-consecutive underscores"),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAgentloopDatasetCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/agentspace/%s/dataset", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("dataset_name"); ok {
		request["datasetName"] = v
	}

	schemaMap := make(map[string]interface{})
	for key, val := range d.Get("schema").(map[string]interface{}) {
		// Decode with UseNumber so integer literals beyond the float64 precision
		// range reach the API unrounded.
		indexKey, err := agentloopDecodeJSON(val.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_dataset", action, "Invalid schema."+key, AlibabaCloudSdkGoERROR)
		}
		schemaMap[key] = indexKey
	}
	request["schema"] = schemaMap
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			// InternalError is retried because the backend intermittently
			// returns 500 right after the hosting AgentSpace is created.
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict.DatasetBusy", "InternalError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_dataset", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["datasetName"]))

	return resourceAliCloudAgentloopDatasetRead(d, meta)
}

func resourceAliCloudAgentloopDatasetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopDataset(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_dataset DescribeAgentloopDataset Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("region_id", objectRaw["regionId"])
	if v, ok := objectRaw["schema"]; ok && v != nil {
		if schemaResp, ok := v.(map[string]interface{}); ok {
			// The state currently holds the user-supplied JSON text (on create it is
			// the configuration verbatim). Re-serializing the API value always
			// produces sorted keys without spaces, which would turn any hand-written
			// non-canonical JSON into a perpetual diff on this ForceNew attribute, so
			// keep the existing text whenever it is semantically equal to the API value.
			configured := make(map[string]interface{})
			if current, ok := d.GetOk("schema"); ok {
				configured = current.(map[string]interface{})
			}
			flatSchema := make(map[string]string, len(schemaResp))
			for key, val := range schemaResp {
				// The API injects a system-defined "agentloop_annotations" IndexKey into
				// the schema. Drop it so the state stays consistent with the configuration.
				if key == "agentloop_annotations" {
					continue
				}
				b, err := json.Marshal(removeDatasetSchemaNullFields(val))
				if err != nil {
					return WrapError(err)
				}
				flat := string(b)
				if orig, ok := configured[key].(string); ok && agentloopJsonStringsEquivalent(orig, flat) {
					flat = orig
				}
				flatSchema[key] = flat
			}
			d.Set("schema", flatSchema)
		} else {
			d.Set("schema", v)
		}
	}
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("agent_space", objectRaw["agentSpace"])
	d.Set("dataset_name", objectRaw["datasetName"])

	return nil
}

func resourceAliCloudAgentloopDatasetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	datasetName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/dataset/%s", agentSpace, datasetName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict.DatasetBusy"}) {
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

	d.Partial(false)
	return resourceAliCloudAgentloopDatasetRead(d, meta)
}

func resourceAliCloudAgentloopDatasetDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	datasetName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/dataset/%s", agentSpace, datasetName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict.DatasetBusy"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"DatasetNotExist", "WorkspaceNotExist", "AgentSpaceNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// agentloopJsonStringsEquivalent reports whether two strings carry the same
// JSON document, ignoring formatting such as key order and whitespace (e.g.
// `{"chn":true,"type":"text"}` vs `{"type":"text", "chn":true}`). A string
// that is not valid JSON only matches a byte-identical counterpart. Numbers
// compare by exact numeric value: `1`, `1.0` and `1e1` are equal, while
// integers beyond the float64 precision range (e.g. 9007199254740992 vs
// 9007199254740993) stay distinct instead of being silently rounded equal.
func agentloopJsonStringsEquivalent(a, b string) bool {
	if a == b {
		return true
	}
	aValue, err := agentloopDecodeJSON(a)
	if err != nil {
		return false
	}
	bValue, err := agentloopDecodeJSON(b)
	if err != nil {
		return false
	}
	return agentloopJSONValueEqual(aValue, bValue)
}

// agentloopDecodeJSON decodes a JSON document keeping number literals exact
// (json.Number) instead of rounding them through float64. The input must
// carry exactly one JSON value: a second value after the first one (for
// example `{"a":1} {"b":2}`) is rejected instead of being silently dropped.
func agentloopDecodeJSON(s string) (interface{}, error) {
	var v interface{}
	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	var extra interface{}
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("multiple JSON values in input")
		}
		return nil, err
	}
	return v, nil
}

// agentloopJSONValueEqual compares two documents decoded by agentloopDecodeJSON.
func agentloopJSONValueEqual(a, b interface{}) bool {
	switch aValue := a.(type) {
	case map[string]interface{}:
		bValue, ok := b.(map[string]interface{})
		if !ok || len(aValue) != len(bValue) {
			return false
		}
		for key, aItem := range aValue {
			bItem, ok := bValue[key]
			if !ok || !agentloopJSONValueEqual(aItem, bItem) {
				return false
			}
		}
		return true
	case []interface{}:
		bValue, ok := b.([]interface{})
		if !ok || len(aValue) != len(bValue) {
			return false
		}
		for i := range aValue {
			if !agentloopJSONValueEqual(aValue[i], bValue[i]) {
				return false
			}
		}
		return true
	case json.Number:
		bValue, ok := b.(json.Number)
		if !ok {
			return false
		}
		if aValue == bValue {
			return true
		}
		// big.Rat parses every legal JSON number (decimal and exponent forms)
		// exactly, so the comparison neither rounds through float64 nor
		// distinguishes numerically equal literals such as 1 and 1.0.
		aRat, aOk := new(big.Rat).SetString(aValue.String())
		bRat, bOk := new(big.Rat).SetString(bValue.String())
		if !aOk || !bOk {
			return false
		}
		return aRat.Cmp(bRat) == 0
	default:
		return reflect.DeepEqual(a, b)
	}
}

// removeDatasetSchemaNullFields strips null-valued keys from the IndexKey objects
// returned by the API, so the flattened JSON strings stay consistent with the
// configuration values produced by jsonencode.
func removeDatasetSchemaNullFields(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(t))
		for k, val := range t {
			if val == nil {
				continue
			}
			m[k] = removeDatasetSchemaNullFields(val)
		}
		return m
	case []interface{}:
		for i, val := range t {
			t[i] = removeDatasetSchemaNullFields(val)
		}
		return t
	default:
		return v
	}
}
