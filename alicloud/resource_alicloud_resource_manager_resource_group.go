package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudResourceManagerResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerResourceGroupCreate,
		Read:   resourceAliCloudResourceManagerResourceGroupRead,
		Update: resourceAliCloudResourceManagerResourceGroupUpdate,
		Delete: resourceAliCloudResourceManagerResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"resource_group_name"},
				Deprecated:    "Field `name` has been deprecated from provider version 1.114.0. New field `resource_group_name` instead.",
			},
			"tags": tagsSchema(),
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_statuses": {
				Type:     schema.TypeList,
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
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'create_date' has been removed from provider version 1.114.0.",
			},
		},
	}
}

func resourceAliCloudResourceManagerResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	action := "CreateResourceGroup"
	request := make(map[string]interface{})
	var err error

	request["DisplayName"] = d.Get("display_name")

	if v, ok := d.GetOk("resource_group_name"); ok {
		request["Name"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "resource_group_name" or "name" must be set one!`))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tag"] = tagsMap
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.ResourceGroup", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_resource_manager_resource_group")
	} else {
		resourceGroupId := resp.(map[string]interface{})["Id"]
		d.SetId(fmt.Sprint(resourceGroupId))
	}

	stateConf := BuildStateConf([]string{}, []string{"OK"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, resourceManagerService.ResourceManagerResourceGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudResourceManagerResourceGroupRead(d, meta)
}

func resourceAliCloudResourceManagerResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}

	object, err := resourceManagerService.DescribeResourceManagerResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_group resourceManagerService.DescribeResourceManagerResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("display_name", object["DisplayName"])
	d.Set("resource_group_name", object["Name"])
	d.Set("name", object["Name"])
	d.Set("account_id", object["AccountId"])
	d.Set("status", object["Status"])

	listTagResourcesObject, err := resourceManagerService.ListTagResources(d.Id(), "ResourceGroup")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	if regionStatuses, ok := object["RegionStatuses"]; ok {
		if regionStatusList, ok := regionStatuses.(map[string]interface{})["RegionStatus"]; ok {
			regionStatusesMaps := make([]map[string]interface{}, 0)
			for _, regionStatus := range regionStatusList.([]interface{}) {
				regionStatusArg := regionStatus.(map[string]interface{})
				regionStatusMap := map[string]interface{}{}

				if regionId, ok := regionStatusArg["RegionId"]; ok {
					regionStatusMap["region_id"] = regionId
				}

				if status, ok := regionStatusArg["Status"]; ok {
					regionStatusMap["status"] = status
				}

				regionStatusesMaps = append(regionStatusesMaps, regionStatusMap)
			}

			d.Set("region_statuses", regionStatusesMaps)
		}
	}

	return nil
}

func resourceAliCloudResourceManagerResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"ResourceGroupId": d.Id(),
	}

	if d.HasChange("display_name") {
		update = true
	}
	request["NewDisplayName"] = d.Get("display_name")

	if update {
		action := "UpdateResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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

		d.SetPartial("display_name")
	}

	if d.HasChange("tags") {
		if err := resourceManagerService.SetResourceTags(d, "ResourceGroup"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliCloudResourceManagerResourceGroupRead(d, meta)
}

func resourceAliCloudResourceManagerResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResourceGroup"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"ResourceGroupId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.ResourceGroup.Resource"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
