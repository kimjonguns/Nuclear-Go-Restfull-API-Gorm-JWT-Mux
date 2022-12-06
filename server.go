package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error
var Response *Result
var Hasil []byte
var BarangBody *Barang
var barangs *[]Barang
var BarangNow *Barang

// model Barang
type Barang struct {
	Id    int             `form:"id" json:"id"`
	Code  string          `form:"code" json:"code"`
	Name  string          `form:"name" json:"name"`
	Price decimal.Decimal `form:"price" json:"price" sql:"type:decimal(16,2);"`
}

// model result
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Env struct {
	USER     string
	PASSWORD string
	DB       string
}

func main() {
	dataenv := Env{}
	err = godotenv.Load()
	if err != nil {
		log.Println("Error")
	} else {
		dataenv.USER = os.Getenv("USER")
		dataenv.PASSWORD = os.Getenv("PASSWORD")
		dataenv.DB = os.Getenv("DB")
		db, err = gorm.Open("mysql", dataenv.USER+":"+dataenv.PASSWORD+"@/"+dataenv.DB+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			log.Println("Error", err)
		} else {

			log.Println("Connection success")
			// generate
			db.AutoMigrate(&Barang{})

			handleRequest()
		}
	}
}

func handleRequest() {
	log.Println("handle in http://127.0.0.1:5000")
	Router := mux.NewRouter().StrictSlash(true)

	// Home
	Router.HandleFunc("/", HomePage).Methods("GET")
	// Get Barangs
	Router.HandleFunc("/barangs", BarangPage).Methods("GET")
	// Get Barang by id
	Router.HandleFunc("/barangs/{id}", GetProductById).Methods("GET")
	// Post Barang
	Router.HandleFunc("/barangs", PostBarang).Methods("POST")
	// Put Barang
	Router.HandleFunc("/barangs/{id}", UpdateBarang).Methods("PUT")
	// Delete Barang
	Router.HandleFunc("/barangs/{id}", DeleteBarang).Methods("DELETE")
	// berjalan
	log.Fatal(http.ListenAndServe(":9999", Router))
}

// GET
func HomePage(w http.ResponseWriter, r *http.Request) {
	Response = &Result{Code: 200, Data: "", Message: "Welcome"}
	Hasil, err = json.Marshal(Response)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Hasil)
	}
}

func BarangPage(w http.ResponseWriter, r *http.Request) {
	barangs = &[]Barang{}
	fmt.Println(&barangs)
	db.Find(&barangs)
	Response = &Result{Code: 200, Data: *barangs, Message: "The datas"}
	Hasil, err = json.Marshal(Response)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Hasil)
	}
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	barangs = &[]Barang{}
	db.Find(&barangs, id)
	Response = &Result{Code: 200, Data: *barangs, Message: "The data"}
	Hasil, err = json.Marshal(Response)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Hasil)
	}
}

// POST
func PostBarang(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(payloads, &BarangBody)
	// fmt.Println(BarangBody) -> &{hasil}
	// fmt.Println(&BarangBody)  -> lokasi
	db.Create(&BarangBody)
	Response = &Result{Code: 200, Data: BarangBody, Message: "Post Data"}
	Hasil, err = json.Marshal(Response)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Hasil)
	}
}

// PUT
func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	payloads, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(payloads, &BarangBody)
	BarangNow = &Barang{}
	db.First(&BarangNow, id)
	db.Model(&BarangNow).Updates(BarangBody)
	Response = &Result{Code: 200, Data: BarangNow, Message: "Update Data"}
	Hasil, err = json.Marshal(Response)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Hasil)
	}
}

// Delete
func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	BarangNow = &Barang{}
	db.First(&BarangNow, id)
	db.Delete(&BarangNow)
	Response = &Result{Code: 200, Message: "Delete Data with id " + id}
	Hasil, err = json.Marshal(Response)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Hasil)
	}
}
