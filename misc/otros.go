// Funciones misceláneas
package misc

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Halla el mínimo de una lista
func Min[T constraints.Ordered](v ...T) T {
	min := v[0]
	for _, x := range v {
		if x < min {
			min = x
		}
	}
	return min
}

// Halla el máximo de una lista
func Max[T constraints.Ordered](v ...T) T {
	max := v[0]
	for _, x := range v {
		if x > max {
			max = x
		}
	}
	return max
}

// Hace lo mismo que la deprecada strings.Title: pasar la primera letra a mayúsculas
func Title(s string) string {
	return cases.Title(language.Und, cases.NoLower).String(s)
}
