// Package polardb: alicloud_polardb_application migrated verbatim from the
// flat alicloud package (resource_alicloud_polardb_application.go) into the
// polardb product package, per the provider-resource-migration standard. This
// is the behavior-preserving migration pass: the schema, CRUD logic and
// create-request mapping are unchanged; only the package, factory name and
// helper qualifications change. The governance pass (property_completeness
// for CreateApplication/TagResources) layers additional schema fields and
// mappings on top afterwards.
package polardb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceApplication is the SDKv2 factory for alicloud_polardb_application,
// registered through ServicePackage. Renamed from the flat-layout
// resourceAlicloudPolarDBApplication per the migration standard
// (factory = Resource<Res>).
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringLenBetween(2, 256),
			},
			"application_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"components": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"component_class": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"component_replica": {
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},
			"pay_type": {
				Type:         schema.TypeString,
				ValidateFunc: validators.StringInSlice([]string{string(helper.PostPaid), string(helper.PrePaid)}, false),
				ForceNew:     true,
				Optional:     true,
				Default:      helper.PostPaid,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validators.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: validators.PolarDBPostPaidDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_from": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ai_db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_api_key": {
				Type:      schema.TypeString,
				ForceNew:  true,
				Optional:  true,
				Sensitive: true,
			},
			"model_base_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_api": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"model_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_version": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"parameters": {
				// Write-only: CreateApplication accepts Parameters, but
				// DescribeApplicationAttribute does not return them
				// (parameter_value is Sensitive), so Read never sets this
				// block. Optional only (not Computed): an unset Optional
				// collection stays absent on both the apply and refresh
				// paths, whereas Optional+Computed makes the SDK inject
				// parameters.#="0" on refresh but omit it on apply — a
				// spurious ID-only refresh diff. This mirrors security_ip_list
				// (another write-only TypeList on this resource).
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter_value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// property_completeness governance: the following CreateApplication
			// parameters are create-only inputs (operatePrivateType=[create] +
			// readOnly, not surfaced by DescribeApplicationAttribute), so they
			// are Optional+ForceNew and retained from config on Read. Tag is the
			// exception: it is mutable via TagResources/UntagResources and read
			// back via ListTagResources (ResourceType="application").
			"target_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"skill_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dnat_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auth_provider": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auth_provider_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_nat_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_ip_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auto_allocate_public_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"agentic_db_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// storages/endpoints/dnat_entries (and the nested
			// visual_llm_modalities) are order-insensitive collections of
			// independent members (storage mounts, service endpoints, DNAT
			// port mappings, LLM modalities): CreateApplication stores them as
			// a set, so member order is not drift. Modeled as TypeSet (not
			// TypeList) so the SDK ignores pure reorders and the resource is
			// exempt from the collection-order convergent-reorder coverage gate.
			"storages": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"storage_instance_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"endpoint_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"mount_path": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"container_mount_path": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"storage_capacity": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"storage_endpoint": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"storage_performance_level": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
					},
				},
			},
			"mem_application_spec": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"db_password": {
							Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true,
						},
						"db_user": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"embedder_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"embedder_model_dimension": {
							Type: schema.TypeInt, Optional: true, ForceNew: true,
						},
						"graph_llm_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"llm_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"memory_type": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"project_name": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"reranker_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"shard": {
							Type: schema.TypeInt, Optional: true, ForceNew: true,
						},
						"visual_embedder_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"visual_embedder_model_dimension": {
							Type: schema.TypeInt, Optional: true, ForceNew: true,
						},
						// visual_llm_modalities is an order-insensitive set of
						// modality strings; TypeSet so reorders are not drift
						// (see the storages comment above).
						"visual_llm_modalities": {
							Type: schema.TypeSet, Optional: true, ForceNew: true,
							Elem: &schema.Schema{Type: schema.TypeString},
						},
						"visual_llm_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
					},
				},
			},
			"agentic_db_branch_spec": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"branch_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"db_cluster_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"fork_from_application_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"fork_from_branch": {
							Type: schema.TypeBool, Optional: true, ForceNew: true,
						},
						"project_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"tenant_id": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
					},
				},
			},
			// endpoints: order-insensitive set; see the storages comment.
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_type": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"description": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
					},
				},
			},
			"knowledge_application_spec": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dashboard_password": {
							Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true,
						},
						"db_password": {
							Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true,
						},
						"llm_model": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
					},
				},
			},
			// dnat_entries: order-insensitive set; see the storages comment.
			"dnat_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_name": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"front_port": {
							Type: schema.TypeInt, Optional: true, ForceNew: true,
						},
					},
				},
			},
			"tags": tags.TagsSchemaWithElements(),
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	request, err := service.buildApplicationCreateRequest(d)
	if err != nil {
		return errs.WrapError(err)
	}
	var response map[string]interface{}
	action := "CreateApplication"

	// Execute CreateApplication only once, no retry
	response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
	if err != nil {
		logs.AddDebug(action, response, request)
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_application", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["ApplicationId"]))
	// wait application status change from Creating to running
	stateConf := helper.BuildStateConf([]string{"Creating"}, []string{"Activated"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, service.applicationStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceApplicationRead(d, meta)
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	if d.HasChange("upgrade_version") && d.Get("upgrade_version").(bool) {
		action := "UpgradeApplicationVersion"
		request := map[string]interface{}{
			"ApplicationId": d.Id(),
		}
		err := retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					return retry.RetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{"MinorVersionUpgrading"}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, service.applicationStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	if d.HasChanges("parameters") {
		action := "ModifyApplicationParameter"
		request := map[string]interface{}{
			"ApplicationId": d.Id(),
		}
		if v, ok := d.GetOk("parameters"); ok {
			parameters := make([]map[string]interface{}, 0)
			params := v.([]interface{})
			for _, param := range params {

				paramMap := make(map[string]interface{})
				item := param.(map[string]interface{})

				if v, ok := item["parameter_name"].(string); ok {
					paramMap["ParameterName"] = v
				}
				if v, ok := item["parameter_value"].(string); ok {
					paramMap["ParameterValue"] = v
				}

				parameters = append(parameters, paramMap)
			}
			jsonData, err := json.Marshal(parameters)
			if err != nil {
				return errs.WrapError(err)
			}
			request["Parameters"] = string(jsonData)
		}

		response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			logs.AddDebug(action, response, request)
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_application", action, errs.AlibabaCloudSdkGoERROR)
		}
		// wait application status change from ClassChanging to running
		stateConf := helper.BuildStateConf([]string{"ClassChanging", "Restarting"}, []string{"Activated"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, service.applicationStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChanges("security_ip_list", "security_ip_array_name", "modify_mode", "component_id") {
		action := "ModifyApplicationWhitelist"
		modifyRequest := map[string]interface{}{
			"ApplicationId": d.Id(),
		}
		if _, ok := d.GetOk("modify_mode"); ok {
			modifyRequest["ModifyMode"] = d.Get("modify_mode").(string)
		}
		if _, ok := d.GetOk("security_ip_array_name"); ok {
			modifyRequest["SecurityIPArrayName"] = d.Get("security_ip_array_name").(string)
		}
		if _, ok := d.GetOk("security_ip_list"); ok {
			ipList := ""
			ids := d.Get("security_ip_list").([]interface{})
			for _, id := range ids {
				ipList += id.(string)
				ipList += ","
			}
			modifyRequest["SecurityIPList"] = strings.TrimRight(ipList, ",")
		}
		if _, ok := d.GetOk("component_id"); ok {
			modifyRequest["ComponentId"] = d.Get("component_id").(string)
		}

		response, err := client.RpcPost("polardb", "2017-08-01", action, nil, modifyRequest, false)

		if err != nil {
			logs.AddDebug(action, response, modifyRequest)
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("tags") {
		// item 65: tags mutation on a live application uses TagResources to
		// add/update and UntagResources to remove, both with
		// ResourceType="application" (the CreateApplication inline Tag param
		// only applies at create time). Untag before Tag: a changed key shows
		// up in both the added set and the removed list (ParsingTags removes a
		// key when its value differs), so untagging first then re-tagging
		// yields the new value instead of deleting the key.
		addedTags, removedTagKeys := tags.ParsingTags(d)
		if len(removedTagKeys) > 0 {
			if err := service.UntagApplicationResources(d.Id(), removedTagKeys); err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "UntagResources", errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(addedTags) > 0 {
			if err := service.TagApplicationResources(d.Id(), addedTags); err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "TagResources", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	d.Partial(false)
	return resourceApplicationRead(d, meta)
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	applicationAttribute, err := service.DescribeApplicationAttribute(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if v, ok := applicationAttribute["RegionId"].(string); ok {
		d.Set("region_id", v)
	}
	if v, ok := applicationAttribute["ZoneId"].(string); ok {
		d.Set("zone_id", v)
	}
	if v, ok := applicationAttribute["VSwitchId"].(string); ok {
		d.Set("vswitch_id", v)
	}
	if v, ok := applicationAttribute["PayType"].(string); ok {
		d.Set("pay_type", helper.GetChargeType(v))
	}
	if v, ok := applicationAttribute["VPCId"].(string); ok {
		d.Set("vpc_id", v)
	}
	if v, ok := applicationAttribute["Description"].(string); ok {
		d.Set("description", v)
	}
	if v, ok := applicationAttribute["ApplicationType"].(string); ok {
		d.Set("application_type", v)
	}
	components, _ := applicationAttribute["Components"].([]interface{})
	componentsList := make([]map[string]interface{}, 0, len(components))
	for i, _ := range components {
		componentsMap := components[i].(map[string]interface{})
		component := make(map[string]interface{})
		if v, ok := componentsMap["ComponentType"]; ok {
			componentType := v.(string)
			component["component_type"] = convertComponentTypeReadResponse(componentType)
		}
		if v, ok := componentsMap["ComponentClass"]; ok {
			component["component_class"] = v.(string)
		}
		if v, ok := componentsMap["ComponentReplica"]; ok {
			if num, ok := v.(json.Number); ok {
				if replica, err := num.Int64(); err == nil {
					component["component_replica"] = int(replica)
				}
			} else if replica, ok := v.(int); ok {
				component["component_replica"] = replica
			}
		}
		componentsList = append(componentsList, component)
	}
	d.Set("components", componentsList)
	if v, ok := applicationAttribute["Status"].(string); ok {
		d.Set("status", v)
	}
	if v, ok := applicationAttribute["DBClusterId"].(string); ok {
		d.Set("db_cluster_id", v)
	}
	// Tags are read back via ListTagResources (ResourceType="application").
	// d.Set is only called when the API actually returns tags. Setting an
	// *empty* collection here is asymmetric across the SDKv2 apply and refresh
	// paths: during apply it is suppressed against an empty config (tags.%
	// stays absent), while during refresh it serializes tags.% to "0" — a
	// spurious ID-only refresh diff. Skipping d.Set when there are no tags
	// leaves tags.% absent on both paths (when the config has none). A
	// not-found from ListTagResources, expected right after create before the
	// tag service has indexed the new application, is treated the same way:
	// tags is left at its config/prior value rather than reset to empty. Any
	// other error class is propagated, so authorization or request-shape
	// failures are not silently hidden by the supplementary tag query. When
	// tags exist, both paths d.Set the non-empty map (a non-empty collection
	// d.Set, like components, serializes consistently).
	if fetched, err := service.ListApplicationTags(d.Id()); err != nil {
		if !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}
		logs.AddDebug("ListTagResources not found; leaving tags at prior value", err, map[string]interface{}{"ApplicationId": d.Id()})
	} else if len(fetched) > 0 {
		if err := d.Set("tags", fetched); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	app, err := service.DescribeApplicationAttribute(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}
	// Pre paid application can not be release.
	payType, _ := app["PayType"].(string)
	if helper.PayType(payType) == helper.Prepaid {
		return errs.WrapError(errs.Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}
	_, err = service.DeleteApplication(d.Id())
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteApplication", errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, service.applicationStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		if strings.HasPrefix(err.Error(), "couldn't find resource") {
			d.SetId("")
			return nil
		}
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	d.SetId("")
	return nil
}

// buildApplicationCreateRequest assembles the CreateApplication RPC request.
// It is a method on *Service so the zone-derivation fallback (when ZoneId is
// omitted but VSwitchId is set) can reuse this package's describeVSwitch
// reimplemented on the generic RpcPost entry, instead of the flat-layout
// VpcService.DescribeVSwitch that is no longer reachable from here.
func (s *Service) buildApplicationCreateRequest(d *schema.ResourceData) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"RegionId":        helper.Trim(d.Get("region_id").(string)),
		"Description":     helper.Trim(d.Get("description").(string)),
		"ApplicationType": helper.Trim(d.Get("application_type").(string)),
		"Architecture":    helper.Trim(d.Get("architecture").(string)),
		"ClientToken":     helper.BuildClientToken("CreateApplication"),
	}

	if v, ok := d.GetOk("db_cluster_id"); ok && v.(string) != "" {
		request["DBClusterId"] = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v.(string)
	}

	if zone, ok := d.GetOk("zone_id"); ok && helper.Trim(zone.(string)) != "" {
		request["ZoneId"] = helper.Trim(zone.(string))
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v.(string)
	}

	if request["VSwitchId"] != nil {
		if request["ZoneId"] == nil {
			// check vswitchId in zone
			zoneId, err := s.describeVSwitch(request["VSwitchId"].(string))
			if err != nil {
				return nil, errs.WrapError(err)
			}
			request["ZoneId"] = zoneId
		}
	}

	if v, ok := d.GetOk("components"); ok {
		components := make([]map[string]interface{}, 0)
		comps := v.([]interface{})
		for _, comp := range comps {
			item := comp.(map[string]interface{})
			component := map[string]interface{}{}
			if v, ok := item["component_type"].(string); ok {
				component["ComponentType"] = v
			}
			if v, ok := item["component_class"].(string); ok {
				component["ComponentClass"] = v
			}
			if v, ok := item["component_replica"].(int); ok {
				component["ComponentReplica"] = v
			}
			components = append(components, component)
		}
		jsonData, err := json.Marshal(components)
		if err != nil {
			return nil, errs.WrapError(err)
		}
		request["Components"] = string(jsonData)
	}

	payType := helper.Trim(d.Get("pay_type").(string))
	request["PayType"] = string(helper.Postpaid)
	if payType == string(helper.PrePaid) {
		request["PayType"] = string(helper.Prepaid)
	}
	if helper.PayType(request["PayType"].(string)) == helper.Prepaid {
		period := d.Get("period").(int)
		if period <= 0 {
			return nil, errs.WrapError(fmt.Errorf("'period' must be greater than 0 when 'pay_type' is 'PrePaid'"))
		}
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = string(helper.Month)
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = string(helper.Year)
		}
		request["AutoRenew"] = d.Get("auto_renew").(bool)
	}

	if _, ok := d.GetOk("security_ip_list"); ok {
		ipList := ""
		ids := d.Get("security_ip_list").([]interface{})
		for _, id := range ids {
			ipList += id.(string)
			ipList += ","
		}
		request["SecurityIPList"] = strings.TrimRight(ipList, ",")
	}

	if v, ok := d.GetOk("model_from"); ok {
		request["ModelFrom"] = v.(string)
	}

	if v, ok := d.GetOk("ai_db_cluster_id"); ok {
		request["AIDBClusterId"] = v.(string)
	}

	if v, ok := d.GetOk("model_api_key"); ok {
		request["ModelApiKey"] = v.(string)
	}

	if v, ok := d.GetOk("model_base_url"); ok {
		request["ModelBaseUrl"] = v.(string)
	}

	if v, ok := d.GetOk("model_api"); ok {
		request["ModelApi"] = v.(string)
	}

	if v, ok := d.GetOk("model_name"); ok {
		request["ModelName"] = v.(string)
	}

	// property_completeness governance: additional CreateApplication scalar
	// parameters. Each is a create-only input sent only when configured.
	if v, ok := d.GetOk("target_version"); ok {
		request["TargetVersion"] = v.(string)
	}
	if v, ok := d.GetOk("skill_template_id"); ok {
		request["SkillTemplateId"] = v.(string)
	}
	if v, ok := d.GetOk("dnat_ip_address"); ok {
		request["DnatIpAddress"] = v.(string)
	}
	if v, ok := d.GetOk("auth_provider"); ok {
		request["AuthProvider"] = v.(string)
	}
	if v, ok := d.GetOk("auth_provider_config"); ok {
		request["AuthProviderConfig"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_nat_gateway_id"); ok {
		request["VpcNatGatewayId"] = v.(string)
	}
	if v, ok := d.GetOk("security_ip_type"); ok {
		request["SecurityIPType"] = v.(string)
	}
	if v, ok := d.GetOk("security_ip_array_name"); ok {
		request["SecurityIPArrayName"] = v.(string)
	}
	if v, ok := d.GetOk("agentic_db_cluster_id"); ok {
		request["AgenticDBClusterId"] = v.(string)
	}
	request["AutoAllocatePublicEip"] = d.Get("auto_allocate_public_eip").(bool)

	// CreateApplication object/array parameters use style:"json", so they are
	// encoded as JSON strings — mirroring Components above. Object params
	// (MaxItems=1) become a single JSON object; array params become a JSON
	// array. Empty/zero sub-fields are omitted by marshalBlockItem.
	if s, ok, err := marshalBlockList(d, "storages", storagesFields); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["Storages"] = s
	}
	if s, ok, err := marshalBlockObject(d, "mem_application_spec", memApplicationSpecFields); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["MemApplicationSpec"] = s
	}
	if s, ok, err := marshalBlockObject(d, "agentic_db_branch_spec", agenticDBBranchSpecFields); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["AgenticDBBranchSpec"] = s
	}
	if s, ok, err := marshalBlockList(d, "endpoints", endpointsFields); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["Endpoints"] = s
	}
	if s, ok, err := marshalBlockObject(d, "knowledge_application_spec", knowledgeApplicationSpecFields); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["KnowledgeApplicationSpec"] = s
	}
	if s, ok, err := marshalBlockList(d, "dnat_entries", dnatEntriesFields); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["DnatEntries"] = s
	}

	// CreateApplication accepts Tag as a repeatList (Tag.N.Key/Tag.N.Value);
	// tags set at create time are applied inline here. Subsequent mutations go
	// through TagResources/UntagResources in Update.
	if v, ok := d.GetOk("tags"); ok {
		if tagMap := v.(map[string]interface{}); len(tagMap) > 0 {
			tags.ExpandTagsToMap(request, tags.ConvertTags(tagMap))
		}
	}

	return request, nil
}

func convertComponentTypeReadResponse(source string) string {
	switch source {
	case "polarclaw":
		return "polarclaw_comp"
	}
	return source
}
