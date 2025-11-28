package cache

func (this *internalCache) putUnique(pk, uk string) {
	if uk == "" {
		return
	}
	// Clean up old unique key mapping if it exists and is different
	if oldUk, exists := this.PrimaryToUnique[pk]; exists && oldUk != uk {
		delete(this.UniqueToPrimary, oldUk)
	}
	this.hasExtraKeys = true
	this.UniqueToPrimary[uk] = pk
	this.PrimaryToUnique[pk] = uk
}

func (this *internalCache) deleteUnique(pk, uk string) {
	// If uk not provided, look it up before deleting
	if uk == "" {
		uk = this.PrimaryToUnique[pk]
	}
	delete(this.UniqueToPrimary, uk)
	delete(this.PrimaryToUnique, pk)
}
