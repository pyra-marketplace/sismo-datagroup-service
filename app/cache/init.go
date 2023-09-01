package cache

import (
	"fmt"
	"github.com/fanjindong/go-cache"
	"sismo-datagroup-service/app/model"
	"time"
)

var _memCache cache.ICache

var MemCache = new(MemoryCache)

type MemoryCache struct{}

func InitCache() {
	_memCache = cache.NewMemCache()
}

func (*MemoryCache) Cached(meta *model.DataGroupMate) {
	_memCache.Set(meta.GroupName, meta, cache.WithEx(24*time.Hour))
}

func (*MemoryCache) GetCachedMeta(groupName string) (interface{}, bool) {
	return _memCache.Get(groupName)
}

func (*MemoryCache) CacheRecord(groupName string, account string, record *model.DataGroupRecord) {
	key := fmt.Sprintf("%s_%s", groupName, account)
	_memCache.Set(key, record, cache.WithEx(24*time.Hour))
}

func (*MemoryCache) GetCachedRecordByGroupName(groupName string, account string) (interface{}, bool) {
	key := fmt.Sprintf("%s_%s", groupName, account)
	return _memCache.Get(key)
}

func (*MemoryCache) CacheGroupMembers(groupName string, groupMembers map[string]string, expiredAt time.Time) {
	key := fmt.Sprintf("%s::%s", groupName, "members")
	_memCache.Set(key, groupMembers, cache.WithExAt(expiredAt))
}

func (*MemoryCache) GetGroupMembers(groupName string) (interface{}, bool) {
	key := fmt.Sprintf("%s::%s", groupName, "members")
	return _memCache.Get(key)
}
