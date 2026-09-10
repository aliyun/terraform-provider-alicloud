package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const ResourceAlicloudCSKubernetesAddon = "resourceAlicloudCSKubernetesAddon"

func resourceAlicloudCSKubernetesAddon() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesAddonCreate,
		Read:   resourceAlicloudCSKubernetesAddonRead,
		Update: resourceAlicloudCSKubernetesAddonUpdate,
		Delete: resourceAlicloudCSKubernetesAddonDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"next_version": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'next_version' has been deprecated from provider version 1.273.0. Please use 'next_version' of DataSource 'alicloud_cs_kubernetes_addons' to replace it",
			},
			"can_upgrade": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "Field 'can_upgrade' has been deprecated from provider version 1.273.0. Please compare 'next_version' and 'current_version' of DataSource 'alicloud_cs_kubernetes_addons' to replace it",
			},
			"required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cleanup_cloud_resources": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCSKubernetesAddonRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	addon, err := csClient.DescribeCsKubernetesAddon(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddon", err)
	}

	d.Set("cluster_id", parts[0])
	d.Set("name", addon.ComponentName)
	d.Set("version", addon.Version)
	d.Set("next_version", addon.NextVersion)
	d.Set("can_upgrade", addon.CanUpgrade)
	d.Set("required", addon.Required)
	if addon.Config != "" {
		config := d.Get("config").(string)
		if config == "" {
			config = "{}"
		}
		localConfig := map[string]interface{}{}
		err := json.Unmarshal([]byte(config), &localConfig)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "json config parse error", err)
		}
		remoteConfig, err := convertJsonStringToMap(addon.Config)
		if err != nil {
			return WrapError(err)
		}
		resultMap := mergeConfigValues(localConfig, remoteConfig)

		result, err := json.Marshal(resultMap)
		if err != nil {
			return WrapError(err)
		}
		if string(result) != "{}" {
			d.Set("config", string(result))
		}
	}

	return nil
}

func resourceAlicloudCSKubernetesAddonCreate(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, name))
	addon, err := csClient.DescribeCsKubernetesAddon(d.Id())
	if err != nil {
		if !NotFoundError(err) {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddonStatus", err)
		} else {
			// Addon not installed
			err := csClient.installAddon(d)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", err)
			}

			// double check addon installed
			_, err = csClient.DescribeCsKubernetesAddon(d.Id())
			if NotFoundError(err) {
				return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", "Unknown Reason. Please check addon config value type first")
			}
		}
	} else {
		// Addon has been installed
		err = updateAddon(csClient, d, addon)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudCSKubernetesAddonRead(d, meta)
}

func resourceAlicloudCSKubernetesAddonUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	clusterId := parts[0]
	addonName := parts[1]

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	addon, err := csClient.GetCsKubernetesAddonInstance(clusterId, addonName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddonStatus", err)
	}

	if addon == nil {
		return WrapError(fmt.Errorf("addon %s not found", d.Id()))
	}

	if d.HasChange("version") || d.HasChange("config") {
		err = updateAddon(csClient, d, addon)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudCSKubernetesAddonRead(d, meta)
}

func resourceAlicloudCSKubernetesAddonDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	clusterId := parts[0]
	addonName := parts[1]

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	addon, err := csClient.GetCsKubernetesAddonInstance(clusterId, addonName)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddon", err)
	}

	if addon.Required || !IsContain(addon.SupportedActions, "Uninstall") {
		log.Printf("[DEBUG] Skip delete system addon %s\n", addonName)
		return nil
	}

	err = csClient.uninstallAddon(d)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "uninstallAddon", err)
	}

	return nil
}

func updateAddon(csClient CsClient, d *schema.ResourceData, addon *Component) error {
	updateVersion, updateConfig, err := needUpgrade(d, addon)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
	}

	if updateVersion {
		err := csClient.upgradeAddon(d, updateVersion, updateConfig)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
		}
	} else if updateConfig {
		err := csClient.updateAddonConfig(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "updateAddonConfig", err)
		}
	}

	return nil
}

func needUpgrade(d *schema.ResourceData, c *Component) (updateVersion bool, updateConfig bool, err error) {
	// Is version changed
	version := d.Get("version").(string)
	if version != "" && c.Version != "" && c.Version != version {
		updateVersion = true
	}

	// Is config changed
	if _, ok := d.GetOk("config"); ok {
		localConfig := map[string]interface{}{}
		err := json.Unmarshal([]byte(d.Get("config").(string)), &localConfig)
		if err != nil {
			return false, false, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "json config parse error", err)
		}

		if c.Config == "" {
			c.Config = "{}"
		}
		remoteConfig := map[string]interface{}{}
		err = json.Unmarshal([]byte(c.Config), &remoteConfig)
		if err != nil {
			return false, false, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "json config parse error", err)
		}
		// if localConfig is subset of remoteConfig, do not call update to do action
		updateConfig = !isSubset(localConfig, remoteConfig)
	}
	return
}

// return true if oldMap is subset of newMap
func isSubset(oldMap, newMap map[string]interface{}) bool {
	for key, oldVal := range oldMap {
		newVal, exists := newMap[key]
		if !exists {
			return false
		}
		oldValMap, oldValIsMap := oldVal.(map[string]interface{})
		newValMap, newValIsMap := newVal.(map[string]interface{})
		if oldValIsMap && newValIsMap {
			if !isSubset(oldValMap, newValMap) {
				return false
			}
			continue
		}

		if oldValIsMap != newValIsMap || !reflect.DeepEqual(oldVal, newVal) {
			return false
		}
	}
	return true
}

func mergeConfigValues(oldMap, newMap map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for key, newValue := range newMap {
		if oldValue, exists := oldMap[key]; exists {
			switch oldValueTyped := oldValue.(type) {
			case map[string]interface{}:
				if newMapTyped, ok := newValue.(map[string]interface{}); ok {
					merged[key] = mergeConfigValues(oldValueTyped, newMapTyped)
					continue
				}
			default:
				merged[key] = newValue
			}
		}
	}

	return merged
}
