package cache

func (this *internalCache) putUnique(pk, uk string) {
	if uk == "" {
		return
	}
	this.hasExtraKeys = true
	this.UniqueToPrimary[uk] = pk
	this.PrimaryToUnique[pk] = uk
}

func (this *internalCache) deleteUnique(pk, uk string) {
	delete(this.UniqueToPrimary, uk)
	delete(this.PrimaryToUnique, pk)
	uk = this.PrimaryToUnique[pk]
	delete(this.UniqueToPrimary, uk)
}
