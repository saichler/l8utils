package cache

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

func (this *Cache) Fetch(start, blockSize int, q ifs.IQuery) ([]interface{}, *l8api.L8MetaData) {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	values, metadata := this.iCache.fetch(start, blockSize, q)
	result := make([]interface{}, len(values))
	for i, v := range values {
		result[i] = cloner.Clone(v)
	}

	if q.Page() == 0 {
		metadataClone := cloner.Clone(metadata).(*l8api.L8MetaData)
		return result, metadataClone
	}
	return result, nil
}
