package pakarbibackend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// <---Function Generate Code QR--->
// func GenerateQRCodeWithLogo(mconn *mongo.Database, dataparkiran Parkiran) (string, error) {
// 	// Convert struct to JSON
// 	dataJSON, err := json.Marshal(dataparkiran)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal JSON: %v", err)
// 	}

// 	// Generate QR code
// 	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate QR code: %v", err)
// 	}

// 	// Create an image from the QR code
// 	qrImage, err := imaging.Decode(bytes.NewReader(qrCode))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode QR code image: %v", err)
// 	}

// 	// Open the ULBI logo file
// 	logoFile, err := os.Open("logo_ulbi.png") // Replace with your ULBI logo file path
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open logo file: %v", err)
// 	}
// 	defer logoFile.Close()

// 	// Decode the ULBI logo
// 	logo, _, err := image.Decode(logoFile)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode logo image: %v", err)
// 	}

// 	// Resize the logo to fit within the QR code
// 	resizedLogo := imaging.Resize(logo, 80, 0, imaging.Lanczos)

// 	// Calculate position to overlay the logo on the QR code
// 	x := (qrImage.Bounds().Dx() - resizedLogo.Bounds().Dx()) / 2
// 	y := (qrImage.Bounds().Dy() - resizedLogo.Bounds().Dy()) / 2

// 	// Draw the logo onto the QR code
// 	result := imaging.Overlay(qrImage, resizedLogo, image.Pt(x, y), 1.0)

// 	// Save the final QR code with logo
// 	fileName := dataparkiran.Parkiranid + "_qrcode.png" // Using Parkiran ID in the file name
// 	outFile, err := os.Create(fileName)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create output file: %v", err)
// 	}
// 	defer outFile.Close()

// 	// Encode the final image into the output file
// 	err = imaging.Encode(outFile, result, imaging.PNG)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to encode image: %v", err)
// 	}

// 	return fileName, nil
// }

// func GenerateQRCodeWithLogo(mconn *mongo.Database, dataparkiran Parkiran) (string, error) {
// 	// Convert struct to JSON
// 	dataJSON, err := json.Marshal(dataparkiran)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal JSON: %v", err)
// 	}

// 	// Generate QR code
// 	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate QR code: %v", err)
// 	}

// 	// Create an image from the QR code
// 	qrImage, err := imaging.Decode(bytes.NewReader(qrCode))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode QR code image: %v", err)
// 	}

// 	// // Open the ULBI logo file from the "qrcode" folder
// 	// logoFile, err := os.Open("qrcode/logo_ulbi.png") // Replace with your ULBI logo file path
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to open logo file: %v", err)
// 	// }
// 	// defer logoFile.Close()

// 	// Open the ULBI logo file from the "qrcode" folder
// 	logoFilePath := filepath.Join("qrcode", "logo_ulbi.png")
// 	logoFile, err := os.Open(logoFilePath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open logo file: %v", err)
// 	}
// 	defer logoFile.Close()

// 	// Get the base name of the ULBI logo file
// 	// logoBaseName := path.Base(logoFile.Name())

// 	// Decode the ULBI logo
// 	logo, _, err := image.Decode(logoFile)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode logo image: %v", err)
// 	}

// 	// Resize the logo to fit within the QR code
// 	resizedLogo := resize.Resize(80, 0, logo, resize.Lanczos3)

// 	// Calculate position to overlay the logo on the QR code
// 	x := (qrImage.Bounds().Dx() - resizedLogo.Bounds().Dx()) / 2
// 	y := (qrImage.Bounds().Dy() - resizedLogo.Bounds().Dy()) / 2

// 	// Draw the logo onto the QR code
// 	result := imaging.Overlay(qrImage, resizedLogo, image.Pt(x, y), 1.0)

// 	// Save the final QR code with logo
// 	// Save the final QR code with logo
// 	fileName := filepath.Join("qrcode", fmt.Sprintf("%s_logo_ulbi_qrcode.png", dataparkiran.Parkiranid))
// 	outFile, err := os.Create(fileName)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create output file: %v", err)
// 	}
// 	defer outFile.Close()

