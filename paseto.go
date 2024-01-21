package pakarbibackend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

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

// senResponse
func sendResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

// func GCFInsertParkiranEmail(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
// 		userdata.Email = checktoken
// 		if checktoken == "" {
// 			response.Message = "Kamu kayaknya belum punya akun"
// 		} else {
// 			user2 := FindUserEmail(mconn, colluser, userdata)
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

// FUNCTION INSERT PARKIRAN 1
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
					// Menggunakan GenerateParkiranID untuk mendapatkan Parkiran ID
					dataparkiran.Parkiranid = GenerateParkiranID(dataparkiran.NPM)

					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
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
					insertParkiran(mconn, collparkiran, Parkiran{
						Parkiranid:     dataparkiran.Parkiranid,
						Nama:           dataparkiran.Nama,
						NPM:            dataparkiran.NPM,
						Prodi:          dataparkiran.Prodi,
						NamaKendaraan:  dataparkiran.NamaKendaraan,
						NomorKendaraan: dataparkiran.NomorKendaraan,
						JenisKendaraan: dataparkiran.JenisKendaraan,
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

// func GCFInsertParkiranNPM2(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
// 					// Create auto-incremented Parkiran ID
// 					parkiranID, err := createParkiranID(mconn)
// 					if err != nil {
// 						response.Message = "Error creating Parkiran ID: " + err.Error()
// 					} else {
// 						// Assign auto-incremented ID to dataparkiran
// 						dataparkiran.Parkiranid = parkiranID

// 						// Insert Parkiran data
// 						insertParkiran(mconn, collparkiran, dataparkiran)

// 						// Generate QR code with logo and base64 encoding
// 						_, err := GenerateQRCodeLogoBase64(mconn, collparkiran, dataparkiran)
// 						if err != nil {
// 							response.Message = "Error generating QR code: " + err.Error()
// 						} else {
// 							response.Status = true
// 							response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
// 						}
// 					}
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Insert data karena bukan user"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

// func GCFInsertParkiranEmail2(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
// 			user2 := FindUserEmail(mconn, colluser, userdata)
// 			if user2.Role == "user" {
// 				var dataparkiran Parkiran
// 				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
// 				if err != nil {
// 					response.Message = "Error parsing application/json: " + err.Error()
// 				} else {
// 					// Create auto-incremented Parkiran ID
// 					parkiranID, err := createParkiranID(mconn)
// 					if err != nil {
// 						response.Message = "Error creating Parkiran ID: " + err.Error()
// 					} else {
// 						// Assign auto-incremented ID to dataparkiran
// 						dataparkiran.Parkiranid = parkiranID

// 						// Insert Parkiran data
// 						insertParkiran(mconn, collparkiran, dataparkiran)

// 						// Generate QR code with logo and base64 encoding
// 						_, err := GenerateQRCodeLogoBase64(mconn, collparkiran, dataparkiran)
// 						if err != nil {
// 							response.Message = "Error generating QR code: " + err.Error()
// 						} else {
// 							response.Status = true
// 							response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
// 						}
// 					}
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Insert data karena bukan user"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

// func GCFInsertParkiranNPM2(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
// 					// Generate QR code with logo and base64 encoding
//                     _, err := GenerateQRCodeLogoBase64(mconn, collparkiran, dataparkiran)
//                     if err != nil {
//                         response.Message = "Error generating QR code: " + err.Error()
//                     } else {
//                         response.Status = true
//                         response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
//                     }
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Insert data karena bukan user"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

// func GCFInsertParkiranEmail2(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
// 			user2 := FindUserEmail(mconn, colluser, userdata)
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
// 					// Generate QR code with logo and base64 encoding
//                     _, err := GenerateQRCodeLogoBase64(mconn, collparkiran, dataparkiran)
//                     if err != nil {
//                         response.Message = "Error generating QR code: " + err.Error()
//                     } else {
//                         response.Status = true
//                         response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
//                     }
// 				}
// 			} else {
// 				response.Message = "Anda tidak dapat Insert data karena bukan user"
// 			}
// 		}
// 	}
// 	return GCFReturnStruct(response)
// }

// GCF untuk Generate Code QR
func GCFGenerateQR(publickey, MONGOCONNSTRINGENV, dbname, colluser, collparkiran string, r *http.Request) string {
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
			// Mendapatkan informasi user berdasarkan NPM dan Email
			userNPM := FindUserByField(mconn, colluser, "npm", userdata.NPM)
			userEmail := FindUserByField(mconn, colluser, "email", userdata.Email)

			// Mengecek peran user berdasarkan NPM atau Email
			var userRole string
			if userNPM.Role == "user" {
				userRole = userNPM.Role
			} else if userEmail.Role == "user" {
				userRole = userEmail.Role
			}

			if userRole == "user" {
				var dataparkiran Parkiran // Change to Parkiran type
				err := json.NewDecoder(r.Body).Decode(&dataparkiran)
				if err != nil {
					response.Message = "Error parsing application/json: " + err.Error()
				} else {
					// Insert data to MongoDB
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

					// Generate QR code with logo and base64 encoding
					qrImagePath, err := GenerateQRCodeLogoBase64(mconn, collparkiran, dataparkiran)
					if err != nil {
						response.Message = "Error generating QR code: " + err.Error()
					} else {
						// Read the PNG image file
						qrImageFile, err := os.Open(qrImagePath)
						if err != nil {
							response.Message = "Error opening QR code image file: " + err.Error()
						} else {
							defer qrImageFile.Close()

							// Read the image file into a byte slice
							qrImageBytes, err := ioutil.ReadAll(qrImageFile)
							if err != nil {
								response.Message = "Error reading QR code image file: " + err.Error()
							} else {
								// Convert the image bytes to base64
								qrImageBase64 := base64.StdEncoding.EncodeToString(qrImageBytes)

								// Use qrImageBase64 as needed, for example, inserting into MongoDB
								err = InsertQRCodeDataToMongoDB(mconn, "qrcodes", dataparkiran.Parkiranid, []byte(qrImageBase64))
								if err != nil {
									response.Message = "Error inserting QR code data to MongoDB: " + err.Error()
								} else {
									response.Status = true
									response.Message = "Berhasil Insert Data Parkiran dan Generate QR Code"
								}
							}
						}
					}
				}
			} else {
				response.Message = "Anda tidak dapat Insert data karena bukan user"
			}
		}
	}
	return GCFReturnStruct(response)
}


