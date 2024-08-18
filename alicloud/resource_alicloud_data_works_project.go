// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudDataWorksProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksProjectCreate,
		Read:   resourceAliCloudDataWorksProjectRead,
		Update: resourceAliCloudDataWorksProjectUpdate,
		Delete: resourceAliCloudDataWorksProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudDataWorksProjectCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateProject"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("project_mode"); ok {
		request["ProjectMode"] = v
	}
	request["ProjectName"] = d.Get("display_name")
	request["ProjectIdentifier"] = d.Get("project_name")
	request["ProjectDescription"] = d.Get("description")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.System"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_project", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	dataWorksServiceV2 := DataWorksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, dataWorksServiceV2.DataWorksProjectStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDataWorksProjectUpdate(d, meta)
}

func resourceAliCloudDataWorksProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksProject(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_project DescribeDataWorksProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["GmtCreate"] != nil {
		d.Set("create_time", objectRaw["GmtCreate"])
	}
	if objectRaw["ProjectDescription"] != nil {
		d.Set("description", objectRaw["ProjectDescription"])
	}
	if objectRaw["ProjectName"] != nil {
		d.Set("display_name", objectRaw["ProjectName"])
	}
	if objectRaw["ProjectMode"] != nil {
		d.Set("project_mode", objectRaw["ProjectMode"])
	}
	if objectRaw["ProjectIdentifier"] != nil {
		d.Set("project_name", objectRaw["ProjectIdentifier"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	return nil
}

func resourceAliCloudDataWorksProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateProject"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ProjectId"] = d.Id()

	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}

	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
	}
	request["ProjectName"] = d.Get("display_name")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["ProjectDescription"] = d.Get("description")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), query, request, &runtime)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		dataWorksServiceV2 := DataWorksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"0", "4", "7"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, dataWorksServiceV2.DataWorksProjectStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudDataWorksProjectRead(d, meta)
}

func resourceAliCloudDataWorksProjectDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteProject"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["ProjectId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.System"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"Invalid.Tenant.ProjectNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	dataWorksServiceV2 := DataWorksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, dataWorksServiceV2.DataWorksProjectStateRefreshFunc(d.Id(), "ProjectId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
