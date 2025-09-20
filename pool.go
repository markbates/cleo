package cleo

import (
	"sync"

	"github.com/markbates/plugins"
)

// pooling for memory efficiency
var (
	pluginSlicePool = sync.Pool{
		New: func() any {
			slice := make(plugins.Plugins, 0, 16)
			return &slice
		},
	}
)

// getPluginSlice gets a plugin slice from the pool
func getPluginSlice() plugins.Plugins {
	slice := pluginSlicePool.Get().(*plugins.Plugins)
	return *slice
}

// putPluginSlice returns a plugin slice to the pool
func putPluginSlice(slice plugins.Plugins) {
	slice = slice[:0] // reset length but keep capacity
	pluginSlicePool.Put(&slice)
}