// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsProjectCreate,
		Read:   resourceAliCloudSlsProjectRead,
		Update: resourceAliCloudSlsProjectUpdate,
		Delete: resourceAliCloudSlsProjectDelete,
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
				Optional: true,
			},
			"project_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"project_name", "name"},
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "The name of the log project. It is the only in one Alicloud account. The project name is globally unique in Alibaba Cloud and cannot be modified after it is created. The naming rules are as follows:- The project name must be globally unique. - The name can contain only lowercase letters, digits, and hyphens (-). - It must start and end with a lowercase letter or number. - The value contains 3 to 63 characters."),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'name' has been deprecated since provider version 1.223.0. New field 'project_name' instead.",
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "The name of the log project. It is the only in one Alicloud account. The project name is globally unique in Alibaba Cloud and cannot be modified after it is created. The naming rules are as follows:- The project name must be globally unique. - The name can contain only lowercase letters, digits, and hyphens (-). - It must start and end with a lowercase letter or number. - The value contains 3 to 63 characters."),
			},
		},
	}
}

func resourceAliCloudSlsProjectCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("project_name"); ok {
		request["projectName"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["projectName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(client.GenRoaParam("CreateProject", "POST", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_project", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["projectName"]))

	return resourceAliCloudSlsProjectUpdate(d, meta)
}

func resourceAliCloudSlsProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsProject(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_project DescribeSlsProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("status", objectRaw["status"])
	d.Set("project_name", objectRaw["projectName"])

	objectRaw, err = slsServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps := objectRaw["tagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("project_name", d.Id())

	logService := LogService{client}
	policy, err := logService.DescribeLogProjectPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if policy != "" && policy != "{}" {
		d.Set("policy", policy)
	}
	d.Set("name", d.Get("project_name"))
	return nil
}

func resourceAliCloudSlsProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	var hostMap map[string]*string
	update := false
	d.Partial(true)
	action := fmt.Sprintf("/")
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["project"] = StringPointer(d.Id())
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["description"] = d.Get("description")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(client.GenRoaParam("UpdateProject", "PUT", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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
	}
	update = false
	action = fmt.Sprintf("/resourcegroup")
	conn, err = client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	request["resourceId"] = d.Id()
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["resourceGroupId"] = d.Get("resource_group_id")
	request["resourceType"] = "PROJECT"
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(client.GenRoaParam("ChangeResourceGroup", "PUT", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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
	}

	if d.HasChange("tags") {
		slsServiceV2 := SlsServiceV2{client}
		if err := slsServiceV2.SetResourceTags(d, "PROJECT"); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("policy") {
		var requestInfo *sls.Client
		policy := ""
		if v, ok := d.GetOk("policy"); ok {
			policy = v.(string)
		}
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			if policy == "" {
				return nil, slsClient.DeleteProjectPolicy(d.Id())
			}
			return nil, slsClient.UpdateProjectPolicy(d.Id(), policy)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProjectPolicy", AliyunLogGoSdkERROR)
		}
		addDebug("UpdateProjectPolicy", raw, requestInfo, request)
		d.SetPartial("policy")
	}
	d.Partial(false)
	return resourceAliCloudSlsProjectRead(d, meta)
}

func resourceAliCloudSlsProjectDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Id())

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.Execute(client.GenRoaParam("DeleteProject", "DELETE", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: nil, HostMap: hostMap}, &util.RuntimeOptions{})

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

	return nil
}
