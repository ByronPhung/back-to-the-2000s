package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Main Program

// In a more formal project, both the client and service would probably get instantiated once-and-only-once
// in a more global context, such as during service startup, so that it can be injected and shared across different services.
var userPostServiceImpl userPostService

func initialize() {
	userPostServiceImpl = userPostService{
		TypicodeClient: typicodeClient{
			// Right now, we don't care about any special customizations for the sake of this project,
			// but if we wanted to customize certain attributes, such as timeouts, redirect policies, etc.,
			// then it would make sense to inject a customized client
			//
			// Also, in a more formal project, the http.Client, typicodeClient, and userPostService would probably
			// get instantiated once-and-only-once in a more global context, such as during service startup, so that
			// they can be shared across different services.
			Client: http.DefaultClient,

			// If we were testing across multiple environments, then this would probably make more sense as
			// a config/environment variable.
			BaseUrl: "https://jsonplaceholder.typicode.com",
		},
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/v1/user-posts/:userId", getUserPostsByUserId)
	return router
}

func main() {
	initialize()
	router := setupRouter()
	router.Run("localhost:8080")
}

/*
	Controller Layer

	I still need to learn Go project organization standards, so I didn't do the actual folder layout
	for this project yet, but this section would represent the controller layer that is closest to the
	end user where the actual API status and response bodies are processed and returned to the consumer.

	Business logic should live in the service layer instead as the controller layer should just be dedicated
	to mapping service logic to API responses.
*/

func getUserPostsByUserId(c *gin.Context) {
	userId := c.Param("userId")

	// Validate input as expected ID type.
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Expected ID in integer format, but got '" + userId + "' instead"})
		return
	}

	userPostsResp, err := userPostServiceImpl.getUserPostsByUserId(userIdInt)

	// We have a defined value. Return a 200 with the JSON response.
	if !reflect.DeepEqual(userPostsResp, userPosts{}) {
		c.IndentedJSON(http.StatusOK, userPostsResp)
	} else if err == nil {
		// No explicit error. Treat this as a 404.
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Could not find userId=" + userId})
	} else {
		// Treat all other errors as 500s. Make sure we log it so that it can be troubleshooted in a live site environment too.
		//
		// Also, in general, in a live site environment, having monitors for general service 500 errors + alerts to page on-call
		// engineers if we have a large burst within a short period of time would be good.
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
}

/*
	Service Layer

	I still need to learn Go project organization standards, so I didn't do the actual folder layout
	for this project yet, but this section would represent the service layer that owns most of the
	internal business logic. The controller layer would primarily rely on this layer to do the
	"heavy lifting" so that the controller layer only has to worry about mapping the service layer
	to API status codes and their response bodies (if applicable).
*/

type userPostService struct {
	TypicodeClient typicodeClient
}

func (userPostService userPostService) getUserPostsByUserId(userId int) (userPosts, error) {
	var userResp user
	var userErr error
	var posts []postSummary
	var postsErr error

	// Run the API requests concurrently since they are independent and don't rely on each other's responses.
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		userResp, userErr = userPostService.TypicodeClient.getUserById(userId)
		waitGroup.Done()
	}()
	waitGroup.Add(1)
	go func() {
		posts, postsErr = userPostService.TypicodeClient.getPostsByUserId(userId)
		waitGroup.Done()
	}()
	waitGroup.Wait()

	// Prioritize getUserById response first since we can't return any relevant information if no
	// user info exists at all.
	if !reflect.DeepEqual(userResp, user{}) {
		// There is user data + no error from fetching posts
		if postsErr == nil {
			return userPosts{
				ID: userResp.ID,
				UserInfo: userInfo{
					Name:     userResp.Name,
					Username: userResp.Username,
					Email:    userResp.Email,
				},
				Posts: posts,
			}, nil
		} else {
			// Plan for 500 with any other error captured specifically while fetching posts
			return userPosts{}, postsErr
		}
	} else if userErr == nil {
		// Empty 'user' + no explicit error = 404 Not Found
		return userPosts{}, nil
	} else {
		// Plan for 500 with any other error captured specifically while fetching user info
		return userPosts{}, userErr
	}
}

// Clients - General

// Needed for injection in both functional code and unit tests.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Clients - typicodeClient
//
// Raw integration with Typicode API. Should only contain pure API integration logic with Typicode and
// should not host any business logic reserved for the service layer.

type typicodeClient struct {
	Client  httpClient
	BaseUrl string
}

