package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack struct {
	A []int
	B []int
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	stack, err := parseInput(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	if hasDuplicates(stack.A) {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	if isSorted(stack.A) {
		return
	}

	instructions := solve(stack)
	for _, instr := range instructions {
		fmt.Println(instr)
	}
}

func parseInput(input string) (*Stack, error) {
	values := strings.Fields(input)
	if len(values) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	stack := &Stack{
		A: make([]int, 0, len(values)),
		B: make([]int, 0, len(values)),
	}

	for _, v := range values {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		stack.A = append(stack.A, num)
	}
	return stack, nil
}

func hasDuplicates(nums []int) bool {
	seen := make(map[int]bool)
	for _, num := range nums {
		if seen[num] {
			return true
		}
		seen[num] = true
	}
	return false
}

func isSorted(nums []int) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i-1] > nums[i] {
			return false
		}
	}
	return true
}

func solve(s *Stack) []string {
	switch len(s.A) {
	case 0, 1:
		return []string{}
	case 2:
		return sortTwo(s)
	case 3:
		return sortThree(s)
	case 4, 5:
		return sortFive(s)
	case 6:
		return sortSix(s)
	default:
		return optimizedRadixSort(s)
	}
}

func sortTwo(s *Stack) []string {
	if s.A[0] > s.A[1] {
		return []string{"sa"}
	}
	return []string{}
}

func sortThree(s *Stack) []string {
	var instructions []string
	a, b, c := s.A[0], s.A[1], s.A[2]

	if a < b && b < c {
		return []string{}
	}

	if a == max(a, b, c) {
		instructions = append(instructions, "ra")
		rotate(&s.A)
	} else if b == max(a, b, c) {
		instructions = append(instructions, "rra")
		reverseRotate(&s.A)
	}

	if s.A[0] > s.A[1] {
		instructions = append(instructions, "sa")
	}
	return instructions
}

func sortFive(s *Stack) []string {
	var instructions []string
	for i := 0; i < 2; i++ {
		minIndex := findMinIndex(s.A)
		rotateToTop(s, &instructions, minIndex, "A")
		instructions = append(instructions, "pb")
		push(&s.A, &s.B)
	}

	instructions = append(instructions, sortThree(s)...)

	for len(s.B) > 0 {
		maxIndex := findMaxIndex(s.B)
		rotateToTop(s, &instructions, maxIndex, "B")
		instructions = append(instructions, "pa")
		push(&s.B, &s.A)
	}
	return instructions
}

func sortSix(s *Stack) []string {
	var instructions []string
	for i := 0; i < 3; i++ {
		minIndex := findMinIndex(s.A)
		rotateToTop(s, &instructions, minIndex, "A")
		instructions = append(instructions, "pb")
		push(&s.A, &s.B)
	}

	instructions = append(instructions, sortThree(s)...)

	for len(s.B) > 0 {
		maxIndex := findMaxIndex(s.B)
		rotateToTop(s, &instructions, maxIndex, "B")
		instructions = append(instructions, "pa")
		push(&s.B, &s.A)
	}
	return instructions
}

func optimizedRadixSort(s *Stack) []string {
	var instructions []string

	// If stack is empty or has only one element, return empty instructions
	if len(s.A) <= 1 {
		return instructions
	}

	offset := adjustNegatives(s.A)
	maxBits := calculateMaxBits(s.A)

	// Safety check to prevent infinite loops
	if maxBits > 32 {
		maxBits = 32 // Cap at 32 bits since we're dealing with integers
	}

	for bit := 0; bit < maxBits; bit++ {
		// Count numbers that will go to stack B (0 bit)
		zeroCount := 0
		for i := 0; i < len(s.A); i++ {
			if ((s.A[i]+offset)>>bit)&1 == 0 {
				zeroCount++
			}
		}

		// Push numbers with 0 at current bit to B
		pushed := 0
		for pushed < zeroCount {
			if ((s.A[0]+offset)>>bit)&1 == 0 {
				instructions = append(instructions, "pb")
				push(&s.A, &s.B)
				pushed++
			} else {
				instructions = append(instructions, "ra")
				rotate(&s.A)
			}
		}

		// Push everything back to A
		for len(s.B) > 0 {
			instructions = append(instructions, "pa")
			push(&s.B, &s.A)
		}
	}

	// Correct for negative numbers if needed
	correctNegatives(s, &instructions, offset)
	return instructions
}

