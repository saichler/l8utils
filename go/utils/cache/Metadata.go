package cache

const (
	Total = "Total"
)

func (this *Cache) AddMetadataFunc(name string, f func(interface{}) (bool, string)) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.iCache.addMetadataFunc(name, f)
}

func addTotalMetadata(cache *Cache) {
	cache.AddMetadataFunc(Total, TotalStat)
}

func TotalStat(any interface{}) (bool, string) {
	return true, ""
}
