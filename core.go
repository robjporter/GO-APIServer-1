package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"path"
	"errors"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
	"runtime"
	"os/signal"
	"path/filepath"
	"./modules/funcs"
	"./modules/pushbullet"
	"encoding/json"
	"github.com/kataras/iris"
	//"github.com/tidwall/redcon"
	"github.com/dgrijalva/jwt-go"
	"github.com/jasonlvhit/gocron"
	"github.com/kataras/iris/config"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/secure"
	conf "github.com/roporter/go-libs/go-config"
	"github.com/iris-contrib/middleware/recovery"
	//"github.com/iris-contrib/plugin/iriscontrol"
	"github.com/iris-contrib/template/django"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/roporter/go-libs/middleware/stats"
	"github.com/roporter/go-libs/middleware/logger"
	"github.com/roporter/go-libs/middleware/headers"
	"github.com/roporter/go-libs/middleware/ipfilter"
	"github.com/roporter/go-libs/middleware/requestid"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/Sirupsen/logrus"
)

// TODO:
// Change paths if internet connection

// TODO:
// https://github.com/cybozu-go/cmd
// https://github.com/tmrts/go-patterns/blob/master/README.md
// https://github.com/tomnomnom/gron
// https://github.com/tidwall/redcon
// https://github.com/uniplaces/carbon
// https://github.com/gravitational/teleconsole
// https://github.com/augustoroman/sandwich
// http://country.io/data/
// https://github.com/uber-go/zap
// https://github.com/firstrow/logvoyage
// https://github.com/mgutz/logxi
// https://github.com/dimiro1/health
// https://github.com/firstrow/tcp_server
// https://github.com/stretchr/testify
// https://github.com/verdverm/frisby
// https://github.com/beefsack/go-rate
// https://github.com/sadlil/go-trigger



var (
	baseURL = "/api/v1"
	theme = ""
	instance    int
	version = "0.6.1"
	logFile *os.File
)

type MyCustomClaims struct {
	Name string 		`json:"name"`
	Username string 	`json:"username"`
	Email string 		`json:"email"`
	Role string 		`json:"role"`
	Capabilities string `json:"capabilities"`
	jwt.StandardClaims
}

type myIris struct {
	core 			iris.Framework
	config 			conf.Config
	stats 			*stats.Stats
	ipfilter		*ipfilter.IPFilter
	errorLogger 	iris.HandlerFunc
	scheduler		*gocron.Scheduler
	consul 			*consulapi.Client
	jwt				*jwtmiddleware.Middleware
	system			SystemData
	logLogger		*logrus.Logger
	log				*logrus.Entry
	paths       	[]string
	mainRoutes  	interface{}
	pluginFolders 	[]string
	plugins     	interface{}
	version     	string
}

type SystemData struct {
	ip				string
	hostname		string
	city			string
	region 			string
	country 		string
	location 		string
	organisation	string
	getRoutes		int
	postRoutes		int
	headRoutes		int
	putRoutes		int
	deleteRoutes    int
}

type myPlugin struct{}

func init() {
	api = myIris{
		config: 		conf.NewConfig(),
		stats: 			stats.New(),
		errorLogger: 	logger.New(iris.Logger),
		scheduler: 		gocron.NewScheduler(),
		mainRoutes: 	make([]interface{}, 0),
		paths:			[]string{},
		version:        version,
		pluginFolders:	[]string{},
		logLogger:		logrus.New(),
	}
	api.logLogger.Formatter = new(logrus.JSONFormatter)
	api.logLogger.Out = os.Stderr
	api.logLogger.Level = logrus.DebugLevel
	api.logLogger.WithFields(logrus.Fields{"Line":"105","Function":"init","File":"core"}).Info("Application intialisation started")
	api.core = *iris.New(api.getIrisConfig())
	api.logLogger.WithFields(logrus.Fields{"Line":"107","Function":"init","File":"core"}).Debug("Web Server initialisation complete")
	api.buildPluginRoutes()
	api.logLogger.WithFields(logrus.Fields{"Line":"109","Function":"init","File":"core"}).Debug("Discovered and indexed external plugins")
	api.buildPluginPublicFolders()
	api.logLogger.WithFields(logrus.Fields{"Line":"111","Function":"init","File":"core"}).Debug("Discovered external plugin directories")
	api.ipfilter = ipfilter.New(api.getIPFilterConfig())
	api.logLogger.WithFields(logrus.Fields{"Line":"113","Function":"init","File":"core"}).Debug("Initialising IP Security filters")
	api.logLogger.WithFields(logrus.Fields{"Line":"114","Function":"init","File":"core"}).Info("Application initialisation complete")
}

