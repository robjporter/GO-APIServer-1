package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)


type MyCustomClaims struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	Capabilities string `json:"capabilities"`
	jwt.StandardClaims
}

func main() {
	mySigningKey := []byte("AllYourBase")
	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		"test@test.com",
		"admin",
		"100010101010101",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(10)).Unix(),
			Issuer:    "test",
			Subject:   "Validation",
			Audience:  "Global",
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("TOKEN: %v\nERROR: %v\n", ss, err)
	
	token2, err2 := jwt.ParseWithClaims(ss, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("AllYourBase"), nil
	})
	if token2.Valid {
		if claims2, ok := token2.Claims.(*MyCustomClaims); ok && token2.Valid {
			fmt.Printf("NAME: %v\n", claims2.Name)
			fmt.Printf("EMAIL: %v\n", claims2.Email)
			fmt.Printf("ROLE: %v\n", claims2.Role)
			fmt.Printf("CAPABILITIES: %v\n", claims2.Capabilities)
			fmt.Printf("EXPIRES: %v\n", claims2.StandardClaims.ExpiresAt)
			fmt.Printf("ISSUER: %v\n", claims2.StandardClaims.Issuer)
			fmt.Printf("SUBJECT: %v\n", claims2.StandardClaims.Subject)
			fmt.Printf("AUDIENCE: %v\n", claims2.StandardClaims.Audience)
			fmt.Printf("BEFORE: %v\n", claims2.StandardClaims.NotBefore)
			fmt.Printf("ISSUED AT: %v\n", claims2.StandardClaims.IssuedAt)
		} else {
			fmt.Println(err2)
		}
	}
}