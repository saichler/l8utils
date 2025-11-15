package cache

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

func (this *Cache) Fetch(start, blockSize int, q ifs.IQuery) ([]interface{}, *l8api.L8Counts) {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	values, counts := this.iCache.fetch(start, blockSize, q)
	result := make([]interface{}, len(values))
	for i, v := range values {
		result[i] = cloner.Clone(v)
	}
	countsClone := cloner.Clone(counts).(*l8api.L8Counts)
	return result, countsClone
}
