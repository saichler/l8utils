package cache

const (
	Total = "Total"
)

func (this *Cache) AddCountFunc(name string, f func(interface{}) (bool, string)) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.iCache.addCountFunc(name, f)
}

func addTotalCount(cache *Cache) {
	cache.AddCountFunc(Total, TotalStat)
}

func TotalStat(any interface{}) (bool, string) {
	return true, ""
}
