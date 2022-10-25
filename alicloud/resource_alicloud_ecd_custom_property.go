package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcdCustomProperty() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdCustomPropertyCreate,
		Read:   resourceAlicloudEcdCustomPropertyRead,
		Update: resourceAlicloudEcdCustomPropertyUpdate,
		Delete: resourceAlicloudEcdCustomPropertyDelete,
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

func resourceAlicloudEcdCustomPropertyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateProperty"
	request := make(map[string]interface{})
	conn, err := client.NewEdsuserClient()
	if err != nil {
		return WrapError(err)
	}
	request["PropertyKey"] = d.Get("property_key")

	if v, ok := d.GetOk("property_values"); ok {
		propertyValuesList := make([]string, 0)
		for _, propertyValue := range v.(*schema.Set).List() {
			propertyValueArg := propertyValue.(map[string]interface{})
			propertyValuesList = append(propertyValuesList, propertyValueArg["property_value"].(string))
		}
		request["PropertyValues"] = propertyValuesList
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_custom_property", action, AlibabaCloudSdkGoERROR)
	}
	responseCreateResult := response["CreateResult"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseCreateResult["PropertyId"]))

	return resourceAlicloudEcdCustomPropertyRead(d, meta)
}
func resourceAlicloudEcdCustomPropertyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edsUserService := EdsUserService{client}
	object, err := edsUserService.DescribeEcdCustomProperty(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_custom_property edsUserService.DescribeEcdCustomProperty Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
func resourceAlicloudEcdCustomPropertyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEdsuserClient()
	if err != nil {
		return WrapError(err)
	}
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudEcdCustomPropertyRead(d, meta)
}
func resourceAlicloudEcdCustomPropertyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveProperty"
	var response map[string]interface{}
	conn, err := client.NewEdsuserClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PropertyId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-03-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
