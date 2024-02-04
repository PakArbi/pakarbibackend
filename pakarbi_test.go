package pakarbibackend

import (
	"fmt"
	// "os"
	// "context"
	// "io/ioutil"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "ULBICAMPUS"
	userdata.NPM = "1214000"
	userdata.Password = "ulbipass"
	userdata.PasswordHash = "ulbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
	userdata.Role = "user"
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	CreateNewUserRole(mconn, "user", userdata)
}

// func TestDeleteUser(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "pasabarapk")
// 	var userdata User
// 	userdata.Email = "1214006@std.ulbi.ac.id"
// 	DeleteUser(mconn, "user", userdata)
// }

func CreateNewUserToken(t *testing.T) {
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "ULBICAMPUS"
	userdata.NPM = "1214000"
	userdata.Password = "ulbipass"
	userdata.PasswordHash = "ulbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
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
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "ULBICAMPUS"
	userdata.NPM = "1214000"
	userdata.Password = "ulbipass"
	userdata.PasswordHash = "ulbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
	userdata.Role = "user"
	CreateNewUserRole(mconn, "user", userdata)
}

func TestParkiran(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var parkirandata Parkiran
	parkirandata.Parkiranid = "D41214000"
	parkirandata.Nama = "ULBI"
	parkirandata.NPM = "1214000"
	parkirandata.Prodi = "D4 Teknik Informatika"
	parkirandata.NamaKendaraan = "ULBI BUS"
	parkirandata.NomorKendaraan = "D 1234 NR"
	parkirandata.JenisKendaraan = "Mobil"
	CreateNewParkiran(mconn, "parkiran", parkirandata)
}

func TestAllParkiran(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	parkiran := GetAllParkiran(mconn, "parkiran")
	fmt.Println(parkiran)
}

func TestGeneratePasswordHash(t *testing.T) {
	passwordhash := "ulbipass"
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
	hasil, err := watoken.Encode("ulbipass", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.NPM = "1214000"
	userdata.PasswordHash = "ulbipass"

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
	userdata.NPM = "1214000"
	userdata.PasswordHash = "ulbipass"

	anu := IsPasswordValidNPM(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "PakArbiApp")
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "ULBICAMPUS"
	userdata.NPM = "1214000"
	userdata.Password = "ulbipass"
	userdata.PasswordHash = "ulbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
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
	userdata.NPM = "1214000"
	userdata.PasswordHash = "ulbipass"
	IsPasswordValidNPM(mconn, "user", userdata)
	fmt.Println(userdata)
}

// func TestInsertQRCodeDataToMongoDB(t *testing.T) {
// 	// Set up your MongoDB connection
// 	mconn := SetConnection("MONGOSTRING", "PakArbiApp")

// 	// Set up a sample Parkiran struct for testing
// 	dataparkiran := Parkiran{
// 		Parkiranid:     "D41214000",
// 		Nama:           "ULBICAMPUS",
// 		NPM:            "1214000",
// 		Prodi:          "D4 Teknik Informatika",
// 		NamaKendaraan:  "ULBIBUS",
// 		NomorKendaraan: "D 1234 ULBI",
// 		JenisKendaraan: "Mobil",
// 		Status:         "Mahasiswa Aktif",
// 		JamMasuk:       "08:00",
// 		JamKeluar:      "15:00",
// 	}

// 	// Generate QR code with logo and insert into MongoDB
// 	fileName, err := GenerateQRCodeBase64(mconn, "parkiran", dataparkiran)
// 	if err != nil {
// 		t.Errorf("Error generating QR code with logo: %v", err)
// 		return
// 	}

// 	// Read the QR code file
// 	_, errRead := ioutil.ReadFile(fileName)
// 	if errRead != nil {
// 		t.Errorf("Error reading QR code file: %v", errRead)
// 		return
// 	}

// 	// Check if the data is inserted into MongoDB correctly
// 	result := mconn.Collection("parkiran").FindOne(context.TODO(), bson.M{"parkiranid": dataparkiran.Parkiranid})
// 	if result.Err() != nil {
// 		t.Errorf("Failed to find inserted data in MongoDB: %v", result.Err())
// 		return
// 	}

// 	t.Log("Successfully generated QR code with logo, inserted data into MongoDB, and verified QR code file existence.")
// }

// func TestGenerateQRCodeBase64WithoutLogo(t *testing.T) {
// 	// Set up MongoDB connection
// 	mconn := SetConnection("Mongostring", "parkabi")
// 	collparkiran := "parkiran"

// 	// Create a sample Parkiran object
// 	dataparkiran := Parkiran{
// 		Parkiranid:     "123",
// 		Nama:           "John Doe",
// 		NPM:            "123456",
// 		Prodi:          "Computer Science",
// 		NamaKendaraan:  "Car",
// 		NomorKendaraan: "ABC123",
// 		JenisKendaraan: "Sedan",
// 		JamMasuk:       "09:00 AM",
// 		JamKeluar:      "05:00 PM",
// 		Status:         "Parked",
// 	}

// 	// Call the function to generate QR code and update MongoDB
// 	qrBase64, err := GenerateQRCodeBase64WithoutLogo(dataparkiran, mconn, collparkiran)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	// Print the generated QR code base64
// 	fmt.Println("Generated QR Code Base64:", qrBase64)
// }
