package alicloud

import (
	"fmt"
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDatahubProjectName,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": &schema.Schema{
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

	_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		return nil, dataHubClient.CreateProject(projectName, projectComment)
	})
	if err != nil {
		return fmt.Errorf("failed to create project '%s' with error: %s", projectName, err)
	}

	d.SetId(strings.ToLower(projectName))
	return resourceAliyunDatahubProjectRead(d, meta)
}

func resourceAliyunDatahubProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	projectName := d.Id()
	raw, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
		return dataHubClient.GetProject(projectName)
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			d.SetId("")
		}
		return fmt.Errorf("failed to create project '%s' with error: %s", projectName, err)
	}
	project, _ := raw.(*datahub.Project)

	d.SetId(strings.ToLower(projectName))

	d.Set("name", projectName)
	d.Set("comment", project.Comment)
	d.Set("create_time", datahub.Uint64ToTimeString(project.CreateTime))
	d.Set("last_modify_time", datahub.Uint64ToTimeString(project.LastModifyTime))
	return nil
}

func resourceAliyunDatahubProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("comment") {
		projectName := d.Id()
		projectComment := d.Get("comment").(string)
		_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return nil, dataHubClient.UpdateProject(projectName, projectComment)
		})
		if err != nil {
			return fmt.Errorf("failed to update project '%s' with error: %s", projectName, err)
		}
	}

	return resourceAliyunDatahubProjectRead(d, meta)
}

func resourceAliyunDatahubProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	projectName := d.Id()
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return dataHubClient.GetProject(projectName)
		})
		if err != nil {
			if isDatahubNotExistError(err) {
				return nil
			}
			if isRetryableDatahubError(err) {
				return resource.RetryableError(fmt.Errorf("when deleting project '%s', failed to access it with error: %s", projectName, err))
			}
			return resource.NonRetryableError(fmt.Errorf("when deleting project '%s', failed to access it with error: %s", projectName, err))
		}

		_, err = client.WithDataHubClient(func(dataHubClient *datahub.DataHub) (interface{}, error) {
			return nil, dataHubClient.DeleteProject(projectName)
		})
		if err == nil || isDatahubNotExistError(err) {
			return nil
		}

		if isRetryableDatahubError(err) {
			return resource.RetryableError(fmt.Errorf("Deleting project '%s' timeout and got an error: %#v.", projectName, err))
		}

		return resource.NonRetryableError(fmt.Errorf("Deleting project '%s' timeout and got an error: %#v.", projectName, err))
	})

}
