package datahub

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "project added by terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"create_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	projectName := d.Get("name").(string)
	projectComment := d.Get("comment").(string)

	var requestInfo *datahub.DataHub

	raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
		requestInfo = dataHubClient.(*datahub.DataHub)

		return dataHubClient.CreateProject(projectName, projectComment)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_datahub_project", "CreateProject", errs.AliyunDatahubSdkGo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["ProjectComment"] = projectComment
		logs.AddDebug("CreateProject", raw, requestInfo, requestMap)
	}

	d.SetId(strings.ToLower(projectName))
	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}
	object, err := datahubService.DescribeDatahubProject(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.SetId(strings.ToLower(d.Id()))

	d.Set("name", d.Id())
	d.Set("comment", object.Comment)
	d.Set("create_time", strconv.FormatInt(object.CreateTime, 10))
	d.Set("last_modify_time", strconv.FormatInt(object.LastModifyTime, 10))
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("comment") {

		projectName := d.Id()
		projectComment := d.Get("comment").(string)

		var requestInfo *datahub.DataHub

		raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
			requestInfo = dataHubClient.(*datahub.DataHub)
			return dataHubClient.UpdateProject(projectName, projectComment)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UpdateProject", errs.AliyunDatahubSdkGo)
		}
		if logs.DebugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["ProjectComment"] = projectComment
			logs.AddDebug("UpdateProject", raw, requestInfo, requestMap)
		}
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			datahubService := DatahubService{client}
			object, err := datahubService.DescribeDatahubProject(d.Id())
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if object.Comment != projectComment {
				return retry.RetryableError(fmt.Errorf("waiting for updating project %s comment finished timwout. "+
					"current comment is %s ", d.Id(), object.Comment))
			}
			return nil
		})
		if err != nil {
			return errs.WrapError(err)
		}
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}

	projectName := d.Id()

	var requestInfo *datahub.DataHub

	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
			requestInfo = dataHubClient.(*datahub.DataHub)
			return dataHubClient.DeleteProject(projectName)
		})
		if err != nil {
			if isRetryableDatahubError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			logs.AddDebug("DeleteProject", raw, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteProject", errs.AliyunDatahubSdkGo)
	}
	return errs.WrapError(datahubService.WaitForDatahubProject(d.Id(), Deleted, DefaultTimeout))
}
