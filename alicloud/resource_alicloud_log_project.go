package alicloud

import (
	"fmt"
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
	invoker := NewInvoker()
	invoker.AddCatcher(SlsClientTimeoutCatcher)
	if err := invoker.Run(func() error {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.CreateProject(d.Get("name").(string), d.Get("description").(string))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "log_project", "CreateProject", AliyunLogGoSdkERROR)
		}
		project, _ := raw.(*sls.LogProject)
		d.SetId(project.Name)
		return nil
	}); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()
	invoker.AddCatcher(SlsClientTimeoutCatcher)
	return invoker.Run(func() error {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetProject(d.Id())
		})
		if err != nil {
			if IsExceptedError(err, ProjectNotExist) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "log_project", "GetProject", AliyunLogGoSdkERROR)
		}
		project, _ := raw.(*sls.LogProject)
		d.Set("name", project.Name)
		d.Set("description", project.Description)

		return nil
	})
}

func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("description") {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.UpdateProject(d.Get("name").(string), d.Get("description").(string))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "log_project", "UpdateProject", AliyunLogGoSdkERROR)
		}
	}

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(d.Id())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{LogClientTimeout, LogRequestTimeout}) {
				return resource.RetryableError(BuildWrapError("DeleteProject", d.Id(), AliyunLogGoSdkERROR, err, ""))
			}
			if !IsExceptedErrors(err, []string{ProjectNotExist}) {
				return resource.NonRetryableError(BuildWrapError("DeleteProject", d.Id(), AliyunLogGoSdkERROR, err, ""))
			}
		}

		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.CheckProjectExist(d.Id())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{LogClientTimeout}) {
				return resource.RetryableError(BuildWrapError("CheckProjectExist", d.Id(), AliyunLogGoSdkERROR, err, ""))
			}
			return resource.NonRetryableError(BuildWrapError("CheckProjectExist", d.Id(), AliyunLogGoSdkERROR, err, ""))
		}
		exist, _ := raw.(bool)
		if !exist {
			return nil
		}

		return resource.RetryableError(BuildWrapError("DeleteProject", d.Id(), ProviderERROR, fmt.Errorf("Timeout"), ""))
	})
}
