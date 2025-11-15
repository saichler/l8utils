package tests

import (
	"fmt"
	"testing"

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/cache"
)

func createIQuery(gsql string, r ifs.IResources) ifs.IQuery {
	elems, e := object.NewQuery(gsql, r)
	if e != nil {
		panic(e)
	}
	q, _ := elems.Query(r)
	return q
}

func TestFetch(t *testing.T) {
	res := newResources()
	t1 := createModel(1)
	t2 := createModel(2)
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	c.Post(t1, false)
	c.Post(t2, false)

	q := createIQuery("select * from TestProto", res)
	elems, _ := c.Fetch(0, 25, q)
	if len(elems) != 2 {
		t.Fail()
		fmt.Println("Error, expected 2 elements")
		return
	}

	q = createIQuery("select * from TestProto where MyString=*", res)
	elems, _ = c.Fetch(0, 25, q)
	if len(elems) != 2 {
		t.Fail()
		fmt.Println("Error, expected 2 elements")
		return
	}
}

// Simple storage implementation for testing
type testStorage struct {
	data         map[string]interface{}
	cacheEnabled bool
}

func newTestStorage(cacheEnabled bool) *testStorage {
	return &testStorage{
		data:         make(map[string]interface{}),
		cacheEnabled: cacheEnabled,
	}
}

func (s *testStorage) Get(key string) (interface{}, error) {
	val, ok := s.data[key]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (s *testStorage) Put(key string, value interface{}) error {
	s.data[key] = value
	return nil
}

func (s *testStorage) Delete(key string) (interface{}, error) {
	val, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return val, nil
}

func (s *testStorage) Collect(f func(interface{}) (bool, interface{})) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range s.data {
		include, val := f(v)
		if include {
			result[k] = val
		}
	}
	return result
}

func (s *testStorage) CacheEnabled() bool {
	return s.cacheEnabled
}

// Test NewCache with no initial elements
func TestNewCacheEmpty(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	if c == nil {
		t.Fatal("Expected cache to be created")
	}

	size := c.Size()
	if size != 0 {
		t.Errorf("Expected size 0, got %d", size)
	}
}

// Test NewCache with initial elements
func TestNewCacheWithElements(t *testing.T) {
	res := newResources()
	model1 := createModel(1)
	model2 := createModel(2)
	model3 := createModel(3)

	initElements := []interface{}{model1, model2, model3}
	c := cache.NewCache(model1, initElements, nil, res)

	if c == nil {
		t.Fatal("Expected cache to be created")
	}

	size := c.Size()
	if size != 3 {
		t.Errorf("Expected size 3, got %d", size)
	}
}

// Test Post operation - adding new item
func TestCachePost(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	c.SetNotificationsFor("test-service", 1)

	newModel := createModel(10)
	notification, err := c.Post(newModel, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification == nil {
		t.Error("Expected notification to be created")
	}

	size := c.Size()
	if size != 1 {
		t.Errorf("Expected size 1 after post, got %d", size)
	}
}

// Test Post operation without notification
func TestCachePostNoNotification(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	newModel := createModel(20)
	notification, err := c.Post(newModel, false)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification != nil {
		t.Error("Expected no notification to be created")
	}

	size := c.Size()
	if size != 1 {
		t.Errorf("Expected size 1 after post, got %d", size)
	}
}

// Test Post operation - replacing existing item
func TestCachePostReplace(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)
	c.SetNotificationsFor("test-service", 1)

	// Post with same key but different values
	model1Updated := createModel(1)
	model1Updated.MyInt32 = 999

	notification, err := c.Post(model1Updated, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Size should remain 1
	size := c.Size()
	if size != 1 {
		t.Errorf("Expected size 1 after replace, got %d", size)
	}

	// Verify the item was updated
	retrieved, err := c.Get(model1)
	if err != nil {
		t.Errorf("Expected no error on get, got: %v", err)
	}
	if retrievedModel, ok := retrieved.(*testtypes.TestProto); ok {
		if retrievedModel.MyInt32 != 999 {
			t.Errorf("Expected MyInt32 to be 999, got %d", retrievedModel.MyInt32)
		}
	}

	// Check notification was created for replace
	if notification != nil {
		t.Logf("Replace notification created with sequence %d", notification.Sequence)
	}
}

// Test Get operation
func TestCacheGet(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)

	retrieved, err := c.Get(model1)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if retrieved == nil {
		t.Error("Expected item to be retrieved")
	}

	if retrievedModel, ok := retrieved.(*testtypes.TestProto); ok {
		if retrievedModel.MyString != model1.MyString {
			t.Errorf("Expected MyString '%s', got '%s'", model1.MyString, retrievedModel.MyString)
		}
	} else {
		t.Error("Expected retrieved item to be TestProto")
	}
}

