// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudElasticsearchInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudElasticsearchInstanceCreate,
		Read:   resourceAliCloudElasticsearchInstanceRead,
		Update: resourceAliCloudElasticsearchInstanceUpdate,
		Delete: resourceAliCloudElasticsearchInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(61 * time.Minute),
			Update: schema.DefaultTimeout(360 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arch_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_renew_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client_node_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_efficiency"}, false),
						},
						"amount": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_node_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_ssd", "cloud_essd", "cloud_efficiency"}, false),
						},
						"disk_encryption": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"amount": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntBetween(2, 50),
						},
						"spec": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_kibana_private_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_kibana_public_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"x-pack", "IS"}, false),
			},
			"kibana_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amount": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: IntInSlice([]int{0, 1}),
						},
						"spec": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disk": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: IntInSlice([]int{0}),
						},
					},
				},
			},
			"kibana_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kibana_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kibana_private_security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kibana_private_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: elasticsearchEnableKibanaPrivateDiffSuppressFunc,
			},
			"kibana_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: elasticsearchEnableKibanaPublicDiffSuppressFunc,
			},
			"master_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_ssd", "cloud_essd"}, false),
						},
						"amount": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: IntInSlice([]int{0, 3}),
						},
						"spec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: IntInSlice([]int{0, 20}),
						},
					},
				},
			},
			"order_action_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"UPGRADE", "DOWNGRADE"}, false),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"private_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"public_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: elasticsearchEnablePublicDiffSuppressFunc,
			},
			"renew_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"renewal_duration_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"setting_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"update_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: esVersionDiffSuppressFunc,
			},
			"warm_node_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_efficiency"}, false),
						},
						"disk_encryption": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"amount": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"zone_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 3),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Computed:     true,
				Optional:     true,
				Deprecated:   "Field 'instance_charge_type' has been deprecated since provider version 1.262.0. New field 'payment_type' instead.",
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Computed:     true,
			},
			"kibana_node_spec": {
				Type:       schema.TypeString,
				Computed:   true,
				Optional:   true,
				Deprecated: "Field 'kibana_node_spec' has been deprecated since provider version 1.262.0. New field 'kibana_configuration.spec' instead.",
			},
			"master_node_spec": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'master_node_spec' has been deprecated since provider version 1.262.0. New field 'master_configuration.spec' instead.",
			},
			"master_node_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"cloud_ssd", "cloud_essd"}, false),
				Deprecated:   "Field 'master_node_disk_type' has been deprecated since provider version 1.262.0. New field 'master_configuration.disk_type' instead.",
			},

			"client_node_amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 25),
				Deprecated:   "Field 'client_node_amount' has been deprecated since provider version 1.262.0. New field 'client_node_configuration.amount' instead.",
			},
			"client_node_spec": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'client_node_spec' has been deprecated since provider version 1.262.0. New field 'client_node_configuration.spec' instead.",
			},

			"data_node_amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 50),
				Deprecated:   "Field 'data_node_amount' has been deprecated since provider version 1.262.0. New field 'data_node_configuration.amount' instead.",
			},
			"data_node_spec": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'data_node_spec' has been deprecated since provider version 1.262.0. New field 'data_node_configuration.spec' instead.",
			},
			"data_node_disk_size": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'data_node_disk_size' has been deprecated since provider version 1.262.0. New field 'data_node_configuration.disk' instead.",
			},
			"data_node_disk_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'data_node_disk_type' has been deprecated since provider version 1.262.0. New field 'data_node_configuration.disk_type' instead.",
			},
			"data_node_disk_encrypted": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'data_node_disk_encrypted' has been deprecated since provider version 1.262.0. New field 'data_node_configuration.disk_encrypted' instead.",
			},
			"data_node_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: esDataNodeDiskPerformanceLevelDiffSuppressFunc,
				Deprecated:       "Field 'data_node_disk_performance_level' has been deprecated since provider version 1.262.0. New field 'data_node_configuration.performance_level' instead.",
			},
			"warm_node_amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(2, 50),
				Deprecated:   "Field 'warm_node_amount' has been deprecated since provider version 1.262.0. New field 'warm_node_configuration.amount' instead.",
			},
			"warm_node_spec": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'warm_node_amount' has been deprecated since provider version 1.262.0. New field 'warm_node_configuration.spec' instead.",
			},
			"warm_node_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(500, 20480),
				Deprecated:   "Field 'warm_node_amount' has been deprecated since provider version 1.262.0. New field 'warm_node_configuration.disk' instead.",
			},
			"warm_node_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"cloud_efficiency"}, false),
				Deprecated:   "Field 'warm_node_amount' has been deprecated since provider version 1.262.0. New field 'warm_node_configuration.disk_type' instead.",
			},
			"warm_node_disk_encrypted": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'warm_node_amount' has been deprecated since provider version 1.262.0. New field 'warm_node_configuration.disk_encrypted' instead.",
			},
		},
	}
}

func resourceAliCloudElasticsearchInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/openapi/instances")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	clientNodeConfiguration := make(map[string]interface{})

	if v := d.Get("client_node_configuration"); !IsNil(v) {
		disk1, _ := jsonpath.Get("$[0].disk", v)
		if disk1 != nil && disk1 != "" {
			clientNodeConfiguration["disk"] = disk1
		}
		amount1, _ := jsonpath.Get("$[0].amount", v)
		if amount1 != nil && amount1 != "" {
			clientNodeConfiguration["amount"] = amount1
		}
		diskType1, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType1 != nil && diskType1 != "" {
			clientNodeConfiguration["diskType"] = diskType1
		}
		spec1, _ := jsonpath.Get("$[0].spec", v)
		if spec1 != nil && spec1 != "" {
			clientNodeConfiguration["spec"] = spec1
		}

		request["clientNodeConfiguration"] = clientNodeConfiguration
	}

	nodeSpec := make(map[string]interface{})

	if v := d.Get("data_node_configuration"); !IsNil(v) {
		disk3, _ := jsonpath.Get("$[0].disk", v)
		if disk3 != nil && disk3 != "" {
			nodeSpec["disk"] = disk3
		}
		diskType3, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType3 != nil && diskType3 != "" {
			nodeSpec["diskType"] = diskType3
		}
		spec3, _ := jsonpath.Get("$[0].spec", v)
		if spec3 != nil && spec3 != "" {
			nodeSpec["spec"] = spec3
		}
		performanceLevel1, _ := jsonpath.Get("$[0].performance_level", v)
		if performanceLevel1 != nil && performanceLevel1 != "" {
			nodeSpec["performanceLevel"] = performanceLevel1
		}
		diskEncryption1, _ := jsonpath.Get("$[0].disk_encryption", v)
		if diskEncryption1 != nil && diskEncryption1 != "" {
			nodeSpec["diskEncryption"] = diskEncryption1
		}

		request["nodeSpec"] = nodeSpec
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["paymentType"] = convertElasticsearchInstancepaymentTypeRequest(v.(string))
	}
	masterConfiguration := make(map[string]interface{})

	if v := d.Get("master_configuration"); !IsNil(v) {
		disk5, _ := jsonpath.Get("$[0].disk", v)
		if disk5 != nil && disk5 != "" {
			masterConfiguration["disk"] = disk5
		}
		amount3, _ := jsonpath.Get("$[0].amount", v)
		if amount3 != nil && amount3 != "" {
			masterConfiguration["amount"] = amount3
		}
		spec5, _ := jsonpath.Get("$[0].spec", v)
		if spec5 != nil && spec5 != "" {
			masterConfiguration["spec"] = spec5
		}
		diskType5, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType5 != nil && diskType5 != "" {
			masterConfiguration["diskType"] = diskType5
		}

		request["masterConfiguration"] = masterConfiguration
	}

	warmNodeConfiguration := make(map[string]interface{})

	if v := d.Get("warm_node_configuration"); !IsNil(v) {
		spec7, _ := jsonpath.Get("$[0].spec", v)
		if spec7 != nil && spec7 != "" {
			warmNodeConfiguration["spec"] = spec7
		}
		amount5, _ := jsonpath.Get("$[0].amount", v)
		if amount5 != nil && amount5 != "" {
			warmNodeConfiguration["amount"] = amount5
		}
		diskEncryption3, _ := jsonpath.Get("$[0].disk_encryption", v)
		if diskEncryption3 != nil && diskEncryption3 != "" {
			warmNodeConfiguration["diskEncryption"] = diskEncryption3
		}
		diskType7, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType7 != nil && diskType7 != "" {
			warmNodeConfiguration["diskType"] = diskType7
		}
		disk7, _ := jsonpath.Get("$[0].disk", v)
		if disk7 != nil && disk7 != "" {
			warmNodeConfiguration["disk"] = disk7
		}

		request["warmNodeConfiguration"] = warmNodeConfiguration
	}

	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}
	if password != "" {
		request["esAdminPassword"] = password
	} else {
		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["esAdminPassword"] = decryptResp
	}
	dataNodeConfigurationAmountJsonPath, err := jsonpath.Get("$[0].amount", d.Get("data_node_configuration"))
	if err == nil && dataNodeConfigurationAmountJsonPath.(int) > 0 {
		request["nodeAmount"] = dataNodeConfigurationAmountJsonPath
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	kibanaConfiguration := make(map[string]interface{})

	if v := d.Get("kibana_configuration"); !IsNil(v) {
		disk9, _ := jsonpath.Get("$[0].disk", v)
		if disk9 != nil && disk9 != "" {
			kibanaConfiguration["disk"] = disk9
		}
		amount7, _ := jsonpath.Get("$[0].amount", v)
		if amount7 != nil && amount7 != "" {
			kibanaConfiguration["amount"] = amount7
		}
		spec9, _ := jsonpath.Get("$[0].spec", v)
		if spec9 != nil && spec9 != "" {
			kibanaConfiguration["spec"] = spec9
		}

		request["kibanaConfiguration"] = kibanaConfiguration
	}

	if v, ok := d.GetOk("instance_category"); ok {
		request["instanceCategory"] = v
	}
	if v, ok := d.GetOkExists("zone_count"); ok && v.(int) > 0 {
		request["zoneCount"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["esVersion"] = d.Get("version")
	vswitchId := d.Get("vswitch_id")
	vpcService := VpcService{client}
	vsw, err := vpcService.DescribeVSwitch(vswitchId.(string))
	if err != nil {
		return WrapError(err)
	}
	network := make(map[string]interface{})
	network["type"] = "vpc"
	network["vpcId"] = vsw.VpcId
	network["vswitchId"] = vswitchId
	network["vsArea"] = vsw.ZoneId
	request["networkConfig"] = network

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["paymentType"] = strings.ToLower(v.(string))
	}
	if v := d.Get("period"); v.(int) > 0 {
		paymentInfo := make(map[string]interface{})
		if d.Get("period").(int) >= 12 {
			paymentInfo["duration"] = d.Get("period").(int) / 12
			paymentInfo["pricingCycle"] = string(Year)
		} else {
			paymentInfo["duration"] = d.Get("period").(int)
			paymentInfo["pricingCycle"] = string(Month)
		}
		request["paymentInfo"] = paymentInfo
	}

	if v, ok := d.GetOk("kibana_node_spec"); ok {
		kibanaNode := make(map[string]interface{})
		kibanaNode["spec"] = v.(string)
		kibanaNode["disk"] = "0"
		kibanaNode["amount"] = 1
		request["haveKibana"] = true
		request["kibanaConfiguration"] = kibanaNode
	}

	if v, ok := d.GetOk("master_node_spec"); ok {
		masterDiskType, diskExits := d.GetOkExists("master_node_disk_type")
		if !diskExits {
			return WrapError(fmt.Errorf("CreateInstance error: master_node_disk_type is null"))
		}
		masterNode := make(map[string]interface{})
		masterNode["spec"] = v.(string)
		masterNode["amount"] = "3"
		masterNode["disk"] = "20"
		masterNode["diskType"] = masterDiskType.(string)
		request["advancedDedicateMaster"] = true
		request["masterConfiguration"] = masterNode
	}

	if v, ok := d.GetOk("client_node_spec"); ok {
		clientNode := make(map[string]interface{})
		clientNode["spec"] = v.(string)
		clientNode["disk"] = "20"
		clientNode["diskType"] = "cloud_efficiency"
		if d.Get("client_node_amount") == nil {
			clientNode["amount"] = 2
		} else {
			clientNode["amount"] = d.Get("client_node_amount")
		}
		request["haveClientNode"] = true
		request["clientNodeConfiguration"] = clientNode
	}

	if v, ok := d.GetOk("data_node_amount"); ok {
		request["nodeAmount"] = v
		dataNodeSpec := make(map[string]interface{})
		dataNodeSpec["spec"] = d.Get("data_node_spec")
		dataNodeSpec["disk"] = d.Get("data_node_disk_size")
		dataNodeSpec["diskType"] = d.Get("data_node_disk_type")
		dataNodeSpec["diskEncryption"] = d.Get("data_node_disk_encrypted")
		performanceLevel := d.Get("data_node_disk_performance_level")
		if performanceLevel != "" {
			dataNodeSpec["performanceLevel"] = performanceLevel
		}
		request["nodeSpec"] = dataNodeSpec
	}

	if d.Get("warm_node_spec") != nil && d.Get("warm_node_spec") != "" && nil != d.Get("warm_node_amount") && d.Get("warm_node_amount").(int) > 0 && nil != d.Get("warm_node_disk_size") && d.Get("warm_node_disk_size").(int) > 0 {
		warmNode := make(map[string]interface{})
		warmNode["spec"] = d.Get("warm_node_spec")
		warmNode["amount"] = d.Get("warm_node_amount")
		warmNode["disk"] = d.Get("warm_node_disk_size")
		warmNode["diskType"] = d.Get("warm_node_disk_type")
		warmNode["diskEncryption"] = d.Get("warm_node_disk_encrypted")
		request["warmNode"] = true
		request["warmNodeConfiguration"] = warmNode
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TokenPreviousRequestProcessError", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_elasticsearch_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Result.instanceId", response)
	d.SetId(fmt.Sprint(id))

	elasticsearchServiceV2 := ElasticsearchServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudElasticsearchInstanceUpdate(d, meta)
}

func resourceAliCloudElasticsearchInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchServiceV2 := ElasticsearchServiceV2{client}

	objectRaw, err := elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_elasticsearch_instance DescribeElasticsearchInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("arch_type", objectRaw["archType"])

	d.Set("vswitch_id", objectRaw["networkConfig"].(map[string]interface{})["vswitchId"])
	d.Set("kibana_node_spec", objectRaw["kibanaConfiguration"].(map[string]interface{})["spec"])
	d.Set("master_node_spec", objectRaw["masterConfiguration"].(map[string]interface{})["spec"])
	d.Set("master_node_disk_type", objectRaw["masterConfiguration"].(map[string]interface{})["diskType"])

	if objectRaw["clientNodeConfiguration"] != nil && objectRaw["clientNodeConfiguration"].(map[string]interface{})["spec"] != nil {
		d.Set("client_node_amount", objectRaw["clientNodeConfiguration"].(map[string]interface{})["amount"])
		d.Set("client_node_spec", objectRaw["clientNodeConfiguration"].(map[string]interface{})["spec"])
	}

	d.Set("data_node_amount", objectRaw["nodeAmount"])
	d.Set("data_node_spec", objectRaw["nodeSpec"].(map[string]interface{})["spec"])
	d.Set("data_node_disk_size", objectRaw["nodeSpec"].(map[string]interface{})["disk"])
	d.Set("data_node_disk_type", objectRaw["nodeSpec"].(map[string]interface{})["diskType"])
	d.Set("data_node_disk_encrypted", objectRaw["nodeSpec"].(map[string]interface{})["diskEncryption"])
	d.Set("data_node_disk_performance_level", objectRaw["nodeSpec"].(map[string]interface{})["performanceLevel"])

	if objectRaw["warmNodeConfiguration"] != nil && objectRaw["warmNodeConfiguration"].(map[string]interface{})["spec"] != nil {
		d.Set("warm_node_amount", objectRaw["warmNodeConfiguration"].(map[string]interface{})["amount"])
		d.Set("warm_node_spec", objectRaw["warmNodeConfiguration"].(map[string]interface{})["spec"])
		d.Set("warm_node_disk_size", objectRaw["warmNodeConfiguration"].(map[string]interface{})["disk"])
		d.Set("warm_node_disk_type", objectRaw["warmNodeConfiguration"].(map[string]interface{})["diskType"])
		d.Set("warm_node_disk_encrypted", objectRaw["warmNodeConfiguration"].(map[string]interface{})["diskEncryption"])
	}

	d.Set("create_time", objectRaw["createdAt"])
	d.Set("description", objectRaw["description"])
	d.Set("domain", objectRaw["domain"])
	d.Set("enable_kibana_private_network", objectRaw["enableKibanaPrivateNetwork"])
	d.Set("enable_kibana_public_network", objectRaw["enableKibanaPublicNetwork"])
	d.Set("enable_public", objectRaw["enablePublic"])
	d.Set("instance_category", objectRaw["instanceCategory"])
	d.Set("kibana_domain", objectRaw["kibanaDomain"])
	d.Set("kibana_port", objectRaw["kibanaPort"])
	d.Set("payment_type", convertElasticsearchInstanceResultpaymentTypeResponse(objectRaw["paymentType"]))
	d.Set("instance_charge_type", getChargeType(objectRaw["paymentType"].(string)))
	d.Set("protocol", objectRaw["protocol"])
	d.Set("public_domain", objectRaw["publicDomain"])
	d.Set("public_port", objectRaw["publicPort"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("setting_config", objectRaw["esConfig"])
	d.Set("status", objectRaw["status"])
	d.Set("version", objectRaw["esVersion"])
	d.Set("zone_count", objectRaw["zoneCount"])

	clientNodeConfigurationMaps := make([]map[string]interface{}, 0)
	clientNodeConfigurationMap := make(map[string]interface{})
	clientNodeConfigurationRaw := make(map[string]interface{})
	if objectRaw["clientNodeConfiguration"] != nil {
		clientNodeConfigurationRaw = objectRaw["clientNodeConfiguration"].(map[string]interface{})
	}
	if len(clientNodeConfigurationRaw) > 0 {
		clientNodeConfigurationMap["amount"] = clientNodeConfigurationRaw["amount"]
		clientNodeConfigurationMap["disk"] = clientNodeConfigurationRaw["disk"]
		clientNodeConfigurationMap["disk_type"] = clientNodeConfigurationRaw["diskType"]
		clientNodeConfigurationMap["spec"] = clientNodeConfigurationRaw["spec"]

		clientNodeConfigurationMaps = append(clientNodeConfigurationMaps, clientNodeConfigurationMap)
	}
	if err := d.Set("client_node_configuration", clientNodeConfigurationMaps); err != nil {
		return err
	}
	dataNodeConfigurationMaps := make([]map[string]interface{}, 0)
	dataNodeConfigurationMap := make(map[string]interface{})

	dataNodeConfigurationMap["amount"] = objectRaw["nodeAmount"]

	nodeSpecRaw := make(map[string]interface{})
	if objectRaw["nodeSpec"] != nil {
		nodeSpecRaw = objectRaw["nodeSpec"].(map[string]interface{})
	}
	if len(nodeSpecRaw) > 0 {
		dataNodeConfigurationMap["disk"] = nodeSpecRaw["disk"]
		dataNodeConfigurationMap["disk_encryption"] = nodeSpecRaw["diskEncryption"]
		dataNodeConfigurationMap["disk_type"] = nodeSpecRaw["diskType"]
		dataNodeConfigurationMap["performance_level"] = nodeSpecRaw["performanceLevel"]
		dataNodeConfigurationMap["spec"] = nodeSpecRaw["spec"]
	}
	dataNodeConfigurationMaps = append(dataNodeConfigurationMaps, dataNodeConfigurationMap)
	if err := d.Set("data_node_configuration", dataNodeConfigurationMaps); err != nil {
		return err
	}
	kibanaConfigurationMaps := make([]map[string]interface{}, 0)
	kibanaConfigurationMap := make(map[string]interface{})
	kibanaConfigurationRaw := make(map[string]interface{})
	if objectRaw["kibanaConfiguration"] != nil {
		kibanaConfigurationRaw = objectRaw["kibanaConfiguration"].(map[string]interface{})
	}
	if len(kibanaConfigurationRaw) > 0 {
		kibanaConfigurationMap["amount"] = kibanaConfigurationRaw["amount"]
		kibanaConfigurationMap["disk"] = kibanaConfigurationRaw["disk"]
		kibanaConfigurationMap["spec"] = kibanaConfigurationRaw["spec"]

		kibanaConfigurationMaps = append(kibanaConfigurationMaps, kibanaConfigurationMap)
	}
	if err := d.Set("kibana_configuration", kibanaConfigurationMaps); err != nil {
		return err
	}
	kibanaPrivateIPWhitelistRaw := make([]interface{}, 0)
	if objectRaw["kibanaPrivateIPWhitelist"] != nil {
		kibanaPrivateIPWhitelistRaw = convertToInterfaceArray(objectRaw["kibanaPrivateIPWhitelist"])
	}

	d.Set("kibana_private_whitelist", kibanaPrivateIPWhitelistRaw)
	kibanaIPWhitelistRaw := make([]interface{}, 0)
	if objectRaw["kibanaIPWhitelist"] != nil {
		kibanaIPWhitelistRaw = convertToInterfaceArray(objectRaw["kibanaIPWhitelist"])
	}

	d.Set("kibana_whitelist", kibanaIPWhitelistRaw)
	masterConfigurationMaps := make([]map[string]interface{}, 0)
	masterConfigurationMap := make(map[string]interface{})
	masterConfigurationRaw := make(map[string]interface{})
	if objectRaw["masterConfiguration"] != nil {
		masterConfigurationRaw = objectRaw["masterConfiguration"].(map[string]interface{})
	}
	if len(masterConfigurationRaw) > 0 {
		masterConfigurationMap["amount"] = masterConfigurationRaw["amount"]
		masterConfigurationMap["disk"] = masterConfigurationRaw["disk"]
		masterConfigurationMap["disk_type"] = masterConfigurationRaw["diskType"]
		masterConfigurationMap["spec"] = masterConfigurationRaw["spec"]

		masterConfigurationMaps = append(masterConfigurationMaps, masterConfigurationMap)
	}
	if err := d.Set("master_configuration", masterConfigurationMaps); err != nil {
		return err
	}
	privateNetworkIpWhiteListRaw := make([]interface{}, 0)
	if objectRaw["privateNetworkIpWhiteList"] != nil {
		privateNetworkIpWhiteListRaw = convertToInterfaceArray(objectRaw["privateNetworkIpWhiteList"])
	}

	d.Set("private_whitelist", privateNetworkIpWhiteListRaw)
	publicIpWhitelistRaw := make([]interface{}, 0)
	if objectRaw["publicIpWhitelist"] != nil {
		publicIpWhitelistRaw = convertToInterfaceArray(objectRaw["publicIpWhitelist"])
	}

	d.Set("public_whitelist", publicIpWhitelistRaw)
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	warmNodeConfigurationMaps := make([]map[string]interface{}, 0)
	warmNodeConfigurationMap := make(map[string]interface{})
	warmNodeConfigurationRaw := make(map[string]interface{})
	if objectRaw["warmNodeConfiguration"] != nil {
		warmNodeConfigurationRaw = objectRaw["warmNodeConfiguration"].(map[string]interface{})
	}
	if len(warmNodeConfigurationRaw) > 0 {
		warmNodeConfigurationMap["amount"] = warmNodeConfigurationRaw["amount"]
		warmNodeConfigurationMap["disk"] = warmNodeConfigurationRaw["disk"]
		warmNodeConfigurationMap["disk_encryption"] = warmNodeConfigurationRaw["diskEncryption"]
		warmNodeConfigurationMap["disk_type"] = warmNodeConfigurationRaw["diskType"]
		warmNodeConfigurationMap["spec"] = warmNodeConfigurationRaw["spec"]

		warmNodeConfigurationMaps = append(warmNodeConfigurationMaps, warmNodeConfigurationMap)
	}
	if err := d.Set("warm_node_configuration", warmNodeConfigurationMaps); err != nil {
		return err
	}

	checkValue00 := d.Get("payment_type")
	if checkValue00 == "Subscription" {
		objectRaw, err = elasticsearchServiceV2.DescribeInstanceQueryAvailableInstances(d)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		d.Set("auto_renew_duration", objectRaw["RenewalDuration"])
		d.Set("renew_status", objectRaw["RenewStatus"])
		d.Set("renewal_duration_unit", objectRaw["RenewalDurationUnit"])

	}
	checkValue00 = d.Get("arch_type")
	if checkValue00 == "public" {
		objectRaw, err = elasticsearchServiceV2.DescribeInstanceListKibanaPvlNetwork(d.Id())
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		securityGroupsRaw, _ := jsonpath.Get("$.securityGroups", objectRaw)
		if securityGroupsRaw != nil {
			d.Set("kibana_private_security_group_id", securityGroupsRaw.([]interface{})[0])
		}

	}

	return nil
}

func resourceAliCloudElasticsearchInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.IsNewResource() && d.HasChange("instance_charge_type") {
		if err := updateInstanceChargeType(d, meta); err != nil {
			return WrapError(err)
		}
	} else if !d.IsNewResource() && d.Get("instance_charge_type").(string) == string(PrePaid) && d.HasChange("period") {
		if err := renewInstance(d, meta); err != nil {
			return WrapError(err)
		}
	}

	if !d.IsNewResource() && d.HasChange("kibana_node_spec") {
		if err := updateKibanaNode(d, meta); err != nil {
			return WrapError(err)
		}
	}

	if !d.IsNewResource() && (d.HasChange("master_node_spec") || d.HasChange("master_node_disk_type")) {
		if err := updateMasterNode(d, meta); err != nil {
			return WrapError(err)
		}
	}

	if !d.IsNewResource() && (d.HasChange("client_node_spec") || d.HasChange("client_node_amount")) {
		if err := updateClientNode(d, meta); err != nil {
			return WrapError(err)
		}
	}

	if !d.IsNewResource() && d.HasChange("data_node_amount") {
		if err := updateDataNodeAmount(d, meta); err != nil {
			return WrapError(err)
		}
	}

	if !d.IsNewResource() && d.HasChanges("data_node_spec", "data_node_disk_size", "data_node_disk_performance_level") {
		if err := updateDataNodeSpec(d, meta); err != nil {
			return WrapError(err)
		}
	}

	if !d.IsNewResource() && d.HasChanges("warm_node_spec", "warm_node_amount", "warm_node_disk_size") {
		if err := updateWarmNode(d, meta); err != nil {
			return WrapError(err)
		}
	}

	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	elasticsearchServiceV2 := ElasticsearchServiceV2{client}
	objectRaw, _ := elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())

	if d.HasChange("protocol") {
		var err error
		target := d.Get("protocol").(string)

		currentStatus, err := jsonpath.Get("protocol", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "protocol", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "HTTPS" {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/open-https", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "HTTP" {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/close-https", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	if d.HasChange("enable_kibana_private_network") {
		var err error
		target := d.Get("enable_kibana_private_network").(bool)

		currentStatus, err := jsonpath.Get("enableKibanaPrivateNetwork", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "enableKibanaPrivateNetwork", objectRaw)
		}
		if formatBool(currentStatus) != target {
			enableTriggerNetworkTrue := false
			checkValue00 := objectRaw["enableKibanaPrivateNetwork"]
			checkValue01 := objectRaw["archType"]
			if !(checkValue00 == true) && !(checkValue01 == "public") {
				enableTriggerNetworkTrue = true
			}
			if enableTriggerNetworkTrue && target == true {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				request["actionType"] = "OPEN"
				request["networkType"] = "PRIVATE"
				request["nodeType"] = "KIBANA"

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			enableTriggerNetworkFalse := false
			checkValue00 = objectRaw["enableKibanaPrivateNetwork"]
			checkValue01 = objectRaw["archType"]
			if !(checkValue00 == false) && !(checkValue01 == "public") {
				enableTriggerNetworkFalse = true
			}
			if enableTriggerNetworkFalse && target == false {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				request["networkType"] = "PRIVATE"
				request["nodeType"] = "KIBANA"
				request["actionType"] = "CLOSE"

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			enableEnableKibanaPvlNetworkTrue := false
			checkValue00 = objectRaw["archType"]
			if checkValue00 == "public" {
				enableEnableKibanaPvlNetworkTrue = true
			}
			if enableEnableKibanaPvlNetworkTrue && target == true {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/enable-kibana-private", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				if v, ok := d.GetOk("kibana_private_security_group_id"); ok {
					localData, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					securityGroupsMapsArray := convertToInterfaceArray(localData)

					request["securityGroups"] = securityGroupsMapsArray
				}

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			enableDisableKibanaPvlNetworkFalse := false
			checkValue00 = objectRaw["archType"]
			if checkValue00 == "public" {
				enableDisableKibanaPvlNetworkFalse = true
			}
			if enableDisableKibanaPvlNetworkFalse && target == false {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/disable-kibana-private", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	if d.HasChange("enable_public") {
		var err error
		target := d.Get("enable_public").(bool)

		currentStatus, err := jsonpath.Get("enablePublic", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "enablePublic", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				request["nodeType"] = "WORKER"
				request["actionType"] = "OPEN"
				request["networkType"] = "PUBLIC"

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == false {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				request["nodeType"] = "WORKER"
				request["networkType"] = "PUBLIC"
				request["actionType"] = "CLOSE"

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	if d.HasChange("enable_kibana_public_network") {
		var err error
		target := d.Get("enable_kibana_public_network").(bool)

		currentStatus, err := jsonpath.Get("enableKibanaPublicNetwork", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "enableKibanaPublicNetwork", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				request["actionType"] = "OPEN"
				request["nodeType"] = "KIBANA"
				request["networkType"] = "PUBLIC"

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == false {
				InstanceId := d.Id()
				action := fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", InstanceId)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})

				request["nodeType"] = "KIBANA"
				request["networkType"] = "PUBLIC"
				request["actionType"] = "CLOSE"

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
				elasticsearchServiceV2 := ElasticsearchServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	InstanceId := d.Id()
	action := fmt.Sprintf("/openapi/instances/%s", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("client_node_configuration") {
		update = true
	}
	clientNodeConfiguration := make(map[string]interface{})

	if v := d.Get("client_node_configuration"); !IsNil(v) || d.HasChange("client_node_configuration") {
		disk1, _ := jsonpath.Get("$[0].disk", v)
		if disk1 != nil && disk1 != "" {
			clientNodeConfiguration["disk"] = disk1
		}
		amount1, _ := jsonpath.Get("$[0].amount", v)
		if amount1 != nil && amount1 != "" {
			clientNodeConfiguration["amount"] = amount1
		}
		diskType1, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType1 != nil && diskType1 != "" {
			clientNodeConfiguration["diskType"] = diskType1
		}
		spec1, _ := jsonpath.Get("$[0].spec", v)
		if spec1 != nil && spec1 != "" {
			clientNodeConfiguration["spec"] = spec1
		}

		request["clientNodeConfiguration"] = clientNodeConfiguration
	}

	if !d.IsNewResource() && d.HasChange("data_node_configuration") {
		update = true
	}
	nodeSpec := make(map[string]interface{})

	if v := d.Get("data_node_configuration"); !IsNil(v) || d.HasChange("data_node_configuration") {
		disk3, _ := jsonpath.Get("$[0].disk", v)
		if disk3 != nil && disk3 != "" {
			nodeSpec["disk"] = disk3
		}
		diskType3, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType3 != nil && diskType3 != "" {
			nodeSpec["diskType"] = diskType3
		}
		spec3, _ := jsonpath.Get("$[0].spec", v)
		if spec3 != nil && spec3 != "" {
			nodeSpec["spec"] = spec3
		}
		performanceLevel1, _ := jsonpath.Get("$[0].performance_level", v)
		if performanceLevel1 != nil && performanceLevel1 != "" {
			nodeSpec["performanceLevel"] = performanceLevel1
		}
		diskEncryption1, _ := jsonpath.Get("$[0].disk_encryption", v)
		if diskEncryption1 != nil && diskEncryption1 != "" {
			nodeSpec["diskEncryption"] = diskEncryption1
		}

		request["nodeSpec"] = nodeSpec
	}

	if !d.IsNewResource() && d.HasChange("warm_node_configuration") {
		update = true
	}
	warmNodeConfiguration := make(map[string]interface{})

	if v := d.Get("warm_node_configuration"); !IsNil(v) || d.HasChange("warm_node_configuration") {
		spec5, _ := jsonpath.Get("$[0].spec", v)
		if spec5 != nil && spec5 != "" {
			warmNodeConfiguration["spec"] = spec5
		}
		amount3, _ := jsonpath.Get("$[0].amount", v)
		if amount3 != nil && amount3 != "" {
			warmNodeConfiguration["amount"] = amount3
		}
		diskEncryption3, _ := jsonpath.Get("$[0].disk_encryption", v)
		if diskEncryption3 != nil && diskEncryption3 != "" {
			warmNodeConfiguration["diskEncryption"] = diskEncryption3
		}
		diskType5, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType5 != nil && diskType5 != "" {
			warmNodeConfiguration["diskType"] = diskType5
		}
		disk5, _ := jsonpath.Get("$[0].disk", v)
		if disk5 != nil && disk5 != "" {
			warmNodeConfiguration["disk"] = disk5
		}

		request["warmNodeConfiguration"] = warmNodeConfiguration
	}

	if !d.IsNewResource() && d.HasChange("master_configuration") {
		update = true
	}
	masterConfiguration := make(map[string]interface{})

	if v := d.Get("master_configuration"); !IsNil(v) || d.HasChange("master_configuration") {
		disk7, _ := jsonpath.Get("$[0].disk", v)
		if disk7 != nil && disk7 != "" {
			masterConfiguration["disk"] = disk7
		}
		amount5, _ := jsonpath.Get("$[0].amount", v)
		if amount5 != nil && amount5 != "" {
			masterConfiguration["amount"] = amount5
		}
		spec7, _ := jsonpath.Get("$[0].spec", v)
		if spec7 != nil && spec7 != "" {
			masterConfiguration["spec"] = spec7
		}
		diskType7, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType7 != nil && diskType7 != "" {
			masterConfiguration["diskType"] = diskType7
		}

		request["masterConfiguration"] = masterConfiguration
	}

	if !d.IsNewResource() && d.HasChange("data_node_configuration.0.amount") {
		update = true
	}
	if v, ok := d.GetOkExists("data_node_configuration"); ok || d.HasChange("data_node_configuration") {
		dataNodeConfigurationAmountJsonPath, err := jsonpath.Get("$[0].amount", v)
		if err == nil && dataNodeConfigurationAmountJsonPath != "" {
			request["nodeAmount"] = dataNodeConfigurationAmountJsonPath
		}
	}
	if !d.IsNewResource() && d.HasChange("kibana_configuration") {
		update = true
	}
	kibanaConfiguration := make(map[string]interface{})

	if v := d.Get("kibana_configuration"); !IsNil(v) || d.HasChange("kibana_configuration") {
		disk9, _ := jsonpath.Get("$[0].disk", v)
		if disk9 != nil && disk9 != "" {
			kibanaConfiguration["disk"] = disk9
		}
		amount7, _ := jsonpath.Get("$[0].amount", v)
		if amount7 != nil && amount7 != "" {
			kibanaConfiguration["amount"] = amount7
		}
		spec9, _ := jsonpath.Get("$[0].spec", v)
		if spec9 != nil && spec9 != "" {
			kibanaConfiguration["spec"] = spec9
		}

		request["kibanaConfiguration"] = kibanaConfiguration
	}

	if v, ok := d.GetOkExists("force"); ok {
		query["force"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("order_action_type"); ok {
		query["orderActionType"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"TheSpecNotEnoughInDetail", "ConcurrencyUpdateInstanceConflict", "InstanceDuplicateScheduledTask", "ServiceUnavailable", "InternalServerError", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{"inactive"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/description", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavaliable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/admin-pwd", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("kms_encrypted_password") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("password") {
		update = true
	}
	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}
	if password != "" {
		request["esAdminPassword"] = password
	} else {
		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["esAdminPassword"] = decryptResp
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/actions/convert-pay-type", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	if v, ok := d.GetOk("payment_type"); ok || d.HasChange("payment_type") {
		request["paymentType"] = convertElasticsearchInstancepaymentTypeRequest(v.(string))
	}
	if request["paymentType"] == convertElasticsearchInstancepaymentTypeRequest("Subscription") {
		if v := d.Get("period"); v.(int) > 0 {
			paymentInfo := make(map[string]interface{})
			if d.Get("period").(int) >= 12 {
				paymentInfo["duration"] = d.Get("period").(int) / 12
				paymentInfo["pricingCycle"] = string(Year)
			} else {
				paymentInfo["duration"] = d.Get("period").(int)
				paymentInfo["pricingCycle"] = string(Month)
			}
			request["paymentInfo"] = paymentInfo
		}
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("payment_type"))}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "paymentType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/instance-settings", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if v, ok := d.GetOk("update_strategy"); ok {
		query["updateStrategy"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOkExists("force"); ok {
		query["force"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if d.HasChange("setting_config") {
		update = true
	}
	if v, ok := d.GetOk("setting_config"); ok || d.HasChange("setting_config") {
		request["esConfig"] = v
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	objectRaw, _ = elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())
	enableModifyWhiteIps3 := false
	checkValue00 := objectRaw["enablePublic"]
	if checkValue00 == true {
		enableModifyWhiteIps3 = true
	}
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/actions/modify-white-ips", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	request["nodeType"] = "WORKER"
	request["modifyMode"] = "Cover"
	if d.HasChange("public_whitelist") {
		update = true
	}
	if v, ok := d.GetOk("public_whitelist"); ok || d.HasChange("public_whitelist") {
		whiteIpListMapsArray := convertToInterfaceArray(v)

		request["whiteIpList"] = whiteIpListMapsArray
	}

	request["networkType"] = "PUBLIC"

	body = request
	if update && enableModifyWhiteIps3 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/actions/modify-white-ips", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	request["nodeType"] = "WORKER"
	request["modifyMode"] = "Cover"
	if d.HasChange("private_whitelist") {
		update = true
	}
	if v, ok := d.GetOk("private_whitelist"); ok || d.HasChange("private_whitelist") {
		whiteIpListMapsArray := convertToInterfaceArray(v)

		request["whiteIpList"] = whiteIpListMapsArray
	}

	request["networkType"] = "PRIVATE"

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	objectRaw, _ = elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())
	enableModifyWhiteIps4 := false
	checkValue00 = objectRaw["enableKibanaPublicNetwork"]
	if checkValue00 == true {
		enableModifyWhiteIps4 = true
	}
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/actions/modify-white-ips", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	request["modifyMode"] = "Cover"
	request["nodeType"] = "KIBANA"
	request["networkType"] = "PUBLIC"
	if d.HasChange("kibana_whitelist") {
		update = true
	}
	if v, ok := d.GetOk("kibana_whitelist"); ok || d.HasChange("kibana_whitelist") {
		whiteIpListMapsArray := convertToInterfaceArray(v)

		request["whiteIpList"] = whiteIpListMapsArray
	}

	body = request
	if update && enableModifyWhiteIps4 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	objectRaw, _ = elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())
	enableModifyWhiteIps5 := false
	checkValue00 = objectRaw["enableKibanaPrivateNetwork"]
	if checkValue00 == true {
		enableModifyWhiteIps5 = true
	}
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/actions/modify-white-ips", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	request["modifyMode"] = "Cover"
	request["networkType"] = "PRIVATE"
	request["nodeType"] = "KIBANA"
	if d.HasChange("kibana_private_whitelist") {
		update = true
	}
	if v, ok := d.GetOk("kibana_private_whitelist"); ok || d.HasChange("kibana_private_whitelist") {
		whiteIpListMapsArray := convertToInterfaceArray(v)

		request["whiteIpList"] = whiteIpListMapsArray
	}

	body = request
	if update && enableModifyWhiteIps5 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	objectRaw, _ = elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())
	enableUpdateKibanaPvlNetwork1 := false
	checkValue00 = objectRaw["archType"]
	checkValue01 := objectRaw["enableKibanaPrivateNetwork"]
	if (checkValue00 == "public") && (checkValue01 == true) {
		enableUpdateKibanaPvlNetwork1 = true
	}
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/actions/update-kibana-private", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("kibana_private_security_group_id") {
		update = true
	}
	if v, ok := d.GetOk("kibana_private_security_group_id"); ok || d.HasChange("kibana_private_security_group_id") {
		localData, err := jsonpath.Get("$", v)
		if err != nil {
			return WrapError(err)
		}
		securityGroupsMapsArray := convertToInterfaceArray(localData)

		request["securityGroups"] = securityGroupsMapsArray
	}

	body = request
	if update && enableUpdateKibanaPvlNetwork1 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchInstanceStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	{
		update = false
		objectRaw, _ = elasticsearchServiceV2.DescribeElasticsearchInstance(d.Id())
		enableSetRenewal1 := false
		checkValue00 = convertElasticsearchInstanceResultpaymentTypeResponse(objectRaw["paymentType"])
		if checkValue00 == "Subscription" {
			enableSetRenewal1 = true
		}
		action = "SetRenewal"
		request = make(map[string]interface{})
		query := make(map[string]interface{})
		query["InstanceIDs"] = d.Id()

		if d.HasChange("auto_renew_duration") {
			update = true
		}
		if v, ok := d.GetOkExists("auto_renew_duration"); ok {
			query["RenewalPeriod"] = strconv.Itoa(v.(int))
		}

		query["ProductType"] = "elasticsearchpre"
		if d.HasChange("renewal_duration_unit") {
			update = true
		}
		if v, ok := d.GetOk("renewal_duration_unit"); ok {
			query["RenewalPeriodUnit"] = v.(string)
		}

		if d.HasChange("renew_status") {
			update = true
		}
		if v, ok := d.GetOk("renew_status"); ok {
			query["RenewalStatus"] = v.(string)
		}

		query["ProductCode"] = "elasticsearch"
		var endpoint string
		request["ProductCode"] = "elasticsearch"
		request["ProductType"] = "elasticsearchpre"
		if client.IsInternationalAccount() {
			request["ProductType"] = "elasticsearchpre_intl"
		}

		if update && enableSetRenewal1 {
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
						request["ProductCode"] = "elasticsearch"
						request["ProductType"] = "elasticsearchpre_intl"
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
	}

	if d.HasChange("tags") {
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		if err := elasticsearchServiceV2.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudElasticsearchInstanceRead(d, meta)
}

func resourceAliCloudElasticsearchInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	InstanceId := d.Id()
	action := fmt.Sprintf("/openapi/instances/%s", InstanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("elasticsearch", "2017-06-13", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TokenPreviousRequestProcessError", "ServiceUnavailable", "InstanceActivating", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertElasticsearchInstanceResultpaymentTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "postpaid":
		return "PayAsYouGo"
	case "prepaid":
		return "Subscription"
	}
	return source
}
func convertElasticsearchInstancepaymentTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "prepaid"
	case "PayAsYouGo":
		return "postpaid"
	}
	return source
}
