// Code generated by goctl. DO NOT EDIT.
package types

type UserInfoRequest struct {
	Id string `path:"id"`
}

type UserInfoResponse struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
