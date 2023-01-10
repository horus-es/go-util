// Funciones de gestión para POSTGRESQL usando el driver pgxpool
package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/dbscan"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/horus-es/go-util/errores"
	"github.com/horus-es/go-util/formato"
	"github.com/horus-es/go-util/logger"
	"github.com/horus-es/go-util/misc"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool
var dbLog *logger.Logger
var inTest bool

// Conecta a la base de datos y establece el logger. Si el logger es nil, todos los mensajes se meustran por consola.
func InitPool(connectString string, logger *logger.Logger) {
	var err error
	dbPool, err = pgxpool.Connect(context.Background(), connectString)
	errores.PanicIfError(err, "Error conectando base de datos")
	dbLog = logger
}

// Comienza una transacción
func StartTX() pgx.Tx {
	tx, err := dbPool.Begin(context.Background())
	errores.PanicIfError(err, "StartTX")
	dbLog.Infof("StartTX")
	return tx
}

// Finaliza una transacción
func CommitTX(tx pgx.Tx) {
	err := tx.Commit(context.Background())
	errores.PanicIfError(err, "CommitTX")
	dbLog.Infof("CommitTX")
}

// Aborta una transacción
func RollbackTX(tx pgx.Tx) {
	err := tx.Rollback(context.Background())
	if err == nil && dbLog != nil {
		dbLog.Warnf("RollbackTX")
	}
}

// Función de utilidad para consultas que devuelven exactamente una fila.
// dst puede ser la direccion de una struct o de una variable simple
// Panic si la query devuelve mas de una fila o no devuelve ninguna fila.
func GetOneRow(dst any, query string, params ...any) {
	limpio := reemplaza(query, params...)
	if strings.HasPrefix(strings.ToLower(limpio), "select * from ") {
		query = replaceAsterisk(query, dst)
	}
	rows, err := dbPool.Query(context.Background(), query, params...)
	errores.PanicIfError(err, "GetOneRow: %s", limpio)
	err = pgxscan.ScanOne(dst, rows)
	errores.PanicIfError(err, "GetOneRow: %s", limpio)
	dbLog.Infof(limpio)
}

// Función de utilidad para consultas que solo pueden devolver una (resultado true)
// o ninguna filas (resultado false).
// Panic si la query devuelve mas de una fila.
func GetOneOrZeroRows(dst any, query string, params ...any) bool {
	limpio := reemplaza(query, params...)
	if strings.HasPrefix(strings.ToLower(limpio), "select * from ") {
		query = replaceAsterisk(query, dst)
	}
	rows, err := dbPool.Query(context.Background(), query, params...)
	errores.PanicIfError(err, "GetOneOrZeroRows: %s", limpio)
	err = pgxscan.ScanOne(dst, rows)
	if pgxscan.NotFound(err) {
		dbLog.Infof(limpio)
		return false
	}
	errores.PanicIfError(err, "GetOneOrZeroRows: %s", limpio)
	dbLog.Infof(limpio)
	return true
}

// Función de utilidad que devuelve la primera fila de una consulta (devuelve true).
// Si no hay filas, devuelve false.
// Panic si la query no contiene la cláusula "order by".
func GetFirstRow(dst any, query string, params ...any) bool {
	limpio := reemplaza(query, params...)
	isOrdered := strings.Contains(strings.ToLower(limpio), " order by ")
	errores.PanicIfTrue(!isOrdered, "GetFirstRow: Debe incluir la cláusula 'order by'")
	if strings.HasPrefix(strings.ToLower(limpio), "select * from ") {
		query = replaceAsterisk(query, dst)
	}
	rows, err := dbPool.Query(context.Background(), query, params...)
	errores.PanicIfError(err, "GetFirstRow: %s", limpio)
	defer rows.Close()
	if !rows.Next() {
		dbLog.Infof(limpio)
		return false
	}
	err = pgxscan.ScanRow(dst, rows)
	errores.PanicIfError(err, "GetFirstRow: %s", limpio)
	dbLog.Infof(limpio)
	return true
}

