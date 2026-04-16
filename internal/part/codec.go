package part

// File overview:
// codec serializes and deserializes heterogeneous part payloads using registry decoders.
// Subsystem: part serialization.
// It connects app/world save-load flows with catalog-specific decode implementations.
// Flow position: persistence adapter between files and in-memory part instances.

import (
	"coilforge/internal/core"
	"encoding/json"
	"fmt"
)

type Record struct {
	Type core.PartTypeID `json:"type"` // part type ID.
	Data json.RawMessage `json:"data"` // data value.
}

// EncodeRecord handles encode record.
func EncodeRecord(p Part) (Record, error) {
	data, err := p.MarshalJSON()
	if err != nil {
		return Record{}, err
	}

	return Record{
		Type: p.Base().TypeID,
		Data: data,
	}, nil
}

// DecodeRecord handles decode record.
func DecodeRecord(record Record) (Part, error) {
	info, ok := Registry[record.Type]
	if !ok {
		return nil, fmt.Errorf("unknown part type %q", record.Type)
	}
	return info.Decode(record.Data)
}
