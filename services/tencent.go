package services

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"hotel/internal/util"
	"net/http"
	"net/url"
)

func StartTencentCos(cfg *util.Config) *cos.Client {

	u, _ := url.Parse(cfg.Cos.Website)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.Tencent.SecretId,
			SecretKey: cfg.Tencent.SecretKey,
		},
	})

	return client
}
