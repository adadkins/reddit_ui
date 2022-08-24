# ABOUT
Reddit is lacking a feature to tell which top level comments have been replied to most recently. 

This is a backend service to pull comments from a specified subreddit and track last updated by shoving them into a Postgres DB. 

A comment is updated when it has a child or any of its childrens children etc have a response.

The Postgres DB will allow us to build an API to return which comments have been replied to most recently.

# TODO
- Need to fix if post is a linked post. We might be able to check "is_self==false", then get  See praw submission object. (https://praw.readthedocs.io/en/latest/code_overview/models/submission.html#praw.models.Submission)

- Docker/dockercompose-ize services.

- There are a lot of "No Data returned for comment" from praw. Investigate cause. Its likely I'm not handling deleted or removed comments correctly.

- Clean up code. Error handling is duplicated. SQL statements and commits are duplicated. Could be abstracted and cleaned up.

- ~~fatal error if columns are not long enough for author or body~~ permanently fix temporary 5000 character limit patch


```

docker buildx build --push --platform linux/arm/v7,linux/amd64,linux/arm/v6 --tag ghcr.io/adadkins/reddit_comment_puller:latest .

```