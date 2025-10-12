package cache

import "github.com/saichler/l8types/go/ifs"

func (this *Cache) Fetch(start, blockSize int, q ifs.IQuery) []interface{} {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	values := this.iCache.fetch(start, blockSize, q)
	result := make([]interface{}, len(values))
	for i, v := range values {
		result[i] = cloner.Clone(v)
	}
	return result
}
