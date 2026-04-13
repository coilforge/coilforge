package part

import (
	"coilforge/internal/core"
	"encoding/json"
	"fmt"
)

type Record struct {
	Type core.PartTypeID `json:"type"`
	Data json.RawMessage `json:"data"`
}

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

func DecodeRecord(record Record) (Part, error) {
	info, ok := Registry[record.Type]
	if !ok {
		return nil, fmt.Errorf("unknown part type %q", record.Type)
	}
	return info.Decode(record.Data)
}
