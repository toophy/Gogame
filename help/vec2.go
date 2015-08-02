package help

import ()

type Vec2 struct {
	X float32
	Y float32
}

// 定制lua脚本
// 禁止以下功能
// dofile
// require
// load
// loadfile
// loadstring
