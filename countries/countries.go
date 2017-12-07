package countries

type Country struct {
	Name,
	FullName,
	English,
	Alpha2,
	Alpha3,
	ISO,
	Location,
	LocationPrecise string
}

func ByAlpha2(alpha2 string) *Country {
	p,_ :=  byAlpha2[alpha2]
	return p
}


var byAlpha2 map[string] *Country

func init() {
	byAlpha2 = make(map[string] *Country)
	for i := range countries{
		byAlpha2[countries[i].Alpha2] = &countries[i]
	}
}