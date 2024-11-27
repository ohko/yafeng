package auth

import "testing"

func TestAuth_EnAuth(t *testing.T) {
	o := &Auth{UserID: 1, IP: "127.0.0.1"}
	token, err := o.EnToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token:", token)
}

func TestAuth_DeAuth(t *testing.T) {
	o := &Auth{UserID: 0, IP: ""}
	if err := o.DeToken("47de78b78b0d8d41787aa8eff01e2d81a77e491ff58a5017b1b44676484fac0b1e8a63ef66d07f7bfc4b63f2bc"); err != nil {
		t.Fatal(err)
	}
	// t.Log(o.UserID)
}
