// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksDataServiceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksDataServiceGroupCreate,
		Read:   resourceAliCloudDataWorksDataServiceGroupRead,
		Delete: resourceAliCloudDataWorksDataServiceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"api_gateway_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_service_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_service_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "This field is deprecated.",
			},
		},
	}
}

func resourceAliCloudDataWorksDataServiceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDataServiceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		request["ProjectId"] = v
	}
	request["RegionId"] = client.RegionId

	request["ApiGatewayGroupId"] = d.Get("api_gateway_group_id")
	request["GroupName"] = d.Get("data_service_group_name")
	request["Description"] = d.Get("description")
	if v, ok := d.GetOkExists("tenant_id"); ok {
		request["TenantId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_data_service_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ProjectId"], response["GroupId"]))

	return resourceAliCloudDataWorksDataServiceGroupRead(d, meta)
}

func resourceAliCloudDataWorksDataServiceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDataServiceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_data_service_group DescribeDataWorksDataServiceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("api_gateway_group_id", objectRaw["ApiGatewayGroupId"])
	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("creator_id", objectRaw["CreatorId"])
	d.Set("data_service_group_name", objectRaw["GroupName"])
	d.Set("description", objectRaw["Description"])
	d.Set("modified_time", objectRaw["ModifiedTime"])
	d.Set("tenant_id", objectRaw["TenantId"])
	d.Set("project_id", objectRaw["ProjectId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("data_service_group_id", parts[1])

	return nil
}

func resourceAliCloudDataWorksDataServiceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Data Service Group. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