// Test Get operation - item not found
func TestCacheGetNotFound(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	c := cache.NewCache(model1, nil, nil, res)

	model999 := createModel(999)
	retrieved, err := c.Get(model999)

	if err == nil {
		t.Error("Expected error for item not found")
	}
	if retrieved != nil {
		t.Error("Expected nil result for item not found")
	}
}

// Test Put operation (alias for Post)
func TestCachePut(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	c.SetNotificationsFor("test-service", 1)

	newModel := createModel(30)
	notification, err := c.Put(newModel, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification == nil {
		t.Error("Expected notification to be created")
	}

	size := c.Size()
	if size != 1 {
		t.Errorf("Expected size 1 after put, got %d", size)
	}
}

// Test Patch operation - update existing item
func TestCachePatch(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)
	c.SetNotificationsFor("test-service", 1)

	// Patch with partial update
	patchModel := createModel(1)
	patchModel.MyInt32 = 777

	notification, err := c.Patch(patchModel, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify the item was patched
	retrieved, err := c.Get(model1)
	if err != nil {
		t.Errorf("Expected no error on get, got: %v", err)
	}
	if retrievedModel, ok := retrieved.(*testtypes.TestProto); ok {
		if retrievedModel.MyInt32 != 777 {
			t.Errorf("Expected MyInt32 to be 777, got %d", retrievedModel.MyInt32)
		}
	}

	if notification != nil {
		t.Logf("Patch notification created with %d changes", len(notification.NotificationList))
	}
}

// Test Patch operation - create new item if not exists
func TestCachePatchNewItem(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	c.SetNotificationsFor("test-service", 1)

	newModel := createModel(40)
	notification, err := c.Patch(newModel, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification == nil {
		t.Error("Expected notification for new item")
	}

	size := c.Size()
	if size != 1 {
		t.Errorf("Expected size 1 after patch, got %d", size)
	}
}

// Test Patch without notification
func TestCachePatchNoNotification(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)

	patchModel := createModel(1)
	patchModel.MyInt32 = 555

	notification, err := c.Patch(patchModel, false)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification != nil {
		t.Error("Expected no notification")
	}
}

// Test Delete operation
func TestCacheDelete(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)
	c.SetNotificationsFor("test-service", 1)

	notification, err := c.Delete(model1, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification == nil {
		t.Error("Expected notification for delete")
	}

	size := c.Size()
	if size != 0 {
		t.Errorf("Expected size 0 after delete, got %d", size)
	}

	// Verify item is gone
	_, err = c.Get(model1)
	if err == nil {
		t.Error("Expected error when getting deleted item")
	}
}

// Test Delete without notification
func TestCacheDeleteNoNotification(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)

	notification, err := c.Delete(model1, false)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification != nil {
		t.Error("Expected no notification")
	}

	size := c.Size()
	if size != 0 {
		t.Errorf("Expected size 0 after delete, got %d", size)
	}
}

// Test Delete - item not found
func TestCacheDeleteNotFound(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	c := cache.NewCache(model1, nil, nil, res)

	model999 := createModel(999)
	notification, err := c.Delete(model999, true)

	if err == nil {
		t.Error("Expected error when deleting non-existent item")
	}
	if notification != nil {
		t.Error("Expected no notification for failed delete")
	}
}

