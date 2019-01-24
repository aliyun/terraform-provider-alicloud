package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform/helper/resource"
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

		Schema: map[string]*schema.Schema{
			// Basic instance information
			"description": &schema.Schema{
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

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},

			"version": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: esVersionDiffSuppressFunc,
			},

			// Life cycle
			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateInstanceChargeType,
				ForceNew:     true,
				Default:      PostPaid,
				Optional:     true,
			},

			"period": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},

			// Data node configuration
			"data_node_amount": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(2, 50),
			},

			"data_node_spec": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"data_node_disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"data_node_disk_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"private_whitelist": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"public_whitelist": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"master_node_spec": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// Kibana node configuration
			"kibana_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"kibana_port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"kibana_whitelist": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		return BuildWrapError("CreateInstance", "", ProviderERROR, err, "")
	}

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.CreateInstance(request)
	})

	if err != nil {
		return BuildWrapError(request.GetActionName(), "", ProviderERROR, err, "")
	}

	resp, _ := raw.(*elasticsearch.CreateInstanceResponse)
	d.SetId(resp.Result.InstanceId)

	if err := elasticsearchService.WaitForElasticsearchInstance(resp.Result.InstanceId, []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout); err != nil {
		return BuildWrapError(request.GetActionName(), resp.Result.InstanceId, ProviderERROR, err, "")
	}

	return resourceAlicloudElasticsearchUpdate(d, meta)
}

func resourceAlicloudElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	resp, err := elasticsearchService.DescribeInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return BuildWrapError("DescribeInstance ", d.Id(), AlibabaCloudSdkGoERROR, err, "")
	}

	d.Set("description", resp.Result.Description)
	d.Set("status", resp.Result.Status)
	d.Set("vswitch_id", resp.Result.NetworkConfig.VswitchId)

	d.Set("private_whitelist", resp.Result.EsIPWhitelist)
	d.Set("public_whitelist", resp.Result.PublicIpWhitelist)
	d.Set("version", resp.Result.EsVersion)
	d.Set("instance_charge_type", getChargeType(resp.Result.PaymentType))

	d.Set("domain", resp.Result.Domain)
	d.Set("port", resp.Result.Port)

	// Kibana configuration
	d.Set("kibana_domain", resp.Result.KibanaDomain)
	d.Set("kibana_port", resp.Result.KibanaPort)
	d.Set("kibana_whitelist", resp.Result.KibanaIPWhitelist)

	// Data node configuration
	d.Set("data_node_amount", resp.Result.NodeAmount)
	d.Set("data_node_spec", resp.Result.NodeSpec.Spec)
	d.Set("data_node_disk_size", resp.Result.NodeSpec.Disk)
	d.Set("data_node_disk_type", resp.Result.NodeSpec.DiskType)
	d.Set("master_node_spec", resp.Result.MasterConfiguration.Spec)

	return nil
}

func resourceAlicloudElasticsearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	d.Partial(true)

	if d.HasChange("description") {
		if err := updateDescription(d, meta); err != nil {
			return err
		}

		d.SetPartial("description")
	}

	if d.HasChange("private_whitelist") {
		if err := updatePrivateWhitelist(d, meta); err != nil {
			return err
		}

		d.SetPartial("private_whitelist")
	}

	if d.HasChange("public_whitelist") {
		if err := updatePublicWhitelist(d, meta); err != nil {
			return err
		}

		d.SetPartial("public_whitelist")
	}

	if d.HasChange("kibana_whitelist") {
		if err := updateKibanaWhitelist(d, meta); err != nil {
			return err
		}

		d.SetPartial("kibana_whitelist")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudElasticsearchRead(d, meta)
	}

	if d.HasChange("data_node_amount") {

		if err := elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout); err != nil {
			return err
		}

		if err := updateDateNodeAmount(d, meta); err != nil {
			return err
		}

		d.SetPartial("data_node_amount")
	}

	if d.HasChange("data_node_spec") || d.HasChange("data_node_disk_size") || d.HasChange("data_node_disk_type") {

		if err := elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout); err != nil {
			return err
		}

		if err := updateDataNodeSpec(d, meta); err != nil {
			return err
		}

		d.SetPartial("data_node_spec")
		d.SetPartial("data_node_disk_size")
		d.SetPartial("data_node_disk_type")
	}

	if d.HasChange("master_node_spec") {

		if err := elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout); err != nil {
			return BuildWrapError("WaitInstanceStatus", d.Id(), ProviderERROR, err, "")
		}

		if err := updateMasterNode(d, meta); err != nil {
			return err
		}

		d.SetPartial("master_node_spec")
	}

	if d.HasChange("password") {

		if err := elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout); err != nil {
			return BuildWrapError("WaitInstanceStatus", d.Id(), ProviderERROR, err, "")
		}

		if err := updatePassword(d, meta); err != nil {
			return err
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
		return BuildWrapError(
			"DeleteInstance",
			d.Id(),
			ProviderERROR,
			fmt.Errorf("At present, 'PrePaid' instance cannot be deleted and must wait it to be expired and release it automatically"),
			"",
		)
	}

	request := elasticsearch.CreateDeleteInstanceRequest()
	request.InstanceId = d.Id()
	request.SetContentType("application/json")

	if err := resource.Retry(2*time.Hour, func() *resource.RetryError {
		if _, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DeleteInstance(request)
		}); err != nil {
			if IsExceptedError(err, ESInstanceNotFound) {
				return nil
			}

			return resource.RetryableError(BuildWrapError("DeleteInstance", d.Id(), ProviderERROR, err, ""))
		}

		if _, err := elasticsearchService.DescribeInstance(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
		}

		return nil
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), ProviderERROR, err, "")
	}

	return nil
}

func buildElasticsearchCreateRequest(d *schema.ResourceData, meta interface{}) (*elasticsearch.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := elasticsearch.CreateCreateInstanceRequest()
	vpcService := VpcService{client}

	content := make(map[string]interface{})

	if v := d.Get("description").(string); v != "" {
		content["description"] = v
	}

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
	content["esAdminPassword"] = d.Get("password")

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
	vsw, err := vpcService.DescribeVswitch(vswitchId.(string))
	if err != nil {
		return nil, BuildWrapError("DescribeVSwitch", vswitchId.(string), ProviderERROR, err, "")
	}

	network := make(map[string]interface{})
	network["type"] = "vpc"
	network["vpcId"] = vsw.VpcId
	network["vswitchId"] = vswitchId
	network["vsArea"] = vsw.ZoneId

	content["networkConfig"] = network

	data, err := json.Marshal(content)
	request.SetContent(data)
	request.SetContentType("application/json")

	return request, nil
}
