package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudOtsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceCreate,
		Read:   resourceAliyunOtsInstanceRead,
		Update: resourceAliyunOtsInstanceUpdate,
		Delete: resourceAliyunOtsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},

			// Expired
			"accessed_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(AnyNetwork), string(VpcOnly), string(VpcOrConsole),
				}, false),
			},

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_type_acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"network_source_acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  OtsHighPerformance,
				ValidateFunc: validation.StringInSlice([]string{
					string(OtsCapacity), string(OtsHighPerformance),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			// lintignore: S006
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"vcu": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"period_in_month": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"enable_auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_elastic_vcu": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"auto_renew_period_in_month": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"elastic_vcu_upper_limit": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	instanceName := d.Get("name").(string)

	// VCU instance (reserved mode): created via CreateVCUInstance when vcu > 0.
	if vcu := d.Get("vcu").(int); vcu > 0 {
		actionPath := "/v2/openapi/createvcuinstance"
		request := buildCreateVCUInstanceRequest(d)
		if _, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request); err != nil {
			return WrapError(err)
		}
		d.SetId(instanceName)
		if err := otsService.WaitForOtsInstance(instanceName, toInstanceInnerStatus(Running), DefaultTimeout); err != nil {
			return WrapError(err)
		}
		return resourceAliyunOtsInstanceUpdate(d, meta)
	}

	instanceTypeStr := d.Get("instance_type").(string)
	instanceType, err := parseAndCheckInstanceType(instanceTypeStr, otsService)
	if err != nil {
		return WrapError(err)
	}

	actionPath, _, request := buildCreateInstanceRoaRequest(d, client.RegionId, instanceType)

	if _, err = OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request); err != nil {
		return WrapError(err)
	}

	d.SetId(instanceName)
	if err := otsService.WaitForOtsInstance(instanceName, toInstanceInnerStatus(Running), DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunOtsInstanceUpdate(d, meta)
}

func parseAndCheckInstanceType(instanceTypeStr string, otsService OtsService) (string, error) {
	instanceType := convertInstanceType(OtsInstanceType(instanceTypeStr))
	types, err := otsService.DescribeOtsInstanceTypes()
	if err != nil {
		return "", WrapError(err)
	}
	valid := false
	for _, t := range types {
		if instanceType == t {
			valid = true
			break
		}
	}
	if valid {
		return instanceType, nil
	}
	return instanceType, WrapError(Error("The instance type %s is not available in the region %s.", instanceTypeStr, otsService.client.RegionId))

}

func buildCreateInstanceRoaRequest(d *schema.ResourceData, regionId string, instanceType string) (string, string, map[string]interface{}) {
	actionPath := "/v2/openapi/createinstance"
	request := make(map[string]interface{})
	request["RegionId"] = StringPointer(regionId)
	request["ClusterType"] = StringPointer(instanceType)
	instanceName := d.Get("name").(string)
	request["InstanceName"] = StringPointer(instanceName)
	request["ResourceGroupId"] = StringPointer(d.Get("resource_group_id").(string))
	request["InstanceDescription"] = StringPointer(d.Get("description").(string))

	hasSetNetwork := false
	if v, ok := d.GetOk("accessed_by"); ok {
		hasSetNetwork = true
		request["Network"] = StringPointer(convertInstanceAccessedBy(InstanceAccessedByType(v.(string))))
	}

	hasSetACL := false
	// LIST or SET cannot set default values in schema in latest terraform version, so do it manually
	// terraform cannot handle nil and[] in list/set: https://github.com/hashicorp/terraform-plugin-sdk/issues/142
	// in terraform the zero value of list/set is [], in golang the zero value of slice is nil
	netTypeList := []string{string(VpcAccess), string(ClassicAccess), string(InternetAccess)}
	// v not nil and [], it will be ok
	if v, ok := d.GetOk("network_type_acl"); ok {
		hasSetACL = true
		netTypeList = expandStringList(v.(*schema.Set).List())
	}
	request["NetworkTypeACL"] = netTypeList

	netSourceList := []string{string(TrustProxyAccess)}
	if v, ok := d.GetOk("network_source_acl"); ok {
		hasSetACL = true
		netSourceList = expandStringList(v.(*schema.Set).List())
	}
	request["NetworkSourceACL"] = netSourceList

	// In order to maintain compatibility, when the Network attribute is set,
	// the ACL attribute cannot have a default value.
	if hasSetNetwork && !hasSetACL {
		request["NetworkTypeACL"] = nil
		request["NetworkSourceACL"] = nil
	}

	if tagMap, ok := d.GetOk("tags"); ok {
		var tags []map[string]string

		for key, value := range tagMap.(map[string]interface{}) {
			tag := map[string]string{"Key": key, "Value": value.(string)}
			tags = append(tags, tag)
		}
		request["Tags"] = tags
	}

	return actionPath, instanceName, request
}

