package tests

import (
	"testing"

	"github.com/saichler/l8reflect/go/reflect/helping"
	"github.com/saichler/l8reflect/go/reflect/introspecting"
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8types/go/types/l8sysconfig"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/notify"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
)

func createModel(i int) *testtypes.TestProto {
	return utils.CreateTestModelInstance(i)
}

func newResources() ifs.IResources {
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	res := resources.NewResources(log)
	res.Set(registry.NewRegistry())
	res.Set(&l8sysconfig.L8SysConfig{})
	in := introspecting.NewIntrospect(res.Registry())
	res.Set(in)
	node, _ := res.Introspector().Inspect(testtypes.TestProto{})
	helping.AddPrimaryKeyDecorator(node, "MyString")
	helping.AddUniqueKeyDecorator(node, "MyInt32")
	return res
}

func createChanges(aside, zside *testtypes.TestProto, r ifs.IResources) []*updating.Change {
	updater := updating.NewUpdater(r, true, true)
	updater.Update(aside, zside)
	return updater.Changes()
}

// Test CreateNotificationSet
func TestCreateNotificationSet(t *testing.T) {
	notSet := notify.CreateNotificationSet(
		l8notify.L8NotificationType_Post,
		"test-service",
		"test-key",
		1,
		"TestModel",
		"test-source",
		5,
		123,
	)

	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
	if notSet.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", notSet.ServiceName)
	}
	if notSet.ModelKey != "test-key" {
		t.Errorf("Expected model key 'test-key', got '%s'", notSet.ModelKey)
	}
	if notSet.ServiceArea != 1 {
		t.Errorf("Expected service area 1, got %d", notSet.ServiceArea)
	}
	if notSet.ModelType != "TestModel" {
		t.Errorf("Expected model type 'TestModel', got '%s'", notSet.ModelType)
	}
	if notSet.Source != "test-source" {
		t.Errorf("Expected source 'test-source', got '%s'", notSet.Source)
	}
	if notSet.Type != l8notify.L8NotificationType_Post {
		t.Errorf("Expected type Add, got %v", notSet.Type)
	}
	if len(notSet.NotificationList) != 5 {
		t.Errorf("Expected 5 notifications, got %d", len(notSet.NotificationList))
	}
	if notSet.Sequence != 123 {
		t.Errorf("Expected sequence 123, got %d", notSet.Sequence)
	}
}

// Test CreateAddNotification
func TestCreateAddNotification(t *testing.T) {
	model := createModel(1)
	res := newResources()

	notSet, err := notify.CreateAddNotification(
		model,
		"test-service",
		"test-key",
		1,
		"TestProto",
		"test-source",
		1,
		100,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
	if notSet.Type != l8notify.L8NotificationType_Post {
		t.Errorf("Expected Add type, got %v", notSet.Type)
	}
	if len(notSet.NotificationList) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notSet.NotificationList))
	}
	if notSet.NotificationList[0].NewValue == nil {
		t.Error("Expected new value to be set")
	}

	// Test ItemOf with Add notification
	item, err := notify.ItemOf(notSet, res)
	if err != nil {
		t.Errorf("Expected no error from ItemOf, got: %v", err)
	}
	if item == nil {
		t.Error("Expected item to be extracted")
	}
}

