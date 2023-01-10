package misc

import "strings"

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
