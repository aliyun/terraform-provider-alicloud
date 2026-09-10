package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudOosParameter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosParameterCreate,
		Read:   resourceAlicloudOosParameterRead,
		Update: resourceAlicloudOosParameterUpdate,
		Delete: resourceAlicloudOosParameterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"constraints": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"has_value_wo": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"parameter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^ALIYUN.*)|(^ACS.*)|(^ALIBABA.*)|(^ALICLOUD.*)|(^OOS.*)`), "It cannot start with `ALIYUN`, `ACS`, `ALIBABA`, `ALICLOUD`, or `OOS`"), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_/-]{2,180}`), "The name must be `2` to `180` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/) and underscores (_).")),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"String", "StringList"}, false),
			},
			"value": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
				ExactlyOneOf: []string{"value", "value_wo"},
			},
			"value_wo": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				WriteOnly:    true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
				ExactlyOneOf: []string{"value", "value_wo"},
				RequiredWith: []string{"value_wo_version"},
			},
			"value_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"value_wo"},
			},
		},
	}
}

func resourceAlicloudOosParameterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateParameter"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("constraints"); ok {
		request["Constraints"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["Name"] = d.Get("parameter_name")
	request["Type"] = d.Get("type")
	value := d.Get("value").(string)
	if woValue, err := getWriteOnlyStringValue(d, cty.GetAttrPath("value_wo")); err != nil {
		return WrapError(err)
	} else if woValue != "" {
		value = woValue
	}
	request["Value"] = value
	request["ClientToken"] = buildClientToken("CreateParameter")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_parameter", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Name"]))

	return resourceAlicloudOosParameterRead(d, meta)
}
func resourceAlicloudOosParameterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosParameter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_parameter oosService.DescribeOosParameter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("parameter_name", d.Id())
	d.Set("constraints", object["Constraints"])
	d.Set("description", object["Description"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("type", object["Type"])
	hasValueWo := false
	if v, ok := d.GetOk("has_value_wo"); ok && v.(bool) {
		hasValueWo = true
	}
	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() {
		woValue, err := getWriteOnlyStringValue(d, cty.GetAttrPath("value_wo"))
		if err != nil {
			return WrapError(err)
		}
		hasValueWo = woValue != ""
	}
	if hasValueWo {
		d.Set("has_value_wo", true)
		d.Set("value", nil)
	} else {
		d.Set("has_value_wo", nil)
		d.Set("value", object["Value"])
	}
	return nil
}
func resourceAlicloudOosParameterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	value := d.Get("value").(string)
	if woValue, err := getWriteOnlyStringValue(d, cty.GetAttrPath("value_wo")); err != nil {
		return WrapError(err)
	} else if woValue != "" {
		value = woValue
	}
	if d.HasChange("value") {
		update = true
	}
	if d.HasChange("value_wo_version") {
		update = true
	}
	request["Value"] = value
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["Tags"] = respJson
		}
	}
	if update {
		action := "UpdateParameter"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudOosParameterRead(d, meta)
}
func resourceAlicloudOosParameterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteParameter"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
