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

func resourceAlicloudGaBasicEndpointGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaBasicEndpointGroupCreate,
		Read:   resourceAlicloudGaBasicEndpointGroupRead,
		Update: resourceAlicloudGaBasicEndpointGroupUpdate,
		Delete: resourceAlicloudGaBasicEndpointGroupDelete,
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
			"endpoint_group_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENI", "SLB", "ECS"}, false),
			},
			"endpoint_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_sub_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"basic_endpoint_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaBasicEndpointGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateBasicEndpointGroup"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateBasicEndpointGroup")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["EndpointGroupRegion"] = d.Get("endpoint_group_region")

	if v, ok := d.GetOk("endpoint_type"); ok {
		request["EndpointType"] = v
	}

	if v, ok := d.GetOk("endpoint_address"); ok {
		request["EndpointAddress"] = v
	}

	if v, ok := d.GetOk("endpoint_sub_address"); ok {
		request["EndpointSubAddress"] = v
	}

	if v, ok := d.GetOk("basic_endpoint_group_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_basic_endpoint_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointGroupId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaBasicEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaBasicEndpointGroupRead(d, meta)
}

func resourceAlicloudGaBasicEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaBasicEndpointGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("endpoint_group_region", object["EndpointGroupRegion"])
	d.Set("endpoint_type", object["EndpointType"])
	d.Set("endpoint_address", object["EndpointAddress"])
	d.Set("endpoint_sub_address", object["EndpointSubAddress"])
	d.Set("basic_endpoint_group_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	return nil
}

func resourceAlicloudGaBasicEndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("UpdateBasicEndpointGroup"),
		"EndpointGroupId": d.Id(),
	}

	if d.HasChange("endpoint_sub_address") {
		update = true
	}
	if v, ok := d.GetOk("endpoint_sub_address"); ok {
		request["EndpointSubAddress"] = v
	}

	if d.HasChange("basic_endpoint_group_name") {
		update = true
	}
	if v, ok := d.GetOk("basic_endpoint_group_name"); ok {
		request["Name"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "UpdateBasicEndpointGroup"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "NotExist.BasicBandwidthPackage"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaBasicEndpointGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudGaBasicEndpointGroupRead(d, meta)
}

func resourceAlicloudGaBasicEndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteBasicEndpointGroup"
	var response map[string]interface{}

	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":     buildClientToken("DeleteBasicEndpointGroup"),
		"EndpointGroupId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup", "ExistBoundEndpoint.EndpointGroup", "NotExist.BasicBandwidthPackage"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBasicEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
