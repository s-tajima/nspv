package main

import (
	"fmt"

	"github.com/s-tajima/nspv"
)

func main() {
	v := nspv.NewValidator()
	v.SetDict([]string{"nspv"})
	v.SetMinLength(2)
	v.SetMaxLength(6)
	v.SetHibpThreshold(100)
	v.SetLevenshteinThreshold(2)

	res, _ := v.Validate("nspv")
	fmt.Println(res.String()) // ViolateDictCheck

	res, _ = v.Validate("n3pv")
	fmt.Println(res.String()) // ViolateDictCheck

	res, _ = v.Validate("n3pw")
	fmt.Println(res.String()) // Ok
}