// Función de utilidad para consultas que pueden devolver varias filas.
// Panic si la query no contiene un "order by".
func GetOrderedRows(dst any, query string, params ...any) {
	limpio := reemplaza(query, params...)
	isOrdered := strings.Contains(strings.ToLower(limpio), " order by ")
	errores.PanicIfTrue(!isOrdered, "GetOrderedRows: Debe incluir la cláusula 'order by'")
	if strings.HasPrefix(strings.ToLower(limpio), "select * from ") {
		query = replaceAsterisk(query, dst)
	}
	err := pgxscan.Select(context.Background(), dbPool, dst, query, params...)
	errores.PanicIfError(err, "GetOrderedRows: %s", limpio)
	dbLog.Infof(limpio)
}

// Función auxiliar de insert y update, que parsea especial en un mapa, puede ser:
//
//	campo => solo se actualiza/inserta este campo y otros explicitamente incluidos.
//	-campo => se excluye este campo de la actualización/inserción.
//	campo=expresion => se actualiza/inserta este campo con esta expresion.
func getMapaEspecial(especiales []string) (map[string]string, bool) {
	result := map[string]string{}
	excludeAll := false
	for _, item := range especiales {
		k, v, f := strings.Cut(item, "=")
		if !f {
			if strings.HasPrefix(k, "-") {
				// -campo
				k = k[1:]
				v = "-"
			} else {
				// campo
				excludeAll = true
			}
		}
		k = strings.ToLower(k)
		_, ok := result[k]
		errores.PanicIfTrue(ok, "especial %q duplicado", k)
		isArray := false
		if f && strings.HasSuffix(k, "]") {
			n := strings.Index(k, "[")
			if n > 0 {
				k = k[:n]
				isArray = true
			}
		}
		if isArray {
			result[k] = "[]"
		} else {
			result[k] = v
		}
	}
	return result, excludeAll
}

// Función auxiliar que obtiene las ordenes para campos de tipo array, p.e. campo[3]=34
func getArrayEspecial(especiales []string, campo string) []string {
	result := []string{}
	campo += "["
	for _, item := range especiales {
		if strings.HasPrefix(strings.ToLower(item), campo) {
			result = append(result, item)
		}
	}
	return result
}

// Inserta una fila en una tabla cuyo nombre sea el del tipo de src (T_nombretabla) y que tenga una pk (id uuid).
// Especial contiene una lista de campos a excluir o incluir de la insercion:
//
//	campo => solo se inserta este campo y otros explicitamente incluidos.
//	-campo => se excluye este campo de la inserción.
//	campo=expresion => se inserta este campo con esta expresion.
//
// Por ejemplo si especial es "-inicio","final=now()","parking=null" se excluye inicio, final=hora actual y parking=nulo.
// Devuelve el id de la fila insertada.
func InsertRow(src any, especiales ...string) string {
	mapaEspecial, excludeAll := getMapaEspecial(especiales)
	valor := reflect.ValueOf(src)
	tipo := valor.Type()
	campos := reflect.VisibleFields(tipo)
	tabla := strings.TrimPrefix(strings.ToLower(tipo.Name()), "t_")
	query := "insert into " + tabla
	var c int // número de campos
	var p int // número de parámetros
	for _, campo := range campos {
		fieldName := dbscan.SnakeCaseMapper(campo.Name)
		especial, ok := mapaEspecial[fieldName]
		if fieldName == "id" || especial == "-" || (excludeAll && !ok) {
			continue
		}
		if c == 0 {
			query += " ("
		} else {
			query += ","
		}
		query += fieldName
		c++
		if especial == "" {
			p++
		}
	}
	errores.PanicIfTrue(c == 0, "InsertRow: No hay campos que insertar")
	query += ") values"
	params := make([]any, p)
	c = 0
	p = 0
	for _, campo := range campos {
		fieldName := dbscan.SnakeCaseMapper(campo.Name)
		especial, ok := mapaEspecial[fieldName]
		if fieldName == "id" || especial == "-" || (excludeAll && !ok) {
			continue
		}
		if c == 0 {
			query += " ("
		} else {
			query += ","
		}
		c++
		if especial == "" {
			params[p] = valor.FieldByName(campo.Name).Interface()
			p++
			query += "$" + strconv.Itoa(p)
		} else {
			query += especial
		}

	}
	query += ") returning id"
	limpio := reemplaza(query, params...)
	row := dbPool.QueryRow(context.Background(), query, params...)
	var result string
	err := row.Scan(&result)
	errores.PanicIfError(err, "InsertRow: %s", limpio)
	dbLog.Infof(limpio)
	return result
}

