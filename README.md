# auction-extension

拍卖防狙击延时：临近结束时出价会顺延截止时间，并受总顺延上限约束，只用 Go 标准库。

```
go run ./cmd/auction-extension --sample=window
go run ./cmd/auction-extension --sample=cap
go run ./cmd/auction-extension --sample=dup
go run ./cmd/auction-extension --sample=late
go run ./cmd/auction-extension --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **触发窗口**：只有在距截止时刻不超过 `windowMs`（含边界）时出价才顺延，否则截止时间不变。
- **顺延步长**：每次触发把截止时刻往后推 `stepMs`，累计顺延不超过 `maxExtMs`，超出部分不再顺延。
- **结束边界**：`atMs` 不小于当前截止时刻的出价一律拒绝，正好等于截止时刻算已结束。
- **幂等**：同一个 `bidID` 只生效一次，重复出价返回 false 且不改变截止时间与累计顺延。
- **不变量**：截止时刻单调不减，累计顺延不超过上限；同一串出价重复执行结果相同。
- **代价**：出价不许全扫历史出价，`Scanned` 不随出价数乘出价次数增长；10 万次出价每秒 10 万次处理。

## 输出契约（不改格式）

```
ends=..
ok=..
scanned=..
```

## 目录

```
auction.go            触发窗口、顺延上限与幂等
auction_test.go       单元测试
cmd/auction-extension/命令行入口
```
