package notify

import (
	"errors"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8types/go/types/l8reflect"
	"reflect"
)

func CreateNotificationSet(t l8notify.L8NotificationType, serviceName, key string, serviceArea byte, modelType, source string,
	changeCount int, notifySequence uint32) *l8notify.L8NotificationSet {
	notificationSet := &l8notify.L8NotificationSet{}
	notificationSet.ServiceName = serviceName
	notificationSet.ServiceArea = int32(serviceArea)
	notificationSet.ModelType = modelType
	notificationSet.Source = source
	notificationSet.Type = t
	notificationSet.NotificationList = make([]*l8notify.L8Notification, changeCount)
	notificationSet.Sequence = notifySequence
	notificationSet.ModelKey = key
	return notificationSet
}

func CreateAddNotification(any interface{}, serviceName, key string, serviceArea byte, modelType, source string, changeCount int, notifySequence uint32) (*l8notify.L8NotificationSet, error) {
	notificationSet := CreateNotificationSet(l8notify.L8NotificationType_Post, serviceName, key, serviceArea, modelType, source, changeCount, notifySequence)
	obj := object.NewEncode()
	err := obj.Add(any)
	if err != nil {
		return nil, err
	}
	n := &l8notify.L8Notification{}
	n.NewValue = obj.Data()
	notificationSet.NotificationList[0] = n
	return notificationSet, nil
}

func CreateReplaceNotification(old, new interface{}, serviceName, key string, serviceArea byte, modelType, source string, changeCount int, notifySequence uint32) (*l8notify.L8NotificationSet, error) {
	notificationSet := CreateNotificationSet(l8notify.L8NotificationType_Put, serviceName, key, serviceArea, modelType, source, 1, notifySequence)
	oldObj := object.NewEncode()
	err := oldObj.Add(old)
	if err != nil {
		return nil, err
	}

	newObj := object.NewEncode()
	err = newObj.Add(new)
	if err != nil {
		return nil, err
	}

	n := &l8notify.L8Notification{}
	n.OldValue = oldObj.Data()
	n.NewValue = newObj.Data()
	notificationSet.NotificationList[0] = n
	return notificationSet, nil
}

func CreateDeleteNotification(any interface{}, serviceName, key string, serviceArea byte, modelType, source string, changeCount int, notifySequence uint32) (*l8notify.L8NotificationSet, error) {
	notificationSet := CreateNotificationSet(l8notify.L8NotificationType_Delete, serviceName, key, serviceArea, modelType, source, 1, notifySequence)
	obj := object.NewEncode()
	err := obj.Add(any)
	if err != nil {
		return nil, err
	}
	n := &l8notify.L8Notification{}
	n.OldValue = obj.Data()
	notificationSet.NotificationList[0] = n
	return notificationSet, nil
}

func CreateUpdateNotification(changes []*updating.Change, serviceName, key string, serviceArea byte, modelType, source string, changeCount int, notifySequence uint32) (*l8notify.L8NotificationSet, error) {
	notificationSet := CreateNotificationSet(l8notify.L8NotificationType_Patch, serviceName, key, serviceArea, modelType, source, changeCount, notifySequence)
	for i, change := range changes {
		n := &l8notify.L8Notification{}
		n.PropertyId = change.PropertyId()
		if change.OldValue() != nil {
			obj := object.NewEncode()
			err := obj.Add(change.OldValue())
			if err != nil {
				return nil, err
			}
			n.OldValue = obj.Data()
		}
		if change.NewValue() != nil {
			obj := object.NewEncode()
			err := obj.Add(change.NewValue())
			if err != nil {
				return nil, err
			}
			n.NewValue = obj.Data()
		}
		notificationSet.NotificationList[i] = n
	}
	return notificationSet, nil
}

func ItemOf(n *l8notify.L8NotificationSet, resources ifs.IResources) (interface{}, error) {
	switch n.Type {
	case l8notify.L8NotificationType_Put:
		fallthrough
	case l8notify.L8NotificationType_Post:
		obj := object.NewDecode(n.NotificationList[0].NewValue, 0, resources.Registry())
		v, e := obj.Get()
		return v, e
	case l8notify.L8NotificationType_Delete:
		obj := object.NewDecode(n.NotificationList[0].OldValue, 0, resources.Registry())
		v, e := obj.Get()
		return v, e
	case l8notify.L8NotificationType_Patch:

		info, err := resources.Registry().Info(n.ModelType)
		if err != nil {
			return nil, err
		}
		root, err := info.NewInstance()
		if err != nil {
			return nil, err
		}

		node, _ := resources.Introspector().Node(n.ModelType)
		fields, err := resources.Introspector().Decorators().Fields(node, l8reflect.L8DecoratorType_Primary)
		if err != nil {
			panic(err)
		}
		reflect.ValueOf(root).Elem().FieldByName(fields[0]).Set(reflect.ValueOf(n.ModelKey))

		for _, notif := range n.NotificationList {
			p, e := properties.PropertyOf(notif.PropertyId, resources)
			var value interface{}
			if notif.NewValue != nil {
				obj := object.NewDecode(notif.NewValue, 0, resources.Registry())
				v, e1 := obj.Get()
				if e1 != nil {
					return nil, e1
				}
				value = v
			}
			_, _, e = p.Set(root, value)
			if e != nil {
				return nil, e
			}
		}
		return root, nil
	}
	return nil, errors.New("Unknown notification type")
}
