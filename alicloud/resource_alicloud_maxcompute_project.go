package alicloud

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/maxcompute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudMaxComputeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunMaxComputeProjectCreate,
		Read:   resourceAliyunMaxComputeProjectRead,
		Delete: resourceAliyunMaxComputeProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 27),
			},

			"specification_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"OdpsStandard"}, false),
			},

			"order_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
			},
		},
	}
}

func resourceAliyunMaxComputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := maxcompute.CreateCreateProjectRequest()

	request.OdpsRegionId = client.RegionId
	request.ProjectName = d.Get("name").(string)
	request.OdpsSpecificationType = d.Get("specification_type").(string)
	request.OrderType = d.Get("order_type").(string)

	raw, err := client.WithMaxComputeClient(func(MaxComputeClient *maxcompute.Client) (interface{}, error) {
		return MaxComputeClient.CreateProject(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_maxcompute_project", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*maxcompute.CreateProjectResponse)

	if response.Code != "200" {
		return WrapError(Error("%v", response))
	}

	d.SetId(request.ProjectName)

	return resourceAliyunMaxComputeProjectRead(d, meta)
}

func resourceAliyunMaxComputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxcomputeService := MaxComputeService{client}
	response, err := maxcomputeService.DescribeMaxComputeProject(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(response.Data), &dat); err != nil {
		return WrapError(Error("%v", response))
	}
	d.Set("order_type", dat["orderType"].(string))
	d.Set("name", dat["projectName"].(string))

	return nil
}

func resourceAliyunMaxComputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	maxcomputeService := MaxComputeService{client}

	request := maxcompute.CreateDeleteProjectRequest()

	request.RegionIdName = client.RegionId
	request.ProjectName = d.Get("name").(string)

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithMaxComputeClient(func(MaxComputeClient *maxcompute.Client) (interface{}, error) {
			return MaxComputeClient.DeleteProject(request)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}

		response := raw.(*maxcompute.DeleteProjectResponse)
		if response.Code == "500" {
			return resource.RetryableError(nil)
		}

		if response.Code != "200" {
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if isProjectNotExistError(response.Data) {
			return nil
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.ProjectName, "DeleteProject", AliyunMaxComputeSdkGo)
	}
	return WrapError(maxcomputeService.WaitForMaxComputeProject(request.ProjectName, Deleted, DefaultTimeout))

}
