package userlogic

import (
	"context"

	"go-zero-study/model"
	"go-zero-study/rpc/internal/svc"
	"go-zero-study/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserResponse, error) {

	db := l.svcCtx.DB
	u := &model.User{}
	err := db.Take(u, "id = ?", in.UserId).Error
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{
		Username: u.Username,
		UserId:   uint32(u.ID),
	}, nil
}
