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

func resourceAliCloudEsaVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaVersionCreate,
		Read:   resourceAliCloudEsaVersionRead,
		Update: resourceAliCloudEsaVersionUpdate,
		Delete: resourceAliCloudEsaVersionDelete,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"origin_version": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaVersionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CloneVersion"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOkExists("origin_version"); ok {
		request["SiteVersion"] = v
	}
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	request["RegionId"] = client.RegionId

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_version", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["CloneVersion"]))

	return resourceAliCloudEsaVersionUpdate(d, meta)
}

func resourceAliCloudEsaVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_version DescribeEsaVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("status", objectRaw["Status"])
	d.Set("site_version", objectRaw["SiteVersion"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprint(parts[0]))

	return nil
}

func resourceAliCloudEsaVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateVersionDesc"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteVersion"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudEsaVersionRead(d, meta)
}

func resourceAliCloudEsaVersionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteVersion"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteVersion"] = parts[1]
	request["SiteId"] = parts[0]

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