// buildCreateVCUInstanceRequest builds the request body for the CreateVCUInstance
// API. A VCU instance is a reserved-mode Tablestore instance; its ClusterType
// is fixed to "vcu" and the instance is billed by VCU quota.
func buildCreateVCUInstanceRequest(d *schema.ResourceData) map[string]interface{} {
	request := make(map[string]interface{})
	request["ClusterType"] = "vcu"
	request["VCU"] = d.Get("vcu").(int)
	request["PeriodInMonth"] = d.Get("period_in_month").(int)
	request["InstanceName"] = d.Get("name").(string)
	if v, ok := d.GetOk("alias_name"); ok {
		request["AliasName"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request["InstanceDescription"] = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v.(string)
	}
	if v, ok := d.GetOk("enable_auto_renew"); ok {
		request["EnableAutoRenew"] = v.(bool)
	}
	if v, ok := d.GetOk("enable_elastic_vcu"); ok {
		request["EnableElasticVCU"] = v.(bool)
	}
	if v, ok := d.GetOk("auto_renew_period_in_month"); ok {
		request["AutoRenewPeriodInMonth"] = v.(int)
	}
	if tagMap, ok := d.GetOk("tags"); ok {
		var tags []map[string]string
		for key, value := range tagMap.(map[string]interface{}) {
			tags = append(tags, map[string]string{"Key": key, "Value": value.(string)})
		}
		request["Tags"] = tags
	}
	return request
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	instance, err := otsService.DescribeOtsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", instance.InstanceName)
	err = d.Set("accessed_by", convertInstanceAccessedByRevert(instance.Network))
	if err != nil {
		return err
	}
	err = d.Set("resource_group_id", instance.ResourceGroupId)
	if err != nil {
		return err
	}
	err = d.Set("network_type_acl", instance.NetworkTypeACL)
	if err != nil {
		return err
	}
	err = d.Set("network_source_acl", instance.NetworkSourceACL)
	if err != nil {
		return err
	}

	err = d.Set("instance_type", convertInstanceTypeRevert(instance.InstanceSpecification))
	if err != nil {
		return err
	}
	err = d.Set("description", instance.InstanceDescription)
	if err != nil {
		return err
	}
	err = d.Set("tags", otsRestTagsToMap(instance.Tags))
	if err != nil {
		return err
	}
	err = d.Set("payment_type", instance.PaymentType)
	if err != nil {
		return err
	}
	err = d.Set("alias_name", instance.AliasName)
	if err != nil {
		return err
	}
	err = d.Set("elastic_vcu_upper_limit", instance.ElasticVCUUpperLimit)
	if err != nil {
		return err
	}
	err = d.Set("policy", instance.Policy)
	if err != nil {
		return err
	}
	err = d.Set("policy_version", instance.PolicyVersion)
	if err != nil {
		return err
	}
	// vcu is the purchased VCU quota for a reserved-mode instance. Reading
	// back from VCUQuota keeps import of an existing VCU instance in sync
	// with the server.
	err = d.Set("vcu", instance.VCUQuota)
	if err != nil {
		return err
	}
	return nil
}

func resourceAliyunOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		actionPath := "/v2/openapi/changeresourcegroup"
		request := make(map[string]interface{})
		request["ResourceId"] = StringPointer(d.Id())
		request["NewResourceGroupId"] = StringPointer(d.Get("resource_group_id").(string))

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		d.SetPartial("resource_group_id")
	}
	hasChangeACL := false
	if !d.IsNewResource() && d.HasChange("network_type_acl") {
		actionPath := "/v2/openapi/updateinstance"
		request := make(map[string]interface{})
		request["RegionId"] = StringPointer(client.RegionId)
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())

		netTypeList := expandStringList(d.Get("network_type_acl").(*schema.Set).List())
		request["NetworkTypeACL"] = netTypeList
		// acl must set together
		netSourceList := expandStringList(d.Get("network_source_acl").(*schema.Set).List())
		request["NetworkSourceACL"] = netSourceList

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		hasChangeACL = true
		d.SetPartial("network_type_acl")
	}

	if !d.IsNewResource() && d.HasChange("network_source_acl") {
		actionPath := "/v2/openapi/updateinstance"
		request := make(map[string]interface{})
		request["RegionId"] = StringPointer(client.RegionId)
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())

		// todo handle lens 0 []
		netSourceList := expandStringList(d.Get("network_source_acl").(*schema.Set).List())
		request["NetworkSourceACL"] = netSourceList
		// acl must set together
		netTypeList := expandStringList(d.Get("network_type_acl").(*schema.Set).List())
		request["NetworkTypeACL"] = netTypeList

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		hasChangeACL = true
		d.SetPartial("network_source_acl")
	}

	// accessed_by is a deprecated attribute, updates on accessed_by will only take effect when the ACL has not been updated.
	if !d.IsNewResource() && (d.HasChange("accessed_by") && !hasChangeACL) {
		actionPath := "/v2/openapi/updateinstance"
		request := make(map[string]interface{})
		request["RegionId"] = StringPointer(client.RegionId)
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())
		request["Network"] = StringPointer(convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string))))

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		d.SetPartial("accessed_by")
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			var removeKeys []string
			for _, t := range remove {
				removeKeys = append(removeKeys, t.Key)
			}

			actionPath := "/v2/openapi/untagresources"
			request := make(map[string]interface{})
			request["RegionId"] = StringPointer(client.RegionId)
			request["ResourceType"] = "INSTANCE"
			// id is instanceName
			request["ResourceIds"] = expandStringList([]interface{}{d.Id()})
			request["TagKeys"] = removeKeys

			response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
			}
			addDebug(actionPath, response, request)
		}

		if len(create) > 0 {
			var insertTags []RestOtsTagInfo
			for _, t := range create {
				insertTags = append(insertTags, RestOtsTagInfo{
					Key:   t.Key,
					Value: t.Value,
				})
			}

			actionPath := "/v2/openapi/tagresources"
			request := make(map[string]interface{})
			request["RegionId"] = StringPointer(client.RegionId)
			request["ResourceType"] = "INSTANCE"
			// id is instanceName
			request["ResourceIds"] = expandStringList([]interface{}{d.Id()})
			request["Tags"] = insertTags

			response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
			}
			addDebug(actionPath, response, request)
		}
		d.SetPartial("tags")
	}
	if !d.IsNewResource() && d.HasChange("alias_name") {
		actionPath := "/v2/openapi/updateinstance"
		request := make(map[string]interface{})
		request["RegionId"] = StringPointer(client.RegionId)
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())
		request["AliasName"] = StringPointer(d.Get("alias_name").(string))

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		d.SetPartial("alias_name")
	}
	if !d.IsNewResource() && d.HasChange("policy") {
		actionPath := "/v2/openapi/updateinstancepolicy"
		request := make(map[string]interface{})
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())
		request["Policy"] = StringPointer(d.Get("policy").(string))
		request["PolicyVersion"] = d.Get("policy_version").(int)

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		d.SetPartial("policy")
	}
	if !d.IsNewResource() && d.HasChange("elastic_vcu_upper_limit") {
		actionPath := "/v2/openapi/updateinstanceelasticvcuupperlimit"
		request := make(map[string]interface{})
		// id is instanceName
		request["InstanceName"] = StringPointer(d.Id())
		request["ElasticVCUUpperLimit"] = d.Get("elastic_vcu_upper_limit").(float64)

		response, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		addDebug(actionPath, response, request)
		d.SetPartial("elastic_vcu_upper_limit")
	}
	if err := otsService.WaitForOtsInstance(d.Id(), toInstanceInnerStatus(Running), DefaultTimeout); err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAliyunOtsInstanceRead(d, meta)
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	// A reserved-mode (VCU) instance is deleted via DeleteVCUInstance.
	// Triggered when vcu > 0. A free instance (vcu == 0) is deleted via
	// the standard deleteinstance API. Read() backfills vcu from
	// VCUQuota so imported VCU instances are also routed correctly.
	// payment_type is not used here: the GetInstance response returns
	// "PayAsYouGo" for both free and VCU instances, so a string-based
	// branch would misroute free instances into DeleteVCUInstance and
	// hit "can't delete a non-VCU instance" (IllegalOp 400).
	if d.Get("vcu").(int) > 0 {
		actionPath := "/v2/openapi/deletevcuinstance"
		urlQuery := make(map[string]*string)
		// id is instanceName
		urlQuery["InstanceName"] = StringPointer(d.Id())

		_, err := OtsRestApiDeleteWithRetry(client, "tablestore", "2020-12-09", actionPath, urlQuery)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
		}
		return WrapError(otsService.WaitForOtsInstance(d.Id(), string(Deleted), DefaultLongTimeout))
	}

	actionPath := "/v2/openapi/deleteinstance"
	request := make(map[string]interface{})
	request["RegionId"] = StringPointer(client.RegionId)
	// id is instanceName
	request["InstanceName"] = StringPointer(d.Id())

	_, err := OtsRestApiPostWithRetry(client, "tablestore", "2020-12-09", actionPath, request)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), actionPath, AlibabaCloudSdkGoERROR)
	}

	return WrapError(otsService.WaitForOtsInstance(d.Id(), string(Deleted), DefaultLongTimeout))
}
