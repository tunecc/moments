package handler

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kingwrcy/moments/vo"
)

// newS3Client 根据系统配置创建 S3 客户端。
// 统一封装,替代此前在 file.go / memo.go 中重复的初始化逻辑,
// 并使用 v2 推荐的 BaseEndpoint 取代已废弃的 WithEndpointResolver。
func newS3Client(ctx context.Context, s3vo vo.S3VO) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(s3vo.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3vo.AccessKey,
				s3vo.SecretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		if s3vo.Endpoint != "" {
			o.BaseEndpoint = aws.String(s3vo.Endpoint)
		}
	}), nil
}