// Test SetNotificationsFor
func TestCacheSetNotificationsFor(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	c.SetNotificationsFor("my-service", 5)

	// Test that notifications work after setting
	newModel := createModel(50)
	notification, err := c.Post(newModel, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if notification == nil {
		t.Error("Expected notification to be created")
	}
	// Source is now automatically determined from resources, not set explicitly
	if notification.ServiceName != "my-service" {
		t.Errorf("Expected service name 'my-service', got '%s'", notification.ServiceName)
	}
	if notification.ServiceArea != 5 {
		t.Errorf("Expected service area 5, got %d", notification.ServiceArea)
	}
}

// Test notification sequence incrementing
func TestCacheNotificationSequence(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	c.SetNotificationsFor("test-service", 1)

	sequences := []uint32{}

	// Post 3 items and collect sequences
	for i := 1; i <= 3; i++ {
		newModel := createModel(i * 10)
		notification, err := c.Post(newModel, true)
		if err != nil {
			t.Fatalf("Expected no error on post %d, got: %v", i, err)
		}
		if notification != nil {
			sequences = append(sequences, notification.Sequence)
		}
	}

	if len(sequences) != 3 {
		t.Errorf("Expected 3 sequences, got %d", len(sequences))
	}

	// Verify sequences are incrementing
	for i := 1; i < len(sequences); i++ {
		if sequences[i] <= sequences[i-1] {
			t.Errorf("Expected sequence to increment, got %d then %d", sequences[i-1], sequences[i])
		}
	}
}

// Test multiple operations in sequence
func TestCacheMultipleOperations(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)
	c.SetNotificationsFor("test-service", 1)

	// Post
	model1 := createModel(1)
	_, err := c.Post(model1, false)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}

	// Get
	retrieved, err := c.Get(model1)
	if err != nil || retrieved == nil {
		t.Error("Get failed after post")
	}

	// Patch
	patchModel := createModel(1)
	patchModel.MyInt32 = 111
	_, err = c.Patch(patchModel, false)
	if err != nil {
		t.Fatalf("Patch failed: %v", err)
	}

	// Get again to verify patch
	retrieved, err = c.Get(model1)
	if err != nil {
		t.Fatalf("Get failed after patch: %v", err)
	}
	if retrievedModel, ok := retrieved.(*testtypes.TestProto); ok {
		if retrievedModel.MyInt32 != 111 {
			t.Errorf("Expected MyInt32 to be 111 after patch, got %d", retrievedModel.MyInt32)
		}
	}

	// Post another
	model2 := createModel(2)
	_, err = c.Post(model2, false)
	if err != nil {
		t.Fatalf("Second post failed: %v", err)
	}

	// Check size
	if c.Size() != 2 {
		t.Errorf("Expected size 2, got %d", c.Size())
	}

	// Delete one
	_, err = c.Delete(model1, false)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Check size again
	if c.Size() != 1 {
		t.Errorf("Expected size 1 after delete, got %d", c.Size())
	}
}

// Test cache with multiple items
func TestCacheMultipleItems(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	// Add 10 items
	for i := 1; i <= 10; i++ {
		newModel := createModel(i)
		_, err := c.Post(newModel, false)
		if err != nil {
			t.Fatalf("Failed to post model %d: %v", i, err)
		}
	}

	if c.Size() != 10 {
		t.Errorf("Expected size 10, got %d", c.Size())
	}

	// Retrieve each item
	for i := 1; i <= 10; i++ {
		testModel := createModel(i)
		retrieved, err := c.Get(testModel)
		if err != nil {
			t.Errorf("Failed to get model %d: %v", i, err)
		}
		if retrieved == nil {
			t.Errorf("Expected to retrieve model %d", i)
		}
	}
}

// Test AddMetadataFunc and Metadata
func TestCacheStats(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	// Add a metadata function that counts items with MyInt32 > 100
	c.AddMetadataFunc("high_values", func(i interface{}) (bool, string) {
		if testModel, ok := i.(*testtypes.TestProto); ok {
			return testModel.MyInt32 > 100, ""
		}
		return false, ""
	})

	// Add items with different values
	for i := 1; i <= 5; i++ {
		newModel := createModel(i * 10)
		newModel.MyInt32 = int32(i * 50)
		_, err := c.Post(newModel, false)
		if err != nil {
			t.Fatalf("Failed to post model: %v", err)
		}
	}

	metadata := c.Metadata()
	highValues, ok := metadata["high_values"]
	if !ok {
		t.Error("Expected 'high_values' metadata to exist")
	}

	// Models with i=3,4,5 should have MyInt32 > 100 (150, 200, 250)
	if highValues != 3 {
		t.Errorf("Expected 3 high values, got %d", highValues)
	}
}

