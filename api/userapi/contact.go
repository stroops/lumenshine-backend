package main

import (
	"net/http"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"github.com/Soneso/lumenshine-backend/pb"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	"github.com/gin-gonic/gin"
)

// Contact - user contact details
// swagger:model
type Contact struct {
	ID             int64  `json:"id"`
	ContactName    string `json:"contact_name"`
	StellarAddress string `json:"stellar_address"`
	PublicKey      string `json:"public_key"`
}

// ContactList returns the contacts of a user
// swagger:route GET /portal/user/dashboard/contact_list contact ContactList
//
// Lists all contacts of a user
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []Contact - Lists all contacts of a user
func ContactList(uc *mw.IcopContext, c *gin.Context) {
	userID := mw.GetAuthUser(c).UserID
	req := &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	r, err := dbClient.GetUserContacts(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading contacts", cerr.GeneralError))
		return
	}

	response := make([]Contact, len(r.UserContacts))
	for i, contact := range r.UserContacts {
		response[i] = Contact{
			ID:             contact.Id,
			ContactName:    contact.ContactName,
			StellarAddress: contact.StellarAddress,
			PublicKey:      contact.PublicKey,
		}
	}

	c.JSON(http.StatusOK, response)
}

//AddContactRequest request-data
//swagger:parameters AddContactRequest AddContact
type AddContactRequest struct {
	// Stellar account public key - identifies the account
	// required: required if stellar address is not specified
	PublicKey string `form:"public_key" json:"public_key"  validate:"omitempty,base64,len=56"`
	// Stellar account federation address - identifies the account
	// required: required if stellar public key is not specified
	StellarAddress string `form:"stellar_address" json:"stellar_address" validate:"max=256"`
	// Contact name
	// required:true
	ContactName string `form:"contact_name" json:"contact_name" validate:"required,max=256"`
}

//AddContactResponse - list of contacts and newly added contact id
// swagger:model
type AddContactResponse struct {
	ID       int64     `json:"id"`
	Contacts []Contact `json:"contacts"`
}

//AddContact - adds a new contact to the user
// swagger:route GET /portal/user/dashboard/add_contact contact AddContact
//
// Adds a new contact to the user
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []AddContactResponse
func AddContact(uc *mw.IcopContext, c *gin.Context) {
	var r AddContactRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	if r.PublicKey == "" && r.StellarAddress == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key", cerr.InvalidArgument, "Public key or stellar address must be specified.", ""))
		return
	}
	userID := mw.GetAuthUser(c).UserID
	req := &pb.AddUserContactRequest{
		Base:           NewBaseRequest(uc),
		UserId:         userID,
		ContactName:    r.ContactName,
		PublicKey:      r.PublicKey,
		StellarAddress: r.StellarAddress,
	}
	idResp, err := dbClient.AddUserContact(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading contacts", cerr.GeneralError))
		return
	}

	response := AddContactResponse{ID: idResp.Id}

	idReq := &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	cResp, err := dbClient.GetUserContacts(c, idReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading contacts", cerr.GeneralError))
		return
	}

	response.Contacts = make([]Contact, len(cResp.UserContacts))
	for i, contact := range cResp.UserContacts {
		response.Contacts[i] = Contact{
			ID:             contact.Id,
			ContactName:    contact.ContactName,
			StellarAddress: contact.StellarAddress,
			PublicKey:      contact.PublicKey,
		}
	}
	c.JSON(http.StatusOK, response)
}

//EditContactRequest request-data
//swagger:parameters EditContactRequest EditContact
type EditContactRequest struct {
	// Contact id
	// required : true
	ID int64 `form:"id" json:"id"`
	// Stellar account public key - identifies the account
	// required: required if stellar address is not specified
	PublicKey string `form:"public_key" json:"public_key"  validate:"omitempty,base64,len=56"`
	// Stellar account federation address - identifies the account
	// required: required if stellar public key is not specified
	StellarAddress string `form:"stellar_address" json:"stellar_address" validate:"max=256"`
	// Contact name
	// required:true
	ContactName string `form:"contact_name" json:"contact_name" validate:"required,max=256"`
}

//EditContact - edits the user contact
// swagger:route GET /portal/user/dashboard/edit_contact contact EditContact
//
// Edits the user contact
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []Contact - Lists all contacts of a user
func EditContact(uc *mw.IcopContext, c *gin.Context) {
	var r EditContactRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, r); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	if r.PublicKey == "" && r.StellarAddress == "" {
		c.JSON(http.StatusBadRequest, cerr.NewIcopError("public_key", cerr.InvalidArgument, "Public key or stellar address must be specified.", ""))
		return
	}

	req := &pb.UpdateUserContactRequest{
		Base:           NewBaseRequest(uc),
		Id:             r.ID,
		ContactName:    r.ContactName,
		PublicKey:      r.PublicKey,
		StellarAddress: r.StellarAddress,
	}
	_, err := dbClient.UpdateUserContact(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating contact", cerr.GeneralError))
		return
	}

	userID := mw.GetAuthUser(c).UserID
	idReq := &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	cResp, err := dbClient.GetUserContacts(c, idReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading contacts", cerr.GeneralError))
		return
	}

	response := make([]Contact, len(cResp.UserContacts))
	for i, contact := range cResp.UserContacts {
		response[i] = Contact{
			ID:             contact.Id,
			ContactName:    contact.ContactName,
			StellarAddress: contact.StellarAddress,
			PublicKey:      contact.PublicKey,
		}
	}
	c.JSON(http.StatusOK, response)
}

//RemoveContactRequest request-data
//swagger:parameters RemoveContactRequest RemoveContact
type RemoveContactRequest struct {
	// Contact id
	// required : true
	ID int64 `form:"id" json:"id"`
}

//RemoveContact - deletes the user contact
// swagger:route GET /portal/user/dashboard/remove_contact contact RemoveContact
//
// Deletes the user contact
//
// Consumes:
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: []Contact - Lists all contacts of a user
func RemoveContact(uc *mw.IcopContext, c *gin.Context) {
	var r RemoveContactRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	req := &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   r.ID,
	}
	_, err := dbClient.DeleteUserContact(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error deleting contact", cerr.GeneralError))
		return
	}

	userID := mw.GetAuthUser(c).UserID
	idReq := &pb.IDRequest{
		Base: NewBaseRequest(uc),
		Id:   userID,
	}
	cResp, err := dbClient.GetUserContacts(c, idReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading contacts", cerr.GeneralError))
		return
	}

	response := make([]Contact, len(cResp.UserContacts))
	for i, contact := range cResp.UserContacts {
		response[i] = Contact{
			ID:             contact.Id,
			ContactName:    contact.ContactName,
			StellarAddress: contact.StellarAddress,
			PublicKey:      contact.PublicKey,
		}
	}
	c.JSON(http.StatusOK, response)
}
