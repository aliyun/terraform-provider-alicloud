package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudResourceManagerSharedResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerSharedResourceCreate,
		Read:   resourceAliCloudResourceManagerSharedResourceRead,
		Update: resourceAliCloudResourceManagerSharedResourceUpdate,
		Delete: resourceAliCloudResourceManagerSharedResourceDelete,
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

func resourceAliCloudResourceManagerSharedResourceCreate(d *schema.ResourceData, meta interface{}) error {

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

	if v, ok := d.GetOk("resource_arn"); ok {
		localData, err := jsonpath.Get("$", v)
		if err != nil {
			return WrapError(err)
		}
		resourceArnsMapsArray := convertToInterfaceArray(localData)

		request["ResourceArns"] = resourceArnsMapsArray
	}

	if v, ok := d.GetOk("permission_name"); ok {
		localData1, err := jsonpath.Get("$", v)
		if err != nil {
			return WrapError(err)
		}
		permissionNamesMapsArray := convertToInterfaceArray(localData1)

		request["PermissionNames"] = permissionNamesMapsArray
	}

	// Only set Resources when resource_arn is not specified to avoid API conflict
	jsonString := convertObjectToJsonString(request)
	if _, ok := d.GetOk("resource_arn"); !ok {
		jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceId", d.Get("resource_id"))
		jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceType", d.Get("resource_type"))
	}
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_shared_resource", action, AlibabaCloudSdkGoERROR)
	}

	ResourceShareAssociationsResourceShareIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].ResourceShareId", response)
	ResourceShareAssociationsEntityIdVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].EntityId", response)
	ResourceShareAssociationsEntityTypeVar, _ := jsonpath.Get("$.ResourceShareAssociations[0].EntityType", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", ResourceShareAssociationsResourceShareIdVar, ResourceShareAssociationsEntityIdVar, ResourceShareAssociationsEntityTypeVar))

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedResourceStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudResourceManagerSharedResourceRead(d, meta)
}

func resourceAliCloudResourceManagerSharedResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerSharedResource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_shared_resource DescribeResourceManagerSharedResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("resource_arn", objectRaw["ResourceArn"])
	d.Set("status", objectRaw["AssociationStatus"])
	d.Set("resource_id", objectRaw["EntityId"])
	d.Set("resource_share_id", objectRaw["ResourceShareId"])
	d.Set("resource_type", objectRaw["EntityType"])

	return nil
}

func resourceAliCloudResourceManagerSharedResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Shared Resource.")
	return nil
}

func resourceAliCloudResourceManagerSharedResourceDelete(d *schema.ResourceData, meta interface{}) error {

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

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceId", parts[1])
	jsonString, _ = sjson.Set(jsonString, "Resources.0.ResourceType", parts[2])
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceSharing", "2020-01-10", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	resourceManagerServiceV2 := ResourceManagerServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Disassociated"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourceManagerServiceV2.ResourceManagerSharedResourceStateRefreshFunc(d.Id(), "AssociationStatus", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