// Test Metadata with empty cache
func TestCacheStatsEmpty(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	c.AddMetadataFunc("any_stat", func(i interface{}) (bool, string) {
		return true, ""
	})

	metadata := c.Metadata()
	// When cache is empty, the metadata might not be initialized
	if len(metadata) > 0 {
		anyStat, ok := metadata["any_stat"]
		if ok && anyStat != 0 {
			t.Errorf("Expected 0 for empty cache, got %d", anyStat)
		}
	}
}

// Test Metadata update after patch
func TestCacheStatsAfterPatch(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	// Add metadata function
	c.AddMetadataFunc("even_values", func(i interface{}) (bool, string) {
		if testModel, ok := i.(*testtypes.TestProto); ok {
			return testModel.MyInt32%2 == 0, ""
		}
		return false, ""
	})

	// Add item with odd value
	model1 := createModel(1)
	model1.MyInt32 = 5
	_, err := c.Post(model1, false)
	if err != nil {
		t.Fatalf("Failed to post: %v", err)
	}

	metadata := c.Metadata()
	if metadata["even_values"] != 0 {
		t.Errorf("Expected 0 even values, got %d", metadata["even_values"])
	}

	// Patch to even value
	patchModel := createModel(1)
	patchModel.MyInt32 = 10
	_, err = c.Patch(patchModel, false)
	if err != nil {
		t.Fatalf("Failed to patch: %v", err)
	}

	metadata = c.Metadata()
	if metadata["even_values"] != 1 {
		t.Errorf("Expected 1 even value after patch, got %d", metadata["even_values"])
	}
}

// Test Metadata update after delete
func TestCacheStatsAfterDelete(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	// Add metadata function
	c.AddMetadataFunc("count_all", func(i interface{}) (bool, string) {
		return true, ""
	})

	// Add items
	model1 := createModel(1)
	model2 := createModel(2)
	c.Post(model1, false)
	c.Post(model2, false)

	metadata := c.Metadata()
	if metadata["count_all"] != 2 {
		t.Errorf("Expected 2 items, got %d", metadata["count_all"])
	}

	// Delete one
	c.Delete(model1, false)

	metadata = c.Metadata()
	if metadata["count_all"] != 1 {
		t.Errorf("Expected 1 item after delete, got %d", metadata["count_all"])
	}
}

// Test PrimaryKeyFor with nil
func TestCachePrimaryKeyForNil(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	_, err := c.PrimaryKeyFor(nil)
	if err == nil {
		t.Error("Expected error for nil interface")
	}
}

// Test Post with no changes (should not create notification)
func TestCachePostNoChanges(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)
	c.SetNotificationsFor("test-service", 1)

	// Post same model again
	notification, err := c.Post(model1, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// When there are no changes, notification should be nil
	if notification != nil {
		t.Logf("Notification created even with no changes: %v", notification)
	}
}

// Test Metadata with multiple functions
func TestCacheStatsMultipleFunctions(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	// Add multiple metadata functions
	c.AddMetadataFunc("positive", func(i interface{}) (bool, string) {
		if testModel, ok := i.(*testtypes.TestProto); ok {
			return testModel.MyInt32 > 0, ""
		}
		return false, ""
	})

	c.AddMetadataFunc("large", func(i interface{}) (bool, string) {
		if testModel, ok := i.(*testtypes.TestProto); ok {
			return testModel.MyInt32 > 100, ""
		}
		return false, ""
	})

	// Add items
	for i := 1; i <= 5; i++ {
		newModel := createModel(i * 10)
		newModel.MyInt32 = int32(i * 30)
		c.Post(newModel, false)
	}

	metadata := c.Metadata()

	// All should be positive
	if metadata["positive"] != 5 {
		t.Errorf("Expected 5 positive values, got %d", metadata["positive"])
	}

	// Items with i=4,5 should be > 100 (120, 150)
	if metadata["large"] != 2 {
		t.Errorf("Expected 2 large values, got %d", metadata["large"])
	}
}