func adjustNegatives(nums []int) int {
	minVal := findMin(nums)
	if minVal >= 0 {
		return 0
	}
	offset := -minVal
	for i := range nums {
		nums[i] += offset
	}
	return offset
}

func calculateMaxBits(nums []int) int {
	maxVal := findMax(nums)
	bits := 0
	for maxVal > 0 {
		bits++
		maxVal >>= 1
	}
	return bits + 1
}

func correctNegatives(s *Stack, instructions *[]string, offset int) {
	if offset == 0 {
		return
	}
	splitIndex := 0
	for splitIndex < len(s.A) && s.A[splitIndex] >= offset {
		splitIndex++
	}
	if splitIndex > 0 {
		rotateCount := splitIndex
		if rotateCount <= len(s.A)/2 {
			addReverseRotateA(s, instructions, rotateCount)
		} else {
			addRotateA(s, instructions, len(s.A)-rotateCount)
		}
	}
	for i := range s.A {
		s.A[i] -= offset
	}
}

func addReverseRotateA(s *Stack, instructions *[]string, count int) {
	for i := 0; i < count; i++ {
		*instructions = append(*instructions, "rra")
		reverseRotate(&s.A)
	}
}

func addRotateA(s *Stack, instructions *[]string, count int) {
	for i := 0; i < count; i++ {
		*instructions = append(*instructions, "ra")
		rotate(&s.A)
	}
}

func findMin(nums []int) int {
	minVal := nums[0]
	for _, num := range nums {
		if num < minVal {
			minVal = num
		}
	}
	return minVal
}

func findMax(nums []int) int {
	maxVal := nums[0]
	for _, num := range nums {
		if num > maxVal {
			maxVal = num
		}
	}
	return maxVal
}

func max(a, b, c int) int {
	if a >= b && a >= c {
		return a
	}
	if b >= a && b >= c {
		return b
	}
	return c
}

func findMinIndex(nums []int) int {
	minIndex := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] < nums[minIndex] {
			minIndex = i
		}
	}
	return minIndex
}

func findMaxIndex(nums []int) int {
	maxIndex := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] > nums[maxIndex] {
			maxIndex = i
		}
	}
	return maxIndex
}

func rotateToTop(s *Stack, instructions *[]string, index int, stackName string) {
	var stack *[]int
	var rot, revRot string

	if stackName == "A" {
		stack = &s.A
		rot, revRot = "ra", "rra"
	} else {
		stack = &s.B
		rot, revRot = "rb", "rrb"
	}

	n := len(*stack)
	if index <= n/2 {
		for i := 0; i < index; i++ {
			*instructions = append(*instructions, rot)
			rotate(stack)
		}
	} else {
		for i := 0; i < n-index; i++ {
			*instructions = append(*instructions, revRot)
			reverseRotate(stack)
		}
	}
}

func rotate(stack *[]int) {
	if len(*stack) == 0 {
		return
	}
	*stack = append((*stack)[1:], (*stack)[0])
}

func reverseRotate(stack *[]int) {
	if len(*stack) == 0 {
		return
	}
	last := len(*stack) - 1
	*stack = append([]int{(*stack)[last]}, (*stack)[:last]...)
}

func swap(stack *[]int) {
	if len(*stack) >= 2 {
		(*stack)[0], (*stack)[1] = (*stack)[1], (*stack)[0]
	}
}

func push(src *[]int, dst *[]int) {
	if len(*src) == 0 {
		return
	}
	*dst = append([]int{(*src)[0]}, *dst...)
	*src = (*src)[1:]
}
