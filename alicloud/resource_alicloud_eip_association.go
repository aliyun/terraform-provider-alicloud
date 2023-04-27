package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipAssociationCreate,
		Read:   resourceAliyunEipAssociationRead,
		Update: resourceAliyunEipAssociationUpdate,
		Delete: resourceAliyunEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AssociateEipAddress"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	vpcService := VpcService{client}
	if v, ok := d.GetOk("activity_id"); ok {
		request["ActivityId"] = v
	}
	request["RegionId"] = client.RegionId
	request["AllocationId"] = Trim(d.Get("allocation_id").(string))
	request["InstanceId"] = Trim(d.Get("instance_id").(string))
	request["InstanceType"] = EcsInstance
	if strings.HasPrefix(request["InstanceId"].(string), "lb-") {
		request["InstanceType"] = SlbInstance
	}
	if strings.HasPrefix(request["InstanceId"].(string), "ngw-") {
		request["InstanceType"] = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = instanceType.(string)
	}
	if privateIPAddress, ok := d.GetOk("private_ip_address"); ok {
		request["PrivateIpAddress"] = privateIPAddress.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v.(string)
	}
	request["ClientToken"] = buildClientToken("AssociateEipAddress")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectEipStatus", "InvalidBindingStatus", "IncorrectInstanceStatus", "IncorrectStatus.NatGateway", "InvalidStatus.EcsStatusNotSupport", "InvalidStatus.InstanceHasBandWidth", "InvalidStatus.EniStatusNotSupport"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eip_association", action, AlibabaCloudSdkGoERROR)
	}

	if err := vpcService.WaitForEip(request["AllocationId"].(string), InUse, 60); err != nil {
		return WrapError(err)
	}

	d.SetId(fmt.Sprint(request["AllocationId"], ":", request["InstanceId"]))

	return resourceAliyunEipAssociationRead(d, meta)
}

func resourceAliyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeEipAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", object["InstanceId"])
	d.Set("allocation_id", object["AllocationId"])
	d.Set("instance_type", object["InstanceType"])
	d.Set("vpc_id", object["VpcId"])
	return nil
}

func resourceAliyunEipAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] The update method is used to ensure that the force parameter does not need to add forcenew.")
	return nil
}

func resourceAliyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	allocationId, instanceId := parts[0], parts[1]
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateEipAddressRequest()
	request.RegionId = client.RegionId
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	request.InstanceType = EcsInstance
	request.Force = requests.NewBoolean(d.Get("force").(bool))
	request.ClientToken = buildClientToken(request.GetActionName())

	if strings.HasPrefix(instanceId, "lb-") {
		request.InstanceType = SlbInstance
	}
	if strings.HasPrefix(instanceId, "ngw-") {
		request.InstanceType = Nat
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = instanceType.(string)
	}
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateEipAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus.%s", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectEipStatus", "InvalidBindingStatus", "IncorrectInstanceStatus", "IncorrectHaVipStatus", "TaskConflict", "InvalidIpStatus.HasBeenUsedBySnatTable", "InvalidIpStatus.HasBeenUsedByForwardEntry", "InvalidStatus.EniStatusNotSupport", "InvalidStatus.EcsStatusNotSupport", "InvalidStatus.NotAllow", "InvalidStatus.SnatOrDnat"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.EipAssociationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
