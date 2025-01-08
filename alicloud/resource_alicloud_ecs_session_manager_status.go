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

func resourceAliCloudEcsSessionManagerStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsSessionManagerStatusCreate,
		Read:   resourceAliCloudEcsSessionManagerStatusRead,
		Update: resourceAliCloudEcsSessionManagerStatusUpdate,
		Delete: resourceAliCloudEcsSessionManagerStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"session_manager_status_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"sessionManagerStatus"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
		},
	}
}

func resourceAliCloudEcsSessionManagerStatusCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "ModifyUserBusinessBehavior"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request["statusKey"] = d.Get("session_manager_status_name")
	request["statusValue"] = convertEcsSessionManagerStatusStatusRequest(d.Get("status"))

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_session_manager_status", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["statusKey"]))

	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(request["statusValue"])}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsSessionManagerStatusStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsSessionManagerStatusRead(d, meta)
}

func resourceAliCloudEcsSessionManagerStatusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsSessionManagerStatus(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_session_manager_status ecsService.DescribeEcsSessionManagerStatus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("session_manager_status_name", d.Id())
	d.Set("status", convertEcsSessionManagerStatusStatusResponse(object["StatusValue"]))

	return nil
}

func resourceAliCloudEcsSessionManagerStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}

	update := false
	request := map[string]interface{}{
		"statusKey": d.Id(),
	}

	if d.HasChange("status") {
		update = true
	}
	request["statusValue"] = convertEcsSessionManagerStatusStatusRequest(d.Get("status"))

	if update {
		action := "ModifyUserBusinessBehavior"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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

		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(request["statusValue"])}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsSessionManagerStatusStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudEcsSessionManagerStatusRead(d, meta)
}

func resourceAliCloudEcsSessionManagerStatusDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAliCloudEcsSessionManagerStatus. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertEcsSessionManagerStatusStatusRequest(source interface{}) interface{} {
	switch source {
	case "Enabled":
		return "enabled"
	case "Disabled":
		return "disabled"
	}

	return source
}

func convertEcsSessionManagerStatusStatusResponse(source interface{}) interface{} {
	switch source {
	case "enabled":
		return "Enabled"
	case "disabled":
		return "Disabled"
	}

	return source
}
