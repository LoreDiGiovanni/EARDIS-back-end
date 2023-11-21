package tools 

import(
    "net/http"
    "encoding/json"
    "crypto/sha256"
    "crypto/rand"
    "encoding/hex"
    "os"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
    w.Header().Add("Content-Type","application/json")
    w.WriteHeader(status)
    if(v != nil){
        return json.NewEncoder(w).Encode(v);
    }else{
        return nil 
    }
}

func GeneratePwd(pwd string) (string,string){
    pepper := os.Getenv("PEPPER")
    saltBytes := make([]byte, 16)
    _, err := rand.Read(saltBytes);if err != nil{
        panic(err)
    }
    salt := hex.EncodeToString(saltBytes[:])
    hash := sha256.Sum256([]byte(pepper+pwd+salt))
    return hex.EncodeToString(hash[:]),salt
}

func RiGeneratePwd(pwd string,salt string) (string){
    pepper := os.Getenv("PEPPER")
    hash := sha256.Sum256([]byte(pepper+pwd+salt))
    return hex.EncodeToString(hash[:])
}

