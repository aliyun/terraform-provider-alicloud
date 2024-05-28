package alicloud

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAlikafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlikafkaInstanceCreate,
		Read:   resourceAliCloudAlikafkaInstanceRead,
		Update: resourceAliCloudAlikafkaInstanceUpdate,
		Delete: resourceAliCloudAlikafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"deploy_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 5}),
			},
			"partition_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"partition_num", "topic_quota"},
			},
			"topic_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, _ := strconv.Atoi(old)
					partitionNum := d.Get("partition_num").(int)
					if o > 0 {
						return o-1000 == partitionNum
					}
					return false
				},
				Deprecated: "Attribute `topic_quota` has been deprecated since 1.194.0 and it will be removed in the next future. Using new attribute `partition_num` instead.",
			},
			"io_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"io_max_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"io_max", "io_max_spec"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"paid_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "normal",
			},
			"eip_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(int) == 5
				},
			},
			"security_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: alikafkaInstanceConfigDiffSuppressFunc,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"selected_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"allowed_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port_range": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"9092/9092", "9094/9094", "9095/9095"}, false),
									},
									"allowed_ip_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"internet_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port_range": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: StringInSlice([]string{"9093/9093"}, false),
									},
									"allowed_ip_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"tags": tagsSchema(),
			"end_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_num_of_buy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"topic_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"topic_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_partition_buy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	vpcService := VpcService{client}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}

	// 1. Create order
	var createOrderAction string
	createOrderResponse := make(map[string]interface{})
	createOrderReq := make(map[string]interface{})
	createOrderReq["RegionId"] = client.RegionId
	if v, ok := d.GetOk("partition_num"); ok {
		createOrderReq["PartitionNum"] = v
	} else if v, ok := d.GetOk("topic_quota"); ok {
		createOrderReq["TopicQuota"] = v
	}
	createOrderReq["DiskType"] = d.Get("disk_type")
	createOrderReq["DiskSize"] = d.Get("disk_size")
	createOrderReq["DeployType"] = d.Get("deploy_type")
	if v, ok := d.GetOk("io_max"); ok {
		createOrderReq["IoMax"] = v
	}
	if v, ok := d.GetOk("io_max_spec"); ok {
		createOrderReq["IoMaxSpec"] = v
	}
	if v, ok := d.GetOk("spec_type"); ok {
		createOrderReq["SpecType"] = v
	}
	if v, ok := d.GetOkExists("eip_max"); ok {
		createOrderReq["EipMax"] = v
	}
	if v, ok := d.GetOk("paid_type"); ok {
		switch v.(string) {
		case "PostPaid":
			createOrderAction = "CreatePostPayOrder"
		case "PrePaid":
			createOrderAction = "CreatePrePayOrder"
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		createOrderResponse, err = conn.DoRequest(StringPointer(createOrderAction), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, createOrderReq, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL", "ONS_SYSTEM_ERROR"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(createOrderAction, createOrderResponse, createOrderReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", createOrderAction, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(createOrderResponse["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", createOrderAction, createOrderResponse))
	}

	alikafkaInstanceVO, err := alikafkaService.DescribeAliKafkaInstanceByOrderId(fmt.Sprint(createOrderResponse["OrderId"]), 60)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(alikafkaInstanceVO["InstanceId"]))

	// 2. Start instance
	startInstanceAction := "StartInstance"
	startInstanceResponse := make(map[string]interface{})
	startInstanceReq := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		startInstanceReq["VpcId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		startInstanceReq["ZoneId"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		startInstanceReq["VSwitchId"] = v
	}

	if (startInstanceReq["ZoneId"] == nil || startInstanceReq["VpcId"] == nil) && startInstanceReq["VSwitchId"] != nil {
		vsw, err := vpcService.DescribeVswitch(startInstanceReq["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := startInstanceReq["VpcId"].(string); !ok || v == "" {
			startInstanceReq["VpcId"] = vsw["VpcId"]
		}
		if v, ok := startInstanceReq["ZoneId"].(string); !ok || v == "" {
			startInstanceReq["ZoneId"] = vsw["ZoneId"]
		}
	}

	startInstanceReq["RegionId"] = client.RegionId
	startInstanceReq["InstanceId"] = alikafkaInstanceVO["InstanceId"]
	if _, ok := d.GetOkExists("eip_max"); ok {
		startInstanceReq["DeployModule"] = "eip"
		startInstanceReq["IsEipInner"] = true
	}
	if v, ok := d.GetOk("name"); ok {
		startInstanceReq["Name"] = v
	}
	if v, ok := d.GetOk("security_group"); ok {
		startInstanceReq["SecurityGroup"] = v
	}
	if v, ok := d.GetOk("service_version"); ok {
		startInstanceReq["ServiceVersion"] = v
	}
	if v, ok := d.GetOk("config"); ok {
		startInstanceReq["Config"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		startInstanceReq["KMSKeyId"] = v
	}
	if v, ok := d.GetOk("selected_zones"); ok {
		startInstanceReq["SelectedZones"] = formatSelectedZonesReq(v.([]interface{}))
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		startInstanceResponse, err = conn.DoRequest(StringPointer(startInstanceAction), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, startInstanceReq, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(startInstanceAction, startInstanceResponse, startInstanceReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alikafka_instance", startInstanceAction, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(startInstanceResponse["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", startInstanceAction, startInstanceResponse))
	}

	// 3. wait until running
	stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAlikafkaInstanceUpdate(d, meta)
}

func resourceAliCloudAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	object, err := alikafkaService.DescribeAliKafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("disk_type", object["DiskType"])
	d.Set("disk_size", object["DiskSize"])
	d.Set("deploy_type", object["DeployType"])
	d.Set("io_max", object["IoMax"])
	d.Set("eip_max", object["EipMax"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("paid_type", PostPaid)
	d.Set("spec_type", object["SpecType"])
	d.Set("security_group", object["SecurityGroup"])
	d.Set("end_point", object["EndPoint"])
	d.Set("status", object["ServiceStatus"])
	// object.UpgradeServiceDetailInfo.UpgradeServiceDetailInfoVO[0].Current2OpenSourceVersion can guaranteed not to be null
	d.Set("service_version", object["UpgradeServiceDetailInfo"].(map[string]interface{})["Current2OpenSourceVersion"])
	d.Set("config", object["AllConfig"])
	d.Set("kms_key_id", object["KmsKeyId"])
	if fmt.Sprint(object["PaidType"]) == "0" {
		d.Set("paid_type", PrePaid)
	}

	allowedIpList, err := alikafkaService.GetAllowedIpList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if allowedList, ok := allowedIpList["AllowedList"]; ok {
		allowedListMaps := make([]map[string]interface{}, 0)
		allowedListMap := map[string]interface{}{}

		if vpcList, ok := allowedList.(map[string]interface{})["VpcList"]; ok {
			vpcListMaps := make([]map[string]interface{}, 0)
			for _, vpcListValue := range vpcList.([]interface{}) {
				vpcListArg := vpcListValue.(map[string]interface{})
				vpcListMap := map[string]interface{}{}

				if portRange, ok := vpcListArg["PortRange"]; ok {
					vpcListMap["port_range"] = portRange
				}

				if vpcAllowedIpList, ok := vpcListArg["AllowedIpList"]; ok {
					vpcListMap["allowed_ip_list"] = vpcAllowedIpList
				}

				vpcListMaps = append(vpcListMaps, vpcListMap)
			}

			allowedListMap["vpc_list"] = vpcListMaps
		}

		if internetList, ok := allowedList.(map[string]interface{})["InternetList"]; ok {
			internetListMaps := make([]map[string]interface{}, 0)
			for _, internetListValue := range internetList.([]interface{}) {
				internetListArg := internetListValue.(map[string]interface{})
				internetListMap := map[string]interface{}{}

				if portRange, ok := internetListArg["PortRange"]; ok {
					internetListMap["port_range"] = portRange
				}

				if internetAllowedIpList, ok := internetListArg["AllowedIpList"]; ok {
					internetListMap["allowed_ip_list"] = internetAllowedIpList
				}

				internetListMaps = append(internetListMaps, internetListMap)
			}

			allowedListMap["internet_list"] = internetListMaps
		}

		allowedListMaps = append(allowedListMaps, allowedListMap)

		d.Set("allowed_list", allowedListMaps)
	}

	quota, err := alikafkaService.GetQuotaTip(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("topic_quota", quota["TopicQuota"])
	d.Set("partition_num", quota["PartitionNumOfBuy"])
	d.Set("topic_num_of_buy", quota["TopicNumOfBuy"])
	d.Set("topic_used", quota["TopicUsed"])
	d.Set("topic_left", quota["TopicLeft"])
	d.Set("partition_used", quota["PartitionUsed"])
	d.Set("partition_left", quota["PartitionLeft"])
	d.Set("group_used", quota["GroupUsed"])
	d.Set("group_left", quota["GroupLeft"])
	d.Set("is_partition_buy", quota["IsPartitionBuy"])

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", alikafkaService.tagsToMap(tags))
	d.Set("io_max_spec", object["IoMaxSpec"])

	return nil
}

func resourceAliCloudAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if err := alikafkaService.setInstanceTags(d, TagResourceInstance); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliCloudAlikafkaInstanceRead(d, meta)
	}

	// Process change instance name.
	if d.HasChange("name") {
		action := "ModifyInstanceName"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("name"); ok {
			request["InstanceName"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("name")
	}

	// Process paid type change, note only support change from post to pre pay.
	if d.HasChange("paid_type") {
		o, n := d.GetChange("paid_type")
		oldPaidType := o.(string)
		newPaidType := n.(string)
		oldPaidTypeInt := 1
		newPaidTypeInt := 1
		if oldPaidType == string(PrePaid) {
			oldPaidTypeInt = 0
		}
		if newPaidType == string(PrePaid) {
			newPaidTypeInt = 0
		}
		if oldPaidTypeInt == 1 && newPaidTypeInt == 0 {
			action := "ConvertPostPayOrder"
			request := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InstanceId": d.Id(),
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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

			stateConf := BuildStateConf([]string{}, []string{strconv.Itoa(newPaidTypeInt)}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "PaidType", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		} else {
			return WrapError(errors.New("paid type only support change from post pay to pre pay"))
		}

		d.SetPartial("paid_type")
	}

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}
	// updating topic_quota only by updating partition_num
	if d.HasChange("partition_num") {
		update = true
	}
	request["PartitionNum"] = d.Get("partition_num")
	if d.HasChange("disk_size") {
		update = true
	}
	request["DiskSize"] = d.Get("disk_size")

	if d.HasChange("io_max") {
		update = true

		if v, ok := d.GetOk("io_max"); ok {
			request["IoMax"] = v
		}
	}

	if d.HasChange("io_max_spec") {
		update = true

		if v, ok := d.GetOk("io_max_spec"); ok {
			request["IoMaxSpec"] = v
		}
	}

	if d.HasChange("spec_type") {
		update = true
	}
	request["SpecType"] = d.Get("spec_type")

	if d.HasChange("deploy_type") {
		update = true
	}
	if d.Get("deploy_type").(int) == 4 {
		request["EipModel"] = true
	} else {
		request["EipModel"] = false
	}
	if d.HasChange("eip_max") {
		update = true
	}
	request["EipMax"] = d.Get("eip_max").(int)

	if update {
		action := "UpgradePostPayOrder"
		if d.Get("paid_type").(string) == string(PrePaid) {
			action = "UpgradePrePayOrder"
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, raw, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("disk_size"))}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "DiskSize", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if d.HasChange("io_max") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("io_max"))}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "IoMax", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if d.HasChange("io_max_spec") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("io_max_spec"))}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "IoMaxSpec", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("eip_max"))}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "EipMax", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("spec_type"))}, d.Timeout(schema.TimeoutUpdate), 0*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "SpecType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("partition_num")
		d.SetPartial("disk_size")
		d.SetPartial("io_max")
		d.SetPartial("io_max_spec")
		d.SetPartial("spec_type")
		d.SetPartial("eip_max")
	}

	if d.HasChange("service_version") {
		action := "UpgradeInstanceVersion"
		request := map[string]interface{}{
			"InstanceId": d.Id(),
			"RegionId":   client.RegionId,
		}

		if v, ok := d.GetOk("service_version"); ok {
			request["TargetVersion"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				// means no need to update version
				if IsExpectedErrors(err, []string{"ONS_INIT_ENV_ERROR"}) {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		time.Sleep(60 * time.Second)
		// upgrade service may be last a long time
		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("service_version")
	}

	if d.HasChange("config") {
		action := "UpdateInstanceConfig"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("config"); ok {
			request["Config"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("config")
	}

	update = false
	if d.HasChange("allowed_list") {
		if d.HasChange("allowed_list.0.vpc_list") {
			removed, added := d.GetChange("allowed_list.0.vpc_list")
			deleteAllowedIpReq := map[string]interface{}{
				"RegionId":        client.RegionId,
				"UpdateType":      "delete",
				"AllowedListType": "vpc",
				"InstanceId":      d.Id(),
			}

			for _, vpcList := range removed.(*schema.Set).List() {
				update = true
				vpcListArg := vpcList.(map[string]interface{})

				if portRange, ok := vpcListArg["port_range"]; ok {
					deleteAllowedIpReq["PortRange"] = portRange
				}

				if allowedIpList, ok := vpcListArg["allowed_ip_list"]; ok {
					for _, vpcAllowedIp := range allowedIpList.(*schema.Set).List() {
						deleteAllowedIpReq["AllowedListIp"] = vpcAllowedIp

						if update {
							action := "UpdateAllowedIp"
							conn, err := client.NewAlikafkaClient()
							if err != nil {
								return WrapError(err)
							}

							runtime := util.RuntimeOptions{}
							runtime.SetAutoretry(true)
							wait := incrementalWait(3*time.Second, 3*time.Second)
							err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
								response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, deleteAllowedIpReq, &runtime)
								if err != nil {
									if NeedRetry(err) {
										wait()
										return resource.RetryableError(err)
									}
									return resource.NonRetryableError(err)
								}
								return nil
							})
							addDebug(action, response, deleteAllowedIpReq)

							if err != nil {
								return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
							}

							if fmt.Sprint(response["Success"]) == "false" {
								return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
							}
						}
					}
				}
			}

			update = false
			addAllowedIpReq := map[string]interface{}{
				"RegionId":        client.RegionId,
				"UpdateType":      "add",
				"AllowedListType": "vpc",
				"InstanceId":      d.Id(),
			}

			for _, vpcList := range added.(*schema.Set).List() {
				update = true
				vpcListArg := vpcList.(map[string]interface{})

				if portRange, ok := vpcListArg["port_range"]; ok {
					addAllowedIpReq["PortRange"] = portRange
				}

				if allowedIpList, ok := vpcListArg["allowed_ip_list"]; ok {
					addAllowedIpReq["AllowedListIp"] = convertListToCommaSeparate(allowedIpList.(*schema.Set).List())
				}

				if update {
					action := "UpdateAllowedIp"
					conn, err := client.NewAlikafkaClient()
					if err != nil {
						return WrapError(err)
					}

					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, addAllowedIpReq, &runtime)
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})
					addDebug(action, response, addAllowedIpReq)

					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
					}

					if fmt.Sprint(response["Success"]) == "false" {
						return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
					}
				}
			}
		}

		if d.HasChange("allowed_list.0.internet_list") {
			removed, added := d.GetChange("allowed_list.0.internet_list")
			deleteAllowedIpReq := map[string]interface{}{
				"RegionId":        client.RegionId,
				"UpdateType":      "delete",
				"AllowedListType": "internet",
				"InstanceId":      d.Id(),
			}

			for _, internetList := range removed.(*schema.Set).List() {
				update = true
				internetListArg := internetList.(map[string]interface{})

				if portRange, ok := internetListArg["port_range"]; ok {
					deleteAllowedIpReq["PortRange"] = portRange
				}

				if allowedIpList, ok := internetListArg["allowed_ip_list"]; ok {
					for _, internetAllowedIp := range allowedIpList.(*schema.Set).List() {
						deleteAllowedIpReq["AllowedListIp"] = internetAllowedIp

						if update {
							action := "UpdateAllowedIp"
							conn, err := client.NewAlikafkaClient()
							if err != nil {
								return WrapError(err)
							}

							runtime := util.RuntimeOptions{}
							runtime.SetAutoretry(true)
							wait := incrementalWait(3*time.Second, 3*time.Second)
							err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
								response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, deleteAllowedIpReq, &runtime)
								if err != nil {
									if NeedRetry(err) {
										wait()
										return resource.RetryableError(err)
									}
									return resource.NonRetryableError(err)
								}
								return nil
							})
							addDebug(action, response, deleteAllowedIpReq)

							if err != nil {
								return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
							}

							if fmt.Sprint(response["Success"]) == "false" {
								return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
							}
						}
					}
				}
			}

			update = false
			addAllowedIpReq := map[string]interface{}{
				"RegionId":        client.RegionId,
				"UpdateType":      "add",
				"AllowedListType": "internet",
				"InstanceId":      d.Id(),
			}

			for _, internetList := range added.(*schema.Set).List() {
				update = true
				internetListArg := internetList.(map[string]interface{})

				if portRange, ok := internetListArg["port_range"]; ok {
					addAllowedIpReq["PortRange"] = portRange
				}

				if allowedIpList, ok := internetListArg["allowed_ip_list"]; ok {
					addAllowedIpReq["AllowedListIp"] = convertListToCommaSeparate(allowedIpList.(*schema.Set).List())
				}

				if update {
					action := "UpdateAllowedIp"
					conn, err := client.NewAlikafkaClient()
					if err != nil {
						return WrapError(err)
					}

					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, addAllowedIpReq, &runtime)
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						return nil
					})
					addDebug(action, response, addAllowedIpReq)

					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
					}

					if fmt.Sprint(response["Success"]) == "false" {
						return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
					}
				}
			}
		}

		d.SetPartial("allowed_list")
	}

	d.Partial(false)

	return resourceAliCloudAlikafkaInstanceRead(d, meta)
}

func resourceAliCloudAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	action := "ReleaseInstance"
	conn, err := client.NewAlikafkaClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId":          d.Id(),
		"RegionId":            client.RegionId,
		"ForceDeleteInstance": true,
	}

	// Pre paid instance can not be release.
	if d.Get("paid_type").(string) == string(PrePaid) {
		return nil
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || NeedRetry(err) {
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	stateConf := BuildStateConf([]string{}, []string{"15"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	action = "DeleteInstance"
	request = map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	stateConf = BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func formatSelectedZonesReq(configured []interface{}) string {
	doubleList := make([][]interface{}, len(configured))
	for i, v := range configured {
		doubleList[i] = []interface{}{v}
	}

	if len(doubleList) < 1 {
		return ""
	}

	if len(doubleList) == 1 {
		return "[[\"" + doubleList[0][0].(string) + "\"],[]]"
	}

	result := "[["

	for i := 0; i < len(doubleList); i++ {
		switch i {
		case len(doubleList) - 2:
			result += "\"" + doubleList[i][0].(string) + "\""
		case len(doubleList) - 1:
			result += "],[\"" + doubleList[i][0].(string) + "\"]"
		default:
			result += "\"" + doubleList[i][0].(string) + "\","
		}
	}

	result += "]"

	return result
}
