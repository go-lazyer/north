package nmap

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
