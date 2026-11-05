package ddosbasic

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDdosBasicDefenseThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceDdosBasicAntiddosCreate,
		Read:   resourceDdosBasicAntiddosRead,
		Update: resourceDdosBasicAntiddosUpdate,
		Delete: resourceDdosBasicAntiddosDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ecs", "slb", "eip"}, false),
				ForceNew:     true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ddos_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"defense"}, false),
				ForceNew:     true,
			},
			"bps": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"pps": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_auto": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_pps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_bps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDdosBasicAntiddosCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyDefenseThreshold"
	request := make(map[string]interface{})
	var err error
	request["DdosRegionId"] = client.RegionId
	request["InstanceId"] = d.Get("instance_id")
	request["InstanceType"] = d.Get("instance_type")
	if v, ok := d.GetOk("bps"); ok {
		request["Bps"] = v
	}
	if v, ok := d.GetOkExists("is_auto"); ok {
		request["IsAuto"] = v
	}
	if v, ok := d.GetOk("pps"); ok {
		request["Pps"] = v
	}
	if v, ok := d.GetOk("internet_ip"); ok {
		request["InternetIp"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ddos_basic_defense_threshold", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["InstanceType"], ":", d.Get("ddos_type")))
	antiddosPublicService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"defense", "blackhole"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, antiddosPublicService.DdosBasicAntiDdosStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceDdosBasicAntiddosRead(d, meta)
}
func resourceDdosBasicAntiddosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	antiddosPublicService := NewService(client)
	object, err := antiddosPublicService.DescribeDdosBasicAntiddos(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddos_basic_antiddos antiddosPublicService.DescribeDdosBasicAntiddos Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("instance_type", parts[1])
	d.Set("ddos_type", object["DdosType"])
	d.Set("bps", object["Bps"])
	d.Set("pps", object["Pps"])
	d.Set("is_auto", object["IsAuto"])
	d.Set("internet_ip", object["InternetIp"])
	d.Set("max_bps", object["MaxBps"])
	d.Set("max_pps", object["MaxPps"])
	return nil
}
func resourceDdosBasicAntiddosUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "ModifyDefenseThreshold"
	update := false
	request := map[string]interface{}{
		"DdosRegionId": client.RegionId,
		"InstanceId":   parts[0],
		"InstanceType": parts[1],
	}

	if d.HasChange("bps") {
		update = true
	}
	if v, ok := d.GetOk("bps"); ok {
		request["Bps"] = v
	}
	if d.HasChange("pps") {
		update = true
	}
	if v, ok := d.GetOk("pps"); ok {
		request["Pps"] = v
	}
	if d.HasChange("internet_ip") {
		update = true
	}
	if v, ok := d.GetOk("internet_ip"); ok {
		request["InternetIp"] = v
	}
	if d.HasChange("is_auto") {
		update = true
	}
	if v, ok := d.GetOkExists("is_auto"); ok {
		request["IsAuto"] = v
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, false)
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
	antiddosPublicService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"defense", "blackhole"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, antiddosPublicService.DdosBasicAntiDdosStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceDdosBasicAntiddosRead(d, meta)
}
func resourceDdosBasicAntiddosDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource alicloud_ddos_basic_antiddos. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
