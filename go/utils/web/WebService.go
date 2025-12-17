package web

import (
	"encoding/base64"
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8web"
	"google.golang.org/protobuf/proto"
)

type WebService struct {
	webService *l8web.L8WebService
	plugin     string
	protos     map[string]proto.Message
}

func New(serviceName string, serviceArea byte, vnet uint32) ifs.IWebService {
	webService := &WebService{}
	webService.webService = &l8web.L8WebService{}
	webService.webService.ServiceName = serviceName
	webService.webService.ServiceArea = int32(serviceArea)
	webService.webService.Endpoints = make(map[string]*l8web.L8WebServiceEndpoint)
	webService.protos = make(map[string]proto.Message)
	webService.webService.Vnet = vnet

	filename := serviceName + "-" + strconv.Itoa(int(serviceArea)) + "-registry.so"
	_, err := os.Stat(filename)
	if err == nil {
		data, err := os.ReadFile(filename)
		if err == nil {
			webService.plugin = base64.StdEncoding.EncodeToString(data)
		}
	}

	return webService
}

func (this *WebService) Protos(bodyType string, action ifs.Action) (proto.Message, proto.Message, error) {
	endpoint, ok := this.webService.Endpoints[bodyType]
	if !ok {
		return nil, nil, errors.New("endpoint not found for " + bodyType)
	}
	respType, ok := endpoint.Actions[int32(action)]
	if !ok {
		return nil, nil, errors.New("action not found for " + bodyType)
	}
	body := this.protos[bodyType]
	resp := this.protos[respType]
	return body, resp, nil
}

func (this *WebService) AddEndpoint(body proto.Message, action ifs.Action, resp proto.Message) {
	if body == nil {
		body = &l8web.L8Empty{}
	}
	if resp == nil {
		resp = &l8web.L8Empty{}
	}
	bodyType := this.typeOf(body)
	bodyEndpoint, ok := this.webService.Endpoints[bodyType]
	if !ok {
		bodyEndpoint = &l8web.L8WebServiceEndpoint{}
		bodyEndpoint.Actions = make(map[int32]string)
		this.webService.Endpoints[bodyType] = bodyEndpoint
	}
	respType := this.typeOf(resp)
	bodyEndpoint.Actions[int32(action)] = respType

	this.protos[bodyType] = body
	this.protos[respType] = resp
}

func (this *WebService) typeOf(pb proto.Message) string {
	if pb == nil {
		return ""
	}
	v := reflect.ValueOf(pb).Elem()
	return v.Type().Name()
}

func (this *WebService) Serialize() *l8web.L8WebService {
	return this.webService
}

func (this *WebService) DeSerialize(ws *l8web.L8WebService) {
	this.webService = ws
}

func (this *WebService) ServiceName() string {
	return this.webService.ServiceName
}

func (this *WebService) ServiceArea() byte {
	return byte(this.webService.ServiceArea)
}

func (this *WebService) Vnet() uint32 {
	return this.webService.Vnet
}

func (this *WebService) Plugin() string {
	return this.plugin
}
