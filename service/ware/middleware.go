package ware

import (
	"context"
	"net/http"
	meta "octopus/metadata"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type HandlerUnit func(ctx context.Context, data *meta.MetaData) error

type AfterHandlerUnit func(ctx context.Context, message proto.Message, response http.ResponseWriter, callbackheader metadata.MD) error

type Middleware func(HandlerUnit) HandlerUnit
