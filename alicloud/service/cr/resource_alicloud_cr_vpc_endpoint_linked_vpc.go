package cr

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cr_vpc_endpoint_linked_vpc", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["IsSuccess"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["InstanceId"], request["VpcId"], request["VswitchId"], request["ModuleName"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, crService.CrVpcEndpointLinkedVpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAlicloudCrVpcEndpointLinkedVpcRead(d, meta)
}

func resourceAlicloudCrVpcEndpointLinkedVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	object, err := crService.DescribeCrVpcEndpointLinkedVpc(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
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
	// The baseline realized Update by borrowing resourceAlicloudHbrHanaBackupClientRead
	// (an HBR read outside package cr's visibility): a copy-paste defect that called
	// HBR's DescribeHbrHanaBackupClient with this resource's CR id and could drop it
	// from state on NotFound. Migration into package cr made the HBR read unreachable.
	// Per the provider's standard no-update pattern (an Update delegates to its own
	// Read), delegate to this resource's own Read to preserve Read-based refresh.
	// Authorized production fix scoped to this Update only; no other behavior changed.
	return resourceAlicloudCrVpcEndpointLinkedVpcRead(d, meta)
}

func resourceAlicloudCrVpcEndpointLinkedVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	action := "DeleteInstanceVpcEndpointLinkedVpc"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"VpcId":      parts[1],
		"VswitchId":  parts[2],
		"ModuleName": parts[3],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, true)
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
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"INSTANCE_ACCESS_NOT_EXIST"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, crService.CrVpcEndpointLinkedVpcStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
