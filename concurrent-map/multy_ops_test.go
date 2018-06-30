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
		shard := m.getShard(key)

		shard.mux.RLock()
		actual, ok := shard.items[key]
		require.True(t, ok, fmt.Sprintf("key: %d, shard index: %d, value: %d", key, m.getShardIndex(key), value))
		require.Equal(t, value, actual)
		shard.mux.RUnlock()
	}
}

func Test_ConcurrentMap_MSetIfAllAbsent_EmptyMap(t *testing.T) {
	const elementsNum = 10000
	m := NewDefault()

	data := map[uint16]interface{}{}

	for i := 0; i < elementsNum; i++ {
		data[uint16(i)] = i << 8
	}

	m.MSetIfAllAbsent(data)

	for key, value := range data {
		shard := m.getShard(key)

		shard.mux.RLock()
		actual, ok := shard.items[key]
		require.True(t, ok, fmt.Sprintf("key: %d, shard index: %d, value: %d", key, m.getShardIndex(key), value))
		require.Equal(t, value, actual)
		shard.mux.RUnlock()
	}
}

func Test_ConcurrentMap_MSetIfAllAbsent_NotEmptyMap(t *testing.T) {
	const elementsNum = 10000
	m := NewDefault()

	data := map[uint16]interface{}{}

	for i := 0; i < elementsNum; i++ {
		data[uint16(i)] = i << 8
	}

	key := uint16(elementsNum - 1)
	value := "not absent"

	m.shards[m.getShardIndex(key)].items[key] = value

	m.MSetIfAllAbsent(data)

	delete(m.shards[m.getShardIndex(key)].items, key)

	for key, value := range data {
		shard := m.getShard(key)

		shard.mux.RLock()
		_, ok := shard.items[key]
		require.False(t, ok, fmt.Sprintf("key: %d, shard index: %d, value: %d", key, m.getShardIndex(key), value))
		shard.mux.RUnlock()
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
		shard := m.getShard(key)

		shard.mux.RLock()
		_, ok := shard.items[key]
		require.False(t, ok, fmt.Sprintf("key: %d, shard index: %d", key, m.getShardIndex(key)))
		shard.mux.RUnlock()
	}
}

func Test_ConcurrentMap_HasAny_EmptyMap(t *testing.T) {
	m := NewDefault()

	require.False(t, m.HasAny([]uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}))
}

func Test_ConcurrentMap_HasAny_NotEmptyMap(t *testing.T) {
	const index = 0

	m := NewDefault()

	m.shards[m.getShardIndex(index)].items[index] = "ok"

	require.True(t, m.HasAny([]uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}))
	require.False(t, m.HasAny([]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9}))
}

func Test_ConcurrentMap_HasAll_EmptyMap(t *testing.T) {
	m := NewDefault()

	require.False(t, m.HasAll([]uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}))
}

func Test_ConcurrentMap_HasAll_NotEmptyMap(t *testing.T) {
	const index = 0

	m := NewDefault()

	m.shards[m.getShardIndex(index)].items[index] = "ok"

	require.False(t, m.HasAll([]uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}))
	require.False(t, m.HasAll([]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	require.True(t, m.HasAll([]uint16{0}))
}
