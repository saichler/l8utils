package web

import (
	"encoding/base64"
	"fmt"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"google.golang.org/protobuf/proto"
	"os"
	"reflect"
	"strconv"
)

type WebService struct {
	serviceName string
	serviceArea byte

	postBody string
	postResp string

	putBody string
	putResp string

	patchBody string
	patchResp string

	deleteBody string
	deleteResp string

	getBody string
	getResp string

	plugin string
	pbs    map[string]proto.Message
}

func New(serviceName string, serviceArea byte,
	postBody, postResp,
	putBody, putResp,
	patchBody, patchResp,
	deleteBody, deleteResp,
	getBody, getResp proto.Message) ifs.IWebService {
	webService := &WebService{}
	webService.serviceName = serviceName
	webService.serviceArea = serviceArea
	webService.pbs = make(map[string]proto.Message)

	webService.postBody = webService.typeOf(postBody)
	webService.postResp = webService.typeOf(postResp)

	webService.putBody = webService.typeOf(putBody)
	webService.putResp = webService.typeOf(putResp)

	webService.patchBody = webService.typeOf(patchBody)
	webService.patchResp = webService.typeOf(patchResp)

	webService.deleteBody = webService.typeOf(deleteBody)
	webService.deleteResp = webService.typeOf(deleteResp)

	webService.getBody = webService.typeOf(getBody)
	webService.getResp = webService.typeOf(getResp)

	fmt.Println("Get Body:", webService.getBody, " Get Resp:", webService.getResp)

	filename := serviceName + "-" + strconv.Itoa(int(serviceArea)) + "-registry.so"
	_, err := os.Stat(filename)
	if err == nil {
		fmt.Print("Loaded registry plugin ", filename)
		data, err := os.ReadFile(filename)
		if err == nil {
			webService.plugin = base64.StdEncoding.EncodeToString(data)
		}
	}

	return webService
}

func (this *WebService) typeOf(pb proto.Message) string {
	if pb == nil {
		return ""
	}
	v := reflect.ValueOf(pb).Elem()
	this.pbs[v.Type().Name()] = pb
	return v.Type().Name()
}

func (this *WebService) Serialize() *types.WebService {
	r := &types.WebService{}
	r.ServiceName = this.serviceName
	r.ServiceArea = int32(this.serviceArea)

	r.PostRespType = this.postResp
	r.PostBodyType = this.postBody

	r.PutRespType = this.putResp
	r.PutBodyType = this.putBody

	r.PatchRespType = this.patchResp
	r.PatchBodyType = this.patchBody

	r.DeleteRespType = this.deleteResp
	r.DeleteBodyType = this.deleteBody

	r.GetRespType = this.getResp
	r.GetBodyType = this.getBody
	if this.plugin != "" {
		r.Plugin = &types.Plugin{Data: this.plugin}
	}
	return r
}

func (this *WebService) DeSerialize(ws *types.WebService) {
	this.serviceArea = byte(ws.ServiceArea)
	this.serviceName = ws.ServiceName

	this.postBody = ws.PostBodyType
	this.postResp = ws.PostRespType

	this.putBody = ws.PutBodyType
	this.putResp = ws.PutRespType

	this.patchBody = ws.PatchBodyType
	this.patchResp = ws.PatchRespType

	this.deleteBody = ws.DeleteBodyType
	this.deleteResp = ws.DeleteRespType

	this.getBody = ws.GetBodyType
	this.getResp = ws.GetRespType
	if ws.Plugin != nil {
		this.plugin = ws.Plugin.Data
	}
}
