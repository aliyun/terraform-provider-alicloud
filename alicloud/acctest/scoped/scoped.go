// Package scoped provides a per-suite test provider with explicit product
// dependencies. Importing it does not import the full provider registry.
package scoped

import (
	"context"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest/checks"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	provider "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/scoped"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Suite struct {
	Provider             *schema.Provider
	SDKProviders         map[string]*schema.Provider
	SDKProviderFactories map[string]func() (*schema.Provider, error)
	DefaultRegionToTest  string
	packages             []conns.ServicePackage
}

// New creates an independent suite. Register dependencies once per product test
// package; individual cases share its client and provider factories.
func New(packages ...conns.ServicePackage) *Suite {
	packages = append([]conns.ServicePackage(nil), packages...)
	p := provider.SDKProvider(packages...)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	return &Suite{Provider: p, SDKProviders: map[string]*schema.Provider{"alicloud": p},
		SDKProviderFactories: map[string]func() (*schema.Provider, error){"alicloud": func() (*schema.Provider, error) { return p, nil }},
		DefaultRegionToTest:  region, packages: packages}
}

func (s *Suite) ProtoV5ProviderFactories(names ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	if len(names) == 0 {
		names = []string{"alicloud"}
	}
	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(names))
	for _, name := range names {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			factory, err := provider.ProtoV5ProviderServerFactory(context.Background(), s.Provider, s.packages...)
			if err != nil {
				return nil, err
			}
			return factory(), nil
		}
	}
	return factories
}
func (s *Suite) PreCheck(t *testing.T) {
	checks.PreCheck(t)
	s.DefaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}
func (s *Suite) PreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	checks.PreCheckWithRegions(t, supported, regions)
	s.DefaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}
func (s *Suite) PreCheckWithAccountSiteType(t *testing.T, site checks.AccountSite) {
	checks.PreCheckWithAccountSiteType(t, site)
}
