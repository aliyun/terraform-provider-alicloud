package ddosbasic

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDdosBasicThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceDdosBasicThresholdCreate,
		Read:   resourceDdosBasicThresholdRead,
		Update: resourceDdosBasicThresholdUpdate,
		Delete: resourceDdosBasicThresholdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"ecs", "slb", "eip"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bps": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pps": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_bps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_pps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDdosBasicThresholdCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyIpDefenseThreshold"
	request := make(map[string]interface{})
	var err error

	request["DdosRegionId"] = client.RegionId
	request["InstanceType"] = d.Get("instance_type")
	request["InstanceId"] = d.Get("instance_id")
	request["InternetIp"] = d.Get("internet_ip")
	request["Bps"] = d.Get("bps")
	request["Pps"] = d.Get("pps")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ddos_basic_threshold", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["InstanceType"], request["InstanceId"], request["InternetIp"]))

	return resourceDdosBasicThresholdRead(d, meta)
}

func resourceDdosBasicThresholdRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	antiddosPublicService := NewService(client)

	object, err := antiddosPublicService.DescribeDdosBasicThreshold(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddos_basic_threshold antiddosPublicService.DescribeDdosBasicThreshold Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("instance_type", parts[0])
	d.Set("instance_id", parts[1])
	d.Set("internet_ip", parts[2])

	if v, ok := object["Bps"]; ok && fmt.Sprint(v) != "0" {
		bpsVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("bps", bpsVal)
	}

	if v, ok := object["Pps"]; ok && fmt.Sprint(v) != "0" {
		ppsVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("pps", ppsVal)
	}

	if v, ok := object["MaxBps"]; ok && fmt.Sprint(v) != "0" {
		max_bpsVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("max_bps", max_bpsVal)
	}

	if v, ok := object["MaxPps"]; ok && fmt.Sprint(v) != "0" {
		max_ppsVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("max_pps", max_ppsVal)
	}

	return nil
}

func resourceDdosBasicThresholdUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DdosRegionId": client.RegionId,
		"InstanceType": parts[0],
		"InstanceId":   parts[1],
		"InternetIp":   parts[2],
	}

	if d.HasChange("bps") {
		update = true
	}
	request["Bps"] = d.Get("bps")

	if d.HasChange("pps") {
		update = true
	}
	request["Pps"] = d.Get("pps")

	if update {
		action := "ModifyIpDefenseThreshold"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, true)
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
	}

	return resourceDdosBasicThresholdRead(d, meta)
}

func resourceDdosBasicThresholdDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceDdosBasicThreshold. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
