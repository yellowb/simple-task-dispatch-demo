package status

import (
	"fmt"
	"strings"
)

func CheckStatus(curStatus Status, statusList ...Status) error {
	for _, status := range statusList {
		if curStatus == status {
			return nil
		}
	}
	// 构造错误信息
	statusStrList := make([]string, 0, len(statusList))
	for _, status := range statusList {
		statusStrList = append(statusStrList, status.String())
	}
	return fmt.Errorf("status must be in %s, current status: %s", strings.Join(statusStrList, "/"), curStatus.String())
}
