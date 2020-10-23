package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"io/ioutil"

	//using  go get gopkg.in/yaml.v2
	"gopkg.in/yaml.v2"
)

var db *sql.DB
var err error

type YamlConfig struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		User     string `yaml:"user"`
	}
}

//struck mahasiswa
type mahasiswas struct {
	Idmahasiswa string `json:"Id_mahasiswa"`
	Nama        string `json:"nama"`
	Alamat      struct {
		Jalan     string `json:"jalan"`
		Kelurahan string `json:"kelurahan"`
		Kecamatan string `json:"kecamatan"`
		Kabupaten string `json:"kabupaten"`
		Provinsi  string `json:"provinsi"`
	} `json:"alamat"`
	Fakultas string         `json:"fakultas"`
	Jurusan  string         `json:"jurusan"`
	Nilai    []nilaidetails `json:"nilai"`
}

type alamatdetails struct {
}
type nilaidetails struct {
	Namamatkul string `json:"nama"`
	Nilai      string `json:"Nilai"`
	Semester   string `json:"semester"`
}

type matkuls struct {
	Idmatkul string `json:"Id_matkul"`
	Nama     string `json:"Nama"`
}

func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mahasiswa mahasiswas
	var nilaidet nilaidetails
	params := mux.Vars(r)

	sql := `select
				mahasiswa.id_mahasiswa,
				mahasiswa.nama,
				fakultas.nama as fakultas,
				jurusan.nama as jurusan,
				mahasiswa.jalan,
				mahasiswa.kelurahan,
				mahasiswa.kecamatan,
				mahasiswa.kabupaten,
				mahasiswa.provinsi 
				FROM
				mahasiswa.mahasiswa
				INNER JOIN mahasiswa.fakultas
				ON (mahasiswa.Id_Fakultas = fakultas.id_fakultas)
				INNER JOIN mahasiswa.jurusan
				ON (mahasiswa.Id_Jurusan = jurusan.id_jurusan) where mahasiswa.id_mahasiswa=?`
	result, err := db.Query(sql, params["id"])

	defer result.Close()
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err := result.Scan(&mahasiswa.Idmahasiswa, &mahasiswa.Nama, &mahasiswa.Fakultas, &mahasiswa.Jurusan,
			&mahasiswa.Alamat.Jalan, &mahasiswa.Alamat.Kelurahan, &mahasiswa.Alamat.Kecamatan, &mahasiswa.Alamat.Kabupaten, &mahasiswa.Alamat.Provinsi)

		Idmahasiswa := &mahasiswa.Idmahasiswa

		if err != nil {
			panic(err.Error())
		}

		sqlnilai := `SELECT
						matkul.nama,nilai.nilai,nilai.semester
						FROM
							mahasiswa.nilai
							INNER JOIN mahasiswa.matkul
								ON (nilai.Id_matkul = matkul.id_matkul) where nilai.id_mahasiswa=?;`

		resultnilai, errnilai := db.Query(sqlnilai, *Idmahasiswa)

		defer resultnilai.Close()

		if errnilai != nil {
			panic(err.Error())
		}

		for resultnilai.Next() {
			err := resultnilai.Scan(&nilaidet.Namamatkul, &nilaidet.Nilai, &nilaidet.Semester)

			if err != nil {
				panic(err.Error())
			}

			mahasiswa.Nilai = append(mahasiswa.Nilai, nilaidet)
		}

	}
	// for result.Next() {
	// 	err := result.Scan(&mahasiswa.Idmahasiswa, &mahasiswa.Nama, &mahasiswa.Fakultas, &mahasiswa.Jurusan)

	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mahasiswa)
}

func createmahasiswa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Id_mahasiswa := r.FormValue("Id_Mahasiswa")
		Nama := r.FormValue("Nama")
		Jalan := r.FormValue("Jalan")
		Kelurahan := r.FormValue("Kelurahan")
		Kecamatan := r.FormValue("Kecamatan")
		Kabupaten := r.FormValue("Kabupaten")
		Provinsi := r.FormValue("Provinsi")
		Id_Fakultas := r.FormValue("Id_Fakultas")
		Id_Jurusan := r.FormValue("Id_Jurusan")

		stmt, err := db.Prepare("INSERT INTO mahasiswa (Id_Mahasiswa,Nama,Jalan,Kelurahan,Kecamatan,Kabupaten,Provinsi,Id_Fakultas,Id_Jurusan) VALUES (?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(Id_mahasiswa, Nama, Jalan, Kelurahan, Kecamatan, Kabupaten, Provinsi, Id_Fakultas, Id_Jurusan)
		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func deletemahasiswa(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM mahasiswa WHERE Id_mahasiswa = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Mahasiswa with ID = %s was deleted", params["id"])
}

func createNilai(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Id_mahasiswa := r.FormValue("Id_Mahasiswa")
		Idmatkul := r.FormValue("Id_Matkul")
		Nilai := r.FormValue("Nilai")
		Semester := r.FormValue("semester")

		stmt, err := db.Prepare("INSERT INTO nilai (Id_Mahasiswa,Id_Matkul,nilai,semester) VALUES (?,?,?,?)")

		_, err = stmt.Exec(Id_mahasiswa, Idmatkul, Nilai, Semester)
		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func creatematkul(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		Idmatkul := r.FormValue("Id_Matkul")
		Nama := r.FormValue("nama")

		stmt, err := db.Prepare("INSERT INTO matkul (Id_Matkul,nama) VALUES (?,?)")

		_, err = stmt.Exec(Idmatkul, Nama)
		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}
func main() {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	db, err = sql.Open("mysql", ""+yamlConfig.Connection.User+":"+yamlConfig.Connection.Password+"@tcp("+yamlConfig.Connection.Host+":"+yamlConfig.Connection.Port+")/mahasiswa")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/mahasiswas/{id}", getMahasiswa).Methods("GET")
	r.HandleFunc("/mahasiswas", createmahasiswa).Methods("POST")
	r.HandleFunc("/mahasiswas/{id}", deletemahasiswa).Methods("DELETE")

	//nilai
	r.HandleFunc("/mahasiswanilai", createNilai).Methods("POST")

	//matkul
	r.HandleFunc("/matkul", creatematkul).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))

}
