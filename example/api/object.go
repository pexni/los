package api

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pexni/los"
	"github.com/pexni/xhttp"
)

var ObjectApi objectApi

type objectApi struct{}

func (a *objectApi) PutObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")

	_, err := los.SClient.HeadBucket(&los.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}

	putObjectOutput, err := los.SClient.PutObject(&los.PutObjectInput{
		Bucket: bucket,
		Key:    key,
		Body:   r.Body,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	xhttp.JSON(w, http.StatusOK, putObjectOutput)
}

func (a *objectApi) ListObjects(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	_, err := los.SClient.HeadBucket(&los.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	listObjectsOutput, err := los.SClient.ListObjects(&los.ListObjectsInput{Bucket: bucket})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	xhttp.JSON(w, http.StatusOK, listObjectsOutput)
}

func (a *objectApi) GetObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	_, err := los.SClient.HeadBucket(&los.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	getObjectOutput, err := los.SClient.GetObject(&los.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	// if the body is nil, can not defer close
	defer getObjectOutput.Body.Close()
	all, err := io.ReadAll(getObjectOutput.Body)
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	http.ServeContent(w, r, key, getObjectOutput.LastModified, bytes.NewReader(all))
}

func (a *objectApi) DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	_, err := los.SClient.HeadBucket(&los.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	deleteObjectOutput, err := los.SClient.DeleteObject(&los.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	xhttp.JSON(w, http.StatusNoContent, deleteObjectOutput)
}

func (a *objectApi) HeadObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	_, err := los.SClient.HeadBucket(&los.HeadBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}

	headObjectOutput, err := los.SClient.HeadObject(&los.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		xhttp.BadRequest(w, 400, err.Error(), nil)
		return
	}
	w.Header().Add("Content-Length", strconv.FormatInt(headObjectOutput.ContentLength, 10))
	w.Header().Add("Date", headObjectOutput.LastModified.String())
	w.Header().Add("Content-Type", headObjectOutput.ContentType)
	w.WriteHeader(http.StatusOK)
}
