package alicloud

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	aliyungoecs "github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesNodePoolCreate,
		Read:   resourceAlicloudCSNodePoolRead,
		Update: resourceAlicloudCSNodePoolUpdate,
		Delete: resourceAlicloudCSNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
			},
			"instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"key_name", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"password", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"labels": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"taints": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"management": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_repair": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"auto_upgrade": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"surge": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"surge_percentage": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"max_unavailable": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"scaling_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"max_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"cpu", "gpu", "gpushare", "spot"}, false),
						},
						"is_bond_eip": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"eip_internet_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
						},
						"eip_bandwidth": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 500),
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCSKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var raw interface{}

	clusterId := d.Get("cluster_id").(string)
	// prepare args and set default value
	args, err := buildNodePoolArgs(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "PrepareKubernetesNodePoolArgs", err)
	}

	if err = invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.CreateNodePool(args, d.Get("cluster_id").(string))
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "CreateKubernetesNodePool", raw)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = args
		addDebug("CreateKubernetesNodePool", raw, requestInfo, requestMap)
	}

	nodePool, ok := raw.(*cs.CreateNodePoolResponse)
	if ok != true {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes_node_pool", "ParseKubernetesNodePoolResponse", raw)
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, nodePool.NodePoolID))

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial", "scaling"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSNodePoolRead(d, meta)
}

func resourceAlicloudCSNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	vpcService := VpcService{client}
	d.Partial(true)
	update := false
	invoker := NewInvoker()
	args := &cs.UpdateNodePoolRequest{
		RegionId:         common.Region(client.RegionId),
		NodePoolInfo:     cs.NodePoolInfo{},
		ScalingGroup:     cs.ScalingGroup{},
		KubernetesConfig: cs.KubernetesConfig{},
		AutoScaling:      cs.AutoScaling{},
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("node_count") {
		oldV, newV := d.GetChange("node_count")

		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("node_count new value can not be parsed"), "parseError %d", newValue)
		}

		if newValue < oldValue {
			_ = removeNodePoolNodes(d, meta, parts)
			// The removal of a node is logically independent.
			// The removal of a node should not involve parameter changes.
			return resourceAlicloudCSNodePoolRead(d, meta)
		}
		update = true
		args.Count = int64(newValue) - int64(oldValue)
	}
	args.NodePoolInfo.Name = d.Get("name").(string)
	if d.HasChange("name") {
		update = true
		args.NodePoolInfo.Name = d.Get("name").(string)
	}
	if d.HasChange("vswitch_ids") {
		update = true
		var vswitchID string
		if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
			vswitchID = list[0]
		} else {
			vswitchID = ""
		}

		var vpcId string
		if vswitchID != "" {
			vsw, err := vpcService.DescribeVSwitch(vswitchID)
			if err != nil {
				return err
			}
			vpcId = vsw.VpcId
		}
		args.ScalingGroup.VpcId = vpcId
		args.ScalingGroup.VswitchIds = expandStringList(d.Get("vswitch_ids").([]interface{}))
	}
	if d.HasChange("instance_types") {
		update = true
		args.ScalingGroup.InstanceTypes = expandStringList(d.Get("instance_types").([]interface{}))
	}

	// password is required by update method
	args.ScalingGroup.LoginPassword = d.Get("password").(string)
	if d.HasChange("password") {
		update = true
		args.ScalingGroup.LoginPassword = d.Get("password").(string)
	}

	args.ScalingGroup.KeyPair = d.Get("key_name").(string)
	if d.HasChange("key_name") {
		update = true
		args.ScalingGroup.KeyPair = d.Get("key_name").(string)
	}

	if d.HasChange("security_group_id") {
		update = true
		args.ScalingGroup.SecurityGroupId = d.Get("security_group_id").(string)
	}

	if d.HasChange("system_disk_category") {
		update = true
		args.ScalingGroup.SystemDiskCategory = aliyungoecs.DiskCategory(d.Get("system_disk_category").(string))
	}

	if d.HasChange("system_disk_size") {
		update = true
		args.ScalingGroup.SystemDiskSize = int64(d.Get("system_disk_size").(int))
	}

	if d.HasChange("image_id") {
		update = true
		args.ScalingGroup.ImageId = d.Get("image_id").(string)
	}

	if d.HasChange("data_disks") {
		update = true
		setNodePoolDataDisks(&args.ScalingGroup, d)
	}

	if d.HasChange("tags") {
		update = true
		setNodePoolTags(&args.ScalingGroup, d)
	}

	if d.HasChange("labels") {
		update = true
		setNodePoolLabels(&args.KubernetesConfig, d)
	}

	if d.HasChange("taints") {
		update = true
		setNodePoolTaints(&args.KubernetesConfig, d)
	}

	if d.HasChange("node_name_mode") {
		update = true
		args.KubernetesConfig.NodeNameMode = d.Get("node_name_mode").(string)
	}

	if d.HasChange("user_data") {
		update = true
		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				args.KubernetesConfig.UserData = v
			} else {
				args.KubernetesConfig.UserData = base64.StdEncoding.EncodeToString([]byte(v))
			}
		}
	}

	if d.HasChange("scaling_config") {
		update = true
		sc := d.Get("scaling_config").([]interface{})
		args.AutoScaling = setAutoScalingConfig(sc)
	}

	if v, ok := d.Get("management").([]interface{}); len(v) > 0 && ok {
		args.Management = setManagedNodepoolConfig(v)
	}

	if update {

		var resoponse interface{}
		if err := invoker.Run(func() error {
			var err error
			resoponse, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				resp, err := csClient.UpdateNodePool(parts[0], parts[1], args)
				return resp, err
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateKubernetesNodePool", DenverdinoAliyungo)
		}
		if debugOn() {
			resizeRequestMap := make(map[string]interface{})
			resizeRequestMap["ClusterId"] = parts[0]
			resizeRequestMap["NodePoolId"] = parts[1]
			resizeRequestMap["Args"] = args
			addDebug("UpdateKubernetesNodePool", resoponse, resizeRequestMap)
		}

		stateConf := BuildStateConf([]string{"scaling", "updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	d.Partial(false)
	return resourceAlicloudCSNodePoolRead(d, meta)
}

func resourceAlicloudCSNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	object, err := csService.DescribeCsKubernetesNodePool(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("node_count", object.TotalNodes)
	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_ids", object.VswitchIds)
	d.Set("instance_types", object.InstanceTypes)
	d.Set("key_name", object.KeyPair)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("system_disk_category", object.SystemDiskCategory)
	d.Set("system_disk_size", object.SystemDiskSize)
	d.Set("image_id", object.ImageId)
	d.Set("node_name_mode", object.NodeNameMode)
	d.Set("user_data", object.UserData)
	d.Set("scaling_group_id", object.ScalingGroupId)

	if passwd, ok := d.GetOk("password"); ok && passwd.(string) != "" {
		d.Set("password", passwd)
	}

	if parts, err := ParseResourceId(d.Id(), 2); err != nil {
		return WrapError(err)
	} else {
		d.Set("cluster_id", string(parts[0]))
	}

	if err := d.Set("data_disks", flattenNodeDataDisksConfig(object.DataDisks)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("taints", flattenTaintsConfig(object.Taints)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("labels", flattenLabelsConfig(object.Labels)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}

	if m, ok := d.GetOk("management"); ok && m != nil {
		if err := d.Set("management", flattenManagementNodepoolConfig(&object.Management)); err != nil {
			return WrapError(err)
		}
	}

	if err := d.Set("scaling_config", flattenAutoScalingConfig(&object.AutoScaling)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudCSNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	// delete all nodes
	_ = removeNodePoolNodes(d, meta, parts)

	var response interface{}
	err = resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.DeleteNodePool(parts[0], parts[1])
			})
			response = raw
			return err
		}); err != nil {
			return resource.RetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = parts[0]
			requestMap["NodePoolId"] = parts[1]
			addDebug("DeleteClusterNodePool", response, d.Id(), requestMap)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNodePoolNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteClusterNodePool", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"active", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildNodePoolArgs(d *schema.ResourceData, meta interface{}) (*cs.CreateNodePoolRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}

	var vswitchID string
	if list := expandStringList(d.Get("vswitch_ids").([]interface{})); len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
	}

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp.Plaintext
		}
	}

	creationArgs := &cs.CreateNodePoolRequest{
		RegionId: common.Region(client.RegionId),
		Count:    int64(d.Get("node_count").(int)),
		NodePoolInfo: cs.NodePoolInfo{
			Name:         d.Get("name").(string),
			NodePoolType: "ess", // hard code the type
		},
		ScalingGroup: cs.ScalingGroup{
			VpcId:              vpcId,
			VswitchIds:         expandStringList(d.Get("vswitch_ids").([]interface{})),
			InstanceTypes:      expandStringList(d.Get("instance_types").([]interface{})),
			LoginPassword:      password,
			KeyPair:            d.Get("key_name").(string),
			SystemDiskCategory: aliyungoecs.DiskCategory(d.Get("system_disk_category").(string)),
			SystemDiskSize:     int64(d.Get("system_disk_size").(int)),
			SecurityGroupId:    d.Get("security_group_id").(string),
			ImageId:            d.Get("image_id").(string),
		},
		KubernetesConfig: cs.KubernetesConfig{
			NodeNameMode: d.Get("node_name_mode").(string),
		},
	}

	setNodePoolDataDisks(&creationArgs.ScalingGroup, d)
	setNodePoolTags(&creationArgs.ScalingGroup, d)
	setNodePoolTaints(&creationArgs.KubernetesConfig, d)
	setNodePoolLabels(&creationArgs.KubernetesConfig, d)

	if v, ok := d.GetOk("user_data"); ok && v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			creationArgs.KubernetesConfig.UserData = v.(string)
		} else {
			creationArgs.KubernetesConfig.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	// set auto scaling config
	if v, ok := d.GetOk("scaling_config"); ok {
		if sc, ok := v.([]interface{}); len(sc) > 0 && ok {
			creationArgs.AutoScaling = setAutoScalingConfig(sc)
		}
	}

	// set manage nodepool params
	if v, ok := d.GetOk("management"); ok {
		if management, ok := v.([]interface{}); len(management) > 0 && ok {
			creationArgs.Management = setManagedNodepoolConfig(management)
		}
	}

	return creationArgs, nil
}

