package resourcemanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

func ResourceSharedResource() *schema.Resource {
	return &schema.Resource{
		Create: ResourceSharedResourceCreate,
		Read:   ResourceSharedResourceRead,
		Update: ResourceSharedResourceUpdate,
		Delete: ResourceSharedResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"properties_resource_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_property": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceSharedResourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateResourceShare"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("resource_share_id"); ok {
		request["ResourceShareId"] = v
	}
	request["RegionId"] = client.RegionId

	resourcePropertiesDataList := make(map[string]interface{})

	if v, ok := d.GetOk("resource_property"); ok {
		resourcePropertiesDataList["Property"] = v
	}

	if v, ok := d.GetOk("properties_resource_arn"); ok {
		resourcePropertiesDataList["ResourceArn"] = v
	}

	ResourcePropertiesMap := make([]interface{}, 0)
	ResourcePropertiesMap = append(ResourcePropertiesMap, resourcePropertiesDataList)
	request["ResourceProperties"] = ResourcePropertiesMap

	if v, ok := d.GetOk("resource_arn"); ok {
		localData, err := jsonpath.Get("$", v)
		if err != nil {
			return errs.WrapError(err)
		}
		resourceArnsMapsArray := convert.ConvertToInterfaceArray(localData)

		request["ResourceArns"] = resourceArnsMapsArray
	}

	if v, ok := d.GetOk("permission_name"); ok {
		localData1, err := jsonpath.Get("$", v)
		if err != nil {
			return errs.WrapError(err)
		}
		permissionNamesMapsArray := convert.ConvertToInterfaceArray(localData1)

		request["PermissionNames"] = permissionNamesMapsArray
	}

	// Only set Resources when resource_arn is not specified to avoid API conflict
	jsonString := convert.ConvertObjectToJsonString(request)
	if _, ok := d.GetOk("resource_arn"); !ok {
		jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceId", d.Get("resource_id"))
		jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceType", d.Get("resource_type"))
	}
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_shared_resource", action, errs.AlibabaCloudSdkGoERROR)
	}

	ResourceShareAssociationsResourceShareIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].ResourceShareId", response)
	ResourceShareAssociationsEntityIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].EntityId", response)
	ResourceShareAssociationsEntityTypeVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].EntityType", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", ResourceShareAssociationsResourceShareIdVar, ResourceShareAssociationsEntityIdVar, ResourceShareAssociationsEntityTypeVar))

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedResourceStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceSharedResourceRead(d, meta)
}

func ResourceSharedResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerSharedResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_shared_resource DescribeResourceManagerSharedResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("resource_arn", objectRaw["ResourceArn"])
	d.Set("resource_property", objectRaw["ResourceProperty"])
	d.Set("status", objectRaw["AssociationStatus"])
	d.Set("resource_id", objectRaw["EntityId"])
	d.Set("resource_share_id", objectRaw["ResourceShareId"])
	d.Set("resource_type", objectRaw["EntityType"])

	return nil
}

func ResourceSharedResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Shared Resource.")
	return nil
}

func ResourceSharedResourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DisassociateResourceShare"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ResourceShareId"] = parts[0]
	request["RegionId"] = client.RegionId

	onlyArnType := []string{"PolarDBBackupSet"}

	if helper.InArray(parts[2], onlyArnType) {
		if v, ok := d.GetOk("resource_arn"); ok {
			localData, _ := jsonpath.Get("$", v)
			resourceArnsMapsArray := convert.ConvertToInterfaceArray(localData)

			request["ResourceArns"] = resourceArnsMapsArray
		}
	} else {
		jsonString := convert.ConvertObjectToJsonString(request)
		jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceId", parts[1])
		jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceType", parts[2])
		_ = json.Unmarshal([]byte(jsonString), &request)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Disassociated"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedResourceStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
