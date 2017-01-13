// Пример построения с6енрвиса для базы данных 
// MS Sql

package main

import (
	  "database/sql"
	  "fmt"
	  "net/http"
  	  "time"
	  "io/ioutil"
	  "strings"
	  "github.com/gorilla/mux"
	_ "github.com/denisenkom/go-mssqldb"      // драйвер MS SQL
)

var db *sql.DB


// 
func main() {
	database, err := sql.Open("mssql","server=KAF; database=Sam;user id=Art;password=123")
	
	if err != nil {
	   panic(err)
	}
	
	db = database
	
	r:= mux.NewRouter()

	r.HandleFunc("/",                                         simpleHandler)
	r.HandleFunc("/insert",                                   dbInsert)
    r.HandleFunc("/delete/{id:[0-9]+}",                       dbDelete)
    r.HandleFunc("/deleteall",                                dbDeleteAll)

    r.HandleFunc("/insert/{id}/{ses}*{ot}*{nn}:{ff}-{ss}",    dbInsertpram)

	r.HandleFunc("/derek",                                    derekHandler)
	r.HandleFunc("/file",                                     ioHandler)
	r.HandleFunc("/write",                                    ioWriteHandler)

	r.HandleFunc("/db/User/{uid:[0-9]+}",                     dbHandler).Name("User")
	r.HandleFunc("/Session/{SessionID}/{uid:[0-9]+}",         sessionHandler).Name("Session")
	
    fmt.Println("Start Service Compute ....")
	http.ListenAndServe(":8082", r)
}


// Стартовая ттраница
func simpleHandler(w http.ResponseWriter, r *http.Request) {
	 fmt.Fprintf(w, "Start Page.")
}


func derekHandler(w http.ResponseWriter, r *http.Request) {
	 fmt.Fprintf(w, "Sup")
	 fmt.Fprintf(w, time.Now().Local().Format("15:04:05"))
}


// Чтение файла
func ioHandler(w http.ResponseWriter, r *http.Request) {
	thingie, err := ioutil.ReadFile("Bas.txt")
	if err != nil{
		panic(err)
	}
	
	fmt.Fprintf(w, string(thingie))
}


// Замена строки в файлике 
// wwww -> www1
func ioWriteHandler(w http.ResponseWriter, r *http.Request) {
	thingie, err := ioutil.ReadFile("Bas.txt")
	
	if err != nil{
	   panic(err)
	}
	
	mystuff := strings.Replace(string(thingie), "www", "Замена строки с поиском", 1)

	
	err = ioutil.WriteFile("Bas.txt", []byte(mystuff), 0644)
	if err != nil{
	   panic(err)
	}
	
	fmt.Fprintf(w, string(mystuff))
}

// Cессия
func sessionHandler(w http.ResponseWriter, r *http.Request) {
	params    := mux.Vars(r)
	sessionID := params["SessionID"]
	uid       := params["uid"]
	fmt.Fprintln(w, sessionID)
	fmt.Fprintf(w, string(uid))
}


// Время юникстайим
func CTU() int64 {
	 return time.Now().UnixNano() / 1000000
}

// Пример вставки записи в базу данных
func dbInsert(w http.ResponseWriter, r *http.Request) {
  
    var ii int 
    ii = 127  

    rows, err := db.Query("INSERT INTO City(Zarp, Fam, Name) VALUES (?,?,?)", ii, "Namess", "Nbamm")

	if err != nil {
	   panic(err)
	}else{
	   fmt.Println("Вставка ОК")
	}
    defer rows.Close()
}



// 
// Пример вставки записи в базу данных по шаблону 
// шаблон выбран для примера и показа изощренных
// методов работы  в строке поиска
// http://localhost:8082/insert/фамилия/имя*отчество*примечание:зарплата-123
//                       insert/{id}/{ses}*{ot}*{nn}:{ff}-{ss}
func dbInsertpram(w http.ResponseWriter, r *http.Request) {
  
    // Параметеры котрорые ловим на входе
	params    := mux.Vars(r)
	se        := params["ses"]
	uid       := params["id"]
	ot        := params["ot"]
	nn        := params["nn"]
	ff        := params["ff"]
    ss        := params["ss"]
	ii        := 123.45
    tm        := time.Now().Local().Format("15:04:05")
    dt        := time.Now().Local().Format("2006-01-02")

	fmt.Println(se,uid,ii, ff)

    rows, err := db.Query("INSERT INTO City(Zarp, Fam, Name, Otch, Address, Summ, Note, Times, Dats) VALUES (?,?,?,?,?,?,?,?,?)", ii, se, uid, ot, nn, ss, ff, tm, dt)

	if err != nil {
	   panic(err)
	}else{
	   fmt.Println("Вставка ОК")
	}
  
    defer rows.Close()
}


// Удаление всех записей
// Без предупреждения
func dbDeleteAll(w http.ResponseWriter, r *http.Request) {
    rows, _ :=db.Query("DELETE FROM City")
	fmt.Println("Произошло удаление ALL записи ")
    defer rows.Close()
}


// Пример вставки записи в базу данных
func dbDelete(w http.ResponseWriter, r *http.Request) {
  
	params := mux.Vars(r)
	id     := params["id"]
    
    rows, err := db.Query("DELETE FROM City WHERE Id=?", id)

	if err != nil {
	   panic(err)
	}else{
		fmt.Println("Произошло удаление записи ", id)
	}
  
    defer rows.Close()
}


// Поиск и высветка данных из базы
func dbHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id     := params["uid"]
	
	// rows, err := db.Query("SELECT AppUserID, Handle FROM AppUser where AppUserID = ? order by AppUserID desc", id)

	fmt.Println(id)
	rows, err := db.Query("SELECT Fam, Name FROM City WHERE Id=? ORDER BY Fam", id)
	
    fmt.Println(rows)

	if err != nil {
	   panic(err)
	}
	defer rows.Close()


    // Цикл по условию из таблицы
	for rows.Next(){
        // Опредление перменных в которые будут заполняться данные из таблицы
		var Fam  string
		var Name string
		
		if err := rows.Scan(&Fam, &Name); err != nil {
		   panic(err)
		}
		
		fmt.Fprintln(w, Fam, Name)
	}


    // Ошибка
	if err := rows.Err(); err != nil {
       panic(err)
    }
}
