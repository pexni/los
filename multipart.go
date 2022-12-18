package los

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

const (
	UploadId = "uploadId"
	Complete = "complete"
)

func (c *Client) CreateMultipartUpload(params *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error) {
	bucketPath := path.Join(c.MultipartLocation, params.Bucket)
	_, err := os.Stat(bucketPath)
	if err != nil {
		if os.IsNotExist(err) {
			mkdirErr := os.Mkdir(bucketPath, 0755)
			if mkdirErr != nil {
				return nil, err
			}
		}
		return nil, err
	}
	objectDirPath := path.Join(bucketPath, params.Key)
	err = os.MkdirAll(objectDirPath, 0755)
	if err != nil {
		return nil, err
	}
	uploadIdFilePath := path.Join(objectDirPath, UploadId)
	uploadIdFile, err := os.Create(uploadIdFilePath)
	if err != nil {
		return nil, err
	}
	defer uploadIdFile.Close()
	hash := sha256.New()
	hash.Write([]byte(objectDirPath + time.Now().String()))
	uploadId := hex.EncodeToString(hash.Sum(nil))
	_, err = uploadIdFile.WriteString(uploadId)
	if err != nil {
		return nil, err
	}
	return &CreateMultipartUploadOutput{
		Bucket:   params.Bucket,
		Key:      params.Key,
		UploadId: uploadId,
	}, nil
}

func (c *Client) UploadPart(params *UploadPartInput) (*UploadPartOutput, error) {
	uploadIdPath := path.Join(c.MultipartLocation, params.Bucket, params.Key, UploadId)
	uploadId, err := os.ReadFile(uploadIdPath)
	if err != nil {
		return nil, err
	}
	if params.UploadId != string(uploadId) {
		return nil, errors.New("upload id invalid")
	}

	partPath := path.Join(c.MultipartLocation, params.Bucket, params.Key, fmt.Sprintf("%03d", params.PartNumber))
	file, err := os.Create(partPath)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, params.Body)
	if err != nil {
		return nil, err
	}
	return &UploadPartOutput{}, nil
}

func (c *Client) ListParts(params *ListPartsInput) (*ListPartsOutput, error) {
	uploadIdPath := path.Join(c.MultipartLocation, params.Bucket, params.Key, UploadId)
	uploadId, err := os.ReadFile(uploadIdPath)
	if err != nil {
		return nil, err
	}
	if params.UploadId != string(uploadId) {
		return nil, errors.New("upload id invalid")
	}

	objectDirPath := path.Join(c.MultipartLocation, params.Bucket, params.Key)
	entries, err := os.ReadDir(objectDirPath)
	if err != nil {
		return nil, err
	}
	parts := make([]Part, 0)
	for _, entry := range entries {
		if entry.Name() == UploadId {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			return nil, err
		}
		partNumber, parseIntErr := strconv.ParseInt(info.Name(), 10, 64)
		if parseIntErr != nil {
			return nil, err
		}
		parts = append(parts, Part{
			LastModified: info.ModTime(),
			PartNumber:   int32(partNumber),
			Size:         info.Size(),
		})
	}
	return &ListPartsOutput{
		Bucket:   params.Bucket,
		Key:      params.Key,
		MaxParts: 1000,
		Parts:    parts,
		UploadId: params.UploadId,
	}, nil
}

func (c *Client) AbortMultipartUpload(params *AbortMultipartUploadInput) (*AbortMultipartUploadOutput, error) {
	err := os.RemoveAll(path.Join(c.MultipartLocation, params.Bucket, params.Key))
	if err != nil {
		return nil, err
	}
	return &AbortMultipartUploadOutput{}, nil
}

func (c *Client) CompleteMultipartUpload(params *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error) {
	uploadIdPath := path.Join(c.MultipartLocation, params.Bucket, params.Key, UploadId)
	uploadId, err := os.ReadFile(uploadIdPath)
	if err != nil {
		return nil, err
	}
	if params.UploadId != string(uploadId) {
		return nil, errors.New("upload id invalid")
	}

	objectDirPath := path.Join(c.MultipartLocation, params.Bucket, params.Key)
	entries, err := os.ReadDir(objectDirPath)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(entries, func(i, j int) bool {
		nameI, parseIntErr := strconv.ParseInt(entries[i].Name(), 10, 64)
		if parseIntErr != nil {
			return false
		}
		nameJ, parseIntErr := strconv.ParseInt(entries[j].Name(), 10, 64)
		if parseIntErr != nil {
			return false
		}
		return nameI < nameJ
	})
	objectFilePath := path.Join(c.BucketsLocation, params.Bucket, params.Key)
	objectFile, err := os.Create(objectFilePath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		// skip not 001
		_, parseIntErr := strconv.ParseInt(entry.Name(), 10, 64)
		if parseIntErr != nil {
			continue
		}
		file, openErr := os.Open(path.Join(c.MultipartLocation, params.Bucket, params.Key, entry.Name()))
		if openErr != nil {
			_ = os.Remove(objectFilePath)
			return nil, openErr
		}
		_, copyErr := io.Copy(objectFile, file)
		if copyErr != nil {
			_ = os.Remove(objectFilePath)
			return nil, copyErr
		}
		file.Close()
	}

	completeObjectPath := path.Join(c.MultipartLocation, params.Bucket, params.Key, Complete)
	completeFile, err := os.Create(completeObjectPath)
	if err != nil {
		return nil, err
	}
	defer completeFile.Close()

	if err != nil {
		_ = os.Remove(completeObjectPath)
		_ = os.Remove(objectFilePath)
		return nil, err
	}

	return &CompleteMultipartUploadOutput{
		Bucket: params.Bucket,
		Key:    params.Key,
	}, nil
}
