package otaru

import (
	"bytes"
	"io"
	"testing"
)

var (
	Key          = []byte("0123456789abcdef")
	HelloWorld   = []byte("Hello, world")
	HogeFugaPiyo = []byte("hogefugapiyo")
)

func testCipher() Cipher {
	c, err := NewCipher(Key)
	if err != nil {
		panic("Failed to init Cipher for testing")
	}
	return c
}

func genTestData(l int) []byte {
	ret := make([]byte, l)
	i := 0
	for i+3 < l {
		x := i / 4
		ret[i+0] = byte((x >> 24) & 0xff)
		ret[i+1] = byte((x >> 16) & 0xff)
		ret[i+2] = byte((x >> 8) & 0xff)
		ret[i+3] = byte((x >> 0) & 0xff)
		i += 4
	}
	if i < l {
		copy(ret[i:], []byte{0xab, 0xcd, 0xef})
	}
	return ret
}

func negateBits(p []byte) []byte {
	ret := make([]byte, len(p))
	for i, x := range p {
		ret[i] = 0xff ^ x
	}
	return ret
}

func Test_genTestData(t *testing.T) {
	td := genTestData(5)
	if !bytes.Equal(td, []byte{0, 0, 0, 0, 0xab}) {
		t.Errorf("unexp testdata: %v", td)
	}

	td = genTestData(11)
	if !bytes.Equal(td, []byte{0, 0, 0, 0, 0, 0, 0, 1, 0xab, 0xcd, 0xef}) {
		t.Errorf("unexp testdata: %v", td)
	}
}

func genFrameByChunkWriter(t *testing.T, p []byte) []byte {
	buf := new(bytes.Buffer)
	cw := NewChunkWriter(buf, testCipher())

	err := cw.WriteHeaderAndPrologue(
		len(p),
		&ChunkPrologue{OrigFilename: "testframe.dat", OrigOffset: 0},
	)
	if err != nil {
		t.Errorf("Failed to write chunk header: %v", err)
		return nil
	}

	if _, err := cw.Write(p); err != nil {
		t.Errorf("Failed to write chunk payload: %v", err)
		return nil
	}

	if err := cw.Close(); err != nil {
		t.Errorf("Failed to close ChunkWriter: %v", err)
		return nil
	}

	return buf.Bytes()
}

func TestChunkIO_Read_HelloWorld(t *testing.T) {
	b := genFrameByChunkWriter(t, HelloWorld)
	if b == nil {
		return
	}
	testbh := &TestBlobHandle{b}
	cio := NewChunkIO(testbh, testCipher())

	readtgt := make([]byte, len(HelloWorld))
	if err := cio.PRead(0, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, HelloWorld) {
		t.Errorf("Read content invalid")
		return
	}

	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}
}

func TestChunkIO_Read_1MB(t *testing.T) {
	td := genTestData(1024*1024 + 123)
	b := genFrameByChunkWriter(t, td)
	if b == nil {
		return
	}
	testbh := &TestBlobHandle{b}
	cio := NewChunkIO(testbh, testCipher())

	// Full read
	readtgt := make([]byte, len(td))
	if err := cio.PRead(0, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, td) {
		t.Errorf("Read content invalid")
		return
	}

	// Partial read
	readtgt = readtgt[:321]
	if err := cio.PRead(1012345, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, td[1012345:1012345+321]) {
		t.Errorf("Read content invalid")
		return
	}

	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}
}

func TestChunkIO_Write_UpdateHello(t *testing.T) {
	b := genFrameByChunkWriter(t, HelloWorld)
	if b == nil {
		return
	}
	testbh := &TestBlobHandle{b}
	cio := NewChunkIO(testbh, testCipher())

	upd := []byte("testin write")
	if err := cio.PWrite(0, upd); err != nil {
		t.Errorf("failed to PWrite to ChunkIO: %v", err)
		return
	}

	readtgt := make([]byte, len(upd))
	if err := cio.PRead(0, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, upd) {
		t.Errorf("Read content invalid")
		return
	}

	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}
}

func TestChunkIO_Write_Update1MB(t *testing.T) {
	td := genTestData(1024*1024 + 123)
	b := genFrameByChunkWriter(t, td)
	if b == nil {
		return
	}
	testbh := &TestBlobHandle{b}
	cio := NewChunkIO(testbh, testCipher())

	// Full update
	td2 := negateBits(td)
	if err := cio.PWrite(0, td2); err != nil {
		t.Errorf("failed to PWrite into ChunkIO: %v", err)
		return
	}
	readtgt := make([]byte, len(td))
	if err := cio.PRead(0, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, td2) {
		t.Errorf("Read content invalid")
		return
	}

	// Partial update
	if err := cio.PWrite(1012345, td[1012345:1012345+321]); err != nil {
		t.Errorf("failed to PWrite into ChunkIO: %v", err)
		return
	}
	td3 := make([]byte, len(td2))
	copy(td3, td2)
	copy(td3[1012345:1012345+321], td[1012345:1012345+321])
	if err := cio.PRead(0, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, td3) {
		t.Errorf("Read content invalid")
		return
	}
	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}
}

