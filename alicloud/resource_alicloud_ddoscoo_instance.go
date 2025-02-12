package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDdoscooInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdoscooInstanceCreate,
		Read:   resourceAliCloudDdoscooInstanceRead,
		Update: resourceAliCloudDdoscooInstanceUpdate,
		Delete: resourceAliCloudDdoscooInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringLenBetween(1, 64),
			},
			"port_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"normal_bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"normal_qps": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"edition_sale": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"coop"}, false),
			},
			"product_plan": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Ipv4", "Ipv6"}, false),
			},
			"bandwidth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"0", "1", "2"}, false),
			},
			"function_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"0", "1"}, false),
			},
			"product_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ddoscoo", "ddoscoo_intl", "ddosDip"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDdoscooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "CreateInstance"
	request := make(map[string]interface{})

	request["ProductCode"] = "ddos"
	request["SubscriptionType"] = "Subscription"

	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "ServicePartner",
		"Value": "coop-line-001",
	})

	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "PortCount",
		"Value": d.Get("port_count").(string),
	})

	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "DomainCount",
		"Value": d.Get("domain_count").(string),
	})

	if v, ok := d.GetOk("base_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BaseBandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Bandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("service_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ServiceBandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("normal_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "NormalBandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("normal_qps"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "NormalQps",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("edition_sale"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Edition",
			"Value": v,
		})
	} else {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Edition",
			"Value": "coop",
		})
	}

	if v, ok := d.GetOk("product_plan"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ProductPlan",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("address_type"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "AddressType",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("bandwidth_mode"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "95BurstBandwidthMode",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("function_version"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "FunctionVersion",
			"Value": v,
		})
	}

	request["Parameter"] = parameterMapList

	if v, ok := d.GetOk("product_type"); ok {
		request["ProductType"] = v
	} else {
		request["ProductType"] = "ddoscoo"
		if client.IsInternationalAccount() {
			request["ProductType"] = "ddoscoo_intl"
		}
	}

	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	} else {
		request["Period"] = 1
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", action, AlibabaCloudSdkGoERROR)
	}

	response = response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))

	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ddoscooService.DdosStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDdoscooInstanceUpdate(d, client)
}

func resourceAliCloudDdoscooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	instanceInfo, err := ddoscooService.DescribeDdoscooInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	instanceSpecInfo, err := ddoscooService.DescribeDdoscooInstanceSpec(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	instanceExtInfo, err := ddoscooService.DescribeDdoscooInstanceExt(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", instanceInfo["Remark"])
	d.Set("port_count", instanceSpecInfo["PortLimit"])
	d.Set("domain_count", instanceSpecInfo["DomainLimit"])
	d.Set("base_bandwidth", instanceSpecInfo["BaseBandwidth"])
	d.Set("bandwidth", instanceSpecInfo["ElasticBandwidth"])
	d.Set("service_bandwidth", instanceSpecInfo["BandwidthMbps"])
	d.Set("normal_bandwidth", instanceExtInfo["NormalBandwidth"])
	d.Set("normal_qps", instanceSpecInfo["QpsLimit"])
	d.Set("edition_sale", convertEditionResponse(formatInt(instanceInfo["Edition"])))
	d.Set("product_plan", instanceExtInfo["ProductPlan"])
	d.Set("address_type", instanceInfo["IpVersion"])
	d.Set("function_version", instanceExtInfo["FunctionVersion"])
	d.Set("ip", instanceInfo["Ip"])

	return nil
}

func resourceAliCloudDdoscooInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if d.HasChange("name") {
		action := "ModifyInstanceRemark"
		request := map[string]interface{}{
			"InstanceId": d.Id(),
			"Remark":     d.Get("name"),
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
			return WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliCloudDdoscooInstanceRead(d, meta)
	}

	if d.HasChange("bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("bandwidth", "Bandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("bandwidth")
	}

	if d.HasChange("base_bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("base_bandwidth", "BaseBandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("base_bandwidth")
	}

	if d.HasChange("domain_count") {
		if err := ddoscooService.UpdateInstanceSpec("domain_count", "DomainCount", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("domain_count")
	}

	if d.HasChange("port_count") {
		if err := ddoscooService.UpdateInstanceSpec("port_count", "PortCount", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("port_count")
	}

	if d.HasChange("service_bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("service_bandwidth", "ServiceBandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("service_bandwidth")
	}

	stateConf := BuildStateConf([]string{""}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ddoscooService.DdosStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.Partial(false)

	return resourceAliCloudDdoscooInstanceRead(d, meta)
}

func resourceAliCloudDdoscooInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseInstance"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"RegionId":   "cn-hangzhou",
		"InstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}
		if IsExpectedErrors(err, []string{"InstanceNotExpire"}) {
			log.Printf("[INFO]  instance cannot be deleted and must wait it to be expired and release it automatically")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertEditionResponse(source int) string {
	switch source {
	case 9:
		return "coop"
	}
	return ""
}
