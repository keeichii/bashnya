package main

import "fmt"

func main() {
	var res int64
	fmt.Scan(&res)

	if res > 12307 {
		fmt.Println("higher than expected")
		return
	}

outer:
	for res < 12307 {
		switch {
		case res < 0:
			res = res * -1
		case res%7 == 0:
			res = res * 39
		case res%9 == 0:
			res = res*13 + 1
			continue
		default:
			res = (res + 2) * 3
		}

		if res%13 == 0 && res%9 == 0 {
			fmt.Println("service error")
			break outer
		} else {
			res++
		}
	}
	fmt.Println(res)
}
