package request

type GetDeviceInfoReq struct {
	ProcessId string `json:"process_id" form:"processId" url:"processId" binding:"required,number"`
}

type GetTaskEnvReq struct {
	DeviceId  string `json:"deviceId" validate:"required,string"`
	ProcessId int    `json:"processId" validate:"required,numeric"`
}
