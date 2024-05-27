package userlogic

import (
	"context"

	"go-zero-study/model"
	"go-zero-study/rpc/internal/svc"
	"go-zero-study/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCreateLogic) UserCreate(in *user.UserRequest) (*user.UserResponse, error) {
	// todo: add your logic here and delete this line
	db := l.svcCtx.DB
	u := &model.User{
		Username: in.Username,
		Password: in.Password,
	}

	if err := db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &user.UserResponse{
		Username: u.Username,
		UserId:   u.ID,
	}, nil
}
