package main

import (
    "fmt"
    "github.com/tidwall/gjson"
)

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
const json2 = `{
  "programmers": [
    {
      "firstName": "Janet",
      "lastName": "McLaughlin"
    }, {
      "firstName": "Elliotte",
      "lastName": "Hunter"
    }, {
      "firstName": "Jason",
      "lastName": "Harold"
    },
    {
        "actions": [
            { "action1": "test" }
        ]
    }
  ]
}`

func main() {
    value := gjson.Get(json, "name.last")
    println(value.String())
fmt.Println("=========================================================")
    result := gjson.Get(json2, "programmers.#.lastName")
    for _,name := range result.Array() {
        println(name.String())
    }
fmt.Println("=========================================================")
    fmt.Println(gjson.Parse(json).Get("name").Get("last"))
    fmt.Println(gjson.Get(json, "name").Get("last"))
    fmt.Println(gjson.Get(json, "name.last"))
fmt.Println("=========================================================")
    value = gjson.Get(json, "name.last")
    if !value.Exists() {
        println("no last name")
    } else {
        println(value.String())
    }

    // Or as one step
    if gjson.Get(json, "name.last").Exists(){
        println("has a last name")
    }
fmt.Println("=========================================================")
    fmt.Println(gjson.Get(json2,"programmers.#"))
    fmt.Println(gjson.Get(json2,"programmers.0.firstName"))
    fmt.Println(gjson.Get(json2,"programmers.3.actions.#"))
    fmt.Println(gjson.Get(json2,"programmers.#.firstName"))
}
