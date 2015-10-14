package main

import (
	"flag"
	"fmt"

	"github.com/nullbus/coupon"
)

var (
	count    int
	parts    int
	seedText string
)

func init() {
	flag.IntVar(&count, "count", 1, "number of generated coupons")
	flag.IntVar(&parts, "part", 3, "number of part")
	flag.StringVar(&seedText, "text", "", "seed for randomize")
}

func main() {
	flag.Parse()

	// store unique coupon code
	codes := map[string]bool{}
	validator := coupon.Validator{
		NumParts: parts,
	}

	for len(codes) < count {
		generator := coupon.Generator{
			NumParts:  parts,
			PlainText: fmt.Sprintf("%s%d", seedText, len(codes)),
		}

		code := generator.Generate()
		if _, err := validator.Validate(code); nil != err {
			fmt.Printf("validate error, %s , %s\n", code, err.Error())
			return
		}

		codes[code] = true
		fmt.Println(code)
	}

	fmt.Println(len(codes), "generated")
}
