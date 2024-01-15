package pakarbibackend

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/base64"
	"unicode/utf8"
	"io/ioutil"
	"fmt"
	"strings"
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	// "github.com/huimingz/mongo-tools/common/json"
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

//functione
func ImageToBase64(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	// Encode the byte slice to base64
	base64Image := base64.StdEncoding.EncodeToString(fileContent)

	return base64Image, nil
}

// // InsertQRCodeDataToMongoDB inserts QR code data into MongoDB
// func InsertQRCodeDataToMongoDB(mconn *mongo.Database, collectionName, parkiranID string, qrCodeData []byte) error {
//     // Convert the base64-encoded string to []byte
//     _, err := base64.StdEncoding.DecodeString(string(qrCodeData))
//     if err != nil {
//         return fmt.Errorf("failed to decode base64: %v", err)
//     }

//     // Your implementation here to insert qrCodeData into MongoDB
//     // You can use the provided MongoDB connection (mconn) and collectionName to perform the insertion

//     // Encode QR code data as Base64
//     qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeData)

//     // For example:
//     collection := mconn.Collection(collectionName)

//     _, err = collection.InsertOne(context.TODO(), bson.M{
//         "parkiranID": parkiranID,
//         "qrCodeData": qrCodeBase64,
//     })

//     return err
// }

// InsertQRCodeDataToMongoDB inserts QR code data into MongoDB
func InsertQRCodeDataToMongoDB(mconn *mongo.Database, collectionName, parkiranID string, qrCodeData []byte) error {
    // Encode QR code data as Base64
    qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeData)

    // Decode the Base64 data to ensure correctness
    _, err := base64.StdEncoding.DecodeString(qrCodeBase64)
    if err != nil {
        return fmt.Errorf("failed to decode Base64 data: %v", err)
    }

    // Your implementation here to insert qrCodeData into MongoDB
    // You can use the provided MongoDB connection (mconn) and collectionName to perform the insertion

    collection := mconn.Collection(collectionName)

    _, err = collection.InsertOne(context.TODO(), bson.M{
        "parkiranID": parkiranID,
        "qrCodeData": qrCodeBase64,
    })

    return err
}

func replaceInvalidUTF8(input string) string {
	result := make([]rune, 0, len(input))
	for _, r := range input {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(string(r))
			if size == 1 {
				continue
			}
		}
		result = append(result, r)
	}
	return string(result)
}

// GenerateQRCodeLogoBase64 generates QR code with logo and inserts data into MongoDB
func GenerateQRCodeLogoBase64(mconn *mongo.Database, collparkiran string, dataparkiran Parkiran) (string, error) {

	// Convert struct to JSON
	dataJSON, err := json.Marshal(dataparkiran)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Replace invalid UTF-8 characters
	dataJSON = []byte(replaceInvalidUTF8(string(dataJSON)))

	// Ensure the dataJSON is valid UTF-8 encoded
	if !utf8.Valid(dataJSON) {
		return "", fmt.Errorf("data contains invalid UTF-8 characters")
	}

	// Insert Parkiran data into MongoDB
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

	// Generate QR code
	qrCode, err := qrcode.Encode(string(dataJSON), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Open the ULBI logo file from the "qrcode" folder
	logoFilePath := filepath.Join("qrcode", "logo_ulbi.png")
	logoBase64, err := ImageToBase64(logoFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to convert logo to base64: %v", err)
	}

	// Decode the ULBI logo
	logo, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(logoBase64)))
	if err != nil {
		return "", fmt.Errorf("failed to decode logo image: %v", err)
	}

	// Resize the logo to fit within the QR code
	resizedLogo := resize.Resize(80, 0, logo, resize.Lanczos3)

	// Create an image from the QR code
	qrImage, err := imaging.Decode(bytes.NewReader(qrCode))
	if err != nil {
		return "", fmt.Errorf("failed to decode QR code image: %v", err)
	}

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

	// Encode the final image into base64
	var base64String string
	base64Writer := &bytes.Buffer{} // Use a buffer instead of base64.NewEncoder
	err = imaging.Encode(base64Writer, result, imaging.PNG)
	if err != nil {
	return "", fmt.Errorf("failed to encode image: %v", err)
	}

	// Convert base64 data to JSON string
	base64Data := base64Writer.Bytes()
	jsonString, err := json.Marshal(map[string]interface{}{"base64Image": base64.StdEncoding.EncodeToString(base64Data)})
	if err != nil {
	return "", fmt.Errorf("failed to convert base64 to JSON string: %v", err)
	}

	// Use jsonString as needed, for example, inserting into MongoDB
	err = InsertQRCodeDataToMongoDB(mconn, "qrcodes", dataparkiran.Parkiranid, []byte(jsonString))
	if err != nil {
	return "", fmt.Errorf("failed to insert QR code data to MongoDB: %v", err)
	}

// Get the base64 representation as a string
base64String = strings.TrimSpace(string(jsonString))

	// Insert QR code data into MongoDB
	err = InsertQRCodeDataToMongoDB(mconn, "qrcodes", dataparkiran.Parkiranid, []byte(base64String))
	if err != nil {
		return "", fmt.Errorf("failed to insert QR code data to MongoDB: %v", err)
	}

	// Update Parkiran data with Base64 image
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "data.base64Image", Value: base64String}}}}
	_, err = mconn.Collection(collparkiran).UpdateOne(context.TODO(), bson.M{"parkiranid": dataparkiran.Parkiranid}, update)
	if err != nil {
		return "", fmt.Errorf("failed to update Parkiran data with Base64 image: %v", err)
	}

	return fileName, nil
}




// // GenerateQRCodeLogoBase64 generates QR code with logo and inserts data into MongoDB
// func GenerateQRCodeLogoBase64(mconn *mongo.Database, collparkiran string, dataparkiran Parkiran) (string, error) {
// 	// Insert Parkiran data into MongoDB
// 	insertParkiran(mconn, collparkiran, dataparkiran)

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

// 	// Open the ULBI logo file from the "qrcode" folder
// 	logoFilePath := filepath.Join("qrcode", "logo_ulbi.png")
// 	logoBase64, err := ImageToBase64(logoFilePath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to convert logo to base64: %v", err)
// 	}

// 	// Insert QR code data into MongoDB
// 	err = InsertQRCodeDataToMongoDB(mconn, "qrcodes", dataparkiran.Parkiranid, qrCode)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to insert QR code data to MongoDB: %v", err)
// 	}

// 	// Decode the ULBI logo
// 	logo, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(logoBase64)))
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

// 	// Encode QR code data as Base64
//     qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

//     // Save the final QR code data into MongoDB
//     err = InsertQRCodeDataToMongoDB(mconn, "qrcodes", dataparkiran.Parkiranid, []byte(qrCodeBase64))
//     if err != nil {
//         return "", fmt.Errorf("failed to insert QR code data to MongoDB: %v", err)
//     }

// 	// Encode the final image into base64
// 	var base64String string
// 	base64Writer := &bytes.Buffer{} // Use a buffer instead of base64.NewEncoder
// 	err = imaging.Encode(base64Writer, result, imaging.PNG)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to encode image: %v", err)
// 	}

// 	// Get the base64 representation as a string
// 	base64String = strings.TrimSpace(base64Writer.String())

// 	// Insert data into MongoDB collection along with base64 image data
// 	dataparkiran.Base64Image = base64String
// 	insertParkiran(mconn, collparkiran, dataparkiran) // Use insertParkiran2 for Parkiran2 type

// 	return fileName, nil
// }


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
