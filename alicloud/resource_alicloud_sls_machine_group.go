// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsMachineGroupCreate,
		Read:   resourceAliCloudSlsMachineGroupRead,
		Delete: resourceAliCloudSlsMachineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_attribute": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_topic": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"external_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"machine_identify_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudSlsMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/machinegroups")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("group_name"); ok {
		request["groupName"] = v
	}

	request["machineIdentifyType"] = d.Get("machine_identify_type")
	if v, ok := d.GetOk("group_type"); ok {
		request["groupType"] = v
	}
	dataList := make(map[string]interface{})

	if v := d.Get("group_attribute"); !IsNil(v) {
		groupTopic1, _ := jsonpath.Get("$[0].group_topic", v)
		if groupTopic1 != nil && groupTopic1 != "" {
			dataList["groupTopic"] = groupTopic1
		}
		externalName1, _ := jsonpath.Get("$[0].external_name", v)
		if externalName1 != nil && externalName1 != "" {
			dataList["externalName"] = externalName1
		}

		request["groupAttribute"] = dataList
	}

	if v, ok := d.GetOk("machine_list"); ok {
		machineListMapsArray := v.([]interface{})
		request["machineList"] = machineListMapsArray
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateMachineGroup", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_machine_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["groupName"]))

	return resourceAliCloudSlsMachineGroupRead(d, meta)
}

func resourceAliCloudSlsMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsMachineGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_machine_group DescribeSlsMachineGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_type", objectRaw["groupType"])
	d.Set("machine_identify_type", objectRaw["machineIdentifyType"])

	groupAttributeMaps := make([]map[string]interface{}, 0)
	groupAttributeMap := make(map[string]interface{})
	groupAttributeRaw := make(map[string]interface{})
	if objectRaw["groupAttribute"] != nil {
		groupAttributeRaw = objectRaw["groupAttribute"].(map[string]interface{})
	}
	if len(groupAttributeRaw) > 0 {
		groupAttributeMap["external_name"] = groupAttributeRaw["externalName"]
		groupAttributeMap["group_topic"] = groupAttributeRaw["groupTopic"]

		groupAttributeMaps = append(groupAttributeMaps, groupAttributeMap)
	}
	if err := d.Set("group_attribute", groupAttributeMaps); err != nil {
		return err
	}
	machineListRaw := make([]interface{}, 0)
	if objectRaw["machineList"] != nil {
		machineListRaw = objectRaw["machineList"].([]interface{})
	}

	d.Set("machine_list", machineListRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])
	d.Set("group_name", parts[1])

	return nil
}

func resourceAliCloudSlsMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	machineGroup := parts[1]
	action := fmt.Sprintf("/machinegroups/%s", machineGroup)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteMachineGroup", action), query, nil, nil, hostMap, false)

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

	return nil
}
