package main

import (
	"database/sql"
	"time"
)

//Test object
type Test struct {
	ID     int
	Nombre string
}

//TipoEvento object
type Tipo_evento struct {
	ID     int
	Nombre string
}

//Icono object
type Icono struct {
	ID     int
	Nombre string
	Imagen []uint8
}

//Evento object
type Evento struct {
	ID               int
	Nombre           string
	Descripcion      string
	Tipo_evento      uint64
	Subtipo_evento   string
	Visibilidad      bool
	Icono            int
	Imagen_profile   string
	Video_background string
	Usuario_mobile   string
	Status           bool
}

//Localizacion_evento object
type Localizacion_evento struct {
	ID     int
	Lat    float32
	Long   float32
	Evento uint64
}

type Localizacion_evento_icono struct {
	ID           int
	Lat          float32
	Long         float32
	Evento       uint64
	IdIcono      int
	NombreEvento string
}

type Localizacion_privado struct {
	ID              int
	Localizacion_Id int
	Remitente       string
	Destinatario    string
	Confirmacion    bool
}

type Usuarios struct {
	ID           int
	Nombre       string
	Email        string
	Mobile       string
	Api_key      string
	Status       bool
	Fecha_creado sql.NullString
	//Fecha_creado  time.Time
	Hash_password string
	Premium       sql.NullBool
}

type Usuarios_privado struct {
	Mobile       string
	Nombre       string
	Confirmacion sql.NullString
}

type Sms_codigos struct {
	ID           int
	Usuario_id   uint64
	Codigo       string
	Status       bool
	Fecha_creado time.Time
}
