package pakarbibackend

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode/qr"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// <---Function Generate Code QR--->
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

// functione
func GenerateQRCodeBase64(mconn *mongo.Database, collparkiran string, dataparkiran Parkiran) (string, error) {
	// Generate QR code
	qrCode, err := generateQRCode(dataparkiran.Parkiranid)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Resize the QR code image
	qrCode = resize.Resize(256, 256, qrCode, resize.Lanczos3)

	// Convert the QR code image to base64
	qrBase64, err := imageToBase64(qrCode)
	if err != nil {
		return "", fmt.Errorf("failed to convert QR code image to base64: %v", err)
	}

	// Update data Parkiran dengan gambar Base64
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "base64Image", Value: qrBase64},
	}}}
	_, err = mconn.Collection(collparkiran).UpdateOne(context.TODO(), bson.M{"parkiranid": dataparkiran.Parkiranid}, update)
	if err != nil {
		return "", fmt.Errorf("failed to update Parkiran data with base64 image: %v", err)
	}

	return "", nil
}

// generateQRCode generates a QR code for the given data.
func generateQRCode(data string) (image.Image, error) {
	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

// imageToBase64 converts an image to base64.
func imageToBase64(img image.Image) (string, error) {
	var imgData []byte
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		return "", err
	}
	imgData = buffer.Bytes()

	base64Str := base64.StdEncoding.EncodeToString(imgData)
	return base64Str, nil
}

// func GenerateQRCodeBase64(mconn *mongo.Database, collparkiran string, dataparkiran Parkiran) (string, error) {
// 	// Convert struct to JSON
// 	dataJSON, err := json.Marshal(dataparkiran)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal JSON: %v", err)
// 	}

// 	// Encode JSON as base64
// 	encodedJSON := base64.StdEncoding.EncodeToString(dataJSON)

//     // Ensure the dataJSON is valid UTF-8 encoded
//     if !utf8.Valid(dataJSON) {
//         return "", fmt.Errorf("data contains invalid UTF-8 characters")
//     }

//      // Generate QR code
// 	 qrCode, err := qrcode.Encode(encodedJSON, qrcode.Medium, 256)
// 	 if err != nil {
// 		 return "", fmt.Errorf("failed to generate QR code: %v", err)
// 	 }

//     // Create an image from the QR code
//     qrImage, err := imaging.Decode(bytes.NewReader(qrCode))
//     if err != nil {
//         return "", fmt.Errorf("failed to decode QR code image: %v", err)
//     }

//     // Convert the final image to base64
//     finalImageBase64, err := ImageToBase64(fileName)
//     if err != nil {
//         return "", fmt.Errorf("failed to convert final image to base64: %v", err)
//     }

//     // Insert Parkiran data into MongoDB
//     err = InsertParkiran(mconn, collparkiran, dataparkiran)
//     if err != nil {
//         return "", fmt.Errorf("failed to insert Parkiran data to MongoDB: %v", err)
//     }

//     // Update data Parkiran dengan gambar Base64
//     update := bson.D{{Key: "$set", Value: bson.D{
//         {Key: "base64Image", Value: finalImageBase64},
//     }}}
//     _, err = mconn.Collection(collparkiran).UpdateOne(context.TODO(), bson.M{"parkiranid": dataparkiran.Parkiranid}, update)
//     if err != nil {
//         return "", fmt.Errorf("gagal memperbarui data Parkiran dengan gambar Base64: %v", err)
//     }

//     return fileName, nil
// }

// <---FUNCTION GENERATE FOR PARKIRANID --->
// Ambil npm 2 belakang.
func GetLastDigitsNPM(npm string) string {
	// Assuming you want the last 4 digits, adjust accordingly if needed
	lastDigits := npm[len(npm)-4:]
	return lastDigits
}

func GenerateParkiranID(npm string, option string) (string, error) {
	lastDigits := GetLastDigitsNPM(npm)
	// Validate the option
	if option != "D3" && option != "D4" {
		return "", errors.New("Invalid option. Use 'D3' or 'D4'")
	}
	return option + lastDigits, nil
}

func ImageToBase64(imagePath string) (string, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer imageFile.Close()

	image, _, err := image.Decode(imageFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, image)
	if err != nil {
		return "", fmt.Errorf("failed to encode image to PNG: %v", err)
	}

	// Convert the buffer to a base64 string
	base64String := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return base64String, nil
}

// fungsi GenerateImageFromBase64 untuk mengembalikan nilai fileName
func GenerateImageFromBase64(base64Data string, fileName string) (string, error) {
	// Decode the base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image data: %v", err)
	}

	// Create the image file
	imageFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %v", err)
	}
	defer imageFile.Close()

	// Write the decoded data to the file
	_, err = imageFile.Write(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to write image data to file: %v", err)
	}

	return fileName, nil
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

