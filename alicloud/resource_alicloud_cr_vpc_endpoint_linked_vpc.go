package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCrVpcEndpointLinkedVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrVpcEndpointLinkedVpcCreate,
		Read:   resourceAlicloudCrVpcEndpointLinkedVpcRead,
		Update: resourceAlicloudCrVpcEndpointLinkedVpcUpdate,
		Delete: resourceAlicloudCrVpcEndpointLinkedVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"module_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Registry", "Chart"}, false),
			},
			"enable_create_dns_record_in_pvzt": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCrVpcEndpointLinkedVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	var response map[string]interface{}
	action := "CreateInstanceVpcEndpointLinkedVpc"
	request := make(map[string]interface{})
	var err error

	request["InstanceId"] = d.Get("instance_id")
	request["VpcId"] = d.Get("vpc_id")
	request["VswitchId"] = d.Get("vswitch_id")
	request["ModuleName"] = d.Get("module_name")

	if v, ok := d.GetOkExists("enable_create_dns_record_in_pvzt"); ok {
		request["EnableCreateDNSRecordInPvzt"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_vpc_endpoint_linked_vpc", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["IsSuccess"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["InstanceId"], request["VpcId"], request["VswitchId"], request["ModuleName"]))

	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, crService.CrVpcEndpointLinkedVpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCrVpcEndpointLinkedVpcRead(d, meta)
}

func resourceAlicloudCrVpcEndpointLinkedVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	object, err := crService.DescribeCrVpcEndpointLinkedVpc(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VswitchId"])
	d.Set("module_name", parts[3])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudCrVpcEndpointLinkedVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudHbrHanaBackupClientRead(d, meta)
}

func resourceAlicloudCrVpcEndpointLinkedVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	action := "DeleteInstanceVpcEndpointLinkedVpc"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"VpcId":      parts[1],
		"VswitchId":  parts[2],
		"ModuleName": parts[3],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, true)
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"INSTANCE_ACCESS_NOT_EXIST"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, crService.CrVpcEndpointLinkedVpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
