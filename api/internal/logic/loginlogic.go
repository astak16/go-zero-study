package logic

import (
	"context"

	"go-zero-study/api/internal/svc"
	"go-zero-study/api/internal/types"
	"go-zero-study/common/jwt"
	"go-zero-study/rpc/client/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserRequest) (resp string, err error) {
	// todo: add your logic here and delete this line

	userResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.UserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return "", err
	}
	auth := l.svcCtx.Config.Auth
	token, err := jwt.GenToken(jwt.JwtPayload{UserName: userResp.Username, UserID: uint(userResp.UserId)}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		return "失败", err
	}
	return token, nil
}
