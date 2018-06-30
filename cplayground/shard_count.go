package cplayground

const defaultShardCount = 32

const (
	sizeSmall  = 30 * 30
	sizeMiddle = 100 * 100
	sizeLarge  = 200 * 200
)

const (
	shardCountSmallMap   = 2
	shardCountMiddleMap  = 16
	shardCountLargeMap   = 32
	shardCountBiggestMap = 64
)

func calcShardCount(size uint16) int {
	if size < sizeSmall {
		return shardCountSmallMap
	}

	if size < sizeMiddle {
		return shardCountMiddleMap
	}

	if size < sizeLarge {
		return shardCountLargeMap
	}

	return shardCountBiggestMap
}
