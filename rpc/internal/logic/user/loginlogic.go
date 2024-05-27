package userlogic

import (
	"context"

	"go-zero-study/model"
	"go-zero-study/rpc/internal/svc"
	"go-zero-study/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.UserRequest) (*user.UserResponse, error) {
	// todo: add your logic here and delete this line
	db := l.svcCtx.DB
	u := &model.User{}
	err := db.Take(u, "username = ? and password = ?", in.Username, in.Password).Error
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{
		Username: u.Username,
		UserId:   uint32(u.ID),
	}, nil
}
