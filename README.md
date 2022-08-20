# ABOUT

Reddit lacks a feature to see which submissions and which comments have had replies most recently. This is an attempt to build a UI with that feature.

reddit comment puller   - pulls comments from a subreddit and puts them in a postgres database, and updates the parent comments timestamp so we can query by last replied to
reddit comment api      - two endpoints that return most commented on submissions, and and endpoint that returns the comments sorted by most recently replied to
reddit react app        - a front end to display submissions and their comments
