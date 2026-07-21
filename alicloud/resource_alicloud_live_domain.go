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

func resourceAliCloudLiveDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudLiveDomainCreate,
		Read:   resourceAliCloudLiveDomainRead,
		Update: resourceAliCloudLiveDomainUpdate,
		Delete: resourceAliCloudLiveDomainDelete,
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
				ValidateFunc: StringInSlice([]string{"domestic", "overseas", "global"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudLiveDomainCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"LockFail"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_live_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))

	liveServiceV2 := LiveServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, liveServiceV2.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudLiveDomainUpdate(d, meta)
}

func resourceAliCloudLiveDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	liveServiceV2 := LiveServiceV2{client}

	objectRaw, err := liveServiceV2.DescribeLiveDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_live_domain DescribeLiveDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["GmtCreated"])
	d.Set("domain_type", objectRaw["LiveDomainType"])
	d.Set("region", objectRaw["Region"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("scope", objectRaw["Scope"])
	d.Set("status", objectRaw["DomainStatus"])
	d.Set("domain_name", objectRaw["DomainName"])

	objectRaw, err = liveServiceV2.DescribeDomainListLiveTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagResourceRawObj, _ := jsonpath.Get("$.TagResources.TagResource[*]", objectRaw)
	tagResourceRaw := make([]interface{}, 0)
	if tagResourceRawObj != nil {
		tagResourceRaw = convertToInterfaceArray(tagResourceRawObj)
	}

	d.Set("tags", tagsToMap(tagResourceRaw))

	return nil
}

func resourceAliCloudLiveDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	liveServiceV2 := LiveServiceV2{client}
	objectRaw, _ := liveServiceV2.DescribeLiveDomain(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("DomainStatus", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "DomainStatus", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "online" {
				action := "StartLiveDomain"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["DomainName"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
				liveServiceV2 := LiveServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, liveServiceV2.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "offline" {
				action := "StopLiveDomain"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["DomainName"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
				liveServiceV2 := LiveServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, liveServiceV2.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{"online"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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

	if d.HasChange("tags") {
		liveServiceV2 := LiveServiceV2{client}
		if err := liveServiceV2.SetLiveResourceTags(d, "DOMAIN"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudLiveDomainRead(d, meta)
}

func resourceAliCloudLiveDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteLiveDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DomainName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("live", "2016-11-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	liveServiceV2 := LiveServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 2*time.Second, liveServiceV2.LiveDomainStateRefreshFunc(d.Id(), "DomainStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
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
