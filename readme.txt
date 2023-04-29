To install dependency -> npm i 
To run the code -> npm start
running port -> 4000

request to hit the :- 
Get http://localhost:4000/degree-of-separation?actor1=amitabh-bachchan&actor2=robert-de-niro

send actor1, actor2 in query

success response 200
{
    "message": "Fetched Successfully",
    "data": {
        "degreeOfSeparation": 219,
        "needToIncludeMovies": [
            "bollywood-the-greatest-love-story-ever-told",
            "ki-and-ka",
            "international-indian-film-awards",
            "gangaa-jamunaa-saraswathi",
            "naseeb",
            "krantiveer-the-revolution",
            "suhaag-1979-hindi",
            "pink-2016-hindi",
            "dafan",
            "toofan-1989-hindi",
            "zameer",
            "49th-manikchand-filmfare-awards-2003",
            "52nd-fair-one-filmfare-awards",
            "ahsaas",
            "chashme-buddoor",
            "charandas",
            "amrithadhare"
        ]
    }
} 

here 
    degreeOfSeparation ->  desired solution 
    needToIncludeMovies ->  while fetching data of some movies end point returns 403 for the 
                            including the response in the degreeOfSeparation and the array of the following 
                            movies has been mentioned in it.