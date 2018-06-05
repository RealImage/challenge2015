# *********************************   ruby version: ruby 2.4.1p111 ************************************************
# ********************************* To run the program --> ruby degree_of_seperation.rb ***************************
# ********************************* Enter 2 vallues --> christian-bale tom-hardy ********************************** 

require 'httparty'
# require 'uri'
# require 'net/http'
require 'json'

def get_input
  puts "Enter 2 actor/actress names"

  actors = gets.chomp.split(' ')
  first_actor =  actors[0]
  second_actor =  actors[1]

  if first_actor.nil? || second_actor.nil?
    puts "Please enter 2 names."
    exit
  end

  if first_actor == second_actor
    puts "Names are same. Please enter different names !!"
    exit
  end
  return first_actor, second_actor
end

def degree_of_seperation
  actor1, actor2 = get_input
  first_actor_data = fetch_details(actor1)
  second_actor_data = fetch_details(actor2)
  
  if first_actor_data.empty? || second_actor_data.empty?
  	puts " Dont have sufficient data to proceed!!"
    exit
  end
  first_actor_movie_names = first_actor_data["movies"].collect{|f| f["name"]}
  second_actor_movie_names =  second_actor_data["movies"].collect{|s| s["name"]}

  matched_movie_names = first_actor_movie_names & second_actor_movie_names
 
  if matched_movie_names.length == 0
    #indirect match?
    find_and_print_indirect_degree_of_seperation first_actor_data, second_actor_data, actor1, actor2
  else
    #direct match
    matched_movie_name = matched_movie_names.first
    first_actor_details = matched_actor_movie_details first_actor_data, matched_movie_name, actor1
    second_actor_details = matched_actor_movie_details second_actor_data, matched_movie_name, actor2
    print_direct_degree_of_seperation matched_movie_name, first_actor_details, second_actor_details
  end
end



def print_direct_degree_of_seperation movie, detail_1, detail_2
  puts "Degree of Seperation: 1"
  puts "Movie: #{movie}"
  puts "#{detail_1['details'].first['role']} : #{detail_1['actor']}"
  puts "#{detail_2['details'].first['role']} : #{detail_2['actor']}" 
end

def fetch_details(item)
  base_url = "http://data.moviebuff.com/"
  begin
  		value = HTTParty.get "#{base_url}#{item}"
      details = JSON.parse(value.body)
    rescue Exception => e  
     	puts "#{e.message}"
      details = {}
    end

	# begin
	# 	puts "1233"
	#    url = URI("#{base_url}#{val}")
	#    http = Net::HTTP.new(url.host, url.port)
	#    http.use_ssl = true
	#    http.verify_mode = OpenSSL::SSL::VERIFY_NONE
	# 	request = Net::HTTP::Get.new(url, 'Content-Type' => 'puts plication/json')

	# 	request["cache-control"] = 'no-cache'
	# 	puts "1111111"
	# 	puts "#{http.request(request)}"
	# details = JSON.parse(http.request(request).read_body)

	# rescue
	# 	details = {}
	# end
end

def find_and_print_indirect_degree_of_seperation first_actor_data, second_actor_data, first_name, second_name
  matched_cast_list = []
  first_actor_data['movies'].each do |movie1|
    movie1_data = fetch_details(movie1['url'])

    unless movie1_data.empty?
      movie1_cast_list = movie1_data['cast']

      second_actor_data['movies'].each do |movie2|
        movie2_data = fetch_details(movie2['url'])
        unless movie2_data.empty?
          movie2_cast_list = movie2_data['cast']
        
          movie1_cast_names = movie1_cast_list.collect{|m| m["name"]}
          movie2_cast_names =  movie2_cast_list.collect{|m| m["name"]}

          matched_cast_list = movie1_cast_names & movie2_cast_names

          if matched_cast_list.count > 0
             puts "Matching cast list is #{matched_cast_list}"
             comman_actor = movie1_cast_list.select {|c| c if c["name"] == matched_cast_list.first}
             print_indirect_degree_of_seperation movie1, movie2, comman_actor.first, second_actor_data
             
             exit
          end
        end
      end
    end
  end
  puts "No Degree of Separation"
  exit
end

def print_indirect_degree_of_seperation movie1, movie2, comman_actor, second_actor_data
  puts " comman actor data is -----> #{comman_actor}"
  puts "Degrees of Separation: 2"
  puts "1.Movie : #{movie1['name']}"
  puts "#{movie1['role']} : #{comman_actor['name']}"
  puts "#{comman_actor['role']} : #{comman_actor['name']}"
  puts "2.Movie : #{movie2['name']}"
  puts "#{comman_actor['role']} : #{comman_actor['name']}"
  puts "#{movie2['role']} : #{second_actor_data['name']}"
end

def matched_actor_movie_details data_set, matched_movie_name, actor
  res = {}
  res["actor"] = actor
  res["details"] = data_set["movies"].select {|d| d if d["name"] == matched_movie_name }
  res
end


#call the main program
degree_of_seperation
