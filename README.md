# Push-Swap Project

## Overview

This project involves implementing a simple **non-comparative sorting algorithm** using two stacks (`a` and `b`) and a restricted set of operations. It includes two separate Go programs:

- `push-swap`: Outputs a list of instructions to sort a given stack `a`.
- `checker`: Verifies if the output instructions correctly sort the stack.

The goal is to sort integers in **ascending order** with the **least number of operations** possible.

---

## Stack Operations

You may use the following stack instructions:

| Command | Description |
|---------|-------------|
| `pa`    | Push top of stack **b** to stack **a** |
| `pb`    | Push top of stack **a** to stack **b** |
| `sa`    | Swap first 2 elements of stack **a** |
| `sb`    | Swap first 2 elements of stack **b** |
| `ss`    | Execute `sa` and `sb` simultaneously |
| `ra`    | Rotate stack **a** upwards |
| `rb`    | Rotate stack **b** upwards |
| `rr`    | Execute `ra` and `rb` simultaneously |
| `rra`   | Reverse rotate **a** |
| `rrb`   | Reverse rotate **b** |
| `rrr`   | Execute `rra` and `rrb` simultaneously |

---

## Usage

### 1. Build the executables:

```bash
cd pushswap
go build -o push-swap

cd checker
go build -o checker

mv pushswap/push-swap checker // To enable you run from the checker directory as required
```


### 2. Run the program: 

```bash
./push-swap "2 1 3 6 5 8" //input your stack here
```

### 3. Run the checker:

```bash
./checker "3 2 1 0" //input your stack here
```

### 4. General usage:

```bash
ARG="4 3 5 1 2"; ./push-swap "$ARG" | wc -l              //gives the word count of the number of instructions used to sort the stack.
ARG="4 3 5 1 2"; ./push-swap "$ARG" | ./checker "$ARG"   //gives the OK if the checker passes and KO if the checker fails.
```