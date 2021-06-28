package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	User_created_successfully = 0
	User_create_failed        = 1
	User_alredy_existed       = 2
)

func createOtp(usuario_id string, otp int, dbinfo string) {
	db, err := sql.Open("postgres", dbinfo)

	//Delete the old otp if exist
	stmt, err := db.Prepare(`delete from "Sms_codigos" where "Usuario_Id"=($1);`)
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(usuario_id)

	stmt, err = db.Prepare(`insert into "Sms_codigos" ("Usuario_Id", "Codigo", "Status") values($1,$2,$3);`)
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(usuario_id, otp, 0)

}

func isUserExists(mobile string, dbinfo string) bool {
	var count int

	db, err := sql.Open("postgres", dbinfo)

	rows, err := db.Query(`select "Id" from "Usuarios" where "Mobile"=($1) and "Status"=($2);`, mobile, 1)
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		count = count + 1
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()

	if count > 0 {
		return true
	} else {
		return false
	}

}

func activateUserStatus(user_Id string, dbinfo string) {

	db, err := sql.Open("postgres", dbinfo)

	stmt, err := db.Prepare(`update "Usuarios" set "Status"=($1) where "Id"=($2);`)
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(1, user_Id)

	stmt, err = db.Prepare(`update "Sms_codigos" set "Status"=($1) where "Usuario_Id"=($2);`)
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(1, user_Id)

}

func activateUser(otp int, dbinfo string) Usuarios {
	var Id string
	var usuario Usuarios

	db, err := sql.Open("postgres", dbinfo)

	err = db.QueryRow(`SELECT "u"."Id", "u"."Nombre", "u"."Email", "u"."Mobile", "u"."Apikey", "u"."Status", "u"."Fecha_creado" FROM "Usuarios" as u left join "Sms_codigos" on "Sms_codigos"."Usuario_Id" = "u"."Id" WHERE "Sms_codigos"."Codigo"=($1);`, otp).Scan(&usuario.ID, &usuario.Nombre, &usuario.Email, &usuario.Mobile, &usuario.Api_key, &usuario.Status, &usuario.Fecha_creado)

	if err != nil {
		panic(err)
	}

	Id = strconv.Itoa(usuario.ID)

	activateUserStatus(Id, dbinfo)

	return usuario
}

func createUser(nombre string, email string, mobile string, otp int, dbinfo string) int {
	var api_key string
	var Id string
	var new_user_id string

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Print(err.Error())
	}

	test := isUserExists(mobile, dbinfo)

	if !test {

		api_key = randToken()

		err := db.QueryRow(`insert into "Usuarios" ("Nombre", "Email", "Mobile", "Apikey", "Status") values($1,$2,$3,$4,$5) RETURNING "Id";`, nombre, email, mobile, api_key, 0).Scan(&Id)

		if err != nil {
			fmt.Println("User_create_failed")
			return User_create_failed
		}

		new_user_id = Id

		createOtp(new_user_id, otp, dbinfo)

		fmt.Println(Id)

		fmt.Println("User_created_successfully")
		return User_created_successfully

	} else {
		fmt.Println("User_alredy")
		return User_alredy_existed
	}

}
