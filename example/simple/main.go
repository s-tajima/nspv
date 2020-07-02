package main

import (
	"fmt"

	"github.com/s-tajima/nspv"
)

func main() {
	v := nspv.NewValidator()

	res, _ := v.Validate("_sup3r_comp1ex_passw0rd_")
	fmt.Println(res.String()) // Ok

	res, _ = v.Validate("short")
	fmt.Println(res.String()) // ViolateMinLengthCheck

	res, _ = v.Validate("password")
	fmt.Println(res.String()) // ViolateBibpCheck
}
