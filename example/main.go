package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pexni/los"
	"github.com/pexni/los/example/api"
)

func main() {
	initClient()
	r := initHttpHandler()
	mountHandlers(r)
	startHttpServer(r)
}

func initClient() {
	los.SClient = los.NewClient("public")
}

func initHttpHandler() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func mountHandlers(r *chi.Mux) {
	r.Post("/buckets/{bucket}", api.BucketApi.CreateBucket)
	r.Get("/buckets", api.BucketApi.ListBuckets)
	r.Delete("/buckets/{bucket}", api.BucketApi.DeleteBucket)
	r.Head("/buckets/{bucket}", api.BucketApi.HeadBucket)

	r.Put("/{bucket}/*", api.ObjectApi.PutObject)
	r.Get("/{bucket}/objects", api.ObjectApi.ListObjects)
	r.Get("/{bucket}/*", api.ObjectApi.GetObject)
	r.Delete("/{bucket}/*", api.ObjectApi.DeleteObject)
	r.Head("/{bucket}/*", api.ObjectApi.HeadObject)

	r.Post("/multipart/{bucket}/*", api.MultipartApi.CreateMultipartUpload)
	r.Put("/multipart/{bucket}/*", api.MultipartApi.UploadPart)
	r.Get("/multipart/{bucket}/*", api.MultipartApi.ListParts)
	r.Delete("/multipart/{bucket}/*", api.MultipartApi.AbortMultipartUpload)
	r.Post("/multipart/complete/{bucket}/*", api.MultipartApi.CompleteMultipartUpload)
}

func startHttpServer(r *chi.Mux) {
	s := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	log.Fatal(s.ListenAndServe())
}