// 	// Encode the final image into the output file
// 	err = imaging.Encode(outFile, result, imaging.PNG)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to encode image: %v", err)
// 	}

// 	return fileName, nil
// }

func GenerateQRCodeWithLogo(mconn *mongo.Database, collparkiran string, dataparkiran Parkiran) (string, error) {
	// Convert struct to JSON
	dataJSON, err := json.Marshal(dataparkiran)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Generate QR code
	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Create an image from the QR code
	qrImage, err := imaging.Decode(bytes.NewReader(qrCode))
	if err != nil {
		return "", fmt.Errorf("failed to decode QR code image: %v", err)
	}

	// Open the ULBI logo file from the "qrcode" folder
	logoFilePath := filepath.Join("qrcode", "logo_ulbi.png")
	logoFile, err := os.Open(logoFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open logo file: %v", err)
	}
	defer logoFile.Close()

	// Decode the ULBI logo
	logo, _, err := image.Decode(logoFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode logo image: %v", err)
	}

	// Resize the logo to fit within the QR code
	resizedLogo := resize.Resize(80, 0, logo, resize.Lanczos3)

	// Calculate position to overlay the logo on the QR code
	x := (qrImage.Bounds().Dx() - resizedLogo.Bounds().Dx()) / 2
	y := (qrImage.Bounds().Dy() - resizedLogo.Bounds().Dy()) / 2

	// Draw the logo onto the QR code
	result := imaging.Overlay(qrImage, resizedLogo, image.Pt(x, y), 1.0)

	// Save the final QR code with logo
	fileName := filepath.Join("qrcode", fmt.Sprintf("%s_logo_ulbi_qrcode.png", dataparkiran.Parkiranid))
	outFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the final image into the output file
	err = imaging.Encode(outFile, result, imaging.PNG)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	// Insert data into MongoDB collection
	insertParkiran(mconn, collparkiran, dataparkiran)

	return fileName, nil
}

func GenerateQRCodeWithLogoULBI(mconn *mongo.Database, collparkiran string, dataparkiran Parkiran) (string, error) {
	// Convert struct to JSON
	dataJSON, err := json.Marshal(dataparkiran)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Generate QR code
	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Create an image from the QR code
	qrImage, err := imaging.Decode(bytes.NewReader(qrCode))
	if err != nil {
		return "", fmt.Errorf("failed to decode QR code image: %v", err)
	}

	// Open the ULBI logo file from the project root directory
	logoFilePath := "logo_ulbi.png"
	logoFile, err := os.Open(logoFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open logo file: %v", err)
	}
	defer logoFile.Close()

	// Decode the ULBI logo
	logo, _, err := image.Decode(logoFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode logo image: %v", err)
	}

	// Resize the logo to fit within the QR code
	resizedLogo := resize.Resize(80, 0, logo, resize.Lanczos3)

	// Calculate position to overlay the logo on the QR code
	x := (qrImage.Bounds().Dx() - resizedLogo.Bounds().Dx()) / 2
	y := (qrImage.Bounds().Dy() - resizedLogo.Bounds().Dy()) / 2

	// Draw the logo onto the QR code
	result := imaging.Overlay(qrImage, resizedLogo, image.Pt(x, y), 1.0)

	// Save the final QR code with logo
	fileName := filepath.Join("qrcode", fmt.Sprintf("%s_logo_ulbi_qrcode.png", dataparkiran.Parkiranid))
	outFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the final image into the output file
	err = imaging.Encode(outFile, result, imaging.PNG)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	// Insert data into MongoDB collection
	insertParkiran(mconn, collparkiran, dataparkiran)

	return fileName, nil
}


// PathQRCode menyimpan path untuk folder QR code.
const PathQRCode = "C:\\Users\\ACER\\Documents\\pakarbibackend\\qrcode"

