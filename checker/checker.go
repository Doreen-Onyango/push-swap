package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var allowedInstructions = map[string]bool{
	"sa":  true,
	"sb":  true,
	"ss":  true,
	"pa":  true,
	"pb":  true,
	"ra":  true,
	"rb":  true,
	"rr":  true,
	"rra": true,
	"rrb": true,
	"rrr": true,
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	stackA, err := parseInput(os.Args[1])
	if err != nil || hasDuplicates(stackA) {
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}

	stackB := []int{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !allowedInstructions[line] {
			fmt.Fprintln(os.Stderr, "Error")
			os.Exit(1)
		}
		applyInstruction(line, &stackA, &stackB)
	}

	if isSorted(stackA) && len(stackB) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

func parseInput(input string) ([]int, error) {
	parts := strings.Fields(input)
	stack := make([]int, 0, len(parts))
	for _, val := range parts {
		num, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		stack = append(stack, num)
	}
	return stack, nil
}

func hasDuplicates(nums []int) bool {
	seen := make(map[int]bool)
	for _, n := range nums {
		if seen[n] {
			return true
		}
		seen[n] = true
	}
	return false
}

func isSorted(stack []int) bool {
	for i := 0; i < len(stack)-1; i++ {
		if stack[i] > stack[i+1] {
			return false
		}
	}
	return true
}

func applyInstruction(instr string, a *[]int, b *[]int) {
	switch instr {
	case "sa":
		swap(a)
	case "sb":
		swap(b)
	case "ss":
		swap(a)
		swap(b)
	case "pa":
		push(b, a)
	case "pb":
		push(a, b)
	case "ra":
		rotate(a)
	case "rb":
		rotate(b)
	case "rr":
		rotate(a)
		rotate(b)
	case "rra":
		reverseRotate(a)
	case "rrb":
		reverseRotate(b)
	case "rrr":
		reverseRotate(a)
		reverseRotate(b)
	}
}

func swap(stack *[]int) {
	if len(*stack) >= 2 {
		(*stack)[0], (*stack)[1] = (*stack)[1], (*stack)[0]
	}
}

func push(src, dst *[]int) {
	if len(*src) == 0 {
		return
	}
	val := (*src)[0]
	*src = (*src)[1:]
	*dst = append([]int{val}, *dst...)
}

func rotate(stack *[]int) {
	if len(*stack) > 0 {
		first := (*stack)[0]
		*stack = append((*stack)[1:], first)
	}
}

func reverseRotate(stack *[]int) {
	if len(*stack) > 0 {
		last := (*stack)[len(*stack)-1]
		*stack = append([]int{last}, (*stack)[:len(*stack)-1]...)
	}
}
