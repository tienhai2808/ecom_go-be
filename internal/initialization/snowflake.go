package initialization

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake/v2"
)

func InitSnowFlake() (*sonyflake.Sonyflake, error) {
	st := sonyflake.Settings{
		StartTime: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (int, error) {
			return 1, nil
		},
	}

	sf, err := sonyflake.New(st)
	if err != nil {
		return nil, fmt.Errorf("khởi tạo Snowflake thất bại: %w", err)
	}

	return sf, nil
}
