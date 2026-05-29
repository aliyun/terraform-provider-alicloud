package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudPolarDBAICluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBAIClusterCreate,
		Read:   resourceAlicloudPolarDBAIClusterRead,
		Delete: resourceAlicloudPolarDBAIClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region ID of the AI cluster.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The zone ID of the AI cluster.",
			},
			"db_node_class": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The DB node class of the AI cluster.",
			},
			"db_cluster_description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the AI DB cluster.",
			},
			"pay_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Postpaid", "Prepaid"}, false),
				Description:  "The billing method. Valid values: `Postpaid`, `Prepaid`.",
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable auto-renewal.",
			},
			"period": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The subscription period.",
			},
			"used_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The subscription duration.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC ID.",
			},
			"vswitch_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vSwitch ID.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The security group ID.",
			},
			"kube_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ainode"}, false),
				Description:  "The type of the Kubernetes cluster.",
			},
			"model_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The model name. Example: `Qwen3.5-9B`.",
			},
			"extension": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"maas", "custom"}, false),
				Description:  "The extension type. Valid values: `maas`, `custom`.",
			},
			"inference_engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"sglang", "vllm"}, false),
				Description:  "The inference engine. Valid values: `sglang`, `vllm`.",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the associated DB cluster.",
			},
			"auto_use_coupon": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Whether to automatically use coupons. Default: `true`.",
			},
			"promotion_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The promotion code. If not specified, the default coupon is used.",
			},
			// Computed attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the AI cluster.",
			},
			"model_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The model type. Example: `public`.",
			},
		},
	}
}

func resourceAlicloudPolarDBAIClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	if d.Get("pay_type").(string) == "Prepaid" {
		if _, ok := d.GetOk("period"); !ok {
			return WrapError(fmt.Errorf("'period' is required when 'pay_type' is 'Prepaid'"))
		}
		if _, ok := d.GetOk("used_time"); !ok {
			return WrapError(fmt.Errorf("'used_time' is required when 'pay_type' is 'Prepaid'"))
		}
	}

	action := "CreateAIDBCluster"
	request := map[string]interface{}{
		"RegionId":    d.Get("region_id").(string),
		"DBNodeClass": d.Get("db_node_class").(string),
		"PayType":     d.Get("pay_type").(string),
		"VPCId":       d.Get("vpc_id").(string),
		"VSwitchId":   d.Get("vswitch_id").(string),
		"ClientToken": buildClientToken("CreateAIDBCluster"),
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v.(string)
	}
	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v.(string)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		request["AutoRenew"] = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v.(string)
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v.(string)
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v.(string)
	}
	if v, ok := d.GetOk("kube_type"); ok {
		request["KubeType"] = v.(string)
	}
	if v, ok := d.GetOk("model_name"); ok {
		request["ModelName"] = v.(string)
	}
	if v, ok := d.GetOk("extension"); ok {
		request["Extension"] = v.(string)
	}
	if v, ok := d.GetOk("inference_engine"); ok {
		request["InferenceEngine"] = v.(string)
	}
	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v.(string)
	}
	request["AutoUseCoupon"] = d.Get("auto_use_coupon").(bool)
	if v, ok := d.GetOk("promotion_code"); ok {
		request["PromotionCode"] = v.(string)
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_aicluster", action, AlibabaCloudSdkGoERROR)
	}

	aiClusterId, ok := response["DBClusterId"].(string)
	if !ok || aiClusterId == "" {
		return WrapError(fmt.Errorf("CreateAIDBCluster returned empty DBClusterId"))
	}
	d.SetId(aiClusterId)

	// Wait for the AI cluster to become activated
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, polarDBService.PolarDBAIClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPolarDBAIClusterRead(d, meta)
}

func resourceAlicloudPolarDBAIClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	attribute, err := polarDBService.DescribePolarDBAIClusterAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if v, ok := attribute["DBClusterStatus"].(string); ok {
		d.Set("status", v)
	}
	if v, ok := attribute["RegionId"].(string); ok {
		d.Set("region_id", v)
	}
	if v, ok := attribute["ZoneId"].(string); ok {
		d.Set("zone_id", v)
	}
	if v, ok := attribute["DBClusterDescription"].(string); ok {
		d.Set("db_cluster_description", v)
	}
	if v, ok := attribute["PayType"].(string); ok {
		d.Set("pay_type", v)
	}
	if v, ok := attribute["VPCId"].(string); ok {
		d.Set("vpc_id", v)
	}
	if v, ok := attribute["VSwitchId"].(string); ok {
		d.Set("vswitch_id", v)
	}
	if v, ok := attribute["RunType"].(string); ok {
		d.Set("kube_type", v)
	}
	if v, ok := attribute["ModelName"].(string); ok {
		d.Set("model_name", v)
	}
	if v, ok := attribute["AiNodeType"].(string); ok {
		d.Set("extension", v)
	}
	if v, ok := attribute["EcsSecurityGroupId"].(string); ok {
		d.Set("security_group_id", v)
	}
	if v, ok := attribute["ModelType"].(string); ok {
		d.Set("model_type", v)
	}
	// Extract db_node_class from DBNodes array
	if dbNodes, ok := attribute["DBNodes"].([]interface{}); ok && len(dbNodes) > 0 {
		if node, ok := dbNodes[0].(map[string]interface{}); ok {
			if nodeClass, ok := node["DBNodeClass"].(string); ok {
				d.Set("db_node_class", nodeClass)
			}
		}
	}
	return nil
}

func resourceAlicloudPolarDBAIClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	action := "DeleteAIDBCluster"
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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

	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, polarDBService.PolarDBAIClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		if strings.HasPrefix(err.Error(), "couldn't find resource") {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.SetId("")
	return nil
}
