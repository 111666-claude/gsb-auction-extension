// Command auction-extension 跑拍卖防狙击延时样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/auction"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("auction-extension", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "window", "window / cap / dup / late / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "window":
		book := auction.New(10000, 2000, 500, 100000)
		book.Bid("b1", "p", 100, 5000)
		fmt.Fprintf(stdout, "ends=%d\n", book.EndsAt())
	case "cap":
		book := auction.New(10000, 2000, 500, 1000)
		book.Bid("b1", "p", 100, 9500)
		book.Bid("b2", "p", 200, 9600)
		book.Bid("b3", "p", 300, 9700)
		fmt.Fprintf(stdout, "ends=%d\n", book.EndsAt())
	case "dup":
		book := auction.New(10000, 2000, 500, 100000)
		book.Bid("b1", "p", 100, 9500)
		book.Bid("b1", "p", 100, 9500)
		fmt.Fprintf(stdout, "ends=%d\n", book.EndsAt())
	case "late":
		book := auction.New(10000, 2000, 500, 100000)
		fmt.Fprintf(stdout, "ok=%v\n", book.Bid("b1", "p", 100, 10000))
	case "scan":
		book := auction.New(1000000, 2000, 500, 100000000)
		for index := 0; index < 3000; index++ {
			book.Bid(fmt.Sprintf("b-%d", index), "p", 100, 0)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", book.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample window|cap|dup|late|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
