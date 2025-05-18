package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
)

func (r *Runner) call(msg *nats.Msg) ([]byte, error) {

	data := msg.Data
	var req request.RunnerRequest
	err1 := json.Unmarshal(data, &req)
	if err1 != nil {
		logrus.Errorf("call  json.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call  json.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
	}

	runResponse, err1 := r.runRequest(context.Background(), req.Request)
	if err1 != nil {
		logrus.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
	}
	marshal, err1 := json.Marshal(runResponse)
	if err1 != nil {
		logrus.Errorf("call json.Marshal err,req:%+v err:%s", req, err1.Error())
		return nil, fmt.Errorf("call json.Marshal err,req:%+v err:%s", req, err1.Error())
	}

	return marshal, nil
}
