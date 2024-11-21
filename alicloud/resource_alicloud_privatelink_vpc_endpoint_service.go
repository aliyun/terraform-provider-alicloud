package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPrivateLinkVpcEndpointService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPrivateLinkVpcEndpointServiceCreate,
		Read:   resourceAliCloudPrivateLinkVpcEndpointServiceRead,
		Update: resourceAliCloudPrivateLinkVpcEndpointServiceUpdate,
		Delete: resourceAliCloudPrivateLinkVpcEndpointServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_accept_connection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"connect_bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(100, 10240),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payer": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Endpoint", "EndpointService"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"slb", "alb", "nlb"}, false),
			},
			"service_support_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPrivateLinkVpcEndpointServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOkExists("service_support_ipv6"); ok {
		request["ServiceSupportIPv6"] = v
	}
	if v, ok := d.GetOk("service_resource_type"); ok {
		request["ServiceResourceType"] = v
	}
	if v, ok := d.GetOkExists("zone_affinity_enabled"); ok {
		request["ZoneAffinityEnabled"] = v
	}
	if v, ok := d.GetOk("service_description"); ok {
		request["ServiceDescription"] = v
	}
	if v, ok := d.GetOkExists("auto_accept_connection"); ok {
		request["AutoAcceptEnabled"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("payer"); ok {
		request["Payer"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceId"]))

	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(d, meta)
}

func resourceAliCloudPrivateLinkVpcEndpointServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privateLinkServiceV2 := PrivateLinkServiceV2{client}

	objectRaw, err := privateLinkServiceV2.DescribePrivateLinkVpcEndpointService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service DescribePrivateLinkVpcEndpointService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AddressIpVersion"] != nil {
		d.Set("address_ip_version", objectRaw["AddressIpVersion"])
	}
	if objectRaw["AutoAcceptEnabled"] != nil {
		d.Set("auto_accept_connection", objectRaw["AutoAcceptEnabled"])
	}
	if objectRaw["ConnectBandwidth"] != nil {
		d.Set("connect_bandwidth", objectRaw["ConnectBandwidth"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Payer"] != nil {
		d.Set("payer", objectRaw["Payer"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["ServiceBusinessStatus"] != nil {
		d.Set("service_business_status", objectRaw["ServiceBusinessStatus"])
	}
	if objectRaw["ServiceDescription"] != nil {
		d.Set("service_description", objectRaw["ServiceDescription"])
	}
	if objectRaw["ServiceDomain"] != nil {
		d.Set("service_domain", objectRaw["ServiceDomain"])
	}
	if objectRaw["ServiceResourceType"] != nil {
		d.Set("service_resource_type", objectRaw["ServiceResourceType"])
	}
	if objectRaw["ServiceSupportIPv6"] != nil {
		d.Set("service_support_ipv6", objectRaw["ServiceSupportIPv6"])
	}
	if objectRaw["ServiceStatus"] != nil {
		d.Set("status", objectRaw["ServiceStatus"])
	}
	if objectRaw["ServiceName"] != nil {
		d.Set("vpc_endpoint_service_name", objectRaw["ServiceName"])
	}
	if objectRaw["ZoneAffinityEnabled"] != nil {
		d.Set("zone_affinity_enabled", objectRaw["ZoneAffinityEnabled"])
	}

	objectRaw, err = privateLinkServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "UpdateVpcEndpointServiceAttribute"
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("service_description") {
		update = true
		request["ServiceDescription"] = d.Get("service_description")
	}

	if !d.IsNewResource() && d.HasChange("auto_accept_connection") {
		update = true
		request["AutoAcceptEnabled"] = d.Get("auto_accept_connection")
	}

	if !d.IsNewResource() && d.HasChange("service_support_ipv6") {
		update = true
		request["ServiceSupportIPv6"] = d.Get("service_support_ipv6")
	}

	if d.HasChange("connect_bandwidth") {
		update = true
		request["ConnectBandwidth"] = d.Get("connect_bandwidth")
	}

	if !d.IsNewResource() && d.HasChange("zone_affinity_enabled") {
		update = true
		request["ZoneAffinityEnabled"] = d.Get("zone_affinity_enabled")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if d.HasChange("address_ip_version") {
		update = true
		request["AddressIpVersion"] = d.Get("address_ip_version")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
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
	action = "ChangeResourceGroup"
	conn, err = client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if !d.IsNewResource() && d.HasChange("region_id") {
		update = true
		request["ResourceRegionId"] = d.Get("region_id")
	}

	request["ResourceType"] = "VpcEndpointService"
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-15"), StringPointer("AK"), query, request, &runtime)
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
		privateLinkServiceV2 := PrivateLinkServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		privateLinkServiceV2 := PrivateLinkServiceV2{client}
		if err := privateLinkServiceV2.SetResourceTags(d, "VpcEndpointService"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudPrivateLinkVpcEndpointServiceRead(d, meta)
}

func resourceAliCloudPrivateLinkVpcEndpointServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewPrivatelinkClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServiceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

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
			if IsExpectedErrors(err, []string{"EndpointServiceConnectionDependence", "ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
