// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaScheduledPreloadJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaScheduledPreloadJobCreate,
		Read:   resourceAliCloudEsaScheduledPreloadJobRead,
		Update: resourceAliCloudEsaScheduledPreloadJobUpdate,
		Delete: resourceAliCloudEsaScheduledPreloadJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insert_way": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"textBox", "oss"}, false),
			},
			"oss_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_preload_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_preload_job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"url_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEsaScheduledPreloadJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateScheduledPreloadJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("url_list"); ok {
		request["UrlList"] = v
	}
	if v, ok := d.GetOk("oss_url"); ok {
		request["OssUrl"] = v
	}
	request["InsertWay"] = d.Get("insert_way")
	request["Name"] = d.Get("scheduled_preload_job_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_scheduled_preload_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["SiteId"], response["Id"]))

	return resourceAliCloudEsaScheduledPreloadJobRead(d, meta)
}

func resourceAliCloudEsaScheduledPreloadJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaScheduledPreloadJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_scheduled_preload_job DescribeEsaScheduledPreloadJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedAt"])
	d.Set("insert_way", objectRaw["InsertWay"])
	d.Set("scheduled_preload_job_name", objectRaw["Name"])
	d.Set("scheduled_preload_job_id", objectRaw["Id"])
	d.Set("site_id", objectRaw["SiteId"])

	return nil
}

func resourceAliCloudEsaScheduledPreloadJobUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Scheduled Preload Job.")
	return nil
}

func resourceAliCloudEsaScheduledPreloadJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteScheduledPreloadJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

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
