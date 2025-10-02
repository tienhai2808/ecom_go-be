package initialization

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake/v2"
)

var Sf *sonyflake.Sonyflake

func InitSnowFlake() error {
	st := sonyflake.Settings{
		StartTime: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (int, error) {
			return 1, nil
		},
	}

	var err error
	Sf, err = sonyflake.New(st)
	if err != nil {
		return fmt.Errorf("khởi tạo Snowflake thất bại: %w", err)
	}

	return nil
}
