package tools 

import(
    "net/http"
    "encoding/json"
    "crypto/sha256"
    "crypto/rand"
    "encoding/hex"
    "os"
    "time"
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
func WriteHttpOnlyCookie(jwt string)(http.Cookie){
    cookie := http.Cookie{
		Name:     "eardis",
		Value:    jwt,
        Path:     "/",
		HttpOnly: true,
		Secure:   false, // true for HTTPS 
		SameSite: http.SameSiteStrictMode,
        Domain:   "localhost",
		Expires:  time.Now().Add(24 * time.Hour), 
	}
    return cookie
}

func ReadHttpOnlyCookie(r *http.Request)(string,error){
    cookie, err := r.Cookie("eardis")
    if err != nil{
        return "",err
    }else{
        return cookie.Value, nil
    }
}
