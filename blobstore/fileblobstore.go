package blobstore

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/dustin/go-humanize"

	fl "github.com/nyaxt/otaru/flags"
	"github.com/nyaxt/otaru/logger"
	"github.com/nyaxt/otaru/util"
)

var mylog = logger.Registry().Category("filebs")

type FileBlobHandle struct {
	Fp *os.File
}

func (h FileBlobHandle) PRead(p []byte, offset int64) error {
	for len(p) > 0 {
		n, err := h.Fp.ReadAt(p, offset)
		if err != nil {
			return err
		}
		if n == 0 {
			logger.Warningf(mylog, "*os.File ReadAt returned len 0!: %v", h.Fp)
			return io.EOF
		}
		p = p[n:]
	}
	return nil
}

func (h FileBlobHandle) PWrite(p []byte, offset int64) error {
	if _, err := h.Fp.WriteAt(p, offset); err != nil {
		return err
	}
	return nil
}

func (h FileBlobHandle) Size() int64 {
	fi, err := h.Fp.Stat()
	if err != nil {
		logger.Panicf(mylog, "Stat failed: %v", err)
	}

	return fi.Size()
}

func (h FileBlobHandle) Truncate(size int64) error {
	if err := h.Fp.Truncate(size); err != nil {
		return err
	}
	return nil
}

func (h FileBlobHandle) Close() error {
	return h.Fp.Close()
}

type FileBlobStore struct {
	base  string
	flags int
	fmask int
}

func NewFileBlobStore(base string, flags int) (*FileBlobStore, error) {
	base = path.Clean(base)

	fi, err := os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("Fstat base \"%s\" failed: %v", base, err)
	}
	if !fi.Mode().IsDir() {
		return nil, fmt.Errorf("Specified base \"%s\" is not a directory")
	}

	fbs := &FileBlobStore{
		base: base,
	}
	fbs.SetFlags(flags)
	return fbs, nil
}

func (f *FileBlobStore) SetFlags(flags int) {
	fmask := fl.O_RDONLY
	if fl.IsWriteAllowed(flags) {
		fmask = fl.O_RDONLY | fl.O_WRONLY | fl.O_RDWR | fl.O_CREATE | fl.O_EXCL
	}

	f.fmask = fmask
	f.flags = flags
}

func (f *FileBlobStore) Open(blobpath string, flags int) (BlobHandle, error) {
	realpath := path.Join(f.base, blobpath)

	fp, err := os.OpenFile(realpath, flags&f.fmask, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, util.ENOENT
		}
		return nil, err
	}
	return &FileBlobHandle{fp}, nil
}

func (f *FileBlobStore) Flags() int {
	return f.flags
}

func (f *FileBlobStore) OpenWriter(blobpath string) (io.WriteCloser, error) {
	if !fl.IsWriteAllowed(f.flags) {
		return nil, util.EACCES
	}

	realpath := path.Join(f.base, blobpath)
	return os.Create(realpath)
}

func (f *FileBlobStore) OpenReader(blobpath string) (io.ReadCloser, error) {
	if !fl.IsReadAllowed(f.flags) {
		return nil, util.EACCES
	}

	realpath := path.Join(f.base, blobpath)
	rc, err := os.Open(realpath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, util.ENOENT
		}
		return nil, err
	}
	return rc, nil
}

var _ = BlobLister(&FileBlobStore{})

func (f *FileBlobStore) ListBlobs() ([]string, error) {
	start := time.Now()

	d, err := os.Open(f.base)
	if err != nil {
		return nil, fmt.Errorf("Open dir failed: %v", err)
	}
	defer d.Close()
	fis, err := d.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("Readdir failed: %v", err)
	}

	blobs := make([]string, 0, len(fis))
	for _, fi := range fis {
		if !fi.Mode().IsRegular() {
			continue
		}
		blobs = append(blobs, fi.Name())
	}

	logger.Infof(mylog, "FileBlobStore.ListBlobs() found %d blobs, took %s.", len(blobs), time.Since(start))
	return blobs, nil
}

var _ = BlobSizer(&FileBlobStore{})

func (f *FileBlobStore) BlobSize(blobpath string) (int64, error) {
	realpath := path.Join(f.base, blobpath)

	fi, err := os.Stat(realpath)
	if err != nil {
		if os.IsNotExist(err) {
			return -1, util.ENOENT
		}
		return -1, err
	}

	return fi.Size(), nil
}

var _ = BlobRemover(&FileBlobStore{})

func (f *FileBlobStore) RemoveBlob(blobpath string) error {
	if !fl.IsWriteAllowed(f.flags) {
		return util.EACCES
	}
	err := os.Remove(path.Join(f.base, blobpath))
	if err != nil {
		if os.IsNotExist(err) {
			return util.ENOENT
		}
		return err
	}
	return nil
}

var _ = TotalSizer(&FileBlobStore{})

func (f *FileBlobStore) TotalSize() (int64, error) {
	start := time.Now()

	d, err := os.Open(f.base)
	if err != nil {
		return 0, fmt.Errorf("Open dir failed: %v", err)
	}
	defer d.Close()
	fis, err := d.Readdir(-1)
	if err != nil {
		return 0, fmt.Errorf("Readdir failed: %v", err)
	}

	totalSize := int64(0)
	for _, fi := range fis {
		if !fi.Mode().IsRegular() {
			continue
		}

		totalSize += fi.Size()
	}

	logger.Debugf(mylog, "FileBlobStore.TotalSize() was %s. took %s.", humanize.Bytes(uint64(totalSize)), time.Since(start))
	return totalSize, nil
}

func (*FileBlobStore) ImplName() string { return "FileBlobStore" }

func (f *FileBlobStore) GetBase() string { return f.base }
