package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{"GET", "/", Index},

	// Preset routes
	Route{"POST", "/presets", CreatePreset},
	Route{"PUT", "/presets", UpdatePreset},
	Route{"GET", "/presets", ListPresets},
	Route{"GET", "/presets/{presetName}", GetPresetDetails},

	// Job routes
	Route{"POST", "/jobs", CreateJob},
	Route{"GET", "/jobs", ListJobs},
	Route{"GET", "/jobs/{jobId}", GetJobDetails},
}
