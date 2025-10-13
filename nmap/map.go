package nmap

import "encoding/json"

// city=5YyX5Lqs; manualCity=5YyX5Lqs; token=7053d25529d55e028c3533f5fc9a0c58
// 转为这个样式
func ToString(m map[string]any) string {
	if m == nil {
		return ""
	}
	str := ""
	for k, v := range m {
		if str != "" {
			str += "; "
		}
		if v == nil {
			str += k + "="
		} else {
			str += k + "=" + v.(string)
		}
	}
	return str
}
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
func ToMap(obj any) (map[string]any, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