func ConvertCsTags(d *schema.ResourceData) ([]cs.Tag, error) {
	tags := make([]cs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, cs.Tag{
						Key:   key,
						Value: v,
					})
				}
			}
		}
	}

	return tags, nil
}

func setNodePoolTags(scalingGroup *cs.ScalingGroup, d *schema.ResourceData) error {
	if _, ok := d.GetOk("tags"); ok {
		if tags, err := ConvertCsTags(d); err == nil {
			scalingGroup.Tags = tags
		}
	}

	return nil
}

func setNodePoolLabels(config *cs.KubernetesConfig, d *schema.ResourceData) error {
	if v, ok := d.GetOk("labels"); ok && len(v.([]interface{})) > 0 {
		vl := v.([]interface{})
		labels := make([]cs.Label, 0)
		for _, i := range vl {
			if m, ok := i.(map[string]interface{}); ok {
				labels = append(labels, cs.Label{
					Key:   m["key"].(string),
					Value: m["value"].(string),
				})
			}

		}
		config.Labels = labels
	}

	return nil
}

func setNodePoolDataDisks(scalingGroup *cs.ScalingGroup, d *schema.ResourceData) error {
	if dds, ok := d.GetOk("data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]cs.NodePoolDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := cs.NodePoolDataDisk{
				Size:                 pack["size"].(int),
				DiskName:             pack["name"].(string),
				Category:             pack["category"].(string),
				Device:               pack["device"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				Encrypted:            pack["encrypted"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		scalingGroup.DataDisks = createDataDisks
	}

	return nil
}

func setNodePoolTaints(config *cs.KubernetesConfig, d *schema.ResourceData) error {
	if v, ok := d.GetOk("taints"); ok && len(v.([]interface{})) > 0 {
		vl := v.([]interface{})
		taints := make([]cs.Taint, 0)
		for _, i := range vl {
			if m, ok := i.(map[string]interface{}); ok {
				taints = append(taints, cs.Taint{
					Key:    m["key"].(string),
					Value:  m["value"].(string),
					Effect: cs.Effect(m["effect"].(string)),
				})
			}

		}
		config.Taints = taints
	}

	return nil
}

func setManagedNodepoolConfig(l []interface{}) (config cs.Management) {
	if len(l) == 0 || l[0] == nil {
		return config
	}

	m := l[0].(map[string]interface{})

	// Once "management" is set, we think of it as creating a managed node pool
	config.Enable = true

	if v, ok := m["auto_repair"].(bool); ok {
		config.AutoRepair = v
	}
	if v, ok := m["auto_upgrade"].(bool); ok {
		config.UpgradeConf.AutoUpgrade = v
	}
	if v, ok := m["surge"].(int); ok {
		config.UpgradeConf.Surge = int64(v)
	}
	if v, ok := m["surge_percentage"].(int); ok {
		config.UpgradeConf.SurgePercentage = int64(v)
	}
	if v, ok := m["max_unavailable"].(int); ok {
		config.UpgradeConf.MaxUnavailable = int64(v)
	}

	return config
}

func setAutoScalingConfig(l []interface{}) (config cs.AutoScaling) {
	if len(l) == 0 || l[0] == nil {
		return config
	}

	m := l[0].(map[string]interface{})

	// Once "scaling_config" is set, we think of it as creating a auto scaling node pool
	config.Enable = true

	if v, ok := m["min_size"].(int); ok {
		config.MinInstances = int64(v)
	}
	if v, ok := m["max_size"].(int); ok {
		config.MaxInstances = int64(v)
	}
	if v, ok := m["type"].(string); ok {
		config.Type = v
	}
	if v, ok := m["is_bond_eip"].(bool); ok {
		config.IsBindEip = &v
	}
	if v, ok := m["eip_internet_charge_type"].(string); ok {
		config.EipInternetChargeType = v
	}
	if v, ok := m["eip_bandwidth"].(int); ok {
		config.EipBandWidth = int64(v)
	}

	return config
}

func flattenAutoScalingConfig(config *cs.AutoScaling) (m []map[string]interface{}) {
	if config == nil {
		return
	}
	m = append(m, map[string]interface{}{
		"min_size":                 config.MinInstances,
		"max_size":                 config.MaxInstances,
		"type":                     config.Type,
		"is_bond_eip":              config.IsBindEip,
		"eip_internet_charge_type": config.EipInternetChargeType,
		"eip_bandwidth":            config.EipBandWidth,
	})

	return
}

func flattenManagementNodepoolConfig(config *cs.Management) (m []map[string]interface{}) {
	if config == nil {
		return
	}
	m = append(m, map[string]interface{}{
		"auto_repair":      config.AutoRepair,
		"auto_upgrade":     config.UpgradeConf.AutoUpgrade,
		"surge":            config.UpgradeConf.Surge,
		"surge_percentage": config.UpgradeConf.SurgePercentage,
		"max_unavailable":  config.UpgradeConf.MaxUnavailable,
	})

	return
}

func flattenNodeDataDisksConfig(config []cs.NodePoolDataDisk) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, disks := range config {
		m = append(m, map[string]interface{}{
			"size":      disks.Size,
			"category":  disks.Category,
			"encrypted": disks.Encrypted,
		})
	}

	return m
}

func flattenTaintsConfig(config []cs.Taint) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, taint := range config {
		m = append(m, map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		})
	}

	return m
}

