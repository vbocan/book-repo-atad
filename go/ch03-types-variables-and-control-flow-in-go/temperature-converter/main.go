package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: tempconv <value> <unit>")
        fmt.Println("Units: C, F, K")
        os.Exit(1)
    }

    value, err := strconv.ParseFloat(os.Args[1], 64)
    if err != nil {
        fmt.Printf("Invalid temperature: %s\n", os.Args[1])
        os.Exit(1)
    }

    unit := strings.ToUpper(os.Args[2])

    var celsius float64
    switch unit {
    case "C":
        celsius = value
    case "F":
        celsius = (value - 32) * 5 / 9
    case "K":
        celsius = value - 273.15
    default:
        fmt.Printf("Unknown unit: %s\n", unit)
        os.Exit(1)
    }

    fahrenheit := celsius*9/5 + 32
    kelvin := celsius + 273.15

    fmt.Printf("%.2f°C = %.2f°F = %.2fK\n", celsius, fahrenheit, kelvin)
}
