package common

import (
	"encoding/json"
)

type Stamp struct {
	ReplicaID ReplicaID
	Timestamp Timestamp
}

func (s *Stamp) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{s.ReplicaID, s.Timestamp})
}

func (s *Stamp) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[0], &s.ReplicaID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw[1], &s.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stamp) IsOlderThan(other *Stamp) bool {
	if s == nil {
		return true
	} else if other == nil {
		return false
	}
	if s.Timestamp == other.Timestamp {
		return s.ReplicaID < other.ReplicaID
	}

	return s.Timestamp < other.Timestamp
}
