package store

import "github.com/bwmarrin/snowflake"

// Snowflake节点
var SnowFlake *snowflake.Node

func InitSnowFlakeNode() {
	var err error
	SnowFlake, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}
