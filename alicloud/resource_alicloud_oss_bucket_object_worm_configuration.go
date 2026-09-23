package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketObjectWormConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketObjectWormConfigurationCreate,
		Read:   resourceAliCloudOssBucketObjectWormConfigurationRead,
		Update: resourceAliCloudOssBucketObjectWormConfigurationUpdate,
		Delete: resourceAliCloudOssBucketObjectWormConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"object_worm_enabled": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Enabled"}, false),
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_retention": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1, // only 1 retention is supported
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"years": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"mode": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"COMPLIANCE"}, false),
									},
									"days": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudOssBucketObjectWormConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?objectWorm")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket_name").(string))

	objectWormConfiguration := make(map[string]interface{})

	if v := d.Get("object_worm_enabled"); !IsNil(v) {
		objectWormConfiguration["ObjectWormEnabled"] = v
	}

	if v, ok := d.GetOk("rule"); ok {
		rule := expandOssBucketObjectWormConfigurationRule(v.([]interface{}))
		if rule != nil {
			objectWormConfiguration["Rule"] = rule
		}
	}

	request["ObjectWormConfiguration"] = objectWormConfiguration

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketObjectWormConfiguration", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_object_worm_configuration", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketObjectWormConfigurationRead(d, meta)
}

func resourceAliCloudOssBucketObjectWormConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketObjectWormConfiguration(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_object_worm_configuration DescribeOssBucketObjectWormConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("object_worm_enabled", objectRaw["ObjectWormEnabled"])

	if err := d.Set("rule", flattenOssBucketObjectWormConfigurationRule(objectRaw["Rule"])); err != nil {
		return WrapError(err)
	}

	d.Set("bucket_name", d.Id())

	return nil
}

func resourceAliCloudOssBucketObjectWormConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?objectWorm")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	if d.HasChange("object_worm_enabled") {
		update = true
	}
	if d.HasChange("rule") {
		update = true
	}
	objectWormConfiguration := make(map[string]interface{})

	if v, ok := d.GetOk("object_worm_enabled"); ok {
		objectWormConfiguration["ObjectWormEnabled"] = v
	}

	if v, ok := d.GetOk("rule"); ok {
		rule := expandOssBucketObjectWormConfigurationRule(v.([]interface{}))
		if rule != nil {
			objectWormConfiguration["Rule"] = rule
		}
	}

	request["ObjectWormConfiguration"] = objectWormConfiguration

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketObjectWormConfiguration", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssBucketObjectWormConfigurationRead(d, meta)
}

func resourceAliCloudOssBucketObjectWormConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Bucket Object Worm Configuration. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func expandOssBucketObjectWormConfigurationRule(ruleList []interface{}) map[string]interface{} {
	if len(ruleList) == 0 || ruleList[0] == nil {
		return nil
	}
	ruleRaw := ruleList[0].(map[string]interface{})
	rule := make(map[string]interface{})

	if v, ok := ruleRaw["default_retention"]; ok {
		retentionList := v.([]interface{})
		if len(retentionList) > 0 && retentionList[0] != nil {
			retRaw := retentionList[0].(map[string]interface{})
			retention := make(map[string]interface{})
			if days, ok := retRaw["days"]; ok {
				retention["Days"] = days
			}
			if mode, ok := retRaw["mode"]; ok {
				retention["Mode"] = mode
			}
			if years, ok := retRaw["years"]; ok {
				retention["Years"] = years
			}
			rule["DefaultRetention"] = retention
		}
	}

	return rule
}

func flattenOssBucketObjectWormConfigurationRule(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ruleRaw, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	ruleMap := make(map[string]interface{})

	if drRaw, ok := ruleRaw["DefaultRetention"]; ok && drRaw != nil {
		drList, ok := drRaw.([]interface{})
		if !ok {
			return nil
		}
		defaultRetentionMaps := make([]map[string]interface{}, 0, len(drList))
		for _, item := range drList {
			dr, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			defaultRetentionMaps = append(defaultRetentionMaps, map[string]interface{}{
				"days":  dr["Days"],
				"mode":  dr["Mode"],
				"years": dr["Years"],
			})
		}
		ruleMap["default_retention"] = defaultRetentionMaps
	}

	return []map[string]interface{}{ruleMap}
}
