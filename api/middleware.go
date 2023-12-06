package api 

import (
	"net/http"
    "eardis/tools"

	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtHandle func(w http.ResponseWriter,r *http.Request,t *jwt.Token) error

func jwtHandleFunc(f jwtHandle) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request){
        tokenString,err := tools.ReadHttpOnlyCookie(r); if err!=nil{
            tools.WriteJSON(w,http.StatusForbidden, ApiError{Error: "Invalid authorization"})
        }else{
            t,err := tools.ValidateJWT(tokenString)
            if err!= nil{
                tools.WriteJSON(w,http.StatusForbidden, ApiError{Error: "Invalid authorization"})
            }else if err := f(w,r,t); err!=nil {
                tools.WriteJSON(w,http.StatusBadRequest, ApiError{Error: err.Error()}) 
            }
        }
    }
}