func GCFGetQRCode(MONGOCONNSTRINGENV, dbname, collparkiran, parkiranID string) (string, error) {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Retrieve QR code data from MongoDB
	qrCodeData, err := GetQRCodeDataFromMongoDB(mconn, "qrcodes", parkiranID)
	if err != nil {
		return "", fmt.Errorf("failed to get QR code data from MongoDB: %v", err)
	}

	// Convert base64 data to JSON string
	jsonString, err := json.Marshal(map[string]interface{}{"base64Image": qrCodeData})
	if err != nil {
		return "", fmt.Errorf("failed to convert base64 to JSON string: %v", err)
	}

	return strings.TrimSpace(string(jsonString)), nil
}

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

func GCFGetAllParkiranID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var dataparkiran Parkiran
	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
	if err != nil {
		return GCFReturnStruct(CreateResponse(false, "Failed to decode request body", nil))
	}

	// Generate Parkiran ID using the provided NPM
	dataparkiran.Parkiranid = GenerateParkiranID(dataparkiran.NPM)

	parkiran := GetAllParkiranID(mconn, collectionname, dataparkiran)
	if parkiran.Parkiranid != "" {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Parkiran", parkiran))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Parkiran", nil))
	}
}

func GCFGetOneParkiran(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var dataparkiran Parkiran
	err := json.NewDecoder(r.Body).Decode(&dataparkiran)
	if err != nil {
		return GCFReturnStruct(CreateResponse(false, "Failed to decode request body", nil))
	}
	// Generate Parkiran ID using the provided NPM
	dataparkiran.Parkiranid = GenerateParkiranID(dataparkiran.NPM)
	// Assuming you have a function to retrieve parkiran by ID, NPM, and Nama Mahasiswa
	parkiran := GetOneParkiranData(mconn, collectionname, dataparkiran.Parkiranid)  // Perubahan pada argumen pemanggilan

	if parkiran != (Parkiran{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get Parkiran Data", parkiran))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get Parkiran Data", nil))
	}
}
