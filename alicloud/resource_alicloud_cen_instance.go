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

func resourceAlicloudCenInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceCreate,
		Read:   resourceAlicloudCenInstanceRead,
		Update: resourceAlicloudCenInstanceUpdate,
		Delete: resourceAlicloudCenInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_instance_name": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"name": {
				Optional:   true,
				Computed:   true,
				Type:       schema.TypeString,
				Deprecated: "attribute 'name' has been deprecated from version 1.98.0. Use 'cen_instance_name' instead.",
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"protection_level": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudCenInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("cen_instance_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("protection_level"); ok {
		request["ProtectionLevel"] = v
	}

	request["ClientToken"] = buildClientToken("CreateCen")
	var response map[string]interface{}
	action := "CreateCen"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_instance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.CenId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cen_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenInstanceStateRefreshFunc(d, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCenInstanceUpdate(d, meta)
}
func resourceAlicloudCenInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_instance cbnService.DescribeCenInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("cen_instance_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("protection_level", object["ProtectionLevel"])
	d.Set("status", object["Status"])
	tagsMap := make(map[string]interface{})
	tagsRaw, _ := jsonpath.Get("$.Tags.Tag", object)
	if tagsRaw != nil {
		for _, value0 := range tagsRaw.([]interface{}) {
			tags := value0.(map[string]interface{})
			key := tags["Key"].(string)
			value := tags["Value"]
			if !ignoredTags(key, value) {
				tagsMap[key] = value
			}
		}
	}
	if len(tagsMap) > 0 {
		d.Set("tags", tagsMap)
	}

	return nil
}

func resourceAlicloudCenInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	d.Partial(true)
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"CenId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChanges("cen_instance_name", "name") {
		update = true
		if v, ok := d.GetOk("cen_instance_name"); ok {
			request["Name"] = v
		} else if v, ok := d.GetOk("name"); ok {
			request["Name"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("protection_level") {
		update = true
		if v, ok := d.GetOk("protection_level"); ok {
			request["ProtectionLevel"] = v
		}
	}

	if update {
		action := "ModifyCenAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cen_instance_name")
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("protection_level")
	}

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "cen"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudCenInstanceRead(d, meta)
}

func resourceAlicloudCenInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CenId": d.Id(),
	}

	action := "DeleteCen"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenInstanceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
