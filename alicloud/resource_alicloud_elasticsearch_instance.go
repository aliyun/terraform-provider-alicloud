package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if reg := regexp.MustCompile(`^[\w\-.]{0,30}$`); !reg.MatchString(value) {
						errors = append(errors, fmt.Errorf("%q be 0 to 30 characters in length and can contain numbers, letters, underscores, (_) and hyphens (-). It must start with a letter, a number or Chinese character.", k))
					}

					return
				},
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

			// Life cycle
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceChargeType,
				ForceNew:     true,
				Default:      PostPaid,
				Optional:     true,
			},

			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},

			// Data node configuration
			"data_node_amount": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(2, 50),
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

			"private_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"public_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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

			"kibana_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"zone_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(1, 3),
				Default:      1,
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

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
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
	d.Set("version", object.Result.EsVersion)
	d.Set("instance_charge_type", getChargeType(object.Result.PaymentType))

	d.Set("domain", object.Result.Domain)
	d.Set("port", object.Result.Port)

	// Kibana configuration
	d.Set("kibana_domain", object.Result.KibanaDomain)
	d.Set("kibana_port", object.Result.KibanaPort)
	d.Set("kibana_whitelist", filterWhitelist(object.Result.KibanaIPWhitelist, d.Get("kibana_whitelist").(*schema.Set)))

	// Data node configuration
	d.Set("data_node_amount", object.Result.NodeAmount)
	d.Set("data_node_spec", object.Result.NodeSpec.Spec)
	d.Set("data_node_disk_size", object.Result.NodeSpec.Disk)
	d.Set("data_node_disk_type", object.Result.NodeSpec.DiskType)
	d.Set("master_node_spec", object.Result.MasterConfiguration.Spec)

	// Cross zone configuration
	d.Set("zone_count", object.Result.ZoneCount)

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
		if err := updatePrivateWhitelist(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("private_whitelist")
	}

	if d.HasChange("public_whitelist") {
		if err := updatePublicWhitelist(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("public_whitelist")
	}

	if d.HasChange("kibana_whitelist") {
		if err := updateKibanaWhitelist(d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("kibana_whitelist")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudElasticsearchRead(d, meta)
	}

	if d.HasChange("data_node_amount") {

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		if err := updateDateNodeAmount(d, meta); err != nil {
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DeleteInstance(request)
		})
		if err != nil {
			if IsExceptedError(err, InstanceActivating) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, ESInstanceNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

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

	content["paymentType"] = strings.ToLower(d.Get("instance_charge_type").(string))
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		paymentInfo := make(map[string]interface{})
		if d.Get("period").(int) >= 12 {
			paymentInfo["duration"] = d.Get("period").(int) / 12
			paymentInfo["pricingCycle"] = strings.ToLower(string(Year))
		} else {
			paymentInfo["duration"] = d.Get("period").(int)
			paymentInfo["pricingCycle"] = strings.ToLower(string(Month))
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
	content["nodeSpec"] = dataNodeSpec

	// Master node configuration
	if d.Get("master_node_spec") != nil && d.Get("master_node_spec") != "" {
		masterNode := make(map[string]interface{})
		masterNode["spec"] = d.Get("master_node_spec")
		masterNode["amount"] = MasterNodeAmount
		masterNode["disk"] = MasterNodeDisk
		masterNode["diskType"] = MasterNodeDiskType
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
