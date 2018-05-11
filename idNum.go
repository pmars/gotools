package gotools

func CheckIDNumberRule(num string) bool {
	if len(num) != 18 {
		return false
	}
	// step1: 将前面的身份证号码17位数分别乘以不同的系数。从第一位到第十七位的系数分别为：7 9 10 5 8 4 2 1 6 3 7 9 10 5 8 4 2
	// step2: 将这17位数字和系数相乘的结果相加。
	sum := 0
	params := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	for i, p := range params {
		sum += int(num[i]-48) * p
		// logs.Debug("sum + num:%v*param:%d = %d", num[i], p, sum)
	}

	// step3: 用加出来和除以11，看余数是多少？
	// step4: 余数只可能有0 1 2 3 4 5 6 7 8 9 10这11个数字。其分别对应的最后一位身份证的号码为1 0 X 9 8 7 6 5 4 3 2。
	results := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	index := sum % 11

	// logs.Debug("sum:%d, index:%d, result[index]:%v, num:%s", sum, index, results[index]-48, num)
	return results[index] == num[17]
}
