# Slimgfast

Here's how you might use it:

    go get github.com/ericflo/slimgfast
    cd github.com/ericflo/slimgfast # Depends where your go env vars are set to
    go run main.go https://i.imgur.com

Now instead of going to:
    
    https://i.imgur.com/EgLrnVL.jpg

Now you can do:

    http://localhost:4400/EgLrnVL.jpg
    http://localhost:4400/EgLrnVL.jpg?w=300
    http://localhost:4400/EgLrnVL.jpg?w=300&h=600

It's easy to write your own fetchers, so instead of proxying, there's also a
built-in fetcher that pulls from S3.  Alternatively you could trivially write
a fetcher to pull from a directory on disk (actually, I'll probably write that
and include it in the package) or from anywhere that you're storing your
images.  More docs coming on how to do that.

Status: Very, very alpha.  I started with some code I'd written for work, but
it ended up being more or less a complete rewrite, and this version hasn't seen
any production traffic, ever.