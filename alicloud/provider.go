package alicloud

import (
	"os"

	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for alicloud
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_REGION", os.Getenv("ALICLOUD_REGION")),
				Description: descriptions["region"],
			},
			"security_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECURITY_TOKEN", os.Getenv("SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ots_instance_name": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'ots_instance_name' has been deprecated from provider version 1.9.7. New field 'instance_name' of resource 'alicloud_ots_table' instead.",
			},
			"log_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOG_ENDPOINT", os.Getenv("LOG_ENDPOINT")),
				Description: descriptions["log_endpoint"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{

			"alicloud_images":         dataSourceAlicloudImages(),
			"alicloud_regions":        dataSourceAlicloudRegions(),
			"alicloud_zones":          dataSourceAlicloudZones(),
			"alicloud_instance_types": dataSourceAlicloudInstanceTypes(),
			"alicloud_instances":      dataSourceAlicloudInstances(),
			"alicloud_vpcs":           dataSourceAlicloudVpcs(),
			"alicloud_vswitches":      dataSourceAlicloudVSwitches(),
			"alicloud_eips":           dataSourceAlicloudEips(),
			"alicloud_key_pairs":      dataSourceAlicloudKeyPairs(),
			"alicloud_kms_keys":       dataSourceAlicloudKmsKeys(),
			"alicloud_dns_domains":    dataSourceAlicloudDnsDomains(),
			"alicloud_dns_groups":     dataSourceAlicloudDnsGroups(),
			"alicloud_dns_records":    dataSourceAlicloudDnsRecords(),
			// alicloud_dns_domain_groups, alicloud_dns_domain_records have been deprecated.
			"alicloud_dns_domain_groups":  dataSourceAlicloudDnsGroups(),
			"alicloud_dns_domain_records": dataSourceAlicloudDnsRecords(),
			// alicloud_ram_account_alias has been deprecated
			"alicloud_ram_account_alias":    dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_aliases":  dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_groups":           dataSourceAlicloudRamGroups(),
			"alicloud_ram_users":            dataSourceAlicloudRamUsers(),
			"alicloud_ram_roles":            dataSourceAlicloudRamRoles(),
			"alicloud_ram_policies":         dataSourceAlicloudRamPolicies(),
			"alicloud_security_groups":      dataSourceAlicloudSecurityGroups(),
			"alicloud_security_group_rules": dataSourceAlicloudSecurityGroupRules(),
			"alicloud_db_instances":         dataSourceAlicloudDBInstances(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_instance":                  resourceAliyunInstance(),
			"alicloud_ram_role_attachment":       resourceAlicloudRamRoleAttachment(),
			"alicloud_disk":                      resourceAliyunDisk(),
			"alicloud_disk_attachment":           resourceAliyunDiskAttachment(),
			"alicloud_security_group":            resourceAliyunSecurityGroup(),
			"alicloud_security_group_rule":       resourceAliyunSecurityGroupRule(),
			"alicloud_db_database":               resourceAlicloudDBDatabase(),
			"alicloud_db_account":                resourceAlicloudDBAccount(),
			"alicloud_db_account_privilege":      resourceAlicloudDBAccountPrivilege(),
			"alicloud_db_backup_policy":          resourceAlicloudDBBackupPolicy(),
			"alicloud_db_connection":             resourceAlicloudDBConnection(),
			"alicloud_db_instance":               resourceAlicloudDBInstance(),
			"alicloud_ess_scaling_group":         resourceAlicloudEssScalingGroup(),
			"alicloud_ess_scaling_configuration": resourceAlicloudEssScalingConfiguration(),
			"alicloud_ess_scaling_rule":          resourceAlicloudEssScalingRule(),
			"alicloud_ess_schedule":              resourceAlicloudEssSchedule(),
			"alicloud_ess_attachment":            resourceAlicloudEssAttachment(),
			"alicloud_vpc":                       resourceAliyunVpc(),
			"alicloud_nat_gateway":               resourceAliyunNatGateway(),
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnet":              resourceAliyunSubnet(),
			"alicloud_vswitch":             resourceAliyunSubnet(),
			"alicloud_route_entry":         resourceAliyunRouteEntry(),
			"alicloud_snat_entry":          resourceAliyunSnatEntry(),
			"alicloud_forward_entry":       resourceAliyunForwardEntry(),
			"alicloud_eip":                 resourceAliyunEip(),
			"alicloud_eip_association":     resourceAliyunEipAssociation(),
			"alicloud_slb":                 resourceAliyunSlb(),
			"alicloud_slb_listener":        resourceAliyunSlbListener(),
			"alicloud_slb_attachment":      resourceAliyunSlbAttachment(),
			"alicloud_slb_server_group":    resourceAliyunSlbServerGroup(),
			"alicloud_slb_rule":            resourceAliyunSlbRule(),
			"alicloud_oss_bucket":          resourceAlicloudOssBucket(),
			"alicloud_oss_bucket_object":   resourceAlicloudOssBucketObject(),
			"alicloud_dns_record":          resourceAlicloudDnsRecord(),
			"alicloud_dns":                 resourceAlicloudDns(),
			"alicloud_dns_group":           resourceAlicloudDnsGroup(),
			"alicloud_key_pair":            resourceAlicloudKeyPair(),
			"alicloud_key_pair_attachment": resourceAlicloudKeyPairAttachment(),
			"alicloud_kms_key":             resourceAlicloudKmsKey(),
			"alicloud_ram_user":            resourceAlicloudRamUser(),
			"alicloud_ram_access_key":      resourceAlicloudRamAccessKey(),
			"alicloud_ram_login_profile":   resourceAlicloudRamLoginProfile(),
			"alicloud_ram_group":           resourceAlicloudRamGroup(),
			"alicloud_ram_role":            resourceAlicloudRamRole(),
			"alicloud_ram_policy":          resourceAlicloudRamPolicy(),
			// alicloud_ram_alias has been deprecated
			"alicloud_ram_alias":                   resourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_alias":           resourceAlicloudRamAccountAlias(),
			"alicloud_ram_group_membership":        resourceAlicloudRamGroupMembership(),
			"alicloud_ram_user_policy_attachment":  resourceAlicloudRamUserPolicyAtatchment(),
			"alicloud_ram_role_policy_attachment":  resourceAlicloudRamRolePolicyAttachment(),
			"alicloud_ram_group_policy_attachment": resourceAlicloudRamGroupPolicyAtatchment(),
			"alicloud_container_cluster":           resourceAlicloudCSSwarm(),
			"alicloud_cs_application":              resourceAlicloudCSApplication(),
			"alicloud_cs_swarm":                    resourceAlicloudCSSwarm(),
			"alicloud_cs_kubernetes":               resourceAlicloudCSKubernetes(),
			"alicloud_cdn_domain":                  resourceAlicloudCdnDomain(),
			"alicloud_router_interface":            resourceAlicloudRouterInterface(),
			"alicloud_ots_table":                   resourceAlicloudOtsTable(),
			"alicloud_ots_instance":                resourceAlicloudOtsInstance(),
			"alicloud_ots_instance_attachment":     resourceAlicloudOtsInstanceAttachment(),
			"alicloud_cms_alarm":                   resourceAlicloudCmsAlarm(),
			"alicloud_log_project":                 resourceAlicloudLogProject(),
			"alicloud_log_store":                   resourceAlicloudLogStore(),
			"alicloud_log_store_index":             resourceAlicloudLogStoreIndex(),
			"alicloud_log_machine_group":           resourceAlicloudLogMachineGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	region, ok := d.GetOk("region")
	if !ok {
		if region == "" {
			region = DEFAULT_REGION
		}
	}
	config := Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Region:    common.Region(region.(string)),
		RegionId:  region.(string),
	}

	if token, ok := d.GetOk("security_token"); ok && token.(string) != "" {
		config.SecurityToken = token.(string)
	}

	if ots_instance_name, ok := d.GetOk("ots_instance_name"); ok && ots_instance_name.(string) != "" {
		config.OtsInstanceName = ots_instance_name.(string)
	}

	if logEndpoint, ok := d.GetOk("log_endpoint"); ok && logEndpoint.(string) != "" {
		config.LogEndpoint = logEndpoint.(string)
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// This is a global MutexKV for use within this plugin.
var alicloudMutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key":     "Access key of alicloud",
		"secret_key":     "Secret key of alicloud",
		"region":         "Region of alicloud",
		"security_token": "Alibaba Cloud Security Token",
		"log_endpoint":   "Alibaba Cloud log service self-define endpoint",
	}
}
