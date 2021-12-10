package core

type Storage interface {
	Ls(path string) (err error)
	Put(src, dst string) (err error)
	//Init() error
}

// Config for storage
type Config struct {
	Endpoint  string
	AccessID  string
	SecretKey string
	SSL       bool
	Region    string
	Path      string
	Bucket    string
	Addr      string
	Driver    string
}

/*
// Storage for s3 and disk
type Storage interface {
	// CreateBucket for create new folder
	CreateBucket(context.Context, string, string) error
	// BucketExists Checks if a bucket exists.
	BucketExists(ctx context.Context, bucketName string) (found bool, err error)
	// UploadFile for upload single file
	UploadFile(context.Context, string, string, []byte, io.Reader) error
	// UploadFileByReader for upload single file
	UploadFileByReader(context.Context, string, string, io.Reader, string, int64) error
	// DeleteFile for delete single file
	DeleteFile(context.Context, string, string) error
	// FilePath for store path + file name
	FilePath(string, string) string
	// GetFile for storage host + bucket + filename
	GetFileURL(string, string) string
	// DownloadFile downloads and saves the object as a file in the local filesystem.
	DownloadFile(context.Context, string, string, string) error
	// DownloadFileByProgress downloads and saves the object as a file in the local filesystem.
	DownloadFileByProgress(context.Context, string, string, string, *pb.ProgressBar) error
	// FileExist check object exist. bucket + filename
	FileExist(context.Context, string, string) bool
	// GetContent for storage bucket + filename
	GetContent(context.Context, string, string) ([]byte, error)
	// Copy Create or replace an object through server-side copying of an existing object.
	CopyFile(context.Context, string, string, string, string) error
	// Client get storage client
	Client() interface{}
	// SignedURL get signed URL
	SignedURL(context.Context, string, string, *SignedURLOptions) (string, error)
}*/
