package types

type HeightQueueItem struct {
	Height    int64
	Overwrite bool
}

// HeightQueue is a simple type alias for a (buffered) channel of block heights.
type HeightQueue chan HeightQueueItem

func NewQueue(size int) HeightQueue {
	return make(chan HeightQueueItem, size)
}
