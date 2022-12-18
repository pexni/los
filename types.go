package los

import (
	"io"
	"time"
)

type BucketCannedACL string

// Enum values for BucketCannedACL
const (
	BucketCannedACLPrivate           BucketCannedACL = "private"
	BucketCannedACLPublicRead        BucketCannedACL = "public-read"
	BucketCannedACLPublicReadWrite   BucketCannedACL = "public-read-write"
	BucketCannedACLAuthenticatedRead BucketCannedACL = "authenticated-read"
)

type ObjectCannedACL string

// Enum values for ObjectCannedACL
const (
	ObjectCannedACLPrivate                ObjectCannedACL = "private"
	ObjectCannedACLPublicRead             ObjectCannedACL = "public-read"
	ObjectCannedACLPublicReadWrite        ObjectCannedACL = "public-read-write"
	ObjectCannedACLAuthenticatedRead      ObjectCannedACL = "authenticated-read"
	ObjectCannedACLAwsExecRead            ObjectCannedACL = "aws-exec-read"
	ObjectCannedACLBucketOwnerRead        ObjectCannedACL = "bucket-owner-read"
	ObjectCannedACLBucketOwnerFullControl ObjectCannedACL = "bucket-owner-full-control"
)

type Bucket struct {
	Name string
}

type CreateBucketInput struct {
	Bucket string
	ACL    BucketCannedACL
}

type CreateBucketOutput struct {
}

type ListBucketsInput struct {
}

type ListBucketsOutput struct {
	Buckets []Bucket
}

type DeleteBucketInput struct {
	Bucket string
}

type DeleteBucketOutput struct {
}

type HeadBucketInput struct {
	Bucket              string
	ExpectedBucketOwner string
}

type HeadBucketOutput struct {
}

type Object struct {
	Key          string
	LastModified time.Time
	Size         int64
}

type PutObjectInput struct {
	Bucket string
	Key    string
	ACL    ObjectCannedACL
	Body   io.Reader
}

type PutObjectOutput struct {
}

type ListObjectsInput struct {
	Bucket string
}

type ListObjectsOutput struct {
	Name     string
	Contents []Object
}

type DeleteObjectInput struct {
	Bucket string
	Key    string
}

type DeleteObjectOutput struct {
}

type GetObjectInput struct {
	Bucket string
	Key    string
	Range  string
}

type GetObjectOutput struct {
	AcceptRanges  string
	Body          io.ReadCloser
	ContentLength int64
	ContentRange  string
	ContentType   string
	LastModified  time.Time
	PartsCount    int32
}

type HeadObjectInput struct {
	Bucket string
	Key    string
}

type HeadObjectOutput struct {
	AcceptRanges  string
	ContentLength int64
	ContentType   string
	LastModified  time.Time
	PartsCount    int32
}

type CreateMultipartUploadInput struct {
	Bucket string
	Key    string
	ACL    ObjectCannedACL
}

type CreateMultipartUploadOutput struct {
	Bucket   string
	Key      string
	UploadId string
}

type UploadPartInput struct {
	Bucket     string
	Key        string
	PartNumber int32
	UploadId   string
	Body       io.Reader
}

type UploadPartOutput struct {
}

type ListPartsInput struct {
	Bucket   string
	Key      string
	UploadId string
	MaxParts int32
}

type Part struct {
	LastModified time.Time
	PartNumber   int32
	Size         int64
}

type ListPartsOutput struct {
	Bucket   string
	Key      string
	MaxParts int32
	Parts    []Part
	UploadId string
}

type CompletedPart struct {
	PartNumber int32
}

type CompletedMultipartUpload struct {
	Parts []CompletedPart
}

type CompleteMultipartUploadInput struct {
	Bucket          string
	Key             string
	UploadId        string
	MultipartUpload CompletedMultipartUpload
}

type CompleteMultipartUploadOutput struct {
	Bucket string
	Key    string
}

type AbortMultipartUploadInput struct {
	Bucket   string
	Key      string
	UploadId string
}

type AbortMultipartUploadOutput struct {
}
