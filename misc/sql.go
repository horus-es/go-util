package misc

import (
	"slices"
	"strings"
)

// Compone parte de una clausula WHERE CAMPO IN (VALORES...)
// Devuelve "in (VALORES...)" o "= VALOR" o "is null"
func SqlIn(valores ...string) string {
	lista := []string{}
	for _, v := range valores {
		if v == "" || slices.Contains(lista, v) {
			continue
		}
		lista = append(lista, v)
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
