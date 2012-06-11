package str2hash

//Hashes
import "crypto/md5"
import "crypto/sha1"
import "encoding/hex"

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
