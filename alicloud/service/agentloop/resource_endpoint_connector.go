package agentloop

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEndpointConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointConnectorCreate,
		Read:   resourceEndpointConnectorRead,
		Update: resourceEndpointConnectorUpdate,
		Delete: resourceEndpointConnectorDelete,
		Importer: &schema.ResourceImporter{
			State: agentloopImportStatePassthrough(2),
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
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connector_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"credential": {
				Type:      schema.TypeMap,
				Required:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceEndpointConnectorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("headers"); ok {
		request["headers"] = expandAgentloopEndpointConnectorHeaders(v)
	}

	if v, ok := d.GetOk("alias"); ok {
		request["alias"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request["tags"] = convertToInterfaceArray(v)
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["credential"] = d.Get("credential")
	request["name"] = d.Get("name")
	if v, ok := d.GetOk("properties"); ok {
		request["properties"] = v
	}
	request["type"] = d.Get("type")
	request["endpoint"] = d.Get("endpoint")
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) || isExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, maskEndpointConnectorCredentialForDebug(request))

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_agentloop_endpoint_connector", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, response["connectorId"]))

	return resourceEndpointConnectorRead(d, meta)
}

func resourceEndpointConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopService := NewService(client)

	objectRaw, err := agentloopService.DescribeEndpointConnector(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_endpoint_connector DescribeEndpointConnector Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if err := d.Set("alias", objectRaw["alias"]); err != nil {
		return err
	}
	if err := d.Set("created_at", objectRaw["createdAt"]); err != nil {
		return err
	}
	if err := d.Set("description", objectRaw["description"]); err != nil {
		return err
	}
	if err := d.Set("endpoint", objectRaw["endpoint"]); err != nil {
		return err
	}
	if err := d.Set("name", objectRaw["name"]); err != nil {
		return err
	}
	// GetEndpointConnector may omit regionId; the connector lives in the
	// region of the client endpoint (agentloop.<region>.aliyuncs.com), so
	// fall back to the client region when the API omits the field.
	if v, ok := objectRaw["regionId"]; ok && fmt.Sprint(v) != "" {
		if err := d.Set("region_id", v); err != nil {
			return err
		}
	} else {
		if err := d.Set("region_id", client.RegionId); err != nil {
			return err
		}
	}
	if err := d.Set("type", objectRaw["type"]); err != nil {
		return err
	}
	if err := d.Set("updated_at", objectRaw["updatedAt"]); err != nil {
		return err
	}
	if err := d.Set("agent_space", objectRaw["agentSpace"]); err != nil {
		return err
	}
	if err := d.Set("connector_id", objectRaw["connectorId"]); err != nil {
		return err
	}
	headersMaps := make([]map[string]interface{}, 0)
	for _, item := range convertToInterfaceArray(objectRaw["headers"]) {
		if header, ok := item.(map[string]interface{}); ok {
			headersMaps = append(headersMaps, map[string]interface{}{
				"key":   header["key"],
				"value": header["value"],
			})
		}
	}
	if err := d.Set("headers", headersMaps); err != nil {
		return err
	}
	configuredProperties, _ := d.Get("properties").(map[string]interface{})
	connectorType, _ := d.Get("type").(string)
	if err := d.Set("properties", flattenAgentloopEndpointConnectorProperties(objectRaw["properties"], configuredProperties, connectorType)); err != nil {
		return err
	}
	if objectRaw["tags"] != nil {
		if err := d.Set("tags", convertToInterfaceArray(objectRaw["tags"])); err != nil {
			return err
		}
	}
	// credential is intentionally NOT read back: the Get API returns a
	// desensitized shell that would overwrite the configured secret.

	return nil
}

// agentloopEndpointConnectorDefaults returns the backend-injected default
// properties for the given connector type.  The qpsLimit default differs:
// agent_app → "20", all others (model_service etc.) → "100".
func agentloopEndpointConnectorDefaults(connectorType string) map[string]string {
	qps := "100"
	if connectorType == "agent_app" {
		qps = "20"
	}
	return map[string]string{
		"channelType": "custom",
		"maxRetries":  "3",
		"qpsLimit":    qps,
		"timeoutMs":   "300000",
	}
}

// flattenAgentloopEndpointConnectorProperties reads properties back while
// tolerating the backend's normalization. GetEndpointConnector returns every
// key it supports (unset ones as null) plus injected defaults, and drops keys
// it does not support:
//   - a configured key absent from the response is not stored by the backend,
//     so the configured value is kept;
//   - a configured key returned as null was removed outside Terraform and is
//     dropped so the drift shows up;
//   - an unconfigured key is read back unless it is null or still equals the
//     injected default;
//   - JSON values keep the configured text when semantically equal.
func flattenAgentloopEndpointConnectorProperties(raw interface{}, configured map[string]interface{}, connectorType string) map[string]interface{} {
	result := make(map[string]interface{})
	rawMap, _ := raw.(map[string]interface{})
	defaults := agentloopEndpointConnectorDefaults(connectorType)
	for k, v := range configured {
		if _, returned := rawMap[k]; !returned {
			result[k] = v
		}
	}
	for k, v := range rawMap {
		if v == nil {
			continue
		}
		flat, ok := v.(string)
		if !ok {
			bs, err := json.Marshal(v)
			if err != nil {
				continue
			}
			flat = string(bs)
		}
		orig, managed := configured[k].(string)
		if !managed {
			if def, ok := defaults[k]; ok && def == flat {
				continue
			}
		} else if orig != flat && agentloopJsonStringsEquivalent(orig, flat) {
			flat = orig
		}
		result[k] = flat
	}
	return result
}

func resourceEndpointConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	connectorId := parts[1]
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", agentSpace, connectorId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("headers") {
		update = true
		if v, ok := d.GetOk("headers"); ok {
			request["headers"] = expandAgentloopEndpointConnectorHeaders(v)
		} else {
			request["headers"] = make([]interface{}, 0)
		}
	}

	if d.HasChange("alias") {
		update = true
	}
	if v, ok := d.GetOk("alias"); ok || d.HasChange("alias") {
		request["alias"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("credential") {
		update = true
	}
	// UpdateEndpointConnector accepts the persisted credential envelope and
	// rejects the create-time api_key shorthand with "credential.providerType is
	// required for credential update". Normalize the shorthand for credential
	// rotation; an existing envelope passes through.
	request["credential"] = expandAgentloopEndpointConnectorUpdateCredential(d.Get("credential"))
	if d.HasChange("name") {
		update = true
	}
	request["name"] = d.Get("name")
	if d.HasChange("endpoint") {
		update = true
	}
	request["endpoint"] = d.Get("endpoint")
	if d.HasChange("properties") {
		update = true
	}
	if v, ok := d.GetOk("properties"); ok || d.HasChange("properties") {
		request["properties"] = v
	}
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			request["tags"] = convertToInterfaceArray(v)
		} else {
			request["tags"] = make([]interface{}, 0)
		}
	}
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) || isExpectedErrors(err, []string{"Conflict."}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, maskEndpointConnectorCredentialForDebug(request))
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceEndpointConnectorRead(d, meta)
}

func resourceEndpointConnectorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	connectorId := parts[1]
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", agentSpace, connectorId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) || isExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if isExpectedErrors(err, []string{"NotFound.EndpointConnector"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func expandAgentloopEndpointConnectorUpdateCredential(value interface{}) interface{} {
	credential, ok := value.(map[string]interface{})
	if !ok || credential["providerType"] != nil {
		return value
	}
	apiKey, ok := credential["api_key"]
	if !ok {
		return value
	}
	sensitiveInfo, _ := json.Marshal(map[string]interface{}{"apiKey": apiKey})
	return map[string]interface{}{
		"providerType":  "plain",
		"sensitiveInfo": string(sensitiveInfo),
	}
}

func expandAgentloopEndpointConnectorHeaders(value interface{}) []interface{} {
	headers := make([]interface{}, 0)
	for _, raw := range convertToInterfaceArray(value) {
		header, ok := raw.(map[string]interface{})
		if !ok {
			// An empty `headers {}` block is represented as a nil list element.
			continue
		}
		headers = append(headers, map[string]interface{}{
			"key":   header["key"],
			"value": header["value"],
		})
	}
	return headers
}

// maskEndpointConnectorCredentialForDebug returns a copy of the request map with
// the credential replaced, so logs.AddDebug never logs the plaintext secret (the
// create/update requests carry the full api key material).
func maskEndpointConnectorCredentialForDebug(request map[string]interface{}) map[string]interface{} {
	if request == nil || request["credential"] == nil {
		return request
	}
	masked := make(map[string]interface{}, len(request))
	for k, v := range request {
		masked[k] = v
	}
	masked["credential"] = "***"
	return masked
}
