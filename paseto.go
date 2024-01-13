package pakarbibackend

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// < --- FUNCTION CEK EMAIL --- >

func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		regexPattern: `^[a-zA-Z0-9._%+-]+@std.ulbi.ac.id$`,
	}
}

// IsValid memeriksa apakah email sesuai dengan pola npm@std.ulbi.ac.id
func (v *EmailValidator) IsValid(email string) bool {
	match, _ := regexp.MatchString(v.regexPattern, email)
	return match
}

// <--- FUNCTION USER --->
func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	userdata := new(User)
	resp.Status = false
	conn := MongoCreateConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true

		// Validasi email sebelum proses pendaftaran
		validator := NewEmailValidator()
		if !validator.IsValid(userdata.Email) {
			resp.Message = "Email is not valid"
			resp.Status = false
			response := ReturnStringStruct(resp)
			return response
		}

		hash, err := HashPass(userdata.PasswordHash)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertUserdata(conn, userdata.UsernameId, userdata.Username, userdata.NPM, userdata.Password, hash, userdata.Email, userdata.Role)
		resp.Message = "Berhasil Input data"
	}
	response := ReturnStringStruct(resp)
	return response
}

func LoginUserNPM(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValidNPM(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.NPM, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "NPM atau Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func LoginUserEmail(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		// Validasi email harus menggunakan npm@std.ulbi.ac.id sesuai dengan email kampus didaftarkan sebelum melakukan login
		validator := NewEmailValidator()
		if !validator.IsValid(datauser.Email) {
			Response.Message = "Email is not valid"
			response := GCFReturnStruct(Response)
			return response
		}

		if IsPasswordValidEmail(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Email, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Email atau Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GetAllDataUser(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data user
			datauser := GetAllUser(conn, colname)
			if datauser == nil {
				req.Status = false
				req.Message = "Data User tidak ada"
			} else {
				req.Status = true
				req.Message = "Data User berhasil diambil"
				req.Data = datauser
			}
		}
	}
	return ReturnStringStruct(req)
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

// <--- FUNCTION ADMIN --->
func LoginAdmin(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var dataadmin Admin
	err := json.NewDecoder(r.Body).Decode(&dataadmin)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		// Validasi email harus menggunakan npm@std.ulbi.ac.id sesuai dengan email kampus didaftarkan sebelum melakukan login
		validator := NewEmailValidator()
		if !validator.IsValid(dataadmin.Email) {
			Response.Message = "Email is not valid"
			response := GCFReturnStruct(Response)
			return response
		}

		if IsPasswordValidEmailAdmin(mconn, collectionname, dataadmin) {
			Response.Status = true
			tokenstring, err := watoken.Encode(dataadmin.Email, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang Admin"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Email atau Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

// <--- FUNCTION PARKIRAN --->
// func GCFInsertParkiranNPM(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
// 	var response Credential
// 	response.Status = false
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
// 	var userdata User
// 	gettoken := r.Header.Get("Login")
// 	if gettoken == "" {
// 		response.Message = "Header Login Not Exist"
// 	} else {
// 		// Process the request with the "Login" token
// 		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
// 		userdata.NPM = checktoken
// 		if checktoken == "" {
// 			response.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			user2 := FindUserNPM(mconn, colluser, userdata)
// 			if user2.Role == "user" {
// 				var dataparkiran Parkiran
// 				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 				if err != nil {
// 					response.Message = "Error parsing application/json: " + err.Error()
// 				} else {
// 					insertParkiran(mconn, collparkiran, Parkiran{
// 						Parkiranid:     dataparkiran.Parkiranid,
// 						Nama:           dataparkiran.Nama,
// 						NPM:            dataparkiran.NPM,
// 						Prodi:          dataparkiran.Prodi,
// 						NamaKendaraan:  dataparkiran.NamaKendaraan,
// 						NomorKendaraan: dataparkiran.NomorKendaraan,
// 						JenisKendaraan: dataparkiran.JenisKendaraan,
// 						Status:         dataparkiran.Status,
// 					})
// 					response.Status = true
// 					response.Message = "Berhasil Insert Data Parkiran"
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Insert data karena bukan user"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

func GCFInsertParkiranNPM(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.NPM = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserNPM(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Generate ParkiranID using the function
					dataparkiran.Parkiranid = GenerateParkiranID(dataparkiran.NPM)

					// Memeriksa apakah sudah ada waktu keluar, jika tidak maka dianggap 'sudah masuk Parkir'
					if dataparkiran.Status.WaktuKeluar == "" {
						dataparkiran.Status = Status{
							Message:    "sudah masuk Parkir",
							WaktuMasuk: time.Now().Format(time.RFC3339),
						}
					} else {
						dataparkiran.Status = Status{
							Message:     "sudah keluar Parkir",
							WaktuKeluar: dataparkiran.Status.WaktuKeluar,
							WaktuMasuk:  dataparkiran.Status.WaktuMasuk,
						}
					}

					// Insert parkiran data
					fileName, err := GenerateQRCodeWithLogo(mconn, dataparkiran)
					if err != nil {
						response.Message = "Error generating QR code: " + err.Error()
					} else {
						response.Status = true
						response.Message = "Berhasil Insert Data Parkiran"
						response.Data = fileName // Add the file name to the response
					}
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}

	return GCFReturnStruct(response)
}

func GCFInsertParkiranEmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserEmail(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Generate ParkiranID using the function
					dataparkiran.Parkiranid = GenerateParkiranID(dataparkiran.NPM)

					// Memeriksa apakah sudah ada waktu keluar, jika tidak maka dianggap 'sudah masuk Parkir'
					if dataparkiran.Status.WaktuKeluar == "" {
						dataparkiran.Status = Status{
							Message:    "sudah masuk Parkir",
							WaktuMasuk: time.Now().Format(time.RFC3339),
						}
					} else {
						dataparkiran.Status = Status{
							Message:     "sudah keluar Parkir",
							WaktuKeluar: dataparkiran.Status.WaktuKeluar,
							WaktuMasuk:  dataparkiran.Status.WaktuMasuk,
						}
					}

					// Insert parkiran data
					fileName, err := GenerateQRCodeWithLogo(mconn, dataparkiran)
					if err != nil {
						response.Message = "Error generating QR code: " + err.Error()
					} else {
						response.Status = true
						response.Message = "Berhasil Insert Data Parkiran"
						response.Data = fileName // Add the file name to the response
					}
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}

	return GCFReturnStruct(response)
}

// <--- FUNCTION INSERT PARKIRAN 2 --->
func GCFInsertParkiranNPM2(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	// Inisialisasi folder QR code
	err := InitQRCodeFolder()
	if err != nil {
		return GCFReturnStruct(Credential{Status: false, Message: "Failed to initialize QR code folder"})
	}

	// Set koneksi MongoDB
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Inisialisasi respons
	var response Credential
	response.Status = false

	// Mendapatkan data token dari header
	gettoken := r.Header.Get("Login")

	// Memeriksa apakah token ada
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Proses permintaan dengan token "Login"
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		var userdata User
		userdata.NPM = checktoken

		// Memeriksa apakah token valid
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserNPM(mconn, colluser, userdata)

			// Memeriksa apakah pengguna memiliki peran "user"
			if user2.Role == "user" {
				var dataparkiran Parkiran

				// Membaca dan mendecode data parkiran dari body request
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Menyisipkan data parkiran ke MongoDB
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
						Status:         dataparkiran.Status,
					})

					// Generate QR code tanpa logo
					qrOutputPath := "qrcode/" + dataparkiran.Parkiranid + "_qrcode.png"
					err := GenerateQRCode(dataparkiran, qrOutputPath)
					if err != nil {
						response.Message = "Failed to generate QR code: " + err.Error()
						return GCFReturnStruct(response)
					}

					// Setel respons berhasil
					response.Status = true
					response.Message = "Berhasil Insert Data Parkiran"
					response.Data = qrOutputPath
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}

	// Mengembalikan respons dalam bentuk string JSON
	return GCFReturnStruct(response)
}

func GCFInsertParkiranNPM3(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	// Inisialisasi folder QR code
	err := InitQRCodeFolder()
	if err != nil {
		return GCFReturnStruct(Credential{Status: false, Message: "Failed to initialize QR code folder"})
	}

	// Set koneksi MongoDB
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Inisialisasi respons
	var response Credential
	response.Status = false

	// Mendapatkan data token dari header
	gettoken := r.Header.Get("Login")

	// Memeriksa apakah token ada
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Proses permintaan dengan token "Login"
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		var userdata User
		userdata.NPM = checktoken

		// Memeriksa apakah token valid
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserNPM(mconn, colluser, userdata)

			// Memeriksa apakah pengguna memiliki peran "user"
			if user2.Role == "user" {
				var dataparkiran Parkiran

				// Membaca dan mendecode data parkiran dari body request
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Menyisipkan data parkiran ke MongoDB
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
						Status:         dataparkiran.Status,
					})
					// Generate QR code tanpa logo
					qrOutputPath := filepath.Join("C:\\Users\\ACER\\Documents\\pakarbibackend\\qrcode", dataparkiran.Parkiranid+"_qrcode.png")
					err := GenerateQRCode(dataparkiran, qrOutputPath)
					if err != nil {
						response.Message = "Failed to generate QR code: " + err.Error()
						return GCFReturnStruct(response)
					}

					// Setel respons berhasil
					response.Status = true
					response.Message = "Berhasil Insert Data Parkiran"
					response.Data = qrOutputPath
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}

	// Mengembalikan respons dalam bentuk string JSON
	return GCFReturnStruct(response)
}

