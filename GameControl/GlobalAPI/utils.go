package GlobalAPI

import "github.com/google/uuid"

// 生成一个新的 uuid 对象并返回
func generateUUID() uuid.UUID {
	for {
		uniqueId, err := uuid.NewUUID()
		if err != nil {
			continue
		}
		return uniqueId
	}
}
