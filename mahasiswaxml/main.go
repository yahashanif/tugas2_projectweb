package main

import (
	"database/sql"
	"encoding/xml"
	"log"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//using  go get gopkg.in/yaml.v2
)

var db *sql.DB
var err error

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
	Jalan     string `json:"jalan"`
	Kelurahan string `json:"kelurahan"`
	Kecamatan string `json:"kecamatan"`
	Kabupaten string `json:"kabupaten"`
	Provinsi  string `json:"provinsi"`
}
type nilaidetails struct {
	Namamatkul string `json:"nama"`
	Nilai      string `json:"Nilai"`
	Semester   string `json:"semester"`
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

	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	xml.NewEncoder(w).Encode(mahasiswa)
}

func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mahasiswa")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/mahasiswas/{id}", getMahasiswa).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))

}
