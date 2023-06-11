package telegram

import "net/http"

type Client struct {
	host     string
	basePath string
	client   http.Client
}
