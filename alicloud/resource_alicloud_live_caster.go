// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudLiveCaster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudLiveCasterCreate,
		Read:   resourceAliCloudLiveCasterRead,
		Update: resourceAliCloudLiveCasterUpdate,
		Delete: resourceAliCloudLiveCasterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_switch_urgent_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_switch_urgent_on": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"callback_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"caster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delay": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"norm_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"program_effect": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"program_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"record_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"side_output_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"side_output_url_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync_groups_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"transcode_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urgent_image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urgent_image_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urgent_live_stream_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urgent_material_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudLiveCasterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCaster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewLiveClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("caster_name"); ok {
		request["CasterName"] = v
	}
	request["NormType"] = d.Get("norm_type")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ChargeType"] = convertLiveCasterChargeTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-11-01"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_live_caster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CasterId"]))

	return resourceAliCloudLiveCasterRead(d, meta)
}

func resourceAliCloudLiveCasterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	liveServiceV2 := LiveServiceV2{client}

	objectRaw, err := liveServiceV2.DescribeLiveCaster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_live_caster DescribeLiveCaster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CasterName"] != nil {
		d.Set("caster_name", objectRaw["CasterName"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["NormType"] != nil {
		d.Set("norm_type", objectRaw["NormType"])
	}
	if objectRaw["ChargeType"] != nil {
		d.Set("payment_type", convertLiveCasterCasterListCasterChargeTypeResponse(objectRaw["ChargeType"]))
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudLiveCasterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "SetCasterConfig"
	conn, err := client.NewLiveClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CasterId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("caster_name") {
		update = true
		request["CasterName"] = d.Get("caster_name")
	}

	if v, ok := d.GetOk("domain_name"); ok {
		request["DomainName"] = v
	}
	if v, ok := d.GetOk("auto_switch_urgent_config"); ok {
		request["AutoSwitchUrgentConfig"] = v
	}
	if v, ok := d.GetOkExists("auto_switch_urgent_on"); ok {
		request["AutoSwitchUrgentOn"] = v
	}
	if v, ok := d.GetOk("transcode_config"); ok {
		request["TranscodeConfig"] = v
	}
	if v, ok := d.GetOk("record_config"); ok {
		request["RecordConfig"] = v
	}
	if v, ok := d.GetOk("delay"); ok {
		request["Delay"] = v
	}
	if v, ok := d.GetOk("urgent_image_id"); ok {
		request["UrgentImageId"] = v
	}
	if v, ok := d.GetOk("urgent_image_url"); ok {
		request["UrgentImageUrl"] = v
	}
	if v, ok := d.GetOk("urgent_material_id"); ok {
		request["UrgentMaterialId"] = v
	}
	if v, ok := d.GetOk("sync_groups_config"); ok {
		request["SyncGroupsConfig"] = v
	}
	if v, ok := d.GetOk("program_name"); ok {
		request["ProgramName"] = v
	}
	if v, ok := d.GetOk("side_output_url"); ok {
		request["SideOutputUrl"] = v
	}
	if v, ok := d.GetOk("side_output_url_list"); ok {
		request["SideOutputUrlList"] = v
	}
	if v, ok := d.GetOk("callback_url"); ok {
		request["CallbackUrl"] = v
	}
	if v, ok := d.GetOkExists("program_effect"); ok {
		request["ProgramEffect"] = v
	}
	if v, ok := d.GetOk("urgent_live_stream_url"); ok {
		request["UrgentLiveStreamUrl"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-11-01"), StringPointer("AK"), query, request, &runtime)
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
	update = false
	action = "UpdateCasterResourceGroup"
	conn, err = client.NewLiveClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CasterId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-11-01"), StringPointer("AK"), query, request, &runtime)
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

	if d.HasChange("tags") {
		liveServiceV2 := LiveServiceV2{client}
		if err := liveServiceV2.SetResourceTags(d, ""); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudLiveCasterRead(d, meta)
}

func resourceAliCloudLiveCasterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCaster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewLiveClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["CasterId"] = d.Id()
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-11-01"), StringPointer("AK"), query, request, &runtime)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertLiveCasterCasterListCasterChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PrePaid":
		return "Subscription"
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
func convertLiveCasterChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PrePaid"
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}
