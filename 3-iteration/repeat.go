package main

func Repeat(character string) string {
	var repated string
	for i := 0; i < 5; i++ {
		repated += character
	}
	return repated
}
