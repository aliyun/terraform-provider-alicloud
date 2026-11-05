package ecs

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEcsHpcCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceEcsHpcClusterCreate,
		Read:   resourceEcsHpcClusterRead,
		Update: resourceEcsHpcClusterUpdate,
		Delete: resourceEcsHpcClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceEcsHpcClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHpcCluster"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["Name"] = d.Get("name")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateHpcCluster")
	response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_hpc_cluster", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["HpcClusterId"]))

	return resourceEcsHpcClusterRead(d, meta)
}
func resourceEcsHpcClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := Service{client}
	object, err := ecsService.DescribeEcsHpcCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_hpc_cluster ecsService.DescribeEcsHpcCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("name", object["Name"])
	return nil
}
func resourceEcsHpcClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"HpcClusterId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}
	if update {
		action := "ModifyHpcClusterAttribute"
		request["ClientToken"] = buildClientToken("ModifyHpcClusterAttribute")
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceEcsHpcClusterRead(d, meta)
}
func resourceEcsHpcClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHpcCluster"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"HpcClusterId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteHpcCluster")
	response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