// InitQRCodeFolder membuat folder QR code jika belum ada.
func InitQRCodeFolder() error {
	if _, err := os.Stat(PathQRCode); os.IsNotExist(err) {
		log.Printf("Creating QR code folder at: %s", PathQRCode)
		err := os.Mkdir(PathQRCode, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		log.Printf("QR code folder already exists at: %s", PathQRCode)
	}
	return nil
}

// GenerateQRCode menghasilkan QR code dari data parkiran dan menyimpannya di path yang ditentukan.
func GenerateQRCode(DataParkir Parkiran, outputFilePath string) error {
	// Convert struct to JSON
	dataJSON, err := json.Marshal(DataParkir)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Generate QR code
	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Decode the QR code image
	qrImage, _, err := image.Decode(bytes.NewReader(qrCode))
	if err != nil {
		return fmt.Errorf("failed to decode QR code image: %v", err)
	}

	// Create the output directory if it doesn't exist
	outputDir := filepath.Dir(outputFilePath)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create the output file
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the final image into the output file
	err = imaging.Encode(outFile, qrImage, imaging.PNG)
	if err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}

	log.Printf("QR code generated successfully and saved to %s", outputFilePath)
	return nil
}

func GenerateQRCodeULBI(dataParkir Parkiran) (string, error) {
	// Convert struct to JSON
	dataJSON, err := json.Marshal(dataParkir)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Generate QR code
	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Create the output file
	qrOutputPath := "C:\\Users\\ACER\\Documents\\pakarbibackend\\qrcode\\" + dataParkir.Parkiranid + "_qrcode.png"
	outFile, err := os.Create(qrOutputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Write the QR code to the file
	_, err = outFile.Write(qrCode)
	if err != nil {
		return "", fmt.Errorf("failed to write QR code to file: %v", err)
	}

	return qrOutputPath, nil
}

// <--- FUNCTION CRUD --->
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (insertedID interface{}) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// <--- FUNCTION USER --->

func CreateNewUserRole(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(userdata.PasswordHash)
	if err != nil {
		return err
	}
	userdata.PasswordHash = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

func CreateUserAndAddToken(privateKeyEnv string, mongoconn *mongo.Database, collection string, userdata User) error {
	// Hash the password before storing it
	hashedPassword, err := HashPass(userdata.PasswordHash)
	if err != nil {
		return err
	}
	userdata.PasswordHash = hashedPassword

	// Create a token for the admin
	tokenstring, err := watoken.Encode(userdata.Email, os.Getenv(privateKeyEnv))
	if err != nil {
		return err
	}

	userdata.Token = tokenstring

	// Insert the admin data into the MongoDB collection
	if err := atdb.InsertOneDoc(mongoconn, collection, userdata.Email); err != nil {
		return nil // Mengembalikan kesalahan yang dikembalikan oleh atdb.InsertOneDoc
	}

	// Return nil to indicate success
	return nil
}

func CreateStatus(status string, message string, data interface{}, requestParkiran RequestParkiran) Status2 {
	return Status2{
		Status:          status,
		Message:         message,
		DataParkir:      data,
		RequestParkiran: requestParkiran,
	}
}

func CreateResponse(status bool, message string, data interface{}) Response {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return response
}

func CreateUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(userdata.PasswordHash)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	userid := userdata.Email
	tokenstring, err := watoken.Encode(userid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get adminid
	useridstring := watoken.DecodeGetId(publicKey, tokenstring)
	if useridstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(useridstring)
	userdata.Private = privateKey
	userdata.Public = publicKey
	userdata.PasswordHash = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

func UpdatedUser(mongoconn *mongo.Database, collection string, filter bson.M, userdata User) interface{} {
	updatedFilter := bson.M{"usernameid": userdata.UsernameId}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, userdata)
}

func GetAllUser(mongoconn *mongo.Database, collection string) []User {
	user := atdb.GetAllDoc[[]User](mongoconn, collection)
	return user
}

// <--- FUNCTION ADMIN --->

func CreateNewAdminRole(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.PasswordHash)
	if err != nil {
		return err
	}
	admindata.PasswordHash = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

func CreateAdminAndAddToken(privateKeyEnv string, mongoconn *mongo.Database, collection string, admindata User) error {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.PasswordHash)
	if err != nil {
		return err
	}
	admindata.PasswordHash = hashedPassword

	// Create a token for the admin
	tokenstring, err := watoken.Encode(admindata.Email, os.Getenv(privateKeyEnv))
	if err != nil {
		return err
	}

	admindata.Token = tokenstring

	// Insert the admin data into the MongoDB collection
	if err := atdb.InsertOneDoc(mongoconn, collection, admindata.Email); err != nil {
		return nil // Mengembalikan kesalahan yang dikembalikan oleh atdb.InsertOneDoc
	}

	// Return nil to indicate success
	return nil
}

func CreateAdmin(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPass(admindata.PasswordHash)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	adminid := admindata.Email
	tokenstring, err := watoken.Encode(adminid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get adminid
	adminidstring := watoken.DecodeGetId(publicKey, tokenstring)
	if adminidstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(adminidstring)
	admindata.Private = privateKey
	admindata.Public = publicKey
	admindata.PasswordHash = hashedPassword

	// Insert the admin data into the database
	return atdb.InsertOneDoc(mongoconn, collection, admindata)
}

// function mekanisme untuk auto-increment
func SequenceAutoIncrement(mongoconn *mongo.Database, sequenceName string) int {
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result struct {
		Seq int `bson:"seq"`
	}

	after := options.After
	opt := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	collection := mongoconn.Collection("counters")
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&result)
	if err != nil {
		// handle error
	}
	return result.Seq
}

// <---FUNCTION GENERATE FOR PARKIRANID --->
func GenerateParkiranID(npm string) string {
	// Contoh: Jika NPM adalah '1214000'. maka yang diambil '4000'
	// Anda dapat menggunakan beberapa digit terakhir dari NPM
	// Misalnya, mengambil 4 digit terakhir (atau lebih sesuai kebutuhan)
	lastDigits := npm[len(npm)-4:] // Mengambil 4 digit terakhir dari NPM
	return "D3/D4" + lastDigits    // Menggabungkan pola dengan digit terakhir dari NPM
}

// <--- FUNCTION CRUD PARKIRAN --->

// parkiran
func CreateNewParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, parkirandata)
}

