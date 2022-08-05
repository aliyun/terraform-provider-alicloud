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

func resourceAlicloudDdosBasicThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdosBasicThresholdCreate,
		Read:   resourceAlicloudDdosBasicThresholdRead,
		Update: resourceAlicloudDdosBasicThresholdUpdate,
		Delete: resourceAlicloudDdosBasicThresholdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bps": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ecs", "slb", "eip"}, false),
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_bps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_pps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pps": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceAlicloudDdosBasicThresholdCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyIpDefenseThreshold"
	request := make(map[string]interface{})
	conn, err := client.NewDdosbasicClient()
	if err != nil {
		return WrapError(err)
	}
	request["Bps"] = d.Get("bps")
	request["Pps"] = d.Get("pps")
	request["DdosRegionId"] = client.RegionId
	request["InstanceId"] = d.Get("instance_id")
	request["InstanceType"] = d.Get("instance_type")
	request["InternetIp"] = d.Get("internet_ip")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddos_basic_threshold", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceType"], ":", request["InstanceId"], ":", request["InternetIp"]))

	return resourceAlicloudDdosBasicThresholdRead(d, meta)
}
func resourceAlicloudDdosBasicThresholdRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	antiddosPublicService := AntiddosPublicService{client}
	object, err := antiddosPublicService.DescribeDdosBasicThreshold(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddos_basic_threshold antiddosPublicService.DescribeDdosBasicThreshold Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", parts[0])
	d.Set("instance_id", parts[1])
	d.Set("internet_ip", parts[2])
	if v, ok := object["Bps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("bps", formatInt(v))
	}
	if v, ok := object["MaxBps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_bps", formatInt(v))
	}
	if v, ok := object["MaxPps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_pps", formatInt(v))
	}
	if v, ok := object["Pps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("pps", formatInt(v))
	}
	return nil
}
func resourceAlicloudDdosBasicThresholdUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewDdosbasicClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{}
	request["DdosRegionId"] = client.RegionId
	request["InstanceId"] = parts[1]
	request["InstanceType"] = parts[0]
	request["InternetIp"] = parts[2]
	request["Bps"] = d.Get("bps")
	if d.HasChange("bps") {
		update = true
	}
	request["Pps"] = d.Get("pps")
	if d.HasChange("pps") {
		update = true
	}
	if update {
		action := "ModifyIpDefenseThreshold"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudDdosBasicThresholdRead(d, meta)
}
func resourceAlicloudDdosBasicThresholdDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudDdosBasicThreshold. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
