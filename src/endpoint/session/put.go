package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/rebel-l/sessionservice/src/endpoint"
	"github.com/rebel-l/sessionservice/src/request"
	"github.com/rebel-l/sessionservice/src/response"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Put handles all PUT requests to session endpoint
type Put struct {
	session *Session
}

func NewPut(session *Session) *Put {
	put := new(Put)
	put.session = session
	return put
}

func (put *Put) Handler(res http.ResponseWriter, req *http.Request) {
	log.Info("Executing session PUT")
	msg := endpoint.NewMessage(res)

	// read request body
	requestBody, err := put.getRequestBody(req)
	if err != nil {
		log.Errorf("Parsing request body failed: %s", err)
		msg.Plain(BadRequestText, http.StatusBadRequest)
		return
	}

	// store session
	sessionResponse := response.NewSession(requestBody.Id, put.session.Config.SessionLifetime)
	status := http.StatusOK
	if requestBody.Id == "" {
		log.Debugf("Create new session: %s", sessionResponse.Id)
		status = http.StatusCreated
	} else {
		// TODO: move code to own method
		log.Debugf("Update session: %s", requestBody.Id)

		// 1. load stored session
		result := put.session.Redis.Get(requestBody.Id)

		// 2. if key not found ==> return error (404)
		storageData, err := result.Result()
		if err != nil {
			log.Errorf("Session Id %s not found or has expired: %s", sessionResponse.Id, err)
			msg.Plain("Session was not found or has expired.",	http.StatusNotFound)
			return
		}

		// 3. merge data with current stored
		log.Debugf("Loaded session data for %s: %s", sessionResponse.Id, storageData)
		oldData := make(map[string]string)
		err = json.Unmarshal([]byte(storageData), &oldData)
		if err != nil {
			log.Errorf("Data loaded for %s can't be turned into map: %s", sessionResponse.Id, err)
			msg.InternalServerError(InternalServerErrorText)
			return
		}

		requestBody.Data = put.mergeData(oldData, requestBody.Data)
	}

	// store data in redis
	err = put.storeData(sessionResponse, requestBody.Data, msg)
	if err != nil {
		log.Error(err)
		msg.InternalServerError(InternalServerErrorText)
		return
	}

	// write response
	msg.SendJson(sessionResponse, status)

	log.Info("Executing session PUT done!")
}

func (put *Put) storeData (response *response.Session, data map[string]string, msg endpoint.Message) error {
	lifetime := time.Duration(response.Lifetime) * time.Second
	dataJson, err := json.Marshal(data)
	if err != nil {
		return errors.New(fmt.Sprintf("Saving Id %s failed: %s", response.Id, err))
	}

	result := put.session.Redis.Set(response.Id, dataJson, lifetime)
	if result.Err() != nil {
		return errors.New(fmt.Sprintf("Saving Id %s failed: %s", response.Id, result.Err().Error()))
	}
	response.Data = data

	log.Debug("Data added ...")
	for key, value := range data {
		log.Debugf("%s: %s", key, value)
	}

	return nil
}

func (put *Put) getRequestBody(req *http.Request) (body request.Update, err error) {
	// read request body
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err = decoder.Decode(&body)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to read request body: %s", err))
		return
	}

	err = put.validateRequestBody(&body)
	return
}

func (put *Put) validateRequestBody(body *request.Update) error {
	// Id must be uuid or empty string
	if body.Id != "" {
		body.Id = uuid.Parse(body.Id).String()
		log.Debugf("UUID parsed: %s", body.Id)
		if len(body.Id) != UUIDLENGTH {
			return errors.New("request body validation failed ==> wrong UUID provided")
		}
	}

	// data field must have entries
	if len(body.Data) < 1 {
		return errors.New("request body validation failed ==> no data field received")
	}

	return nil
}

func (put *Put) mergeData(old map[string]string, new map[string]string) map[string]string {
	result := make(map[string]string)
	for key, value := range old {
		result[key] = value
	}

	for key, value := range new {
		result[key] = value
	}

	return result
}
