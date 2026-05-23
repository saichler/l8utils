package events

import (
	"fmt"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8events"
	"time"
)

const (
	EventsServiceName = "Events"
	EventsServiceArea = byte(76)
)

// Events implements ifs.IEvents, routing events to the Events service via VNic unicast.
type Events struct {
	vnic ifs.IVNic
}

func (this *Events) SetVnic(vnic ifs.IVNic) {
	this.vnic = vnic
}

func (this *Events) PostEvent(category l8events.EventCategory, eventType string,
	severity l8events.Severity, sourceId, sourceName, sourceType, message string,
	attributes map[string]string) {

	event := &l8events.EventRecord{
		Category:   category,
		EventType:  eventType,
		Severity:   severity,
		SourceId:   sourceId,
		SourceName: sourceName,
		SourceType: sourceType,
		Message:    message,
		OccurredAt: time.Now().Unix(),
		Attributes: attributes,
	}
	this.post(event)
}

func (this *Events) PostAuditEvent(evt *l8events.AuditEvent) {
	this.post(evt)
}

func (this *Events) PostSystemEvent(evt *l8events.SystemEvent) {
	this.post(evt)
}

func (this *Events) PostMonitoringEvent(evt *l8events.MonitoringEvent) {
	this.post(evt)
}

func (this *Events) PostSecurityEvent(evt *l8events.SecurityEvent) {
	this.post(evt)
}

func (this *Events) PostIntegrationEvent(evt *l8events.IntegrationEvent) {
	this.post(evt)
}

func (this *Events) PostNetworkEvent(evt *l8events.NetworkEvent) {
	this.post(evt)
}

func (this *Events) PostKubernetesEvent(evt *l8events.KubernetesEvent) {
	this.post(evt)
}

func (this *Events) PostPerformanceEvent(evt *l8events.PerformanceEvent) {
	this.post(evt)
}

func (this *Events) PostSyslogEvent(evt *l8events.SyslogEvent) {
	this.post(evt)
}

func (this *Events) PostTrapEvent(evt *l8events.TrapEvent) {
	this.post(evt)
}

func (this *Events) PostComputeEvent(evt *l8events.ComputeEvent) {
	this.post(evt)
}

func (this *Events) PostStorageEvent(evt *l8events.StorageEvent) {
	this.post(evt)
}

func (this *Events) PostPowerEvent(evt *l8events.PowerEvent) {
	this.post(evt)
}

func (this *Events) PostGpuEvent(evt *l8events.GpuEvent) {
	this.post(evt)
}

func (this *Events) PostTopologyEvent(evt *l8events.TopologyEvent) {
	this.post(evt)
}

func (this *Events) PostAutomationEvent(evt *l8events.AutomationEvent) {
	this.post(evt)
}

func (this *Events) post(payload interface{}) {
	if this.vnic == nil {
		fmt.Println("Events services does not have a VNic for:", payload)
		return
	}
	err := this.vnic.Unicast("", EventsServiceName, EventsServiceArea, ifs.POST, payload)
	if err != nil {
		this.vnic.Resources().Logger().Warning("PostEvent: " + err.Error())
	}
}