func GCFInsertParkiranNPM4(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	// Inisialisasi folder QR code
	err := InitQRCodeFolder()
	if err != nil {
		return GCFReturnStruct(Credential{Status: false, Message: "Failed to initialize QR code folder"})
	}

	// Set koneksi MongoDB
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Inisialisasi respons
	var response Credential
	response.Status = false

	// Mendapatkan data token dari header
	gettoken := r.Header.Get("Login")

	// Memeriksa apakah token ada
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Proses permintaan dengan token "Login"
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		var userdata User
		userdata.NPM = checktoken

		// Memeriksa apakah token valid
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserNPM(mconn, colluser, userdata)

			// Memeriksa apakah pengguna memiliki peran "user"
			if user2.Role == "user" {
				var dataparkiran Parkiran

				// Membaca dan mendecode data parkiran dari body request
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Menyisipkan data parkiran ke MongoDB
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
						Status:         dataparkiran.Status,
					})

					// Generate QR code tanpa logo
					qrOutputPath := filepath.Join("C:\\Users\\ACER\\Documents\\pakarbibackend\\qrcode", dataparkiran.Parkiranid+"_qrcode.png")
					err := GenerateQRCode(dataparkiran, qrOutputPath)
					if err != nil {
						response.Message = "Failed to generate QR code: " + err.Error()
						return GCFReturnStruct(response)
					}

					// Setel respons berhasil
					response.Status = true
					response.Message = "Berhasil Insert Data Parkiran"
					response.Data = qrOutputPath
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}

	// Mengembalikan respons dalam bentuk string JSON
	return GCFReturnStruct(response)
}

