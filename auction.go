// Package auction 是拍卖防狙击延时：临近结束时出价会顺延截止时间，并受总顺延上限约束。
package auction

// Auction 是一口价拍卖的截止时间台账。
type Auction struct {
	windowMs int
	stepMs   int
	maxExtMs int
	endsAt   int
	extended int
	seen     map[string]struct{}
	scanned  int
}

// New 建拍卖：endsAt 是初始截止时刻，windowMs 是触发顺延的窗口，stepMs 是每次顺延步长，
// maxExtMs 是累计顺延上限。
func New(endsAt, windowMs, stepMs, maxExtMs int) *Auction {
	return &Auction{
		windowMs: windowMs,
		stepMs:   stepMs,
		maxExtMs: maxExtMs,
		endsAt:   endsAt,
		seen:     make(map[string]struct{}),
	}
}

// Bid 记录一次出价，返回是否被接受。
func (a *Auction) Bid(bidID, player string, amount, atMs int) bool {
	a.scanned++
	if _, dup := a.seen[bidID]; dup {
		return false
	}
	if atMs >= a.endsAt {
		return false
	}
	a.seen[bidID] = struct{}{}
	if a.endsAt-atMs <= a.windowMs {
		ext := a.stepMs
		if remaining := a.maxExtMs - a.extended; ext > remaining {
			ext = remaining
		}
		a.endsAt += ext
		a.extended += ext
	}
	return true
}

// EndsAt 返回当前截止时刻。
func (a *Auction) EndsAt() int { return a.endsAt }

// Extended 返回累计顺延的毫秒数。
func (a *Auction) Extended() int { return a.extended }

// Scanned 是累计扫描的历史出价数（规模观测）。
func (a *Auction) Scanned() int { return a.scanned }