// Actualiza una fila en una tabla cuyo nombre sea el del tipo de src (T_nombretabla) y que tenga una pk (id uuid)
// Especial contiene una lista de campos a incluir o excluir de la actualización:
//
//	campo => solo se actualiza este campo y otros explicitamente incluidos.
//	-campo => se excluye este campo de la actualización.
//	campo=expresion => se actualiza este campo con esta expresion.
//
// Por ejemplo si especial es "-inicio","final=now()","parking=null" se excluye inicio, final=hora actual, parking=nulo y la tabla a actualizar es otra
// Panic si la fila no existe
func UpdateRow(src any, especiales ...string) {
	mapaEspecial, excludeAll := getMapaEspecial(especiales)
	valor := reflect.ValueOf(src)
	tipo := valor.Type()
	campos := reflect.VisibleFields(tipo)
	tabla := strings.TrimPrefix(strings.ToLower(tipo.Name()), "t_")
	query := "update " + tabla + " set "
	var c int // número de campos
	var p int // número de parámetros
	for _, campo := range campos {
		fieldName := dbscan.SnakeCaseMapper(campo.Name)
		especial, ok := mapaEspecial[fieldName]
		if fieldName == "id" || especial == "-" || (excludeAll && !ok) {
			continue
		}
		if c > 0 {
			query += ","
		}
		c++
		if especial == "" {
			p++
			query += fieldName + "=$" + strconv.Itoa(p)
		} else if especial == "[]" {
			for k, a := range getArrayEspecial(especiales, fieldName) {
				if k > 0 {
					query += ","
				}
				query += a
			}
		} else {
			query += fieldName + "=" + especial
		}
	}
	errores.PanicIfTrue(c == 0, "UpdateRow: No hay campos que actualizar")
	p++
	query += " where id=$" + strconv.Itoa(p)
	params := make([]any, p)
	var id any
	c = 0
	p = 0
	for _, campo := range campos {
		fieldName := dbscan.SnakeCaseMapper(campo.Name)
		if fieldName == "id" {
			id = valor.FieldByName(campo.Name).Interface()
			continue
		}
		especial, ok := mapaEspecial[fieldName]
		if especial == "-" || (excludeAll && !ok) {
			continue
		}
		if especial == "" {
			params[p] = valor.FieldByName(campo.Name).Interface()
			p++
		}
	}
	errores.PanicIfTrue(id == nil, "UpdateRow: Falta el campo 'id'")
	params[p] = id
	limpio := reemplaza(query, params...)
	tag, err := dbPool.Exec(context.Background(), query, params...)
	errores.PanicIfError(err, "UpdateRow: %s", limpio)
	errores.PanicIfTrue(tag.RowsAffected() == 0, "UpdateRow: Ninguna fila actualizada: %s", limpio)
	errores.PanicIfTrue(tag.RowsAffected() >= 2, "UpdateRow: %d filas actualizadas: %s", tag.RowsAffected(), limpio)
	dbLog.Infof(limpio)
}

// Elimina una fila en una tabla cuya pk sea 'id uuid'
// Panic si la fila no existe
func DeleteRow(id string, table string) {
	query := "delete from " + table + " where id=$1"
	limpio := reemplaza(query, id)
	if inTest {
		// Truco para mantener el log invariante en los tests
		limpio = reemplaza(query, "81c11fc2-0439-4ae5-baa4-3d40716bdce3")
	}
	tag, err := dbPool.Exec(context.Background(), query, id)
	errores.PanicIfError(err, "DeleteRow: %s", limpio)
	errores.PanicIfTrue(tag.RowsAffected() == 0, "DeleteRow: Ninguna fila eliminada: %s", limpio)
	errores.PanicIfTrue(tag.RowsAffected() >= 2, "DeleteRow: %d filas eliminadas: %s", tag.RowsAffected(), limpio)
	dbLog.Infof(limpio)
}

// auxiliar reemplaza()
var singleSpacePattern = regexp.MustCompile(`\s+`)

