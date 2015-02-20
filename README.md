#Degrees of Separation

With cinema going global these days, every one of the [A-Z]ollywoods are now connected. Use the wealth of data available at [Moviebuff](http://www.moviebuff.com) to see how. 

Write a Go program that behaves the following way:

```
$ degrees amitabh-bachchan robert-de-niro

Degrees of Separation: 3

1. Movie: The Great Gatsby
Supporting Actor: Amitabh Bachchan
Actor: Leonardo DiCaprio

2. Movie: The Wolf of Wall Street
Actor: Leonardo DiCaprio
Director: Martin Scorsese

3. Movie: Taxi Driver
Director: Martin Scorsese
Actor: Robert De Niro
```

Your solution should use the Moviebuff data available to figure out the smallest degree of separation between the two people. 

All the inputs should be Moviebuff URLs for their respective people: For Amitabh Bachchan, his page is on http://www.moviebuff.com/amitabh-bachchan and his Moviebuff URL is `amitabh-bachchan`.

Please do not attempt to scrape the Moviebuff website - All the data is available on an S3 bucket in an easy to parse JSON format here: https://data.moviebuff.com/<moviebuff_url>

To solve the example above, you solution would fetch at least the following:

https://data.moviebuff.com/amitabh-bachchan

https://data.moviebuff.com/the-great-gatsby

https://data.moviebuff.com/leonardo-dicaprio

https://data.moviebuff.com/the-wolf-of-wall-street

https://data.moviebuff.com/martin-scorsese

https://data.moviebuff.com/taxi-driver

##Submissions
Feel free to fork this repo and submit a pull request with your code. If you would like to use a private repo or hide your solution, do add `rajeshr` and `sudhirj` as collaborators and let us know at sudhir.j@moviebuff.com

###Judging Criteria
* Accuracy
* Efficiency
* Additional Features (UI / other options)



