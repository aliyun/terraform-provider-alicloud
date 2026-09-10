// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDasSqlLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDasSqlLogConfigCreate,
		Read:   resourceAliCloudDasSqlLogConfigRead,
		Update: resourceAliCloudDasSqlLogConfigUpdate,
		Delete: resourceAliCloudDasSqlLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cold_retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"hot_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sql_log_visible_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDasSqlLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	action := "ModifySqlLogConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	instanceId := d.Get("instance_id").(string)
	request["InstanceId"] = instanceId

	if v, ok := d.GetOkExists("enable"); ok {
		request["Enable"] = v
	}
	if v, ok := d.GetOkExists("request_enable"); ok {
		request["RequestEnable"] = v
	}
	if v, ok := d.GetOkExists("retention"); ok {
		request["Retention"] = v
	}
	if v, ok := d.GetOkExists("hot_retention"); ok {
		request["HotRetention"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("DAS", "2020-01-16", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_das_sql_log_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(instanceId)

	return resourceAliCloudDasSqlLogConfigRead(d, meta)
}

func resourceAliCloudDasSqlLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasServiceV2 := DasServiceV2{client}

	objectRaw, err := dasServiceV2.DescribeDasSqlLogConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_das_sql_log_config DescribeDasSqlLogConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	dataRawObj, _ := jsonpath.Get("$.Data", objectRaw)
	dataRaw := make(map[string]interface{})
	if dataRawObj != nil {
		dataRaw = dataRawObj.(map[string]interface{})
	}
	d.Set("instance_id", d.Id())
	d.Set("cold_retention", dataRaw["ColdRetention"])
	d.Set("hot_retention", dataRaw["HotRetention"])
	d.Set("log_filter", dataRaw["LogFilter"])
	d.Set("request_enable", dataRaw["RequestEnable"])
	d.Set("retention", dataRaw["Retention"])
	d.Set("sql_log_visible_time", dataRaw["SqlLogVisibleTime"])
	d.Set("version", dataRaw["Version"])
	d.Set("enable", dataRaw["SqlLogEnable"])

	return nil
}

func resourceAliCloudDasSqlLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifySqlLogConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	if d.HasChange("enable") {
		update = true
	}
	if d.HasChange("request_enable") {
		update = true
	}
	if d.HasChange("retention") {
		update = true
	}
	if d.HasChange("hot_retention") {
		update = true
	}

	if update {
		if v, ok := d.GetOkExists("enable"); ok {
			request["Enable"] = v
		}
		if v, ok := d.GetOkExists("request_enable"); ok {
			request["RequestEnable"] = v
		}
		if v, ok := d.GetOkExists("retention"); ok {
			request["Retention"] = v
		}
		if v, ok := d.GetOkExists("hot_retention"); ok {
			request["HotRetention"] = v
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DAS", "2020-01-16", action, query, request, true)
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

	return resourceAliCloudDasSqlLogConfigRead(d, meta)
}

func resourceAliCloudDasSqlLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Sql Log Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
