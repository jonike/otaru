package otaru

type MockBlobStoreOperation struct {
	Type      rune
	Offset    int64
	Length    int64
	FirstByte byte
}

type MockBlobHandle struct {
	Log        []MockBlobStoreOperation
	PayloadLen int64
}

func NewMockBlobHandle() *MockBlobHandle {
	return &MockBlobHandle{
		Log:        []MockBlobStoreOperation{},
		PayloadLen: 0,
	}
}

func (bh *MockBlobHandle) PRead(offset int64, p []byte) error {
	return nil
}

func (bh *MockBlobHandle) PWrite(offset int64, p []byte) error {
	right := offset + int64(len(p))
	if right > bh.PayloadLen {
		bh.PayloadLen = right
	}
	return nil
}

func (bh *MockBlobHandle) Size() int64 {
	return bh.PayloadLen
}

func (bh *MockBlobHandle) Close() error {
	return nil
}

type MockBlobStore struct {
	Paths map[string]*MockBlobHandle
}

func NewMockBlobStore() *MockBlobStore {
	return &MockBlobStore{make(map[string]*MockBlobHandle)}
}

func (bs *MockBlobStore) Open(blobpath string) (BlobHandle, error) {
	bh := bs.Paths[blobpath]
	if bh == nil {
		bh = NewMockBlobHandle()
		bs.Paths[blobpath] = bh
	}
	return bh, nil
}
