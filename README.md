# Coupon
Simple coupon generating library written in [Go](https://golang.org). The algorithm for code generation is from [Grantm's CouponCode](https://github.com/grantm/Algorithm-CouponCode)

## Command line usage
    $ go get github.com/nullbus/coupon/coupon-gen
    $ coupon-gen

## Code usage
```
package main

import (
	"fmt"
	"github.com/nullbus/coupon"
)

func main() {
	generator := coupon.Generator{}
	fmt.Println("Yay! coupon code:", generator.Generate())
}
```

For more information about code, see [GoDoc](https://godoc.org/github.com/nullbus/coupon)

