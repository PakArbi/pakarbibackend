package pakarbibackend

import (
	"context"
	"fmt"
	"os"
	"time"
	"encoding/base64"


	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// <--- FUNCTION MONGODB --->

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func MongoCreateConnection(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func IsExist(Tokenstr, PublicKey string) bool {
	id := watoken.DecodeGetId(PublicKey, Tokenstr)
	return id != ""
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

// func untuk update status auto
func UpdateStatusInMongoDB(mconn *mongo.Database, collname, parkiranID string, updatedStatus Status) error {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: updatedStatus}}}}
	_, err := mconn.Collection(collname).UpdateOne(context.TODO(), bson.M{"parkiranid": parkiranID}, update)
	return err
}

// <---FUNCTION UNTUK MENCARI DAN Mengambil GAMBAR CODE QR DI MONGODB --->
func GetQRCodeDataFromMongoDB(mconn *mongo.Database, collname, parkiranID string) (string, error) {
	result := mconn.Collection(collname).FindOne(context.TODO(), bson.M{"parkiranid": parkiranID})
	if result.Err() != nil {
		return "", fmt.Errorf("failed to find QR code data: %v", result.Err())
	}

	var data map[string]interface{}
	if err := result.Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode QR code data: %v", err)
	}

	base64Image, ok := data["base64Image"].(string)
	if !ok {
		return "", fmt.Errorf("base64Image not found or not a string")
	}

	return base64Image, nil
}
func ScanQRCode(mconn *mongo.Database, collparkiran, qrcodeData string) (ScanResult, error) {
	var result ScanResult

	// Decode QR code data
	qrCodeBytes, err := base64.StdEncoding.DecodeString(qrcodeData)
	if err != nil {
		return result, fmt.Errorf("failed to decode QR code data: %v", err)
	}

	// Convert QR code data to string
	qrCodeString := string(qrCodeBytes)

	// Find Parkiran entry by QR code data
	var parkiran Parkiran
	err = mconn.Collection(collparkiran).FindOne(context.TODO(), bson.M{"qrcode.base64image": qrCodeString}).Decode(&parkiran)
	if err != nil {
		return result, fmt.Errorf("failed to find Parkiran by QR code: %v", err)
	}

	// Check if the entry timestamp is available
	if parkiran.Status.DataParkir.WaktuMasuk != "" {
		waktuMasuk, err := time.Parse(time.RFC3339, parkiran.Status.DataParkir.WaktuMasuk)
		if err != nil {
			return result, fmt.Errorf("failed to parse entry timestamp: %v", err)
		}
		result.WaktuMasuk = waktuMasuk
	}

	// Update data Parkiran with timestamp for waktu keluar
	currentTime := time.Now()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status.dataparkir.waktukeluar", Value: currentTime.Format(time.RFC3339)},
		}},
	}
	_, err = mconn.Collection(collparkiran).UpdateOne(context.TODO(), bson.M{"parkiranid": parkiran.Parkiranid}, update)
	if err != nil {
		return result, fmt.Errorf("failed to update data Parkiran with timestamp for waktu keluar: %v", err)
	}

	result.Message = "QR code scanned successfully"
	result.WaktuKeluar = currentTime

	return result, nil
}

// <--- FUNCTION USER --->
func InsertUserdata(MongoConn *mongo.Database, usernameid, username, npm, password, passwordhash, email, role string) (InsertedID interface{}) {
	req := new(User)
	req.UsernameId = usernameid
	req.Username = username
	req.NPM = npm
	req.Password = password
	req.PasswordHash = passwordhash
	req.Email = email
	req.Role = role
	return InsertOneDoc(MongoConn, "user", req)
}

func DeleteUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	filter := bson.M{"npm": userdata.NPM}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func FindUserByField(mongoconn *mongo.Database, collection, searchField, searchValue string) User {
    filter := bson.M{searchField: searchValue}
    return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func FindUserNPM(mongoconn *mongo.Database, collection string, userdata User) User {
	filter := bson.M{"npm": userdata.NPM}
	return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func FindUserEmail(mongoconn *mongo.Database, collection string, userdata User) User {
	filter := bson.M{"email": userdata.Email}
	return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func IsPasswordValidNPM(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{
		"$or": []bson.M{
			{"npm": userdata.NPM},
			{"email": userdata.Email},
		},
	}

	var res User
	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)

	if err == nil {
		return CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
	}
	return false
}

func IsPasswordValidEmail(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{
		"$or": []bson.M{
			{"email": userdata.Email},
			{"npm": userdata.NPM},
		},
	}

	var res User
	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)

	if err == nil {
		return CheckPasswordHash(userdata.PasswordHash, res.PasswordHash)
	}
	return false
}

func GetOneUser(MongoConn *mongo.Database, colname string, userdata User) User {
	filter := bson.M{"npm": userdata.NPM}
	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
	return data
}

// <--- FUNCTION ADMIN --->
func InsertAdmindata(MongoConn *mongo.Database, usernameid, username, password, passwordhash, email, role string) (InsertedID interface{}) {
	req := new(Admin)
	req.UsernameId = usernameid
	req.Username = username
	req.Password = password
	req.PasswordHash = passwordhash
	req.Email = email
	req.Role = role
	return InsertOneDoc(MongoConn, "admin", req)
}

func IsPasswordValidEmailAdmin(mongoconn *mongo.Database, collection string, admindata Admin) bool {
	filter := bson.M{
		"$or": []bson.M{
			{"email": admindata.Email},
		},
	}

	var res Admin
	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)

	if err == nil {
		return CheckPasswordHash(admindata.PasswordHash, res.PasswordHash)
	}
	return false
}

