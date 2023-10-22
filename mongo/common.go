package mongo

const (
	tokenTmp = "w\u0000\u0000\u0000\u0002_data\u0000g\u0000\u0000\u0000%s\u0000\u0000"
)

// mongodb 操作类型
const (
	OperationTypeInsert  = "insert"
	OperationTypeDelete  = "delete"
	OperationTypeUpdate  = "update"
	OperationTypeReplace = "replace"
)

func SliceToSet(slice []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, val := range slice {
		result[val] = struct{}{}
	}
	return result
}
