package frame

import (
	"strconv"
	"time"

	"github.com/yomorun/y3"
)

// MetaFrame is a Y3 encoded bytes, SeqID is a fixed value of TYPE_ID_TRANSACTION.
// used for describes metadata for a DataFrame.
type MetaFrame struct {
	tid      string
	metadata []byte
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

// SetMetadata set the metadata
func (m *MetaFrame) SetMetadata(md []byte) {
	m.metadata = md
}

// TransactionID returns transactionID
func (m *MetaFrame) TransactionID() string {
	return m.tid
}

// Metadata returns metadata
func (m *MetaFrame) Metadata() []byte {
	return m.metadata
}

// Encode implements Frame.Encode method.
func (m *MetaFrame) Encode() []byte {
	meta := y3.NewNodePacketEncoder(byte(TagOfMetaFrame))

	transactionID := y3.NewPrimitivePacketEncoder(byte(TagOfTransactionID))
	transactionID.SetStringValue(m.tid)

	metadata := y3.NewPrimitivePacketEncoder(byte(TagOfMetadata))
	transactionID.SetBytesValue(m.metadata)

	meta.AddPrimitivePacket(transactionID)
	meta.AddPrimitivePacket(metadata)
	return meta.Encode()
}

// DecodeToMetaFrame decode a MetaFrame instance from given buffer.
func DecodeToMetaFrame(buf []byte) (*MetaFrame, error) {
	nodeBlock := y3.NodePacket{}
	_, err := y3.DecodeToNodePacket(buf, &nodeBlock)
	if err != nil {
		return nil, err
	}

	meta := &MetaFrame{}
	// for _, v := range nodeBlock.PrimitivePackets {
	// 	val, _ := v.ToUTF8String()
	// 	meta.tid = val
	// 	break
	// }

	// tid
	if tidPacket, ok := nodeBlock.PrimitivePackets[byte(TagOfTransactionID)]; ok {
		val, err := tidPacket.ToUTF8String()
		if err != nil {
			return nil, err
		}
		meta.tid = val
	}
	// metadata
	if metadataPacket, ok := nodeBlock.PrimitivePackets[byte(TagOfMetadata)]; ok {
		meta.metadata = metadataPacket.ToBytes()
	}

	return meta, nil
}
