# Product-scoped acceptance providers

Use one suite per product test package. The selected products must cover both the
resource under test and the resources/data sources in its HCL fixtures:

```go
var testSuite = scoped.New(
    fcv3.ServicePackage(),
    vpc.ServicePackage(),
)
```

Use `testSuite.ProtoV5ProviderFactories()` in `resource.TestCase`. It serves both
SDK and Framework declarations from that selection through the production mux.
Use `testSuite.Provider.Meta()` in API checks so they observe the configured SDK
client shared with Framework resources. Independent suites do not change a global
provider selection. `testSuite.SDKProviderFactories` remains available for helpers
that require the SDK-only factory type.

Import generic assertions and HCL builders from `alicloud/acctest/helpers`, and
prechecks from `alicloud/acctest/checks` or the suite's methods. Importing the legacy
`alicloud/acctest`, root `alicloud`, or the full `provider/registry` restores the
complete compilation graph even if a test only executes one resource.

The FCV3 suite demonstrates a shared selection, a dependency-boundary regression,
and a fixture-registration coverage check. When adding cross-product fixtures,
update that shared selection. Dependencies imported internally by product code
still contribute to compilation time.

Existing `alicloud.Provider()`, `provider.ProtoV5ProviderServerFactory`,
`framework.NewProvider`, and `alicloud/acctest` retain their full selections.
