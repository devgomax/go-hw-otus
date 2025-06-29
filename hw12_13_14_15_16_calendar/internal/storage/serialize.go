package storage

import (
	"encoding/json"

	"github.com/pkg/errors" //nolint:depguard
)

// Serialize сериализует модель БД в мапу.
func Serialize(e *Event) (map[string]any, error) {
	var dest map[string]any

	data, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrap(err, "[storage::Serialize]: can't marshal data")
	}

	if err = json.Unmarshal(data, &dest); err != nil {
		return nil, errors.Wrap(err, "[storage::Serialize]: can't unmarshal data")
	}

	dest["notify_interval"] = e.NotifyInterval

	return dest, nil
}
