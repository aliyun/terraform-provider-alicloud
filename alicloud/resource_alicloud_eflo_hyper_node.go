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

func resourceAliCloudEfloHyperNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloHyperNodeCreate,
		Read:   resourceAliCloudEfloHyperNodeRead,
		Update: resourceAliCloudEfloHyperNodeUpdate,
		Delete: resourceAliCloudEfloHyperNodeDelete,
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
			"hpn_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_arch": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stage_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEfloHyperNodeCreate(d *schema.ResourceData, meta interface{}) error {

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
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "PaymentRatio",
		"Value": "0",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Classify",
		"Value": "gpuserver",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "discountlevel",
		"Value": "0",
	})
	if v, ok := d.GetOk("machine_type"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "computingserver",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "ProductForm",
		"Value": "Hypernode",
	})
	if v, ok := d.GetOk("zone_id"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Zone",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["SubscriptionType"] = d.Get("payment_type")
	if v, ok := d.GetOkExists("renewal_duration"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("payment_duration"); ok {
		request["Period"] = v
	}
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_computinginstance_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_hyper_node", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"HealthyUnused"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloHyperNodeStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEfloHyperNodeUpdate(d, meta)
}

func resourceAliCloudEfloHyperNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloHyperNode(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_hyper_node DescribeEfloHyperNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("hpn_zone", objectRaw["HpnZone"])
	d.Set("machine_type", objectRaw["MachineType"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("zone_id", objectRaw["ZoneId"])

	objectRaw, err = efloServiceV2.DescribeHyperNodeQueryAvailableInstances(d)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("region_id", objectRaw["Region"])
	if fmt.Sprint(objectRaw["RenewalDurationUnit"]) == "Y" {
		d.Set("renewal_duration", formatInt(objectRaw["RenewalDuration"])*12)
	} else {
		d.Set("renewal_duration", objectRaw["RenewalDuration"])
	}
	d.Set("renewal_status", objectRaw["RenewStatus"])

	objectRaw, err = efloServiceV2.DescribeHyperNodeListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEfloHyperNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	if !d.IsNewResource() && d.HasChange("renewal_duration") {
		update = true
		request["RenewalPeriod"] = d.Get("renewal_duration")
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renewal_status")
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["SubscriptionType"] = v
	}
	if request["SubscriptionType"] == "" {
		request["SubscriptionType"] = "Subscription"
	}
	if request["SubscriptionType"] == "Subscription" {
		v, ok := d.GetOk("renewal_duration")
		if !ok {
			return WrapError(Error("renewal_duration is required when renewal_status is set to AutoRenewal."))
		}
		request["RenewalPeriod"] = v
		if v.(int) < 12 {
			request["RenewalPeriod"] = v
			request["RenewalPeriodUnit"] = "M"
		} else {
			if v.(int)%12 != 0 {
				return WrapError(Error("renewal_duration must be a multiple of 12 when renewal_duration more than 12."))
			}
			renewPeriod := v.(int) / 12
			request["RenewalPeriod"] = renewPeriod
			request["RenewalPeriodUnit"] = "Y"
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_eflocomputing_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "bccluster"
						request["ProductType"] = "bccluster_computinginstance_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
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
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "HyperNode"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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
		efloServiceV2 := EfloServiceV2{client}
		if err := efloServiceV2.SetResourceTags(d, "HyperNode"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEfloHyperNodeRead(d, meta)
}

func resourceAliCloudEfloHyperNodeDelete(d *schema.ResourceData, meta interface{}) error {

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
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_computinginstance_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
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
