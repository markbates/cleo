package cleo

import (
	"sync"

	"github.com/markbates/plugins"
)

// pooling for memory efficiency
var (
	pluginSlicePool = sync.Pool{
		New: func() interface{} {
			return make(plugins.Plugins, 0, 16)
		},
	}
	stringSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 8)
		},
	}
)

// getPluginSlice gets a plugin slice from the pool
func getPluginSlice() plugins.Plugins {
	return pluginSlicePool.Get().(plugins.Plugins)
}

// putPluginSlice returns a plugin slice to the pool
func putPluginSlice(slice plugins.Plugins) {
	slice = slice[:0] // reset length but keep capacity
	pluginSlicePool.Put(slice)
}

// getStringSlice gets a string slice from the pool
func getStringSlice() []string {
	return stringSlicePool.Get().([]string)
}

// putStringSlice returns a string slice to the pool
func putStringSlice(slice []string) {
	slice = slice[:0] // reset length but keep capacity
	stringSlicePool.Put(slice)
}