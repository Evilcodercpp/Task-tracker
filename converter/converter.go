package converter

func CelsiusToFahrenheit(c float64) float64 { 
	F := (c * 9 / 5) + 32
	return F
}

func KilometersToMiles(km float64) float64 {
	miles := km * 0.621371
	return miles
}