func Test_ChunkIOWrite_NewHello_ChunkReaderRead(t *testing.T) {
	testbh := &TestBlobHandle{}
	cio := NewChunkIO(testbh, testCipher())
	if err := cio.PWrite(0, HelloWorld); err != nil {
		t.Errorf("failed to PWrite to ChunkIO: %v", err)
		return
	}
	readtgt := make([]byte, len(HelloWorld))
	if err := cio.PRead(0, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, HelloWorld) {
		t.Errorf("Read content invalid")
		return
	}
	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}

	cr := NewChunkReader(bytes.NewBuffer(testbh.Buf), testCipher())
	if err := cr.ReadHeader(); err != nil {
		t.Errorf("failed to read header: %v", err)
		return
	}
	if err := cr.ReadPrologue(); err != nil {
		t.Errorf("failed to read prologue: %v", err)
		return
	}
	if cr.Length() != len(HelloWorld) {
		t.Errorf("failed to recover payload len")
	}
	readtgt2 := make([]byte, len(HelloWorld))
	if _, err := io.ReadFull(cr, readtgt2); err != nil {
		t.Errorf("failed to Read from ChunkReader: %v", err)
		return
	}
	if !bytes.Equal(readtgt2, HelloWorld) {
		t.Errorf("Read content invalid")
		return
	}
}

func checkZero(t *testing.T, p []byte, off int, length int) {
	i := 0
	for i < length {
		if p[off+i] != 0 {
			t.Errorf("Given slice non-0 at idx: %d", off+i)
		}
		i++
	}
}

func Test_ChunkIOWrite_ZeroFillPadding(t *testing.T) {
	testbh := &TestBlobHandle{}
	cio := NewChunkIO(testbh, testCipher())

	// [ zero ][ hello ]
	//    10      12
	if err := cio.PWrite(10, HelloWorld); err != nil {
		t.Errorf("failed to PWrite to ChunkIO: %v", err)
		return
	}
	readtgt := make([]byte, len(HelloWorld))
	if err := cio.PRead(10, readtgt); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	if !bytes.Equal(readtgt, HelloWorld) {
		t.Errorf("Read content invalid")
		return
	}
	readtgt2 := make([]byte, 10+len(HelloWorld))
	if err := cio.PRead(0, readtgt2); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	checkZero(t, readtgt2, 0, 10)
	if !bytes.Equal(readtgt2[10:10+12], HelloWorld) {
		t.Errorf("Read content invalid: hello1 %v != %v", readtgt2[10:10+12], HelloWorld)
		return
	}

	// [ zero ][ hello ][ zero ][ hello ]
	//    10      12      512k      12
	if err := cio.PWrite(10+12+512*1024, HelloWorld); err != nil {
		t.Errorf("failed to PWrite to ChunkIO: %v", err)
		return
	}
	readtgt3 := make([]byte, 10+12+512*1024+12)
	if err := cio.PRead(0, readtgt3); err != nil {
		t.Errorf("failed to PRead from ChunkIO: %v", err)
		return
	}
	checkZero(t, readtgt3, 0, 10)
	checkZero(t, readtgt3, 10+12, 512*1024)
	if !bytes.Equal(readtgt3[10:10+12], HelloWorld) {
		t.Errorf("Read content invalid: hello1")
		return
	}
	if !bytes.Equal(readtgt3[10+12+512*1024:10+12+512*1024+12], HelloWorld) {
		t.Errorf("Read content invalid: hello2")
		return
	}

	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}
}

func Test_ChunkIOWrite_OverflowUpdate(t *testing.T) {
	testbh := &TestBlobHandle{}
	cio := NewChunkIO(testbh, testCipher())
	if err := cio.PWrite(0, HelloWorld); err != nil {
		t.Errorf("failed to PWrite to ChunkIO: %v", err)
		return
	}
	if err := cio.PWrite(7, HogeFugaPiyo); err != nil {
		t.Errorf("failed to PWrite to ChunkIO: %v", err)
		return
	}
	if err := cio.Close(); err != nil {
		t.Errorf("failed to Close ChunkIO: %v", err)
		return
	}

	cr := NewChunkReader(bytes.NewBuffer(testbh.Buf), testCipher())
	if err := cr.ReadHeader(); err != nil {
		t.Errorf("failed to read header: %v", err)
		return
	}
	if err := cr.ReadPrologue(); err != nil {
		t.Errorf("failed to read prologue: %v", err)
		return
	}
	exp := []byte("Hello, hogefugapiyo")
	if cr.Length() != len(exp) {
		t.Errorf("failed to recover payload len")
	}
	readtgt := make([]byte, len(exp))
	if _, err := io.ReadFull(cr, readtgt); err != nil {
		t.Errorf("failed to Read from ChunkReader: %v", err)
		return
	}
	if !bytes.Equal(readtgt, exp) {
		t.Errorf("Read content invalid")
		return
	}
}