func (a *myIris) alert() {
	n := pushbullet.Params{
		Title:   "title",
		Message: "mesg",
		API:     pushbullet.API,
		Token:  a.config.GetString("notification.pushbullet.token"),
	}
	pushbullet.Notify(n)
}

//myIRIS////////////////////////////////////////////////////////////////
func (a *myIris) run(configFile string) {
	a.logLogger.WithFields(logrus.Fields{"Line":"119","Function":"run","File":"core"}).Info("Starting run")
	a.setConfigDefaultParameters() // Set up some variables for loading config etc
	if a.readConfigFile(configFile) {
		a.setConfigDefaultParameters() // Reset any variables not included in loaded config with defaults.
		a.logLogger.WithFields(logrus.Fields{"Line":"122","Function":"run","File":"core"}).Debug("Configuration file loadded and defaults set successfully")
		a.setConfigOverrideParameters()
		a.logLogger.WithFields(logrus.Fields{"Line":"124","Function":"run","File":"core"}).Debug("Configuration overrides changed successfully")
		a.initialiseLog()
		a.log.WithFields(logrus.Fields{"Line":"126","Function":"run","File":"core"}).Debug("Logger initialised successfully")
		a.setupMiddleware()
		a.log.WithFields(logrus.Fields{"Line":"128","Function":"run","File":"core"}).Debug("Middleware loaded and initialised successfully")
		a.addStaticRoutes()
		a.log.WithFields(logrus.Fields{"Line":"130","Function":"run","File":"core"}).Debug("All static routes applied successfully")
		a.setFavIcon()
		a.log.WithFields(logrus.Fields{"Line":"132","Function":"run","File":"core"}).Debug("Site wide favourite icon loaded successfully")
		a.addStaticThemeRoutes()
		a.log.WithFields(logrus.Fields{"Line":"134","Function":"run","File":"core"}).Debug("Theme and template static routes applied successfully")
		a.addTempStaticRoutes()
		a.log.WithFields(logrus.Fields{"Line":"136","Function":"run","File":"core"}).Debug("Temporary static routes applied successfully")
		a.loadExternalPlugins()
		a.log.WithFields(logrus.Fields{"Line":"138","Function":"run","File":"core"}).Debug("Indexed all external plugins successfully")
		a.addErrorHandlers()
		a.log.WithFields(logrus.Fields{"Line":"140","Function":"run","File":"core"}).Debug("Added global error handlers successfully")
		a.getMachineData()
		a.log.WithFields(logrus.Fields{"Line":"142","Function":"run","File":"core"}).Debug("Discovered Server public data successfully")
		a.setTemplateConfig()
		a.log.WithFields(logrus.Fields{"Line":"144","Function":"run","File":"core"}).Debug("Completed setting up template information successfully")
		a.addCorePlugins()
		a.log.WithFields(logrus.Fields{"Line":"146","Function":"run","File":"core"}).Debug("Loaded and applied core plugins successfully")
		a.addControlCenter()
		a.log.WithFields(logrus.Fields{"Line":"148","Function":"run","File":"core"}).Debug("Initialised Iris Control Center successfully")
		a.buildCronJobs()
		a.log.WithFields(logrus.Fields{"Line":"150","Function":"run","File":"core"}).Debug("Cron jobs initialised successfully")
		a.addSecurityFilters()
		a.log.WithFields(logrus.Fields{"Line":"152","Function":"run","File":"core"}).Debug("Applied security filters successfully")
		a.startSSHServer()
		a.log.WithFields(logrus.Fields{"Line":"154","Function":"run","File":"core"}).Debug("Applied SSH server configuration successfully")
		a.log.WithFields(logrus.Fields{"Line":"155","Function":"run","File":"core"}).Debug("Starting Server....")
		a.startServer()
	} else {
		fmt.Println("Config file does not exist.  Exiting")
		a.log.WithFields(logrus.Fields{"Line":"159","Function":"run","File":"core"}).Panic("Configuration file was not found and is needed to continue loading the application.")
	}
	a.log.WithFields(logrus.Fields{"Line":"142","Function":"run","File":"core"}).Info("Run complete")
	logFile.Close()
}