func GetAllParkirans(mongoconn *mongo.Database, collection string) []Parkiran {
	parkiran := atdb.GetAllDoc[[]Parkiran](mongoconn, collection)
	return parkiran
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

func DeleteQRCodeData(mconn *mongo.Database, collparkiran, parkiranID string) error {
	collection := mconn.Collection(collparkiran)

	// Delete QR code data based on Parkiranid
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"parkiranid": parkiranID},
		bson.D{{Key: "$unset", Value: bson.D{{Key: "base64Image"}}}},
	)
	if err != nil {
		return fmt.Errorf("failed to delete QR code data: %v", err)
	}

	return nil
}

// DeleteQRCodeData menghapus data QR code dari MongoDB
func DeleteQRCodeData2(mconn *mongo.Database, collparkiran, parkiranID string) error {
	// Hapus data parkiran dari MongoDB
	_, err := mconn.Collection(collparkiran).DeleteOne(context.TODO(), bson.M{"parkiranid": parkiranID})
	if err != nil {
		return fmt.Errorf("failed to delete parkiran data: %v", err)
	}

	// Implementasikan logika penghapusan QR code jika diperlukan

	return nil
}

func UpdateParkiran(mconn *mongo.Database, collparkiran string, newData Parkiran) error {
	update := bson.D{{Key: "$set", Value: newData}}
	_, err := mconn.Collection(collparkiran).UpdateOne(context.TODO(), bson.M{"parkiranid": newData.Parkiranid}, update)
	return err
}

func UpdatedParkiran(mongoconn *mongo.Database, collection string, filter bson.M, parkirandata Parkiran) interface{} {
	updatedFilter := bson.M{"parkiranid": parkirandata.Parkiranid}
	return atdb.ReplaceOneDoc(mongoconn, collection, updatedFilter, parkirandata)
}

func GetAllParkiran(mongoconn *mongo.Database, collection string) []Parkiran {
	parkiran := atdb.GetAllDoc[[]Parkiran](mongoconn, collection)
	return parkiran
}

func GetOneData(mongoconn *mongo.Database, collection string, parkirandata Parkiran) interface{} {
	filter := bson.M{"parkiranid": parkirandata.Parkiranid}
	return atdb.GetOneDoc[Parkiran](mongoconn, collection, filter)
}

// GetOneDataParkiranByID mengambil satu data parkiran dari MongoDB berdasarkan parkiranID
func GetOneDataParkiranByID(mconn *mongo.Database, collparkiran, parkiranID string) (*Parkiran, error) {
	var result Parkiran

	err := mconn.Collection(collparkiran).FindOne(context.TODO(), bson.M{"parkiranid": parkiranID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("parkiran data not found")
		}
		return nil, fmt.Errorf("failed to get parkiran data: %v", err)
	}

	return &result, nil
}

func GetParkiranFromID(db *mongo.Database, col string, _id primitive.ObjectID) (*Parkiran, error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}

	reportlist := new(Parkiran)

	err := cols.FindOne(context.Background(), filter).Decode(reportlist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no data found for ID %s", _id.Hex())
		}
		return nil, fmt.Errorf("error retrieving data for ID %s: %s", _id.Hex(), err.Error())
	}

	return reportlist, nil
}

// Percobaan ganti alurnya untuk Generate code qr ketika sudah inputkan datanya
func GenerateQRCodeBase64WithoutLogo(dataparkiran Parkiran, mconn *mongo.Database, collparkiran string) (string, error) {
	// Construct the data to be encoded in the QR code
	qrCodeData := fmt.Sprintf(
		// "Parkiran ID: %s\nNama: %s\nNPM: %s\nProdi: %s\nNama Kendaraan: %s\nNomor Kendaraan: %s\nJenis Kendaraan: %s\nJam Masuk: %s\nJam Keluar: %s\nStatus: %s",
		dataparkiran.Parkiranid,
		dataparkiran.Nama,
		dataparkiran.NPM,
		dataparkiran.Prodi,
		dataparkiran.NamaKendaraan,
		dataparkiran.NomorKendaraan,
		dataparkiran.JenisKendaraan,
		dataparkiran.JamMasuk,
		dataparkiran.JamKeluar,
		dataparkiran.Status,
	)

	// Generate QR code
	qrCode, err := generateQRCode(qrCodeData)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Resize the QR code image
	qrCode = resize.Resize(1080, 1080, qrCode, resize.Lanczos3)

	// Convert the QR code image to base64
	qrBase64, err := imageToBase64(qrCode)
	if err != nil {
		return "", fmt.Errorf("failed to convert QR code image to base64: %v", err)
	}

	// Update data Parkiran dengan gambar Base64
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "base64Image", Value: qrBase64},
	}}}
	_, err = mconn.Collection(collparkiran).UpdateOne(context.TODO(), bson.M{"parkiranid": dataparkiran.Parkiranid}, update)
	if err != nil {
		return "", fmt.Errorf("failed to update Parkiran data with base64 image: %v", err)
	}

	return qrBase64, nil
}
