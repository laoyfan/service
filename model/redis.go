package model

type RedisData struct {
	Id    string `json:"id" binding:"required" validate:"numeric"`
	Value string `json:"value" binding:"required"`
}

type GetRedisData struct {
	Id string `json:"id" form:"id" url:"id" binding:"required" validate:"numeric"`
}
