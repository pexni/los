package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pexni/los"
	"github.com/pexni/xhttp"
)

var BucketApi bucketApi

type bucketApi struct {
}

func (a *bucketApi) CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	createBucketOutput, err := los.SClient.CreateBucket(&los.CreateBucketInput{Bucket: bucket})
	if err != nil {
		xhttp.BadRequest(w, 200, err.Error(), nil)
		return
	}
	xhttp.JSON(w, http.StatusOK, createBucketOutput)
}

func (a *bucketApi) ListBuckets(w http.ResponseWriter, r *http.Request) {
	listBucketsOutput, err := los.SClient.ListBuckets(&los.ListBucketsInput{})
	if err != nil {
		xhttp.BadRequest(w, 200, err.Error(), nil)
		return
	}
	xhttp.JSON(w, http.StatusOK, listBucketsOutput)
}

func (a *bucketApi) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	deleteBucketOutput, err := los.SClient.DeleteBucket(&los.DeleteBucketInput{Bucket: bucket})
	if err != nil {
		xhttp.BadRequest(w, 200, err.Error(), nil)
		return
	}
	xhttp.JSON(w, http.StatusOK, deleteBucketOutput)
}

func (a *bucketApi) HeadBucket(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	_, err := los.SClient.HeadBucket(&los.HeadBucketInput{Bucket: bucket})
	if err != nil {
		xhttp.BadRequest(w, 200, err.Error(), nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}
