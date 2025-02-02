package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMseEngineNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMseEngineNamespaceCreate,
		Read:   resourceAlicloudMseEngineNamespaceRead,
		Update: resourceAlicloudMseEngineNamespaceUpdate,
		Delete: resourceAlicloudMseEngineNamespaceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace_show_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace_desc": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudMseEngineNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	var response map[string]interface{}
	action := "CreateEngineNamespace"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}

	if v, ok := d.GetOk("namespace_desc"); ok {
		request["Desc"] = v
	}
	var instanceId = d.Get("instance_id")
	var clusterId = d.Get("cluster_id").(string)
	if instanceId == nil {
		object, err := mseService.GetInstanceIdBYClusterId(clusterId)
		if err != nil {
			instanceId = object["InstanceId"]
		} else {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_engine_namespace", action, AlibabaCloudSdkGoERROR)
		}
	}
	request["Id"] = d.Get("namespace_id")
	request["Name"] = d.Get("namespace_show_name")
	request["InstanceId"] = instanceId

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_engine_namespace", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Id"]))

	return resourceAlicloudMseEngineNamespaceRead(d, meta)
}
func resourceAlicloudMseEngineNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := mseService.DescribeMseEngineNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_engine_namespace mseService.DescribeMseEngineNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	clusterObject, err := mseService.DescribeMseCluster(parts[0])
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_engine_namespace mseService.DescribeMseCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("namespace_id", object["Namespace"])
	d.Set("cluster_id", clusterObject["ClusterId"])
	d.Set("instance_id", parts[0])
	d.Set("namespace_show_name", object["NamespaceShowName"])
	d.Set("namespace_desc", object["NamespaceDesc"])
	d.Set("quota", object["Quota"])
	d.Set("config_count", object["ConfigCount"])
	d.Set("service_count", object["ServiceCount"])
	d.Set("source_type", object["SourceType"])
	d.Set("type", object["Type"])
	return nil
}
func resourceAlicloudMseEngineNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": parts[0],
		"ClusterId":  parts[0],
		"Id":         parts[1],
	}

	if d.HasChanges("namespace_show_name", "dsec") {
		update = true
		if v, ok := d.GetOk("namespace_show_name"); ok {
			request["Name"] = v
		}
		if v, ok := d.GetOk("namespace_desc"); ok {
			request["Desc"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("accept_language"); ok {
			request["AcceptLanguage"] = v
		}
		action := "UpdateEngineNamespace"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
	return resourceAlicloudMseEngineNamespaceRead(d, meta)
}
func resourceAlicloudMseEngineNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEngineNamespace"
	var response map[string]interface{}
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id":         parts[1],
		"InstanceId": parts[0],
	}

	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
	return nil
}
