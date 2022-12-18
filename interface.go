package los

type ObjectStorageInterface interface {
	CreateBucket(params *CreateBucketInput) (*CreateBucketOutput, error)
	ListBuckets(params *ListBucketsInput) (*ListBucketsOutput, error)
	DeleteBucket(params *DeleteBucketInput) (*DeleteBucketOutput, error)
	HeadBucket(params *HeadBucketInput) (*HeadBucketOutput, error)

	PutObject(params *PutObjectInput) (*PutObjectOutput, error)
	ListObjects(params *ListObjectsInput) (*ListObjectsOutput, error)
	DeleteObject(params *DeleteObjectInput) (*DeleteObjectOutput, error)
	GetObject(params *GetObjectInput) (*GetObjectOutput, error)
	HeadObject(params *HeadObjectInput) (*HeadObjectOutput, error)

	CreateMultipartUpload(params *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error)
	UploadPart(params *UploadPartInput) (*UploadPartOutput, error)
	ListParts(params *ListPartsInput) (*ListPartsOutput, error)
	AbortMultipartUpload(params *AbortMultipartUploadInput) (*AbortMultipartUploadOutput, error)
	CompleteMultipartUpload(params *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error)
}