// parkiran function
func insertParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
	return atdb.InsertOneDoc(mongoconn, collection, parkirandata)
}

func DeleteParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
	filter := bson.M{"parkiranid": parkirandata.Parkiranid}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func UpdatedParkiran(mongoconn *mongo.Database, collection string, filter bson.M, parkirandata Parkiran) interface{} {
	updatedFilter := bson.M{"parkiranid": parkirandata.Parkiranid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, parkirandata)
}

func GetAllParkiran(mongoconn *mongo.Database, collection string) []Parkiran {
	parkiran := atdb.GetAllDoc[[]Parkiran](mongoconn, collection)
	return parkiran
}

func GetOneParkiran(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
	filter := bson.M{"parkiranid": parkirandata.Parkiranid}
	return atdb.GetOneDoc[Parkiran](mongoconn, collection, filter)
}

func GetAllParkiranID(mongoconn *mongo.Database, collection string, parkirandata Parkiran) Parkiran {
	filter := bson.M{
		"parkiranid":     parkirandata.Parkiranid,
		"nama":           parkirandata.Nama,
		"npm":            parkirandata.NPM,
		"prodi":          parkirandata.Prodi,
		"namakendaraan":  parkirandata.NamaKendaraan,
		"nomorkendaraan": parkirandata.NomorKendaraan,
		"jeniskendaraan": parkirandata.JenisKendaraan,
	}
	parkiranID := atdb.GetOneDoc[Parkiran](mongoconn, collection, filter)
	return parkiranID
}
