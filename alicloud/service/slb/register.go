package slb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Slb service.
// The factories live in this package and are wired into the provider through
// the registry; the SDKResources/SDKDataSources below replace the 29 lines
// previously maintained inline in alicloud/provider.go. The alicloud_slb and
// alicloud_slbs entries are the historical aliases of alicloud_slb_load_balancer
// and alicloud_slb_load_balancers and share their factories.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Slb",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_slb",
				Factory:  ResourceLoadBalancer,
			},
			{
				TypeName: "alicloud_slb_listener",
				Factory:  ResourceListener,
			},
			{
				TypeName: "alicloud_slb_attachment",
				Factory:  ResourceAttachment,
			},
			{
				TypeName: "alicloud_slb_backend_server",
				Factory:  ResourceBackendServer,
			},
			{
				TypeName: "alicloud_slb_domain_extension",
				Factory:  ResourceDomainExtension,
			},
			{
				TypeName: "alicloud_slb_server_group",
				Factory:  ResourceServerGroup,
			},
			{
				TypeName: "alicloud_slb_master_slave_server_group",
				Factory:  ResourceMasterSlaveServerGroup,
			},
			{
				TypeName: "alicloud_slb_rule",
				Factory:  ResourceRule,
			},
			{
				TypeName: "alicloud_slb_acl",
				Factory:  ResourceAcl,
			},
			{
				TypeName: "alicloud_slb_ca_certificate",
				Factory:  ResourceCaCertificate,
			},
			{
				TypeName: "alicloud_slb_server_certificate",
				Factory:  ResourceServerCertificate,
			},
			{
				TypeName: "alicloud_slb_load_balancer",
				Factory:  ResourceLoadBalancer,
			},
			{
				TypeName: "alicloud_slb_tls_cipher_policy",
				Factory:  ResourceTlsCipherPolicy,
			},
			{
				TypeName: "alicloud_slb_acl_entry_attachment",
				Factory:  ResourceAclEntryAttachment,
			},
			{
				TypeName: "alicloud_slb_server_group_server_attachment",
				Factory:  ResourceServerGroupServerAttachment,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_slb_attachments",
				Factory:  DataSourceAttachments,
			},
			{
				TypeName: "alicloud_slb_backend_servers",
				Factory:  DataSourceBackendServers,
			},
			{
				TypeName: "alicloud_slb_listeners",
				Factory:  DataSourceListeners,
			},
			{
				TypeName: "alicloud_slb_rules",
				Factory:  DataSourceRules,
			},
			{
				TypeName: "alicloud_slb_server_groups",
				Factory:  DataSourceServerGroups,
			},
			{
				TypeName: "alicloud_slb_master_slave_server_groups",
				Factory:  DataSourceMasterSlaveServerGroups,
			},
			{
				TypeName: "alicloud_slb_acls",
				Factory:  DataSourceAcls,
			},
			{
				TypeName: "alicloud_slb_server_certificates",
				Factory:  DataSourceServerCertificates,
			},
			{
				TypeName: "alicloud_slb_ca_certificates",
				Factory:  DataSourceCaCertificates,
			},
			{
				TypeName: "alicloud_slb_domain_extensions",
				Factory:  DataSourceDomainExtensions,
			},
			{
				TypeName: "alicloud_slb_zones",
				Factory:  DataSourceZones,
			},
			{
				TypeName: "alicloud_slb_load_balancers",
				Factory:  DataSourceLoadBalancers,
			},
			{
				TypeName: "alicloud_slb_tls_cipher_policies",
				Factory:  DataSourceTlsCipherPolicies,
			},
			{
				TypeName: "alicloud_slbs",
				Factory:  DataSourceLoadBalancers,
			},
		},
	}
}
