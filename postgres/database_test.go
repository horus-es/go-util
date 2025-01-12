package postgres

import (
	"strings"
	"testing"
	"time"

	"github.com/horus-es/go-util/v2/errores"
	"github.com/horus-es/go-util/v2/formato"
	"github.com/horus-es/go-util/v2/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

type T_personal struct {
	ID            string
	Operador      pgtype.UUID
	Codigo        string
	Nombre        string
	Hash          pgtype.Text
	Activo        bool
	Administrador bool
}

func init() {
	InitPool(`
	    host=desarrollo.horus.es
		port=5433
		user=jesus.san2
		password=she1amee3Yus
		dbname=SPARK2
		sslmode=disable`, nil)
	inTest = true
}

var (
	UUIDempleado = "fe90b961-0646-4f8e-a698-d3a153abf7d2"
	UUIDoperador = "0cec7694-eb8d-4ab2-95bb-d5d733a3be94"
	UUIDnoexiste = "fe90b951-9999-9999-9999-999999999999"
)

func ExampleGetOneRow() {
	u := struct {
		Codigo string
		Activo bool
	}{}
	GetOneRow(&u, "select codigo,activo from personal where id=$1", UUIDempleado)
	logger.Infof("%v", u)
	// Output:
	// INFO: select codigo,activo from personal where id='fe90b961-0646-4f8e-a698-d3a153abf7d2'
	// INFO: {pablo7 true}
}

func TestGetOneRowSummary(t *testing.T) {
	var n int
	GetOneRow(&n, "select count(*) from personal where id=$1", UUIDempleado)
	assert.Equal(t, n, 1, "Carga sumaria incorrecta")
}

func TestGetOneRowPanicNone(t *testing.T) {
	u := T_personal{}
	defer func() { recover() }()
	GetOneRow(&u, "select * from personal where id=$1", UUIDnoexiste)
	t.Error("Sin pánico sin filas")
}

func TestGetOneRowPanicMany(t *testing.T) {
	p := T_personal{}
	defer func() { recover() }()
	GetOneRow(&p, "select * from personal")
	t.Error("Sin pánico con varias filas")
}

func ExampleGetOneOrZeroRows() {
	u := struct {
		Codigo string
		Activo bool
	}{}
	if GetOneOrZeroRows(&u, "select codigo,activo from personal where id=$1", UUIDempleado) {
		logger.Infof("Hallado usuario: %q\n", u.Codigo)
	}
	if !GetOneOrZeroRows(&u, "select codigo,activo from personal where id=$1", UUIDnoexiste) {
		logger.Infof("Usuario no hallado.\n")
	}
	// Output:
	// INFO: select codigo,activo from personal where id='fe90b961-0646-4f8e-a698-d3a153abf7d2'
	// INFO: Hallado usuario: "pablo7"
	// INFO: select codigo,activo from personal where id='fe90b951-9999-9999-9999-999999999999'
	// INFO: Usuario no hallado.
}

func TestGetOneOrZeroRowsPanicMany(t *testing.T) {
	p := T_personal{}
	defer func() { recover() }()
	GetOneOrZeroRows(&p, "select * from personal")
	t.Error("Sin pánico con varias filas")
}

func ExampleGetOrderedRows() {
	us := []string{}
	GetOrderedRows(&us, "select codigo from personal where operador=$1 and codigo>'dad' order by codigo limit 3", UUIDoperador)
	logger.Infof("Primeros 3 usuarios hallados: %s\n", strings.Join(us, ", "))
	// Output:
	// INFO: select codigo from personal where operador='0cec7694-eb8d-4ab2-95bb-d5d733a3be94' and codigo>'dad' order by codigo limit 3
	// INFO: Primeros 3 usuarios hallados: dadiz, emple, emple100E
}

func TestGetJoin(t *testing.T) {
	type t_operador struct {
		Id     string
		Razon  string
		Idioma string
	}
	type t_datos struct {
		Personal T_personal
		Operador t_operador
	}
	var datos []t_datos
	GetOrderedRows(&datos, "select * from personal,operadores operador where personal.operador=operador.id order by personal.id")
	assert.Greater(t, len(datos), 1, "Filas no cargadas")
}

func TestGetOrderedRowsPanic(t *testing.T) {
	var ps []*T_personal
	defer func() { recover() }()
	GetOrderedRows(&ps, "select * from personal where operador=$1", UUIDoperador)
	t.Error("Sin pánico sin 'order by'")
}

func TestInsertUpdateDelete(t *testing.T) {
	p1 := T_personal{}
	p1.Operador = formato.MustParseUUID(UUIDoperador)
	p1.Nombre = "InsertRow"
	p1.Codigo = "TestInsertUpdateDelete " + time.Now().Format("01-02-2006 15:04:05")
	p1.Hash.Valid = true
	p1.ID = InsertRow(p1)
	p2 := T_personal{}
	GetOneRow(&p2, "select * from personal where id=$1", p1.ID)
	assert.Equal(t, p1, p2, "Insert falló")
	p1.Nombre = "UpdateRow"
	p1.Activo = true
	UpdateRow(p1)
	GetOneRow(&p2, "select * from personal where id=$1", p1.ID)
	assert.Equal(t, p1, p2, "Update falló")
	DeleteRow(p1.ID, "personal")
	f := GetOneOrZeroRows(&p2, "select * from personal where id=$1", p1.ID)
	assert.False(t, f, "Delete falló")
}

func ExampleInsertRow() {
	u := T_personal{}
	u.Operador, _ = formato.ParseUUID(UUIDoperador)
	u.Codigo = "TestInsert"
	u.Nombre = "Usuario de prueba"
	u.Activo = true
	u.ID = InsertRow(u, "hash='zecreto2023'")
	logger.Infof("Una fila insertada")
	DeleteRow(u.ID, "personal")
	logger.Infof("Una fila eliminada")
	// Output:
	// INFO: insert into personal (operador,codigo,nombre,hash,activo,administrador) values ('0cec7694-eb8d-4ab2-95bb-d5d733a3be94','TestInsert','Usuario de prueba','zecreto2023',true,false) returning id
	// INFO: Una fila insertada
	// INFO: delete from personal where id='81c11fc2-0439-4ae5-baa4-3d40716bdce3'
	// INFO: Una fila eliminada
}

func ExampleDeleteRow() {
	u := T_personal{}
	u.Operador, _ = formato.ParseUUID(UUIDoperador)
	u.Codigo = "TestDelete"
	u.Nombre = "Usuario de prueba"
	u.Activo = true
	u.ID = InsertRow(u, "hash='zecreto2023'")
	logger.Infof("Una fila insertada")
	DeleteRow(u.ID, "personal")
	logger.Infof("Una fila eliminada")
	// Output:
	// INFO: insert into personal (operador,codigo,nombre,hash,activo,administrador) values ('0cec7694-eb8d-4ab2-95bb-d5d733a3be94','TestDelete','Usuario de prueba','zecreto2023',true,false) returning id
	// INFO: Una fila insertada
	// INFO: delete from personal where id='81c11fc2-0439-4ae5-baa4-3d40716bdce3'
	// INFO: Una fila eliminada
}

func ExampleUpdateRow() {
	u := T_personal{}
	GetOneRow(&u, "select * from personal where id=$1", UUIDempleado)
	logger.Infof("Usuario cargado")
	u.Codigo = "pablo7"
	UpdateRow(u, "codigo")
	logger.Infof("Nombre actualizado")
	// Output:
	// INFO: select * from personal where id='fe90b961-0646-4f8e-a698-d3a153abf7d2'
	// INFO: Usuario cargado
	// INFO: update personal set codigo='pablo7' where id='fe90b961-0646-4f8e-a698-d3a153abf7d2'
	// INFO: Nombre actualizado
}

func TestInsertUpdateDeleteEspecial(t *testing.T) {
	p1 := T_personal{}
	p1.Operador = formato.MustParseUUID(UUIDoperador)
	p1.Nombre = "InsertRow"
	p1.Codigo = "TestInsertUpdateDeleteExclude " + time.Now().Format("01-02-2006 15:04:05")
	p1.Activo = true
	p1.ID = InsertRow(p1, "-hash", "activo=false")
	p1.Activo = false
	p1.Hash.Valid = true
	p2 := T_personal{}
	GetOneRow(&p2, "select * from personal where id=$1", p1.ID)
	assert.Equal(t, p1, p2, "Insert falló")
	p1.Nombre = "UpdateRow"
	UpdateRow(p1, "-codigo", "-hash", "-operador", "activo=true")
	GetOneRow(&p2, "select * from personal where id=$1", p1.ID)
	p1.Activo = true
	assert.Equal(t, p1, p2, "Update falló")
	UpdateRow(p1, "activo=false")
	GetOneRow(&p2, "select * from personal where id=$1", p1.ID)
	p1.Activo = false
	assert.Equal(t, p1, p2, "Update falló")
	p1.Activo = true
	UpdateRow(p1, "activo")
	GetOneRow(&p2, "select * from personal where id=$1", p1.ID)
	assert.Equal(t, p1, p2, "Update falló")
	DeleteRow(p1.ID, "personal")
	f := GetOneOrZeroRows(&p2, "select * from personal where id=$1", p1.ID)
	assert.False(t, f, "Delete falló")
}

func TestUpdateNonExistant(t *testing.T) {
	p1 := T_personal{}
	p1.ID = UUIDnoexiste
	p1.Operador = formato.MustParseUUID(UUIDoperador)
	p1.Nombre = "UpdateRow"
	p1.Codigo = "TestUpdateNoExistant " + time.Now().Format("01-02-2006 15:04:05")
	defer func() { recover() }()
	UpdateRow(p1)
	t.Error("Sin pánico no existe")
}

func TestDeleteNonExistant(t *testing.T) {
	defer func() { recover() }()
	DeleteRow(UUIDnoexiste, "personal")
	t.Error("Sin pánico no existe")
}

func ExampleStartTX() {
	StartTX()
	defer RollbackTX()
	// Inicio del bloque protegido con transacción
	logger.Infof("... órdenes SQL contenidas en la transacción ...")
	// Fin del bloque protegido con transacción
	CommitTX()
	// Output:
	// INFO: StartTX
	// INFO: ... órdenes SQL contenidas en la transacción ...
	// INFO: CommitTX
}

func ExampleCommitTX() {
	StartTX()
	defer RollbackTX()
	// Inicio del bloque protegido con transacción
	logger.Infof("... órdenes SQL contenidas en la transacción ...")
	// Fin del bloque protegido con transacción
	CommitTX()
	// Output:
	// INFO: StartTX
	// INFO: ... órdenes SQL contenidas en la transacción ...
	// INFO: CommitTX
}

func ExampleRollbackTX() {
	defer func() { recover() }() // Capturamos panic
	StartTX()
	defer RollbackTX()
	// Inicio del bloque protegido con transacción
	logger.Infof("... órdenes SQL contenidas en la transacción ...")
	errores.PanicIfTrue(true, "... algo produce un panic ...")
	logger.Infof("... mas órdenes SQL ...")
	// Fin del bloque protegido con transacción
	CommitTX()
	// Output:
	// INFO: StartTX
	// INFO: ... órdenes SQL contenidas en la transacción ...
	// WARN: RollbackTX
}

func TestTX(t *testing.T) {
	// Iniciamos transacción
	StartTX()
	// Cargamos un usuario
	u := T_personal{}
	GetOneRow(&u, "select * from personal where id=$1", UUIDempleado)
	códigoOriginal := u.Codigo
	// Actualizamos la fila
	u.Codigo = "Nombre " + formato.PrintFechaHora(time.Now(), formato.ISO)
	UpdateRow(u, "codigo")
	// Deshacemos transacción
	RollbackTX()
	// Comprobamos si se ha modificado la fila
	GetOneRow(&u, "select * from personal where id=$1", UUIDempleado)
	assert.Equal(t, códigoOriginal, u.Codigo)
}
