package bing

func GetPublicFolders() []string {
	return nil
}

func GetRoutes() map[string]map[string]interface{} {
	Routes := map[string]map[string]interface{}{
		"bing_index_get": map[string]interface{}{
			"handler": Index, //Handler(controller in other frameworks)
			"path":    "/bing/",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "bing",
		},
		"bing_index_post": map[string]interface{}{
			"handler": Index, //Handler(controller in other frameworks)
			"path":    "/bing/",        //Relative URL
			"method":  "POST",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "bing",
		},
		"bing_piccy_get": map[string]interface{}{
			"handler": Piccy,
			"path": "/bing/pic",
			"method": "GET",
			"module": "bing",
		},
		"bing_json_get": map[string]interface{}{
			"handler": Json,
			"path": "/bing/json",
			"method": "GET",
			"module": "bing",
		},
	}
	return Routes
}
