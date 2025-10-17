package cache

const (
	Total = "Total"
)

func (this *Cache) AddStatFunc(name string, f func(interface{}) bool) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.iCache.addStatsFunc(name, f)
}

func (this *Cache) Stats() map[string]int32 {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make(map[string]int32)
	for k, v := range this.iCache.stats {
		result[k] = v
	}
	return result
}

func addTotalStat(cache *Cache) {
	cache.AddStatFunc(Total, TotalStat)
}

func TotalStat(any interface{}) bool {
	return true
}
