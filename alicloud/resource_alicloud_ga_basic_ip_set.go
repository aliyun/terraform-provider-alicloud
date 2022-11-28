package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceAlicloudGaBasicIpSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaBasicIpSetCreate,
		Read:   resourceAlicloudGaBasicIpSetRead,
		Update: resourceAlicloudGaBasicIpSetUpdate,
		Delete: resourceAlicloudGaBasicIpSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accelerate_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"isp_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BGP", "BGP_PRO", "ChinaTelecom", "ChinaUnicom", "ChinaMobile", "ChinaTelecom_L2", "ChinaUnicom_L2", "ChinaMobile_L2"}, false),
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaBasicIpSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateBasicIpSet"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateBasicIpSet")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["AccelerateRegionId"] = d.Get("accelerate_region_id")

	if v, ok := d.GetOk("isp_type"); ok {
		request["IspType"] = v
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotExist.BasicBandwidthPackage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_basic_ip_set", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpSetId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaBasicIpSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaBasicIpSetRead(d, meta)
}

func resourceAlicloudGaBasicIpSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaBasicIpSet(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("accelerate_region_id", object["AccelerateRegionId"])
	d.Set("isp_type", object["IspType"])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("status", object["State"])

	return nil
}

func resourceAlicloudGaBasicIpSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("UpdateBasicIpSet"),
		"IpSetId":     d.Id(),
	}

	if d.HasChange("bandwidth") {
		update = true
	}
	request["Bandwidth"] = d.Get("bandwidth")

	if update {
		action := "UpdateBasicIpSet"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaBasicIpSetStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudGaBasicIpSetRead(d, meta)
}

func resourceAlicloudGaBasicIpSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteBasicIpSet"
	var response map[string]interface{}

	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": buildClientToken("DeleteBasicIpSet"),
		"IpSetId":     d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "NotActive.IpSet", "NotExist.BasicBandwidthPackage", "ExistBoundAccelerateIp.IpSet"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBasicIpSetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
