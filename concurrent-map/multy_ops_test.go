package cmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConcurrentMap_MSet(t *testing.T) {
	const elementsNum = 10000

	m := NewDefault()

	data := map[uint16]interface{}{}

	for i := 0; i < elementsNum; i++ {
		data[uint16(i)] = i << 8
	}

	m.MSet(data)

	for key, value := range data {
		shard := m.GetShard(key)

		shard.RLock()
		actual, ok := shard.items[key]
		require.True(t, ok, fmt.Sprintf("key: %d, shard index: %d, value: %d", key, m.getShardIndex(key), value))
		require.Equal(t, value, actual)
		shard.RUnlock()
	}
}

func Test_ConcurrentMap_MSetIfAbsent_EmptyMap(t *testing.T) {
	const elementsNum = 10000
	m := NewDefault()

	data := map[uint16]interface{}{}

	for i := 0; i < elementsNum; i++ {
		data[uint16(i)] = i << 8
	}

	m.MSetIfAbsent(data)

	for key, value := range data {
		shard := m.GetShard(key)

		shard.RLock()
		actual, ok := shard.items[key]
		require.True(t, ok, fmt.Sprintf("key: %d, shard index: %d, value: %d", key, m.getShardIndex(key), value))
		require.Equal(t, value, actual)
		shard.RUnlock()
	}
}

func Test_ConcurrentMap_MSetIfAbsent_NotEmptyMap(t *testing.T) {
	const elementsNum = 10000
	m := NewDefault()

	data := map[uint16]interface{}{}

	for i := 0; i < elementsNum; i++ {
		data[uint16(i)] = i << 8
	}

	key := uint16(elementsNum - 1)
	value := "not absent"

	m.shards[m.getShardIndex(key)].items[key] = value

	m.MSetIfAbsent(data)

	delete(m.shards[m.getShardIndex(key)].items, key)

	for key, value := range data {
		shard := m.GetShard(key)

		shard.RLock()
		_, ok := shard.items[key]
		require.False(t, ok, fmt.Sprintf("key: %d, shard index: %d, value: %d", key, m.getShardIndex(key), value))
		shard.RUnlock()
	}
}

func Test_ConcurrentMap_MRemove(t *testing.T) {
	const elementsNum = 10000
	m := NewDefault()

	keys := make([]uint16, 0)

	for i := 0; i < elementsNum; i++ {
		value := i << 8
		keys = append(keys, uint16(i))
		m.shards[m.getShardIndex(uint16(i))].items[uint16(i)] = value
	}

	m.MRemove(keys)

	for _, key := range keys {
		shard := m.GetShard(key)

		shard.RLock()
		_, ok := shard.items[key]
		require.False(t, ok, fmt.Sprintf("key: %d, shard index: %d", key, m.getShardIndex(key)))
		shard.RUnlock()
	}
}