func (a *myIris) initialiseLog() {
	if a.config.GetBool("server.development") {
		_, filename, _, _ := runtime.Caller(0)
		logPath := path.Join(path.Dir(filename),"/logs")
		fmt.Println(logPath)
		//os.RemoveAll(logPath)
		//os.Mkdir(logPath, os.ModeDir)
	}
	switch a.config.GetString("logging.destination") {
		case "std":
			a.logLogger.Out = os.Stderr
		case "file":
			err := errors.New("")
			filename := "log-" + strconv.FormatInt(makeTimestamp(), 10)
			logFile, err = os.OpenFile("logs/"+filename+".log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err)
			}
			a.logLogger.Out = logFile
	}
	a.log = a.logLogger.WithFields(logrus.Fields{"host": a.config.GetString("server.host") + ":" + a.config.GetString("server.port"),"server":a.config.GetString("server.name"),})
	a.log.Debug("Log configuration complete.")
//	a.log.WithFields(logrus.Fields{"animal": "walrus","size":   10,}).Info("A group of walrus emerges from the ocean")
}

func (a *myIris) startSSHServer() {
	a.core.SSH.Host = a.config.GetString("server.ssh.host") + ":" + strconv.Itoa(a.config.GetInt("server.ssh.port"))
	a.core.SSH.KeyPath = "./keys/iris_rsa" // it's auto-generated if not exists
	a.core.SSH.Users = iris.Users{"roporter": []byte("pass")}
}

func (a *myIris) getIrisConfig() config.Iris {
	return config.Iris {
		DisablePathCorrection: false,
		IsDevelopment: a.config.GetBool("server.development"),
		DisableBanner: false,
		Gzip: false,
	}
}

func (a *myIris) getMachineData() {

// TODO
// GET FULL MACHINE INVENTORY


	a.hasActiveInternetConnection()
	if a.config.GetBool("server.connected") {
		url := "http://ipinfo.io/json"
	    res, err := http.Get(url)
	    if err != nil {
	        panic(err.Error())
	    }

	    body, err := ioutil.ReadAll(res.Body)
	    if err != nil {
	        panic(err.Error())
	    }

	    var data map[string]interface{}
	    json.Unmarshal(body, &data)
		a.system.ip = data["ip"].(string)
		a.system.hostname = data["hostname"].(string)
		a.system.city = data["city"].(string)
		a.system.region = data["region"].(string)
		a.system.country = data["country"].(string)
		a.system.location = data["loc"].(string)
		a.system.organisation = data["org"].(string)
	}
}

func (a *myIris) getJWTConfig() jwtmiddleware.Config {
	return jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(a.config.GetString("middleware.jwt.key")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	}
}

func (a *myIris) getIPFilterConfig() ipfilter.Options {
	return ipfilter.Options{
		Refresh: false,
		IPDBNoFetch: false,
		BlockByDefault: false,
	}
}

