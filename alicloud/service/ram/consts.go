package ram

// defaultTimeout is the default state-refresh wait timeout (minutes) for RAM
// resources, mirroring the root alicloud defaultTimeout constant used by the
// WaitFor* helpers before the product moved into its own subpackage.
const defaultTimeout = 120

// colonSeparated is the id-join separator, mirroring the root alicloud
// colonSeparated constant used to compose multi-segment resource ids.
const colonSeparated = ":"
