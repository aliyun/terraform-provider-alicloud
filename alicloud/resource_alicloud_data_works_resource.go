package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDataWorksResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDataWorksResourceCreate,
		Read:   resourceAlicloudDataWorksResourceRead,
		Update: resourceAlicloudDataWorksResourceUpdate,
		Delete: resourceAlicloudDataWorksResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"resource_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"script": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_command": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"script_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudDataWorksResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateResource"
	query := make(map[string]interface{})
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ProjectId"] = d.Get("project_id")
	if v, ok := d.GetOk("spec"); ok {
		request["Spec"] = v
	}
	if v, ok := d.GetOk("resource_file"); ok {
		request["ResourceFile"] = v
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"9990020002", "9990040003"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_resource", action, AlibabaCloudSdkGoERROR)
	}
	resourceId := fmt.Sprint(response["Id"])
	d.SetId(fmt.Sprint(request["ProjectId"], ":", resourceId))
	return resourceAlicloudDataWorksResourceRead(d, meta)
}

func resourceAlicloudDataWorksResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}
	object, err := dataWorksServiceV2.DescribeDataWorksResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_resource DescribeDataWorksResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	resourceRaw, err := jsonpath.Get("$.Resource", object)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.Resource", object)
	}
	resourceMap := make(map[string]interface{})
	if resourceRaw != nil {
		resourceMap, _ = resourceRaw.(map[string]interface{})
	}
	d.Set("resource_id", resourceMap["Id"])
	d.Set("spec", resourceMap["Spec"])
	d.Set("resource_name", resourceMap["Name"])
	d.Set("create_time", resourceMap["CreateTime"])
	d.Set("project_id", resourceMap["ProjectId"])
	d.Set("owner", resourceMap["Owner"])
	d.Set("modify_time", resourceMap["ModifyTime"])
	d.Set("type", resourceMap["Type"])
	d.Set("source_type", resourceMap["SourceType"])
	d.Set("source_path", resourceMap["SourcePath"])
	d.Set("target_type", resourceMap["TargetType"])
	d.Set("target_path", resourceMap["TargetPath"])
	if dataSourceRaw, ok := resourceMap["DataSource"].(map[string]interface{}); ok {
		d.Set("data_source", []map[string]interface{}{{
			"name": dataSourceRaw["Name"],
			"type": dataSourceRaw["Type"],
		}})
	}
	if scriptRaw, ok := resourceMap["Script"].(map[string]interface{}); ok {
		scriptMap := map[string]interface{}{
			"path":      scriptRaw["Path"],
			"script_id": scriptRaw["Id"],
		}
		if runtimeRaw, ok := scriptRaw["Runtime"].(map[string]interface{}); ok {
			scriptMap["runtime_command"] = runtimeRaw["Command"]
		}
		d.Set("script", []map[string]interface{}{scriptMap})
	}
	return nil
}

func resourceAlicloudDataWorksResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	projectId := parts[0]
	resourceId := parts[1]

	// UpdateResource: update spec and/or resource_file
	if d.HasChange("spec") || d.HasChange("resource_file") {
		action := "UpdateResource"
		query := make(map[string]interface{})
		request := make(map[string]interface{})
		request["Id"] = resourceId
		request["ProjectId"] = projectId
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOk("spec"); ok {
			request["Spec"] = v
		}
		if v, ok := d.GetOk("resource_file"); ok {
			request["ResourceFile"] = v
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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

	// RenameResource: update resource_name
	if d.HasChange("resource_name") {
		action := "RenameResource"
		query := make(map[string]interface{})
		request := make(map[string]interface{})
		request["Id"] = resourceId
		request["Name"] = d.Get("resource_name")
		request["ProjectId"] = projectId
		request["RegionId"] = client.RegionId
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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

	// MoveResource: update path
	if d.HasChange("path") {
		action := "MoveResource"
		query := make(map[string]interface{})
		request := make(map[string]interface{})
		request["Id"] = resourceId
		request["ProjectId"] = projectId
		request["RegionId"] = client.RegionId
		request["Path"] = d.Get("path")
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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

	return resourceAlicloudDataWorksResourceRead(d, meta)
}

func resourceAlicloudDataWorksResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteResource"
	query := make(map[string]interface{})
	request := make(map[string]interface{})
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ProjectId"] = parts[0]
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"5801488441564725318"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
