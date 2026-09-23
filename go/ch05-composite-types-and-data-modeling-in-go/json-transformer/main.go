// Command json-transformer is a reference solution for Exercise 5.3. It
// reads a JSON array of products from a file, keeps only the in-stock ones,
// converts their prices from USD to EUR, and writes the result as
// pretty-printed JSON to another file.
//
// Usage:
//
//	json-transformer <input.json> <output.json>
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// usdToEur is the fixed exchange rate used to convert prices. A real
// program would look this up from a live source; a constant keeps this
// example self-contained and reproducible.
const usdToEur = 0.92

// Product mirrors one element of the input JSON array. The struct tags map
// the idiomatic Go field names to the snake_case keys used in the file.
type Product struct {
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: json-transformer <input.json> <output.json>")
		os.Exit(1)
	}
	inputPath, outputPath := os.Args[1], os.Args[2]

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		fmt.Printf("Error parsing %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	// Filter to in-stock items and convert their price, in place, into a
	// freshly made slice (rather than a nil one) so an all-filtered result
	// still marshals to "[]" instead of "null".
	kept := make([]Product, 0, len(products))
	for _, p := range products {
		if !p.InStock {
			continue
		}
		p.Price = roundToCents(p.Price * usdToEur)
		kept = append(kept, p)
	}

	out, err := json.MarshalIndent(kept, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding result: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Kept %d of %d in-stock items, converted USD to EUR at %.2f, wrote %s\n",
		len(kept), len(products), usdToEur, outputPath)
}

// roundToCents rounds a price to two decimal places, which also absorbs the
// floating-point noise that a raw multiplication like 15.50 * 0.92 leaves
// behind (see the floating-point caveats in Section 3.1.2).
func roundToCents(price float64) float64 {
	return math.Round(price*100) / 100
}
