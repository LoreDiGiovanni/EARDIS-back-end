package main

import(
    "fmt"
    "os"
    "net/http"
    jwt "github.com/golang-jwt/jwt/v5"
)

type jwtHandle func(w http.ResponseWriter,r *http.Request,t *jwt.Token) error

func jwtHandleFunc(f jwtHandle) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request){
        tokenString := r.Header.Get("x-jwt-token")
        t,err := validateJWT(tokenString)
        if err!= nil{
            WriteJSON(w,http.StatusForbidden, ApiError{Error: "Invalid authorization"})
        }else if err := f(w,r,t); err!=nil {
            WriteJSON(w,http.StatusBadRequest, ApiError{Error: err.Error()}) 
        }
    }
}

func createUserJWT(account *User) (string, error){
    claims:= &jwt.MapClaims{
        "id": account.id,
        "pwd": account.pwd,
    }
    secret:= []byte(os.Getenv("JWT_SECRET"))
    token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
    return token.SignedString(secret)
} 

func validateJWT(tokenString string) (*jwt.Token,error){
    secret := os.Getenv("JWT_SECRET")
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		    return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	    }
	    return []byte(secret), nil
    })
}

func decodeJWT(tokenString string) (*jwt.Token,error){
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil})
}
