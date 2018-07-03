#------ This program used to find degrees of separtion between two persons -------#
#------ Ruby proraming language used ----- #
#--- To run this use following commend --- ruby find_degree.rb ---#

require 'open-uri'
require 'json'
require 'open_uri_redirections'

def get_input

	puts "Enter two names with a space"
	inputs = gets.chomp.split(' ')

	if inputs.size < 2
		puts"Wrong number of inputs. Kindly try agin"
		get_input
	else
		actor1,actor2 = inputs[0],inputs[1]
	end # if inputs.size < 2

	#puts"--- #{actor1} --- #{actor2} ---"

    actor1_details = fetch_details(actor1)
    actor2_details = fetch_details(actor2)

	actor1_movies = movie_list(actor1)
	actor2_movies = movie_list(actor2)

	common_movie = actor1_movies & actor2_movies

	if common_movie.length != 0

		puts "Matched :#{common_movie[0]}\n"
        puts "\nDegrees of Separation: 1\n\n"

        movie = fetch_details(common_movie[0])

        movie["cast"] = [] if movie["cast"] == nil
		
		movie["crew"] = [] if movie["crew"] == nil
	
		cast1 = (movie["cast"] + movie["crew"]).select{|d| d if d['url'] == actor1}
		actor1_role = cast1[0]["role"]
	
		cast2 = (movie["cast"] + movie["crew"]).select{|d| d if d['url'] == actor2}
		actor2_role = cast1[0]["role"]

        puts "1.Movie : #{movie["name"]}\n"
        puts "#{actor1_role} : #{actor1}\n"
        puts "#{actor2_role} : #{actor2}\n\n"

        exit

	else

		find_next_degree(actor1_details,actor2_details,actor1,actor2)

	end # if common_movie.length != 0

end # def get_input

def find_next_degree(first_movie_list, second_movie_list,actor1,actor2)
  
  matched_cast_list = []
  
  first_movie_list['movies'].each do |first_movie_hash| 
    
    first_movie_cast_list = fetch_details(first_movie_hash['url'])
    if !first_movie_cast_list.empty?
      first_cast_list = first_movie_cast_list['cast']
      second_movie_list['movies'].each do |second_movie_hash| 
        
        second_movie_cast_list = fetch_details(second_movie_hash['url'])
        if !second_movie_cast_list.empty?
          second_cast_list = second_movie_cast_list['cast']
          first_degree_names = first_cast_list.collect{|f| f["name"]}
          second_degree_names =  second_cast_list.collect{|s| s["name"]}
          matched_cast_list = first_degree_names & second_degree_names
          if matched_cast_list.count > 0
            puts "Matched in both #{matched_cast_list}\n"
            first_seperation = first_cast_list.select {|c| c if c["name"] == matched_cast_list[0]}
            
            puts "\nDegrees of Separation: 2\n\n"

            puts "1.Movie : #{first_movie_hash['name']}\n"
            puts "#{first_movie_hash['role']} : #{first_movie_list['name']}\n"
            puts "#{first_seperation[0]['role']} : #{first_seperation[0]['name']}\n\n"

            puts "2.Movie : #{second_movie_hash['name']}\n"
            puts "#{first_seperation[0]['role']} : #{first_seperation[0]['name']}\n"
            puts "#{second_movie_hash['role']} : #{second_movie_list['name']}\n"

            exit
          end
        end
      end
    end
  end
  puts "\nNo Degree of Separation found for this 2 names"

end # def find_next_degree(actor1_details,actor2_details,actor1,actor2)

def fetch_details(name)

	url = "http://data.moviebuff.com/"
    
    begin
        data = JSON.load(open(url+name,  :allow_redirections => :safe))
    rescue
        data = {}
    end

end # def fetch_details(name)


def movie_list(actor)

	movie_details = fetch_details(actor)

	if movie_details['movies'] != nil	
		movie_details['movies'].map{|movie| movie['url']}.compact
	else
		{}
	end

end # def movie_list(actor1)

def actor_list(movie)

	actor_details = fetch_details(movie)

	actor_details["cast"] = [] if actor_details["cast"] == nil
		
	actor_details["crew"] = [] if actor_details["crew"] == nil
	
	(actor_details["cast"] + actor_details["crew"]).map{|actor| actor['url']}.compact

end # def actor_list(movie)

def call_output()
	puts"-- output --- #{@degree} ---"
end

get_input