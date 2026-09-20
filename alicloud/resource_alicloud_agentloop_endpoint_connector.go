// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAgentloopEndpointConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopEndpointConnectorCreate,
		Read:   resourceAliCloudAgentloopEndpointConnectorRead,
		Update: resourceAliCloudAgentloopEndpointConnectorUpdate,
		Delete: resourceAliCloudAgentloopEndpointConnectorDelete,
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
				ForceNew: true,
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

func resourceAliCloudAgentloopEndpointConnectorCreate(d *schema.ResourceData, meta interface{}) error {

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
		headersMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["value"] = dataLoopTmp["value"]
			dataLoopMap["key"] = dataLoopTmp["key"]
			headersMapsArray = append(headersMapsArray, dataLoopMap)
		}
		request["headers"] = headersMapsArray
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, maskEndpointConnectorCredentialForDebug(request))

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_endpoint_connector", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, response["connectorId"]))

	return resourceAliCloudAgentloopEndpointConnectorRead(d, meta)
}

func resourceAliCloudAgentloopEndpointConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopEndpointConnector(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_endpoint_connector DescribeAgentloopEndpointConnector Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alias", objectRaw["alias"])
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("description", objectRaw["description"])
	d.Set("endpoint", objectRaw["endpoint"])
	d.Set("name", objectRaw["name"])
	// GetEndpointConnector may omit regionId; the connector lives in the
	// region of the client endpoint (agentloop.<region>.aliyuncs.com), so
	// fall back to the client region when the API omits the field.
	if v, ok := objectRaw["regionId"]; ok && fmt.Sprint(v) != "" {
		d.Set("region_id", v)
	} else {
		d.Set("region_id", client.RegionId)
	}
	d.Set("type", objectRaw["type"])
	d.Set("updated_at", objectRaw["updatedAt"])
	d.Set("agent_space", objectRaw["agentSpace"])
	d.Set("connector_id", objectRaw["connectorId"])
	// The following attributes are intentionally NOT read back:
	// - credential: the Get API returns a desensitized shell that would
	//   overwrite the configured secret and cause a perpetual diff.
	// - properties: the Get API injects server-side default keys
	//   (channelType/qpsLimit/timeoutMs/maxRetries/...) and re-formats JSON
	//   string values such as modelList, both of which would diff.
	// - headers: the Get API does not return the configured headers.
	// - tags: the Get API does not return the configured tags.
	// They are listed in ImportStateVerifyIgnore of the acceptance test.

	return nil
}

func resourceAliCloudAgentloopEndpointConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
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
	connectorId := parts[1]
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", agentSpace, connectorId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("headers") {
		update = true
		if v, ok := d.GetOk("headers"); ok {
			headersMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["value"] = dataLoopTmp["value"]
				dataLoopMap["key"] = dataLoopTmp["key"]
				headersMapsArray = append(headersMapsArray, dataLoopMap)
			}
			request["headers"] = headersMapsArray
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
	// The update API requires a different credential shape than create: it
	// rejects a bare {"api_key": ...} map and demands providerType plus
	// sensitiveInfo as a JSON-encoded string whose apiKey member holds the
	// key material (probe-verified). Transform the create-style shorthand
	// accordingly; a credential that already carries providerType is passed
	// through unchanged.
	credential := d.Get("credential")
	if credentialMap, ok := credential.(map[string]interface{}); ok && credentialMap["providerType"] == nil {
		if apiKey, ok := credentialMap["api_key"]; ok {
			sensitiveInfo, _ := json.Marshal(map[string]interface{}{"apiKey": apiKey})
			credential = map[string]interface{}{
				"providerType":  "plain",
				"sensitiveInfo": string(sensitiveInfo),
			}
		}
	}
	request["credential"] = credential
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
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, maskEndpointConnectorCredentialForDebug(request))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudAgentloopEndpointConnectorRead(d, meta)
}

func resourceAliCloudAgentloopEndpointConnectorDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.EndpointConnector"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// maskEndpointConnectorCredentialForDebug returns a copy of the request map with
// the credential replaced, so addDebug never logs the plaintext secret (the
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
