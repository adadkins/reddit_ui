# ABOUT
API to return comments ordered by which have had a comment or child comment most recently replied to  ie which child has been updated last.

Two endpoints:

```bash
GET     /submissions/
GET     /submissions/t3_wq25kc/comments
```

Get posts ordered by latest comments. 

Get comments by latest response.

```

docker buildx build --push --platform linux/arm/v7,linux/amd64,linux/arm/v6 --tag ghcr.io/adadkins/reddit_comment_api:latest .

```