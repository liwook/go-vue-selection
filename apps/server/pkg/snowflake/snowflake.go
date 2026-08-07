package snowflake

import (
	"sync"

	sf "github.com/bwmarrin/snowflake"
)

var (
	node *sf.Node
	once sync.Once
)

// Init 初始化雪花 ID 节点，使用 bwmarrin/snowflake 的默认配置
// （Epoch 默认 2010-11-04，NodeBits=10，StepBits=12）。
// 通过 sync.Once 保证只初始化一次；多次调用只会生效第一次。
// 前后端 ID 均以 string 传输，无需压缩位数来兼容 JS 数值精度。
func Init(machinedID int64) (err error) {
	once.Do(func() {
		node, err = sf.NewNode(machinedID)
	})
	return
}

// GenID 生成下一个雪花 ID。必须先调用 Init 初始化节点，否则会 panic。
func GenID() int64 {
	if node == nil {
		panic("snowflake: GenID called before Init")
	}
	return node.Generate().Int64()
}

//func main() {
//	if err := Init("2020-07-01", 1); err != nil {
//		fmt.Printf("init failed, err:%v \n", err)
//		return
//	}
//	id := GenID()
//	fmt.Println(id)
//}
