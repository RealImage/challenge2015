#Setup :

1. Open terminal and set the gopath (just above gophercon)

	example - export GOPATH=/home/mypc/work/projects/gopath/

2. Make sure limit of open files is high enough (>50000) for better performance.

	ubuntu os: http://askubuntu.com/questions/162229/how-do-i-increase-the-open-files-limit-for-a-non-root-user

3. Make sure you have fast internet access for better performance.

#Building the application

1. Go inside degree directory

	example - cd /home/mypc/work/projects/gopath/src/gophercon/degrees

2. build and install

	go get

	go build

	go install

3. run the application

	USASE: degrees <first-person-name><space><second-person-name>

	Example: degrees Amitabh-Bachchan Robert-De-Niro

#NOTE :

1. If you get any tcp error then decrease the rate-count value in the config file conf.json.

#Examples
```
$ ./degrees Jennifer-Lawrence meryl-streep
Time Taken:  855.773733ms

Degree of saperation:  2

1. Movie: Julie & Julia
Actress: Meryl Streep
Supporting Actor: Stanley Tucci

2. Movie: The Hunger Games: Mockingjay - Part 1
Supporting Actor: Stanley Tucci
Actress: Jennifer Lawrence


$ ./degrees Amitabh-Bachchan Kristen-Wiig
Time Taken:  3.489884099s

Degree of saperation:  3

1. Movie: The Secret Life of Walter Mitty
Actress: Kristen Wiig
Supporting Actor: Gurdeep Singh

2. Movie: Fateh

Executive Producer: Gurdeep Singh
Supporting Actress: Navneet Nishan

3. Movie: Kalpvriksh
Supporting Actor: Navneet Nishan
Actor: Amitabh Bachchan


$ ./degrees winfield-scott-mattraw welker-white
Time Taken:  19.760508061s

Degree of saperation:  7

1. Movie: Snow White and the Seven Dwarfs
Supporting Actor: Winfield Scott Mattraw
Supporting Actor: William Gilbert Barron

2. Movie: The Great Dictator
Supporting Actor: William Gilbert Barron
Editor / Editorial: Willard Nico

3. Movie: City Lights
Editor: Willard Nico
Music Director / Music: Carl Davis

4. Movie: The General
Music Director: Carl Davis
Editor / Editorial: Sherman Kell

5. Movie: The Wind Rises
Music Director: Joe Hisaishi
Supporting Actor: Stanley Tucci

6. Movie: Transformers: Age Of Extinction
Actor: Stanley Tucci
Compositor / Visual Effects: Brian N Bentley

7. Movie: The Wolf of Wall Street
Supporting Actor: Kenneth Choi
Supporting Actress: Welker White
```
