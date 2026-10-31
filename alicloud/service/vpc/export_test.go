package vpc

// RouteEntryTaskLockForTest exports the unexported package-level
// routeEntryTaskLock (see resource_alicloud_route_entry.go) to the external
// vpc_test package so the migrated route-entry sweep/acceptance helpers share
// the provider's real CreateRouteEntry/DeleteRouteEntry serialization mutex
// instead of a fresh, disconnected one. This is the standard Go export_test.go
// bridge: a test-variant file in the internal package (package vpc) exposes an
// unexported symbol to the external test package (vpc_test) within the same
// test binary.
//
// It takes the ADDRESS of the var (&routeEntryTaskLock) so callers bind the
// shared *sync.Mutex pointer rather than copying the mutex value — a
// sync.Mutex must not be copied after first use. No production behavior changes;
// the production var and its Lock/Unlock call sites in
// resourceAliyunRouteEntryCreate/DeleteRouteEntry are untouched.
var RouteEntryTaskLockForTest = &routeEntryTaskLock