func (a *myIris) addSecurityFilters() {
	a.ipfilter.AllowIP("127.0.0.1")
	block := a.config.GetKeyAsStringArray("security.blockedip")
	for i := range block {
		a.ipfilter.BlockIP(block[i])
	}
}

func (a *myIris) setupMiddleware() {
	if a.config.GetBool("logging.live") {
		a.core.Use(logger.New(iris.Logger))
	}
	a.core.Use(logger.NewLogFile(a.log))
	a.core.Use(requestid.New())
	a.core.Use(recovery.New())
	a.core.Use(cors.New(a.getCorsConfig()))
	a.core.Use(a.stats)
	a.core.Use(a.ipfilter)
	a.core.Use(headers.New(a.version))
	a.core.UseFunc(func(c *iris.Context) {
		err := a.getSecureConfig().Process(c)
		if err != nil {return}
		c.Next()
	})
}

func (a *myIris) addStaticRoutes() {
	a.core.Get("/stats", a.stats.Handle) // Endpoint to get stats
	a.core.Get("/health", func(ctx *iris.Context) {
		ctx.Text(iris.StatusOK, "OK")
	})
	a.addMethodPathStatsNumber("GET",2)
}

func (a *myIris) addControlCenter() {
	//port := a.config.GetInt("iris.control.port")
	//a.core.Plugins.Add(iriscontrol.New(port, map[string]string{
	//	a.config.GetString("iris.control.username"): a.config.GetString("iris.control.password"),
	//}))
}

func (a *myIris) addCorePlugins() {
	plugin := myPlugin{}
	a.core.Plugins.Add(plugin)
}

func (a *myIris) buildCronJobs() {
	a.scheduler.Every(30).Seconds().Do(a.hasActiveInternetConnection)
}

func (a *myIris) startServer() {
	go a.startCronJobs()
// FIX
	a.alert()


	if a.config.GetBool("consul.active") {
		go func() {
			a.core.ListenTo(a.getServerConfig())
		}()
		paths := funcs.RemoveDuplicates(a.config.GetStringArray("server.paths"))
		a.registerConsultServer(paths)
		a.displayOutput()
		a.runUntilSignal()
	} else {
		a.displayOutput()
		a.core.Must(a.core.ListenTo(a.getServerConfig()))
	}
}

func (a *myIris) displayOutput() {
	fmt.Printf("IP:             %s\n",a.system.ip)
	fmt.Printf("HOSTNAME:       %s\n",a.system.hostname)
	fmt.Printf("CITY:           %s\n",a.system.city)
	fmt.Printf("REGION:         %s\n",a.system.region)
	fmt.Printf("COUNTRY:        %s\n",a.system.country)
	fmt.Printf("LOCATION:       %s\n",a.system.location)
	fmt.Printf("COUNTRY:        %s\n",a.system.country)
	fmt.Printf("TEST:           %s\n",a.config.GetString("TEST"))
	fmt.Printf("Loaded: GET:    %d routes\n", a.system.getRoutes)
	fmt.Printf("Loaded: HEAD:   %d routes\n", a.system.headRoutes)
	fmt.Printf("Loaded: POST:   %d routes\n", a.system.postRoutes)
	fmt.Printf("Loaded: PUT:    %d routes\n", a.system.putRoutes)
	fmt.Printf("Loaded: DELETE: %d routes\n", a.system.deleteRoutes)
}

func (a *myIris) startCronJobs() {
	 <-a.scheduler.Start()
}

func (a *myIris) getServerConfig() config.Server {
	return config.Server{
		Name: 			a.config.GetString("server.name"),
		WriteTimeout: 	time.Duration(a.config.GetInt("server.writetimeout"))*time.Second,
		ReadTimeout: 	time.Duration(a.config.GetInt("server.readtimeout"))*time.Second,
		ListeningAddr:	a.config.GetString("server.host") + ":" + strconv.Itoa(a.config.GetInt("server.port")),
	}
}

