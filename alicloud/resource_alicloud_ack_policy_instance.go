// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/blues/jsonata-go"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAckPolicyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAckPolicyInstanceCreate,
		Read:   resourceAliCloudAckPolicyInstanceRead,
		Update: resourceAliCloudAckPolicyInstanceUpdate,
		Delete: resourceAliCloudAckPolicyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Optional: true,
				//Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"warn", "deny"}, false),
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAckPolicyInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	cluster_id := d.Get("cluster_id")
	policy_name := d.Get("policy_name")
	action := fmt.Sprintf("/clusters/%s/policies/%s", cluster_id, policy_name)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesMapsArray := convertToInterfaceArray(v)

		request["namespaces"] = namespacesMapsArray
	}

	if v, ok := d.GetOk("policy_action"); ok {
		request["action"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		parametersObject, err := convertYamlToObject(v)
		if err != nil {
			return WrapError(err)
		}

		request["parameters"] = convertObjectToJsonString(parametersObject)
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("CS", "2015-12-15", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ack_policy_instance", action, AlibabaCloudSdkGoERROR)
	}

	instancesVar, _ := jsonpath.Get("$.instances[0]", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", cluster_id, policy_name, instancesVar))

	return resourceAliCloudAckPolicyInstanceRead(d, meta)
}

func resourceAliCloudAckPolicyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackServiceV2 := AckServiceV2{client}

	objectRaw, err := ackServiceV2.DescribeAckPolicyInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ack_policy_instance DescribeAckPolicyInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_action", objectRaw["policy_action"])
	d.Set("cluster_id", objectRaw["cluster_id"])
	d.Set("instance_name", objectRaw["instance_name"])
	d.Set("policy_name", objectRaw["policy_name"])

	e := jsonata.MustCompile("$split($.policy_scope, \",\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("namespaces", evaluation)
	e = jsonata.MustCompile("$merge($map($.ApiOutput, function($v, $k) {{\"repos\":$split($replace($replace($.policy_parameters, \"repos:\n- \", \"\"), \"\n\", \"\"), \"- \")}}))")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("parameters", evaluation)

	return nil
}

func resourceAliCloudAckPolicyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	cluster_id := parts[0]
	policy_name := parts[1]
	action := fmt.Sprintf("/clusters/%s/policies/%s", cluster_id, policy_name)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["instance_name"] = parts[2]

	if d.HasChange("namespaces") {
		update = true
	}
	if v, ok := d.GetOk("namespaces"); ok || d.HasChange("namespaces") {
		namespacesMapsArray := convertToInterfaceArray(v)

		request["namespaces"] = namespacesMapsArray
	}

	if d.HasChange("policy_action") {
		update = true
	}
	if v, ok := d.GetOk("policy_action"); ok || d.HasChange("policy_action") {
		request["action"] = v
	}
	if d.HasChange("parameters") {
		update = true
	}
	if v, ok := d.GetOk("parameters"); ok || d.HasChange("parameters") {
		parametersObject, err := convertYamlToObject(v)
		if err != nil {
			return WrapError(err)
		}

		request["parameters"] = convertObjectToJsonString(parametersObject)
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("CS", "2015-12-15", action, query, nil, body, true)
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

	return resourceAliCloudAckPolicyInstanceRead(d, meta)
}

func resourceAliCloudAckPolicyInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	cluster_id := parts[0]
	policy_name := parts[1]
	action := fmt.Sprintf("/clusters/%s/policies/%s", cluster_id, policy_name)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["instance_name"] = StringPointer(parts[2])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("CS", "2015-12-15", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
