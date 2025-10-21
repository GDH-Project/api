package grpc

import (
	"errors"

	"google.golang.org/grpc/status"
)

func errorFromGrpcError(err error) error {
	s, ok := status.FromError(err)
	if !ok {
		return errors.New("gRPC 에러 메시지 파일 오류")
	}
	return errors.New(s.Message())
}
