// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dev_environment_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"dev_role_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\\w.,;/@-]+$"), "Workspace Display Name"),
			},
			"pai_task_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\\w.,;/@-]+$"), "Workspace Name"),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Available", "Initializing", "Forbidden", "InitFailed", "Deleting", "DeleteFailed", "Frozen", "Updating", "UpdateFailed"}, false),
			},
			"tags": tagsSchema(),
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
	request["RegionId"] = client.RegionId

	request["Name"] = d.Get("project_name")
	request["DisplayName"] = d.Get("display_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["AliyunResourceGroupId"] = v
	}
	request["PaiTaskEnabled"] = d.Get("pai_task_enabled")
	if v, ok := d.GetOkExists("dev_environment_enabled"); ok {
		request["DevEnvironmentEnabled"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOkExists("dev_role_disabled"); ok {
		request["DevRoleDisabled"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"9990020002", "9990040003"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_project", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ProjectId"]))

	dataWorksServiceV2 := DataWorksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, dataWorksServiceV2.DataWorksProjectStateRefreshFunc(d.Id(), "$.Project.Status", []string{}))
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

	project1RawObj, _ := jsonpath.Get("$.Project", objectRaw)
	project1Raw := make(map[string]interface{})
	if project1RawObj != nil {
		project1Raw = project1RawObj.(map[string]interface{})
	}
	if project1Raw["Description"] != nil {
		d.Set("description", project1Raw["Description"])
	}
	if project1Raw["DevEnvironmentEnabled"] != nil {
		d.Set("dev_environment_enabled", project1Raw["DevEnvironmentEnabled"])
	}
	if project1Raw["DevRoleDisabled"] != nil {
		d.Set("dev_role_disabled", project1Raw["DevRoleDisabled"])
	}
	if project1Raw["DisplayName"] != nil {
		d.Set("display_name", project1Raw["DisplayName"])
	}
	if project1Raw["PaiTaskEnabled"] != nil {
		d.Set("pai_task_enabled", project1Raw["PaiTaskEnabled"])
	}
	if project1Raw["Name"] != nil {
		d.Set("project_name", project1Raw["Name"])
	}
	if project1Raw["AliyunResourceGroupId"] != nil {
		d.Set("resource_group_id", project1Raw["AliyunResourceGroupId"])
	}
	if project1Raw["Status"] != nil {
		d.Set("status", project1Raw["Status"])
	}

	tagsMaps, _ := jsonpath.Get("$.Project.AliyunResourceTags", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudDataWorksProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ChangeResourceManagerResourceGroup"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "project"
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceManagerResourceGroupId"] = d.Get("resource_group_id")
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "UpdateProject"
	conn, err = client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("status") {
		update = true
		request["Status"] = d.Get("status")
	}

	if !d.IsNewResource() && d.HasChange("dev_environment_enabled") {
		update = true
		request["DevEnvironmentEnabled"] = d.Get("dev_environment_enabled")
	}

	if !d.IsNewResource() && d.HasChange("dev_role_disabled") {
		update = true
		request["DevRoleDisabled"] = d.Get("dev_role_disabled")
	}

	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
	}
	request["DisplayName"] = d.Get("display_name")
	if !d.IsNewResource() && d.HasChange("pai_task_enabled") {
		update = true
	}
	request["PaiTaskEnabled"] = d.Get("pai_task_enabled")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
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
		dataWorksServiceV2 := DataWorksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available", "Forbidden"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, dataWorksServiceV2.DataWorksProjectStateRefreshFunc(d.Id(), "$.Project.Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		dataWorksServiceV2 := DataWorksServiceV2{client}
		if err := dataWorksServiceV2.SetResourceTags(d, "project"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
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
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)

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
		if IsExpectedErrors(err, []string{"1101080008"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	dataWorksServiceV2 := DataWorksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, dataWorksServiceV2.DataWorksProjectStateRefreshFunc(d.Id(), "$.Project.Id", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
