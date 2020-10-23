package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

func main() {

	url := "http://localhost:8080/mahasiswas/1811082007"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	mahasiswa := mahasiswas{}
	jsonErr := json.Unmarshal(body, &mahasiswa)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println("ID Mahasiswa :", mahasiswa.Idmahasiswa)
	fmt.Println("Nama :", mahasiswa.Nama)
	fmt.Println("Fakultas :", mahasiswa.Fakultas)
	fmt.Println("Jurusan :", mahasiswa.Jurusan)
	fmt.Println("Alamat :")

	fmt.Println("Jalan : ", mahasiswa.Alamat.Jalan)
	fmt.Println("Kelurahan : ", mahasiswa.Alamat.Kelurahan)
	fmt.Println("Kecamatan : ", mahasiswa.Alamat.Kecamatan)
	fmt.Println("Kabupaten : ", mahasiswa.Alamat.Kabupaten)
	fmt.Println("Provinsi : ", mahasiswa.Alamat.Provinsi)

	for _, nilai := range mahasiswa.Nilai {
		fmt.Println("Nama Matkul : ", nilai.Namamatkul)
		fmt.Println("Nilai : ", nilai.Nilai)
		fmt.Println("Semester : ", nilai.Semester)

	}

}
