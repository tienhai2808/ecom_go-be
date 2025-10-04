package snowflake

type SnowflakeGenerator interface {
	NextID() (int64, error)
}