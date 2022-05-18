package frame

import (
	"strconv"
	"time"

	"github.com/yomorun/y3"
)

// State frame state
type State byte

const (
	// Normal data frame is normal(defaults)
	Normal State = 0x0
	// Dispatched data frame is dispatched
	Dispatched State = 0x1
	// Processed data frame is processed(with sfn)
	Processed State = 0x2
)

func (s State) String() string{
	switch s {
	case Dispatched:
		return "Dispatched"
	case Processed:
		return "Processed"
	case Normal:
		return "Normal"
	default:
		return "Unknown"
	}
}

// MetaFrame is a Y3 encoded bytes, SeqID is a fixed value of TYPE_ID_TRANSACTION.
// used for describes metadata for a DataFrame.
type MetaFrame struct {
	tid      string
	sourceID string
	state    State // data process state
}

// NewMetaFrame creates a new MetaFrame instance.
func NewMetaFrame() *MetaFrame {
	return &MetaFrame{
		tid: strconv.FormatInt(time.Now().Unix(), 10),
	}
}

// SetTransactinID set the transaction ID.
func (m *MetaFrame) SetTransactionID(transactionID string) {
	m.tid = transactionID
}

// TransactionID returns transactionID
func (m *MetaFrame) TransactionID() string {
	return m.tid
}

// SetSourceID set the source ID.
func (m *MetaFrame) SetSourceID(sourceID string) {
	m.sourceID = sourceID
}

// SourceID returns source ID
func (m *MetaFrame) SourceID() string {
	return m.sourceID
}

// SetState set the frame state
func (m *MetaFrame) SetState(state State) {
	m.state = state
}

// State returns the frame state
func (m *MetaFrame) State() State {
	return m.state
}

// Encode implements Frame.Encode method.
func (m *MetaFrame) Encode() []byte {
	meta := y3.NewNodePacketEncoder(byte(TagOfMetaFrame))
	// transaction id
	transactionID := y3.NewPrimitivePacketEncoder(byte(TagOfTransactionID))
	transactionID.SetStringValue(m.tid)
	meta.AddPrimitivePacket(transactionID)
	// source id
	sourceID := y3.NewPrimitivePacketEncoder(byte(TagOfSourceID))
	sourceID.SetStringValue(m.sourceID)
	meta.AddPrimitivePacket(sourceID)
	// state
	state := y3.NewPrimitivePacketEncoder(byte(TagOfState))
	state.SetBytesValue([]byte{byte(m.state)})
	meta.AddPrimitivePacket(state)

	return meta.Encode()
}

// DecodeToMetaFrame decode a MetaFrame instance from given buffer.
func DecodeToMetaFrame(buf []byte) (*MetaFrame, error) {
	node := y3.NodePacket{}
	_, err := y3.DecodeToNodePacket(buf, &node)
	if err != nil {
		return nil, err
	}

	meta := &MetaFrame{}
	// for _, v := range node.PrimitivePackets {
	// 	val, _ := v.ToUTF8String()
	// 	meta.tid = val
	// 	break
	// }
	// transaction id
	if transactionIDBlock, ok := node.PrimitivePackets[byte(TagOfTransactionID)]; ok {
		tid, err := transactionIDBlock.ToUTF8String()
		if err != nil {
			return nil, err
		}
		meta.tid = tid
	}
	// source id
	if sourceIDBlock, ok := node.PrimitivePackets[byte(TagOfSourceID)]; ok {
		sourceID, err := sourceIDBlock.ToUTF8String()
		if err != nil {
			return nil, err
		}
		meta.sourceID = sourceID
	}
	// state
	if stateBlock, ok := node.PrimitivePackets[byte(TagOfState)]; ok {
		if state := stateBlock.ToBytes(); len(state) > 0 {
			meta.state = State(state[0])
		}
	}

	return meta, nil
}
