package str2hash

//Hashes
import "crypto/md5"
import "crypto/sha1"
import "encoding/hex"
import "errors"

type Hash struct {
	Hash     string
	HashFunc func(string) string
}

func New(hashString string) (hash *Hash, err error) {

	hash = new(Hash)

	//Validate hash by trying to decode it as hexadecimal
	_, err = hex.DecodeString(hashString)
	if err != nil {
		return nil, err
	}

	hash.Hash = hashString

	switch len(hashString) {
	case 32:
		hash.HashFunc = MD5
	case 40:
		hash.HashFunc = SHA1
	default:
		return nil, errors.New("Invalid hash - length mismatch")
	}

	return hash, nil
}

//Hashes string input to MD5 string output 
func MD5(input string) string {

	//Initiate MD5 type
	h := md5.New()

	//Set input
	h.Write([]byte(input))

	//Return hash as string
	return hex.EncodeToString(h.Sum(nil))
}

//Hashes string input to SHA-1 string output
func SHA1(input string) string {

	//Initiate SHA-1 type
	h := sha1.New()

	//Set input
	h.Write([]byte(input))

	//Return hash as string
	return hex.EncodeToString(h.Sum(nil))
}
