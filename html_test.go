package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	//<h1> Tag with text
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}

	//<h1> tag without text
	inputBody2 := "<html><body><h1></h1></body></html>"
	actual2 := getHeadingFromHTML(inputBody2)
	expected2 := ""

	if actual2 != expected2 {
		t.Errorf("expected %q, got %q", expected2, actual2)
	}

	//no <h1> tag but <h2> tag with text
	inputBody3 := "<html><body><h1></h1><h2>Test H2 tag</h2></body></html>"
	actual3 := getHeadingFromHTML(inputBody3)
	expected3 := "Test H2 tag"

	if actual3 != expected3 {
		t.Errorf("expected %q, got %q", expected3, actual3)
	}

}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	// test <p> present in <main>
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}

	// fallback if <p> in <main> is empty. Find first instance of <p>
	inputBody2 := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p></p>
		</main>
	</body></html>`
	actual2 := getFirstParagraphFromHTML(inputBody2)
	expected2 := "Outside paragraph."

	if actual2 != expected2 {
		t.Errorf("expected %q, got %q", expected2, actual2)
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	inputBody2 := `<html><body><a href="/blog/path"><span>Path</span></a></body></html>`

	actual2, err := getURLsFromHTML(inputBody2, baseURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
	}

	expected2 := []string{"https://crawler-test.com/blog/path"}
	if !reflect.DeepEqual(actual2, expected2) {
		t.Errorf("expected %v, got %v", expected2, actual2)
	}

	inputBody3 := `<html><body><a href="https://crawler-test.com/blog/path"><span>Path</span></a></body></html>`

	actual3, err := getURLsFromHTML(inputBody3, baseURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
	}
	expected3 := []string{"https://crawler-test.com/blog/path"}
	if !reflect.DeepEqual(actual3, expected3) {
		t.Errorf("expected %v, got %v", expected3, actual3)
	}

	inputBody4 := `<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`
	actual4, err := getURLsFromHTML(inputBody4, baseURL)
	if err != nil {
		t.Errorf("couldn't parse URL: %v", err)
	}
	expected4 := []string{"https://crawler-test.com/path/one", "https://other.com/path/one"}
	if !reflect.DeepEqual(actual4, expected4) {
		t.Errorf("expected %v, got %v", expected4, actual4)
	}

}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	inputBody2 := `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`

	actual2, err := getImagesFromHTML(inputBody2, baseURL)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual2, expected) {
		t.Errorf("expected %v, got %v", expected, actual2)
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
