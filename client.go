package los

import "path"

var SClient *Client

type Client struct {
	BaseLocation      string
	BucketsLocation   string
	MultipartLocation string
}

func NewClient(baseLocation string) *Client {
	return &Client{
		BaseLocation:      path.Join(baseLocation),
		BucketsLocation:   path.Join(baseLocation, "buckets"),
		MultipartLocation: path.Join(baseLocation, "multipart"),
	}
}
