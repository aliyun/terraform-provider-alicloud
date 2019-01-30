package alicloud

import (
	"regexp"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDnsDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsDomainsRead,

		Schema: map[string]*schema.Schema{
			"domain_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ali_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"version_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ali_domain": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"puny_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudDnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := &dns.DescribeDomainsArgs{}

	var allDomains []dns.DomainType
	pagination := getPagination(1, 50)
	for {
		args.Pagination = pagination
		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomains(args)
		})
		if err != nil {
			return err
		}
		domains, _ := raw.([]dns.DomainType)
		allDomains = append(allDomains, domains...)

		if len(domains) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	var filteredDomains []dns.DomainType

	for _, domain := range allDomains {
		if v, ok := d.GetOk("ali_domain"); ok && domain.AliDomain != v.(bool) {
			continue
		}

		if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" && domain.InstanceId != v.(string) {
			continue
		}

		if v, ok := d.GetOk("version_code"); ok && v.(string) != "" && domain.VersionCode != v.(string) {
			continue
		}

		if v, ok := d.GetOk("domain_name_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(domain.DomainName) {
				continue
			}
		}

		if v, ok := d.GetOk("group_name_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(domain.GroupName) {
				continue
			}
		}

		filteredDomains = append(filteredDomains, domain)
	}

	return domainsDecriptionAttributes(d, filteredDomains, meta)
}

func domainsDecriptionAttributes(d *schema.ResourceData, domainTypes []dns.DomainType, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, domain := range domainTypes {
		mapping := map[string]interface{}{
			"domain_id":    domain.DomainId,
			"domain_name":  domain.DomainName,
			"group_id":     domain.GroupId,
			"group_name":   domain.GroupName,
			"ali_domain":   domain.AliDomain,
			"instance_id":  domain.InstanceId,
			"version_code": domain.VersionCode,
			"puny_code":    domain.PunyCode,
			"dns_servers":  domain.DnsServers.DnsServer,
		}
		names = append(names, domain.DomainName)
		ids = append(ids, domain.DomainId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("domains", s); err != nil {
		return err
	}
	if err := d.Set("names", names); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
