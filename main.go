package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/go-sessions"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	// "os"
)

var db *sql.DB
var err error

type M map[string]interface{}

type user struct {
	ID        int
	Username  string
	Password  string
	FirstName string
	Email     string
	Akses     string
}

type Artikels struct {
	ID       int
	Judul    string
	Isi      string
	Berkas   string
	Publish  string
	Penerbit string
}

type pesan struct {
	ID    int
	Nama  string
	Email string
	Judul string
	Isi   string
}

//@rachman: start front-end --------------------------------------------------------------------
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := QueryWallArtikelAll()
	data2 := struct{ Data []Artikels }{Data: p}

	var tmpl = template.Must(template.ParseFiles(
		"views/index.html",
		"views/menu.html",
		"views/login.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	data := struct {
		A map[string]interface{}
		B []Artikels
	}{M{"name": "MKRI", "judul": "USERS"}, data2.Data}

	var err3 = tmpl.ExecuteTemplate(w, "index", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles(
		"views/about.html",
		"views/menu.html",
		"views/login.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	var data = M{"name": "MKRI", "judul": "ABOUT"}
	var err3 = tmpl.ExecuteTemplate(w, "about", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func kontak(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles(
		"views/kontak.html",
		"views/menu.html",
		"views/login.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	var data = M{"name": "MKRI", "judul": "CONTACT US"}
	var err3 = tmpl.ExecuteTemplate(w, "kontak", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func kontakUser(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}
	var tmpl = template.Must(template.ParseFiles(
		"views/kontakUser.html",
		"views/menuUser.html",
		"views/menuAdmin.html",
		"views/login.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	var data = M{"name": "MKRI", "judul": "CONTACT US", "nama": session.GetString("name"), "mailx": session.GetString("email"), "ssid": session.GetString("akses")}
	var err3 = tmpl.ExecuteTemplate(w, "kontakUser", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func homeAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	r.ParseForm()
	p := QueryWallArtikelAll()
	data2 := struct{ Data []Artikels }{Data: p}

	var tmpl = template.Must(template.ParseFiles(
		"views/homeAdmin.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	data := struct {
		A map[string]interface{}
		B []Artikels
	}{M{"name": "Mahkamah Konstitusi Republik Indonesia", "judul": "", "ssid": session.GetString("akses")}, data2.Data}
	var err3 = tmpl.ExecuteTemplate(w, "homeAdmin", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func aboutAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}
	var tmpl = template.Must(template.ParseFiles(
		"views/aboutAdmin.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	var data = M{"name": "MKRI", "judul": "ABOUT", "ssid": session.GetString("akses")}
	var err3 = tmpl.ExecuteTemplate(w, "aboutAdmin", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func kontakAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	r.ParseForm()
	p := QueryPesanAll()
	data2 := struct{ Data []pesan }{Data: p}

	var funcMap = template.FuncMap{
		"inc": func(i, j int) int {
			return i + j
		},
	}
	var tmpl = template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"views/kontakAdmin.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	data := struct {
		A map[string]interface{}
		B []pesan
	}{M{"name": "MKRI", "judul": "CONTACT US", "ssid": session.GetString("akses")}, data2.Data}

	var err3 = tmpl.ExecuteTemplate(w, "kontakAdmin", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func artikelAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	r.ParseForm()
	p := QueryArtikelAll()
	data2 := struct{ Data []Artikels }{Data: p}

	var funcMap = template.FuncMap{
		"inc": func(i, j int) int {
			return i + j
		},
	}

	var tmpl = template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"views/artikelAdmin.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	data := Artikel{
		A: M{"name": "MKRI", "judul": "ARTICLES", "ssid": session.GetString("akses")},
		B: data2.Data,
		// C: func(i, j int) int {
		// 	return i + j
		// },
	}
	// 	A map[string]interface{}
	// 	B []Artikels
	// 	C Inc
	// }{M{"name": "MKRI", "judul": "ARTICLES", "ssid": session.GetString("akses")}, data2.Data, func(i, j int) int { return i + j }}

	var err3 = tmpl.ExecuteTemplate(w, "artikel", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

type Artikel struct {
	A map[string]interface{}
	B []Artikels
	C func(int, int) int
}

type /*(a artikel)*/ Inc func(int, int) int //{
// i = i + 1
// 	return i+j
// }
func pengguna(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	r.ParseForm()
	p := QueryUserAll()
	data2 := struct{ Data []user }{Data: p}
	//mr.MExecute(w, t, data)
	var funcMap = template.FuncMap{
		"inc": func(i, j int) int {
			return i + j
		},
	}
	var tmpl = template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"views/pengguna.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//fmt.Println(data2.Data)

	data := struct {
		A map[string]interface{}
		B []user
	}{M{"name": "MKRI", "judul": "USERS", "ssid": session.GetString("akses")}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "pengguna", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func editPengguna(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idguna := r.FormValue("id")
	r.ParseForm()
	p := Query1User(idguna)
	data2 := struct{ Data []user }{Data: p}
	//mr.MExecute(w, t, data)

	var tmpl = template.Must(template.ParseFiles(
		"views/editPengguna.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := struct {
		A map[string]interface{}
		B []user
	}{M{"name": "MKRI", "judul": "EDIT ARTICLES", "ssid": session.GetString("akses")}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "editPengguna", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func detailPesan(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idpsn := r.FormValue("id")
	r.ParseForm()
	p := Query1Pesan(idpsn)
	data2 := struct{ Data []pesan }{Data: p}
	//mr.MExecute(w, t, data)

	var tmpl = template.Must(template.ParseFiles(
		"views/detailPesan.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := struct {
		A map[string]interface{}
		B []pesan
	}{M{"name": "MKRI", "judul": "EDIT ARTICLES", "ssid": session.GetString("akses")}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "detailPesan", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func detailArtikelAdmin(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idarti := r.FormValue("idart")
	r.ParseForm()
	p := Query1Artikel(idarti)
	data2 := struct{ Data []Artikels }{Data: p}
	//mr.MExecute(w, t, data)

	var tmpl = template.Must(template.ParseFiles(
		"views/detailArtikelLogin.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/menu.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := struct {
		A map[string]interface{}
		B []Artikels
	}{M{"name": "MKRI", "judul": "DETAIL ARTICLES", "ssid": session.GetString("akses")}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "detailArtikelAdmin", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func detailArtikel(w http.ResponseWriter, r *http.Request) {

	idarti := r.FormValue("idart")
	r.ParseForm()
	p := Query1Artikel(idarti)
	data2 := struct{ Data []Artikels }{Data: p}
	//mr.MExecute(w, t, data)

	var tmpl = template.Must(template.ParseFiles(
		"views/detailArtikel.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/menu.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := struct {
		A map[string]interface{}
		B []Artikels
	}{M{"name": "MKRI", "judul": "DETAIL ARTICLES"}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "detailArtikel", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func editArtikel(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idart := r.FormValue("id")
	r.ParseForm()
	p := Query1Artikel(idart)
	data2 := struct{ Data []Artikels }{Data: p}
	//mr.MExecute(w, t, data)

	var tmpl = template.Must(template.ParseFiles(
		"views/editArtikel.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := struct {
		A map[string]interface{}
		B []Artikels
	}{M{"name": "MKRI", "judul": "EDIT ARTICLES", "ssid": session.GetString("akses")}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "editArtikel", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func editFoto(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idart := r.FormValue("dift")
	r.ParseForm()
	p := Query1Artikel(idart)
	data2 := struct{ Data []Artikels }{Data: p}
	//mr.MExecute(w, t, data)

	var tmpl = template.Must(template.ParseFiles(
		"views/editFoto.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := struct {
		A map[string]interface{}
		B []Artikels
	}{M{"name": "MKRI", "judul": "EDIT ARTICLES", "ssid": session.GetString("akses")}, data2.Data}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "editFoto", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

func tambahArtikel(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/tambahArtikel.html",
		"views/menuAdmin.html",
		"views/menuUser.html",
		"views/register.html",
		"views/footer.html",
		"views/js.html",
		"views/css.html",
	))

	//data3 := struct{Data []hks}{Data:privl}

	//fmt.Println(privl)

	data := M{"name": "MKRI", "judul": "EDIT USERS", "ssid": session.GetString("akses")}
	//var data = struct{, []user{}}
	var err3 = tmpl.ExecuteTemplate(w, "tambahArtikel", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
}

//@rachman: end front-end ------------------------------------------------------------------------------------------------

//@rachman: start query database -----------------------------------------------------------------------------------------
func connect_db() {
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/rakaweb")

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

func QueryUser(username string) user {
	var users = user{}
	err = db.QueryRow(`
		SELECT id, 
		usname,
		pswd, 
		nama, 
		email, 
		akses 
		FROM aptx WHERE usname=?
		`, username).
		Scan(
			&users.ID,
			&users.Username,
			&users.Password,
			&users.FirstName,
			&users.Email,
			&users.Akses,
		)
	return users
}

func QueryUserAll() []user { //mengambil data dari database
	rows, _ := db.Query(`SELECT aptx.id as ID,aptx.usname as Username,aptx.pswd as Password, aptx.nama as FirstName, 
						aptx.email as Email, aptx.akses as Akses FROM aptx where st = 1`)

	mis := []user{}
	mi := user{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Username, &mi.Password, &mi.FirstName, &mi.Email, &mi.Akses)
		mis = append(mis, mi)
	}
	return mis
}

func QueryPesanAll() []pesan { //mengambil data dari database
	rows, _ := db.Query(`SELECT id as ID, nama as Nama,email as Email, judul as Judul, SUBSTRING(isi, 1, 200) as Isi FROM pesan order by id desc`)

	mis := []pesan{}
	mi := pesan{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Nama, &mi.Email, &mi.Judul, &mi.Isi)
		mis = append(mis, mi)
	}
	return mis
}

func QueryArtikelAll() []Artikels { //mengambil data dari database
	rows, _ := db.Query(`SELECT id as ID, judul_artikel as Judul,SUBSTRING(isi_artikel, 1, 800) as Isi, 
						berkas as Berkas, publish as Publish, 
						penerbit as Penerbit FROM artikel order by id desc`)

	mis := []Artikels{}
	mi := Artikels{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Judul, &mi.Isi, &mi.Berkas, &mi.Publish, &mi.Penerbit)
		mis = append(mis, mi)
	}
	return mis
}

func QueryWallArtikelAll() []Artikels { //mengambil data dari database
	rows, _ := db.Query(`SELECT id as ID, judul_artikel as Judul,SUBSTRING(isi_artikel, 1, 200) as Isi, 
						berkas as Berkas, publish as Publish, 
						penerbit as Penerbit FROM artikel WHERE publish = 1 order by id desc`)

	mis := []Artikels{}
	mi := Artikels{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Judul, &mi.Isi, &mi.Berkas, &mi.Publish, &mi.Penerbit)
		mis = append(mis, mi)
	}
	return mis
}

func Query1Artikel(id string) []Artikels {
	rows, _ := db.Query(`SELECT id as ID, judul_artikel as Judul, isi_artikel as Isi, berkas as Berkas, publish as Publish, 
						  penerbit as Penerbit FROM artikel WHERE id =` + id)
	mis := []Artikels{}
	mi := Artikels{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Judul, &mi.Isi, &mi.Berkas, &mi.Publish, &mi.Penerbit)
		mis = append(mis, mi)
	}
	return mis
}

func Query1Pesan(id string) []pesan {
	rows, _ := db.Query(`SELECT id as ID, nama as Nama, email as Email, judul as Judul, isi as ISi 
						 FROM pesan WHERE id =` + id)
	mis := []pesan{}
	mi := pesan{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Nama, &mi.Email, &mi.Judul, &mi.Isi)
		mis = append(mis, mi)
	}
	return mis
}

func Query1User(id string) []user { //mengambil data dari database
	rows, _ := db.Query(`SELECT id as ID,usname as Username,pswd as Password, nama as FirstName, 
		email as Email, akses as Akses FROM aptx WHERE id =` + id)

	mis := []user{}
	mi := user{}

	for rows.Next() {
		rows.Scan(&mi.ID, &mi.Username, &mi.Password, &mi.FirstName, &mi.Email, &mi.Akses)
		mis = append(mis, mi)
	}
	return mis
}

//@rachman: end query database -----------------------------------------------------------------------------------------

//@rachman: start function ---------------------------------------------------------------------------------------------
func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}

	return true
}

func login(w http.ResponseWriter, r *http.Request) {
	/*
		session := sessions.Start(w, r)
		if len(session.GetString("username")) != 0 && checkErr(w, r, err) {
			http.Redirect(w, r, "/", 302)
		}
	*/
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	username := r.FormValue("uname")
	password := r.FormValue("pwd")

	users := QueryUser(username)

	//deskripsi dan compare password
	var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

	if password_tes == nil {
		//login success
		session := sessions.Start(w, r)
		session.Set("username", users.Username)
		session.Set("name", users.FirstName)
		session.Set("email", users.Email)
		session.Set("akses", users.Akses)
		http.Redirect(w, r, "/homeAdmin", 302)
	} else {
		//login failed
		http.Redirect(w, r, "/", 302)
	}

}

func simpanEditArtikel(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	judul := r.FormValue("jdl")
	isi := r.FormValue("isi")
	publ := r.FormValue("pub")
	idart := r.FormValue("id")

	stmt, err := db.Prepare(`UPDATE artikel SET judul_artikel=?, isi_artikel=?, publish=? WHERE id=?`)
	if err == nil {
		_, err := stmt.Exec(&judul, &isi, &publ, &idart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/artikelAdmin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/artikelAdmin", 302)
	}

}

func simpanPengguna(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	first_name := r.FormValue("ftname")
	email := r.FormValue("email")
	privi := r.FormValue("priv")
	idguna := r.FormValue("id")

	stmt, err := db.Prepare("UPDATE aptx SET nama=?, email=?, akses=? WHERE id=?")
	if err == nil {
		_, err := stmt.Exec(&first_name, &email, &privi, &idguna)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/pengguna", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/pengguna", 302)
	}
}

func deletePengguna(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idguna := r.FormValue("iddl")

	stmt, err := db.Prepare("DELETE FROM aptx WHERE id=?")
	if err == nil {
		_, err := stmt.Exec(&idguna)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/pengguna", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/pengguna", 302)
	}
}

func deleteArtikel(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	idart := r.FormValue("idrt")

	stmt, err := db.Prepare("DELETE FROM artikel WHERE id=?")
	if err == nil {
		_, err := stmt.Exec(&idart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/artikelAdmin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/artikelAdmin", 302)
	}
}

func simpanTambahArtikel(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := xid.New()
	uploadedFile, handler, err := r.FormFile("flgm")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	//if alias != "" {
	filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	//}

	fileLocation := filepath.Join(dir, "/static/files", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.Write([]byte("done"))
	var now = time.Now()
	//layoutFormat := "2006-01-02 15:04:05"
	judul := r.FormValue("jdlart")
	isi := r.FormValue("isi")
	publis := r.FormValue("pub")
	pnrbt := session.GetString("name")
	usr := session.GetString("username")
	//date, _ := now.Parse(layoutFormat, now)

	wkt := now

	stmt, err := db.Prepare(`INSERT INTO artikel SET judul_artikel=?, isi_artikel=?, 
					berkas=?, publish=?, penerbit=?, user=?, tanggal=?`)
	if err == nil {
		_, err := stmt.Exec(&judul, &isi, &filename, &publis, &pnrbt, &usr, &wkt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/artikelAdmin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/artikelAdmin", 302)
	}

}

func simpanFotoArtikel(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := xid.New()
	uploadedFile, handler, err := r.FormFile("flgm")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	//if alias != "" {
	filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	//}

	fileLocation := filepath.Join(dir, "/static/files", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.Write([]byte("done"))
	idtkl := r.FormValue("idftrt")
	stmt, err := db.Prepare("UPDATE artikel SET berkas=? WHERE id=?")
	if err == nil {
		_, err := stmt.Exec(&filename, &idtkl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/artikelAdmin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/artikelAdmin", 302)
	}

}

func simpanTambahPesanUser(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	namax := session.GetString("name")
	emailx := session.GetString("email")
	judulx := r.FormValue("jdl")
	isix := r.FormValue("isi")

	stmt, err := db.Prepare("INSERT INTO pesan SET nama=?, email=?, judul=?, isi=?")
	if err == nil {
		_, err := stmt.Exec(&namax, &emailx, &judulx, &isix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/homeAdmin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/homeAdmin", 302)
	}
}

func simpanTambahPesan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	//return false

	namax := r.FormValue("nm")
	emailx := r.FormValue("mail")
	judulx := r.FormValue("jdl")
	isix := r.FormValue("isi")

	stmt, err := db.Prepare("INSERT INTO pesan SET nama=?, email=?, judul=?, isi=?")
	if err == nil {
		_, err := stmt.Exec(&namax, &emailx, &judulx, &isix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	//return false

	username := r.FormValue("uname")
	first_name := r.FormValue("nama")
	email := r.FormValue("email")
	password := r.FormValue("pwd")
	privi := r.FormValue("priv")

	users := QueryUser(username)

	if (user{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && checkErr(w, r, err) {
			stmt, err := db.Prepare("INSERT INTO aptx SET usname=?, pswd=?, nama=?, email=?, akses=?")
			if err == nil {
				_, err := stmt.Exec(&username, &hashedPassword, &first_name, &email, &privi)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, "/pengguna", http.StatusSeeOther)
				return
			}
		}
	} else {
		http.Redirect(w, r, "/pengguna", 302)
	}

}

func routes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/homeAdmin", homeAdmin)
	http.HandleFunc("/pengguna", pengguna)
	http.HandleFunc("/editPengguna", editPengguna)
	http.HandleFunc("/editArtikel", editArtikel)
	http.HandleFunc("/aboutAdmin", aboutAdmin)
	http.HandleFunc("/kontakAdmin", kontakAdmin)
	http.HandleFunc("/artikelAdmin", artikelAdmin)
	http.HandleFunc("/about", about)
	http.HandleFunc("/simpanPengguna", simpanPengguna)
	http.HandleFunc("/simpanEditArtikel", simpanEditArtikel)
	http.HandleFunc("/deletePengguna", deletePengguna)
	http.HandleFunc("/deleteArtikel", deleteArtikel)
	http.HandleFunc("/kontak", kontak)
	http.HandleFunc("/editFoto", editFoto)
	http.HandleFunc("/detailPesan", detailPesan)
	http.HandleFunc("/kontakUser", kontakUser)
	http.HandleFunc("/simpanTambahPesan", simpanTambahPesan)
	http.HandleFunc("/simpanFotoArtikel", simpanFotoArtikel)
	http.HandleFunc("/tambahArtikel", tambahArtikel)
	http.HandleFunc("/detailArtikel", detailArtikel)
	http.HandleFunc("/detailArtikelAdmin", detailArtikelAdmin)
	http.HandleFunc("/simpanTambahPesanUser", simpanTambahPesanUser)
	http.HandleFunc("/simpanTambahArtikel", simpanTambahArtikel)
	http.HandleFunc("/logout", logout)
}

//@rachman: end function ---------------------------------------------------------------------------------------------

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//fsimg := http.FileServer(http.Dir("./img"))
	//http.Handle("/img/", http.StripPrefix("/img", fsimg))

	connect_db()
	routes()
	defer db.Close()

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("Error running service: ", err)
	}

}
