package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"spectrocloud.com/spectromate/mock"
)

func TestNewCache(t *testing.T) {
	connectionString := "localhost:6379"
	redisUser := ""
	redisPassword := ""
	tlsConfig := &tls.Config{}

	cache := NewCache(connectionString, redisUser, redisPassword, tlsConfig)

	assert.NotNil(t, cache, "Expected non-nil cache object")
}

func TestStoreHashMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mock.NewMockCache(ctrl)

	ctx := context.Background()
	primaryKey := "testPrimaryKey"
	item := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	// Test successful hash map storage
	cache.EXPECT().StoreHashMap(ctx, primaryKey, item).Return(nil)
	err := cache.StoreHashMap(ctx, primaryKey, item)
	assert.NoError(t, err, "Expected no error when storing valid hash map")

	// Test storing hash map with empty primary key
	cache.EXPECT().StoreHashMap(ctx, "", item).Return(fmt.Errorf("invalid primary key"))
	err = cache.StoreHashMap(ctx, "", item)
	assert.Error(t, err, "Expected error when storing hash map with empty primary key")

	// Test storing hash map with nil item
	cache.EXPECT().StoreHashMap(ctx, primaryKey, nil).Return(fmt.Errorf("invalid item"))
	err = cache.StoreHashMap(ctx, primaryKey, nil)
	assert.Error(t, err, "Expected error when storing hash map with nil item")

	// Test error when storing hash map due to Redis client error
	cache.EXPECT().StoreHashMap(ctx, primaryKey, item).Return(fmt.Errorf("Redis client error"))
	err = cache.StoreHashMap(ctx, primaryKey, item)
	assert.Error(t, err, "Expected error when Redis client encounters an error")

	// Test error when storing hash map due to context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	cache.EXPECT().StoreHashMap(cancelCtx, primaryKey, item).Return(fmt.Errorf("context canceled"))
	err = cache.StoreHashMap(cancelCtx, primaryKey, item)
	assert.Error(t, err, "Expected error when context is canceled")
}

func TestStoreHashMapDifferentValueTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mock.NewMockCache(ctrl)

	ctx := context.Background()
	primaryKey := "testPrimaryKey"
	item := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": 45.67,
		"key4": true,
	}

	// Test successful hash map storage with different value types
	cache.EXPECT().StoreHashMap(ctx, primaryKey, item).Return(nil)
	err := cache.StoreHashMap(ctx, primaryKey, item)
	assert.NoError(t, err, "Expected no error when storing hash map with different value types")
}

func TestExpireKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mock.NewMockCache(ctrl)

	ctx := context.Background()
	primaryKey := "testPrimaryKey"

	// Test successful expiration setting
	cache.EXPECT().ExpireKey(ctx, primaryKey, 5*time.Second).Return(nil)
	err := cache.ExpireKey(ctx, primaryKey, 5*time.Second)
	assert.NoError(t, err, "Expected no error when setting expiration on valid key")

	// Test setting expiration with empty primary key
	cache.EXPECT().ExpireKey(ctx, "", 5*time.Second).Return(fmt.Errorf("invalid key"))
	err = cache.ExpireKey(ctx, "", 5*time.Second)
	assert.Error(t, err, "Expected error when setting expiration with empty primary key")

	// Test error when setting expiration due to Redis client error
	cache.EXPECT().ExpireKey(ctx, primaryKey, 5*time.Second).Return(fmt.Errorf("Redis client error"))
	err = cache.ExpireKey(ctx, primaryKey, 5*time.Second)
	assert.Error(t, err, "Expected error when Redis client encounters an error")

	// Test error when setting expiration due to context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	cache.EXPECT().ExpireKey(cancelCtx, primaryKey, 5*time.Second).Return(fmt.Errorf("context canceled"))
	err = cache.ExpireKey(cancelCtx, primaryKey, 5*time.Second)
	assert.Error(t, err, "Expected error when context is canceled")
}

func TestGetHashMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mock.NewMockCache(ctrl)

	ctx := context.Background()
	primaryKey := "testPrimaryKey"
	item := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Test successful hash map retrieval
	cache.EXPECT().GetHashMap(ctx, primaryKey).Return(true, item, nil)
	found, resultMap, err := cache.GetHashMap(ctx, primaryKey)
	assert.True(t, found, "Expected the hash map to be found")
	assert.NoError(t, err, "Expected no error when retrieving valid hash map")
	require.NotNil(t, resultMap, "Expected non-nil result map")
	assert.Equal(t, "value1", resultMap["key1"], "Expected value1 for key1")
	assert.Equal(t, "value2", resultMap["key2"], "Expected value2 for key2")

	// Test retrieving non-existent hash map
	cache.EXPECT().GetHashMap(ctx, primaryKey).Return(false, nil, nil)
	found, resultMap, err = cache.GetHashMap(ctx, primaryKey)
	assert.False(t, found, "Expected the hash map not to be found")
	assert.NoError(t, err, "Expected no error when retrieving non-existent hash map")
	assert.Nil(t, resultMap, "Expected nil result map for non-existent hash map")

	// Test error when retrieving hash map due to Redis client error
	cache.EXPECT().GetHashMap(ctx, primaryKey).Return(false, nil, fmt.Errorf("Redis client error"))
	found, _, err = cache.GetHashMap(ctx, primaryKey)
	assert.False(t, found, "Expected the hash map not to be found")
	assert.Error(t, err, "Expected error when Redis client encounters an error")

	// Test error when retrieving hash map due to context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()
	cache.EXPECT().GetHashMap(cancelCtx, primaryKey).Return(false, nil, fmt.Errorf("context canceled"))
	found, _, err = cache.GetHashMap(cancelCtx, primaryKey)
	assert.False(t, found, "Expected the hash map not to be found")
	assert.Error(t, err, "Expected error when context is canceled")
}

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mock.NewMockCache(ctrl)

	// Test successful ping
	cache.EXPECT().Ping().Return(nil)
	err := cache.Ping()
	assert.NoError(t, err, "Expected no error when pinging Redis with a valid connection")

	// Test error when pinging Redis due to Redis client error
	cache.EXPECT().Ping().Return(fmt.Errorf("Redis client error"))
	err = cache.Ping()
	assert.Error(t, err, "Expected error when Redis client encounters an error")
}
