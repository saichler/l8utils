package cache

import (
	"github.com/saichler/l8types/go/types/l8notify"
)

func (this *Cache) Put(v interface{}, createNotification bool) (*l8notify.L8NotificationSet, error) {
	//Seems that the post is handling also a put situation, where the item is replaced
	return this.Post(v, createNotification)
}
