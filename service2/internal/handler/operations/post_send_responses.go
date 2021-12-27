// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"service2/models"
)

// PostSendOKCode is the HTTP code returned for type PostSendOK
const PostSendOKCode int = 200

/*PostSendOK Success

swagger:response postSendOK
*/
type PostSendOK struct {

	/*
	  In: Body
	*/
	Payload models.ArrayOfHash `json:"body,omitempty"`
}

// NewPostSendOK creates PostSendOK with default headers values
func NewPostSendOK() *PostSendOK {

	return &PostSendOK{}
}

// WithPayload adds the payload to the post send o k response
func (o *PostSendOK) WithPayload(payload models.ArrayOfHash) *PostSendOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post send o k response
func (o *PostSendOK) SetPayload(payload models.ArrayOfHash) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostSendOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.ArrayOfHash{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// PostSendBadRequestCode is the HTTP code returned for type PostSendBadRequest
const PostSendBadRequestCode int = 400

/*PostSendBadRequest Bad request

swagger:response postSendBadRequest
*/
type PostSendBadRequest struct {
}

// NewPostSendBadRequest creates PostSendBadRequest with default headers values
func NewPostSendBadRequest() *PostSendBadRequest {

	return &PostSendBadRequest{}
}

// WriteResponse to the client
func (o *PostSendBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// PostSendInternalServerErrorCode is the HTTP code returned for type PostSendInternalServerError
const PostSendInternalServerErrorCode int = 500

/*PostSendInternalServerError Internal Server Error

swagger:response postSendInternalServerError
*/
type PostSendInternalServerError struct {
}

// NewPostSendInternalServerError creates PostSendInternalServerError with default headers values
func NewPostSendInternalServerError() *PostSendInternalServerError {

	return &PostSendInternalServerError{}
}

// WriteResponse to the client
func (o *PostSendInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
