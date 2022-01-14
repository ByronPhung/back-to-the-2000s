package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// I'm traditionally used to using mocking libraries for everything, but the resources I referenced
// make it appear that Go developers tend to prefer just regular dependency injection rather than
// using mocks, so I stuck with that for now. However, having mocks would've substantially made
// it easier for me to set up my unit tests with true isolation as I test each of my components
// by layer (i.e. controller can more explicitly mock the services, service can mock the clients,
// clients can mock the HTTP requests, etc.).

// Controller - getUserPostsByUserId

func TestGetUserPostsByUserIdSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "userId", Value: fmt.Sprint(userId)}}

	userPostServiceImpl = userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			userBytes, userBytesErr := json.Marshal(testUser)
			assert.Nil(t, userBytesErr)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(userBytes)),
			}, nil
		} else {
			postsBytes, postsBytesErr := json.Marshal(posts)
			assert.Nil(t, postsBytesErr)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(postsBytes)),
			}, nil
		}
	}

	getUserPostsByUserId(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var userPostsResp userPosts
	userPostsRespErr := json.NewDecoder(w.Body).Decode(&userPostsResp)
	assert.Nil(t, userPostsRespErr)
	assert.Equal(t, userPostsResp, testUserPosts)
}

func TestGetUserPostsByUserId400(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "userId", Value: "test-123"}}

	getUserPostsByUserId(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedBody := "{\n    \"message\": \"Expected ID in integer format, but got 'test-123' instead\"\n}"
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestGetUserPostsByUserId404(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "userId", Value: fmt.Sprint(userId)}}

	userPostServiceImpl = userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		} else {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		}
	}

	getUserPostsByUserId(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	expectedBody := fmt.Sprint("{\n    \"message\": \"Could not find userId=", userId, "\"\n}")
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestGetUserPostsByUserId500(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "userId", Value: fmt.Sprint(userId)}}

	userPostServiceImpl = userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(strings.NewReader(errMsg500)),
			}, nil
		} else {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		}
	}

	getUserPostsByUserId(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expectedBody := fmt.Sprint("{\n    \"message\": \"Unexpected server error occurred trying to fetch userId=", userId, " from Cool Vendor: world.execute (me);\"\n}")
	assert.Equal(t, expectedBody, w.Body.String())
}

// userPostService.getUserPostsByUserId

func TestUserPostServiceGetUserPostsByIdSuccess(t *testing.T) {
	userPostService := userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			userBytes, userBytesErr := json.Marshal(testUser)
			assert.Nil(t, userBytesErr)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(userBytes)),
			}, nil
		} else {
			postsBytes, postsBytesErr := json.Marshal(posts)
			assert.Nil(t, postsBytesErr)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(postsBytes)),
			}, nil
		}
	}

	resp, respErr := userPostService.getUserPostsByUserId(userId)
	assert.NotNil(t, resp)
	assert.Nil(t, respErr)
	assert.Equal(t, testUserPosts, resp)
}

func TestUserPostServiceGetUserPostsByIdNoUserData(t *testing.T) {
	userPostService := userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			return &http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		} else {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		}
	}

	resp, respErr := userPostService.getUserPostsByUserId(userId)
	assert.Equal(t, resp, userPosts{})
	assert.Nil(t, respErr)
}

func TestUserPostServiceGetUserPostsByIdUserError(t *testing.T) {
	userPostService := userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(strings.NewReader(errMsg500)),
			}, nil
		} else {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		}
	}

	resp, respErr := userPostService.getUserPostsByUserId(userId)
	assert.Equal(t, resp, userPosts{})
	assert.NotNil(t, respErr)
	// We don't care about asserting the error value because that's tested in their typiscodeClient specs below instead.
}

func TestUserPostServiceGetUserPostsByIdPostsError(t *testing.T) {
	userPostService := userPostService{
		TypicodeClient: typicodeClient{
			Client:  &mockHTTPClient{},
			BaseUrl: mockBaseURL,
		},
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == fmt.Sprint(mockBaseURL, "/users/", userId) {
			userBytes, userBytesErr := json.Marshal(testUser)
			assert.Nil(t, userBytesErr)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(userBytes)),
			}, nil
		} else {
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(strings.NewReader(errMsg500)),
			}, nil
		}
	}

	resp, respErr := userPostService.getUserPostsByUserId(userId)
	assert.Equal(t, resp, userPosts{})
	assert.NotNil(t, respErr)
	// We don't care about asserting the error value because that's tested in their typiscodeClient specs below instead.
}

// typicodeClient.getUserById

func TestTypicodeClientGetUserByIdSuccess(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, r.URL.String(), fmt.Sprint(mockBaseURL, "/users/", userId))
		userBytes, userBytesErr := json.Marshal(testUser)
		assert.Nil(t, userBytesErr)
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(userBytes)),
		}, nil
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.NotNil(t, resp)
	assert.Nil(t, respErr)
	assert.Equal(t, testUser, resp)
}

func TestTypicodeClientGetUserByIdNewRequestErr(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: "   %#%badURL",
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.Equal(t, resp, user{})
	assert.NotNil(t, respErr)
	assert.Contains(t, respErr.Error(), "Unexpected error creating client request for Cool Vendor's Get User API: error=")
}

func TestTypicodeClientGetUserByIdResponseExecutionErr(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return nil, errFoo
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.Equal(t, resp, user{})
	assert.NotNil(t, respErr)
	assert.Equal(t, respErr.Error(), fmt.Sprint("Unexpected communication or client policy error occurred trying to fetch userId=", userId, " from Cool Vendor: ", errFoo.Error()))
}

