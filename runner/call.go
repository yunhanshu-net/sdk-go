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

	//router, exist := r.getRouter(req.Request.Route, req.Request.Method)
	//if !exist {
	//	logrus.Errorf("call getRouter 404,%s %s", req.Request.Method, req.Request.Route)
	//	return nil, fmt.Errorf("call getRouter 404,%s %s", req.Request.Method, req.Request.Route)
	//}
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
