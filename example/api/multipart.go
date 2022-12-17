package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pexni/los"
	"github.com/pexni/xhttp"
)

var MultipartApi multipartApi

type multipartApi struct{}

func (a *multipartApi) CreateMultipartUpload(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")

	createMultipartUploadOutput, err := los.SClient.CreateMultipartUpload(&los.CreateMultipartUploadInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	data := make(map[string]interface{}, 0)
	data["bucket"] = createMultipartUploadOutput.Bucket
	data["key"] = createMultipartUploadOutput.Key
	data["uploadId"] = createMultipartUploadOutput.UploadId
	xhttp.JSON(w, http.StatusOK, data)
}

func (a *multipartApi) UploadPart(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	partNumberStr := r.URL.Query().Get("partNumber")
	uploadId := r.URL.Query().Get("uploadId")
	partNumber, err := strconv.ParseInt(partNumberStr, 10, 64)
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	uploadPartOutput, err := los.SClient.UploadPart(&los.UploadPartInput{
		Bucket:     bucket,
		Key:        key,
		PartNumber: int32(partNumber),
		UploadId:   uploadId,
		Body:       r.Body,
	})
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	xhttp.JSON(w, http.StatusOK, uploadPartOutput)
}

func (a *multipartApi) ListParts(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	uploadId := r.URL.Query().Get("uploadId")

	listPartsOutput, err := los.SClient.ListParts(&los.ListPartsInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: uploadId,
		MaxParts: 1000,
	})
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	xhttp.JSON(w, http.StatusOK, listPartsOutput)
}

func (a *multipartApi) AbortMultipartUpload(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	uploadId := r.URL.Query().Get("uploadId")

	abortMultipartUploadOutput, err := los.SClient.AbortMultipartUpload(&los.AbortMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: uploadId,
	})
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	xhttp.JSON(w, http.StatusNoContent, abortMultipartUploadOutput)
}

func (a *multipartApi) CompleteMultipartUpload(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	key := chi.URLParam(r, "*")
	uploadId := r.URL.Query().Get("uploadId")

	listPartsOutput, err := los.SClient.ListParts(&los.ListPartsInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: uploadId,
	})
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	completeParts := make([]los.CompletedPart, 0)
	for _, part := range listPartsOutput.Parts {
		completeParts = append(completeParts, los.CompletedPart{PartNumber: part.PartNumber})
	}

	completeMultipartUploadOutput, err := los.SClient.CompleteMultipartUpload(&los.CompleteMultipartUploadInput{
		Bucket:   bucket,
		Key:      key,
		UploadId: uploadId,
		MultipartUpload: los.CompletedMultipartUpload{
			Parts: completeParts,
		},
	})
	if err != nil {
		xhttp.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	xhttp.JSON(w, http.StatusOK, completeMultipartUploadOutput)
}
