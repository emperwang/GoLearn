package statement

// 定义一个全局匿名函数
var (
	// 匿名函数返回 两个值中较小的那个
	Func1 = func (num1, num2 int) int {
		if num1 < num2 {
			return num1
		}else {
			return num2
		}
	}
)

// 返回两个值中较大的那个
func Max(num1, num2 int) int {

	if num1 > num2 {
		return num1
	}else {
		return num2
	}
}

// 闭包函数
func GetAddFunc() func (int) int {
	beginVal := 0
	return func(step int) int {
		beginVal += step
		return beginVal
	}
}

// 接收函数为参数的函数
func ProcessNum(num1, num2 int, process func (int, int) int) int {
	res := process(num1,num2)
	return res
}




