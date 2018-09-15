package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"math"
	"sort"
)

type Bank struct {
	coins   []int
	amounts []int
}

func (bank *Bank) getCoins(sum int) []int {
	coins := change(bank.coins, bank.amounts, sum)
	withdrawCoins(bank, coins)
	return coins
}

func (bank *Bank) print() {
	for i, coin := range bank.coins {
		fmt.Println(coin, ":  ", bank.amounts[i])
	}
}

func withdrawCoins(bank *Bank, coins []int) {
	if coins == nil {
		return
	}
	for _, withdCoin := range coins {
		for j, coin := range bank.coins {
			if withdCoin == coin {
				bank.amounts[j]--
			}
		}
	}
}

func calcChange(T []int, amount []int, N int) []int {
	revAm := make([]int, len(T))
	copy(revAm, amount)
	funk.ReverseInt(revAm)

	rT := make([]int, len(T))
	copy(rT, T)
	sort.Sort(sort.Reverse(sort.IntSlice(rT)))

	A := make([]int, len(rT))
	for i, rTi := range rT {
		A[i] = N / rTi
		A[i] = int(math.Min(float64(A[i]), float64(revAm[i])))
		N = N - A[i]*rTi
	}
	return A
}

func change(T []int, amount []int, N int) []int {
	S := calcChange(T, amount, N)
	var res []int

	lenT := len(T) - 1
	for i := len(S) - 1; i >= 0; i-- {
		for j := 0; j < S[i]; j++ {
			res = append(res, T[lenT-i])
		}
	}
	var sum int
	for _, e := range res {
		sum += e
	}
	if sum == N {
		return res
	} else {
		return nil
	}
}
