package cache

import (
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8utils/go/utils/notify"
)

func (this *Cache) createNotificationSet(t l8notify.L8NotificationType, key string, changeCount int) *l8notify.L8NotificationSet {
	defer func() { this.notifySequence++ }()
	return notify.CreateNotificationSet(t, this.serviceName, key, this.serviceArea, this.modelType, this.source, changeCount, this.notifySequence)
}

func (this *Cache) createAddNotification(any interface{}, key string) (*l8notify.L8NotificationSet, error) {
	defer func() { this.notifySequence++ }()
	return notify.CreateAddNotification(any, this.serviceName, key, this.serviceArea, this.modelType, this.source, 1, this.notifySequence)
}

func (this *Cache) createSyncNotification(any interface{}, key string) (*l8notify.L8NotificationSet, error) {
	defer func() { this.notifySequence++ }()
	return notify.CreateSyncNotification(any, this.serviceName, key, this.serviceArea, this.modelType, this.source, 1, this.notifySequence)
}

func (this *Cache) createReplaceNotification(old, new interface{}, key string) (*l8notify.L8NotificationSet, error) {
	defer func() { this.notifySequence++ }()
	return notify.CreateReplaceNotification(old, new, this.serviceName, key, this.serviceArea, this.modelType, this.source, 1, this.notifySequence)
}

func (this *Cache) createDeleteNotification(any interface{}, key string) (*l8notify.L8NotificationSet, error) {
	defer func() { this.notifySequence++ }()
	return notify.CreateDeleteNotification(any, this.serviceName, key, this.serviceArea, this.modelType, this.source, 1, this.notifySequence)
}

func (this *Cache) createUpdateNotification(changes []*updating.Change, key string) (*l8notify.L8NotificationSet, error) {
	defer func() { this.notifySequence++ }()
	return notify.CreateUpdateNotification(changes, this.serviceName, key, this.serviceArea, this.modelType, this.source, len(changes), this.notifySequence)
}