func GCFInsertParkiranEmail2(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserEmail(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
						Status:         dataparkiran.Status,
					})
					response.Status = true
					response.Message = "Berhasil Insert Data Parkiran"
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

// <--- FUNCTION UPDATE PARKIRAN 3--->
// func GCFInsertParkiranNPM3(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
//     var response Credential
//     response.Status = false
//     mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
//     var userdata User
//     gettoken := r.Header.Get("Login")
//     if gettoken == "" {
//         response.Message = "Header Login Not Exist"
//     } else {
//         // Process the request with the "Login" token
//         checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
//         userdata.NPM = checktoken
//         if checktoken == "" {
//             response.Message = "Kamu kayaknya belum punya akun"
//         } else {
//             user2 := FindUserNPM(mconn, colluser, userdata)
//             if user2.Role == "user" {
//                 var dataparkiran Parkiran
//                 err := json.NewDecoder(r.Body).Decode(&dataparkiran)
//                 if err != nil {
//                     response.Message = "Error parsing application/json: " + err.Error()
//                 } else {
//                     // Generate QR code without logo
//                     qrOutputPath := filepath.Join("C:\\Users\\Muhammad Faisal A\\OneDrive\\Pictures\\Code QR" + dataparkiran.Parkiranid + "_qrcode.png")
//                     err := GenerateQRCode(dataparkiran, qrOutputPath)
//                     if err != nil {
//                         response.Message = "Failed to generate QR code: " + err.Error()
//                         return GCFReturnStruct(response)
//                     }

//                     // Simpan gambar kode QR ke MongoDB
//                     err = SaveQRCodeToMongoDB(qrOutputPath, MONGOCONNSTRINGENV, dbname, "qrcode")
//                     if err != nil {
//                         response.Message = "Failed to save QR code to MongoDB: " + err.Error()
//                         return GCFReturnStruct(response)
//                     }

