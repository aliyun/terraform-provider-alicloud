package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcVbrHa() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcVbrHaCreate,
		Read:   resourceAlicloudVpcVbrHaRead,
		Update: resourceAlicloudVpcVbrHaUpdate,
		Delete: resourceAlicloudVpcVbrHaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"peer_vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vbr_ha_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudVpcVbrHaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVbrHa"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["PeerVbrId"] = d.Get("peer_vbr_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vbr_ha_name"); ok {
		request["Name"] = v
	}
	request["VbrId"] = d.Get("vbr_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateVbrHa")
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_vbr_ha", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VbrHaId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcVbrHaStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcVbrHaRead(d, meta)
}
func resourceAlicloudVpcVbrHaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcVbrHa(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_vbr_ha vpcService.DescribeVpcVbrHa Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("peer_vbr_id", object["PeerVbrId"])
	d.Set("status", object["Status"])
	d.Set("vbr_ha_name", object["Name"])
	d.Set("vbr_id", object["VbrId"])
	return nil
}
func resourceAlicloudVpcVbrHaUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudVpcVbrHaRead(d, meta)
}
func resourceAlicloudVpcVbrHaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteVbrHa"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteVbrHa")
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcVbrHaStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
