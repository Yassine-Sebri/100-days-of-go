# Templating

If we continue our journey of writing blog software, we want to generate two kinds of page:

1. **View post**. Renders a specific post. The Body field in Post is a string containing markdown so that should be converted to HTML.
2. **Index**. Lists all of the posts, with hyperlinks to view the specific post.

We'll also want a consistent look and feel across our site, so for each page we'll have the usual HTML furniture like `<html>` and a `<head>` containing links to CSS stylesheets and whatever else we may want.

We'll design our code so it accepts an `io.Writer`. This means the caller of our code has the flexibility to:

- Write them to an `os.File`, so they can be served statically.
- Write out the HTML directly to a `http.ResponseWriter`

## Write the test first

At this stage we're not overly concerned with the specific markup, and an easy first step would be just to check we can render the post's title as an `<h1>`. This feels like the smallest first step that can move us forward a bit.

```go
package blogrenderer_test

import (
 "bytes"
 "blogrenderer"
 "testing"
)

func TestRender(t *testing.T) {
 var (
  aPost = blogrenderer.Post{
   Title:       "hello world",
   Body:        "This is a post",
   Description: "This is a description",
   Tags:        []string{"go", "tdd"},
  }
 )

 t.Run("it converts a single post into HTML", func(t *testing.T) {
  buf := bytes.Buffer{}
  err := blogrenderer.Render(&buf, aPost)

  if err != nil {
   t.Fatal(err)
  }

  got := buf.String()
  want := `<h1>hello world</h1>`
  if got != want {
   t.Errorf("got '%s' want '%s'", got, want)
  }
 })
}
```

Our decision to accept an io.Writer also makes testing simple, in this case we're writing to a `bytes.Buffer` which we can then later inspect the contents.

## Write enough code to make it pass

```go
func Render(w io.Writer, p Post) error {
 _, err := fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
 return err
}
```

Now we have a very basic version working, we can now iterate on the test to expand on the functionality. In this case, rendering more information from the `Post`.

## Write the test

```go
 t.Run("it converts a single post into HTML", func(t *testing.T) {
  buf := bytes.Buffer{}
  err := blogrenderer.Render(&buf, aPost)

  if err != nil {
   t.Fatal(err)
  }

  got := buf.String()
  want := `<h1>hello world</h1>
<p>This is a description</p>
Tags: <ul><li>go</li><li>tdd</li></ul>`

  if got != want {
   t.Errorf("got '%s' want '%s'", got, want)
  }
 })
```

Notice that writing this, feels awkward. Seeing all that markup in the test feels bad, and we haven't even put the body in, or the actual HTML we'd want with all of the `<head>` content and whatever page furniture we need.

## Write code to make it pass

```go
func Render(w io.Writer, p Post) error {
 _, err := fmt.Fprintf(w, "<h1>%s</h1><p>%s</p>", p.Title, p.Description)
 if err != nil {
  return err
 }

 _, err = fmt.Fprint(w, "Tags: <ul>")
 if err != nil {
  return err
 }

 for _, tag := range p.Tags {
  _, err = fmt.Fprintf(w, "<li>%s</li>", tag)
  if err != nil {
   return err
  }
 }

 _, err = fmt.Fprint(w, "</ul>")
 if err != nil {
  return err
 }

 return nil
}
```

## Refactor

We can use templating. Here is a template for our blog:

`<h1>{{.Title}}</h1><p>{{.Description}}</p>Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`

Where do we define this string? Well, we have a few options, but to keep the steps small, let's just start with a plain old string

```go
package blogrenderer

import (
 "html/template"
 "io"
)

const (
 postTemplate = `<h1>{{.Title}}</h1><p>{{.Description}}</p>Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`
)

func Render(w io.Writer, p Post) error {
 templ, err := template.New("blog").Parse(postTemplate)
 if err != nil {
  return err
 }

 if err := templ.Execute(w, p); err != nil {
  return err
 }

 return nil
}
```

We create a new template with a name, and then parse our template string. We can then use the `Execute` method on it, passing in our data, in this case the `Post`.

The template will substitute things like `{{.Description}}` with the content of `p.Description`. Templates also give you some programming primitives like `range` to loop over values, and `if`.

### More refactoring

Using the `html/template` has definitely been an improvement, but having it as a string constant in our code isn't great:

- It's still quite difficult to read.
- It's not IDE/editor friendly. No syntax highlighting, ability to reformat, refactor, etc.

What we'd like to do is have our templates live in separate files so we can better organise them, and work with them as if they're HTML files.

```go
package blogrenderer

import (
 "embed"
 "html/template"
 "io"
)

var (
 postTemplates embed.FS
)

func Render(w io.Writer, p Post) error {
 templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
 if err != nil {
  return err
 }

 if err := templ.Execute(w, p); err != nil {
  return err
 }

 return nil
}
```

By embedding a "file system" into our code, we can load multiple templates and combine them freely. This will become useful when we want to share rendering logic across different templates, such as a header for the top of the HTML page and a footer.

Why would we want to use this? Well the alternative is that we can load our templates from a "normal" file system. However this means we'd have to make sure that the templates are in the correct file path wherever we want to use this software.

We don't really want our template to be defined as a one line string. We want to be able to space it out to make it easier to read and work with, something like this:

```html
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
```
