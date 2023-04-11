# This project uses memcached to store usernames in cache, by doing this we can alleviate a large amount of queries in the database, because first we will search for the user in memcached, and only in case the user does not exist in Memcahed that we will search in the database of data
I made this simple example, but of course we can use this same structure to replicate in a larger project with more requests.

## Project made using:
GO (Programming Language)

> Gin (Web Framework)

> Gorm

> Memcached

### Here is an example of how memcached works in practice:

![Exemplo de GIF animado](http://g.recordit.co/pOVPIv14bl.gif)


#### Here's an explanation of what happened above:
First: We start the Go server
Second: We start the memcached server
Third: We insert a user into the database
Fourth: We can view our cached user by running the command "get user_43" (43 is the user ID, because in our program we cached our user ID)



