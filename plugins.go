package main

import (
	//"fmt"
	"./plugins/admin"
	"./plugins/badges"
	"./plugins/bing"
	"./plugins/core"
	"./plugins/home"
	"./plugins/placeholder"
)

func (a *myIris) buildPluginPublicFolders() {
	mainFolders := []string{}
	folders := admin.GetPublicFolders()
	for folder := range folders {
		mainFolders = append(mainFolders,folders[folder])
	}
	folders = badges.GetPublicFolders()
	for folder := range folders {
		mainFolders = append(mainFolders,folders[folder])
	}
	folders = bing.GetPublicFolders()
	for folder := range folders {
		mainFolders = append(mainFolders,folders[folder])
	}
	folders = core.GetPublicFolders()
	for folder := range folders {
		mainFolders = append(mainFolders,folders[folder])
	}
	folders = home.GetPublicFolders()
	for folder := range folders {
		mainFolders = append(mainFolders,folders[folder])
	}
	folders = placeholder.GetPublicFolders()
	for folder := range folders {
		mainFolders = append(mainFolders,folders[folder])
	}
	a.pluginFolders = mainFolders
}

func (a *myIris) buildPluginRoutes() {
	mainRoutes := make([]interface{}, 0)
	routes := admin.GetRoutes()
	for route := range routes {
		mainRoutes = append(mainRoutes, routes[route])
	}
	routes = badges.GetRoutes()
	for route := range routes {
		mainRoutes = append(mainRoutes, routes[route])
	}
	routes = bing.GetRoutes()
	for route := range routes {
		mainRoutes = append(mainRoutes, routes[route])
	}
	routes = core.GetRoutes()
	for route := range routes {
		mainRoutes = append(mainRoutes, routes[route])
	}
	routes = home.GetRoutes()
	for route := range routes {
		mainRoutes = append(mainRoutes, routes[route])
	}
	routes = placeholder.GetRoutes()
	for route := range routes {
		mainRoutes = append(mainRoutes, routes[route])
	}
	a.mainRoutes = mainRoutes
}