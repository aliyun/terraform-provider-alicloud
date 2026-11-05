package eflo

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNodeGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeGroupAttachmentCreate,
		Read:   resourceNodeGroupAttachmentRead,
		Update: resourceNodeGroupAttachmentUpdate,
		Delete: resourceNodeGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3605 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3605 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_with_node": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNodeGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ExtendCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	request["RegionId"] = client.RegionId

	nodeGroups := make([]map[string]interface{}, 1)
	nodeGroup := make(map[string]interface{})
	if v, ok := d.GetOk("node_group_id"); ok {
		nodeGroup["NodeGroupId"] = v
	}
	if v, ok := d.GetOk("user_data"); ok {
		nodeGroup["UserData"] = v
	}
	nodes := make([]map[string]interface{}, 1)
	node := make(map[string]interface{})
	if v, ok := d.GetOk("node_id"); ok {
		node["NodeId"] = v
	}
	if v, ok := d.GetOk("hostname"); ok {
		node["Hostname"] = v
	}
	if v, ok := d.GetOk("login_password"); ok {
		node["LoginPassword"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		node["VpcId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		node["VSwitchId"] = v
	}
	if v, ok := d.GetOk("data_disk"); ok {
		dataDisks := make([]map[string]interface{}, 0)
		for _, disk := range v.([]interface{}) {
			diskMap := disk.(map[string]interface{})
			dataDisk := make(map[string]interface{})

			if val, ok := diskMap["delete_with_node"]; ok {
				dataDisk["DeleteWithNode"] = val
			}
			if val, ok := diskMap["category"]; ok {
				dataDisk["Category"] = val
			}
			if val, ok := diskMap["size"]; ok {
				dataDisk["Size"] = val
			}
			if val, ok := diskMap["performance_level"]; ok {
				dataDisk["PerformanceLevel"] = val
			}

			dataDisks = append(dataDisks, dataDisk)
		}
		node["DataDisk"] = dataDisks
	}

	nodes[0] = node
	nodeGroup["Nodes"] = nodes
	nodeGroups[0] = nodeGroup
	objectDataLocalMapJson, err := json.Marshal(nodeGroups)
	if err != nil {
		return errs.WrapError(err)
	}
	request["NodeGroups"] = string(objectDataLocalMapJson)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_eflo_node_group_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["ClusterId"], nodeGroup["NodeGroupId"], node["NodeId"]))

	efloServiceV2 := EfloServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"execution_success"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, efloServiceV2.DescribeAsyncEfloNodeGroupAttachmentStateRefreshFunc(d, response, "$.TaskState", []string{"execution_fail"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return resourceNodeGroupAttachmentRead(d, meta)
}

func resourceNodeGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloNodeGroupAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_node_group_attachment DescribeEfloNodeGroupAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("hostname", objectRaw["Hostname"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("node_group_id", objectRaw["NodeGroupId"])
	d.Set("node_id", objectRaw["NodeId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("cluster_id", parts[0])

	return nil
}

func resourceNodeGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Node Group Attachment.")
	return nil
}

func resourceNodeGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "ShrinkCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = parts[0]
	request["RegionId"] = client.RegionId

	nodeGroups := make([]map[string]interface{}, 1)
	nodeGroup := make(map[string]interface{})
	if v, ok := d.GetOk("node_group_id"); ok {
		nodeGroup["NodeGroupId"] = v
	}
	nodes := make([]map[string]interface{}, 1)
	node := make(map[string]interface{})
	if v, ok := d.GetOk("node_id"); ok {
		node["NodeId"] = v
	}
	nodes[0] = node
	nodeGroup["Nodes"] = nodes
	nodeGroups[0] = nodeGroup
	objectDataLocalMapJson, err := json.Marshal(nodeGroups)
	if err != nil {
		return errs.WrapError(err)
	}
	request["NodeGroups"] = string(objectDataLocalMapJson)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

	efloServiceV2 := EfloServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"execution_success"}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, efloServiceV2.DescribeAsyncEfloNodeGroupAttachmentStateRefreshFunc(d, response, "$.TaskState", []string{"execution_fail"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return nil
}
