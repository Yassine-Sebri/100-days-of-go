# Reading files

Pretend you're working with your friend to create some blog software. The idea is an author will write their posts in markdown, with some metadata at the top of the file. On startup, the web server will read a folder to create some `Post`s, and then a separate `NewHandler` function will use those `Post`s as a datasource for the blog's webserver.

We've been asked to create the package that converts a given folder of blog post files into a collection of `Post`s.

## Write the test first

```go
package blogposts_test

import (
 "testing"
 "testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
 fs := fstest.MapFS{
  "hello world.md":  {Data: []byte("hi")},
  "hello-world2.md": {Data: []byte("hola")},
 }

 posts := blogposts.NewPostsFromFS(fs)

 if len(posts) != len(fs) {
  t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
 }
}
```

Notice that the package of our test is `blogposts_test`. We don't want to test internal details because consumers don't care about them. By appending _test to our intended package name, we only access exported members from our package - just like a real user of our package.

We've imported `testing/fstest` which gives us access to the `fstest.MapFS` type. Our fake file system will pass `fstest.MapFS` to our package.

## Write enough code to make it pass

All we need to do is read the directory and create a post for each file we encounter. We don't have to worry about opening files and parsing them just yet.

```go
func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
 dir, _ := fs.ReadDir(fileSystem, ".")
 var posts []Post
 for range dir {
  posts = append(posts, Post{})
 }
 return posts
}
```

`fs.ReadDir` reads a directory inside a given `fs.FS` returning `[]DirEntry`.

## Refactor

Even though our tests are passing, we can't use our new package outside of this context, because it is coupled to a concrete implementation `fstest.MapFS`. But, it doesn't have to be. Change the argument to our `NewPostsFromFS` function to accept the interface from the standard library.

```go
func NewPostsFromFS(fileSystem fs.FS) []Post {
 dir, _ := fs.ReadDir(fileSystem, ".")
 var posts []Post
 for range dir {
  posts = append(posts, Post{})
 }
 return posts
}
```

### Error handling

We parked error handling earlier when we focused on making the happy-path work. Before continuing to iterate on the functionality, we should acknowledge that errors can happen when working with files. Beyond reading the directory, we can run into problems when we open individual files. Let's change our API (via our tests first, naturally) so that it can return an `error`.

```go
func TestNewBlogPosts(t *testing.T) {
 fs := fstest.MapFS{
  "hello world.md":  {Data: []byte("hi")},
  "hello-world2.md": {Data: []byte("hola")},
 }

 posts, err := blogposts.NewPostsFromFS(fs)

 if err != nil {
  t.Fatal(err)
 }

 if len(posts) != len(fs) {
  t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
 }
}
```

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
 dir, err := fs.ReadDir(fileSystem, ".")
 if err != nil {
  return nil, err
 }
 var posts []Post
 for range dir {
  posts = append(posts, Post{})
 }
 return posts, nil
}
```

Logically, our next iterations will be around expanding our `Post` type so that it has some useful data.

## Write the test

We'll start with the first line in the proposed blog post schema, the title field.

We need to change the contents of the test files so they match what was specified, and then we can make an assertion that it is parsed correctly.

```go
func TestNewBlogPosts(t *testing.T) {
 fs := fstest.MapFS{
  "hello world.md":  {Data: []byte("Title: Post 1")},
  "hello-world2.md": {Data: []byte("Title: Post 2")},
 }

 // rest of test code cut for brevity
 got := posts[0]
 want := blogposts.Post{Title: "Post 1"}

 if !reflect.DeepEqual(got, want) {
  t.Errorf("got %+v, want %+v", got, want)
 }
}
```

## Write code to make it pass

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
 dir, err := fs.ReadDir(fileSystem, ".")
 if err != nil {
  return nil, err
 }
 var posts []Post
 for _, f := range dir {
  post, err := getPost(fileSystem, f)
  if err != nil {
   return nil, err
  }
  posts = append(posts, post)
 }
 return posts, nil
}

func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
 postFile, err := fileSystem.Open(f.Name())
 if err != nil {
  return Post{}, err
 }
 defer postFile.Close()

 postData, err := io.ReadAll(postFile)
 if err != nil {
  return Post{}, err
 }

 post := Post{Title: string(postData)[7:]}
 return post, nil
}
```

`fs.FS` gives us a way of opening a file within it by name with its Open method. From there we read the data from the file and, for now, we do not need any sophisticated parsing, just cutting out the Title: text by slicing the string.

## Refactoring

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
 dir, err := fs.ReadDir(fileSystem, ".")
 if err != nil {
  return nil, err
 }
 var posts []Post
 for _, f := range dir {
  post, err := getPost(fileSystem, f.Name())
  if err != nil {
   return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
  }
  posts = append(posts, post)
 }
 return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
 postFile, err := fileSystem.Open(fileName)
 if err != nil {
  return Post{}, err
 }
 defer postFile.Close()
 return newPost(postFile)
}

func newPost(postFile io.Reader) (Post, error) {
 postData, err := io.ReadAll(postFile)
 if err != nil {
  return Post{}, err
 }

 post := Post{Title: string(postData)[7:]}
 return post, nil
}
```

We try to eliminate unnecessary coupling by using interfaces and more general types.
