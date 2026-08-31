package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDdosBasicThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosBasicThresholdCreate,
		Read:   resourceAliCloudDdosBasicThresholdRead,
		Update: resourceAliCloudDdosBasicThresholdUpdate,
		Delete: resourceAliCloudDdosBasicThresholdDelete,
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
				ValidateFunc: StringInSlice([]string{"ecs", "slb", "eip"}, false),
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

func resourceAliCloudDdosBasicThresholdCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, true)
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

	d.SetId(fmt.Sprintf("%v:%v:%v", request["InstanceType"], request["InstanceId"], request["InternetIp"]))

	return resourceAliCloudDdosBasicThresholdRead(d, meta)
}

func resourceAliCloudDdosBasicThresholdRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	antiddosPublicService := AntiddosPublicService{client}

	object, err := antiddosPublicService.DescribeDdosBasicThreshold(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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

	if v, ok := object["Pps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("pps", formatInt(v))
	}

	if v, ok := object["MaxBps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_bps", formatInt(v))
	}

	if v, ok := object["MaxPps"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_pps", formatInt(v))
	}

	return nil
}

func resourceAliCloudDdosBasicThresholdUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, true)
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

	return resourceAliCloudDdosBasicThresholdRead(d, meta)
}

func resourceAliCloudDdosBasicThresholdDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAliCloudDdosBasicThreshold. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
