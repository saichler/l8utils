package cache

func (this *Cache) Collect(f func(interface{}) (bool, interface{})) map[string]interface{} {
	result := map[string]interface{}{}
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	if this.cacheEnabled() {
		for k, v := range this.iCache.cache {
			vClone := cloner.Clone(v)
			ok, elem := f(vClone)
			if ok {
				result[k] = elem
			}
		}
		return result
	}
	return this.store.Collect(f)
}
