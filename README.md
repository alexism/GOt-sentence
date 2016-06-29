


# Goals

Using Golang, create a web page with a single PPF sentence that refreshes on click. 


See http://burgundy.io/

http://agiliq.com/blog/2009/06/generating-pseudo-random-text-with-markov-chains-u/


# Impl Details

- A web server exposes an HTTP endpoint "GET /ppf".
- The web server has the corpus of words in memory.

- Iris will be used
http://iris-go.com/
