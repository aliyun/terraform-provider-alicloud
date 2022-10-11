package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMseCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMseClusterCreate,
		Read:   resourceAlicloudMseClusterRead,
		Update: resourceAlicloudMseClusterUpdate,
		Delete: resourceAlicloudMseClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_entry_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("pub_network_flow").(string) == "0" && d.Get("net_type").(string) == "privatenet" {
						return true
					}
					return false
				},
			},
			"cluster_alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_specification": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MSE_SC_1_2_60_c", "MSE_SC_2_4_60_c", "MSE_SC_4_8_60_c", "MSE_SC_8_16_60_c", "MSE_SC_16_32_60_c", "MSE_SC_1_2_200_c", "MSE_SC_2_4_200_c", "MSE_SC_4_8_200_c", "MSE_SC_8_16_200_c"}, false),
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Eureka", "Nacos-Ans", "ZooKeeper"}, false),
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb"}, false),
			},
			"disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"privatenet", "pubnet"}, false),
			},
			"private_slb_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pub_network_flow": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pub_slb_specification": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"request_pars": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mse_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"mse_dev", "mse_basic", "mse_pro"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudMseClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	var response map[string]interface{}
	action := "CreateCluster"
	request := make(map[string]interface{})
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	request["ClusterSpecification"] = d.Get("cluster_specification")
	request["ClusterType"] = d.Get("cluster_type")
	request["ClusterVersion"] = d.Get("cluster_version")
	if v, ok := d.GetOk("connection_type"); ok {
		request["ConnectionType"] = v
	}
	if v, ok := d.GetOk("disk_type"); ok {
		request["DiskType"] = v
	}

	request["InstanceCount"] = d.Get("instance_count")
	request["NetType"] = d.Get("net_type")
	if v, ok := d.GetOk("private_slb_specification"); ok {
		request["PrivateSlbSpecification"] = v
	}

	request["PubNetworkFlow"] = d.Get("pub_network_flow")

	if v, ok := d.GetOk("pub_slb_specification"); ok {
		request["PubSlbSpecification"] = v
	}
	if v, ok := d.GetOk("mse_version"); ok {
		request["MseVersion"] = v
	}
	if v, ok := d.GetOk("request_pars"); ok {
		request["RequestPars"] = v
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}

	request["Region"] = client.RegionId

	if request["VpcId"] == nil && request["VSwitchId"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVswitch(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := request["VpcId"].(string); !ok || v == "" {
			request["VpcId"] = vsw["VpcId"]
		}
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mse_cluster", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	d.SetId(fmt.Sprint(response["InstanceId"]))
	stateConf := BuildStateConf([]string{}, []string{"INIT_SUCCESS"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, mseService.MseClusterStateRefreshFunc(d.Id(), []string{"INIT_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudMseClusterUpdate(d, meta)
}

func resourceAlicloudMseClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	object, err := mseService.DescribeMseCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_cluster mseService.DescribeMseCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("cluster_type", object["ClusterType"])
	d.Set("cluster_specification", object["ClusterSpecification"])
	d.Set("instance_count", formatInt(object["InstanceCount"]))
	d.Set("pub_network_flow", object["PubNetworkFlow"])
	d.Set("status", object["InitStatus"])
	d.Set("cluster_id", object["ClusterId"])
	d.Set("mse_version", object["MseVersion"])
	d.Set("net_type", object["NetType"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("cluster_version", object["OrderClusterVersion"])
	d.Set("cluster_alias_name", object["ClusterAliasName"])
	d.Set("connection_type", object["ConnectionType"])
	d.Set("vpc_id", object["VpcId"])
	return nil
}

func resourceAlicloudMseClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	var response map[string]interface{}
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("acl_entry_list") {
		request := map[string]interface{}{
			"InstanceId": d.Id(),
		}
		request["AclEntryList"] = convertListToCommaSeparate(d.Get("acl_entry_list").(*schema.Set).List())
		action := "UpdateAcl"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("acl_entry_list")
	}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("cluster_alias_name") {
		update = true
		request["ClusterAliasName"] = d.Get("cluster_alias_name")
	}
	if update {
		action := "UpdateCluster"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cluster_alias_name")
	}

	update = false
	object, err := mseService.DescribeMseCluster(d.Id())
	updateClusterSpecReq := map[string]interface{}{
		"InstanceId": d.Id(),
		"ClusterId":  object["ClusterId"],
	}

	if !d.IsNewResource() && d.HasChange("cluster_specification") {
		update = true
	}
	if v, ok := d.GetOk("cluster_specification"); ok {
		updateClusterSpecReq["ClusterSpecification"] = v
	}

	if !d.IsNewResource() && d.HasChange("instance_count") {
		update = true
	}
	if v, ok := d.GetOk("instance_count"); ok {
		updateClusterSpecReq["InstanceCount"] = v
	}

	if update {
		action := "UpdateClusterSpec"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, updateClusterSpecReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, updateClusterSpecReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"SCALE_SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, mseService.MseClusterStateRefreshFunc(d.Id(), []string{"INIT_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("cluster_specification")
		d.SetPartial("instance_count")
	}

	d.Partial(false)
	return resourceAlicloudMseClusterRead(d, meta)
}

func resourceAlicloudMseClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := MseService{client}
	action := "DeleteCluster"
	var response map[string]interface{}
	conn, err := client.NewMseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapErrorf(fmt.Errorf("%s failed, response: %v", action, response), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"DESTROY_SUCCESS"}, d.Timeout(schema.TimeoutDelete), 60*time.Second, mseService.MseClusterStateRefreshFunc(d.Id(), []string{"DESTROY_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
