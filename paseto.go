package pakarbibackend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// <---GCF untuk Generate Code QR--->
func GCFGenerateCodeQR(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			// Mengecek token yang digunakan
			userdata.NPM = checktoken
			user := FindUserByField(mconn, colluser, "npm", userdata.NPM)
			if user.NPM == "" {
				// Jika tidak menemukan user menggunakan npm, cobain menggunakan email
				userdata.Email = checktoken
				user = FindUserByField(mconn, colluser, "email", userdata.Email)
			}

			if user.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Generate Parkiran ID
					// Insert Parkiran data
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
						Status:         dataparkiran.Status,
						JamMasuk:       dataparkiran.JamMasuk,
						JamKeluar:      dataparkiran.JamKeluar,
					})

					// Generate QR code with logo and base64 encoding
					_, err := GenerateQRCodeBase64(mconn, collparkiran, dataparkiran)
					if err != nil {
						response.Message = "Error generating QR code: " + err.Error()
					} else {
						response.Status = true
						response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
					}
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFGenerateCodeQREmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			// Mengecek token yang digunakan
			userdata.Email = checktoken
			user := FindUserByField(mconn, colluser, "email", userdata.Email)
			if user.Email == "" {
				// Jika tidak menemukan user menggunakan npm, cobain menggunakan email
				userdata.NPM = checktoken
				user = FindUserByField(mconn, colluser, "npm", userdata.NPM)
			}

			if user.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Generate Parkiran ID
					// Insert Parkiran data
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
						Status:         dataparkiran.Status,
						JamMasuk:       dataparkiran.JamMasuk,
						JamKeluar:      dataparkiran.JamKeluar,
					})

					// Generate QR code with logo and base64 encoding
					_, err := GenerateQRCodeBase64(mconn, collparkiran, dataparkiran)
					if err != nil {
						response.Message = "Error generating QR code: " + err.Error()
					} else {
						response.Status = true
						response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
					}
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFUpdateGenerateCodeQR(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			// Mengecek token yang digunakan
			userdata.NPM = checktoken
			user := FindUserByField(mconn, colluser, "npm", userdata.NPM)
			if user.NPM == "" {
				// Jika tidak menemukan user menggunakan npm, cobain menggunakan email
				userdata.Email = checktoken
				user = FindUserByField(mconn, colluser, "email", userdata.Email)
			}

			if user.Role == "user" {
				var newDataParkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&newDataParkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Update Parkiran data
					err := UpdateParkiran(mconn, collparkiran, newDataParkiran)
					if err != nil {
						response.Message = "Error updating Parkiran data: " + err.Error()
					} else {
						// Generate QR code with logo and base64 encoding
						_, err := GenerateQRCodeBase64(mconn, collparkiran, newDataParkiran)
						if err != nil {
							response.Message = "Error generating QR code: " + err.Error()
						} else {
							response.Status = true
							response.Message = "Berhasil Update Data Parkiran dan Generate Ulang QR Code"
						}
					}
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFUpdateGenerateCodeQREmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			// Mengecek token yang digunakan
			userdata.Email = checktoken
			user := FindUserByField(mconn, colluser, "email", userdata.Email)
			if user.Email == "" {
				// Jika tidak menemukan user menggunakan npm, cobain menggunakan email
				userdata.NPM = checktoken
				user = FindUserByField(mconn, colluser, "npm", userdata.NPM)
			}

			if user.Role == "user" {
				var newDataParkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&newDataParkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Update Parkiran data
					err := UpdateParkiran(mconn, collparkiran, newDataParkiran)
					if err != nil {
						response.Message = "Error updating Parkiran data: " + err.Error()
					} else {
						// Generate QR code with logo and base64 encoding
						_, err := GenerateQRCodeBase64(mconn, collparkiran, newDataParkiran)
						if err != nil {
							response.Message = "Error generating QR code: " + err.Error()
						} else {
							response.Status = true
							response.Message = "Berhasil Update Data Parkiran dan Generate Ulang QR Code"
						}
					}
				}
			} else {
				response.Message = "Anda tidak dapat Update data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFDeleteGenerateCodeQR(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Proses request dengan token "Login"
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			// Mengecek token yang digunakan
			userdata.NPM = checktoken
			user := FindUserByField(mconn, colluser, "npm", userdata.NPM)
			if user.NPM == "" {
				// Jika tidak menemukan user menggunakan npm, cobain menggunakan email
				userdata.Email = checktoken
				user = FindUserByField(mconn, colluser, "email", userdata.Email)
			}

			if user.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Menghapus data parkiran dan QR code
					err := DeleteQRCodeData2(mconn, collparkiran, dataparkiran.Parkiranid)
					if err != nil {
						response.Message = "Error deleting parkiran data and QR code: " + err.Error()
					} else {
						response.Status = true
						response.Message = "Berhasil Hapus Data Parkiran dan QR Code"
					}
				}
			} else {
				response.Message = "Anda tidak dapat Hapus data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFDeleteGenerateCodeQREmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	gettoken := r.Header.Get("Login")
	if gettoken == "" {
		response.Message = "Header Login Not Exist"
	} else {
		// Proses request dengan token "Login"
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), gettoken)
		if checktoken == "" {
			response.Message = "Kamu kayaknya belum punya akun"
		} else {
			// Mengecek token yang digunakan
			userdata.Email = checktoken
			user := FindUserByField(mconn, colluser, "email", userdata.Email)
			if user.Email == "" {
				// Jika tidak menemukan user menggunakan npm, cobain menggunakan email
				userdata.NPM = checktoken
				user = FindUserByField(mconn, colluser, "npm", userdata.NPM)
			}

			if user.Role == "user" {
				var dataparkiran Parkiran
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Menghapus data parkiran dan QR code
					err := DeleteQRCodeData2(mconn, collparkiran, dataparkiran.Parkiranid)
					if err != nil {
						response.Message = "Error deleting parkiran data and QR code: " + err.Error()
					} else {
						response.Status = true
						response.Message = "Berhasil Hapus Data Parkiran dan QR Code"
					}
				}
			} else {
				response.Message = "Anda tidak dapat Hapus data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}

func GCFGetAllDataParkiran(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
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

func GCFGetOneDataParkiran(PublicKey, MongoEnv, dbname, colname, parkiranID string, r *http.Request) string {
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
			// Langsung ambil satu data parkiran
			dataparkiran, err := GetOneDataParkiranByID(conn, colname, parkiranID)
			if err != nil {
				req.Status = false
				req.Message = fmt.Sprintf("Data Parkiran dengan ID %s tidak ditemukan", parkiranID)
			} else {
				req.Status = true
				req.Message = "Data Parkiran berhasil diambil"
				req.Data = dataparkiran
			}
		}
	}
	return ReturnStringStruct(req)
}

func GetIDDataParkiran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	parkirandata := new(Parkiran)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&parkirandata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	parkirandata.ID = ID

	// Menggunakan fungsi GetProdukFromID untuk mendapatkan data produk berdasarkan ID
	parkirandata, err = GetParkiranFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Data1 = []Parkiran{*parkirandata}

	return GCFReturnStruct(resp)
}
