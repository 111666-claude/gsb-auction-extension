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
