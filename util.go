package nathttp

import "encoding/json"

func GetJsonString(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	switch v := data.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		b, err := json.Marshal(data)
		return string(b), err
	}
}