func (a *myIris) addErrorHandlers() {
	a.core.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.SetStatusCode(iris.StatusNotFound)
		ctx.Write("404 error page ")
		a.stats.Serve(ctx)
		if a.config.GetBool("logging.live") {
			a.errorLogger.Serve(ctx)
		}
	})
	a.core.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		ctx.SetStatusCode(iris.StatusInternalServerError)
		ctx.Write("500 error page ")
		a.stats.Serve(ctx)
		if a.config.GetBool("logging.live") {
			a.errorLogger.Serve(ctx)
		}
	})
	a.core.OnError(iris.StatusForbidden, func(ctx *iris.Context) {
		ctx.SetStatusCode(iris.StatusForbidden)
		ctx.Write("403 error page ")
		a.stats.Serve(ctx)
		if a.config.GetBool("logging.live") {
			a.errorLogger.Serve(ctx)
		}
	})
}

func (a *myIris) getCorsConfig() cors.Options {
	return cors.Options{}
}

func (a *myIris) getSecureConfig() *secure.Secure {
	return secure.New(secure.Options{
		AllowedHosts:            []string{"ssl.example.com"},                                                                                                                         // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		SSLRedirect:             true,                                                                                                                                                // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                                                                                                                               // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "ssl.example.com",                                                                                                                                   // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"},                                                                                                     // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              315360000,                                                                                                                                           // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                                                                                                                                // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                                                                                                                                // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                                                                                                                               // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                                                                                                                                // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                                                                                                                        // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option.
		ContentTypeNosniff:      true,                                                                                                                                                // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXSSFilter:        true,                                                                                                                                                // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		ContentSecurityPolicy:   "default-src 'none'; script-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com data:; connect-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com data:;",                                                                                                                                // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
		PublicKey:               `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`, // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".

		IsDevelopment: a.config.GetBool("server.development"), // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})
}

func (a *myIris) setFavIcon() {
	a.core.Favicon("templates/" + theme + "/resources/images/favicon.ico")
}

func (a *myIris) readConfigFile(filename string) bool {
	var err error
	a.config,err = conf.ReadFromFile(filename)
	if err != nil {
		fmt.Println("Failed to read config file successfully.  Exiting application....")
		return false
	}
	theme = a.config.GetString("web.theme.name")
	return true
}

func (a *myIris) hasActiveInternetConnection() {
	_, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		a.config.AddDefaultBoolOverride("server.connected",false)
	}
	a.config.AddDefaultBoolOverride("server.connected",true)
}

func (a *myIris) addStaticThemeRoutes() {
	theme := a.config.GetString("web.theme.name")
	a.core.Static("/"+theme+"/css", "./templates/"+theme+"/resources/", 1)
	a.core.Static("/"+theme+"/js", "./templates/"+theme+"/resources/", 1)
	a.core.Static("/"+theme+"/images", "./templates/"+theme+"/resources/", 1)
	a.addMethodPathStatsNumber("GET",3)
	a.paths = append(a.paths,"/"+theme+"/css")
	a.paths = append(a.paths,"/"+theme+"/js")
	a.paths = append(a.paths,"/"+theme+"/images")
}

func (a *myIris) setConfigDefaultParameters() {
	a.config.AddDefault("server.host","0.0.0.0")
	a.config.AddDefaultInt("server.port",8080)
	a.config.AddDefaultBool("server.development",development)
	a.config.AddDefaultString("iris.control.username","username")
	a.config.AddDefaultString("iris.control.password","password")
	a.config.AddDefault("consul.host","")
	a.config.AddDefaultInt("consul.port",0)
	a.config.AddDefaultBool("consul.active",false)
	a.config.AddDefaultString("consul.protocol","http")
	a.config.AddDefaultString("consul.token","")
	a.config.AddDefaultString("consul.tagprefix", "")
	a.config.AddDefaultString("consul.namesep","-")
	a.config.AddDefaultString("consul.prefixsep",",")
	a.config.AddDefaultString("consul.intervaltimer","5s")
	a.config.AddDefaultString("consul.timeouttimer","5s")
}

func (a *myIris) setConfigOverrideParameters() {
	name := a.config.GetString("server.name")
	a.config.AddDefaultStringOverride("server.name", name + "-" + strconv.Itoa(instance))
}

func (a *myIris) setTemplateConfig() {
	templateConfig := django.DefaultConfig()
	templateConfig.DebugTemplates = a.config.GetBool("server.development")
	//templateConfig.Globals["boldme"] = funcs.Bold
	//templateConfig.Globals["italicme"] = funcs.Italic
	//templateConfig.Globals["underlineme"] = funcs.Underline
	a.core.UseTemplate(django.New(templateConfig)).Directory("./", ".html")
}

func (a *myIris) loadExternalPlugins() {
	found := map[string]bool{}

	for _, route := range a.mainRoutes.([]interface{}) {
		r := route.(map[string]interface{})
		module := strings.TrimSpace(r["module"].(string))
		path := strings.TrimSpace(r["path"].(string))
		handler := r["handler"]
		method := strings.TrimSpace(r["method"].(string))

		//fmt.Printf("MODULE: %s\nPATH: %s\nMETHOD: %s\n",module, path, method)

		if !found[module] {
			found[module] = false
			a.addNewModuleToAPI(module)
			found[module] = true
		}
		a.addMethodPathStats(method)
		a.core.HandleFunc(method, path, handler.(func(*iris.Context)))
		a.paths = append(a.paths,path)

	}
	a.recordFoundPlugins(found)
	a.recordAllPaths(a.paths)
	a.addPluginPublicFolders()
}

func (a *myIris) addMethodPathStats(method string) {
	a.addMethodPathStatsNumber(method,1)
}

func (a *myIris) addMethodPathStatsNumber(method string,number int) {
	getCount,postCount,putCount,deleteCount,headCount := 0,0,0,0,0
	switch method {
		case "GET":
			getCount += number
		case "HEAD":
			headCount += number
		case "POST":
			postCount += number
		case "PUT":
			putCount += number
		case "DELETE":
			deleteCount += number
	}
	a.system.getRoutes = a.system.getRoutes + getCount
	a.system.headRoutes = a.system.headRoutes + headCount
	a.system.postRoutes = a.system.postRoutes + postCount
	a.system.putRoutes = a.system.putRoutes + putCount
	a.system.deleteRoutes = a.system.deleteRoutes + deleteCount
}

func (a *myIris) addPluginPublicFolders() {
	for folder := range a.pluginFolders {
		a.core.Static(a.pluginFolders[folder], "./plugins" + filepath.Dir(a.pluginFolders[folder]) + "/",1)
		a.addMethodPathStats("GET")
	}
}

func (a *myIris) addNewModuleToAPI(module string) {
	a.core.Static("/"+module+"/css", "./plugins/"+module+"/resources/", 1)
	a.core.Static("/"+module+"/js", "./plugins/"+module+"/resources/", 1)
	a.core.Static("/"+module+"/images", "./plugins/"+module+"/resources/", 1)
	a.addMethodPathStatsNumber("GET",3)
	a.config.AddFile("./plugins/"+module+"/config/config.json",module,true)
	a.paths = append(a.paths,"/"+module+"/css")
	a.paths = append(a.paths,"/"+module+"/js")
	a.paths = append(a.paths,"/"+module+"/images")
}

func (a *myIris) recordAllPaths(paths []string) {
	a.config.AddDefaultStringArrayOverride("server.paths",paths)
}

func (a *myIris) recordFoundPlugins(found map[string]bool) {
	allplugins := []string{}
	allloaded := []string{}
	for k,val := range found {
		allplugins = append(allplugins, k)
		if val { allloaded = append(allloaded,k)}
	}
	a.config.AddDefaultStringArrayOverride("plugins.all",allplugins)
	a.config.AddDefaultStringArrayOverride("plugins.loaded",allloaded)
}

//CONSUL///////////////////////////////////////////////////////////////////////////////////////
func (a *myIris) registerConsultServer(urls []string) {
	var err error
	tags := a.getConsulTagsFromSlice(urls)
	service := a.getAgentServiceRegistration(tags)
	config := a.getConsulAPIConfig()
	a.consul, err = consulapi.NewClient(config)
	if err != nil {fmt.Println(err)}
	if err := a.consul.Agent().ServiceRegister(service); err != nil {fmt.Println(err)}
	//fmt.Printf("Registered service %q in consul with tags %q\n", a.config.GetString("server.name"), strings.Join(tags, ","))
	fmt.Printf("Registered service %q in consul with %d tags\n",a.config.GetString("server.name"),len(tags))
}

func (a *myIris) runUntilSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	if err := a.consul.Agent().ServiceDeregister(a.getServiceID()); err != nil {fmt.Println(err)}
	fmt.Printf("Deregistered service %q in consul\n", a.config.GetString("server.name"))
}

