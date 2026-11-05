package mse

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceMseEngineNamespace() *schema.Resource {
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
			"quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"config_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudMseEngineNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := NewService(client)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mse_engine_namespace", action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	request["Id"] = d.Get("namespace_id")
	request["Name"] = d.Get("namespace_show_name")
	request["InstanceId"] = instanceId

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mse_engine_namespace", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Id"]))

	return resourceAlicloudMseEngineNamespaceRead(d, meta)
}
func resourceAlicloudMseEngineNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := NewService(client)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := mseService.DescribeMseEngineNamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_engine_namespace mseService.DescribeMseEngineNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	clusterObject, err := mseService.DescribeMseCluster(parts[0])
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_engine_namespace mseService.DescribeMseCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("namespace_id", object["Namespace"])
	d.Set("cluster_id", clusterObject["ClusterId"])
	d.Set("instance_id", parts[0])
	d.Set("namespace_show_name", object["NamespaceShowName"])
	d.Set("namespace_desc", object["NamespaceDesc"])
	quotaVal, err := convert.FormatInt(object["Quota"])
	if err != nil {
		return err
	}
	d.Set("quota", quotaVal)
	config_countVal, err := convert.FormatInt(object["ConfigCount"])
	if err != nil {
		return err
	}
	d.Set("config_count", config_countVal)
	d.Set("service_count", object["ServiceCount"])
	d.Set("source_type", object["SourceType"])
	typeVal, err := convert.FormatInt(object["Type"])
	if err != nil {
		return err
	}
	d.Set("type", typeVal)
	return nil
}
func resourceAlicloudMseEngineNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudMseEngineNamespaceRead(d, meta)
}
func resourceAlicloudMseEngineNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEngineNamespace"
	var response map[string]interface{}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"Id":         parts[1],
		"InstanceId": parts[0],
	}

	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
