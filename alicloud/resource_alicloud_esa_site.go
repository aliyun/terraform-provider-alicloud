// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaSiteCreate,
		Read:   resourceAliCloudEsaSiteRead,
		Update: resourceAliCloudEsaSiteUpdate,
		Delete: resourceAliCloudEsaSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"add_client_geolocation_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"add_real_client_ip_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_architecture_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"coverage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"site_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudEsaSiteCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSite"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["SiteName"] = d.Get("site_name")
	request["Coverage"] = d.Get("coverage")
	request["AccessType"] = d.Get("access_type")
	request["InstanceId"] = d.Get("instance_id")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_site", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SiteId"]))

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"pending"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, esaServiceV2.EsaSiteStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEsaSiteUpdate(d, meta)
}

func resourceAliCloudEsaSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaSite(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_site DescribeEsaSite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_type", objectRaw["AccessType"])
	d.Set("coverage", objectRaw["Coverage"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("site_name", objectRaw["SiteName"])
	d.Set("status", objectRaw["Status"])

	objectRaw, err = esaServiceV2.DescribeSiteListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = esaServiceV2.DescribeSiteGetManagedTransform(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("add_client_geolocation_header", objectRaw["AddClientGeolocationHeader"])
	d.Set("add_real_client_ip_header", objectRaw["AddRealClientIpHeader"])
	d.Set("site_version", objectRaw["SiteVersion"])

	objectRaw, err = esaServiceV2.DescribeSiteGetIPv6(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("ipv6_enable", objectRaw["Enable"])

	objectRaw, err = esaServiceV2.DescribeSiteGetTieredCache(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("cache_architecture_mode", objectRaw["CacheArchitectureMode"])

	return nil
}

func resourceAliCloudEsaSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateSiteCoverage"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("coverage") {
		update = true
	}
	request["Coverage"] = d.Get("coverage")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy"}) || NeedRetry(err) {
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
		esaServiceV2 := EsaServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"pending"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, esaServiceV2.EsaSiteStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "UpdateIPv6"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("ipv6_enable") {
		update = true
	}
	request["Enable"] = d.Get("ipv6_enable")
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
	update = false
	action = "UpdateTieredCache"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("cache_architecture_mode") {
		update = true
	}
	request["CacheArchitectureMode"] = d.Get("cache_architecture_mode")
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
	update = false
	action = "UpdateManagedTransform"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("add_client_geolocation_header") {
		update = true
		request["AddClientGeolocationHeader"] = d.Get("add_client_geolocation_header")
	}

	if d.HasChange("add_real_client_ip_header") {
		update = true
		request["AddRealClientIpHeader"] = d.Get("add_real_client_ip_header")
	}

	if d.HasChange("site_version") {
		update = true
		request["SiteVersion"] = d.Get("site_version")
	}

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

	if d.HasChange("tags") {
		esaServiceV2 := EsaServiceV2{client}
		if err := esaServiceV2.SetResourceTags(d, "Site"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEsaSiteRead(d, meta)
}

func resourceAliCloudEsaSiteDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSite"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy"}) || NeedRetry(err) {
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

	esaServiceV2 := EsaServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, esaServiceV2.DescribeAsyncEsaSiteStateRefreshFunc(d, response, "$.SiteModel.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
