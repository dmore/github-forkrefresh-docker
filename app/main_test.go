package main

import (
	"testing"
	//"os",
)

import (
	"os"
)

func Test1(t *testing.T) {
}

func Test_KEYCHAIN_APP_SERVICE(t *testing.T) {

	got := os.Getenv("KEYCHAIN_APP_SERVICE")
	want := "github-forkrefresh"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func Test_KEYCHAIN_USERNAME(t *testing.T) {

	got := os.Getenv("KEYCHAIN_USERNAME")
	want := "dmore"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGITHUB_TOKEN(t *testing.T) {

	//got := os.Getenv("GITHUB_TOKEN")

	//if (len(got) == 0){
	//	result = main.retrieve_secret_from_keychain("test", "test")
	//}

	//if (len(got) == 0) {
	//    t.Errorf("token not passed as env var")
	// }
}
