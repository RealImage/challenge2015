#!/usr/bin/ruby

BEGIN {
	ServiceName = File.basename($0)
	$0 = ServiceName
}

require 'json'
require 'zlib'
require 'uri'
require 'net/http'
require 'stringio'

class Time
	def to_ms
		(self.to_f * 1000.0).to_i
	end
end

class DegreeSeparation

	Person = Struct.new(:url, :name, :role, :ref) # It will hold movie list reference
	Movie  = Struct.new(:url, :name, :role, :ref) # It will hold person list reference
	DOS    = Struct.new(:movie, :role1, :name1, :role2, :name2) # It will hold degreeOfSeparation details

	def initialize(source, destination)
		@RootURL  = 'http://data.moviebuff.com/'
		@Source   = source.strip
		@Destination = destination.strip

		@VisitedPerson = Hash.new(0)
		@VisitedMovies = Hash.new(0)

		# Root Element
		@DataTree = Person.new(source,'root','root',nil)

		@degreeOfSeparation = []
		@needToBuildTree = []
		@totalNoRequests = 0
		@timeTaken = 0
	end

	def handleInput()
		if (@Source != @Destination)
			# Is this really a Person?

			# at very first time
			if (@DataTree.ref == nil)
				@needToBuildTree[0], @DataTree.ref = buildSubTree(@DataTree)
				return true
			end
		end
		return false
	end

	def getHTTPResponse(urlStr)
		begin
			@totalNoRequests += 1
			rStart = Time.now  # debug
			res = Net::HTTP.get_response(URI.parse("#{@RootURL}#{urlStr}"))
			rEnd = Time.now    # debug
			@timeTaken = @timeTaken + ( rEnd.to_ms - rStart.to_ms )
			if ( res.is_a?(Net::HTTPSuccess) && res.code == '200' )
				gz = Zlib::GzipReader.new(StringIO.new(res.body.to_s)) 
				jResponse = gz.read
				return jResponse
			end
		rescue => err
			puts "Error : Not able to fetch data for #{urlStr}"
			return false
		end
		return false
	end

	def getMovieDetails(movieName)
		personList = []

		return false if @VisitedMovies.has_key?(movieName)
		@VisitedMovies[movieName] += 1

		json = getHTTPResponse(movieName)
		return false if json == false

		jsonOut = JSON.parse(json)

		# "cast" details
		jsonOut['cast'].each do |h|
			personList << Person.new(h['url'], h['name'], h['role'],nil)
			if (h['url'] == @Destination)
				return false, personList 
			end
		end

		# "crew" details
		jsonOut['crew'].each do |h|
			personList << Person.new(h['url'], h['name'], h['role'],nil)
			if (h['url'] == @Destination)
				return false, personList 
			end
		end

		return true, personList
	end

	def getPersonDetails(personName)
		moviesList = []

		return false if @VisitedPerson.has_key?(personName)
		@VisitedPerson[personName] += 1

		json = getHTTPResponse(personName)
		return false if json == false

		jsonOut = JSON.parse(json)

		# "movies" details
		jsonOut['movies'].each do |h|
			moviesList << Movie.new(h['url'], h['name'], h['role'],nil)
		end

		return moviesList
	end

	def getPersonRoleInMovie(personURL, movieName)
		json = getHTTPResponse(personURL)
		return false if json == false

		jsonOut = JSON.parse(json)

		jsonOut['movies'].each do |h|
			if (movieName == h['name'])
				return jsonOut['name'], h['role']
			end
		end
	end

	def buildSubTree(person)
		if (movieList = getPersonDetails(person.url))
			movieList.each_with_index do |movie,idx|
				buildTree, personList = getMovieDetails(movie.url)
				if( personList )
					movie.ref = personList
					return buildTree, movieList if buildTree == false
				else
					movieList.delete_at(idx)
				end
			end
		end
		return true, movieList
	end

	def findSmallestDegree(treeStruct,depth=0)
		return false if treeStruct.ref == nil

		# accessing each movies one by one
		treeStruct.ref.each do |movieRef|
			next if movieRef.ref == nil

			# initializing the degree of separation
			@degreeOfSeparation[depth] = DOS.new(movieRef.name, nil, nil, nil, nil)

			# accessing each persons one by one
			movieRef.ref.each do |personRef|

				# Degrees of Separation
				if (depth > 0)
					@degreeOfSeparation[depth].role1 = @degreeOfSeparation[depth-1].role2
					@degreeOfSeparation[depth].name1 = @degreeOfSeparation[depth-1].name2
				else
					if (personRef.url == @Source)
						@degreeOfSeparation[depth].role1 = personRef.role
						@degreeOfSeparation[depth].name1 = personRef.name
					end
				end
				@degreeOfSeparation[depth].role2 = personRef.role
				@degreeOfSeparation[depth].name2 = personRef.name

				return true if (personRef.url == @Destination)

				if (personRef.ref == nil && @needToBuildTree[depth])
					@needToBuildTree[depth+1], personRef.ref = buildSubTree(personRef)
				elsif (personRef.ref != nil)
					return findSmallestDegree(personRef,depth+1)
				end
			end
		end
		return false
	end

	def printResult
		# if source person details is not found yet, then
		if( @degreeOfSeparation[0].role1 == nil or @degreeOfSeparation[0].name1 == nil )
			@degreeOfSeparation[0].name1, @degreeOfSeparation[0].role1 = getPersonRoleInMovie(@Source,@degreeOfSeparation[0].movie)
		end

		puts "Total no.of.request sent : #{@totalNoRequests}"
		puts "Total time taken (in ms) : #{@timeTaken}"
		puts "\nDegrees of Separation  : #{@degreeOfSeparation.size}\n\n"

		@degreeOfSeparation.each do |details|
			puts "Movie : #{details.movie}"
			puts "#{details.role1} : #{details.name1}"
			puts "#{details.role2} : #{details.name2}"
			puts "\n\n"
		end
	end

	def startDegreeOfSeparation
		if (! handleInput())
			puts "Usage : ruby #{$0} <PersonURL-1> <PersonURL-2>"
			exit(2)
		end

		while(true)
			if (findSmallestDegree(@DataTree))
				printResult
				break
			end
		end
	end
end

if (ARGV.size != 2)
	puts "Usage : ruby #{$0} <PersonURL-1> <PersonURL-2>"
	exit(2)
end

obj = DegreeSeparation.new(ARGV[0], ARGV[1])
obj.startDegreeOfSeparation
