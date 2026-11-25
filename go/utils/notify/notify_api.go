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
	// Assign the id here (the service keeps a preset id) so the persisted
	// record can be read back after a transactional commit.
	record := &ntf.NotifyRecord{
		NotifyId: ifs.NewUuid(),
		Channel:  channel, Endpoint: endpoint, Subject: subject, Message: message, Attributes: attributes,
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
		// A commit only means the record was persisted — the record itself
		// carries the dispatch outcome (a failed SMTP send is still committed,
		// with Status=FAILED). Read it back instead of assuming success.
		return this.persistedResult(record.NotifyId)
	}
	return &ntf.DeliveryResult{Status: ntf.DeliveryStatus_DELIVERY_STATUS_FAILED, ErrorMessage: "unexpected response type"}
}

// persistedResult reads the committed NotifyRecord and returns its delivery
// outcome. If the record cannot be read, the outcome is reported as unknown
// (PENDING with an explanation) rather than claimed as SENT.
func (this *Notify) persistedResult(notifyId string) *ntf.DeliveryResult {
	resp := this.vnic.Request("", NotifyServiceName, NotifyServiceArea, ifs.GET,
		&ntf.NotifyRecord{NotifyId: notifyId}, 30)
	if resp.Error() == nil {
		if stored, ok := resp.Element().(*ntf.NotifyRecord); ok && stored.NotifyId == notifyId {
			return &ntf.DeliveryResult{
				Status: stored.Status, HttpStatus: stored.HttpStatus,
				ErrorMessage: stored.ErrorMessage, Attempt: stored.Attempt, SentAt: stored.SentAt,
			}
		}
	}
	reason := "record not returned"
	if resp.Error() != nil {
		reason = resp.Error().Error()
	}
	return &ntf.DeliveryResult{
		Status:       ntf.DeliveryStatus_DELIVERY_STATUS_PENDING,
		ErrorMessage: ifs.App("notification committed but delivery outcome unknown (", notifyId, "): ", reason),
		SentAt:       time.Now().Unix(),
	}
}
