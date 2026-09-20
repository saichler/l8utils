package notify

import (
	"time"

	"github.com/saichler/l8types/go/ifs"
	ntf "github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8types/go/types/l8services"
)

const (
	NotifyServiceName = "Notify"
	NotifyServiceArea = byte(78)
)

// Notify implements ifs.INotify, routing notification requests to the Notify service via VNic.
type Notify struct {
	vnic ifs.IVNic
}

func (this *Notify) SetVNic(vnic ifs.IVNic) {
	this.vnic = vnic
}

func (this *Notify) Send(channel ntf.NotifyChannel, endpoint, subject, message string,
	attributes map[string]string) *ntf.DeliveryResult {
	if this.vnic == nil {
		return &ntf.DeliveryResult{Status: ntf.DeliveryStatus_DELIVERY_STATUS_FAILED, ErrorMessage: "no VNic"}
	}
	record := &ntf.NotifyRecord{
		Channel: channel, Endpoint: endpoint, Subject: subject, Message: message, Attributes: attributes,
	}
	// Unlike Events' fire-and-forget Unicast, Send needs the dispatch OUTCOME back —
	// use a synchronous request, same pattern l8common.GetEntity/PostEntity use.
	resp := this.vnic.Request("", NotifyServiceName, NotifyServiceArea, ifs.POST, record, 30)
	if resp.Error() != nil {
		return &ntf.DeliveryResult{Status: ntf.DeliveryStatus_DELIVERY_STATUS_FAILED, ErrorMessage: resp.Error().Error()}
	}
	if posted, ok := resp.Element().(*ntf.NotifyRecord); ok {
		return &ntf.DeliveryResult{
			Status: posted.Status, HttpStatus: posted.HttpStatus,
			ErrorMessage: posted.ErrorMessage, Attempt: posted.Attempt, SentAt: posted.SentAt,
		}
	}
	// The Notify service is transactional, so a POST answers with the
	// transaction rather than the persisted record.
	if tr, ok := resp.Element().(*l8services.L8Transaction); ok {
		if ifs.TransactionState(tr.State) != ifs.Committed {
			return &ntf.DeliveryResult{
				Status:       ntf.DeliveryStatus_DELIVERY_STATUS_FAILED,
				ErrorMessage: ifs.App("notification was not committed (", ifs.TransactionState(tr.State).String(), "): ", tr.ErrMsg),
				SentAt:       time.Now().Unix(),
			}
		}
		// Committed says the record was written; the dispatch outcome it
		// recorded is only in that record, which this response does not
		// carry. Report accepted rather than claiming a delivery.
		return &ntf.DeliveryResult{Status: ntf.DeliveryStatus_DELIVERY_STATUS_PENDING, SentAt: time.Now().Unix()}
	}
	return &ntf.DeliveryResult{Status: ntf.DeliveryStatus_DELIVERY_STATUS_FAILED, ErrorMessage: "unexpected response type"}
}