// Fetch the general user info from Cool Vendor.
func (typicodeClient typicodeClient) getUserById(userId int) (user, error) {
	// Note that even though the expected ID type is an integer, the Typicode API can actually handle
	// any string and will just return a 404 with a generic empty JSON {} response, so we can save
	// ourselves from having to actually validate the user's input here.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprint(typicodeClient.BaseUrl, "/users/", userId), nil)
	if err != nil {
		return user{}, errors.New("Unexpected error creating client request for Cool Vendor's Get User API: error=" + err.Error())
	}

	// Execute request.
	resp, err := typicodeClient.Client.Do(req)
	if err != nil {
		return user{}, errors.New(fmt.Sprint("Unexpected communication or client policy error occurred trying to fetch userId=", userId, " from Cool Vendor: ", err.Error()))
	}
	defer resp.Body.Close()

	// Standard API contract by Typicode. If we get a 200 Ok from them, then it's considered successful
	// fetch of user data by ID.
	//
	// We should just make sure we explicitly handle the JSON parsing of the response body just to be safe.
	if resp.StatusCode == http.StatusOK {
		var userObj user
		if err := json.NewDecoder(resp.Body).Decode(&userObj); err != nil {
			return user{}, errors.New("Unable to parse response body as 'user' JSON for Cool Vendor's Get User By ID API: error=" + err.Error())
		}
		return userObj, nil
	} else if resp.StatusCode == http.StatusNotFound {
		// Cool Vendor API returns a 404 when a user could not found in their system, so we can rely
		// on that status code rather than trying to check for an empty {} object response.
		return user{}, nil
	} else {
		// Any non 200 or 404 is considered a general error that we should at least log.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return user{}, errors.New(fmt.Sprint("Unexpected error trying to read response body for server error trying to fetch userId=", userId, " from Cool Vendor: error=", err.Error()))
		}
		return user{}, errors.New(fmt.Sprint("Unexpected server error occurred trying to fetch userId=", userId, " from Cool Vendor: ", string(body)))
	}
}

// Fetch the posts for a given user ID.
func (typicodeClient typicodeClient) getPostsByUserId(userId int) ([]postSummary, error) {
	// Form request.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprint(typicodeClient.BaseUrl, "/posts"), nil)
	if err != nil {
		return []postSummary{}, errors.New("Unexpected error creating client request for Cool Vendor's Get Posts API: error=" + err.Error())
	}

	// Attach query params.
	q := req.URL.Query()
	q.Add("userId", fmt.Sprint(userId))
	req.URL.RawQuery = q.Encode()

	// Execute request.
	resp, err := typicodeClient.Client.Do(req)
	if err != nil {
		return []postSummary{}, errors.New(fmt.Sprint("Unexpected communication or client policy error occurred trying to fetch posts for userId=", userId, " from Cool Vendor: ", err.Error()))
	}
	defer resp.Body.Close()

	// Since Typicode's API contract will always return a 200 Ok for any valid request, including empty array []
	// for any requests where there are no matching user IDs, then it is safe to assume that we can just map any
	// 200 response from Typicode to the []postSummary value directly.
	if resp.StatusCode == http.StatusOK {
		var posts []postSummary
		if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
			return []postSummary{}, errors.New("Unable to parse response body as '[]postSummary' JSON for Cool Vendor's Get Posts API: error=" + err.Error())
		}
		return posts, nil
	} else {
		// Any non-200 response should be processed as an error, though, since we are not expecting it.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []postSummary{}, errors.New(fmt.Sprint("Unexpected error trying to read response body for server error trying to fetch posts for userId=", userId, " from Cool Vendor: error=", err.Error()))
		}
		return []postSummary{}, errors.New(fmt.Sprint("Unexpected server error occurred trying to fetch posts for userId=", userId, " from Cool Vendor: ", string(body)))
	}
}

/*
	Models

	I still need to learn Go project organization standards, so I didn't do the actual folder layout
	for this project yet, but this section would represent what would usually be a separate "models"
	package that individual models could be imported from as needed without bloating the main service
	layer and also promotes reusability across the entire project.
*/

/*
	Represents a user from Cool Vendor's Users API.

	Please note that this model is a bare-minimum version of their full response model
	because we still want to explicitly define a model for usage in our internal code,
	but we only care about consuming these subset of fields.

	@see https://coolvendor.com/api-docs/models/#user
*/
type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Represents a combination of relevant user info and their current posts.
type userPosts struct {
	ID       int           `json:"id"`
	UserInfo userInfo      `json:"userInfo"`
	Posts    []postSummary `json:"posts"`
}

// Represents a summary of user info to be used in "userPosts".
type userInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Represents a summary of raw post data to be used in "userPosts".
// Any user-associated data is removed for this model.
type postSummary struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
