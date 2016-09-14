package home

import (
	//"fmt"
)

func GetPublicFolders() []string {
	return nil
}

func GetRoutes() map[string]map[string]interface{}  {
	Routes := map[string]map[string]interface{}{
		"home_index_get": map[string]interface{}{
			"handler": Index, //Handler(controller in other frameworks)
			"path":    "/home/",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "home",
		},
	}
	return Routes
}