func DeleteAdmin(mongoconn *mongo.Database, collection string, admindata Admin) interface{} {
	filter := bson.M{"email": admindata.Email}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func FindAdmin(mongoconn *mongo.Database, collection string, admindata Admin) Admin {
	filter := bson.M{"email": admindata.Email}
	return atdb.GetOneDoc[Admin](mongoconn, collection, filter)
}

func IsPasswordValidAdmin(mongoconn *mongo.Database, collection string, admindata Admin) bool {
	filter := bson.M{"email": admindata.Email}
	res := atdb.GetOneDoc[Admin](mongoconn, collection, filter)
	return CheckPasswordHash(admindata.PasswordHash, res.PasswordHash)
}

func GetOneAdmin(MongoConn *mongo.Database, colname string, admindata Admin) Admin {
	filter := bson.M{"email": admindata.Email}
	data := atdb.GetOneDoc[Admin](MongoConn, colname, filter)
	return data
}

func GetOneParkiranData(mongoconn *mongo.Database, colname, Pkrid string) (dest Parkiran) {
	filter := bson.M{"parkiranid": Pkrid}
	dest = atdb.GetOneDoc[Parkiran](mongoconn, colname, filter)
	return
}

func GetParkiranById(mconn *mongo.Database, collectionname, parkiranID string) (Parkiran, error) {
	collection := mconn.Collection(collectionname)

	var result Parkiran
	filter := bson.D{{Key: "parkiranid", Value: parkiranID}}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return Parkiran{}, err
	}

	return result, nil
}

// function Parkiran
func InsertDataParkir(MongoConn *mongo.Database, npm string, nama, prodi, namaKendaraan, nomorKendaraan, jenisKendaraan, statusMessage, waktuMasuk, waktuKeluar string) (InsertedID interface{}) {
	// Generate Parkiran ID
	parkiranID := GenerateParkiranID(npm)
	req := Parkiran{
		Parkiranid:     parkiranID,
		Nama:           nama,
		NPM:            npm,
		Prodi:          prodi,
		NamaKendaraan:  namaKendaraan,
		NomorKendaraan: nomorKendaraan,
		JenisKendaraan: jenisKendaraan,
	}
	return InsertOneDoc(MongoConn, "user", req)
}

// fungi untuk mengurutkan id
func SequenceAutoIncrement(mongoconn *mongo.Database, sequenceName string) int {
	// Membuat filter untuk mencari dokumen dengan _id sesuai sequenceName
	filter := bson.M{"_id": sequenceName}
	
	// Membuat update untuk menambahkan nilai sequence
	update := bson.M{"$inc": bson.M{"seq": 1}}

	// Membuat variabel untuk menyimpan hasil sequence
	var result struct {
		Seq int `bson:"seq"`
	}

	// Membuat opsi untuk mengembalikan dokumen setelah update
	after := options.After
	opt := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	// Mengambil koleksi untuk sequence dari database dan melakukan operasi FindOneAndUpdate
	collection := mongoconn.Collection("counters")
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&result)
	if err != nil {
		// Menghandle kesalahan jika terjadi
	}

	// Mengembalikan nilai sequence yang telah diincrement
	return result.Seq
}

func CreateParkiranNomor(mongoconn *mongo.Database, sequenceName string, parkiran Parkiran) error {
	// Mendapatkan nilai sequence berikutnya
	nextSequence := SequenceAutoIncrement(mongoconn, sequenceName)

	// Menetapkan nilai Nomor pada struktur data Parkiran
	parkiran.Nomor = &nextSequence

	// Menyisipkan data Parkiran ke dalam database
	_, err := mongoconn.Collection("parkiran").InsertOne(context.TODO(), parkiran)
	if err != nil {
		// Mengembalikan kesalahan jika gagal menyisipkan data Parkiran
		return fmt.Errorf("gagal menyisipkan data Parkiran: %v", err)
	}

	// Mengembalikan nilai-nilai tanpa kesalahan jika berhasil
	return nil
}
