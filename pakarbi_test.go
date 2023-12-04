package pakarbibackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "pakarbi"
	userdata.NPM = "1214000"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
	userdata.Role = "user"
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	CreateNewUserRole(mconn, "user", userdata)
}

// func TestDeleteUser(t *testing.T) {
// 	mconn := SetConnection("MONGOSTRING", "pasabarapk")
// 	var userdata User
// 	userdata.Email = "1214000@std.ulbi.ac.id"
// 	DeleteUser(mconn, "user", userdata)
// }

func CreateNewUserToken(t *testing.T) {
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "pakarbi"
	userdata.NPM = "1214000"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
	userdata.Role = "user"

	// Create a MongoDB connection
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")

	// Call the function to create a admin and generate a token
	err := CreateUserAndAddToken("your_private_key_env", mconn, "user", userdata)

	if err != nil {
		t.Errorf("Error creating user and token: %v", err)
	}
}

func TestGFCPostHandlerUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "pakarbi"
	userdata.NPM = "1214000"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
	userdata.Role = "user"
	CreateNewUserRole(mconn, "user", userdata)
}

func TestParkiran(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	var parkirandata Parkiran
	parkirandata.ParkiranId = "D4TI150CC"
	parkirandata.Nama = "Farhan Rizki Maulana"
	parkirandata.NPM = "1214020"
	parkirandata.Prodi = "D4 Teknik Informatika"
	parkirandata.NamaKendaraan = "Supra X 125"
	parkirandata.NomorKendaraan = "F 1234 NR"
	parkirandata.JenisKendaraan = "Motor"
	CreateNewParkiran(mconn, "parkiran", parkirandata)
}

func TestAllParkiran(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
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
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "pakarbi"
	userdata.PasswordHash = "pakarbipass"

	filter := bson.M{"usernameid": userdata.UsernameId}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPass(userdata.PasswordHash)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "pakarbi"
	userdata.PasswordHash = "pakarbipass"

	anu := IsPasswordValidNPM(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	var userdata User
	userdata.UsernameId = "D4TI1214000"
	userdata.Username = "pakarbi"
	userdata.NPM = "1214000"
	userdata.Password = "pakarbipass"
	userdata.PasswordHash = "pakarbipass"
	userdata.Email = "1214000@std.ulbi.ac.id"
	userdata.Role = "user"
	CreateUser(mconn, "user", userdata)
}

func TestAdminFix(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "pakarbiappdb")
	var admindata Admin
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