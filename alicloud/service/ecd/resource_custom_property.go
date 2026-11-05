package ecd

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
)

func ResourceCustomProperty() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomPropertyCreate,
		Read:   resourceCustomPropertyRead,
		Update: resourceCustomPropertyUpdate,
		Delete: resourceCustomPropertyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"property_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_values": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"property_value_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCustomPropertyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateProperty"
	request := make(map[string]interface{})
	var err error
	request["PropertyKey"] = d.Get("property_key")

	if v, ok := d.GetOk("property_values"); ok {
		propertyValuesList := make([]string, 0)
		for _, propertyValue := range v.(*schema.Set).List() {
			propertyValueArg := propertyValue.(map[string]interface{})
			propertyValuesList = append(propertyValuesList, propertyValueArg["property_value"].(string))
		}
		request["PropertyValues"] = propertyValuesList
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_custom_property", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseCreateResult := response["CreateResult"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseCreateResult["PropertyId"]))

	return resourceCustomPropertyRead(d, meta)
}
func resourceCustomPropertyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edsUserService := EdsUserService{client}
	object, err := edsUserService.DescribeEcdCustomProperty(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_custom_property edsUserService.DescribeEcdCustomProperty Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("property_key", object["PropertyKey"])
	if propertyValuesList, ok := object["PropertyValues"]; ok && propertyValuesList != nil {
		propertyValuesMaps := make([]map[string]interface{}, 0)
		for _, propertyValuesListItem := range propertyValuesList.([]interface{}) {
			if propertyValuesListItemMapArg, ok := propertyValuesListItem.(map[string]interface{}); ok {
				propertyValuesMaps = append(propertyValuesMaps, map[string]interface{}{
					"property_value":    propertyValuesListItemMapArg["PropertyValue"],
					"property_value_id": propertyValuesListItemMapArg["PropertyValueId"],
				})
			}
			d.Set("property_values", propertyValuesMaps)
		}
	}

	return nil
}
func resourceCustomPropertyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PropertyId": d.Id(),
	}

	if d.HasChange("property_key") {
		update = true
	}
	request["PropertyKey"] = d.Get("property_key")
	if d.HasChange("property_values") {
		update = true
	}
	if v, ok := d.GetOk("property_values"); ok {
		for propertyValuesPtr, propertyValues := range v.(*schema.Set).List() {
			propertyValuesArg := propertyValues.(map[string]interface{})
			request["PropertyValues."+fmt.Sprint(propertyValuesPtr+1)+".PropertyValue"] = propertyValuesArg["property_value"]
		}
	}
	if update {
		action := "UpdateProperty"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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
	return resourceCustomPropertyRead(d, meta)
}
func resourceCustomPropertyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveProperty"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PropertyId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eds-user", "2021-03-08", action, nil, request, false)
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
	return nil
}