// Reemplaza parámetros y sanitiza la orden, a efectos de mostrarla en los logs
func reemplaza(query string, params ...any) string {
	query = singleSpacePattern.ReplaceAllString(strings.TrimSpace(query), " ")
	for k := len(params); k > 0; k-- {
		var valor string
		switch v := params[k-1].(type) {
		case string:
			valor = misc.EscapeSQL(v)
		case []byte:
			valor = misc.EscapeSQL(string(v))
		case pgtype.UUID:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = misc.EscapeSQL(formato.PrintUUID(v))
			default:
				valor = "undefined"
			}
		case time.Time:
			valor = misc.EscapeSQL(formato.PrintFechaHora(v, formato.ISO))
		case pgtype.Date:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = misc.EscapeSQL(formato.PrintDate(v, formato.ISO))
			default:
				valor = "undefined"
			}
		case pgtype.Timestamp:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = misc.EscapeSQL(formato.PrintTimestamp(v, formato.ISO))
			default:
				valor = "undefined"
			}
		case pgtype.Time:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = misc.EscapeSQL(formato.PrintTime(v, true))
			default:
				valor = "undefined"
			}
		case pgtype.Interval:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = misc.EscapeSQL(formato.PrintIntervalIso(v))
			default:
				valor = "undefined"
			}
		case pgtype.Text:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = misc.EscapeSQL(v.String)
			default:
				valor = "undefined"
			}
		case pgtype.Float8:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = strconv.FormatFloat(v.Float, 'g', 4, 64)
			default:
				valor = "undefined"
			}
		case pgtype.Int2:
			switch v.Status {
			case pgtype.Null:
				valor = "null"
			case pgtype.Present:
				valor = strconv.Itoa(int(v.Int))
			default:
				valor = "undefined"
			}
		default:
			valor = fmt.Sprintf("%v", v)
		}
		query = strings.ReplaceAll(query, "$"+strconv.Itoa(k), valor)
	}
	return query
}

// Si dst solo tiene estructuras anidadas (no embebidas), sustituye "select * ..." por la lista cualificada de campos (select alias.campo as "alias.campo",...)
func replaceAsterisk(query string, dst any) string {
	tipo := reflect.TypeOf(dst).Elem()
	if tipo.Kind() == reflect.Slice {
		tipo = tipo.Elem()
	}
	if tipo.Kind() != reflect.Struct {
		return query
	}
	lista := []string{}
	for _, f1 := range reflect.VisibleFields(tipo) {
		if f1.Type.Kind() != reflect.Struct || f1.Anonymous {
			return query
		}
		s := f1.Type.String()
		if !strings.HasPrefix(s, "time.") && !strings.HasPrefix(s, "pgtype.") {
			tabla := strings.ToLower(f1.Name)
			for _, f2 := range reflect.VisibleFields(f1.Type) {
				campo := dbscan.SnakeCaseMapper(f2.Name)
				lista = append(lista, fmt.Sprintf(`%s.%s as "%s.%s"`, tabla, campo, tabla, campo))
			}
			continue
		}
	}
	return strings.Replace(query, "*", strings.Join(lista, ","), 1)
}

type TipoErrorSQL int

const (
	NON_SQL                        TipoErrorSQL = 0
	SQL_OTHER                      TipoErrorSQL = 1
	INTEGRITY_CONSTRAINT_VIOLATION TipoErrorSQL = 2
	PL_PGSQL_RAISE_EXCEPTION       TipoErrorSQL = 3
)

// Determina el tipo de error SQL
// https://www.postgresql.org/docs/13/errcodes-appendix.html
func GetErrorSQL(err error) (TipoErrorSQL, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if strings.HasPrefix(pgErr.SQLState(), "23") {
			return INTEGRITY_CONSTRAINT_VIOLATION, pgErr.Detail
		}
		if pgErr.SQLState() == "P0001" {
			return PL_PGSQL_RAISE_EXCEPTION, pgErr.Message
		}
		return SQL_OTHER, ""
	}
	return NON_SQL, ""
}

// Obtiene una conexión del pool
func AcquireConnection() (conn *pgxpool.Conn, err error) {
	return dbPool.Acquire(context.Background())
}

// Devuelve una conexión al pool
func ReleaseConnection(conn *pgxpool.Conn) {
	conn.Release()
}
