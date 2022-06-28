package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDdosBasicDefenseThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdosBasicAntiddosCreate,
		Read:   resourceAlicloudDdosBasicAntiddosRead,
		Update: resourceAlicloudDdosBasicAntiddosUpdate,
		Delete: resourceAlicloudDdosBasicAntiddosDelete,
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
				ValidateFunc: validation.StringInSlice([]string{"defense", "blackhole"}, false),
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

func resourceAlicloudDdosBasicAntiddosCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyDefenseThreshold"
	request := make(map[string]interface{})
	conn, err := client.NewDdosbasicClient()
	if err != nil {
		return WrapError(err)
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddos_basic_defense_threshold", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["InstanceType"], ":", d.Get("ddos_type")))
	antiddosPublicService := AntiddosPublicService{client}
	stateConf := BuildStateConf([]string{}, []string{"defense", "blackhole"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, antiddosPublicService.DdosBasicAntiDdosStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudDdosBasicAntiddosRead(d, meta)
}
func resourceAlicloudDdosBasicAntiddosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	antiddosPublicService := AntiddosPublicService{client}
	object, err := antiddosPublicService.DescribeDdosBasicAntiddos(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddos_basic_antiddos antiddosPublicService.DescribeDdosBasicAntiddos Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
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
func resourceAlicloudDdosBasicAntiddosUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "ModifyDefenseThreshold"
	conn, err := client.NewDdosbasicClient()
	if err != nil {
		return WrapError(err)
	}
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
	antiddosPublicService := AntiddosPublicService{client}
	stateConf := BuildStateConf([]string{}, []string{"defense", "blackhole"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, antiddosPublicService.DdosBasicAntiDdosStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudDdosBasicAntiddosRead(d, meta)
}
func resourceAlicloudDdosBasicAntiddosDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource alicloud_ddos_basic_antiddos. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
