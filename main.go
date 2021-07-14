package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

//Modulo para conectarse a MYSQL
//go get -u github.com/go-sql-driver/mysql
//Cargar con un guion, cuando se trata de un driver
//Fin

//CONEXION BASE DE DATOS
func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

//Sirve para buscar plantillas en una carpeta determinada
var plantillas = template.Must(template.ParseGlob("Public/*"))

//Estutira empleado
type Empleados struct {
	ID     int
	Nombre string
	Correo string
}

func main() {
	//Enviar ruta y funcion a resivir el sevidor
	http.HandleFunc("/", Index)
	http.HandleFunc("/Form", Form)
	http.HandleFunc("/Insertar", Insertar)
	http.HandleFunc("/Borrar", Borrar)
	http.HandleFunc("/Editar", Editar)
	http.HandleFunc("/Actualizar", Actualizar)
	//Log mensaje por consola
	log.Println("Servidor corriendo")
	//Iniciando el seevirdo con la funcion ListenAnd Serve
	http.ListenAndServe(":3000", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	//Enviando repuesta
	//fmt.Fprint(w, "Hola Develoteca")
	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	empleado := Empleados{}
	arregloEmpleado := []Empleados{}

	for registros.Next() {
		var ID int
		var nombre, correo string
		err = registros.Scan(&ID, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = ID
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	if err != nil {
		panic(err.Error())
	}
	plantillas.ExecuteTemplate(w, "index", arregloEmpleado)
	plantillas.ExecuteTemplate(w, "Form", nil)
}

func Form(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "Form", nil)

}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("name")
		gmail := r.FormValue("email")
		conexionEstablecida := conexionBD()
		//Los signo de interrogacion sirve para cambiar el texto
		insertarRegistro, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		//Exec ejecuta la sentencia de insercion
		insertarRegistro.Exec(nombre, gmail)
		http.Redirect(w, r, "/", 301)

	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("ID")

	conexionEstablecida := conexionBD()
	//Los signo de interrogacion sirve para cambiar el texto
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE ID=?")
	if err != nil {
		panic(err.Error())
	}
	//Exec ejecuta la sentencia de insercion
	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("ID")
	conexionEstablecida := conexionBD()
	query, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE ID=?", idEmpleado)

	empleado := Empleados{}
	//arregloEmpleado := []Empleados{}

	for query.Next() {
		var ID int
		var nombre, correo string
		err = query.Scan(&ID, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = ID
		empleado.Nombre = nombre
		empleado.Correo = correo

	}
	//Enviar datos al frontend
	plantillas.ExecuteTemplate(w, "Editar", empleado)

}

func Actualizar(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		idEmpleado := r.FormValue("ID")
		nombre := r.FormValue("name")
		gmail := r.FormValue("email")
		conexionEstablecida := conexionBD()
		//Los signo de interrogacion sirve para cambiar el texto
		actualizarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET Nombre=?, Correo=? WHERE ID=?")
		if err != nil {
			panic(err.Error())
		}
		//Exec ejecuta la sentencia de insercion
		actualizarRegistro.Exec(nombre, gmail, idEmpleado)
		http.Redirect(w, r, "/", 301)

	}
}
