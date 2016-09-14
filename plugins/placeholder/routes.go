package placeholder

func GetPublicFolders() []string {
	folders := []string {
		"/placeholder/output",
	}
	return folders
}

func GetRoutes() map[string]map[string]interface{} {
	Routes := map[string]map[string]interface{}{
		"placeholder_blank_square_get": map[string]interface{}{
			"handler": Blank, //Handler(controller in other frameworks)
			"path":    "/placeholder/blank/:width",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_blank_custom_get": map[string]interface{}{
			"handler": Blank, //Handler(controller in other frameworks)
			"path":    "/placeholder/blank/:width/:height",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_draw_square_get": map[string]interface{}{
			"handler": Draw, //Handler(controller in other frameworks)
			"path":    "/placeholder/draw/:width",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_draw_custom_get": map[string]interface{}{
			"handler": Draw, //Handler(controller in other frameworks)
			"path":    "/placeholder/draw/:width/:height",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_draw_text_get": map[string]interface{}{
			"handler": Draw, //Handler(controller in other frameworks)
			"path":    "/placeholder/draw/:width/:height/:text",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_save_square_get": map[string]interface{}{
			"handler": Save, //Handler(controller in other frameworks)
			"path":    "/placeholder/save/:width",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_save_custom_get": map[string]interface{}{
			"handler": Save, //Handler(controller in other frameworks)
			"path":    "/placeholder/save/:width/:height",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_save_text_get": map[string]interface{}{
			"handler": Save, //Handler(controller in other frameworks)
			"path":    "/placeholder/save/:width/:height/:text",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_base64_square_get": map[string]interface{}{
			"handler": Base64, //Handler(controller in other frameworks)
			"path":    "/placeholder/base64/:width",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_base64_custom_get": map[string]interface{}{
			"handler": Base64, //Handler(controller in other frameworks)
			"path":    "/placeholder/base64/:width/:height",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_base64_text_get": map[string]interface{}{
			"handler": Base64, //Handler(controller in other frameworks)
			"path":    "/placeholder/base64/:width/:height/:text",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
		"placeholder_index_get": map[string]interface{}{
			"handler": Index, //Handler(controller in other frameworks)
			"path":    "/placeholder/",        //Relative URL
			"method":  "GET",      //GET, POST, PUT, DELETE or empty string for ANY
			"module": "placeholder",
		},
	}
	return Routes
}
