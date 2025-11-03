package web

func (this *WebService) ServiceName() string {
	return this.serviceName
}
func (this *WebService) ServiceArea() byte {
	return this.serviceArea
}
func (this *WebService) PostBody() string {
	return this.postBody
}
func (this *WebService) PostResp() string {
	return this.postResp
}
func (this *WebService) PutBody() string {
	return this.putBody
}
func (this *WebService) PutResp() string {
	return this.putResp
}
func (this *WebService) PatchBody() string {
	return this.patchBody
}
func (this *WebService) PatchResp() string {
	return this.patchResp
}
func (this *WebService) DeleteBody() string {
	return this.deleteBody
}
func (this *WebService) DeleteResp() string {
	return this.deleteResp
}
func (this *WebService) GetBody() string {
	return this.getBody
}
func (this *WebService) GetResp() string {
	return this.getResp
}
func (this *WebService) Plugin() string {
	return this.plugin
}
func (this *WebService) Vnet() uint32 {
	return this.vnet
}
func (this *WebService) SetVnet(vnet uint32) {
	this.vnet = vnet
}
