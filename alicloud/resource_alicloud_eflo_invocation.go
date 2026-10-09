package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEfloInvocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloInvocationCreate,
		Read:   resourceAliCloudEfloInvocationRead,
		Update: resourceAliCloudEfloInvocationUpdate,
		Delete: resourceAliCloudEfloInvocationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"command_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_parameter": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"invocation_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RunCommand",
				ForceNew: true,
			},
			"launcher": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// lintignore: S006
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"repeat_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"termination_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"working_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEfloInvocationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	invocationType := d.Get("invocation_type").(string)
	if invocationType == "SendFile" {
		return resourceAliCloudEfloInvocationCreateSendFile(d, meta, client)
	}
	return resourceAliCloudEfloInvocationCreateRunCommand(d, meta, client)
}

func resourceAliCloudEfloInvocationCreateRunCommand(d *schema.ResourceData, meta interface{}, client *connectivity.AliyunClient) error {
	action := "RunCommand"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("content_encoding"); ok {
		request["ContentEncoding"] = v
	}
	if v, ok := d.GetOk("username"); ok {
		request["Username"] = v
	}
	if v, ok := d.GetOk("command_id"); ok {
		request["CommandId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("working_dir"); ok {
		request["WorkingDir"] = v
	}
	if v, ok := d.GetOk("command_content"); ok {
		request["CommandContent"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v
	}
	if v, ok := d.GetOk("node_id_list"); ok {
		nodeIdListMapsArray := v.([]interface{})
		nodeIdListMapsJson, err := json.Marshal(nodeIdListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["NodeIdList"] = string(nodeIdListMapsJson)
	}

	if v, ok := d.GetOkExists("enable_parameter"); ok {
		request["EnableParameter"] = v
	}
	if v, ok := d.GetOkExists("timeout"); ok {
		request["Timeout"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("termination_mode"); ok {
		request["TerminationMode"] = v
	}
	if v, ok := d.GetOk("frequency"); ok {
		request["Frequency"] = v
	}
	if v, ok := d.GetOk("launcher"); ok {
		request["Launcher"] = v
	}
	if v, ok := d.GetOk("repeat_mode"); ok {
		request["RepeatMode"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_invocation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InvokeId"]))

	return resourceAliCloudEfloInvocationUpdate(d, meta)
}

func resourceAliCloudEfloInvocationCreateSendFile(d *schema.ResourceData, meta interface{}, client *connectivity.AliyunClient) error {
	action := "SendFile"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("target_dir"); ok {
		request["TargetDir"] = v
	}
	if v, ok := d.GetOk("content_type"); ok {
		request["ContentType"] = v
	}
	if v, ok := d.GetOk("content"); ok {
		request["Content"] = v
	}
	if v, ok := d.GetOk("file_mode"); ok {
		request["FileMode"] = v
	}
	if v, ok := d.GetOk("file_owner"); ok {
		request["FileOwner"] = v
	}
	if v, ok := d.GetOk("file_group"); ok {
		request["FileGroup"] = v
	}
	if v, ok := d.GetOk("node_id_list"); ok {
		nodeIdListMapsArray := v.([]interface{})
		nodeIdListMapsJson, err := json.Marshal(nodeIdListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["NodeIdList"] = string(nodeIdListMapsJson)
	}
	if v, ok := d.GetOkExists("overwrite"); ok {
		request["Overwrite"] = v
	}
	if v, ok := d.GetOkExists("timeout"); ok {
		request["Timeout"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_invocation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InvokeId"]))

	return resourceAliCloudEfloInvocationRead(d, meta)
}

func resourceAliCloudEfloInvocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	invocationType := d.Get("invocation_type").(string)

	// SendFile 分支
	if invocationType == "SendFile" {
		object, err := efloServiceV2.DescribeEfloSendFileResults(d.Id())
		if err != nil {
			if !d.IsNewResource() && NotFoundError(err) {
				log.Printf("[DEBUG] Resource alicloud_eflo_invocation DescribeEfloSendFileResults Failed!!! %s", err)
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		setEfloSendFileAttributes(d, object)
		return nil
	}

	// RunCommand 分支（invocationType == "RunCommand" 或空即 import 场景）
	_, err := efloServiceV2.DescribeEfloInvocation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			// RunCommand 未找到；若 invocation_type 为空（import），尝试 SendFile
			if invocationType == "" {
				object, sendErr := efloServiceV2.DescribeEfloSendFileResults(d.Id())
				if sendErr == nil {
					d.Set("invocation_type", "SendFile")
					setEfloSendFileAttributes(d, object)
					return nil
				}
			}
			log.Printf("[DEBUG] Resource alicloud_eflo_invocation DescribeEfloInvocation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	return nil
}

// setEfloSendFileAttributes 把 DescribeSendFileResults 的 Invocation 记录字段反向映射到 state。
// NodeId 来自 InvokeNodes.InvokeNode[*]（取首个 InvokeNode）。
func setEfloSendFileAttributes(d *schema.ResourceData, object map[string]interface{}) {
	d.Set("content", object["Content"])
	d.Set("content_type", object["ContentType"])
	d.Set("description", object["Description"])
	d.Set("file_mode", object["FileMode"])
	d.Set("file_group", object["FileGroup"])
	d.Set("file_owner", object["FileOwner"])
	d.Set("name", object["Name"])
	d.Set("target_dir", object["TargetDir"])
	d.Set("overwrite", object["Overwrite"])
	nodeIdVal, err := jsonpath.Get("$.InvokeNodes.InvokeNode[*].NodeId", object)
	if err == nil && nodeIdVal != nil {
		if nodeIds, ok := nodeIdVal.([]interface{}); ok && len(nodeIds) > 0 {
			d.Set("node_id", fmt.Sprint(nodeIds[0]))
		}
	}
}

func resourceAliCloudEfloInvocationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	invocationType := d.Get("invocation_type").(string)
	// SendFile 是一次性下发，不支持 update；直接 Read
	if invocationType == "SendFile" {
		return resourceAliCloudEfloInvocationRead(d, meta)
	}

	// RunCommand 分支（现有 StopInvocation 逻辑）
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "StopInvocation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InvokeId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("node_id_list") {
		update = true
		if v, ok := d.GetOk("node_id_list"); ok || d.HasChange("node_id_list") {
			nodeIdListMapsArray := v.([]interface{})
			nodeIdListMapsJson, err := json.Marshal(nodeIdListMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["NodeIdList"] = string(nodeIdListMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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

	return resourceAliCloudEfloInvocationRead(d, meta)
}

func resourceAliCloudEfloInvocationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Invocation. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
