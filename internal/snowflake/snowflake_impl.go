package snowflake

import (
	"fmt"

	"github.com/sony/sonyflake/v2"
)

type snowflakeGeneratorImpl struct {
	sf *sonyflake.Sonyflake
}

func NewSnowflakeGenerator(sf *sonyflake.Sonyflake) SnowflakeGenerator {
	return &snowflakeGeneratorImpl{sf}
}

func (g *snowflakeGeneratorImpl) NextID() (int64, error) {
	id, err := g.sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("tạo Snowflake ID thất bại: %w", err)
	}

	return id, nil
}