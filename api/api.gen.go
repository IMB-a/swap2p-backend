// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// Balance defines model for balance.
type Balance []SingleBalance

// Error defines model for error.
type Error struct {
	Error string `json:"error"`
}

// PersonalData defines model for personalData.
type PersonalData struct {
	Balance       Balance `json:"balance"`
	State         string  `db:"state" json:"state"`
	WalletAddress string  `db:"wallet_address" json:"walletAddress"`
}

// SingleBalance defines model for singleBalance.
type SingleBalance struct {
	Address  string `db:"asset_address" json:"address"`
	Amount   string `db:"amount" json:"amount"`
	Asset    string `db:"asset_name" json:"asset"`
	Decimals int    `db:"asset_decimals" json:"decimals"`
}

// Trade defines model for trade.
type Trade struct {
	Closed    bool   `json:"closed"`
	Id        int    `json:"id"`
	XAddress  string `json:"xAddress"`
	XAmount   string `json:"xAmount"`
	XAsset    string `json:"xAsset"`
	XDecimals int    `json:"xDecimals"`
	YAddress  string `json:"yAddress"`
	YAmount   string `json:"yAmount"`
	YAsset    string `json:"yAsset"`
	YDecimals int    `json:"yDecimals"`
}

// TradeList defines model for tradeList.
type TradeList []Trade

// PChatID defines model for pChatID.
type PChatID string

// QChatState defines model for qChatState.
type QChatState string

// QLimit defines model for qLimit.
type QLimit int

// QOffset defines model for qOffset.
type QOffset int

// QWalletAddress defines model for qWalletAddress.
type QWalletAddress string

// ErrorResp defines model for errorResp.
type ErrorResp Error

// PersonalDataResp defines model for personalDataResp.
type PersonalDataResp PersonalData

// TradesResp defines model for tradesResp.
type TradesResp TradeList

// GetAllTradesParams defines parameters for GetAllTrades.
type GetAllTradesParams struct {
	Offset *QOffset `json:"offset,omitempty"`
	Limit  *QLimit  `json:"limit,omitempty"`
}

// UpdateStateParams defines parameters for UpdateState.
type UpdateStateParams struct {
	State QChatState `json:"state"`
}

