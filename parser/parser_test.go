package parser_test

import (
	"testing"

	"github.com/html-link-parser/parser"
)

func TestGetResponseFor(t *testing.T) {
	p := parser.New()
	res, err := p.GetResponseFor("http://example.com")
	if err != nil {
		t.Errorf("Failed to get links err was not nil: %s.", err)
	}

	if res == "" {
		t.Errorf("Unexpected empty result %s.", res)
	}
}

// TESTING ANCHOR EXTRACTOR
func TestNestedText(t *testing.T) {
	// Should include all text inside children of anchor tags.
	s := `<a href="/dog">
          <span>Something in a span</span>
          Text not in a span
          <b>Bold text!</b>
        </a>`

	p := parser.New()

	err := p.ExtractAnchorTagsFrom(s)
	if err != nil {
		t.Errorf("Failed to extract anchor tags when testing nested text %s.", err)
	}

	entry := p.Links[0]
	if entry.Text != "Something in a span Text not in a span Bold text!" ||
		entry.Href != "/dog" {
		t.Errorf("Entry did not match expected result when checking nested text.")
	}
}

func TestExcludeCommentsInText(t *testing.T) {
	// Should exclude all comments inside of the text.
	s := `<html>
          <body>
            <a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
          </body>
        </html>`

	p := parser.New()

	err := p.ExtractAnchorTagsFrom(s)
	if err != nil {
		t.Errorf("Failed to extract anchor tags when testing excluded comments in text %s.", err)
	}

	entry := p.Links[0]
	if entry.Text != "dog cat" || entry.Href != "/dog-cat" {
		t.Errorf("Entry did not match expected result when checking excluded comments in text.")
	}
}

func TestMultipleAnchorTagsInDocument(t *testing.T) {
	// Should add all links when there is more than one.
	s := `<html>
        <head>
          <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
        </head>
        <body>
          <h1>Social stuffs</h1>
          <div>
            <a href="https://www.twitter.com/joncalhoun">
              Check me out on twitter
              <i class="fa fa-twitter" aria-hidden="true"></i>
            </a>
            <a href="https://github.com/gophercises">
              Gophercises is on <strong>Github</strong>!
            </a>
          </div>
        </body>
        </html>`

	p := parser.New()

	err := p.ExtractAnchorTagsFrom(s)
	if err != nil {
		t.Errorf("Failed to extract anchor tags when testing multiple anchor tags %s.", err)
	}

	var runTestsOn func([]parser.Link)
	runTestsOn = func(links []parser.Link) {
		if len(links) < 2 {
			t.Errorf("Unexpected result count when testing multiple anchor tags")
		}

		entry1, entry2 := links[0], links[1]
		if entry1.Href != "https://www.twitter.com/joncalhoun" || entry1.Text != "Check me out on twitter" {
			t.Errorf("Entry1 failed when testing multiple \n EXPECTED -> href:%s, text:%s \n GOT -> href:%s, text:%s",
				"https://www.twitter.com/joncalhoun", "Check me out on twitter", entry1.Href, entry1.Text)
		}

		if entry2.Href != "https://github.com/gophercises" || entry2.Text != "Gophercises is on Github !" {
			t.Errorf("Entry1 failed when testing multiple \n EXPECTED -> href:%s, text:%s \n GOT -> href:%s, text:%s",
				"https://github.com/gophercises", "Gophercises is on Github !", entry2.Href, entry2.Text)
		}
	}

	runTestsOn(p.Links)
}

