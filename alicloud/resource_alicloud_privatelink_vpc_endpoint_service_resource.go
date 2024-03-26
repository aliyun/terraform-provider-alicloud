// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudPrivateLinkVpcEndpointServiceResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPrivateLinkVpcEndpointServiceResourceCreate,
		Read:   resourceAliCloudPrivateLinkVpcEndpointServiceResourceRead,
		Update: resourceAliCloudPrivateLinkVpcEndpointServiceResourceUpdate,
		Delete: resourceAliCloudPrivateLinkVpcEndpointServiceResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"slb", "alb", "nlb"}, false),
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPrivateLinkVpcEndpointServiceResourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AttachResourceToVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["ResourceId"] = d.Get("resource_id")
	query["ServiceId"] = d.Get("service_id")
	query["ZoneId"] = d.Get("zone_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["ResourceType"] = d.Get("resource_type")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_resource", action, AlibabaCloudSdkGoERROR)
	}

	if d.Get("resource_type") == "slb" {
		d.SetId(fmt.Sprintf("%v:%v", query["ServiceId"], query["ResourceId"]))
	} else {
		d.SetId(fmt.Sprintf("%v:%v:%v", query["ServiceId"], query["ResourceId"], query["ZoneId"]))
	}

	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_id"))}, d.Timeout(schema.TimeoutCreate), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceResourceStateRefreshFunc(d.Id(), "ResourceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPrivateLinkVpcEndpointServiceResourceRead(d, meta)
}

func resourceAliCloudPrivateLinkVpcEndpointServiceResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privateLinkServiceV2 := PrivateLinkServiceV2{client}

	objectRaw, err := privateLinkServiceV2.DescribePrivateLinkVpcEndpointServiceResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service_resource DescribePrivateLinkVpcEndpointServiceResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("resource_type", objectRaw["ResourceType"])
	d.Set("resource_id", objectRaw["ResourceId"])
	d.Set("zone_id", objectRaw["ZoneId"])
	parts, err := ParseResourceIds(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("resource_id", parts[1])
	d.Set("service_id", parts[0])
	if len(parts) == 3 {
		d.Set("zone_id", parts[2])
	}

	return nil
}

func resourceAliCloudPrivateLinkVpcEndpointServiceResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Vpc Endpoint Service Resource.")
	return nil
}

func resourceAliCloudPrivateLinkVpcEndpointServiceResourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DetachResourceFromVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["ResourceId"] = parts[1]
	query["ServiceId"] = parts[0]
	if len(parts) == 3 {
		query["ZoneId"] = parts[2]
	}
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	request["ResourceType"] = d.Get("resource_type")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceResourceStateRefreshFunc(d.Id(), "ResourceId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