func TestTypicodeClientGetUserById404(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 404,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}, nil
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.Equal(t, resp, user{})
	assert.Nil(t, respErr)
}

func TestTypicodeClientGetUserByIdBadJson(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id":"im-a-string"}`)),
		}, nil
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.Equal(t, resp, user{})
	assert.NotNil(t, respErr)
	assert.Contains(t, respErr.Error(), "Unable to parse response body as 'user' JSON for Cool Vendor's Get User By ID API: error=")
}

func TestTypicodeClientGetUserById500BadResponseBody(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}

	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       badReadCloser{},
		}, nil
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.Equal(t, resp, user{})
	assert.NotNil(t, respErr)
	assert.Equal(t, respErr.Error(), fmt.Sprint("Unexpected error trying to read response body for server error trying to fetch userId=", userId, " from Cool Vendor: error=", badReadCloserErrMsg))
}

func TestTypicodeClientGetUserById500(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}

	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(strings.NewReader(errMsg500)),
		}, nil
	}

	resp, respErr := typicodeClient.getUserById(userId)
	assert.Equal(t, resp, user{})
	assert.NotNil(t, respErr)
	assert.Equal(t, respErr.Error(), fmt.Sprint("Unexpected server error occurred trying to fetch userId=", userId, " from Cool Vendor: ", errMsg500))
}

// typicodeClient.getPostsByUserId

func TestTypicodeClientGetPostsByUserIdSuccess(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, r.URL.String(), fmt.Sprint(mockBaseURL, "/posts?userId=", userId))
		postsBytes, postsBytesErr := json.Marshal(posts)
		assert.Nil(t, postsBytesErr)
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(postsBytes)),
		}, nil
	}

	resp, respErr := typicodeClient.getPostsByUserId(userId)
	assert.NotNil(t, resp)
	assert.Nil(t, respErr)
	assert.Equal(t, posts, resp)
}

func TestTypicodeClientGetPostsByUserIdNewRequestErr(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: "   %#%badURL",
	}

	resp, respErr := typicodeClient.getPostsByUserId(userId)
	assert.Equal(t, resp, []postSummary{})
	assert.NotNil(t, respErr)
	assert.Contains(t, respErr.Error(), "Unexpected error creating client request for Cool Vendor's Get Posts API: error=")
}

func TestTypicodeClientGetPostsByUserIdResponseExecutionErr(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return nil, errFoo
	}

	resp, respErr := typicodeClient.getPostsByUserId(userId)
	assert.Equal(t, resp, []postSummary{})
	assert.NotNil(t, respErr)
	assert.Equal(t, respErr.Error(), fmt.Sprint("Unexpected communication or client policy error occurred trying to fetch posts for userId=", userId, " from Cool Vendor: ", errFoo.Error()))
}

func TestTypicodeClientGetPostsByUserIdBadJson(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}
	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id":"im-a-string"}`)),
		}, nil
	}

	resp, respErr := typicodeClient.getPostsByUserId(userId)
	assert.Equal(t, resp, []postSummary{})
	assert.NotNil(t, respErr)
	assert.Contains(t, respErr.Error(), "Unable to parse response body as '[]postSummary' JSON for Cool Vendor's Get Posts API: error=")
}

func TestTypicodeClientGetPostsByUserId500BadResponseBody(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}

	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       badReadCloser{},
		}, nil
	}

	resp, respErr := typicodeClient.getPostsByUserId(userId)
	assert.Equal(t, resp, []postSummary{})
	assert.NotNil(t, respErr)
	assert.Equal(t, respErr.Error(), fmt.Sprint("Unexpected error trying to read response body for server error trying to fetch posts for userId=", userId, " from Cool Vendor: error=", badReadCloserErrMsg))
}

func TestTypicodeClientGetPostsByUserId500(t *testing.T) {
	typicodeClient := typicodeClient{
		Client:  &mockHTTPClient{},
		BaseUrl: mockBaseURL,
	}

	mockHTTPClientDo = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(strings.NewReader(errMsg500)),
		}, nil
	}

	resp, respErr := typicodeClient.getPostsByUserId(userId)
	assert.Equal(t, resp, []postSummary{})
	assert.NotNil(t, respErr)
	assert.Equal(t, respErr.Error(), fmt.Sprint("Unexpected server error occurred trying to fetch posts for userId=", userId, " from Cool Vendor: ", errMsg500))
}

// Test Helpers
//

type mockHTTPClient struct{}

var mockHTTPClientDo func(req *http.Request) (*http.Response, error)

func (mockHTTPClient *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return mockHTTPClientDo(req)
}

type badReadCloser struct{}

func (badReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New(badReadCloserErrMsg)
}
func (badReadCloser) Close() error {
	return errors.New(badReadCloserErrMsg)
}

// Test Variables

var mockBaseURL = "https://cool-vendor.com"
var userId = 987654
var badReadCloserErrMsg = "THE WORLD IS OVER!"
var errMsg500 = "world.execute (me);"
var errFoo = errors.New("ooga-booga")

var testUser = user{
	ID:       userId,
	Name:     "Chacha",
	Username: "chacha22",
	Email:    "chacha22@gmail.com",
}
var posts = []postSummary{
	{ID: 42, Title: "How to Adult", Body: "N/A"},
}
var testUserPosts = userPosts{
	ID: testUser.ID,
	UserInfo: userInfo{
		Name:     testUser.Name,
		Username: testUser.Username,
		Email:    testUser.Email,
	},
	Posts: posts,
}
