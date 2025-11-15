package cache

const (
	Total = "Total"
)

func (this *Cache) AddMetadataFunc(name string, f func(interface{}) (bool, string)) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.iCache.addMetadataFunc(name, f)
}

func (this *Cache) Metadata() map[string]int32 {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make(map[string]int32)
	for k, v := range this.iCache.globalMetadata.KeyCount.Counts {
		result[k] = v
	}
	return result
}

func addTotalMetadata(cache *Cache) {
	cache.AddMetadataFunc(Total, TotalStat)
}

func TotalStat(any interface{}) (bool, string) {
	return true, ""
}
