package pakarbibackend

import (
	"fmt"
	// "os"
	// "path/filepath"
	"io/ioutil"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.UsernameId = "D4TI1214001"
	userdata.Username = "pakarbi2"
	userdata.NPM = "1214001"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214001@std.ulbi.ac.id"
	userdata.Role = "user"
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	CreateNewUserRole(mconn, "user", userdata)
}

// func TestDeleteUser(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "pasabarapk")
// 	var userdata User
// 	userdata.Email = "1214001@std.ulbi.ac.id"
// 	DeleteUser(mconn, "user", userdata)
// }

func CreateNewUserToken(t *testing.T) {
	var userdata User
	userdata.UsernameId = "D4TI1214001"
	userdata.Username = "pakarbi2"
	userdata.NPM = "1214001"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214001@std.ulbi.ac.id"
	userdata.Role = "user"

	// Create a MongoDB connection
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")

	// Call the function to create a admin and generate a token
	err := CreateUserAndAddToken("", mconn, "user", userdata)

	if err != nil {
		t.Errorf("Error creating user and token: %v", err)
	}
}

func TestGFCPostHandlerUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.UsernameId = "D4TI1214001"
	userdata.Username = "pakarbi2"
	userdata.NPM = "1214001"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214001@std.ulbi.ac.id"
	userdata.Role = "user"
	CreateNewUserRole(mconn, "user", userdata)
}

func TestParkiran(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var parkirandata Parkiran
	parkirandata.Parkiranid = "D41214020"
	parkirandata.Nama = "Farhan Rizki Maulana"
	parkirandata.NPM = "1214020"
	parkirandata.Prodi = "D4 Teknik Informatika"
	parkirandata.NamaKendaraan = "Supra X 125"
	parkirandata.NomorKendaraan = "F 1234 NR"
	parkirandata.JenisKendaraan = "Motor"
	CreateNewParkiran(mconn, "parkiran", parkirandata)
}

func TestAllParkiran(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	parkiran := GetAllParkiran(mconn, "parkiran")
	fmt.Println(parkiran)
}

