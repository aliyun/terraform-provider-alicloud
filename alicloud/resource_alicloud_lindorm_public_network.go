// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudLindormPublicNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudLindormPublicNetworkCreate,
		Read:   resourceAliCloudLindormPublicNetworkRead,
		Update: resourceAliCloudLindormPublicNetworkUpdate,
		Delete: resourceAliCloudLindormPublicNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enable_public_network": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1),
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudLindormPublicNetworkCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "SwitchInstancePublicNetwork"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["ActionType"] = d.Get("enable_public_network")
	request["EngineType"] = d.Get("engine_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Instance.IsNotValid"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_lindorm_public_network", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	lindormServiceV2 := LindormServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, lindormServiceV2.LindormPublicNetworkStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudLindormPublicNetworkRead(d, meta)
}

func resourceAliCloudLindormPublicNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	lindormServiceV2 := LindormServiceV2{client}

	objectRaw, err := lindormServiceV2.DescribeLindormPublicNetwork(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_lindorm_public_network DescribeLindormPublicNetwork Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("instance_id", objectRaw["InstanceId"])

	objectRaw, err = lindormServiceV2.DescribePublicNetworkGetLindormInstanceEngineList(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("instance_id", objectRaw["InstanceId"])

	d.Set("instance_id", d.Id())

	return nil
}

func resourceAliCloudLindormPublicNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Public Network.")
	return nil
}

func resourceAliCloudLindormPublicNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Public Network. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
