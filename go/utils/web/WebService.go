package web

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/saichler/l8utils/go/utils/strings"
	"os"
	"reflect"
	"strconv"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8web"
	"google.golang.org/protobuf/encoding/protojson"
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
	webService.webService.Endpoints = make(map[int32]*l8web.L8WebServiceEndpoint)
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

func (this *WebService) Protos(bodyData string, action ifs.Action) (proto.Message, proto.Message, error) {
	endpoint, ok := this.webService.Endpoints[int32(action)]
	if !ok {
		msg := strings.New("endpoint not found for action ", this.ServiceName(), " area ", this.ServiceArea(), " ", strconv.Itoa(int(action)))
		return nil, nil, errors.New(msg.String())
	}

	bodyType := endpoint.PrimaryBody

	body := this.protos[bodyType]
	body = reflect.New(reflect.ValueOf(body).Elem().Type()).Interface().(proto.Message)
	err := protojson.Unmarshal([]byte(bodyData), body)
	var resp proto.Message

	if err != nil {
		fmt.Println("Error for ", bodyType, " is ", err.Error())
		bodyType = ""
		for k, v := range endpoint.Body2Response {
			if k != endpoint.PrimaryBody {
				body = this.protos[k]
				body = reflect.New(reflect.ValueOf(body).Elem().Type()).Interface().(proto.Message)
				if k != endpoint.PrimaryBody {
					err = protojson.Unmarshal([]byte(bodyData), body)
					if err == nil {
						resp = this.protos[v]
						resp = reflect.New(reflect.ValueOf(resp).Elem().Type()).Interface().(proto.Message)
						bodyType = k
						break
					}
					fmt.Println("Error for ", bodyType, " is ", err.Error())
				}
			}
		}
	}
	if bodyType == "" {
		msg := strings.New("cannot find any matching body type ", this.ServiceName(), " area ", this.ServiceArea(), " ", strconv.Itoa(int(action)))
		return nil, nil, errors.New(msg.String())
	}
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

	actionEndpoint, ok := this.webService.Endpoints[int32(action)]
	if !ok {
		actionEndpoint = &l8web.L8WebServiceEndpoint{}
		actionEndpoint.Body2Response = make(map[string]string)
		actionEndpoint.PrimaryBody = bodyType
		this.webService.Endpoints[int32(action)] = actionEndpoint
	}
	respType := this.typeOf(resp)
	actionEndpoint.Body2Response[bodyType] = respType

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

func (this *WebService) DeSerialize(ws *l8web.L8WebService, r ifs.IRegistry) error {
	for k, v := range ws.Endpoints {
		for bodyType, respType := range v.Body2Response {
			info, e := r.Info(bodyType)
			if e != nil {
				return e
			}
			body, e := info.NewInstance()
			if e != nil {
				return e
			}

			info, e = r.Info(respType)
			if e != nil {
				return e
			}
			resp, e := info.NewInstance()
			if e != nil {
				return e
			}
			this.AddEndpoint(body.(proto.Message), ifs.Action(k), resp.(proto.Message))
		}
	}
	return nil
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
