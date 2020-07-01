package route

// content tpye
const (
	JSON = "application/json"
)

// checkRole 检查是否越权操作 todo
func checkRole(role int, path string) bool {
	return true
}
