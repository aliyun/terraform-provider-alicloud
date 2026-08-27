package intercept

// ChainOf merges the three sources in fixed order: global, declared,
// registry[name]. An empty result makes every wrapper the identity function.
func ChainOf(name string, declared []Interceptor) []Interceptor {
	chain := make([]Interceptor, 0, len(global)+len(declared)+len(registry[name]))
	chain = append(chain, global...)
	chain = append(chain, declared...)
	chain = append(chain, registry[name]...)
	return chain
}
