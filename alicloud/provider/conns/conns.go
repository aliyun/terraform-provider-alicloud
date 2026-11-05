// Package conns defines the types a service package fills in to state what it
// offers: one ServicePackage per product, plus the per-category item types (SDK v2
// in sdkv2.go, Framework in framework.go). Separate from package registry to avoid
// the cycle registry -> product -> registry. The name mirrors
// terraform-provider-aws; the client itself is alicloud/connectivity.
package conns

type ServicePackage struct {
	Name string // API product code, e.g. "Ims", "Vpc"

	SDKResources       []SDKResource
	SDKDataSources     []SDKDataSource
	Resources          []Resource
	DataSources        []DataSource
	Functions          []Function
	EphemeralResources []EphemeralResource
	ListResources      []ListResource
	Actions            []Action
}
