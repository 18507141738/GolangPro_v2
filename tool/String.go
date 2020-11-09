package tool

import "fmt"

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("首字母不是小写")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func Lowercase(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 {
				vv[i] += 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("首字母不是大写")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
