package mathpack


func Pow(base, exponent int) int{
	result := 1
	for i := 1; i <= exponent; i++ {
		result *= base
	}
	return result
}

func IsPrime(n int) bool{
	y := 0
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			y++
		}
	}
	return y > 2
}