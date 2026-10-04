package auction

import "testing"

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
	book := New(10000, 2000, 500, 100000)
	if !book.Bid("b1", "p", 100, 5000) {
		t.Fatal("窗口外的出价仍应被接受")
	}
	if got := book.EndsAt(); got != 10000 {
		t.Fatalf("窗口外出价不应顺延，得到 ends=%d", got)
	}
	if got := book.Extended(); got != 0 {
		t.Fatalf("窗口外出价累计顺延应为 0，得到 %d", got)
	}
}

func TestBidAtWindowBoundaryExtends(t *testing.T) {
	book := New(10000, 2000, 500, 100000)
	book.Bid("b1", "p", 100, 8000)
	if got := book.EndsAt(); got != 10500 {
		t.Fatalf("正好距截止 windowMs 的出价应顺延，得到 ends=%d", got)
	}
}

func TestExtensionCappedByMaxExt(t *testing.T) {
	book := New(10000, 2000, 500, 1000)
	book.Bid("b1", "p", 100, 9500)
	book.Bid("b2", "p", 200, 9600)
	book.Bid("b3", "p", 300, 9700)
	if got := book.EndsAt(); got != 11000 {
		t.Fatalf("累计顺延不应超过 maxExtMs，得到 ends=%d", got)
	}
	if got := book.Extended(); got != 1000 {
		t.Fatalf("累计顺延应正好到上限 1000，得到 %d", got)
	}
}

func TestExtensionPartiallyAppliedUpToCap(t *testing.T) {
	book := New(10000, 2000, 500, 1200)
	book.Bid("b1", "p", 100, 9500)
	book.Bid("b2", "p", 200, 9600)
	book.Bid("b3", "p", 300, 9700)
	if got := book.Extended(); got != 1200 {
		t.Fatalf("超出上限的部分不再顺延，累计应为 1200，得到 %d", got)
	}
	if got := book.EndsAt(); got != 11200 {
		t.Fatalf("截止时刻应为 11200，得到 %d", got)
	}
}

func TestBidAtDeadlineIsRejected(t *testing.T) {
	book := New(10000, 2000, 500, 100000)
	if book.Bid("b1", "p", 100, 10000) {
		t.Fatal("正好等于截止时刻的出价应被拒绝")
	}
	if book.Bid("b2", "p", 100, 10001) {
		t.Fatal("超过截止时刻的出价应被拒绝")
	}
	if got := book.EndsAt(); got != 10000 {
		t.Fatalf("被拒绝的出价不应改变截止时刻，得到 ends=%d", got)
	}
}

func TestDuplicateBidIDOnlyCountsOnce(t *testing.T) {
	book := New(10000, 2000, 500, 100000)
	if !book.Bid("b1", "p", 100, 9500) {
		t.Fatal("首次出价应被接受")
	}
	if book.Bid("b1", "p", 100, 9500) {
		t.Fatal("同一 bidID 重复出价应返回 false")
	}
	if got := book.EndsAt(); got != 10500 {
		t.Fatalf("重复出价不应改变截止时刻，得到 ends=%d", got)
	}
	if got := book.Extended(); got != 500 {
		t.Fatalf("重复出价不应改变累计顺延，得到 %d", got)
	}
}

func TestReplaySameBidsGivesSameResult(t *testing.T) {
	run := func() (int, int) {
		book := New(10000, 2000, 500, 1000)
		book.Bid("b1", "p", 100, 9500)
		book.Bid("b2", "p", 200, 9600)
		book.Bid("b1", "p", 100, 9600)
		book.Bid("b3", "p", 300, 9700)
		return book.EndsAt(), book.Extended()
	}
	ends1, ext1 := run()
	ends2, ext2 := run()
	if ends1 != ends2 || ext1 != ext2 {
		t.Fatalf("同一串出价重复执行结果应相同：(%d,%d) vs (%d,%d)", ends1, ext1, ends2, ext2)
	}
}

func TestScannedGrowsLinearlyNotQuadratically(t *testing.T) {
	book := New(1000000, 2000, 500, 100000000)
	for index := 0; index < 3000; index++ {
		book.Bid(string(rune('a'+index%26))+string(rune('a'+index/26)), "p", 100, 0)
	}
	if got := book.Scanned(); got > 6000 {
		t.Fatalf("Scanned 不应随出价数乘出价次数增长，得到 %d", got)
	}
}
