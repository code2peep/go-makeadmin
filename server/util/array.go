package util

var ArrayUtil = arrayUtil{}

//arrayUtil 数组工具类
type arrayUtil struct{}

//ListToTree 字典列表转树形结构
func (au arrayUtil) ListToTree(arr []map[string]interface{}, id string, pid string, child string) (mapList []interface{}) {
	mapList = []interface{}{}
	// 遍历以id_为key生成map
	idValMap := make(map[uint64]map[string]interface{}, len(arr))
	for _, m := range arr {
		if idVal, ok := treeUint(m[id]); ok {
			idValMap[idVal] = m
		}
	}
	// 遍历
	for _, m := range arr {
		// 获取父节点
		if pidVal, ok := treeUint(m[pid]); ok && pidVal != 0 {
			if pNode, pok := idValMap[pidVal]; pok {
				// 有父节点则添加到父节点子集
				children, _ := pNode[child].([]interface{})
				pNode[child] = append(children, m)
				continue
			}
		}
		mapList = append(mapList, m)
	}
	return
}

func treeUint(value interface{}) (uint64, bool) {
	switch v := value.(type) {
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case int:
		if v >= 0 {
			return uint64(v), true
		}
	case int8:
		if v >= 0 {
			return uint64(v), true
		}
	case int16:
		if v >= 0 {
			return uint64(v), true
		}
	case int32:
		if v >= 0 {
			return uint64(v), true
		}
	case int64:
		if v >= 0 {
			return uint64(v), true
		}
	}
	return 0, false
}
