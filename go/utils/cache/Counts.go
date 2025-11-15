package cache

const (
	Total = "Total"
)

func (this *Cache) AddCountFunc(name string, f func(interface{}) (bool, string)) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.iCache.addCountFunc(name, f)
}

func (this *Cache) Counts() map[string]int32 {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make(map[string]int32)
	for k, v := range this.iCache.globalCount.Counts.Counts {
		result[k] = v
	}
	return result
}

func addTotalCount(cache *Cache) {
	cache.AddCountFunc(Total, TotalStat)
}

func TotalStat(any interface{}) (bool, string) {
	return true, ""
}
