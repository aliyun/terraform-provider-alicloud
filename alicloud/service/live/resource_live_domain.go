package live

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLiveDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiveDomainCreate,
		Read:   resourceLiveDomainRead,
		Update: resourceLiveDomainUpdate,
		Delete: resourceLiveDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"domestic", "overseas", "global"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceLiveDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddLiveDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("domain_name"); ok {
		request["DomainName"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("check_url"); ok {
		request["CheckUrl"] = v
	}
	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	request["Region"] = d.Get("region")
	request["LiveDomainType"] = d.Get("domain_type")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"LockFail"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_live_domain", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))

	liveService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, liveService.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceLiveDomainUpdate(d, meta)
}

func resourceLiveDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	liveService := NewService(client)

	objectRaw, err := liveService.DescribeLiveDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_live_domain DescribeLiveDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["GmtCreated"])
	d.Set("domain_type", objectRaw["LiveDomainType"])
	d.Set("region", objectRaw["Region"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scope", objectRaw["Scope"])
	d.Set("status", objectRaw["DomainStatus"])
	d.Set("domain_name", objectRaw["DomainName"])

	objectRaw, err = liveService.DescribeDomainListLiveTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagResourceRawObj, _ := jsonpath.Get("$.TagResources.TagResource[*]", objectRaw)
	tagResourceRaw := make([]interface{}, 0)
	if tagResourceRawObj != nil {
		tagResourceRaw = convert.ConvertToInterfaceArray(tagResourceRawObj)
	}

	d.Set("tags", tags.TagsToMap(tagResourceRaw))

	return nil
}

func resourceLiveDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	liveService := NewService(client)
	objectRaw, _ := liveService.DescribeLiveDomain(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("DomainStatus", objectRaw)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "DomainStatus", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "online" {
				action := "StartLiveDomain"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["DomainName"] = d.Id()

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
				liveService := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, liveService.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "offline" {
				action := "StopLiveDomain"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["DomainName"] = d.Id()

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
				liveService := NewService(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, liveService.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"online"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ChangeLiveDomainResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DomainName"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
	update = false
	action = "ModifyLiveDomainSchdmByProperty"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DomainName"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("scope") {
		update = true
	}
	request["Property"] = convertLiveDomainPropertyRequest(d.Get("scope").(string))
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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

	if d.HasChange("tags") {
		liveService := NewService(client)
		if err := liveService.SetLiveResourceTags(d, "DOMAIN"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceLiveDomainRead(d, meta)
}

func resourceLiveDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLiveDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DomainName"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	liveService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 2*time.Second, liveService.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertLiveDomainPropertyRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "global":
		return "{\"coverage\":\"global\"}"
	case "domestic":
		return "{\"coverage\":\"domestic\"}"
	case "overseas":
		return "{\"coverage\":\"overseas\"}"
	}

	return source
}
