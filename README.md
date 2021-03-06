nspv - NIST SP 800-63B Validator
---
![Go](https://github.com/s-tajima/nspv/workflows/Go/badge.svg) [![GoDoc](https://godoc.org/github.com/s-tajima/nspv?status.svg)](https://godoc.org/github.com/s-tajima/nspv) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/3d5d752339c54d3dba8b71665f9b06c0)](https://www.codacy.com/manual/tajima1989/nspv?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=s-tajima/nspv&amp;utm_campaign=Badge_Grade)

nspv is a password validation library for Go, compatible with NIST Special Publication 800-63B.

## Description

[NIST Special Publication 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html) is a notable guideline for digital identity / authentication.  
nspv validates a password with policies based on this guideline. The policies are described below.

* Ensure the password length. (min 8 characters, max 64 characters, by default)
* Compare the password against a list that contains values known to be commonly-used, expected, or compromised. (use [Have I Been Pwned](https://haveibeenpwned.com/) internally)
* Judge whether the password could be predicable in the request context. (use Levenshtein Distance)

## Installation

```bash
go get -u github.com/s-tajima/nspv
```

## Usage

```go
v := nspv.NewValidator()
v.SetDict([]string{"nist-sp-800-63"})

res, err := v.Validate("_sup3r_comp1ex_passw0rd_")
if err != nil {
    // Something wrong with validataion. (e.g. HIBP API Error)
}
if res != nspv.Ok {
    // Validation failure.
}
fmt.Println(result.String()) // Ok

res, _ = v.Validate("short")
fmt.Println(res.String()) // ViolateMinLengthCheck

res, _ = v.Validate("password")
fmt.Println(res.String()) // ViolateHibpCheck

res, _ = v.Validate("n1st-sp-800-63")
fmt.Println(res.String()) // ViolateDictCheck
```

## License

[MIT](./LICENSE.md)

## Author

[Satoshi Tajima](https://github.com/s-tajima)
