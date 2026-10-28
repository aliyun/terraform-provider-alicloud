package cddc

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedHostCreate,
		Read:   resourceDedicatedHostRead,
		Update: resourceDedicatedHostUpdate,
		Delete: resourceDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		DeprecationMessage: "This resource has been deprecated since v1.225.1 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"allocation_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Allocatable", "Suspended"}, false),
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dedicated_host_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{1,63}`), "The name must be `1` to `64` characters in length and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter."),
			},
			"image_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"WindowsWithMssqlEntAlwaysonLicense", "WindowsWithMssqlStdLicense", "WindowsWithMssqlEntLicense", "WindowsWithMssqlWebLicense", "AliLinux"}, false),
			},
			"os_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year", "Week"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"used_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDedicatedHost"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["DedicatedHostGroupId"] = d.Get("dedicated_host_group_id")
	request["HostClass"] = d.Get("host_class")
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	if v, ok := d.GetOk("image_category"); ok {
		request["ImageCategory"] = v
	}
	if v, ok := d.GetOk("os_password"); ok {
		request["OsPassword"] = v
	}
	request["PayType"] = convertDedicatedPaymentTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = fmt.Sprint(v)
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("CreateDedicatedHost")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cddc_dedicated_host", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseDedicateHostList := response["DedicateHostList"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["DedicatedHostGroupId"], ":", responseDedicateHostList["DedicateHostList"].([]interface{})[0].(map[string]interface{})["DedicatedHostId"]))
	cddcService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cddcService.DedicatedHostStateRefreshFunc(d.Id(), []string{"5"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDedicatedHostUpdate(d, meta)
}
func resourceDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := NewService(client)
	object, err := cddcService.DescribeDedicatedHost(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_host cddcService.DescribeDedicatedHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("dedicated_host_group_id", object["DedicatedHostGroupId"])
	d.Set("host_class", object["HostClass"])
	d.Set("host_name", object["HostName"])
	d.Set("status", object["HostStatus"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("allocation_status", convertAllocationStatusResponse(object["AllocationStatus"].(string)))
	d.Set("dedicated_host_id", object["DedicatedHostId"])
	listTagResourcesObject, err := cddcService.ListTagResources(d.Id(), "DEDICATEDHOST")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(listTagResourcesObject))
	return nil
}
func resourceDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := NewService(client)
	var err error
	var response map[string]interface{}
	d.Partial(true)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	if d.HasChange("tags") {
		if err := cddcService.SetResourceTags(d, "DEDICATEDHOST"); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.HasChange("host_name") || d.HasChange("allocation_status") {
		request := map[string]interface{}{
			"DedicatedHostId": parts[1],
		}
		if v, ok := d.GetOk("host_name"); ok {
			request["HostName"] = v
		}
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOk("allocation_status"); ok {
			request["AllocationStatus"] = convertAllocationStatusRequest(v.(string))
		}
		action := "ModifyDedicatedHostAttribute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cddcService.DedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("host_class") {
		request := map[string]interface{}{
			"DedicatedHostId": parts[1],
		}
		request["TargetClassCode"] = d.Get("host_class")
		request["RegionId"] = client.RegionId
		action := "ModifyDedicatedHostClass"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cddcService.DedicatedHostStateRefreshFunc(d.Id(), []string{"5"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceDedicatedHostRead(d, meta)
}
func resourceDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AlicloudResourceCddcDedicatedHost. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertDedicatedPaymentTypeRequest(source string) string {
	switch source {
	case "Subscription":
		return "Prepaid"
	}
	return source
}
func convertAllocationStatusRequest(source string) string {
	switch source {
	case "Allocatable":
		return "1"
	case "Suspended":
		return "0"
	}
	return source
}
func convertAllocationStatusResponse(source string) string {
	switch source {
	case "1":
		return "Allocatable"
	case "0":
		return "Suspended"
	}
	return source
}
