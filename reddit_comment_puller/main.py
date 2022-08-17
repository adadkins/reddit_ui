import sys
import time
import praw
import os
from dotenv import load_dotenv
import psycopg2
from psycopg2 import OperationalError, errorcodes, errors

load_dotenv()

subreddits=os.environ["subreddits"]

conn = psycopg2.connect(
        host="localhost",
        database="postgres",
        user="postgres",
        password="postgres")

def main() -> None:
    print(conn.poll())

    reddit: praw.Reddit = praw.Reddit(
        client_id=os.environ["client_id"],
        client_secret=os.environ["client_secret"],
        user_agent=os.environ["user_agent"],
        username=os.environ["username"],
        password=os.environ["password"]
    )
    while True:
        print("Listening...")
        listen(reddit)
        print("Exiting...")

def listen(reddit: praw.Reddit) -> None:
    for comment in reddit.subreddit("NBA").stream.comments(skip_existing=True):
        #type prefixes
        # t1_	Comment
        # t2_	Account
        # t3_	Link
        # t4_	Message
        # t5_	Subreddit
        # t6_	Award
        # if 't1_' in comment.name:
        #     insert_comment(comment)
        if 't3_ ' in comment.name:
            insert_link(comment, reddit)
        if 't1_' in comment.name:
            insert_comment(comment, reddit)

def insert_link(submission: praw.reddit.models.Comment) -> None:
    cur = conn.cursor()
    postgres_insert_query = """ INSERT INTO submission (ID, BODY, AUTHOR, UPDATED_AT) VALUES (%s,%s,%s,%s)"""
    now = int(time.time())
    if submission.author == None:
        author_name= "Deleted"
    else: 
        author_name=submission.author.name
    record_to_insert = (submission.name, submission.selftext, author_name, now)
    cur.execute(postgres_insert_query, record_to_insert)
    try:
        conn.commit()
    except errors.InFailedSqlTransaction as err:
        print_psycopg2_exception(err)
        if cur is not None:
            cur.rollback()
        return 

def insert_comment(comment: praw.reddit.models.Comment, reddit: praw.Reddit) -> None:
        cur = conn.cursor()
        postgres_insert_query = """ INSERT INTO comments (ID, PARENT_ID, POST_ID, BODY, AUTHOR, UPDATED_AT) VALUES (%s,%s,%s,%s,%s,%s)"""
        now = int(time.time())
        if comment.author == None:
            author_name= "Deleted"
        else: 
            author_name=comment.author.name
        record_to_insert = (comment.name, comment.parent_id,comment.link_id, comment.body, author_name, now)
        cur.execute(postgres_insert_query, record_to_insert)
        try:
            conn.commit()
        except errors.InFailedSqlTransaction as err:
            print_psycopg2_exception(err)
            if cur is not None:
                cur.rollback()
            return
        update_parents(comment.parent_id, now, reddit)

#this could be improved with more clarity on the actual message bsides the pg code
def print_psycopg2_exception(err:Exception):
    # get details about the exception
    err_type, err_obj, traceback = sys.exc_info()

    # get the line number when exception occured
    line_num = traceback.tb_lineno

    # print the connect() error
    print ("\npsycopg2 ERROR:", err, "on line number:", line_num)
    print ("psycopg2 traceback:", traceback, "-- type:", err_type)

    # print the pgcode and pgerror exceptions
    print ("pgerror:", err.pgerror)
    print ("pgcode:", err.pgcode, "\n")

def update_parents(id: str, updated_at: int, reddit: praw.Reddit) -> None:
        while True:
            postgres_update_query = """ UPDATE comments SET UPDATED_AT = %s where ID = %s RETURNING PARENT_ID"""
            record_to_insert = (updated_at, id)

            cur = conn.cursor()
            cur.execute(postgres_update_query, record_to_insert)
            try:
                conn.commit()
            except errors.InFailedSqlTransaction as err:
                print_psycopg2_exception(err)
                if cur is not None:
                    cur.rollback()
                return
            count = cur.rowcount
            # if we do not have the parent in the database table, get it from PRAW
            if count == 0:
                comment_id = id.split("_")[1]
                try:
                    comment = reddit.comment(id=comment_id)
                    insert_comment(comment, reddit)
                    id=comment.parent_id 
                except Exception as ex:
                    print(ex)
                    return
            else:
                data = cur.fetchone()
                id = data[0]
            #print(id, 'now updating comment')

            #need to handle submissions(prefixed with t3_) differently.
            if 't3_' in id: 
                postgres_update_query = """ UPDATE submission SET UPDATED_AT = %s where ID = %s"""
                record_to_insert = (updated_at, id)
                cur = conn.cursor()
                cur.execute(postgres_update_query, record_to_insert)
                try:
                    conn.commit()
                except errors.InFailedSqlTransaction as err:
                    print_psycopg2_exception(err)
                    if cur is not None:
                        cur.rollback()
                    return
                count = cur.rowcount
                if count == 0:
                    #if we do not have the submission in the database get it from PRAW
                    submission_id = id.split("_")[1]
                    try:
                        submission = reddit.submission(id=submission_id)
                        insert_link(submission)
                    except Exception as ex:
                        print(ex)
                        return
                return
    
if __name__ == "__main__":
    main()