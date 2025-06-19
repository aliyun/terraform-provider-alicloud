// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudEfloNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloNodeCreate,
		Read:   resourceAliCloudEfloNodeRead,
		Update: resourceAliCloudEfloNodeUpdate,
		Delete: resourceAliCloudEfloNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"billing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"classify": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"computing_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discount_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hpn_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_ratio": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"product_form": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_arch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stage_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEfloNodeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("server_arch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ServerArch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("hpn_zone"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "HpnZone",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("stage_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "StageNum",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("payment_ratio"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "PaymentRatio",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	if v, ok := d.GetOk("classify"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Classify",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("discount_level"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "discountlevel",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("billing_cycle"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BillingCycle",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("computing_server"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "computingserver",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("zone"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Zone",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("product_form"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ProductForm",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["SubscriptionType"] = "Subscription"
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOkExists("renew_period"); ok {
		request["RenewPeriod"] = v
	}
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if IsExpectedErrors(err, []string{"CSS_CHECK_ORDER_ERROR", "InternalError", "SYSTEM.CONCURRENT_OPERATE"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_node", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Unused"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEfloNodeUpdate(d, meta)
}

func resourceAliCloudEfloNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloNode(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_node DescribeEfloNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["OperatingState"])

	objectRaw, err = efloServiceV2.DescribeNodeListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEfloNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "Node"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NeedRetry(err) {
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
		efloServiceV2 := EfloServiceV2{client}
		if err := efloServiceV2.SetResourceTags(d, "Node"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudEfloNodeRead(d, meta)
}

func resourceAliCloudEfloNodeDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ImmediatelyRelease"] = "1"
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if client.IsInternationalAccount() {
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"RESOURCE_NOT_FOUND"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, efloServiceV2.EfloNodeStateRefreshFunc(d.Id(), "$.NodeId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