//                     // Insert parkiran data
//                     insertParkiran(mconn, collparkiran, Parkiran{
//                         Parkiranid:     dataparkiran.Parkiranid,
//                         Nama:           dataparkiran.Nama,
//                         NPM:            dataparkiran.NPM,
//                         Prodi:          dataparkiran.Prodi,
//                         NamaKendaraan:  dataparkiran.NamaKendaraan,
//                         NomorKendaraan: dataparkiran.NomorKendaraan,
//                         JenisKendaraan: dataparkiran.JenisKendaraan,
//                         Status:         dataparkiran.Status,
//                     })

//                     response.Status = true
//                     response.Message = "Berhasil Insert Data Parkiran"
//                     response.Data = qrOutputPath // Menambahkan path QR code ke respons
//                 }
//             } else {
//                 response.Message = "Anda tidak dapat Insert data karena bukan user"
//             }
//         }
//     }
//     return GCFReturnStruct(response)
// }

// GCF Update Data
func GCFUpdateParkiranNPM(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.NPM = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserNPM(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					UpdatedParkiran(mconn, collparkiran, bson.M{"id": dataparkiran.ID}, dataparkiran)
					response.Status = true
					response.Message = "Berhasil Update Parkiran"
					GCFReturnStruct(CreateResponse(true, "Success Update Parkiran", dataparkiran))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFUpdateParkiranEmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.Email = checktoken
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserEmail(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					UpdatedParkiran(mconn, collparkiran, bson.M{"id": dataparkiran.ID}, dataparkiran)
					response.Status = true
					response.Message = "Berhasil Update Parkiran"
					GCFReturnStruct(CreateResponse(true, "Success Update Parkiran", dataparkiran))
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

//GCF Hapus Data

func GCFDeleteParkiranNPM(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.NPM = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserNPM(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteParkiran(mconn, collparkiran, dataparkiran)
					respon.Status = true
					respon.Message = "Berhasil Delete Parkiran"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(respon)
}

func GCFDeleteParkiranEmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var respon Credential
	respon.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User

	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		respon.Message = "Header Login Not Exist"
	} else {
		// Process the request with the "Login" token
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		userdata.Email = checktoken
		if checktoken == "" {
			respon.Message = "Kamu kayaknya belum punya akun"
		} else {
			user2 := FindUserEmail(mconn, colluser, userdata)
			if user2.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					respon.Message = "Error parsing application/json: " + err.Error()
				} else {
					DeleteParkiran(mconn, collparkiran, dataparkiran)
					respon.Status = true
					respon.Message = "Berhasil Delete Parkiran"
				}
			} else {
				respon.Message = "Anda tidak dapat Delete data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(respon)
}

func GetAllDataParkiran(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Response)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		// Dekode token untuk mendapatkan
		_, err := DecodeGetParkiran(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Data Tersebut tidak ada" + tokenlogin
		} else {
			// Langsung ambil data catalog
			dataparkiran := GetAllParkiran(conn, colname)
			if dataparkiran == nil {
				req.Status = false
				req.Message = "Data Parkiran tidak ada"
			} else {
				req.Status = true
				req.Message = "Data Parkiran berhasil diambil"
				req.Data = dataparkiran
			}
		}
	}
	return ReturnStringStruct(req)
}

func GetOneDataParkiran(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseParkiran)
	resp := new(RequestParkiran)
	conn := MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			dataparkiran := GetOneParkiranData(conn, colname, resp.Parkiranid)
			req.Status = true
			req.Message = "data Parkiran berhasil diambil"
			req.Data = dataparkiran
		}
	}
	return ReturnStringStruct(req)
}

// func GCFGetAllParkiranID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	var dataparkiran Parkiran
// 	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	parkiran := GetAllParkiranID(mconn, collectionname, dataparkiran)
// 	if parkiran != (Parkiran{}) {
// 		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Parkiran", dataparkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Parkiran", dataparkiran))
// 	}
// }

// func GCFGetParkiranById(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	var dataparkiran Parkiran
// 	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 	if err != nil {
// 		return GCFReturnStruct(CreateResponse(false, "Error parsing JSON: "+err.Error(), nil))
// 	}

// 	parkiran, err := GetParkiranById(mconn, collectionname, dataparkiran.Parkiranid)
// 	if err != nil {
// 		return GCFReturnStruct(CreateResponse(false, "Failed to Get Parkiran by ID: "+err.Error(), nil))
// 	}

// 	if parkiran != (Parkiran{}) {
// 		return GCFReturnStruct(CreateResponse(true, "Success: Get Parkiran by ID", parkiran))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "No parkiran found with ID: "+dataparkiran.Parkiranid, nil))
// 	}
// }