func flattenLabelsConfig(config []cs.Label) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, label := range config {
		m = append(m, map[string]interface{}{
			"key":   label.Key,
			"value": label.Value,
		})
	}

	return m
}

func flattenTagsConfig(config []cs.Tag) map[string]string {
	m := make(map[string]string, len(config))
	if len(config) < 0 {
		return m
	}

	for _, tag := range config {
		m[tag.Key] = tag.Value
	}

	return m
}

func removeNodePoolNodes(d *schema.ResourceData, meta interface{}, parseId []string) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var response interface{}
	// list all nodes of the nodepool
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			nodes, _, err := csClient.GetKubernetesClusterNodes(parseId[0], common.Pagination{PageNumber: 1, PageSize: PageSizeLarge}, parseId[1])
			return nodes, err
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
	}

	ret := response.([]cs.KubernetesNodeType)
	// filter out nodes
	var allNodeName []string
	for _, value := range ret {
		allNodeName = append(allNodeName, value.NodeName)
	}

	// remove nodes
	removeNodesName := allNodeName
	if d.HasChange("node_count") {
		o, n := d.GetChange("node_count")
		count := o.(int) - n.(int)
		removeNodesName = allNodeName[:count]
	}

	removeNodesArgs := &cs.DeleteKubernetesClusterNodesRequest{
		Nodes:       removeNodesName,
		ReleaseNode: true,
		DrainNode:   false,
	}
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			resp, err := csClient.DeleteKubernetesClusterNodes(parseId[0], removeNodesArgs)
			return resp, err
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteKubernetesClusterNodes", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetPartial("node_count")

	return nil
}
