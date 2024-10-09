package model

type RedisData struct {
	Id    string `json:"id" binding:"required,numeric"`
	Value string `json:"value" binding:"required"`
	Code  int    `json:"code" binding:"-" validate:"numeric"`
}

type GetRedisData struct {
	Id string `json:"id" form:"id" url:"id" binding:"required,numeric"`
}
