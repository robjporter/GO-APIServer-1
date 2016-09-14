package badges

func GetPublicFolders() []string {
	return nil
}

func GetRoutes() map[string]map[string]interface{} {
	Routes := map[string]map[string]interface{}{
		"badge_color_get": map[string]interface{}{
			"handler": Draw, //Handler(controller in other frameworks)
			"path":    "/badge/draw/:subject/:status/:color",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "badge",
		},
		"badge_get": map[string]interface{}{
			"handler": Draw, //Handler(controller in other frameworks)
			"path":    "/badge/draw/:subject/:status",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "badge",
		},

	}
	return Routes
}