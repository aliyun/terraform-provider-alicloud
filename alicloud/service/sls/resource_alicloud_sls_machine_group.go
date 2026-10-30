package sls

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeSet,
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
	hostMap["project"] = helper.StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("group_name"); ok {
		request["groupName"] = v
	}

	request["machineIdentifyType"] = d.Get("machine_identify_type")
	if v, ok := d.GetOk("group_type"); ok {
		request["groupType"] = v
	}
	dataList := make(map[string]interface{})

	if v := d.Get("group_attribute"); !helper.IsNil(v) {
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
		machineListMapsArray := v.(*schema.Set).List()
		request["machineList"] = machineListMapsArray
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Sls", helper.RoaParam("POST", "2020-12-30", "CreateMachineGroup", action), query, body, nil, hostMap, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sls_machine_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["groupName"]))

	return resourceAliCloudSlsMachineGroupRead(d, meta)
}

func resourceAliCloudSlsMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := Service{client}

	objectRaw, err := slsServiceV2.DescribeSlsMachineGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_machine_group DescribeSlsMachineGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	hostMap["project"] = helper.StringPointer(parts[0])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Sls", helper.RoaParam("DELETE", "2020-12-30", "DeleteMachineGroup", action), query, nil, nil, hostMap, false)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
