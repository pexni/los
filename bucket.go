package los

import (
	"os"
	"path"
)

func (c *Client) CreateBucket(params *CreateBucketInput) (*CreateBucketOutput, error) {
	if params == nil {
		params = &CreateBucketInput{}
	}
	bucketPath := path.Join(c.BucketsLocation, params.Bucket)
	err := os.Mkdir(bucketPath, 0755)
	if err != nil {
		return nil, err
	}
	return &CreateBucketOutput{}, nil
}

func (c *Client) ListBuckets(params *ListBucketsInput) (*ListBucketsOutput, error) {
	if params == nil {
		params = &ListBucketsInput{}
	}
	buckets := make([]Bucket, 0)
	dirEntries, err := os.ReadDir(c.BucketsLocation)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		buckets = append(buckets, Bucket{Name: entry.Name()})
	}
	return &ListBucketsOutput{Buckets: buckets}, nil
}

func (c *Client) DeleteBucket(params *DeleteBucketInput) (*DeleteBucketOutput, error) {
	if params == nil {
		params = &DeleteBucketInput{}
	}
	bucketPath := path.Join(c.BucketsLocation, params.Bucket)
	err := os.RemoveAll(bucketPath)
	if err != nil {
		return nil, err
	}
	return &DeleteBucketOutput{}, nil
}

func (c *Client) HeadBucket(params *HeadBucketInput) (*HeadBucketOutput, error) {
	if params == nil {
		params = &HeadBucketInput{}
	}
	bucketPath := path.Join(c.BucketsLocation, params.Bucket)
	_, err := os.Stat(bucketPath)
	if err != nil {
		return nil, err
	}
	return &HeadBucketOutput{}, nil
}
