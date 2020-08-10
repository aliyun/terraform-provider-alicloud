package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudResourceManagerResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerResourceGroupCreate,
		Read:   resourceAlicloudResourceManagerResourceGroupRead,
		Update: resourceAlicloudResourceManagerResourceGroupUpdate,
		Delete: resourceAlicloudResourceManagerResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_statuses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}

	request := resourcemanager.CreateCreateResourceGroupRequest()
	request.DisplayName = d.Get("display_name").(string)
	request.Name = d.Get("name").(string)
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.CreateResourceGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.CreateResourceGroupResponse)
	d.SetId(fmt.Sprintf("%v", response.ResourceGroup.Id))
	stateConf := BuildStateConf([]string{}, []string{"OK"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourcemanagerService.ResourceManagerResourceGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudResourceManagerResourceGroupRead(d, meta)
}
func resourceAlicloudResourceManagerResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerResourceGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_id", object.AccountId)
	d.Set("create_date", object.CreateDate)
	d.Set("display_name", object.DisplayName)
	d.Set("name", object.Name)

	regionStatus := make([]map[string]interface{}, len(object.RegionStatuses.RegionStatus))
	for i, v := range object.RegionStatuses.RegionStatus {
		regionStatus[i] = make(map[string]interface{})
		regionStatus[i]["region_id"] = v.RegionId
		regionStatus[i]["status"] = v.Status
	}
	if err := d.Set("region_statuses", regionStatus); err != nil {
		return WrapError(err)
	}
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudResourceManagerResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := resourcemanager.CreateUpdateResourceGroupRequest()
	request.ResourceGroupId = d.Id()
	if d.HasChange("display_name") {
		update = true
	}
	request.NewDisplayName = d.Get("display_name").(string)
	if update {
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.UpdateResourceGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudResourceManagerResourceGroupRead(d, meta)
}
func resourceAlicloudResourceManagerResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := resourcemanager.CreateDeleteResourceGroupRequest()
	request.ResourceGroupId = d.Id()
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.DeleteResourceGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