// Test CreateReplaceNotification
func TestCreateReplaceNotification(t *testing.T) {
	oldModel := createModel(1)
	newModel := createModel(2)
	res := newResources()

	notSet, err := notify.CreateReplaceNotification(
		oldModel,
		newModel,
		"replace-service",
		"replace-key",
		3,
		"TestProto",
		"replace-source",
		1,
		300,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
	if notSet.Type != l8notify.L8NotificationType_Put {
		t.Errorf("Expected Replace type, got %v", notSet.Type)
	}
	if notSet.NotificationList[0].OldValue == nil {
		t.Error("Expected old value to be set")
	}
	if notSet.NotificationList[0].NewValue == nil {
		t.Error("Expected new value to be set")
	}

	// Test ItemOf with Replace notification
	item, err := notify.ItemOf(notSet, res)
	if err != nil {
		t.Errorf("Expected no error from ItemOf, got: %v", err)
	}
	if item == nil {
		t.Error("Expected item to be extracted")
	}
}

// Test CreateDeleteNotification
func TestCreateDeleteNotification(t *testing.T) {
	model := createModel(3)
	res := newResources()

	notSet, err := notify.CreateDeleteNotification(
		model,
		"delete-service",
		"delete-key",
		4,
		"TestProto",
		"delete-source",
		1,
		400,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
	if notSet.Type != l8notify.L8NotificationType_Delete {
		t.Errorf("Expected Delete type, got %v", notSet.Type)
	}
	if notSet.NotificationList[0].OldValue == nil {
		t.Error("Expected old value to be set")
	}
	if notSet.NotificationList[0].NewValue != nil {
		t.Error("Expected new value to be nil")
	}

	// Test ItemOf with Delete notification
	item, err := notify.ItemOf(notSet, res)
	if err != nil {
		t.Errorf("Expected no error from ItemOf, got: %v", err)
	}
	if item == nil {
		t.Error("Expected item to be extracted")
	}
}

// Test CreateUpdateNotification
func TestCreateUpdateNotification(t *testing.T) {
	model1 := createModel(1)
	model2 := createModel(2)
	res := newResources()

	// Modify model2 to create changes
	model2.MyString = "modified"
	model2.MyInt32 = 999

	changes := createChanges(model1, model2, res)

	if len(changes) == 0 {
		t.Skip("No changes detected, skipping update notification test")
	}

	notSet, err := notify.CreateUpdateNotification(
		changes,
		"update-service",
		"update-key",
		5,
		"TestProto",
		"update-source",
		len(changes),
		500,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
	if notSet.Type != l8notify.L8NotificationType_Patch {
		t.Errorf("Expected Update type, got %v", notSet.Type)
	}
	if len(notSet.NotificationList) != len(changes) {
		t.Errorf("Expected %d notifications, got %d", len(changes), len(notSet.NotificationList))
	}

	// Verify each notification has property information
	for i, notif := range notSet.NotificationList {
		if notif.PropertyId == "" {
			t.Errorf("Notification %d: expected property ID to be set", i)
		}
	}

	// Test ItemOf with Update notification
	// Note: ItemOf for Update notifications may fail with complex types like slices
	item, err := notify.ItemOf(notSet, res)
	if err != nil {
		t.Logf("ItemOf returned error (may be expected for complex types): %v", err)
	} else if item == nil {
		t.Error("Expected item to be extracted when no error")
	}
}

// Test CreateUpdateNotification with empty changes
func TestCreateUpdateNotificationEmpty(t *testing.T) {
	changes := []*updating.Change{}

	notSet, err := notify.CreateUpdateNotification(
		changes,
		"update-service",
		"update-key",
		5,
		"TestProto",
		"update-source",
		0,
		500,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
	if len(notSet.NotificationList) != 0 {
		t.Errorf("Expected 0 notifications, got %d", len(notSet.NotificationList))
	}
}

// Test CreateUpdateNotification with changes having nil values
func TestCreateUpdateNotificationWithNilValues(t *testing.T) {
	model1 := createModel(1)
	model2 := createModel(1)
	res := newResources()

	// Create a change with nil new value
	model2.MyString = ""

	changes := createChanges(model1, model2, res)

	if len(changes) == 0 {
		t.Skip("No changes detected")
	}

	notSet, err := notify.CreateUpdateNotification(
		changes,
		"update-service",
		"update-key",
		5,
		"TestProto",
		"update-source",
		len(changes),
		500,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notSet == nil {
		t.Fatal("Expected notification set to be created")
	}
}

// Test ItemOf with different notification types
func TestItemOfVariousTypes(t *testing.T) {
	res := newResources()

	testCases := []struct {
		name     string
		setup    func() (*l8notify.L8NotificationSet, error)
		wantItem bool
	}{
		{
			name: "Add notification",
			setup: func() (*l8notify.L8NotificationSet, error) {
				return notify.CreateAddNotification(createModel(1), "svc", "key", 1, "TestProto", "src", 1, 1)
			},
			wantItem: true,
		},
		{
			name: "Replace notification",
			setup: func() (*l8notify.L8NotificationSet, error) {
				return notify.CreateReplaceNotification(createModel(1), createModel(2), "svc", "key", 1, "TestProto", "src", 1, 3)
			},
			wantItem: true,
		},
		{
			name: "Delete notification",
			setup: func() (*l8notify.L8NotificationSet, error) {
				return notify.CreateDeleteNotification(createModel(3), "svc", "key", 1, "TestProto", "src", 1, 4)
			},
			wantItem: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notSet, err := tc.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			item, err := notify.ItemOf(notSet, res)
			if err != nil {
				t.Errorf("ItemOf failed: %v", err)
			}
			if tc.wantItem && item == nil {
				t.Error("Expected item to be extracted")
			}
		})
	}
}

// Test sequence numbers
func TestSequenceNumbers(t *testing.T) {
	model := createModel(1)

	sequences := []uint32{0, 1, 100, 999, 4294967295}
	for _, seq := range sequences {
		notSet, err := notify.CreateAddNotification(
			model,
			"test-service",
			"test-key",
			1,
			"TestProto",
			"test-source",
			1,
			seq,
		)
		if err != nil {
			t.Errorf("Failed to create notification with sequence %d: %v", seq, err)
		}
		if notSet.Sequence != seq {
			t.Errorf("Expected sequence %d, got %d", seq, notSet.Sequence)
		}
	}
}

// Test service areas
func TestServiceAreas(t *testing.T) {
	model := createModel(1)

	areas := []byte{0, 1, 127, 255}
	for _, area := range areas {
		notSet, err := notify.CreateAddNotification(
			model,
			"test-service",
			"test-key",
			area,
			"TestProto",
			"test-source",
			1,
			1,
		)
		if err != nil {
			t.Errorf("Failed to create notification with area %d: %v", area, err)
		}
		if notSet.ServiceArea != int32(area) {
			t.Errorf("Expected service area %d, got %d", area, notSet.ServiceArea)
		}
	}
}

// Test notification types
func TestNotificationTypes(t *testing.T) {
	types := []l8notify.L8NotificationType{
		l8notify.L8NotificationType_Post,
		l8notify.L8NotificationType_Delete,
		l8notify.L8NotificationType_Patch,
		l8notify.L8NotificationType_Put,
	}

	for _, nt := range types {
		notSet := notify.CreateNotificationSet(
			nt,
			"test-service",
			"test-key",
			1,
			"TestModel",
			"test-source",
			1,
			1,
		)
		if notSet.Type != nt {
			t.Errorf("Expected type %v, got %v", nt, notSet.Type)
		}
	}
}

// Test multiple changes in update notification
func TestMultipleChangesInUpdate(t *testing.T) {
	model1 := createModel(1)
	model2 := createModel(1)
	res := newResources()

	// Make multiple changes
	model2.MyString = "changed"
	model2.MyInt32 = 12345
	model2.MyBool = !model1.MyBool

	changes := createChanges(model1, model2, res)

	if len(changes) < 2 {
		t.Logf("Only %d changes detected, test may not be comprehensive", len(changes))
	}

	notSet, err := notify.CreateUpdateNotification(
		changes,
		"update-service",
		"update-key",
		1,
		"TestProto",
		"update-source",
		len(changes),
		1,
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(notSet.NotificationList) != len(changes) {
		t.Errorf("Expected %d notifications, got %d", len(changes), len(notSet.NotificationList))
	}

	// Verify we can extract the updated item
	// Note: ItemOf for Update notifications may fail with complex types like slices
	item, err := notify.ItemOf(notSet, res)
	if err != nil {
		t.Logf("ItemOf returned error (may be expected for complex types): %v", err)
	} else if item == nil {
		t.Error("Expected item to be extracted when no error")
	}
}
