package logic

import (
	"context"

	"dex/app/api/internal/svc"
	"dex/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type APILogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAPILogic(ctx context.Context, svcCtx *svc.ServiceContext) *APILogic {
	return &APILogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *APILogic) API(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	logx.Infof("APILogic|req=%+v", req)
	return
}
