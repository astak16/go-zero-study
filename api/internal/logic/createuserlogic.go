package logic

import (
	"context"
	"strconv"

	"go-zero-study/api/internal/svc"
	"go-zero-study/api/internal/types"
	"go-zero-study/rpc/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.UserRequest) (resp *types.UserInfoResponse, err error) {
	userRpc := l.svcCtx.UserRpc

	userRes, err := userRpc.UserCreate(l.ctx, &user.UserRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &types.UserInfoResponse{
		Username: userRes.Username,
		UserId:   strconv.FormatUint(uint64(userRes.UserId), 10),
	}, nil
}
