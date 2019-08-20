package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDatahubProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDatahubProjectCreate,
		Read:   resourceAliyunDatahubProjectRead,
		Update: resourceAliyunDatahubProjectUpdate,
		Delete: resourceAliyunDatahubProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "project added by terraform",
				ValidateFunc: validateStringLengthInRange(0, 255),
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

func resourceAliyunDatahubProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	projectName := d.Get("name").(string)
	projectComment := d.Get("comment").(string)

	var requestInfo *datahub.DataHub

	raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		requestInfo = dataHubClient
		return nil, dataHubClient.CreateProject(projectName, projectComment)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_datahub_project", "CreateProject", AliyunDatahubSdkGo)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["ProjectComment"] = projectComment
		addDebug("CreateProject", raw, requestInfo, requestMap)
	}

	d.SetId(strings.ToLower(projectName))
	return resourceAliyunDatahubProjectRead(d, meta)
}

func resourceAliyunDatahubProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}
	object, err := datahubService.DescribeDatahubProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(strings.ToLower(d.Id()))

	d.Set("name", d.Id())
	d.Set("comment", object.Comment)
	d.Set("create_time", datahub.Uint64ToTimeString(object.CreateTime))
	d.Set("last_modify_time", datahub.Uint64ToTimeString(object.LastModifyTime))
	return nil
}

func resourceAliyunDatahubProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("comment") {

		projectName := d.Id()
		projectComment := d.Get("comment").(string)

		var requestInfo *datahub.DataHub

		raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			requestInfo = dataHubClient
			return nil, dataHubClient.UpdateProject(projectName, projectComment)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProject", AliyunDatahubSdkGo)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["ProjectComment"] = projectComment
			addDebug("UpdateProject", raw, requestInfo, requestMap)
		}
	}

	return resourceAliyunDatahubProjectRead(d, meta)
}

func resourceAliyunDatahubProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	datahubService := DatahubService{client}

	projectName := d.Id()

	var requestInfo *datahub.DataHub

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			requestInfo = dataHubClient
			return nil, dataHubClient.DeleteProject(projectName)
		})
		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			addDebug("DeleteProject", raw, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProject", AliyunDatahubSdkGo)
	}
	return WrapError(datahubService.WaitForDatahubProject(d.Id(), Deleted, DefaultTimeout))
}
