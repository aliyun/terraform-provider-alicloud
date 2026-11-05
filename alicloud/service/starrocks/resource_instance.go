// Package starrocks. This file is generated automatically. Please do not modify it manually, thank you!
package starrocks

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create:        resourceInstanceCreate,
		Read:          resourceInstanceRead,
		Update:        resourceInstanceUpdate,
		Delete:        resourceInstanceDelete,
		CustomizeDiff: resourceInstanceCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"admin_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backend_node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resident_node_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"local_storage_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"cluster_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"frontend_node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"storage_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resident_node_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"local_storage_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"observer_node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"storage_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resident_node_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"local_storage_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"oss_accessing_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"promotion_option_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"run_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitches": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/webapi/cluster/createV1")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["clientToken"] = helper.BuildClientToken(action)

	request["RegionId"] = client.RegionId
	request["InstanceName"] = d.Get("instance_name")
	request["AdminPassword"] = d.Get("admin_password")
	request["Version"] = d.Get("version")
	request["RunMode"] = d.Get("run_mode")
	request["PackageType"] = d.Get("package_type")
	request["PayType"] = d.Get("pay_type")
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("oss_accessing_role_name"); ok {
		request["OssAccessingRoleName"] = v
	}
	if v, ok := d.GetOkExists("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("frontend_node_groups"); ok {
		frontendNodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["cu"] = dataLoopTmp["cu"]
			dataLoopMap["storageSize"] = dataLoopTmp["storage_size"]
			dataLoopMap["storagePerformanceLevel"] = dataLoopTmp["storage_performance_level"]
			dataLoopMap["diskNumber"] = dataLoopTmp["disk_number"]
			dataLoopMap["residentNodeNumber"] = dataLoopTmp["resident_node_number"]
			dataLoopMap["specType"] = dataLoopTmp["spec_type"]
			dataLoopMap["localStorageInstanceType"] = dataLoopTmp["local_storage_instance_type"]
			dataLoopMap["zoneId"] = dataLoopTmp["zone_id"]
			frontendNodeGroupsMapsArray = append(frontendNodeGroupsMapsArray, dataLoopMap)
		}
		request["FrontendNodeGroups"] = frontendNodeGroupsMapsArray
	}

	if v, ok := d.GetOk("backend_node_groups"); ok {
		backendNodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["cu"] = dataLoop1Tmp["cu"]
			dataLoop1Map["storageSize"] = dataLoop1Tmp["storage_size"]
			dataLoop1Map["storagePerformanceLevel"] = dataLoop1Tmp["storage_performance_level"]
			dataLoop1Map["diskNumber"] = dataLoop1Tmp["disk_number"]
			dataLoop1Map["residentNodeNumber"] = dataLoop1Tmp["resident_node_number"]
			dataLoop1Map["specType"] = dataLoop1Tmp["spec_type"]
			dataLoop1Map["localStorageInstanceType"] = dataLoop1Tmp["local_storage_instance_type"]
			dataLoop1Map["zoneId"] = dataLoop1Tmp["zone_id"]
			backendNodeGroupsMapsArray = append(backendNodeGroupsMapsArray, dataLoop1Map)
		}
		request["BackendNodeGroups"] = backendNodeGroupsMapsArray
	}

	if v, ok := d.GetOk("vswitches"); ok {
		vSwitchesMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range v.([]interface{}) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["VswId"] = dataLoop2Tmp["vswitch_id"]
			dataLoop2Map["ZoneId"] = dataLoop2Tmp["zone_id"]
			vSwitchesMapsArray = append(vSwitchesMapsArray, dataLoop2Map)
		}
		request["VSwitches"] = vSwitchesMapsArray
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["ZoneId"] = d.Get("cluster_zone_id")
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["Encrypted"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}
	if v, ok := d.GetOk("promotion_option_no"); ok {
		request["PromotionOptionNo"] = v
	}
	if v, ok := d.GetOk("observer_node_groups"); ok {
		observerNodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop4 := range v.([]interface{}) {
			dataLoop4Tmp := dataLoop4.(map[string]interface{})
			dataLoop4Map := make(map[string]interface{})
			dataLoop4Map["cu"] = dataLoop4Tmp["cu"]
			dataLoop4Map["storageSize"] = dataLoop4Tmp["storage_size"]
			dataLoop4Map["storagePerformanceLevel"] = dataLoop4Tmp["storage_performance_level"]
			dataLoop4Map["diskNumber"] = dataLoop4Tmp["disk_number"]
			dataLoop4Map["residentNodeNumber"] = dataLoop4Tmp["resident_node_number"]
			dataLoop4Map["specType"] = dataLoop4Tmp["spec_type"]
			dataLoop4Map["localStorageInstanceType"] = dataLoop4Tmp["local_storage_instance_type"]
			dataLoop4Map["zoneId"] = dataLoop4Tmp["zone_id"]
			observerNodeGroupsMapsArray = append(observerNodeGroupsMapsArray, dataLoop4Map)
		}
		request["ObserverNodeGroups"] = observerNodeGroupsMapsArray
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_star_rocks_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	starrocksService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, starrocksService.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceInstanceUpdate(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	starrocksService := NewService(client)

	objectRaw, err := starrocksService.DescribeInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_star_rocks_instance DescribeInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["BeginTime"])
	d.Set("encrypted", objectRaw["Encrypted"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("kms_key_id", objectRaw["KmsKeyId"])
	d.Set("package_type", objectRaw["PackageType"])
	d.Set("pay_type", objectRaw["PayType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("run_mode", objectRaw["RunMode"])
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("version", objectRaw["Version"])
	d.Set("vpc_id", objectRaw["VpcId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))
	vSwitchesRaw := objectRaw["VSwitches"]
	vswitchesMaps := make([]map[string]interface{}, 0)
	if vSwitchesRaw != nil {
		for _, vSwitchesChildRaw := range vSwitchesRaw.([]interface{}) {
			vswitchesMap := make(map[string]interface{})
			vSwitchesChildRaw := vSwitchesChildRaw.(map[string]interface{})
			vswitchesMap["vswitch_id"] = vSwitchesChildRaw["VswId"]
			vswitchesMap["zone_id"] = vSwitchesChildRaw["ZoneId"]

			vswitchesMaps = append(vswitchesMaps, vswitchesMap)
		}
	}
	if err := d.Set("vswitches", vswitchesMaps); err != nil {
		return err
	}

	return nil
}

// resourceInstanceCustomizeDiff keeps replacement semantics for a smaller backend
// node group disk: ModifyDiskSize can only expand a disk, so shrinking or removing
// storage_size must still rebuild the instance instead of failing at apply time.
func resourceInstanceCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() == "" || !diff.HasChange("backend_node_groups") || !diff.NewValueKnown("backend_node_groups") {
		return nil
	}
	oldRaw, newRaw := diff.GetChange("backend_node_groups")
	oldGroups, _ := oldRaw.([]interface{})
	newGroups, _ := newRaw.([]interface{})
	if len(oldGroups) != len(newGroups) {
		return nil
	}
	for i := range newGroups {
		oldStorageSize, oldOk := backendNodeGroupStorageSize(oldGroups[i])
		newStorageSize, newOk := backendNodeGroupStorageSize(newGroups[i])
		if backendNodeGroupStorageSizeChangeForcesReplacement(oldStorageSize, oldOk, newStorageSize, newOk) {
			// ForceNew has to name the changed nested attribute: on the parent list it
			// does not reach the element diff, and the plan stays an in-place update.
			return diff.ForceNew(fmt.Sprintf("backend_node_groups.%d.storage_size", i))
		}
	}
	return nil
}

// backendNodeGroupStorageSizeChangeForcesReplacement reports whether a storage_size
// change cannot be applied in place. Removing the value is treated like a shrink
// because the API cannot detach or reduce a disk, and adding a value to a group that
// never had one cannot be compared with the remote size at plan time, so it keeps the
// historical replacement semantics instead of silently skipping the resize.
func backendNodeGroupStorageSizeChangeForcesReplacement(oldSize int, oldOk bool, newSize int, newOk bool) bool {
	if !oldOk {
		return newOk
	}
	if !newOk {
		return true
	}
	return newSize < oldSize
}

func backendNodeGroupStorageSize(raw interface{}) (int, bool) {
	group, ok := raw.(map[string]interface{})
	if !ok {
		return 0, false
	}
	storageSize, ok := group["storage_size"].(int)
	if !ok || storageSize <= 0 {
		return 0, false
	}
	return storageSize, true
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	deadline := time.Now().Add(d.Timeout(schema.TimeoutUpdate))
	starrocksService := NewService(client)

	// Resolve every backend node group match before the first modifying call: an
	// unmatched or ambiguous group has to fail while the instance is still untouched,
	// name and resource group included.
	var resizes []nodeGroupResize
	if !d.IsNewResource() && d.HasChange("backend_node_groups") {
		var err error
		resizes, err = resourceInstanceResolveBackendNodeGroupResizes(d, starrocksService, deadline)
		if err != nil {
			return errs.WrapError(err)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := fmt.Sprintf("/webapi/cluster/update_name")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = helper.StringPointer(d.Id())
	query["RegionId"] = helper.StringPointer(client.RegionId)
	if d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		query["ClusterName"] = helper.StringPointer(v.(string))
	}

	body = request
	if update {
		budget, budgetErr := remainingUpdateDuration(deadline)
		if budgetErr != nil {
			return errs.WrapError(budgetErr)
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(budget, func() *retry.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = fmt.Sprintf("/webapi/resourceGroup/change")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = helper.StringPointer(d.Id())
	query["RegionId"] = helper.StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["NewResourceGroupId"] = helper.StringPointer(v.(string))
	}

	query["ResourceType"] = helper.StringPointer("instance")

	body = request
	if update {
		budget, budgetErr := remainingUpdateDuration(deadline)
		if budgetErr != nil {
			return errs.WrapError(budgetErr)
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(budget, func() *retry.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		if err := starrocksService.SetResourceTagsWithDeadline(d, "instance", deadline); err != nil {
			return errs.WrapError(err)
		}
	}

	if len(resizes) > 0 {
		if err := resourceInstanceApplyBackendNodeGroupResizes(d, starrocksService, resizes, deadline); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceInstanceRead(d, meta)
}

// nodeGroupResize is one resolved in-place disk expansion.
type nodeGroupResize struct {
	index       int
	nodeGroupId string
	target      int
}

// remainingUpdateDuration returns the time left in the shared update budget. Every
// step consumes the same budget so the number of node groups cannot extend the
// configured timeout.
func remainingUpdateDuration(deadline time.Time) (time.Duration, error) {
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return 0, fmt.Errorf("the update did not finish within the configured timeout")
	}
	return remaining, nil
}

// resourceInstanceResolveBackendNodeGroupResizes validates the whole backend node
// group change without modifying anything. The instance resource does not persist
// node group ids and the API changes one node group per call, so every changed
// configuration entry is matched to its remote group by the attributes that stay
// ForceNew. An unmatched or ambiguous entry fails here, before any modifying call.
func resourceInstanceResolveBackendNodeGroupResizes(d *schema.ResourceData, service *Service, deadline time.Time) ([]nodeGroupResize, error) {
	oldRaw, newRaw := d.GetChange("backend_node_groups")
	oldGroups, _ := oldRaw.([]interface{})
	newGroups, _ := newRaw.([]interface{})
	if len(oldGroups) != len(newGroups) {
		return nil, fmt.Errorf("backend_node_groups can only be modified in place, got %d -> %d group(s)", len(oldGroups), len(newGroups))
	}

	groups, err := service.DescribeBackendNodeGroups(d.Id(), deadline)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DescribeNodeGroups", errs.AlibabaCloudSdkGoERROR)
	}

	resizes := make([]nodeGroupResize, 0)
	matchedIndex := make(map[string]int)
	for i := range newGroups {
		oldGroup, _ := oldGroups[i].(map[string]interface{})
		newGroup, _ := newGroups[i].(map[string]interface{})
		oldStorageSize, oldOk := backendNodeGroupStorageSize(oldGroup)
		newStorageSize, newOk := backendNodeGroupStorageSize(newGroup)
		// CustomizeDiff already forces replacement for these transitions. The guard keeps
		// a direct Update call, or a future plan change, from silently skipping the resize.
		if backendNodeGroupStorageSizeChangeForcesReplacement(oldStorageSize, oldOk, newStorageSize, newOk) {
			switch {
			case !newOk:
				return nil, fmt.Errorf("backend_node_groups.%d storage_size cannot be removed in place, the instance has to be replaced", i)
			case !oldOk:
				return nil, fmt.Errorf("backend_node_groups.%d storage_size cannot be set in place without a previous value, the instance has to be replaced", i)
			default:
				return nil, fmt.Errorf("backend_node_groups.%d storage_size can only be expanded in place, got %d -> %d GiB", i, oldStorageSize, newStorageSize)
			}
		}
		if !newOk || newStorageSize == oldStorageSize {
			continue
		}
		nodeGroupId, err := matchBackendNodeGroup(groups, oldGroup)
		if err != nil {
			return nil, fmt.Errorf("backend_node_groups.%d: %s", i, err)
		}
		if previous, ok := matchedIndex[nodeGroupId]; ok {
			return nil, fmt.Errorf("backend_node_groups.%d and backend_node_groups.%d resolve to the same node group %s", previous, i, nodeGroupId)
		}
		matchedIndex[nodeGroupId] = i
		resizes = append(resizes, nodeGroupResize{index: i, nodeGroupId: nodeGroupId, target: newStorageSize})
	}
	return resizes, nil
}

// resourceInstanceApplyBackendNodeGroupResizes expands the resolved node groups. The
// shared deadline covers the write retries and the RUNNING waits of every group.
func resourceInstanceApplyBackendNodeGroupResizes(d *schema.ResourceData, service *Service, resizes []nodeGroupResize, deadline time.Time) error {
	for _, resize := range resizes {
		if err := resourceInstanceApplyBackendNodeGroupResize(d, service, resize, deadline); err != nil {
			return err
		}
	}
	return nil
}

func resourceInstanceApplyBackendNodeGroupResize(d *schema.ResourceData, service *Service, resize nodeGroupResize, deadline time.Time) error {
	action := fmt.Sprintf("/webapi/resourceChange/modifyDiskSize")
	request := make(map[string]interface{})
	query := make(map[string]*string)
	query["InstanceId"] = helper.StringPointer(d.Id())
	query["NodeGroupId"] = helper.StringPointer(resize.nodeGroupId)
	query["RegionId"] = helper.StringPointer(service.client.RegionId)
	query["Target"] = helper.StringPointer(strconv.Itoa(resize.target))
	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = helper.StringPointer(v.(string))
	}
	body := request
	var response map[string]interface{}
	var err error
	if _, budgetErr := remainingUpdateDuration(deadline); budgetErr != nil {
		return budgetErr
	}
	// The SDK level auto retry is disabled for this billable change: it would resend the
	// request without looking at the remote disk.
	response, err = service.client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, false)
	logs.AddDebug(action, response, request)
	if err != nil && backendNodeGroupResizeRejected(err) {
		// A rejection (throttling or another 4xx) means the service did not execute the
		// change, so sending it again is safe while the remaining budget allows it.
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		for err != nil && backendNodeGroupResizeRejected(err) {
			if _, budgetErr := remainingUpdateDuration(deadline); budgetErr != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			wait()
			response, err = service.client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, false)
			logs.AddDebug(action, response, request)
		}
	}
	if err != nil {
		if !helper.NeedRetry(err) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		// The request was accepted but the response was lost, or the service failed: the
		// resize may or may not be applied. Sending the billable request again could
		// expand the same disk twice, so wait for the reported size instead and report an
		// unknown result when it never converges.
		reached, checkErr := backendNodeGroupStorageSizeReached(service, d.Id(), resize.nodeGroupId, resize.target, deadline)
		if !reached {
			return errs.WrapErrorf(fmt.Errorf("the result of the disk expansion for node group %s is unknown, the request may or may not have been applied: %v (last size check: %v)", resize.nodeGroupId, err, checkErr), errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	nodeGroupStateId := d.Id() + ":" + resize.nodeGroupId
	stateBudget, budgetErr := remainingUpdateDuration(deadline)
	if budgetErr != nil {
		return budgetErr
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"RUNNING"}, stateBudget, 20*time.Second, service.NodeGroupStateRefreshFuncWithApi(nodeGroupStateId, "Status", []string{}, func(id string) (map[string]interface{}, error) {
		return service.DescribeNodeGroupWithDeadline(id, deadline)
	}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, nodeGroupStateId)
	}
	return nil
}

// backendNodeGroupResizeRejected reports whether the service refused the resize before
// executing it, which makes an immediate repeat safe. Throttling and other 4xx answers
// are rejections; 5xx and transport errors may have executed the change and must not be
// repeated blindly.
var backendNodeGroupResizeRejectedRegex = regexp.MustCompile(`(?i)throttling|^code: 4[0-9]{2}`)

func backendNodeGroupResizeRejected(err error) bool {
	if err == nil {
		return false
	}
	return backendNodeGroupResizeRejectedRegex.MatchString(err.Error())
}

// backendNodeGroupResizeConvergenceWindow bounds how long a lost resize response waits
// for the remote disk size to converge before the result is reported as unknown.
var backendNodeGroupResizeConvergenceWindow = 30 * time.Second

// BackendNodeGroupResizeConvergenceWindowForTest overrides the convergence window used
// before a lost resize response is reported as unknown, and returns the previous value.
func BackendNodeGroupResizeConvergenceWindowForTest(window time.Duration) time.Duration {
	previous := backendNodeGroupResizeConvergenceWindow
	backendNodeGroupResizeConvergenceWindow = window
	return previous
}

// backendNodeGroupStorageSizeReached waits briefly for the remote group to report at
// least the requested size, so a lost resize response can be resolved without sending
// the billable change again.
// The wait is bounded by the shared update deadline and by a short convergence window.
func backendNodeGroupStorageSizeReached(service *Service, instanceId, nodeGroupId string, target int, deadline time.Time) (bool, error) {
	waitUntil := time.Now().Add(backendNodeGroupResizeConvergenceWindow)
	if remaining := time.Until(deadline); remaining < time.Until(waitUntil) {
		waitUntil = deadline
	}
	var lastErr error
	for {
		if time.Until(waitUntil) <= 0 {
			return false, lastErr
		}
		object, err := service.DescribeNodeGroupWithDeadline(instanceId+":"+nodeGroupId, waitUntil)
		if err == nil {
			if current, convErr := strconv.Atoi(fmt.Sprint(object["StorageSize"])); convErr == nil && current >= target {
				return true, nil
			}
		} else {
			lastErr = err
		}
		if time.Until(waitUntil) <= time.Second {
			return false, lastErr
		}
		time.Sleep(5 * time.Second)
	}
}

// matchBackendNodeGroup resolves the single remote backend node group described by
// the pre-update configuration entry. Node group ids are never persisted, so an
// entry that matches zero or several groups cannot be resolved safely and must not
// be guessed. The zone is deliberately not part of the match: the API reports it for
// some compute groups and omits or relocates it for others, so a group that differs
// only by the reported zone could still be the configured one.
func matchBackendNodeGroup(groups []map[string]interface{}, config map[string]interface{}) (string, error) {
	candidates := make([]string, 0)
	for _, group := range groups {
		if backendNodeGroupAttributesMatch(group, config) {
			candidates = append(candidates, fmt.Sprint(group["NodeGroupId"]))
		}
	}
	switch len(candidates) {
	case 0:
		return "", fmt.Errorf("no backend node group matches the configured cu/spec/disk attributes, cannot determine which group to expand")
	case 1:
		return candidates[0], nil
	default:
		return "", fmt.Errorf("multiple backend node groups match the configured attributes (%s), refusing to guess", strings.Join(candidates, ", "))
	}
}

func backendNodeGroupAttributesMatch(group map[string]interface{}, config map[string]interface{}) bool {
	stringAttributes := map[string]string{
		"spec_type":                   "SpecType",
		"storage_performance_level":   "StoragePerformanceLevel",
		"local_storage_instance_type": "LocalStorageInstanceType",
	}
	for configKey, remoteKey := range stringAttributes {
		value, _ := config[configKey].(string)
		if value == "" || value == "null" {
			continue
		}
		if fmt.Sprint(group[remoteKey]) != value {
			return false
		}
	}
	intAttributes := map[string]string{
		"cu":                   "Cu",
		"disk_number":          "DiskNumber",
		"resident_node_number": "ResidentNodeNumber",
	}
	for configKey, remoteKey := range intAttributes {
		value, _ := config[configKey].(int)
		if value <= 0 {
			continue
		}
		if fmt.Sprint(group[remoteKey]) != fmt.Sprint(value) {
			return false
		}
	}
	return true
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/webapi/cluster/release")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = helper.StringPointer(d.Id())
	query["RegionId"] = helper.StringPointer(client.RegionId)

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	starrocksService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 20*time.Second, starrocksService.InstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
