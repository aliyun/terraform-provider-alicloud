package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEfloNodeGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloNodeGroupAttachmentCreate,
		Read:   resourceAliCloudEfloNodeGroupAttachmentRead,
		Update: resourceAliCloudEfloNodeGroupAttachmentUpdate,
		Delete: resourceAliCloudEfloNodeGroupAttachmentDelete,
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

func resourceAliCloudEfloNodeGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

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
		return WrapError(err)
	}
	request["NodeGroups"] = string(objectDataLocalMapJson)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_node_group_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["ClusterId"], nodeGroup["NodeGroupId"], node["NodeId"]))

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"execution_success"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, efloServiceV2.DescribeAsyncEfloNodeGroupAttachmentStateRefreshFunc(d, response, "$.TaskState", []string{"execution_fail"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudEfloNodeGroupAttachmentRead(d, meta)
}

func resourceAliCloudEfloNodeGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloNodeGroupAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_node_group_attachment DescribeEfloNodeGroupAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudEfloNodeGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Node Group Attachment.")
	return nil
}

func resourceAliCloudEfloNodeGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

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
		return WrapError(err)
	}
	request["NodeGroups"] = string(objectDataLocalMapJson)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("eflo-controller", "2022-12-15", action, query, request, true)

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

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"execution_success"}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, efloServiceV2.DescribeAsyncEfloNodeGroupAttachmentStateRefreshFunc(d, response, "$.TaskState", []string{"execution_fail"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return nil
}
