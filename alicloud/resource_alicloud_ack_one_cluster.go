// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAckOneCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAckOneClusterCreate,
		Read:   resourceAliCloudAckOneClusterRead,
		Update: resourceAliCloudAckOneClusterUpdate,
		Delete: resourceAliCloudAckOneClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
			Delete: schema.DefaultTimeout(25 * time.Minute),
			Update: schema.DefaultTimeout(25 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"argocd_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitches": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"profile": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Default", "XFlow"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAckOneClusterCreate(d *schema.ResourceData, meta interface{}) error {

	if profile, ok := d.GetOk("profile"); ok && profile.(string) == "XFlow" {
		if _, ok := d.GetOkExists("argocd_enabled"); ok {
			return fmt.Errorf("argocd_enabled is not supported when profile is XFlow")
		}
	}

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHubCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	jsonPathResult, err := jsonpath.Get("$[0].vpc_id", d.Get("network"))
	if err == nil {
		request["VpcId"] = jsonPathResult
	}

	jsonPathResult1, err := jsonpath.Get("$[0].vswitches", d.Get("network"))
	if err == nil {
		request["VSwitches"] = convertListToJsonString(jsonPathResult1.([]interface{}))
	}

	profile := "Default"
	if v, ok := d.GetOk("profile"); ok {
		request["Profile"] = v
		profile = v.(string)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		request["Name"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ack_one_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClusterId"]))

	ackOneServiceV2 := AckOneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackOneServiceV2.AckOneClusterStateRefreshFunc(d.Id(), "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// Wait until argocd_enabled attribute is available
	if profile != "XFlow" {
		stateConf = BuildStateConf([]string{"<nil>"}, []string{"true", "false"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackOneServiceV2.AckOneClusterStateRefreshFunc(d.Id(), "$.ClusterInfo.MetaData.ACKOne.GitOps.Enabled", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, "Failed to get argocd_enabled attribute after cluster creation")
		}

		if v, ok := d.GetOkExists("argocd_enabled"); ok && !v.(bool) {
			if err := resourceAliCloudAckOneClusterUpdate(d, meta); err != nil {
				return WrapErrorf(err, "Failed to set argocd_enabled correctly when creating cluster")
			}
		}
	}

	return resourceAliCloudAckOneClusterRead(d, meta)
}

func resourceAliCloudAckOneClusterRead(d *schema.ResourceData, meta interface{}) error {

	if profile, ok := d.GetOk("profile"); ok && profile.(string) == "XFlow" {
		if _, ok := d.GetOkExists("argocd_enabled"); ok {
			return fmt.Errorf("argocd_enabled is not supported when profile is XFlow")
		}
	}

	client := meta.(*connectivity.AliyunClient)
	ackOneServiceV2 := AckOneServiceV2{client}

	objectRaw, err := ackOneServiceV2.DescribeAckOneCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ack_one_cluster DescribeAckOneCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	clusterInfo1RawObj, _ := jsonpath.Get("$.ClusterInfo", objectRaw)
	clusterInfo1Raw := make(map[string]interface{})
	if clusterInfo1RawObj != nil {
		clusterInfo1Raw = clusterInfo1RawObj.(map[string]interface{})
	}
	d.Set("cluster_name", clusterInfo1Raw["Name"])
	d.Set("create_time", clusterInfo1Raw["CreationTime"])
	d.Set("profile", clusterInfo1Raw["Profile"])
	d.Set("status", clusterInfo1Raw["State"])

	networkMaps := make([]map[string]interface{}, 0)
	networkMap := make(map[string]interface{})
	network1Raw := make(map[string]interface{})
	if objectRaw["Network"] != nil {
		network1Raw = objectRaw["Network"].(map[string]interface{})
	}
	if len(network1Raw) > 0 {
		networkMap["vpc_id"] = network1Raw["VpcId"]
		securityGroupIDs1Raw := make([]interface{}, 0)
		if network1Raw["SecurityGroupIDs"] != nil {
			securityGroupIDs1Raw = network1Raw["SecurityGroupIDs"].([]interface{})
		}

		networkMap["security_group_ids"] = securityGroupIDs1Raw
		vSwitches1Raw := make([]interface{}, 0)
		if network1Raw["VSwitches"] != nil {
			vSwitches1Raw = network1Raw["VSwitches"].([]interface{})
		}

		networkMap["vswitches"] = vSwitches1Raw
		networkMaps = append(networkMaps, networkMap)
	}
	d.Set("network", networkMaps)

	if d.Get("profile") == "Default" {
		argocdEnabled, err := jsonpath.Get("$.ClusterInfo.MetaData.ACKOne.GitOps.Enabled", objectRaw)
		if err != nil {
			d.Set("argocd_enabled", false)
		} else {
			if v, ok := argocdEnabled.(bool); ok {
				d.Set("argocd_enabled", v)
			} else {
				d.Set("argocd_enabled", false)
			}
		}
	}
	return nil
}

func resourceAliCloudAckOneClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	if profile, ok := d.GetOk("profile"); ok && profile.(string) == "XFlow" {
		if _, ok := d.GetOkExists("argocd_enabled"); ok {
			return fmt.Errorf("argocd_enabled is not supported when profile is XFlow")
		}
	}

	client := meta.(*connectivity.AliyunClient)
	action := "UpdateHubClusterFeature"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = d.Id()
	if d.Get("profile") != "XFlow" {
		if v, ok := d.GetOkExists("argocd_enabled"); ok {
			request["ArgoCDEnabled"] = v.(bool)
		} else {
			request["ArgoCDEnabled"] = true
		}
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ack_one_cluster", action, AlibabaCloudSdkGoERROR)
	}

	ackOneServiceV2 := AckOneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ackOneServiceV2.AckOneClusterStateRefreshFunc(d.Id(), "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	time.Sleep(4 * time.Minute)

	return resourceAliCloudAckOneClusterRead(d, meta)
}

func resourceAliCloudAckOneClusterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteHubCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

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

	ackOneServiceV2 := AckOneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ackOneServiceV2.AckOneClusterStateRefreshFunc(d.Id(), "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
