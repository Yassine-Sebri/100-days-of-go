# Reading files (continued)

Let's extend our test further to extract the next line from the file, the description.

## Write the test first

```go
func TestNewBlogPosts(t *testing.T) {
 const (
  firstBody = `Title: Post 1
Description: Description 1`
  secondBody = `Title: Post 2
Description: Description 2`
 )

 fs := fstest.MapFS{
  "hello world.md":  {Data: []byte(firstBody)},
  "hello-world2.md": {Data: []byte(secondBody)},
 }

 // rest of test code cut for brevity
 assertPost(t, posts[0], blogposts.Post{
  Title:       "Post 1",
  Description: "Description 1",
 })

}
```

## Write enough code to make it pass

```go
func newPost(postFile io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postFile)

 scanner.Scan()
 titleLine := scanner.Text()

 scanner.Scan()
 descriptionLine := scanner.Text()

 return Post{Title: titleLine[7:], Description: descriptionLine[13:]}, nil
}
```

`bufio.Scanner` will help us scan the data, line by line.

## Refactor

We have repetition around scanning a line and then reading the text. We know we're going to do this operation at least one more time, it's a simple refactor to DRY up so let's start with that.

```go
func newPost(postFile io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postFile)

 readLine := func() string {
  scanner.Scan()
  return scanner.Text()
 }

 title := readLine()[7:]
 description := readLine()[13:]

 return Post{Title: title, Description: description}, nil
}
```

Whilst the magic numbers of 7 and 13 get the job done, they're not awfully descriptive.

```go
const (
 titleSeparator       = "Title: "
 descriptionSeparator = "Description: "
)

func newPost(postFile io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postFile)

 readLine := func() string {
  scanner.Scan()
  return scanner.Text()
 }

 title := readLine()[len(titleSeparator):]
 description := readLine()[len(descriptionSeparator):]

 return Post{Title: title, Description: description}, nil
}
```

We can use `strings.TrimPrefix` to remove the tags.

```go
func newPost(postBody io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postBody)

 readMetaLine := func(tagName string) string {
  scanner.Scan()
  return strings.TrimPrefix(scanner.Text(), tagName)
 }

 return Post{
  Title:       readMetaLine(titleSeparator),
  Description: readMetaLine(descriptionSeparator),
 }, nil
}
```

The next requirement is extracting the post's tags.

```go
func TestNewBlogPosts(t *testing.T) {
 const (
  firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go`
  secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker`
 )

 // rest of test code cut for brevity
 assertPost(t, posts[0], blogposts.Post{
  Title:       "Post 1",
  Description: "Description 1",
  Tags:        []string{"tdd", "go"},
 })
}
```

```go
const (
 titleSeparator       = "Title: "
 descriptionSeparator = "Description: "
 tagsSeparator        = "Tags: "
)

func newPost(postBody io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postBody)

 readMetaLine := func(tagName string) string {
  scanner.Scan()
  return strings.TrimPrefix(scanner.Text(), tagName)
 }

 return Post{
  Title:       readMetaLine(titleSeparator),
  Description: readMetaLine(descriptionSeparator),
  Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
 }, nil
}
```

We've read the first 3 lines already. We then need to read one more line, discard it and then the remainder of the file contains the post's body.

## Write the test

```go
 const (
  firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
  secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
 )
```

```go
 assertPost(t, posts[0], blogposts.Post{
  Title:       "Post 1",
  Description: "Description 1",
  Tags:        []string{"tdd", "go"},
  Body: `Hello
World`,
 })
```  

## Write code to make it pass

```go
func newPost(postBody io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postBody)

 readMetaLine := func(tagName string) string {
  scanner.Scan()
  return strings.TrimPrefix(scanner.Text(), tagName)
 }

 title := readMetaLine(titleSeparator)
 description := readMetaLine(descriptionSeparator)
 tags := strings.Split(readMetaLine(tagsSeparator), ", ")

 scanner.Scan() // ignore a line

 buf := bytes.Buffer{}
 for scanner.Scan() {
  fmt.Fprintln(&buf, scanner.Text())
 }
 body := strings.TrimSuffix(buf.String(), "\n")

 return Post{
  Title:       title,
  Description: description,
  Tags:        tags,
  Body:        body,
 }, nil
}
```

## Refactoring

```go
func newPost(postBody io.Reader) (Post, error) {
 scanner := bufio.NewScanner(postBody)

 readMetaLine := func(tagName string) string {
  scanner.Scan()
  return strings.TrimPrefix(scanner.Text(), tagName)
 }

 return Post{
  Title:       readMetaLine(titleSeparator),
  Description: readMetaLine(descriptionSeparator),
  Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
  Body:        readBody(scanner),
 }, nil
}

func readBody(scanner *bufio.Scanner) string {
 scanner.Scan() // ignore a line
 buf := bytes.Buffer{}
 for scanner.Scan() {
  fmt.Fprintln(&buf, scanner.Text())
 }
 return strings.TrimSuffix(buf.String(), "\n")
}
```
