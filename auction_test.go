package auction

import (
	"strconv"
	"testing"
)

func TestNewKeepsDeadline(t *testing.T) {
	if got := New(10000, 2000, 500, 1000).EndsAt(); got != 10000 {
		t.Fatalf("新建时截止时刻应该就是传入值，得到 %d", got)
	}
}

func TestBidIsAccepted(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	if !book.Bid("b1", "p", 100, 1000) {
		t.Fatal("拍卖进行中的出价应该被接受")
	}
}

func TestExtendedStartsAtZero(t *testing.T) {
	if New(10000, 2000, 500, 1000).Extended() != 0 {
		t.Fatal("还没有出价时累计顺延应该是 0")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New(10000, 2000, 500, 1000).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestBidOutsideWindowDoesNotExtend(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	if !book.Bid("early", "p", 100, 7999) {
		t.Fatal("窗口外的正常出价仍然应该被接受")
	}
	if got := book.EndsAt(); got != 10000 {
		t.Fatalf("窗口外出价不应顺延，ends=%d", got)
	}
	if got := book.Extended(); got != 0 {
		t.Fatalf("窗口外出价不应增加累计顺延，extended=%d", got)
	}
}

func TestBidAtWindowBoundaryExtends(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	// 距截止正好 windowMs（含边界）应触发。
	if !book.Bid("edge", "p", 100, 8000) {
		t.Fatal("窗口边界上的出价应该被接受")
	}
	if got := book.EndsAt(); got != 10500 {
		t.Fatalf("窗口边界出价应顺延一步，ends=%d", got)
	}
	if got := book.Extended(); got != 500 {
		t.Fatalf("累计顺延应为 500，得到 %d", got)
	}
}

func TestExtensionCappedAtMaxExt(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	book.Bid("b1", "p", 100, 9500)
	book.Bid("b2", "p", 200, 9600)
	book.Bid("b3", "p", 300, 9700)
	if got := book.EndsAt(); got != 11000 {
		t.Fatalf("累计顺延应被上限截断在 11000，得到 %d", got)
	}
	if got := book.Extended(); got != 1000 {
		t.Fatalf("累计顺延不应超过上限，extended=%d", got)
	}
	// 上限用尽后，窗口内出价仍然有效，但不再顺延。
	if !book.Bid("b4", "p", 400, 9800) {
		t.Fatal("达到顺延上限后的出价仍然应该被接受")
	}
	if got := book.EndsAt(); got != 11000 {
		t.Fatalf("上限用尽后截止时刻不应再变，ends=%d", got)
	}
	if got := book.Extended(); got != 1000 {
		t.Fatalf("上限用尽后累计顺延不应再变，extended=%d", got)
	}
}

func TestPartialExtensionWhenNearCap(t *testing.T) {
	book := New(10000, 2000, 500, 700)
	book.Bid("b1", "p", 100, 9500) // 已延 500，余 200
	if !book.Bid("b2", "p", 200, 9600) {
		t.Fatal("窗口内出价应该被接受")
	}
	if got := book.EndsAt(); got != 10700 {
		t.Fatalf("最后一步应只顺延剩余的 200，ends=%d", got)
	}
	if got := book.Extended(); got != 700 {
		t.Fatalf("累计顺延应正好等于上限 700，得到 %d", got)
	}
}

func TestBidAtOrAfterDeadlineRejected(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	if book.Bid("at", "p", 100, 10000) {
		t.Fatal("正好等于截止时刻应算已结束，出价必须拒绝")
	}
	if book.Bid("after", "p", 100, 10001) {
		t.Fatal("晚于截止时刻的出价必须拒绝")
	}
	if got := book.EndsAt(); got != 10000 {
		t.Fatalf("被拒绝的出价不应改变截止时刻，ends=%d", got)
	}
	if got := book.Extended(); got != 0 {
		t.Fatalf("被拒绝的出价不应增加累计顺延，extended=%d", got)
	}
}

func TestDuplicateBidIDIsIdempotent(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	if !book.Bid("b1", "p", 100, 9500) {
		t.Fatal("首次出价应该被接受")
	}
	if book.Bid("b1", "p", 100, 9500) {
		t.Fatal("同一 bidID 重复出价必须返回 false")
	}
	if got := book.EndsAt(); got != 10500 {
		t.Fatalf("重复出价不应再次顺延，ends=%d", got)
	}
	if got := book.Extended(); got != 500 {
		t.Fatalf("重复出价不应改变累计顺延，extended=%d", got)
	}
}

func TestRepeatedSequenceHasSameResult(t *testing.T) {
	run := func() (ends, extended int) {
		book := New(10000, 2000, 500, 1000)
		book.Bid("b1", "p", 100, 1000) // 窗口外，不顺延
		book.Bid("b2", "p", 200, 9500) // 延 500
		book.Bid("b2", "p", 200, 9500) // 重复，无效
		book.Bid("b3", "p", 300, 10500)
		book.Bid("b3", "p", 300, 10500) // 结束后，拒绝
		book.Bid("b4", "p", 400, 9600)  // 再延 500，触顶
		book.Bid("b5", "p", 500, 9700)  // 上限用尽，不顺延
		return book.EndsAt(), book.Extended()
	}
	ends1, ext1 := run()
	ends2, ext2 := run()
	if ends1 != ends2 || ext1 != ext2 {
		t.Fatalf("同一串出价重复执行结果应相同：(%d,%d) vs (%d,%d)", ends1, ext1, ends2, ext2)
	}
	if ends1 != 11000 || ext1 != 1000 {
		t.Fatalf("重复执行结果应为 (11000,1000)，得到 (%d,%d)", ends1, ext1)
	}
}

func TestScannedIsLinearInBidCount(t *testing.T) {
	book := New(1000000, 2000, 500, 100000000)
	const n = 3000
	for index := 0; index < n; index++ {
		book.Bid("b-"+strconv.Itoa(index), "p", 100, 0)
	}
	if got := book.Scanned(); got > 2*n {
		t.Fatalf("Scanned 不应随出价数平方增长，得到 %d，上限 %d", got, 2*n)
	}
}

func Benchmark100kBids(b *testing.B) {
	for i := 0; i < b.N; i++ {
		book := New(100000000, 2000, 500, 100000000)
		for index := 0; index < 100000; index++ {
			book.Bid("b-"+strconv.Itoa(index), "p", 100, 0)
		}
	}
}
