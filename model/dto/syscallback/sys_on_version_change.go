package syscallback

import "github.com/yunhanshu-net/sdk-go/model/dto/api"

type SysOnVersionChangeReq struct {
	OldVersion string `json:"old_version"`
	NewVersion string `json:"new_version"`
}
type SysOnVersionChangeResp struct {
	CurrentVersion string      `json:"current_version"` //此次更新的版本
	AddApi         []*api.Info `json:"add_api"`         //此次新增的api
	DelApi         []*api.Info `json:"del_api"`         //此次删除的api
	UpdateApi      []*api.Info `json:"update_api"`      //此次变更的api
}