func TestGeneratePasswordHash(t *testing.T) {
	passwordhash := "pakarbipass"
	hash, _ := HashPass(passwordhash) // ignore error for the sake of simplicity

	fmt.Println("Password:", passwordhash)
	fmt.Println("Hash:    ", hash)
	match := CheckPasswordHash(passwordhash, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("pakarbipass", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.NPM = "1214001"
	userdata.PasswordHash = "pakarbipass"

	filter := bson.M{"npm": userdata.NPM}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPass(userdata.PasswordHash)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.NPM = "1214001"
	userdata.PasswordHash = "pakarbipass"

	anu := IsPasswordValidNPM(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.UsernameId = "D4TI1214001"
	userdata.Username = "pakarbi2"
	userdata.NPM = "1214001"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214001@std.ulbi.ac.id"
	userdata.Role = "user"
	CreateUser(mconn, "user", userdata)
}

func TestAdminFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var admindata Admin
	admindata.UsernameId = "Pakarbisa2023"
	admindata.Username = "adminpakarbi"
	admindata.Password = "adminpakarbipass"
	admindata.PasswordHash = "adminpakarbipass"
	admindata.Email = "PakArbi2023@std.ulbi.ac.id"
	admindata.Role = "admin"
	CreateAdmin(mconn, "admin", admindata)
}

func TestGeneratePrivateKeyPasetoAdmin(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("adminpakarbipass", privateKey)
	fmt.Println(hasil, err)
}

func TestLoginn(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.NPM = "1214001"
	userdata.PasswordHash = "pakarbipass"
	IsPasswordValidNPM(mconn, "user", userdata)
	fmt.Println(userdata)
}

//proses untuk generate code qr
// func TestGenerateQRCodeWithLogo(t *testing.T) {
//     // Initialize your MongoDB connection here
//     mconn := SetConnection("MONGOSTRING", "PakArbiApp")

//     // Initialize a sample Parkiran struct for testing
//     dataparkiran := Parkiran{
//         Parkiranid:     "D3/D412345", // Use the same value as in the expectation
//         Nama:           "John Doe",
//         NPM:            "12345",     //
//         Prodi:          "Computer Science",
//         NamaKendaraan:  "Car",
//         NomorKendaraan: "AB 1234 CD",
//         JenisKendaraan: "Sedan",
//         Status: Status{
//             Message:    "sudah masuk Parkir",
//             WaktuMasuk: time.Now().Format(time.RFC3339),
//         },
//     }

//     fileName, err := GenerateQRCodeWithLogo(mconn, dataparkiran)
//     if err != nil {
//         t.Errorf("Error generating QR code: %v", err)
//         return
//     }

//     if fileName != "qrcode/D3/D412345_qrcode.png" {
//         t.Errorf("Expected file name 'qrcode/D3/D412345_qrcode.png', got '%s'", fileName)
//         return
//     }

//     t.Log("Berhasil generate code qr")
// }

// func TestGenerateQRCodeWithLogo(t *testing.T) {
//     // Set up your MongoDB connection
//     mconn := SetConnection("MONGOSTRING", "PakArbiApp")

//     // Set up a sample Parkiran struct for testing
//     dataparkiran := Parkiran{
//         Parkiranid:     "D31214000",
//         Nama:           "John Doe",
//         NPM:            "12345",
//         Prodi:          "Computer Science",
//         NamaKendaraan:  "Car",
//         NomorKendaraan: "AB 1234 CD",
//         JenisKendaraan: "Sedan",
//         Status: Status{
//             Message:    "sudah masuk Parkir",
//             WaktuMasuk: time.Now().Format(time.RFC3339),
//         },
//     }

//     // Generate QR code with logo
//     fileName, err := GenerateQRCodeWithLogo(mconn, dataparkiran)
//     if err != nil {
//         t.Errorf("Error generating QR code: %v", err)
//         return
//     }

//     // Check if the generated file exists
//     if _, err := os.Stat(fileName); os.IsNotExist(err) {
//         t.Errorf("Expected file '%s' not found", fileName)
//         return
//     }

//     // Check if the file name follows the expected pattern
//     expectedFileName := filepath.Base(dataparkiran.Parkiranid + "_logo_ulbi_qrcode.png")
//     if fileName != expectedFileName {
//         t.Errorf("Expected file name '%s', got '%s'", expectedFileName, fileName)
//         return
//     }

//     t.Log("Successfully generated QR code with logo")
// }

func TestInsertQRCodeDataToMongoDB(t *testing.T) {
	// Set up your MongoDB connection
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")

	// Set up a sample Parkiran struct for testing
	dataparkiran := Parkiran{
		Parkiranid:     "D41214041",
		Nama:           "ULBI",
		NPM:            "1214041",
		Prodi:          "Computer Science",
		NamaKendaraan:  "Car",
		NomorKendaraan: "AB 1234 CD",
		JenisKendaraan: "Sedan",
		Status: Status{
			Message:     "sudah masuk Parkir",
			WaktuMasuk:  "some_waktu_masuk",
			WaktuKeluar: "some_waktu_keluar",
			// WaktuMasuk: time.Now().Format(time.RFC3339),

		},
	}

	// Generate QR code with logo and insert into MongoDB
	fileName, err := GenerateQRCodeLogoBase64(mconn, "parkiran", dataparkiran)
	if err != nil {
		t.Errorf("Error generating QR code with logo: %v", err)
		return
	}

	// Read the QR code file
	qrCodeData, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("Error reading QR code file: %v", err)
		return
	}

	// Insert QR code data into MongoDB
	err = InsertQRCodeDataToMongoDB(mconn, "qrcodes", dataparkiran.Parkiranid, qrCodeData)
	if err != nil {
		t.Errorf("Error inserting QR code data to MongoDB: %v", err)
		return
	}

	t.Log("Successfully generated QR code with logo and inserted data into MongoDB")
}


