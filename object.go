package los

import (
	"io"
	"os"
	"path"
	"strings"
)

func (c *Client) PutObject(params *PutObjectInput) (*PutObjectOutput, error) {
	if params == nil {
		params = &PutObjectInput{}
	}
	objectPath := path.Join(c.BucketsLocation, params.Bucket, params.Key)
	keyDir, _ := path.Split(params.Key)
	if keyDir != "" {
		err := os.MkdirAll(path.Join(c.BucketsLocation, params.Bucket, keyDir), 0755)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(objectPath)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, params.Body)
	if err != nil {
		return nil, err
	}
	return &PutObjectOutput{}, nil
}

func (c *Client) ListObjects(params *ListObjectsInput) (*ListObjectsOutput, error) {
	if params == nil {
		params = &ListObjectsInput{}
	}
	bucketPath := path.Join(c.BucketsLocation, params.Bucket)
	contents := make([]Object, 0)
	c.getContents(bucketPath, params.Bucket, &contents)
	return &ListObjectsOutput{
		Name:     params.Bucket,
		Contents: contents,
	}, nil
}

func (c *Client) getContents(bucketPath, bucket string, objects *[]Object) {
	entries, err := os.ReadDir(bucketPath)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			c.getContents(path.Join(bucketPath, entry.Name()), bucket, objects)
		} else {
			fileInfo, _ := entry.Info()
			key := strings.Replace(path.Join(bucketPath, fileInfo.Name()), path.Join(c.BucketsLocation, bucket), "", 1)
			*objects = append(*objects, Object{
				Key:          key,
				LastModified: fileInfo.ModTime(),
				Size:         fileInfo.Size(),
			})
		}
	}
}

func (c *Client) DeleteObject(params *DeleteObjectInput) (*DeleteObjectOutput, error) {
	if params == nil {
		params = &DeleteObjectInput{}
	}
	objectPath := path.Join(c.BucketsLocation, params.Bucket, params.Key)
	err := os.Remove(objectPath)
	if err != nil {
		return nil, err
	}
	return &DeleteObjectOutput{}, nil
}

func (c *Client) GetObject(params *GetObjectInput) (*GetObjectOutput, error) {
	if params == nil {
		params = &GetObjectInput{}
	}
	objectPath := path.Join(c.BucketsLocation, params.Bucket, params.Key)
	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &GetObjectOutput{
		Body:          file,
		ContentLength: fileInfo.Size(),
		LastModified:  fileInfo.ModTime(),
	}, nil
}

func (c *Client) HeadObject(params *HeadObjectInput) (*HeadObjectOutput, error) {
	if params == nil {
		params = &HeadObjectInput{}
	}
	objectPath := path.Join(c.BucketsLocation, params.Bucket, params.Key)
	fileInfo, err := os.Stat(objectPath)
	if err != nil {
		return nil, err
	}
	return &HeadObjectOutput{
		ContentLength: fileInfo.Size(),
		LastModified:  fileInfo.ModTime(),
	}, nil
}
