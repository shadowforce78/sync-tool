package main

import "fmt"

func addition(a, b int) int {
    return a + b
}

func main() {
    var x, y int
    fmt.Println("Enter first number:")
    fmt.Scan(&x)
    fmt.Println("Enter second number:")
    fmt.Scan(&y)
    result := addition(x, y)
    fmt.Println("The sum is:", result)
}
