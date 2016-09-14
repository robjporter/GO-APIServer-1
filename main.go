package main

import (
	"fmt"
	"flag"
	//"./classes"
	"time"
	"strings"
	"github.com/kataras/iris"
	"github.com/dgrijalva/jwt-go"
)

var (
	api 		myIris
	development = true
	debug 		= true
)

func main() {
	flag.IntVar(&instance, "instance",1,"Instance ID Number for this application execution cycle.")
	flag.Parse()
	api.run("./config/config.json")
}

func (a *myIris) addTempStaticRoutes() {
	a.core.Get(baseURL + "/demo",demo)
	a.core.Get(baseURL + "/demo2",demo2)
	a.core.Get(baseURL + "/theming",theming)
	a.core.Get(baseURL + "/menu",menu)
	a.core.Get(baseURL + "/menu2",menu2)
	a.core.Get(baseURL + "/test",h)
	a.core.Get(baseURL + "/tester",h)
	a.core.Get("/secured/jwt",a.getJWT)
	a.core.Get("/secured/ping", a.SecuredPingHandler)
	a.addMethodPathStatsNumber("GET",9)
	a.paths = append(a.paths,baseURL + "/")
	a.paths = append(a.paths,baseURL + "/demo")
	a.paths = append(a.paths,baseURL + "/menu")
	a.paths = append(a.paths,baseURL + "/menu2")
	a.paths = append(a.paths,baseURL + "/size/:width")
	a.paths = append(a.paths,baseURL + "/size/:width/:height")
	a.paths = append(a.paths,baseURL + "/test")
	a.paths = append(a.paths,baseURL + "/tester")
	a.paths = append(a.paths,baseURL + "/secured/jwt")
	a.paths = append(a.paths,baseURL + "/secured/ping")
}

func theming(ctx *iris.Context) {
	ctx.Render("templates/" + theme + "/pages/theming.html", map[string]interface{}{
		"page_title": "Iris",
		"page_head": "",
		"page_content": "CONTENT",
		"page_script": "",
		"theme": theme,
	})
}

func menu(ctx *iris.Context) {
	ctx.Render(theme + "/pages/menu.html", map[string]interface{}{
		"page_title": "Iris",
		"page_head": "",
		"page_content": "CONTENT",
		"page_script": "",
		"theme": theme,
	})
}

func menu2(ctx *iris.Context) {
	ctx.Render(theme + "/pages/menu2.html", map[string]interface{}{
		"page_title": "Iris",
		"page_head": "",
		"page_content": "CONTENT",
		"page_script": "",
		"theme": theme,
	})
}

func demo(ctx *iris.Context) {
	ctx.Render(theme + "/pages/demo.html", map[string]interface{}{
		"page_title": "Iris",
		"page_head": "",
		"page_content": "CONTENT",
		"page_script": "",
		"theme": theme,
	})
}

func demo2(ctx *iris.Context) {
	ctx.Render(theme + "/pages/demo2.html", map[string]interface{}{
		"page_title": "Iris",
		"page_head": "",
		"page_content": "CONTENT",
		"page_script": "",
		"theme": theme,
	})
}

func h(ctx *iris.Context) {
	ctx.Write(ctx.PathString())
}

func (a *myIris) getJWT(ctx *iris.Context) {
	mySigningKey := []byte(a.config.GetString("middleware.jwt.key"))
	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		"barbar0011",
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
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		ctx.Write(err.Error())
	} else {
		ctx.Write(tokenString)
	}
}

func (a *myIris) SecuredPingHandler(ctx *iris.Context) {
	response := "All good. You only get this message if you're authenticated"
	//mySigningKey := []byte(a.config.GetString("middleware.jwt.key"))

	//fmt.Println(ctx.RequestHeader("Authorization"))
	tokenString,_ := FromAuthHeader(ctx)
	tokenString = strings.TrimSpace(tokenString)

	token2, err2 := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.GetString("middleware.jwt.key")), nil
	})
	if err2 == nil {
		if token2.Valid {
			if claims2, ok := token2.Claims.(*MyCustomClaims); ok && token2.Valid {
				fmt.Printf("NAME: %v\n", claims2.Name)
				fmt.Printf("USERNAME: %v\n", claims2.Username)
				fmt.Printf("EMAIL: %v\n", claims2.Email)
				fmt.Printf("ROLE: %v\n", claims2.Role)
				fmt.Printf("CAPABILITIES: %v\n", claims2.Capabilities)
				fmt.Printf("EXPIRES: %v\n", claims2.StandardClaims.ExpiresAt)
				fmt.Printf("ISSUER: %v\n", claims2.StandardClaims.Issuer)
				fmt.Printf("SUBJECT: %v\n", claims2.StandardClaims.Subject)
				fmt.Printf("AUDIENCE: %v\n", claims2.StandardClaims.Audience)
				fmt.Printf("BEFORE: %v\n", claims2.StandardClaims.NotBefore)
				fmt.Printf("ISSUED AT: %v\n", claims2.StandardClaims.IssuedAt)
				ctx.JSON(iris.StatusOK, response)
			} else {
				ctx.JSON(iris.StatusForbidden, err2.Error())
			}
		}
	} else {
		ctx.JSON(iris.StatusForbidden, err2.Error())
	}
}

func FromAuthHeader(ctx *iris.Context) (string, error) {

	authHeader := ctx.RequestHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
