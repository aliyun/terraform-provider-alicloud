package esa

import (
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

func ResourceSiteEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceSiteEnvironmentCreate,
		Read:   ResourceSiteEnvironmentRead,
		Update: ResourceSiteEnvironmentUpdate,
		Delete: ResourceSiteEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"next_environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceSiteEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateEnvironment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["EnvironmentName"] = d.Get("environment_name")
	request["SiteId"] = d.Get("site_id")
	request["Rule"] = d.Get("rule")
	request["NextEnvironmentName"] = d.Get("next_environment_name")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "SiteVersion.Configuring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_esa_site_environment", action, errs.AlibabaCloudSdkGoERROR)
	}

	if v, ok := d.GetOk("read_only"); ok {
		var err error
		action := "UpdateEnvironment"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["SiteId"] = d.Get("site_id")
		request["EnvironmentName"] = d.Get("environment_name")
		request["Rule"] = d.Get("rule")
		request["ReadOnly"] = v

		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "SiteVersion.Configuring"}) || helper.NeedRetry(err) {
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

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], request["EnvironmentName"]))

	return ResourceSiteEnvironmentRead(d, meta)
}

func ResourceSiteEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaService := NewService(client)

	objectRaw, err := esaService.DescribeEsaSiteEnvironment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_site_environment DescribeEsaSiteEnvironment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("environment_name", objectRaw["EnvironmentName"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("next_environment_name", d.Get("next_environment_name"))
	d.Set("read_only", objectRaw["ReadOnly"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("update_time", objectRaw["UpdateTime"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", parts[0])

	return nil
}

func ResourceSiteEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateEnvironment"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["EnvironmentName"] = parts[1]

	if d.HasChange("read_only") {
		update = true
	}
	request["ReadOnly"] = d.Get("read_only")

	if d.HasChange("rule") {
		update = true
	}
	request["Rule"] = d.Get("rule")

	if update {
		wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests", "SiteVersion.Configuring"}) || helper.NeedRetry(err) {
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

	return ResourceSiteEnvironmentRead(d, meta)
}

func ResourceSiteEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteEnvironment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["EnvironmentName"] = parts[1]

	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SiteEnvironment.NotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
