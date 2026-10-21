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
		// storages hashes every attribute (so configuration changes stay
		// visible) while Read back-fills the readable ones; the diff callback
		// drops the back-fill-only difference at plan time. See
		// resourceApplicationCustomizeDiff.
		CustomizeDiff: resourceApplicationCustomizeDiff,
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
				ValidateFunc: validators.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				ForceNew:     true,
				Optional:     true,
				Default:      PostPaid,
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
				DiffSuppressFunc: postPaidDiffSuppressFunc,
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
				// spurious ID-only refresh diff. Unlike parameters,
				// security_ip_list is Optional+Computed because the whitelist
				// is reported and Read back-fills it.
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
			// component_id, the whitelist array name and its IP list are all
			// reported by DescribeApplicationAttribute, so they are Computed as
			// well: an omitted configuration value takes the reported one and an
			// import restores the whitelist state.
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modify_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// property_completeness governance: the following CreateApplication
			// parameters are create-only inputs. Most are not surfaced by
			// DescribeApplicationAttribute, so they stay Optional+ForceNew and
			// keep their configured value on Read. The exceptions — the ones the
			// query API returns (VpcNatGatewayId plus the readable sub-fields of
			// storages/mem_application_spec) — are Optional+Computed+ForceNew
			// and are back-filled in Read, so an imported application restores
			// them instead of planning a replacement against an empty state;
			// Computed keeps a server-side default from turning an omitted
			// configuration value into permanent drift.
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
				// CreateApplication accepts DnatIpAddress only together with
				// VpcNatGatewayId; express that at plan time instead of
				// building a request the API rejects.
				RequiredWith: []string{"vpc_nat_gateway_id"},
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
				Computed: true,
				ForceNew: true,
			},
			// security_ip_type is the IP type of the application-level
			// whitelist entry; DescribeApplicationAttribute reports it through
			// SecurityIPArrays[].SecurityIPType, so it is Computed as well.
			"security_ip_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			// storages is also Computed: the query API reports a default
			// storage even when the configuration declares none, so without
			// Computed the back-filled member would turn "not configured" into
			// a permanent replacement diff.
			"storages": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				// All attributes are part of the member hash (see
				// storageMemberHash): SDKv2 compares the set hashes and skips
				// the per-member diff when they match, so an attribute left
				// out of the hash would hide its configuration change. The
				// readable attributes are back-filled by Read, which makes
				// them differ between a configuration that omits them and the
				// state; resourceApplicationCustomizeDiff removes exactly
				// that back-fill-only difference.
				Set: storageMemberHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// DescribeApplicationAttribute returns StorageType/
						// StorageInstanceId/StorageCapacity/StoragePerformanceLevel;
						// Read merges them into the matching member so the four
						// mount-only fields below (never returned) survive.
						"storage_type": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"storage_instance_id": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
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
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"storage_endpoint": {
							Type: schema.TypeString, Optional: true, ForceNew: true,
						},
						"storage_performance_level": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
					},
				},
			},
			// mem_application_spec is Computed for the same reason as storages:
			// Read back-fills MemApplicationAttribute, and an application whose
			// type does not require the block must not gain a replacement diff
			// from a server-reported member.
			"mem_application_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// DescribeApplicationAttribute.MemApplicationAttribute
						// returns DbName/UserName/*ModelName/ProjectName/MemoryType
						// and the visual-model fields; Read merges them into the
						// configured member. db_password, shard and
						// embedder_model_dimension are write-only and keep their
						// configured value.
						"db_name": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"db_password": {
							Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true,
						},
						"db_user": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"embedder_model": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"embedder_model_dimension": {
							Type: schema.TypeInt, Optional: true, ForceNew: true,
						},
						"graph_llm_model": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"llm_model": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"memory_type": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"project_name": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"reranker_model": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"shard": {
							Type: schema.TypeInt, Optional: true, ForceNew: true,
						},
						"visual_embedder_model": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
						},
						"visual_embedder_model_dimension": {
							Type: schema.TypeInt, Optional: true, Computed: true, ForceNew: true,
						},
						// visual_llm_modalities is an order-insensitive set of
						// modality strings; TypeSet so reorders are not drift
						// (see the storages comment above).
						"visual_llm_modalities": {
							Type: schema.TypeSet, Optional: true, Computed: true, ForceNew: true,
							Elem: &schema.Schema{Type: schema.TypeString},
						},
						"visual_llm_model": {
							Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true,
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
				Computed: true,
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
			// Read leaves it at the configured value: DescribeApplicationAttribute
			// returns DnatMappings[].FrontPort, but the API allocates a front port
			// when the request omits it, and this attribute set is
			// Optional+ForceNew without a Computed member, so writing an
			// allocated port back would make an omitted option a permanent
			// replacement diff.
			"dnat_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				// CreateApplication accepts DnatEntries only together with
				// VpcNatGatewayId; express that at plan time instead of
				// building a request the API rejects.
				RequiredWith: []string{"vpc_nat_gateway_id"},
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
		if jsonData, ok, err := marshalApplicationParameters(d); err != nil {
			return errs.WrapError(err)
		} else if ok {
			request["Parameters"] = jsonData
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
		d.Set("pay_type", getChargeType(v))
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
	// The query API reports the CPU architecture; without it an imported
	// application would plan a replacement against an empty attribute.
	if v, ok := applicationAttribute["Architecture"].(string); ok {
		d.Set("architecture", v)
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
	// A single-component application reports its component id; restoring it lets
	// an import reproduce the component reference instead of planning a diff.
	if len(components) == 1 {
		if component, ok := components[0].(map[string]interface{}); ok {
			if componentID, _ := component["ComponentId"].(string); componentID != "" {
				if err := d.Set("component_id", componentID); err != nil {
					return errs.WrapError(err)
				}
			}
		}
	}
	if v, ok := applicationAttribute["Status"].(string); ok {
		d.Set("status", v)
	}
	if v, ok := applicationAttribute["DBClusterId"].(string); ok {
		d.Set("db_cluster_id", v)
	}
	// Create-only inputs that DescribeApplicationAttribute does return are
	// written back, so an imported application restores them instead of
	// planning a replacement against an empty state. An explicitly returned
	// empty VpcNatGatewayId clears the field (a gateway detached out-of-band
	// then shows up as a replacement); an omitted key is left alone, because
	// "not returned yet" must not erase a configured value right after create.
	if v, ok := applicationAttribute["VpcNatGatewayId"].(string); ok {
		if err := d.Set("vpc_nat_gateway_id", v); err != nil {
			return errs.WrapError(err)
		}
	}
	// The storage list is the complete current set: reconcile even when the
	// response is empty so a member removed in the console leaves the state.
	if v, ok := applicationAttribute["Storages"].([]interface{}); ok {
		if err := d.Set("storages", mergeApplicationStorages(d.Get("storages"), v)); err != nil {
			return errs.WrapError(err)
		}
	}
	// The application-level whitelist is reported through SecurityIPArrays;
	// restoring the array name, IP list and IP type lets an import reproduce
	// the whitelist state instead of planning an avoidable update.
	if v, ok := applicationAttribute["SecurityIPArrays"].([]interface{}); ok {
		if name, ipList, ipListPresent, ipType, found := applicationWhitelistFromResponse(v, d.Get("security_ip_array_name").(string)); found {
			if name != "" {
				if err := d.Set("security_ip_array_name", name); err != nil {
					return errs.WrapError(err)
				}
			}
			if ipType != "" {
				if err := d.Set("security_ip_type", ipType); err != nil {
					return errs.WrapError(err)
				}
			}
			// The array was located and the response carries an IP list: an
			// explicitly empty list clears stale state, while a missing key
			// leaves the configured value alone.
			if ipListPresent {
				if err := d.Set("security_ip_list", ipList); err != nil {
					return errs.WrapError(err)
				}
			}
		}
	}
	if v, ok := applicationAttribute["Endpoints"].([]interface{}); ok {
		if endpoints := mergeApplicationEndpoints(d.Get("endpoints"), v); endpoints != nil {
			if err := d.Set("endpoints", endpoints); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	if v, ok := applicationAttribute["MemApplicationAttribute"].(map[string]interface{}); ok {
		if merged, found := mergeApplicationMemSpec(d.Get("mem_application_spec"), v); found {
			if err := d.Set("mem_application_spec", []interface{}{merged}); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	// Tags are read back via ListTagResources (ResourceType="application").
	// A successful empty response must be written back whenever state still
	// holds tags: after removing every tag, or after the tags were cleared
	// out-of-band, skipping d.Set would keep the stale map, keep tags.% at its
	// old value and hide the drift. NotFound from ListTagResources, expected
	// right after create before the tag service has indexed the new
	// application, remains a supplementary source that leaves tags at their
	// config/prior value instead; any other error class is propagated, so
	// authorization or request-shape failures are not silently hidden.
	// The one case left unset is "no tags in the response, none in state and
	// no tags change": writing an empty map there is asymmetric across the
	// SDKv2 apply and refresh paths (apply normalizes it away while refresh
	// serializes tags.% to "0"), which would add a spurious ID-only refresh
	// diff without clearing anything.
	if fetched, err := service.ListApplicationTags(d.Id()); err != nil {
		if !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}
		logs.AddDebug("ListTagResources not found; leaving tags at prior value", err, map[string]interface{}{"ApplicationId": d.Id()})
	} else {
		priorTags, _ := d.Get("tags").(map[string]interface{})
		if len(fetched) > 0 || len(priorTags) > 0 || d.HasChange("tags") {
			if err := d.Set("tags", fetched); err != nil {
				return errs.WrapError(err)
			}
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
	if PayType(payType) == Prepaid {
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
	if err := validateApplicationTypeSpecs(d); err != nil {
		return nil, errs.WrapError(err)
	}
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
			// An all-optional, non-computed block can yield a nil element when
			// the user declares the block without any sub-field.
			item, ok := comp.(map[string]interface{})
			if !ok || item == nil {
				continue
			}
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
	request["PayType"] = string(Postpaid)
	if payType == string(PrePaid) {
		request["PayType"] = string(Prepaid)
	}
	if PayType(request["PayType"].(string)) == Prepaid {
		period := d.Get("period").(int)
		if period <= 0 {
			return nil, errs.WrapError(fmt.Errorf("'period' must be greater than 0 when 'pay_type' is 'PrePaid'"))
		}
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = string(Month)
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = string(Year)
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

	// CreateApplication accepts Parameters (style:"json") at create time; the
	// same list is later maintained through ModifyApplicationParameter. Sending
	// it here means a first apply whose configuration carries parameters really
	// writes them instead of keeping them in state only.
	if s, ok, err := marshalApplicationParameters(d); err != nil {
		return nil, errs.WrapError(err)
	} else if ok {
		request["Parameters"] = s
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