func (a *myIris) getServiceID() string {
	server := a.config.GetString("server.host") + ":" + strconv.Itoa(a.config.GetInt("server.port"))
	return a.config.GetString("server.name") + a.config.GetString("consul.namesep") + server
}

func (a *myIris) getConsulAPIConfig() *consulapi.Config {
	consul := a.config.GetString("consul.host") + ":" + strconv.Itoa(a.config.GetInt("consul.port"))
	proto := a.config.GetString("consul.protocol")
	token := a.config.GetString("consul.token")
	return &consulapi.Config{Address: consul, Scheme: proto, Token: token}
}

func (a *myIris) getConsulTagsFromSlice(urls []string) []string {
	var tags []string
	for _, p := range urls {
		tags = append(tags, a.config.GetString("consul.tagprefix")+p)
	}
	return tags
}

func (a *myIris) getConsulTagsFromString(urls string) []string {
	prefixes := strings.Split(urls, a.config.GetString("consul.prefixsep"))
	var tags []string
	for _, p := range prefixes {
		tags = append(tags, a.config.GetString("consul.tagprefix")+p)
	}
	return tags
}

func (a *myIris) getAgentServiceRegistration(tags []string) *consulapi.AgentServiceRegistration {
	return &consulapi.AgentServiceRegistration{
		ID:      a.getServiceID(),
		Name:    a.config.GetString("server.name"),
		Port:    a.config.GetInt("server.port"),
		Address: a.config.GetString("server.host"),
		Tags:    tags,
		Check:   a.getAgentServiceCheck(),
	}
}

func (a *myIris) getAgentServiceCheck() *consulapi.AgentServiceCheck {
	consul := a.config.GetString("server.host") + ":" + strconv.Itoa(a.config.GetInt("server.port"))
	return &consulapi.AgentServiceCheck{
		HTTP:     "http://" + consul + "/health",
		Interval: a.config.GetString("consul.intervaltimer"),
		Timeout:  a.config.GetString("consul.timeouttimer"),
	}
}
//PLUGINS///////////////////////////////////////////////////////////////////////////////////////
func (pl myPlugin) PreListen(s *iris.Framework) {
	for _, route := range s.Lookups() {
		curMethod := strings.TrimSpace(route.Method())
		if development && debug {
			if curMethod != "HEAD" {
				fmt.Printf("Func: %s | Subdomain %s | Path: %s is going to be registed with %d handler(s). \n", route.Method(), route.Subdomain(), route.Path(), len(route.Middleware()))
			}
		}
	}
}
//FUNCTIONS///////////////////////////////////////////////////////////////////////////////////////
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