func TestMultipleTagsWithComments(t *testing.T) {
	s := `<!DOCTYPE html>
  <!--[if lt IE 7]> <html class="ie ie6 lt-ie9 lt-ie8 lt-ie7" lang="en"> <![endif]-->
  <!--[if IE 7]>    <html class="ie ie7 lt-ie9 lt-ie8"        lang="en"> <![endif]-->
  <!--[if IE 8]>    <html class="ie ie8 lt-ie9"               lang="en"> <![endif]-->
  <!--[if IE 9]>    <html class="ie ie9"                      lang="en"> <![endif]-->
  <!--[if !IE]><!-->
  <html lang="en" class="no-ie">
  <!--<![endif]-->
  
  <head>
    <title>Gophercises - Coding exercises for budding gophers</title>
  </head>
  
  <body>
    <section class="header-section">
      <div class="jumbo-content">
        <div class="pull-right login-section">
          Already have an account?
          <a href="#" class="btn btn-login">Login <i class="fa fa-sign-in" aria-hidden="true"></i></a>
        </div>
        <center>
          <img src="https://gophercises.com/img/gophercises_logo.png" style="max-width: 85%; z-index: 3;">
          <h1>coding exercises for budding gophers</h1>
          <br/>
          <form action="/do-stuff" method="post">
            <div class="input-group">
              <input type="email" id="drip-email" name="fields[email]" class="btn-input" placeholder="Email Address" required>
              <button class="btn btn-success btn-lg" type="submit">Sign me up!</button>
              <a href="/lost">Lost? Need help?</a>
            </div>
          </form>
          <p class="disclaimer disclaimer-box">Gophercises is 100% FREE, but is currently in beta. There will be bugs, and things will be changing significantly over the coming weeks.</p>
        </center>
      </div>
    </section>
    <section class="footer-section">
      <div class="row">
        <div class="col-md-6 col-md-offset-1 vcenter">
          <div class="quote">
            "Success is no accident. It is hard work, perseverance, learning, studying, sacrifice and most of all, love of what you are doing or learning to do." - Pele
          </div>
        </div>
        <div class="col-md-4 col-md-offset-0 vcenter">
          <center>
            <img src="https://gophercises.com/img/gophercises_lifting.gif" style="width: 80%">
            <br/>
            <br/>
          </center>
        </div>
      </div>
      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          <center>
            <p class="disclaimer">
              Artwork created by Marcus Olsson (<a href="https://twitter.com/marcusolsson">@marcusolsson</a>), animated by Jon Calhoun (that's me!), and inspired by the original Go Gopher created by Renee French.
            </p>
          </center>
        </div>
      </div>
    </section>
  </body>
  </html>`

	p := parser.New()

	err := p.ExtractAnchorTagsFrom(s)
	if err != nil {
		t.Errorf("Failed to extract anchor tags when stress testing %s.", err)
	}

	var runTestsOn func([]parser.Link)
	runTestsOn = func(links []parser.Link) {
		if len(links) < 3 {
			t.Errorf("Unexpected result count when stress testing multiple anchor tags")
		}

		entry1, entry2, entry3 := links[0], links[1], links[2]
		if entry1.Href != "#" || entry1.Text != "Login" {
			t.Errorf("Entry1 failed when stress testing multiple \n EXPECTED -> href:%s, text:%s \n GOT -> href:%s, text:%s",
				"#", "Login", entry1.Href, entry1.Text)
		}

		if entry2.Href != "/lost" || entry2.Text != "Lost? Need help?" {
			t.Errorf("Entry2 failed when stress testing multiple \n EXPECTED -> href:%s, text:%s \n GOT -> href:%s, text:%s",
				"/lost", "Lost? Need help?", entry2.Href, entry2.Text)
		}

		if entry3.Href != "https://twitter.com/marcusolsson" || entry3.Text != "@marcusolsson" {
			t.Errorf("Entry3 failed when stress testing multiple \n EXPECTED -> href:%s, text:%s \n GOT -> href:%s, text:%s",
				"https://twitter.com/marcusolsson", "@marcusolsson", entry2.Href, entry2.Text)
		}
	}

	runTestsOn(p.Links)
}

func TestNoAnchorTagsPresent(t *testing.T) {
	// Should return an empty result.
	s := `<html>
          <body>
            <p>dog cat <!-- commented text SHOULD NOT be included! --></p>
            <div>
              <span>Something in a span</span>
              Text not in a span
            </div>
            <b>Bold text!</b>
          </body>
        </html>`

	p := parser.New()

	err := p.ExtractAnchorTagsFrom(s)
	if err != nil {
		t.Errorf("Failed to extract anchor tags when no anchor tags present %s.", err)
	}

	if len(p.Links) != 0 {
		t.Errorf("Unexpected result when testing no anchor tags present %d", len(p.Links))
	}
}
