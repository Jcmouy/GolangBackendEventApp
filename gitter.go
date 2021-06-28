package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	mathrand "math/rand"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dBUser     = "postgres"
	dBPassword = "root"
	dBName     = "EventApp"
)

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getPrivateContacts(id string, dbinfo string) []Usuarios_privado {
	var (
		usuario_privado  Usuarios_privado
		usuarios_privado []Usuarios_privado
	)

	db, err := sql.Open("postgres", dbinfo)

	rows, err := db.Query(`select "Destinatario", "Nombre", "Confirmacion" from "Localizacion_privado" left join "Localizacion_evento" on "Localizacion_evento"."Id" = "Localizacion_privado"."Localizacion_Id" left join "Usuarios" on "Usuarios"."Mobile"="Localizacion_privado"."Destinatario" where "Localizacion_evento"."Evento"=($1) and "Usuarios"."Status"=($2);`, id, 1)

	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {

		err = rows.Scan(&usuario_privado.Mobile, &usuario_privado.Nombre, &usuario_privado.Confirmacion)
		usuarios_privado = append(usuarios_privado, usuario_privado)

		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()

	return usuarios_privado

}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dBUser, dBPassword, dBName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	type Person struct {
		Id         int
		First_Name string
		Last_Name  string
	}
	router := gin.Default()

	// PUT - update a person details
	router.PUT("/updateperson", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Query("id")
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources
	router.DELETE("/deleteperson", func(c *gin.Context) {
		id := c.Query("id")
		stmt, err := db.Prepare("delete from person where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})

	// POST new Localizacion_evento
	router.POST("/localizacion/insert", func(c *gin.Context) {
		var buffer bytes.Buffer
		var Id string

		lat := c.PostForm("Lat")
		long := c.PostForm("Long")
		evento := c.PostForm("Evento")

		err := db.QueryRow(`insert into "Localizacion_evento" ("Lat", "Long", "Evento") values($1,$2,$3) RETURNING "Id";`, lat, long, evento).Scan(&Id)

		if err != nil {
			panic(err)
		}

		fmt.Println(Id)

		// Fastest way to append strings
		buffer.WriteString(lat)
		buffer.WriteString(" ")
		buffer.WriteString(long)
		buffer.WriteString(" ")
		buffer.WriteString(evento)

		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"ID":      fmt.Sprintf(Id),
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// GET all eventos
	router.GET("/localizaciones", func(c *gin.Context) {
		var (
			localizacion   Localizacion_evento
			localizaciones []Localizacion_evento
		)
		rows, err := db.Query(`select * from "Localizacion_evento";`)

		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&localizacion.ID, &localizacion.Lat, &localizacion.Long, &localizacion.Evento)
			localizaciones = append(localizaciones, localizacion)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()

		c.JSON(http.StatusOK, localizaciones)

	})

	// GET all eventos
	router.GET("/localizaciones/getIconos", func(c *gin.Context) {
		var (
			localizacion   Localizacion_evento_icono
			localizaciones []Localizacion_evento_icono
		)
		//rows, err := db.Query(`select "Localizacion_evento".*, "Evento"."Icono" AS "IdIcono" from "Localizacion_evento" left join "Evento" on "Evento"."Id" = "Localizacion_evento"."Evento";`)

		//Get locations by visibility
		rows, err := db.Query(`select "Localizacion_evento".*, "Evento"."Icono" AS "IdIcono", "Evento"."Nombre" AS "NombreEvento" from "Localizacion_evento" left join "Evento" on "Evento"."Id" = "Localizacion_evento"."Evento" where "Evento"."Visibilidad"='true';`)

		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&localizacion.ID, &localizacion.Lat, &localizacion.Long, &localizacion.Evento, &localizacion.IdIcono, &localizacion.NombreEvento)
			localizaciones = append(localizaciones, localizacion)

			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()

		c.JSON(http.StatusOK, localizaciones)

	})

	// POST new Localizacion_evento
	router.POST("/localizacion_privado/insert", func(c *gin.Context) {
		var buffer bytes.Buffer
		var Id string

		localizacion_id := c.PostForm("Localizacion_Id")
		remitente := c.PostForm("Remitente")
		destinatario := c.PostForm("Destinatario")

		err := db.QueryRow(`insert into "Localizacion_privado" ("Localizacion_Id", "Remitente", "Destinatario") values($1,$2,$3) RETURNING "Id";`, localizacion_id, remitente, destinatario).Scan(&Id)

		if err != nil {
			panic(err)
		}

		fmt.Println(Id)

		// Fastest way to append strings
		buffer.WriteString(localizacion_id)
		buffer.WriteString(" ")
		buffer.WriteString(remitente)
		buffer.WriteString(" ")
		buffer.WriteString(destinatario)

		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"ID":      fmt.Sprintf(Id),
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// GET all eventos
	router.GET("/localizacion_privado/getIconos/:myNumber", func(c *gin.Context) {
		var (
			localizacion   Localizacion_evento_icono
			localizaciones []Localizacion_evento_icono
		)
		myNumber := c.Param("myNumber")

		rows, err := db.Query(`select distinct "Localizacion_evento".*, "Evento"."Icono" AS "IdIcono", "Evento"."Nombre" AS "NombreEvento" from "Localizacion_evento" left join "Localizacion_privado" on "Localizacion_privado"."Localizacion_Id"="Localizacion_evento"."Id" left join "Evento" on "Evento"."Id" = "Localizacion_evento"."Evento" where "Localizacion_privado"."Remitente"=($1) or "Localizacion_privado"."Confirmacion"=($2);`, myNumber, `Y`)

		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&localizacion.ID, &localizacion.Lat, &localizacion.Long, &localizacion.Evento, &localizacion.IdIcono, &localizacion.NombreEvento)
			localizaciones = append(localizaciones, localizacion)

			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()

		c.JSON(http.StatusOK, localizaciones)

	})

	// GET all eventos
	router.GET("/localizaciones_notificacion/get/:myNumber", func(c *gin.Context) {
		var (
			localizacion   Localizacion_evento_icono
			localizaciones []Localizacion_evento_icono
		)
		myNumber := c.Param("myNumber")

		rows, err := db.Query(`select distinct "Localizacion_evento".*, "Evento"."Icono" AS "IdIcono", "Evento"."Nombre" AS "NombreEvento" from "Localizacion_evento" left join "Evento" on "Evento"."Id" = "Localizacion_evento"."Evento" left join "Localizacion_privado" on "Localizacion_privado"."Localizacion_Id"="Localizacion_evento"."Id" where "Evento"."Usuario_mobile"!=($1) and "Evento"."Visibilidad"=($2) or "Evento"."Usuario_mobile"!=($3) and "Localizacion_privado"."Destinatario"=($4) and "Localizacion_privado"."Confirmacion"=($5);`, myNumber, true, myNumber, myNumber, `Y`)

		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&localizacion.ID, &localizacion.Lat, &localizacion.Long, &localizacion.Evento, &localizacion.IdIcono, &localizacion.NombreEvento)
			localizaciones = append(localizaciones, localizacion)

			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()

		c.JSON(http.StatusOK, localizaciones)

	})

	// GET a Tipo_evento detail
	router.GET("/tipo_evento/get/:id", func(c *gin.Context) {
		var (
			tipo_evento Tipo_evento
			result      gin.H
		)
		id := c.Param("id")
		row := db.QueryRow(`select * from "Tipo_evento" where "Id"=($1);`, id)
		err = row.Scan(&tipo_evento.ID, &tipo_evento.Nombre)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": tipo_evento,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET a Evento detail
	router.GET("/evento/get/:id", func(c *gin.Context) {
		var (
			evento      Evento
			tipo_evento Tipo_evento
			result      gin.H
			//list_contacts []string
			//usuarios []Usuarios

			usuarios_privado []Usuarios_privado
		)
		id := c.Param("id")
		//row := db.QueryRow(`select * from "Evento" where "Id"=($1);`, id)

		row := db.QueryRow(`Select "Evento".*, "Tipo_evento"."Nombre" AS "NombreTipoEvento" from "Evento" left join "Tipo_evento" on "Tipo_evento"."Id" = "Evento"."Tipo_evento" where "Evento"."Id"=($1);`, id)

		err = row.Scan(&evento.ID, &evento.Nombre, &evento.Descripcion, &evento.Tipo_evento, &evento.Subtipo_evento, &evento.Visibilidad, &evento.Icono, &evento.Imagen_profile, &evento.Video_background, &evento.Usuario_mobile, &evento.Status, &tipo_evento.Nombre)

		if evento.Visibilidad == false {
			//list_contacts = getPrivateContacts(id, dbinfo)
			//usuarios = getPrivateContacts(id, dbinfo)
			usuarios_privado = getPrivateContacts(id, dbinfo)
		}

		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
			fmt.Println(err.Error())
		} else {
			result = gin.H{
				"result":           evento,
				"NombreTipoEvento": tipo_evento.Nombre,
				//"List_Contacts":    list_contacts,
				//"List_Contacts": usuarios,
				"List_Contacts": usuarios_privado,
				"count":         1,
			}
		}

		c.JSON(http.StatusOK, result)

	})

	// GET all eventos
	router.GET("/eventos", func(c *gin.Context) {
		var (
			evento  Evento
			eventos []Evento
		)
		rows, err := db.Query(`select * from "Evento";`)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&evento.ID, &evento.Nombre, &evento.Descripcion, &evento.Tipo_evento, &evento.Subtipo_evento, &evento.Visibilidad, &evento.Icono, &evento.Imagen_profile, &evento.Video_background, &evento.Usuario_mobile, &evento.Status)
			eventos = append(eventos, evento)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": eventos,
			"count":  len(eventos),
		})
	})

	// POST new evento details
	router.POST("/evento/insert", func(c *gin.Context) {
		var buffer bytes.Buffer
		var Id string

		nombre := c.PostForm("Nombre")
		descripcion := c.PostForm("Descripcion")
		tipo_evento := c.PostForm("Tipo_evento")
		subtipo_evento := c.PostForm("Subtipo_evento")
		visibilidad := c.PostForm("Visibilidad")
		icono := c.PostForm("Icono")
		imagen_profile := c.PostForm("Imagen_profile")
		video_background := c.PostForm("Video_background")
		usuario_mobile := c.PostForm("Usuario_mobile")
		status := c.PostForm("Status")

		err := db.QueryRow(`insert into "Evento" ("Nombre", "Descripcion", "Tipo_evento", "Subtipo_evento", "Visibilidad", "Icono", "Imagen_profile", "Video_background", "Usuario_mobile", "Status") values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "Id";`, nombre, descripcion, tipo_evento, subtipo_evento, visibilidad, icono, imagen_profile, video_background, usuario_mobile, status).Scan(&Id)

		if err != nil {
			panic(err)
		}

		fmt.Println(Id)

		// Fastest way to append strings
		buffer.WriteString(nombre)
		buffer.WriteString(" ")
		buffer.WriteString(descripcion)
		buffer.WriteString(" ")
		buffer.WriteString(tipo_evento)
		buffer.WriteString(" ")
		buffer.WriteString(subtipo_evento)
		buffer.WriteString(" ")
		buffer.WriteString(visibilidad)
		buffer.WriteString(" ")
		buffer.WriteString(icono)
		buffer.WriteString(" ")
		buffer.WriteString(imagen_profile)
		buffer.WriteString(" ")
		buffer.WriteString(video_background)
		buffer.WriteString(" ")
		buffer.WriteString(usuario_mobile)
		buffer.WriteString(" ")
		buffer.WriteString(status)

		//defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			//"message": fmt.Sprintf("%s", Id),
			"ID":      fmt.Sprintf(Id),
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	})

	// get all Icono
	router.GET("/icono/getAll", func(c *gin.Context) {
		var (
			icono  Icono
			iconos []Icono
		)
		rows, err := db.Query(`select * from "Icono";`)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&icono.ID, &icono.Nombre, &icono.Imagen)
			iconos = append(iconos, icono)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": iconos,
			"count":  len(iconos),
		})
	})

	// get all Icono
	router.GET("/tipo_evento/getAll", func(c *gin.Context) {
		var (
			tipo_evento  Tipo_evento
			tipo_eventos []Tipo_evento
		)

		rows, err := db.Query(`select * from "Tipo_evento";`)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&tipo_evento.ID, &tipo_evento.Nombre)
			tipo_eventos = append(tipo_eventos, tipo_evento)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		/*
			c.JSON(http.StatusOK, gin.H{
				"result": tipo_eventos,
				"count":  len(tipo_eventos),
			})
		*/
		c.JSON(http.StatusOK, tipo_eventos)
	})

	// GET check if user has an active public event
	router.GET("/usuario/get_evento/:mobile", func(c *gin.Context) {
		var buffer bytes.Buffer

		var (
			evento Evento
			result gin.H
		)
		mobile := c.Param("mobile")
		row := db.QueryRow(`select * from "Evento" where "Usuario_mobile"=($1) and "Visibilidad"=($2) and "Status"=($3);`, mobile, true, true)
		err = row.Scan(&evento.ID, &evento.Nombre, &evento.Descripcion, &evento.Tipo_evento, &evento.Subtipo_evento, &evento.Visibilidad, &evento.Icono, &evento.Imagen_profile, &evento.Video_background, &evento.Usuario_mobile, &evento.Status)
		if err != nil {
			fmt.Print(err.Error())
		}

		buffer.WriteString(mobile)

		if err != nil {
			// If no results send null
			result = gin.H{
				"result": false,
			}
		} else {
			result = gin.H{
				"result": true,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET usuario Status,Premium by mobile
	router.GET("/usuario/get_mobile/:mobile", func(c *gin.Context) {
		var buffer bytes.Buffer

		var (
			usuario Usuarios
			result  gin.H
		)
		mobile := c.Param("mobile")
		row := db.QueryRow(`select "Usuarios"."Status", "Usuarios"."Premium" from "Usuarios" left join "Evento" on "Evento"."Usuario_mobile"="Usuarios"."Mobile" where "Usuarios"."Mobile"=($1);`, mobile)
		err = row.Scan(&usuario.Status, &usuario.Premium)
		if err != nil {
			fmt.Print(err.Error())
		}

		buffer.WriteString(mobile)

		if err != nil {
			// If no results send null
			result = gin.H{
				"status":  usuario.Status,
				"premium": usuario.Premium,
				"count":   0,
			}
		} else {
			result = gin.H{
				"status":  usuario.Status,
				"premium": usuario.Premium,
				"count":   1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// POST new Usuario
	router.POST("/usuario/pedir_sms", func(c *gin.Context) {
		var mensaje string

		nombre := c.PostForm("Nombre")
		email := c.PostForm("Email")
		mobile := c.PostForm("Mobile")

		min := 100000
		max := 999999

		mathrand.Seed(time.Now().Unix())
		otp := mathrand.Intn(max-min) + min

		res := createUser(nombre, email, mobile, otp, dbinfo)

		if res == User_created_successfully {

			//fmt.Println("Llega aquí")

			send(mobile, otp)

			mensaje = "SMS request is initiated You will be receiving it shortly"

		} else if res == User_create_failed {
			mensaje = "Sorry, error in registration"

			if err != nil {
				panic(err)
			}

		} else {
			mensaje = "Mobile numer already existed!"
			if err != nil {
				panic(err)
			}

		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", mensaje),
		})
	})

	router.POST("/usuario/verificar_otp", func(c *gin.Context) {
		var mensaje string

		otp := c.PostForm("Otp")

		otp_int, err := strconv.Atoi(otp)

		if err != nil {
			panic(err)
		}

		usuario := activateUser(otp_int, dbinfo)

		if usuario.Nombre != "" {

			mensaje = "User created successfully!"
		} else {
			mensaje = "Sorry! Failed to create your account!"
		}

		fmt.Println(mensaje)

		c.JSON(http.StatusOK, usuario)

	})

	router.Run("192.168.1.11:3000")
}
