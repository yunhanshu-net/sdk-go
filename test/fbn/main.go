package main

import (
	"fmt"
	"math/big"
	"time"
)

// 快速倍增法计算斐波那契数 F(n)
func fib(n int) *big.Int {
	a := big.NewInt(0)
	b := big.NewInt(1)
	var c, d *big.Int

	// 分解 n 的二进制位，存储到栈中
	stack := []int{}
	for n > 0 {
		stack = append(stack, n)
		n >>= 1 // 右移一位（等价于除以2）
	}

	// 逆序遍历二进制位
	for i := len(stack) - 1; i >= 0; i-- {
		k := stack[i]
		c = new(big.Int).Mul(a, new(big.Int).Sub(new(big.Int).Mul(b, big.NewInt(2)), a)) // c = a*(2b - a)
		d = new(big.Int).Add(new(big.Int).Mul(a, a), new(big.Int).Mul(b, b))             // d = a² + b²

		if k&1 == 1 { // 如果当前位是奇数
			a.Set(d)    // a = d
			b.Add(c, d) // b = c + d
		} else { // 如果当前位是偶数
			a.Set(c) // a = c
			b.Set(d) // b = d
		}
	}
	return a
}

func main() {
	n := 70000

	// 计算并计时
	start := time.Now()
	result := fib(n)
	elapsed := time.Since(start)

	// 输出结果摘要（避免打印超长数字）
	strResult := result.String()
	length := len(strResult)
	fmt.Printf("F(%d) 的位数：%d\n", n, length)
	fmt.Printf("首 10 位：%s...\n", strResult[:10])
	fmt.Printf("末 10 位：...%s\n", strResult[length-10:])
	fmt.Printf("计算耗时：%v\n", elapsed)
}
