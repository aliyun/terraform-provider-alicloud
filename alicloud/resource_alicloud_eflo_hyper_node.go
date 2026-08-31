package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudEfloHyperNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEfloHyperNodeCreate,
		Read:   resourceAliCloudEfloHyperNodeRead,
		Update: resourceAliCloudEfloHyperNodeUpdate,
		Delete: resourceAliCloudEfloHyperNodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(38 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
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
						"bursting_enabled": {
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
						"provisioned_iops": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hpn_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"login_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Subscription"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_arch": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stage_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEfloHyperNodeCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("server_arch"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "ServerArch",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("hpn_zone"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "HpnZone",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("stage_num"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "StageNum",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "PaymentRatio",
		"Value": "0",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Classify",
		"Value": "gpuserver",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "discountlevel",
		"Value": d.Get("payment_duration"),
	})
	if v, ok := d.GetOk("machine_type"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "computingserver",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "ProductForm",
		"Value": "Hypernode",
	})
	if v, ok := d.GetOk("zone_id"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Zone",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["SubscriptionType"] = d.Get("payment_type")
	if v, ok := d.GetOkExists("renewal_duration"); ok {
		request["RenewPeriod"] = v
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOkExists("payment_duration"); ok {
		request["Period"] = v
	}
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_computinginstance_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_hyper_node", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	efloServiceV2 := EfloServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"HealthyUnused"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloServiceV2.EfloHyperNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEfloHyperNodeUpdate(d, meta)
}

func resourceAliCloudEfloHyperNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloServiceV2 := EfloServiceV2{client}

	objectRaw, err := efloServiceV2.DescribeEfloHyperNode(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_hyper_node DescribeEfloHyperNode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_id", objectRaw["ClusterId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("hostname", objectRaw["Hostname"])
	d.Set("hpn_zone", objectRaw["HpnZone"])
	d.Set("machine_type", objectRaw["MachineType"])
	d.Set("node_group_id", objectRaw["NodeGroupId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["OperatingState"])
	d.Set("zone_id", objectRaw["ZoneId"])

	objectRaw, err = efloServiceV2.DescribeHyperNodeQueryAvailableInstances(d)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payment_type", objectRaw["SubscriptionType"])
	d.Set("region_id", objectRaw["Region"])
	if fmt.Sprint(objectRaw["RenewalDurationUnit"]) == "Y" {
		d.Set("renewal_duration", formatInt(objectRaw["RenewalDuration"])*12)
	} else {
		d.Set("renewal_duration", objectRaw["RenewalDuration"])
	}
	d.Set("renewal_status", objectRaw["RenewStatus"])

	objectRaw, err = efloServiceV2.DescribeHyperNodeListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))
	if d.Get("status") == "HealthyUsing" {
		objectRaw, err = efloServiceV2.DescribeHyperNodeListClusterHyperNodes(fmt.Sprint(objectRaw["ClusterId"]), d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
		d.Set("vswitch_id", objectRaw["VSwitchId"])
		d.Set("vpc_id", objectRaw["VpcId"])
	}

	return nil
}

func resourceAliCloudEfloHyperNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "SetRenewal"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceIDs"] = d.Id()

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	request["SubscriptionType"] = d.Get("payment_type")
	if !d.IsNewResource() && d.HasChange("renewal_duration") {
		update = true
		request["RenewalPeriod"] = d.Get("renewal_duration")
	}

	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	request["RenewalStatus"] = d.Get("renewal_status")
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["SubscriptionType"] = v
	}
	if request["SubscriptionType"] == "" {
		request["SubscriptionType"] = "Subscription"
	}
	if request["SubscriptionType"] == "Subscription" && request["RenewalStatus"] == "AutoRenewal" {
		v, ok := d.GetOk("renewal_duration")
		if !ok {
			return WrapError(Error("renewal_duration is required when renewal_status is set to AutoRenewal."))
		}
		request["RenewalPeriod"] = v
		if v.(int) < 12 {
			request["RenewalPeriod"] = v
			request["RenewalPeriodUnit"] = "M"
		} else {
			if v.(int)%12 != 0 {
				return WrapError(Error("renewal_duration must be a multiple of 12 when renewal_duration more than 12."))
			}
			renewPeriod := v.(int) / 12
			request["RenewalPeriod"] = renewPeriod
			request["RenewalPeriodUnit"] = "Y"
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_eflocomputing_public_intl"
					if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
						request["ProductCode"] = "bccluster"
						request["ProductType"] = "bccluster_computinginstance_public_intl"
					}
					endpoint = connectivity.BssOpenAPIEndpointInternational
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
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "HyperNode"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		efloServiceV2 := EfloServiceV2{client}
		if err := efloServiceV2.SetResourceTags(d, "HyperNode"); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChanges("cluster_id", "node_group_id") {
		oldEntry, newEntry := d.GetChange("cluster_id")
		oldClusterId := oldEntry.(string)
		newClusterId := newEntry.(string)
		oldEntry, newEntry = d.GetChange("node_group_id")
		oldNodeGroupId := oldEntry.(string)
		newNodeGroupId := newEntry.(string)

		if oldNodeGroupId != "" {
			action := "ShrinkCluster"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionId"] = client.RegionId
			request["ClusterId"] = oldClusterId
			nodeGroupsDataList := make(map[string]interface{})
			nodeGroupsDataList["NodeGroupId"] = oldNodeGroupId

			NodeGroupsMap := make([]interface{}, 0)
			NodeGroupsMap = append(NodeGroupsMap, nodeGroupsDataList)
			request["NodeGroups"] = NodeGroupsMap

			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.HyperNodeId", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)
			request["NodeGroups"] = convertObjectToJsonString(request["NodeGroups"])

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			efloServiceV2 := EfloServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"HealthyUnused"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, efloServiceV2.EfloHyperNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if newNodeGroupId != "" {
			action := "ExtendCluster"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RegionId"] = client.RegionId
			request["ClusterId"] = newClusterId
			nodeGroupsDataList := make(map[string]interface{})
			nodeGroupsDataList["NodeGroupId"] = newNodeGroupId
			if v, ok := d.GetOk("user_data"); ok && v != "" {
				nodeGroupsDataList["UserData"] = v
			}

			NodeGroupsMap := make([]interface{}, 0)
			NodeGroupsMap = append(NodeGroupsMap, nodeGroupsDataList)
			request["NodeGroups"] = NodeGroupsMap

			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.HyperNodeId", d.Id())
			if v, ok := d.GetOk("hostname"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.Hostname", d.Id())
			}
			if v, ok := d.GetOk("login_password"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.LoginPassword", v)
			}
			if v, ok := d.GetOk("vpc_id"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.VpcId", v)
			}
			if v, ok := d.GetOk("vswitch_id"); ok && v != "" {
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.VSwitchId", v)
			}
			if v, ok := d.GetOk("data_disk"); ok {
				dataDisks := make([]map[string]interface{}, 0)
				for _, disk := range v.([]interface{}) {
					diskMap := disk.(map[string]interface{})
					dataDisk := make(map[string]interface{})

					if val, ok := diskMap["delete_with_node"]; ok {
						dataDisk["DeleteWithNode"] = val
					}
					if val, ok := diskMap["bursting_enabled"]; ok {
						dataDisk["BurstingEnabled"] = val
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
					if val, ok := diskMap["provisioned_iops"]; ok {
						dataDisk["ProvisionedIops"] = val
					}

					dataDisks = append(dataDisks, dataDisk)
				}
				jsonString, _ = sjson.Set(jsonString, "NodeGroups.0.HyperNodes.0.DataDisk", dataDisks)
			}
			_ = json.Unmarshal([]byte(jsonString), &request)
			request["NodeGroups"] = convertObjectToJsonString(request["NodeGroups"])
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			efloServiceV2 := EfloServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"HealthyUsing"}, d.Timeout(schema.TimeoutUpdate), 9*time.Minute, efloServiceV2.EfloHyperNodeStateRefreshFunc(d.Id(), "OperatingState", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}
	}
	d.Partial(false)
	return resourceAliCloudEfloHyperNodeRead(d, meta)
}

func resourceAliCloudEfloHyperNodeDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "RefundInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	request["ImmediatelyRelease"] = "1"
	var endpoint string
	request["ProductCode"] = "bccluster"
	request["ProductType"] = "bccluster_eflocomputing_public_cn"
	if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_computinginstance_public_cn"
	}
	if client.IsInternationalAccount() {
		request["ProductCode"] = "bccluster"
		request["ProductType"] = "bccluster_eflocomputing_public_intl"
		if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
			request["ProductCode"] = "bccluster"
			request["ProductType"] = "bccluster_computinginstance_public_intl"
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "bccluster"
				request["ProductType"] = "bccluster_eflocomputing_public_intl"
				if v, ok := d.GetOk("payment_type"); ok && v == "PayAsYouGo" {
					request["ProductCode"] = "bccluster"
					request["ProductType"] = "bccluster_computinginstance_public_intl"
				}
				endpoint = connectivity.BssOpenAPIEndpointInternational
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
