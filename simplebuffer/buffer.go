package simplebuffer

import "sync"

type Buffer interface {
	Add(key any, data any) bool
	AddWithCallback(key any, data any, callback func(data map[any]any) error) (error, bool)
	GetData() map[any]any
	Reset()
}

type bufferCtx struct {
	mutex sync.Mutex
	data  map[any]any
	size  int
}

// New - buffer creation method
func New(size int) Buffer {
	return &bufferCtx{size: size, data: map[any]any{}}
}

// GetData - gets value copy of buffer data, mutex support
func (bc *bufferCtx) GetData() map[any]any {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	return bc.getData()
}

// Reset - creates new empty buffer, mutex support
func (bc *bufferCtx) Reset() {
	bc.mutex.Lock()
	bc.resetBuffer()
	bc.mutex.Unlock()
}

// Add - append key value pair to buffer, mutex support, false will be returned if buffer is full
func (bc *bufferCtx) Add(key any, data any) bool {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	if len(bc.data) == bc.size {
		return false
	}
	bc.data[key] = data
	return true
}

// AddWithCallback - append key value pair to buffer, calls callback function when buffer is full
func (bc *bufferCtx) AddWithCallback(key any, data any, callback func(data map[any]any) error) (error, bool) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	if len(bc.data) == bc.size {
		d := bc.getData()
		if err := callback(d); err != nil {
			return err, false
		}
		bc.resetBuffer()
	}
	bc.data[key] = data
	return nil, true
}

func (bc *bufferCtx) getData() map[any]any {
	vh := make(map[any]any, len(bc.data))
	for k, v := range bc.data {
		vh[k] = v
	}
	return vh
}

func (bc *bufferCtx) resetBuffer() {
	bc.data = make(map[any]any)
}
