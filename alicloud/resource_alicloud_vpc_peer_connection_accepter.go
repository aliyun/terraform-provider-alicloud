package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcPeerConnectionAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcPeerConnectionAccepterCreate,
		Read:   resourceAlicloudVpcPeerConnectionAccepterRead,
		Update: resourceAlicloudVpcPeerConnectionAccepterUpdate,
		Delete: resourceAlicloudVpcPeerConnectionAccepterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accepting_owner_uid": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"accepting_region_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"accepting_vpc_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"bandwidth": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"dry_run": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peer_connection_accepter_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vpc_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudVpcPeerConnectionAccepterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcPeerService := VpcPeerService{client}
	request := make(map[string]interface{})
	conn, err := client.NewVpcPeerClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	request["ClientToken"] = buildClientToken("AcceptVpcPeerConnection")
	var response map[string]interface{}
	action := "AcceptVpcPeerConnection"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_peer_connection_accepter", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	stateConf := BuildStateConf([]string{}, []string{"Activated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcPeerService.VpcPeerConnectionAccepterStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudVpcPeerConnectionAccepterRead(d, meta)
}

func resourceAlicloudVpcPeerConnectionAccepterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcPeerService := VpcPeerService{client}

	object, err := vpcPeerService.DescribeVpcPeerConnectionAccepter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_peer_connection_accepter vpcPeerService.DescribeVpcPeerConnectionAccepter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("accepting_owner_uid", formatInt(object["AcceptingOwnerUid"]))
	d.Set("accepting_region_id", object["AcceptingRegionId"])
	d.Set("accepting_vpc_id", object["AcceptingVpc"].(map[string]interface{})["VpcId"])
	if v, ok := object["Bandwidth"]; ok && fmt.Sprint(v) != "0" {
		d.Set("bandwidth", formatInt(v))
	}
	d.Set("description", object["Description"])
	d.Set("peer_connection_accepter_name", object["Name"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["Vpc"].(map[string]interface{})["VpcId"])
	d.Set("instance_id", object["InstanceId"])
	return nil
}

func resourceAlicloudVpcPeerConnectionAccepterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudVpcPeerConnectionAccepterRead(d, meta)
}

func resourceAlicloudVpcPeerConnectionAccepterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcPeerService := VpcPeerService{client}
	conn, err := client.NewVpcPeerClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteVpcPeerConnection")
	action := "DeleteVpcPeerConnection"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.RouteEntryExist"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcPeerService.VpcPeerConnectionAccepterStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
