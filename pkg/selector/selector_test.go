package selector

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"sync"
	"syscall"
	"testing"
)

func TestRace(t *testing.T) {
	fileName := "tmp.yml"
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	cfg := `
    websites:
        -   name: "test-example-1"
            url-template: "example1-%d"
            start-index: 0
            end-index: 5
    `
	ioutil.WriteFile(f.Name(), []byte(cfg), 0644)

	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(6)

	reqs := New(f.Name())
	for i := 0; i < 6; i++ {
		go func(w *sync.WaitGroup, ss *States) {
			defer w.Done()
			ss.Next()
		}(&wg, reqs)
	}
}

func TestOneStateWebsite(t *testing.T) {
	fileName := "tmp.yml"
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	cfg := `
    websites:
        -   name: "test-example-1"
            url-template: "example1-%d"
            start-index: 0
            end-index: 5
    `
	ioutil.WriteFile(f.Name(), []byte(cfg), 0644)

	reqs := New(f.Name())
	for i := 0; i < 6; i++ {
		req := reqs.Next()
		exp := &HttpRequest{
			Id:   i,
			Name: "test-example-1",
			Url:  fmt.Sprintf("example1-%d", i),
		}

		if !reflect.DeepEqual(exp, req) {
			t.Error("Expected: ", exp, " got: ", req)
		}
	}
}

func TestRoundRobin(t *testing.T) {
	fileName := "tmp2.yml"
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	cfg := `
    websites:
    -   name: "test-example-1"
    url-template: "example1-%d"
    start-index: 0
    end-index: 2

    -   name: "test-example-2"
    url-template: "example2-%d"
    start-index: 10
    end-index: 11
    `
	ioutil.WriteFile(f.Name(), []byte(cfg), 0644)

	reqs := New(f.Name())

	req := reqs.Next()
	exp := &HttpRequest{
		Id:   0,
		Name: "test-example-1",
		Url:  "example1-0",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = reqs.Next()
	exp = &HttpRequest{
		Id:   10,
		Name: "test-example-2",
		Url:  "example2-10",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = reqs.Next()
	exp = &HttpRequest{
		Id:   1,
		Name: "test-example-1",
		Url:  "example1-1",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = reqs.Next()
	exp = &HttpRequest{
		Id:   11,
		Name: "test-example-2",
		Url:  "example2-11",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = reqs.Next()
	exp = &HttpRequest{
		Id:   2,
		Name: "test-example-1",
		Url:  "example1-2",
	}

	if !reflect.DeepEqual(exp, req) {
		t.Error("Expected: ", exp, " got: ", req)
	}

	req = reqs.Next()
	if req != nil {
		t.Error("Request not nil")
	}
}
