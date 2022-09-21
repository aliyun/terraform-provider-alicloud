package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcdDesktopGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdDesktopGroupCreate,
		Read:   resourceAlicloudEcdDesktopGroupRead,
		Update: resourceAlicloudEcdDesktopGroupUpdate,
		Delete: resourceAlicloudEcdDesktopGroupDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allow_auto_setup": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allow_buffer_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bind_amount": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bundle_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_init_desktop_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"desktop_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"keep_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_desktops_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_desktops_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"own_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reset_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scale_strategy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcdDesktopGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDesktopGroup"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("allow_auto_setup"); ok {
		request["AllowAutoSetup"] = v
	}
	if v, ok := d.GetOk("allow_buffer_count"); ok {
		request["AllowBufferCount"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("bind_amount"); ok {
		request["BindAmount"] = v
	}
	request["BundleId"] = d.Get("bundle_id")
	if v, ok := d.GetOk("charge_type"); ok {
		request["ChargeType"] = v
	}
	if v, ok := d.GetOk("comments"); ok {
		request["Comments"] = v
	}
	if v, ok := d.GetOk("default_init_desktop_count"); ok {
		request["DefaultInitDesktopCount"] = v
	}
	if v, ok := d.GetOk("desktop_group_name"); ok {
		request["DesktopGroupName"] = v
	}
	if v, ok := d.GetOk("directory_id"); ok {
		request["DirectoryId"] = v
	}
	if v, ok := d.GetOk("keep_duration"); ok {
		request["KeepDuration"] = v
	}
	if v, ok := d.GetOk("load_policy"); ok {
		request["LoadPolicy"] = v
	}
	if v, ok := d.GetOk("max_desktops_count"); ok {
		request["MaxDesktopsCount"] = v
	}
	if v, ok := d.GetOk("min_desktops_count"); ok {
		request["MinDesktopsCount"] = v
	}
	request["OfficeSiteId"] = d.Get("office_site_id")
	if v, ok := d.GetOk("own_type"); ok {
		request["OwnType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("policy_group_id"); ok {
		request["PolicyGroupId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("reset_type"); ok {
		request["ResetType"] = v
	}
	if v, ok := d.GetOk("scale_strategy_id"); ok {
		request["ScaleStrategyId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["ClientToken"] = buildClientToken("CreateDesktopGroup")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_desktop_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DesktopGroupId"]))

	return resourceAlicloudEcdDesktopGroupRead(d, meta)
}
func resourceAlicloudEcdDesktopGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdDesktopGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_desktop_group ecdService.DescribeEcdDesktopGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("comments", object["Comments"])
	d.Set("status", object["Status"])
	d.Set("bundle_id", object["OwnBundleId"])
	d.Set("charge_type", object["PayType"])
	d.Set("desktop_group_name", object["DesktopGroupName"])
	d.Set("directory_id", object["DirectoryId"])
	d.Set("keep_duration", fmt.Sprint(formatInt(object["KeepDuration"])))
	if v, ok := object["MaxDesktopsCount"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_desktops_count", formatInt(v))
	}
	if v, ok := object["MinDesktopsCount"]; ok && fmt.Sprint(v) != "0" {
		d.Set("min_desktops_count", formatInt(v))
	}
	d.Set("office_site_id", object["OfficeSiteId"])
	d.Set("policy_group_id", object["PolicyGroupId"])
	d.Set("image_id", object["ImageId"])
	d.Set("load_policy", object["LoadPolicy"])

	objectDetail, err := ecdService.DescribeEcdDesktopGroupDetail(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("allow_buffer_count", objectDetail["AllowBufferCount"])
	d.Set("allow_auto_setup", objectDetail["AllowAutoSetup"])
	d.Set("bind_amount", objectDetail["BindAmount"])
	d.Set("reset_type", objectDetail["ResetType"])
	d.Set("own_type", objectDetail["OwnType"])
	return nil
}
func resourceAlicloudEcdDesktopGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DesktopGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("allow_auto_setup") {
		update = true
		if v, ok := d.GetOk("allow_auto_setup"); ok {
			request["AllowAutoSetup"] = v
		}
	}
	if d.HasChange("allow_buffer_count") {
		update = true
		if v, ok := d.GetOk("allow_buffer_count"); ok {
			request["AllowBufferCount"] = v
		}
	}
	if d.HasChange("bundle_id") {
		update = true
		request["OwnBundleId"] = d.Get("bundle_id")
	}
	if d.HasChange("comments") {
		update = true
		if v, ok := d.GetOk("comments"); ok {
			request["Comments"] = v
		}
	}
	if d.HasChange("desktop_group_name") {
		update = true
		if v, ok := d.GetOk("desktop_group_name"); ok {
			request["DesktopGroupName"] = v
		}
	}
	if d.HasChange("keep_duration") {
		update = true
		if v, ok := d.GetOk("keep_duration"); ok {
			request["KeepDuration"] = v
		}
	}
	if d.HasChange("max_desktops_count") {
		update = true
		if v, ok := d.GetOk("max_desktops_count"); ok {
			request["MaxDesktopsCount"] = v
		}
	}
	if d.HasChange("min_desktops_count") {
		update = true
		if v, ok := d.GetOk("min_desktops_count"); ok {
			request["MinDesktopsCount"] = v
		}
	}
	if d.HasChange("policy_group_id") {
		update = true
		if v, ok := d.GetOk("policy_group_id"); ok {
			request["PolicyGroupId"] = v
		}
	}
	if d.HasChange("scale_strategy_id") {
		update = true
		if v, ok := d.GetOk("scale_strategy_id"); ok {
			request["ScaleStrategyId"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("bind_amount"); ok {
			request["BindAmount"] = v
		}
		if v, ok := d.GetOk("image_id"); ok {
			request["ImageId"] = v
		}
		if v, ok := d.GetOk("load_policy"); ok {
			request["LoadPolicy"] = v
		}
		if v, ok := d.GetOk("reset_type"); ok {
			request["ResetType"] = v
		}
		action := "ModifyDesktopGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		ecdService := EcdService{client}
		stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecdService.EcdDesktopGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudEcdDesktopGroupRead(d, meta)
}
func resourceAlicloudEcdDesktopGroupDelete(d *schema.ResourceData, meta interface{}) error {
	time.Sleep(30 * time.Second)
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDesktopGroup"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DesktopGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	ecdService := EcdService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecdService.EcdDesktopGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