// AddWalletParams defines parameters for AddWallet.
type AddWalletParams struct {
	Wallet QWalletAddress `json:"wallet"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetAllTrades request
	GetAllTrades(ctx context.Context, params *GetAllTradesParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetPersonalData request
	GetPersonalData(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// InitPersonalData request
	InitPersonalData(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateState request
	UpdateState(ctx context.Context, chatID PChatID, params *UpdateStateParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTradesByChatID request
	GetTradesByChatID(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AddWallet request
	AddWallet(ctx context.Context, chatID PChatID, params *AddWalletParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAllTrades(ctx context.Context, params *GetAllTradesParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAllTradesRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetPersonalData(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetPersonalDataRequest(c.Server, chatID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) InitPersonalData(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewInitPersonalDataRequest(c.Server, chatID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateState(ctx context.Context, chatID PChatID, params *UpdateStateParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateStateRequest(c.Server, chatID, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetTradesByChatID(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTradesByChatIDRequest(c.Server, chatID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddWallet(ctx context.Context, chatID PChatID, params *AddWalletParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddWalletRequest(c.Server, chatID, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAllTradesRequest generates requests for GetAllTrades
func NewGetAllTradesRequest(server string, params *GetAllTradesParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trades")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if params.Offset != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "offset", runtime.ParamLocationQuery, *params.Offset); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	if params.Limit != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetPersonalDataRequest generates requests for GetPersonalData
func NewGetPersonalDataRequest(server string, chatID PChatID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "chatID", runtime.ParamLocationPath, chatID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewInitPersonalDataRequest generates requests for InitPersonalData
func NewInitPersonalDataRequest(server string, chatID PChatID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "chatID", runtime.ParamLocationPath, chatID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateStateRequest generates requests for UpdateState
func NewUpdateStateRequest(server string, chatID PChatID, params *UpdateStateParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "chatID", runtime.ParamLocationPath, chatID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/%s/state", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "state", runtime.ParamLocationQuery, params.State); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetTradesByChatIDRequest generates requests for GetTradesByChatID
func NewGetTradesByChatIDRequest(server string, chatID PChatID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "chatID", runtime.ParamLocationPath, chatID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/%s/trades", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAddWalletRequest generates requests for AddWallet
func NewAddWalletRequest(server string, chatID PChatID, params *AddWalletParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "chatID", runtime.ParamLocationPath, chatID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/%s/wallet", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "wallet", runtime.ParamLocationQuery, params.Wallet); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetAllTrades request
	GetAllTradesWithResponse(ctx context.Context, params *GetAllTradesParams, reqEditors ...RequestEditorFn) (*GetAllTradesResponse, error)

	// GetPersonalData request
	GetPersonalDataWithResponse(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*GetPersonalDataResponse, error)

	// InitPersonalData request
	InitPersonalDataWithResponse(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*InitPersonalDataResponse, error)

	// UpdateState request
	UpdateStateWithResponse(ctx context.Context, chatID PChatID, params *UpdateStateParams, reqEditors ...RequestEditorFn) (*UpdateStateResponse, error)

	// GetTradesByChatID request
	GetTradesByChatIDWithResponse(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*GetTradesByChatIDResponse, error)

	// AddWallet request
	AddWalletWithResponse(ctx context.Context, chatID PChatID, params *AddWalletParams, reqEditors ...RequestEditorFn) (*AddWalletResponse, error)
}

type GetAllTradesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *TradeList
}

// Status returns HTTPResponse.Status
func (r GetAllTradesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAllTradesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetPersonalDataResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PersonalData
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetPersonalDataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetPersonalDataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type InitPersonalDataResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r InitPersonalDataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r InitPersonalDataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateStateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r UpdateStateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateStateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTradesByChatIDResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *TradeList
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetTradesByChatIDResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTradesByChatIDResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AddWalletResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r AddWalletResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddWalletResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAllTradesWithResponse request returning *GetAllTradesResponse
func (c *ClientWithResponses) GetAllTradesWithResponse(ctx context.Context, params *GetAllTradesParams, reqEditors ...RequestEditorFn) (*GetAllTradesResponse, error) {
	rsp, err := c.GetAllTrades(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAllTradesResponse(rsp)
}

// GetPersonalDataWithResponse request returning *GetPersonalDataResponse
func (c *ClientWithResponses) GetPersonalDataWithResponse(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*GetPersonalDataResponse, error) {
	rsp, err := c.GetPersonalData(ctx, chatID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetPersonalDataResponse(rsp)
}

// InitPersonalDataWithResponse request returning *InitPersonalDataResponse
func (c *ClientWithResponses) InitPersonalDataWithResponse(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*InitPersonalDataResponse, error) {
	rsp, err := c.InitPersonalData(ctx, chatID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseInitPersonalDataResponse(rsp)
}

// UpdateStateWithResponse request returning *UpdateStateResponse
func (c *ClientWithResponses) UpdateStateWithResponse(ctx context.Context, chatID PChatID, params *UpdateStateParams, reqEditors ...RequestEditorFn) (*UpdateStateResponse, error) {
	rsp, err := c.UpdateState(ctx, chatID, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateStateResponse(rsp)
}

// GetTradesByChatIDWithResponse request returning *GetTradesByChatIDResponse
func (c *ClientWithResponses) GetTradesByChatIDWithResponse(ctx context.Context, chatID PChatID, reqEditors ...RequestEditorFn) (*GetTradesByChatIDResponse, error) {
	rsp, err := c.GetTradesByChatID(ctx, chatID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTradesByChatIDResponse(rsp)
}

// AddWalletWithResponse request returning *AddWalletResponse
func (c *ClientWithResponses) AddWalletWithResponse(ctx context.Context, chatID PChatID, params *AddWalletParams, reqEditors ...RequestEditorFn) (*AddWalletResponse, error) {
	rsp, err := c.AddWallet(ctx, chatID, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddWalletResponse(rsp)
}

// ParseGetAllTradesResponse parses an HTTP response from a GetAllTradesWithResponse call
func ParseGetAllTradesResponse(rsp *http.Response) (*GetAllTradesResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAllTradesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest TradeList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetPersonalDataResponse parses an HTTP response from a GetPersonalDataWithResponse call
func ParseGetPersonalDataResponse(rsp *http.Response) (*GetPersonalDataResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetPersonalDataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PersonalData
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseInitPersonalDataResponse parses an HTTP response from a InitPersonalDataWithResponse call
func ParseInitPersonalDataResponse(rsp *http.Response) (*InitPersonalDataResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &InitPersonalDataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseUpdateStateResponse parses an HTTP response from a UpdateStateWithResponse call
func ParseUpdateStateResponse(rsp *http.Response) (*UpdateStateResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateStateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetTradesByChatIDResponse parses an HTTP response from a GetTradesByChatIDWithResponse call
func ParseGetTradesByChatIDResponse(rsp *http.Response) (*GetTradesByChatIDResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetTradesByChatIDResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest TradeList
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseAddWalletResponse parses an HTTP response from a AddWalletWithResponse call
func ParseAddWalletResponse(rsp *http.Response) (*AddWalletResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddWalletResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /trades)
	GetAllTrades(w http.ResponseWriter, r *http.Request, params GetAllTradesParams)

	// (GET /{chatID})
	GetPersonalData(w http.ResponseWriter, r *http.Request, chatID PChatID)

	// (POST /{chatID})
	InitPersonalData(w http.ResponseWriter, r *http.Request, chatID PChatID)

	// (POST /{chatID}/state)
	UpdateState(w http.ResponseWriter, r *http.Request, chatID PChatID, params UpdateStateParams)

	// (GET /{chatID}/trades)
	GetTradesByChatID(w http.ResponseWriter, r *http.Request, chatID PChatID)

	// (POST /{chatID}/wallet)
	AddWallet(w http.ResponseWriter, r *http.Request, chatID PChatID, params AddWalletParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetAllTrades operation middleware
func (siw *ServerInterfaceWrapper) GetAllTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAllTradesParams

	// ------------- Optional query parameter "offset" -------------
	if paramValue := r.URL.Query().Get("offset"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------
	if paramValue := r.URL.Query().Get("limit"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAllTrades(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetPersonalData operation middleware
func (siw *ServerInterfaceWrapper) GetPersonalData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chatID" -------------
	var chatID PChatID

	err = runtime.BindStyledParameter("simple", false, "chatID", chi.URLParam(r, "chatID"), &chatID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chatID", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPersonalData(w, r, chatID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// InitPersonalData operation middleware
func (siw *ServerInterfaceWrapper) InitPersonalData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chatID" -------------
	var chatID PChatID

	err = runtime.BindStyledParameter("simple", false, "chatID", chi.URLParam(r, "chatID"), &chatID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chatID", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.InitPersonalData(w, r, chatID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// UpdateState operation middleware
func (siw *ServerInterfaceWrapper) UpdateState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chatID" -------------
	var chatID PChatID

	err = runtime.BindStyledParameter("simple", false, "chatID", chi.URLParam(r, "chatID"), &chatID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chatID", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params UpdateStateParams

	// ------------- Required query parameter "state" -------------
	if paramValue := r.URL.Query().Get("state"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "state"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "state", r.URL.Query(), &params.State)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "state", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateState(w, r, chatID, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetTradesByChatID operation middleware
func (siw *ServerInterfaceWrapper) GetTradesByChatID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chatID" -------------
	var chatID PChatID

	err = runtime.BindStyledParameter("simple", false, "chatID", chi.URLParam(r, "chatID"), &chatID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chatID", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTradesByChatID(w, r, chatID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// AddWallet operation middleware
func (siw *ServerInterfaceWrapper) AddWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "chatID" -------------
	var chatID PChatID

	err = runtime.BindStyledParameter("simple", false, "chatID", chi.URLParam(r, "chatID"), &chatID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "chatID", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params AddWalletParams

	// ------------- Required query parameter "wallet" -------------
	if paramValue := r.URL.Query().Get("wallet"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "wallet"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "wallet", r.URL.Query(), &params.Wallet)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "wallet", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddWallet(w, r, chatID, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/trades", wrapper.GetAllTrades)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/{chatID}", wrapper.GetPersonalData)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/{chatID}", wrapper.InitPersonalData)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/{chatID}/state", wrapper.UpdateState)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/{chatID}/trades", wrapper.GetTradesByChatID)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/{chatID}/wallet", wrapper.AddWallet)
	})

	return r
}
