package web

import (
	"github.com/saichler/l8types/go/ifs"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type WebService struct {
	serviceName string
	serviceArea uint16

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

func New(serviceName string, serviceArea uint16,
	postBody, postResp,
	putBody, putResp,
	patchBody, patchResp,
	deleteBody, deleteResp,
	getBody, getResp proto.Message) ifs.IWebService {
	webService := &WebService{}
	webService.serviceName = serviceName
	webService.serviceArea = serviceArea

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
