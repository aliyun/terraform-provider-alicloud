package alicloud

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudElasticsearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudElasticsearchCreate,
		Read:   resourceAlicloudElasticsearchRead,
		Update: resourceAlicloudElasticsearchUpdate,
		Delete: resourceAlicloudElasticsearchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Basic instance information
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w\-.]{0,30}$`), "be 0 to 30 characters in length and can contain numbers, letters, underscores, (_) and hyphens (-). It must start with a letter, a number or Chinese character."),
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
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
			"version": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: esVersionDiffSuppressFunc,
				ForceNew:         true,
			},
			"tags": tagsSchema(),

			// Life cycle
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
				Optional:     true,
			},

			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},

			// Data node configuration
			"data_node_amount": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(2, 50),
			},

			"data_node_spec": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_node_disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"data_node_disk_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_node_disk_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"private_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"enable_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Computed:         true,
				DiffSuppressFunc: elasticsearchEnablePublicDiffSuppressFunc,
			},

			"master_node_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Kibana node configuration
			"kibana_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kibana_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"enable_kibana_public_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"kibana_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Computed:         true,
				DiffSuppressFunc: elasticsearchEnableKibanaPublicDiffSuppressFunc,
			},

			"enable_kibana_private_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kibana_private_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Computed:         true,
				DiffSuppressFunc: elasticsearchEnableKibanaPrivateDiffSuppressFunc,
			},

			"zone_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3),
				Default:      1,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudElasticsearchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	request, err := buildElasticsearchCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"TokenPreviousRequestProcessError"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.CreateInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_elasticsearch_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*elasticsearch.CreateInstanceResponse)
	d.SetId(response.Result.InstanceId)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudElasticsearchUpdate(d, meta)
}

func resourceAlicloudElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	object, err := elasticsearchService.DescribeElasticsearchInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", object.Result.Description)
	d.Set("status", object.Result.Status)
	d.Set("vswitch_id", object.Result.NetworkConfig.VswitchId)

	d.Set("private_whitelist", filterWhitelist(object.Result.EsIPWhitelist, d.Get("private_whitelist").(*schema.Set)))
	d.Set("public_whitelist", filterWhitelist(object.Result.PublicIpWhitelist, d.Get("public_whitelist").(*schema.Set)))
	d.Set("enable_public", object.Result.EnablePublic)
	d.Set("version", object.Result.EsVersion)
	d.Set("instance_charge_type", getChargeType(object.Result.PaymentType))

	d.Set("domain", object.Result.Domain)
	d.Set("port", object.Result.Port)

	// Kibana configuration
	d.Set("enable_kibana_public_network", object.Result.EnableKibanaPublicNetwork)
	d.Set("kibana_whitelist", filterWhitelist(object.Result.KibanaIPWhitelist, d.Get("kibana_whitelist").(*schema.Set)))
	if object.Result.EnableKibanaPublicNetwork {
		d.Set("kibana_domain", object.Result.KibanaDomain)
		d.Set("kibana_port", object.Result.KibanaPort)
	}

	d.Set("enable_kibana_private_network", object.Result.EnableKibanaPrivateNetwork)
	d.Set("kibana_private_whitelist", filterWhitelist(object.Result.KibanaPrivateIPWhitelist, d.Get("kibana_private_whitelist").(*schema.Set)))

	// Data node configuration
	d.Set("data_node_amount", object.Result.NodeAmount)
	d.Set("data_node_spec", object.Result.NodeSpec.Spec)
	d.Set("data_node_disk_size", object.Result.NodeSpec.Disk)
	d.Set("data_node_disk_type", object.Result.NodeSpec.DiskType)
	d.Set("data_node_disk_encrypted", object.Result.NodeSpec.DiskEncryption)
	d.Set("master_node_spec", object.Result.MasterConfiguration.Spec)

	// Cross zone configuration
	d.Set("zone_count", object.Result.ZoneCount)
	d.Set("resource_group_id", object.Result.ResourceGroupId)

	// tags
	tags, err := elasticsearchService.DescribeElasticsearchTags(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", tags)
	}

	return nil
}

func resourceAlicloudElasticsearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if d.HasChange("description") {
		if err := updateDescription(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("description")
	}

	if d.HasChange("private_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(WORKER)
		content["whiteIpList"] = d.Get("private_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("private_whitelist")
	}

	if d.HasChange("enable_public") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(WORKER)
		content["actionType"] = elasticsearchService.getActionType(d.Get("enable_public").(bool))
		if err := elasticsearchService.TriggerNetwork(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("enable_public")
	}

	if d.Get("enable_public").(bool) == true && d.HasChange("public_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(WORKER)
		content["whiteIpList"] = d.Get("public_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("public_whitelist")
	}

	if d.HasChange("enable_kibana_public_network") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(KIBANA)
		content["actionType"] = elasticsearchService.getActionType(d.Get("enable_kibana_public_network").(bool))
		if err := elasticsearchService.TriggerNetwork(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("enable_kibana_public_network")
	}

	if d.Get("enable_kibana_public_network").(bool) == true && d.HasChange("kibana_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(KIBANA)
		content["whiteIpList"] = d.Get("kibana_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("kibana_whitelist")
	}

	if d.HasChange("enable_kibana_private_network") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(KIBANA)
		content["actionType"] = elasticsearchService.getActionType(d.Get("enable_kibana_private_network").(bool))
		if err := elasticsearchService.TriggerNetwork(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("enable_kibana_private_network")
	}

	if d.Get("enable_kibana_private_network").(bool) == true && d.HasChange("kibana_private_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(KIBANA)
		content["whiteIpList"] = d.Get("kibana_private_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("kibana_private_whitelist")
	}

	if d.HasChange("tags") {
		if err := updateInstanceTags(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudElasticsearchRead(d, meta)
	}

	if d.HasChange("instance_charge_type") {
		if err := updateInstanceChargeType(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("instance_charge_type")
	}

	if d.Get("instance_charge_type").(string) == string(PrePaid) && d.HasChange("period") {
		if err := renewInstance(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("period")
	}

	if d.HasChange("data_node_amount") {

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if err := updateDataNodeAmount(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("data_node_amount")
	}

	if d.HasChange("data_node_spec") || d.HasChange("data_node_disk_size") || d.HasChange("data_node_disk_type") {

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if err := updateDataNodeSpec(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("data_node_spec")
		d.SetPartial("data_node_disk_size")
		d.SetPartial("data_node_disk_type")
		d.SetPartial("data_node_disk_encrypted")
	}

	if d.HasChange("master_node_spec") {

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if err := updateMasterNode(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("master_node_spec")
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if err := updatePassword(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudElasticsearchRead(d, meta)
}

func resourceAlicloudElasticsearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	if strings.ToLower(d.Get("instance_charge_type").(string)) == strings.ToLower(string(PrePaid)) {
		return WrapError(Error("At present, 'PrePaid' instance cannot be deleted and must wait it to be expired and release it automatically"))
	}

	request := elasticsearch.CreateDeleteInstanceRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"InstanceActivating", "TokenPreviousRequestProcessError"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.DeleteInstance(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating", "inactive", "active"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{}))
	stateConf.PollInterval = 5 * time.Second

	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// Instance will be completed deleted in 5 minutes, so deleting vswitch is available after the time.
	time.Sleep(5 * time.Minute)

	return nil
}

func buildElasticsearchCreateRequest(d *schema.ResourceData, meta interface{}) (*elasticsearch.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := elasticsearch.CreateCreateInstanceRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	vpcService := VpcService{client}

	content := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		content["resourceGroupId"] = v.(string)
	}

	content["paymentType"] = strings.ToLower(d.Get("instance_charge_type").(string))
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		paymentInfo := make(map[string]interface{})
		if d.Get("period").(int) >= 12 {
			paymentInfo["duration"] = d.Get("period").(int) / 12
			paymentInfo["pricingCycle"] = string(Year)
		} else {
			paymentInfo["duration"] = d.Get("period").(int)
			paymentInfo["pricingCycle"] = string(Month)
		}

		content["paymentInfo"] = paymentInfo
	}

	content["nodeAmount"] = d.Get("data_node_amount")
	content["esVersion"] = d.Get("version")

	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return nil, WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		content["esAdminPassword"] = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return request, WrapError(err)
		}
		content["esAdminPassword"] = decryptResp.Plaintext
	}

	// Data node configuration
	dataNodeSpec := make(map[string]interface{})
	dataNodeSpec["spec"] = d.Get("data_node_spec")
	dataNodeSpec["disk"] = d.Get("data_node_disk_size")
	dataNodeSpec["diskType"] = d.Get("data_node_disk_type")
	dataNodeSpec["diskEncryption"] = d.Get("data_node_disk_encrypted")
	content["nodeSpec"] = dataNodeSpec

	// Master node configuration
	if d.Get("master_node_spec") != nil && d.Get("master_node_spec") != "" {
		masterNode := make(map[string]interface{})
		masterNode["spec"] = d.Get("master_node_spec")
		masterNode["amount"] = "3"
		masterNode["disk"] = "20"
		masterNode["diskType"] = "cloud_ssd"
		content["advancedDedicateMaster"] = true
		content["masterConfiguration"] = masterNode
	}

	// Network configuration
	vswitchId := d.Get("vswitch_id")
	vsw, err := vpcService.DescribeVSwitch(vswitchId.(string))
	if err != nil {
		return nil, WrapError(err)
	}

	network := make(map[string]interface{})
	network["type"] = "vpc"
	network["vpcId"] = vsw.VpcId
	network["vswitchId"] = vswitchId
	network["vsArea"] = vsw.ZoneId

	content["networkConfig"] = network

	if d.Get("zone_count") != nil && d.Get("zone_count") != "" {
		content["zoneCount"] = d.Get("zone_count")
	}

	data, err := json.Marshal(content)
	if err != nil {
		return nil, WrapError(err)
	}
	request.SetContent(data)
	request.SetContentType("application/json")

	return request, nil
}
