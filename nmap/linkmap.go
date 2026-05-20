package nmap

import "fmt"

type LinkMap[K comparable, V any] struct {
	m    map[K]V
	keys []K
}

func NewLinkMap[K comparable, V any]() *LinkMap[K, V] {
	return &LinkMap[K, V]{m: make(map[K]V)}
}

// Set：已存在则更新值，不改变顺序；不存在则追加
func (o *LinkMap[K, V]) Set(k K, v V) {
	if _, ok := o.m[k]; ok {
		o.m[k] = v
		return
	}
	o.m[k] = v
	o.keys = append(o.keys, k)
}

func (o *LinkMap[K, V]) Get(k K) (V, bool) {
	v, ok := o.m[k]
	return v, ok
}

func (o *LinkMap[K, V]) Delete(k K) {
	if _, ok := o.m[k]; !ok {
		return
	}
	delete(o.m, k)
	// 从 keys 中移除
	for i, key := range o.keys {
		if key == k {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
}

// Range：按插入顺序遍历
func (o *LinkMap[K, V]) Range(fn func(k K, v V)) {
	for _, k := range o.keys {
		fn(k, o.m[k])
	}
}

// ToString：按插入顺序返回字符串表示
// 例如 a:1; b:2; c:3; 其中:和;可自定义
// seps[0]：kv 分隔符，默认为:；seps[1]：项分隔符，默认为“;”；
func (o *LinkMap[K, V]) ToString(seps ...string) string {
	var kvSep, itemSep string
	if len(seps) > 0 {
		kvSep = seps[0]
	} else {
		kvSep = ":"
	}
	if len(seps) > 1 {
		itemSep = seps[1]
	} else {
		itemSep = ";"
	}
	var result string
	for i, k := range o.keys {
		result += fmt.Sprintf("%v%v%v", k, kvSep, o.m[k])
		if i < len(o.keys)-1 {
			result += itemSep
		}
	}
	return result
}