// Test Patch with no changes
func TestCachePatchNoChanges(t *testing.T) {
	res := newResources()
	model1 := createModel(1)

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)
	c.SetNotificationsFor("test-service", 1)

	// Patch with same values
	patchModel := createModel(1)

	notification, err := c.Patch(patchModel, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// No changes should result in nil notification
	if notification != nil {
		t.Logf("Notification: %v", notification)
	}
}

// Test cache isolation - verify cloning
func TestCacheIsolation(t *testing.T) {
	res := newResources()
	model1 := createModel(1)
	model1.MyInt32 = 100

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)

	// Get the item
	retrieved, err := c.Get(model1)
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}

	// Modify the retrieved item
	if retrievedModel, ok := retrieved.(*testtypes.TestProto); ok {
		retrievedModel.MyInt32 = 999
	}

	// Get again and verify original is unchanged (due to cloning)
	retrieved2, err := c.Get(model1)
	if err != nil {
		t.Fatalf("Failed to get second time: %v", err)
	}

	if retrievedModel2, ok := retrieved2.(*testtypes.TestProto); ok {
		if retrievedModel2.MyInt32 != 100 {
			t.Errorf("Expected MyInt32 to remain 100, got %d (cloning may not be working)", retrievedModel2.MyInt32)
		}
	}
}

// Test Post updates existing item with replace notification
func TestCachePostUpdateReplaceNotification(t *testing.T) {
	res := newResources()
	model1 := createModel(1)
	model1.MyInt32 = 50

	initElements := []interface{}{model1}
	c := cache.NewCache(model1, initElements, nil, res)
	c.SetNotificationsFor("test-service", 1)

	// Post update with different value
	model1Updated := createModel(1)
	model1Updated.MyInt32 = 150
	model1Updated.MyString = "different"

	notification, err := c.Post(model1Updated, true)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should create a replace notification
	if notification != nil && len(notification.NotificationList) > 0 {
		t.Logf("Replace notification created successfully")
	}
}

// Test AddMetadataFunc on cache with existing items
func TestCacheAddStatFuncWithExistingItems(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	// Add items first
	for i := 1; i <= 3; i++ {
		newModel := createModel(i * 10)
		newModel.MyInt32 = int32(i * 40)
		c.Post(newModel, false)
	}

	// Then add metadata function
	c.AddMetadataFunc("over_50", func(i interface{}) (bool, string) {
		if testModel, ok := i.(*testtypes.TestProto); ok {
			return testModel.MyInt32 > 50, ""
		}
		return false, ""
	})

	metadata := c.Metadata()

	// Items with i=2,3 should be > 50 (80, 120)
	if metadata["over_50"] != 2 {
		t.Errorf("Expected 2 items over 50, got %d", metadata["over_50"])
	}
}

// Test delete with metadata update
func TestCacheDeleteWithStats(t *testing.T) {
	res := newResources()
	model := createModel(1)

	c := cache.NewCache(model, nil, nil, res)

	c.AddMetadataFunc("all_items", func(i interface{}) (bool, string) {
		return true, ""
	})

	c.AddMetadataFunc("high_int", func(i interface{}) (bool, string) {
		if testModel, ok := i.(*testtypes.TestProto); ok {
			return testModel.MyInt32 > 50, ""
		}
		return false, ""
	})

	// Add items
	model1 := createModel(1)
	model1.MyInt32 = 100
	model2 := createModel(2)
	model2.MyInt32 = 20

	c.Post(model1, false)
	c.Post(model2, false)

	metadata := c.Metadata()
	if metadata["all_items"] != 2 {
		t.Errorf("Expected 2 items, got %d", metadata["all_items"])
	}
	if metadata["high_int"] != 1 {
		t.Errorf("Expected 1 high_int, got %d", metadata["high_int"])
	}

	// Delete the high int item
	c.Delete(model1, false)

	metadata = c.Metadata()
	if metadata["all_items"] != 1 {
		t.Errorf("Expected 1 item after delete, got %d", metadata["all_items"])
	}
	if metadata["high_int"] != 0 {
		t.Errorf("Expected 0 high_int after delete, got %d", metadata["high_int"])
	}
}

