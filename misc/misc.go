// Funciones misceláneas
package misc

import (
	"strings"

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

// Compone parte de una clausula WHERE CAMPO IN (VALORES...)
// Devuelve "in (VALORES...)" o "= VALOR" o "is null"
func SqlIn(valores ...string) string {
	lista := []string{}
	vistos := map[string]any{}
	for _, v := range valores {
		_, visto := vistos[v]
		if !visto {
			lista = append(lista, v)
			vistos[v] = nil
		}
	}
	for k, v := range lista {
		lista[k] = EscapeSQL(v)
	}
	if len(lista) == 0 {
		return " is null"
	}
	if len(lista) == 1 {
		return " = " + lista[0]
	}
	return " in (" + strings.Join(lista, ",") + ")"
}

// Escapa un texto para su uso en órdenes SQL
func EscapeSQL(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
