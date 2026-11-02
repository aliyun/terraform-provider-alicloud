package sls

import (
	"fmt"
	"log"
	"regexp"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "The name of the log project. It is the only in one Alicloud account. The project name is globally unique in Alibaba Cloud and cannot be modified after it is created. The naming rules are as follows:- The project name must be globally unique. - The name can contain only lowercase letters, digits, and hyphens (-). - It must start and end with a lowercase letter or number. - The value contains 3 to 63 characters."),
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
			"tags": tags.TagsSchema(),
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "Field 'name' has been deprecated since provider version 1.223.0. New field 'project_name' instead.",
				ForceNew:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile("^[0-9a-zA-Z_-]+$"), "The name of the log project. It is the only in one Alicloud account. The project name is globally unique in Alibaba Cloud and cannot be modified after it is created. The naming rules are as follows:- The project name must be globally unique. - The name can contain only lowercase letters, digits, and hyphens (-). - It must start and end with a lowercase letter or number. - The value contains 3 to 63 characters."),
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
	var err error
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Sls", helper.RoaParam("POST", "2020-12-30", "CreateProject", action), query, body, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_project", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["projectName"]))

	return resourceAliCloudSlsProjectUpdate(d, meta)
}

func resourceAliCloudSlsProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := Service{client}

	objectRaw, err := slsServiceV2.DescribeSlsProject(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_project DescribeSlsProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("status", objectRaw["status"])
	d.Set("project_name", objectRaw["projectName"])

	objectRaw, err = slsServiceV2.DescribeListTagResources(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	tagsMaps := objectRaw["tagResources"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	d.Set("project_name", d.Id())

	logService := Service{client}
	policy, err := logService.DescribeLogProjectPolicy(d.Id())
	if err != nil {
		return errs.WrapError(err)
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
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["project"] = helper.StringPointer(d.Id())
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["description"] = d.Get("description")
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Sls", helper.RoaParam("PUT", "2020-12-30", "UpdateProject", action), query, body, nil, hostMap, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = fmt.Sprintf("/resourcegroup")
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Sls", helper.RoaParam("PUT", "2020-12-30", "ChangeResourceGroup", action), query, body, nil, hostMap, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		slsServiceV2 := Service{client}
		if err := slsServiceV2.SetResourceTags(d, "PROJECT"); err != nil {
			return errs.WrapError(err)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UpdateProjectPolicy", errs.AliyunLogGoSdkERROR)
		}
		logs.AddDebug("UpdateProjectPolicy", raw, requestInfo, request)
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
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = helper.StringPointer(d.Id())

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Sls", helper.RoaParam("DELETE", "2020-12-30", "DeleteProject", action), query, nil, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