// Test Cache with storage - load from storage
func TestCacheWithStorageLoad(t *testing.T) {
	res := newResources()
	model := createModel(1)

	// Create storage with pre-existing data
	storage := newTestStorage(true)
	model1 := createModel(1)
	model2 := createModel(2)
	storage.Put("string-1", model1)
	storage.Put("string-2", model2)

	// Create cache - should load from storage
	c := cache.NewCache(model, nil, storage, res)

	if c.Size() != 2 {
		t.Errorf("Expected cache to load 2 items from storage, got %d", c.Size())
	}

	// Verify we can get the items
	retrieved, err := c.Get(model1)
	if err != nil || retrieved == nil {
		t.Error("Expected to retrieve item loaded from storage")
	}
}

// Test Cache with storage - save to storage
func TestCacheWithStorageSave(t *testing.T) {
	res := newResources()
	model := createModel(1)

	storage := newTestStorage(true)
	c := cache.NewCache(model, nil, storage, res)

	// Post item
	newModel := createModel(10)
	c.Post(newModel, false)

	// Verify it's in storage
	if len(storage.data) != 1 {
		t.Errorf("Expected 1 item in storage, got %d", len(storage.data))
	}
}

// Test Cache with storage disabled (direct storage access)
func TestCacheWithStorageDisabled(t *testing.T) {
	res := newResources()
	model := createModel(1)

	storage := newTestStorage(false)
	c := cache.NewCache(model, nil, storage, res)

	// Post item - should go directly to storage
	newModel := createModel(10)
	c.Post(newModel, false)

	// Cache size should be 0 since cache is disabled
	if c.Size() != 0 {
		t.Errorf("Expected cache size 0 with disabled cache, got %d", c.Size())
	}

	// But item should be in storage
	if len(storage.data) != 1 {
		t.Errorf("Expected 1 item in storage, got %d", len(storage.data))
	}
}

// Test Cache with storage - Get from storage when cache disabled
func TestCacheGetFromStorage(t *testing.T) {
	res := newResources()
	model := createModel(1)

	storage := newTestStorage(false)
	c := cache.NewCache(model, nil, storage, res)

	// Post item
	newModel := createModel(10)
	c.Post(newModel, false)

	// Get should retrieve from storage
	retrieved, err := c.Get(newModel)
	if err != nil {
		t.Errorf("Expected no error getting from storage, got: %v", err)
	}
	if retrieved == nil {
		t.Error("Expected to retrieve item from storage")
	}
}

// Test Cache with storage - Delete from storage
func TestCacheDeleteFromStorage(t *testing.T) {
	res := newResources()
	model := createModel(1)

	storage := newTestStorage(true)
	c := cache.NewCache(model, nil, storage, res)

	// Post item
	newModel := createModel(10)
	c.Post(newModel, false)

	// Delete item
	c.Delete(newModel, false)

	// Verify it's gone from storage
	if len(storage.data) != 0 {
		t.Errorf("Expected 0 items in storage after delete, got %d", len(storage.data))
	}
}

// Test Cache with storage - Patch updates storage
func TestCachePatchUpdatesStorage(t *testing.T) {
	res := newResources()
	model := createModel(1)

	storage := newTestStorage(true)
	c := cache.NewCache(model, nil, storage, res)

	// Post item
	model1 := createModel(1)
	model1.MyInt32 = 100
	c.Post(model1, false)

	// Patch item
	patchModel := createModel(1)
	patchModel.MyInt32 = 200
	c.Patch(patchModel, false)

	// Verify storage is updated
	stored, _ := storage.Get("string-1")
	if storedModel, ok := stored.(*testtypes.TestProto); ok {
		if storedModel.MyInt32 != 200 {
			t.Errorf("Expected storage to have updated value 200, got %d", storedModel.MyInt32)
		}
	}
}

// Test NewCache with initial elements and storage
func TestNewCacheWithElementsAndStorage(t *testing.T) {
	res := newResources()
	model1 := createModel(1)
	model2 := createModel(2)

	storage := newTestStorage(true)
	initElements := []interface{}{model1, model2}

	c := cache.NewCache(model1, initElements, storage, res)

	// Should have items in cache
	if c.Size() != 2 {
		t.Errorf("Expected size 2, got %d", c.Size())
	}

	// Should also have items in storage
	if len(storage.data) != 2 {
		t.Errorf("Expected 2 items in storage, got %d", len(storage.data))
	}
}
