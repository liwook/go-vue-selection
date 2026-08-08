package snowflake

import (
	"testing"
)

func TestGenID_BeforeInitPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when GenID called before Init")
		}
	}()
	// 确保 node 为 nil：通过新进程状态无法可靠重置 once，
	// 故仅在未初始化时验证 panic。若 Init 已被其它测试调用则跳过。
	if node != nil {
		t.Skip("node already initialized by another test; skipping panic check")
	}
	_ = GenID()
}

func TestInitAndGenID(t *testing.T) {
	if err := Init(1); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if node == nil {
		t.Fatal("node is nil after Init")
	}

	id := GenID()
	if id == 0 {
		t.Error("GenID returned 0")
	}
}

func TestInit_MultipleCallsOnce(t *testing.T) {
	// Init 使用 sync.Once，多次调用应安全且不改节点
	if err := Init(2); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	id1 := GenID()
	id2 := GenID()
	if id1 == id2 {
		t.Error("expected distinct IDs across calls")
	}
}

func TestGenID_Uniqueness(t *testing.T) {
	if err := Init(3); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	const n = 1000
	seen := make(map[int64]struct{}, n)
	for i := 0; i < n; i++ {
		id := GenID()
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate ID generated: %d", id)
		}
		seen[id] = struct{}{}
	}
}
