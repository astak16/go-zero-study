package logic

import (
	"context"
	"strconv"

	"go-zero-study/api/internal/svc"
	"go-zero-study/api/internal/types"
	"go-zero-study/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	userRpc := l.svcCtx.UserRpc

	userId, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, err
	}

	userInfo, err := userRpc.UserInfo(l.ctx, &user.UserInfoRequest{
		UserId: uint32(userId),
	})
	if err != nil {
		return nil, err
	}
	return &types.UserInfoResponse{
		Username: userInfo.Username,
		UserId:   strconv.Itoa(int(userInfo.UserId)),
	}, nil
}
