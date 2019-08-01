package alicloud

import (
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectCreate,
		Read:   resourceAlicloudLogProjectRead,
		Update: resourceAlicloudLogProjectUpdate,
		Delete: resourceAlicloudLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.CreateProject(d.Get("name").(string), d.Get("description").(string))
		})
		if err != nil {
			if IsExceptedError(err, LogClientTimeout) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateProject", raw)
		response, _ := raw.(*sls.LogProject)
		d.SetId(response.Name)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_project", "CreateProject", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	object, err := logService.DescribeLogProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.Name)
	d.Set("description", object.Description)

	return nil
}

func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("description") {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.UpdateProject(d.Get("name").(string), d.Get("description").(string))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProject", AliyunLogGoSdkERROR)
		}
		addDebug("UpdateProject", raw)
	}

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(d.Id())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{LogClientTimeout, LogRequestTimeout}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteProject", raw)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ProjectNotExist}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProject", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogProject(d.Id(), Deleted, DefaultTimeout))
}
