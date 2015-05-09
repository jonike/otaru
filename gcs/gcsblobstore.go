package gcs

import (
	"io"

	"golang.org/x/net/context"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"

	"github.com/nyaxt/otaru"
)

type GCSBlobStore struct {
	projectName string
	bucketName  string
	flags       int
	clisrc      ClientSource
}

func NewGCSBlobStore(projectName string, bucketName string, credentialsFilePath string, tokenCacheFilePath string, flags int) (*GCSBlobStore, error) {
	clisrc, err := GetGCloudClientSource(credentialsFilePath, tokenCacheFilePath, false)
	if err != nil {
		return nil, err
	}

	return &GCSBlobStore{
		projectName: projectName,
		bucketName:  bucketName,
		flags:       flags,
		clisrc:      clisrc,
	}, nil
}

type Writer struct {
	gcsw *storage.Writer
}

func (bs *GCSBlobStore) newAuthedContext(basectx context.Context) context.Context {
	return cloud.NewContext(bs.projectName, bs.clisrc(context.TODO()))
}

func (bs *GCSBlobStore) OpenWriter(blobpath string, flags int) (io.WriteCloser, error) {
	if !otaru.IsWriteAllowed(bs.flags) || !otaru.IsWriteAllowed(flags) {
		return nil, otaru.EPERM
	}

	ctx := bs.newAuthedContext(context.TODO())

	gcsw := storage.NewWriter(ctx, bs.bucketName, blobpath)
	gcsw.ContentType = "application/octet-stream"
	return &Writer{gcsw}, nil
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.gcsw.Write(p)
}

func (w *Writer) Close() error {
	if err := w.gcsw.Close(); err != nil {
		return err
	}

	// obj := w.gcsw.Object()
	// do something???

	return nil
}

func (bs *GCSBlobStore) OpenReader(blobpath string, flags int) (io.ReadCloser, error) {
	ctx := bs.newAuthedContext(context.TODO())
	return storage.NewReader(ctx, bs.bucketName, blobpath)
}

func (bs *GCSBlobStore) Flags() int {
	return bs.flags
}

func (bs *GCSBlobStore) Delete(blobpath string) error {
	ctx := bs.newAuthedContext(context.TODO())
	if err := storage.DeleteObject(ctx, bs.bucketName, blobpath); err != nil {
		return err
	}
	return nil
}