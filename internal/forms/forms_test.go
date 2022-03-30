package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest(http.MethodPost, "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when ir does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")

	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.MinLength("a", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existing field")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}


	postedValues = url.Values{}
	postedValues.Add("some_field", "some value")

	form = New(postedValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minlength of 100 met when data is shorter")
	}

	postedValues.Add("another_field", "abc")
	form = New(postedValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows minlenght of 1 is not match when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existing field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "a@c.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "a@c")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email")
	}
}
