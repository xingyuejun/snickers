package snickers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/rest"
	"github.com/flavioribeiro/snickers/types"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest API", func() {
	Context("helper functions", func() {
		It("should write the error as json", func() {
			w := httptest.NewRecorder()
			rest.HTTPError(w, http.StatusOK, "error here", errors.New("database broken"))

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(`{"error": "error here: database broken"}`))
		})
	})

	Context("/presets location", func() {
		var (
			response   *httptest.ResponseRecorder
			server     *mux.Router
			dbInstance db.DatabaseInterface
		)

		BeforeEach(func() {
			response = httptest.NewRecorder()
			server = rest.NewRouter()
			dbInstance, _ = db.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("GET should return application/json on its content type", func() {
			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("GET should return stored presets", func() {
			Skip("we should sort the arrays before compare")
			examplePreset1 := types.Preset{Name: "a"}
			examplePreset2 := types.Preset{Name: "b"}
			dbInstance.StorePreset(examplePreset1)
			dbInstance.StorePreset(examplePreset2)

			expected, _ := json.Marshal(`[{"name":"a","video":{},"audio":{}},{"name":"b","video":{},"audio":{}}]`)

			request, _ := http.NewRequest("GET", "/presets", nil)
			server.ServeHTTP(response, request)
			responseBody, _ := json.Marshal(string(response.Body.String()))

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(responseBody).To(Equal(expected))
		})

		It("POST should save a new preset", func() {
			preset := []byte(`{"name": "storedPreset", "video": {},"audio": {}}`)
			request, _ := http.NewRequest("POST", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			presets, _ := dbInstance.GetPresets()
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(len(presets)).To(Equal(1))
		})

		It("POST with malformed preset should return bad request", func() {
			preset := []byte(`{"neime: "badPreset}}`)
			request, _ := http.NewRequest("POST", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("PUT with a new preset should update the preset", func() {
			dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
			preset := []byte(`{"name":"examplePreset","Description": "new description","video": {},"audio": {}}`)

			request, _ := http.NewRequest("PUT", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			presets, _ := dbInstance.GetPresets()
			newPreset := presets[0]
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(newPreset.Description).To(Equal("new description"))
		})

		It("PUT with malformed preset should return bad request", func() {
			dbInstance.StorePreset(types.Preset{Name: "examplePreset"})
			preset := []byte(`{"name":"examplePreset","Description: "new description","video": {},"audio": {}}`)

			request, _ := http.NewRequest("PUT", "/presets", bytes.NewBuffer(preset))
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("GET for a given preset should return preset details", func() {
			examplePreset := types.Preset{
				Name:         "examplePreset",
				Description:  "This is an example of preset",
				Container:    "mp4",
				Profile:      "high",
				ProfileLevel: "3.1",
				RateControl:  "VBR",
				Video: types.VideoPreset{
					Width:         "720",
					Height:        "1080",
					Codec:         "h264",
					Bitrate:       "10000",
					GopSize:       "90",
					GopMode:       "fixed",
					InterlaceMode: "progressive",
				},
				Audio: types.AudioPreset{
					Codec:   "aac",
					Bitrate: "64000",
				},
			}
			dbInstance.StorePreset(examplePreset)
			expected, _ := json.Marshal(examplePreset)

			request, _ := http.NewRequest("GET", "/presets/examplePreset", nil)
			server.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(response.Body.String()).To(Equal(string(expected)))
		})
	})

	Context("/jobs location", func() {
		var (
			response   *httptest.ResponseRecorder
			server     *mux.Router
			dbInstance db.DatabaseInterface
		)

		BeforeEach(func() {
			response = httptest.NewRecorder()
			server = rest.NewRouter()
			dbInstance, _ = db.GetDatabase()
			dbInstance.ClearDatabase()
		})

		It("GET should return application/json on its content type", func() {
			request, _ := http.NewRequest("GET", "/jobs", nil)
			server.ServeHTTP(response, request)
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})

		It("GET should return stored jobs", func() {
			Skip("we should sort the arrays before compare")

			exampleJob1 := types.Job{ID: "123"}
			exampleJob2 := types.Job{ID: "321"}
			dbInstance.StoreJob(exampleJob1)
			dbInstance.StoreJob(exampleJob2)

			expected, _ := json.Marshal(`[{"id":"123","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""},{"id":"321","source":"","destination":"","preset":{"video":{},"audio":{}},"status":"","progress":""}]`)

			request, _ := http.NewRequest("GET", "/jobs", nil)
			server.ServeHTTP(response, request)
			responseBody, _ := json.Marshal(string(response.Body.String()))

			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(responseBody).To(Equal(expected))
		})

		It("POST should create a new job", func() {
			dbInstance.StorePreset(types.Preset{Name: "presetName"})
			jobJSON := []byte(`{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`)
			request, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jobJSON))
			server.ServeHTTP(response, request)

			jobs, _ := dbInstance.GetJobs()
			Expect(response.Code).To(Equal(http.StatusOK))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
			Expect(len(jobs)).To(Equal(1))
			job := jobs[0]
			Expect(job.Source).To(Equal("http://flv.io/src.mp4"))
			Expect(job.Destination).To(Equal("s3://l@p:google.com"))
			Expect(job.Preset.Name).To(Equal("presetName"))
		})

		It("POST should return BadRequest if preset is not set when creating a new job", func() {
			jobJSON := []byte(`{"source": "http://flv.io/src.mp4", "destination": "s3://l@p:google.com", "preset": "presetName"}`)
			request, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jobJSON))
			server.ServeHTTP(response, request)
			responseBody, _ := json.Marshal(string(response.Body.String()))

			expected, _ := json.Marshal(`{"error": "retrieving preset: preset not found"}`)
			Expect(responseBody).To(Equal(expected))
			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(response.HeaderMap["Content-Type"][0]).To(Equal("application/json; charset=UTF-8"))
		})
	})
})
