// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAligreenOssStockTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAligreenOssStockTaskCreate,
		Read:   resourceAliCloudAligreenOssStockTaskRead,
		Delete: resourceAliCloudAligreenOssStockTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"audio_antispam_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"audio_auto_freeze_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"audio_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"audio_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"audio_scan_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"audio_scenes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"auto_freeze_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"biz_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"buckets": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"callback_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_ad_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"image_auto_freeze_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"image_live_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"image_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"image_porn_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"image_scan_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"image_scenes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"image_terrorism_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"scan_image_no_file_type": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"video_ad_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"video_auto_freeze_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"video_frame_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"video_live_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"video_max_frames": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"video_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"video_opened": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"video_porn_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"video_scan_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"video_scenes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"video_terrorism_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"video_voice_antispam_freeze_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceAliCloudAligreenOssStockTaskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOssStockTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("start_date"); ok {
		request["StartDate"] = v
	}
	if v, ok := d.GetOk("end_date"); ok {
		request["EndDate"] = v
	}
	if v, ok := d.GetOk("image_opened"); ok {
		request["ImageOpened"] = v
	}
	if v, ok := d.GetOk("scan_image_no_file_type"); ok {
		request["ScanImageNoFileType"] = v
	}
	if v, ok := d.GetOk("video_opened"); ok {
		request["VideoOpened"] = v
	}
	if v, ok := d.GetOk("video_scenes"); ok {
		request["VideoScenes"] = v
	}
	if v, ok := d.GetOk("audio_opened"); ok {
		request["AudioOpened"] = v
	}
	if v, ok := d.GetOk("audio_scenes"); ok {
		request["AudioScenes"] = v
	}
	if v, ok := d.GetOk("image_auto_freeze_opened"); ok {
		request["ImageAutoFreezeOpened"] = v
	}
	if v, ok := d.GetOk("image_porn_freeze_config"); ok {
		request["ImagePornFreezeConfig"] = v
	}
	if v, ok := d.GetOk("image_terrorism_freeze_config"); ok {
		request["ImageTerrorismFreezeConfig"] = v
	}
	if v, ok := d.GetOk("image_ad_freeze_config"); ok {
		request["ImageAdFreezeConfig"] = v
	}
	if v, ok := d.GetOk("image_live_freeze_config"); ok {
		request["ImageLiveFreezeConfig"] = v
	}
	if v, ok := d.GetOk("video_frame_interval"); ok {
		request["VideoFrameInterval"] = v
	}
	if v, ok := d.GetOk("video_max_frames"); ok {
		request["VideoMaxFrames"] = v
	}
	if v, ok := d.GetOk("video_max_size"); ok {
		request["VideoMaxSize"] = v
	}
	if v, ok := d.GetOk("video_auto_freeze_opened"); ok {
		request["VideoAutoFreezeOpened"] = v
	}
	if v, ok := d.GetOk("video_porn_freeze_config"); ok {
		request["VideoPornFreezeConfig"] = v
	}
	if v, ok := d.GetOk("video_terrorism_freeze_config"); ok {
		request["VideoTerrorismFreezeConfig"] = v
	}
	if v, ok := d.GetOk("video_ad_freeze_config"); ok {
		request["VideoAdFreezeConfig"] = v
	}
	if v, ok := d.GetOk("video_live_freeze_config"); ok {
		request["VideoLiveFreezeConfig"] = v
	}
	if v, ok := d.GetOk("video_voice_antispam_freeze_config"); ok {
		request["VideoVoiceAntispamFreezeConfig"] = v
	}
	if v, ok := d.GetOk("audio_auto_freeze_opened"); ok {
		request["AudioAutoFreezeOpened"] = v
	}
	if v, ok := d.GetOk("audio_max_size"); ok {
		request["AudioMaxSize"] = v
	}
	if v, ok := d.GetOk("audio_antispam_freeze_config"); ok {
		request["AudioAntispamFreezeConfig"] = v
	}
	if v, ok := d.GetOk("auto_freeze_type"); ok {
		request["AutoFreezeType"] = v
	}
	if v, ok := d.GetOk("callback_id"); ok {
		request["CallbackId"] = v
	}
	if v, ok := d.GetOk("biz_type"); ok {
		request["BizType"] = v
	}
	if v, ok := d.GetOk("image_scan_limit"); ok {
		request["ImageScanLimit"] = v
	}
	if v, ok := d.GetOk("video_scan_limit"); ok {
		request["VideoScanLimit"] = v
	}
	if v, ok := d.GetOk("audio_scan_limit"); ok {
		request["AudioScanLimit"] = v
	}
	if v, ok := d.GetOk("image_scenes"); ok {
		jsonPathResult31, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult31 != "" {
			request["ImageScenes"] = convertListToJsonString(jsonPathResult31.([]interface{}))
		}
	}
	if v, ok := d.GetOk("buckets"); ok {
		request["Buckets"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_aligreen_oss_stock_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudAligreenOssStockTaskRead(d, meta)
}

func resourceAliCloudAligreenOssStockTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aligreenServiceV2 := AligreenServiceV2{client}

	objectRaw, err := aligreenServiceV2.DescribeAligreenOssStockTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_aligreen_oss_stock_task DescribeAligreenOssStockTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AudioAutoFreezeOpened"] != nil {
		d.Set("audio_auto_freeze_opened", objectRaw["AudioAutoFreezeOpened"])
	}
	if objectRaw["AudioMaxSize"] != nil {
		d.Set("audio_max_size", objectRaw["AudioMaxSize"])
	}
	if objectRaw["AudioOpened"] != nil {
		d.Set("audio_opened", objectRaw["AudioOpened"])
	}
	if objectRaw["audioScanLimit"] != nil {
		d.Set("audio_scan_limit", objectRaw["audioScanLimit"])
	}
	if objectRaw["AutoFreezeType"] != nil {
		d.Set("auto_freeze_type", objectRaw["AutoFreezeType"])
	}
	if objectRaw["CallbackId"] != nil {
		d.Set("callback_id", objectRaw["CallbackId"])
	}
	if v, ok := d.GetOk("end_date"); ok {
		d.Set("end_date", v)
	} else if objectRaw["EndDate"] != nil {
		d.Set("end_date", objectRaw["EndDate"])
	}
	if objectRaw["ImageAutoFreezeOpened"] != nil {
		d.Set("image_auto_freeze_opened", objectRaw["ImageAutoFreezeOpened"])
	}
	if objectRaw["ImageOpened"] != nil {
		d.Set("image_opened", objectRaw["ImageOpened"])
	}
	if objectRaw["imageScanLimit"] != nil {
		d.Set("image_scan_limit", objectRaw["imageScanLimit"])
	}
	if objectRaw["ScanImageNoFileType"] != nil {
		d.Set("scan_image_no_file_type", objectRaw["ScanImageNoFileType"])
	}
	if objectRaw["StartDate"] != nil {
		d.Set("start_date", objectRaw["StartDate"])
	}
	if objectRaw["VideoAutoFreezeOpened"] != nil {
		d.Set("video_auto_freeze_opened", objectRaw["VideoAutoFreezeOpened"])
	}
	if objectRaw["VideoFrameInterval"] != nil {
		d.Set("video_frame_interval", objectRaw["VideoFrameInterval"])
	}
	if objectRaw["VideoMaxFrames"] != nil {
		d.Set("video_max_frames", objectRaw["VideoMaxFrames"])
	}
	if objectRaw["videoMaxSize"] != nil {
		d.Set("video_max_size", objectRaw["videoMaxSize"])
	}
	if objectRaw["VideoOpened"] != nil {
		d.Set("video_opened", objectRaw["VideoOpened"])
	}
	if objectRaw["videoScanLimit"] != nil {
		d.Set("video_scan_limit", objectRaw["videoScanLimit"])
	}

	bizTypeTemplate1RawObj, _ := jsonpath.Get("$.BizTypeTemplate", objectRaw)
	bizTypeTemplate1Raw := make(map[string]interface{})
	if bizTypeTemplate1RawObj != nil {
		bizTypeTemplate1Raw = bizTypeTemplate1RawObj.(map[string]interface{})
	}
	if bizTypeTemplate1Raw["BizType"] != nil {
		d.Set("biz_type", bizTypeTemplate1Raw["BizType"])
	}

	imageScenes1Raw := make([]interface{}, 0)
	if objectRaw["ImageScenes"] != nil {
		imageScenes1Raw = objectRaw["ImageScenes"].([]interface{})
	}

	d.Set("image_scenes", imageScenes1Raw)

	e := jsonata.MustCompile("$.AudioAntispamFreezeConfig")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("audio_antispam_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.AudioScenes")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("audio_scenes", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.Buckets")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("buckets", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.ImageAdFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("image_ad_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.ImageLiveFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("image_live_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.ImagePornFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("image_porn_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.ImageTerrorismFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("image_terrorism_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.VideoAdFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("video_ad_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.VideoLiveFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("video_live_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.VideoPornFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("video_porn_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.VideoScenes")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("video_scenes", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.VideoTerrorismFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("video_terrorism_freeze_config", convertObjectToJsonString(evaluation))
	e = jsonata.MustCompile("$.VideoVoiceAntispamFreezeConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("video_voice_antispam_freeze_config", convertObjectToJsonString(evaluation))

	return nil
}

func resourceAliCloudAligreenOssStockTaskDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Oss Stock Task. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
