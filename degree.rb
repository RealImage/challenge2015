require 'httparty'
require 'json'
require 'awesome_print'

def set_actors(actor1,actor2)

    first_degree_data = fetch_actor_data(actor1)
    second_degree_data = fetch_actor_data(actor2)
    
    if first_degree_data.empty? || second_degree_data.empty?
    	return
    end
    fd_names = first_degree_data["movies"].collect{|f| f["name"]}
    sd_names =  second_degree_data["movies"].collect{|s| s["name"]}

    matched_name = fd_names & sd_names

    if matched_name.length != 0
        first_seperation = first_degree_data["movies"].select {|c| c if c["name"] == matched_name[0]}
        second_seperation = second_degree_data["movies"].select{|s| s if s["name"] == matched_name[0] }
    end

    match_cast(first_degree_data, second_degree_data, actor1,actor2)

end

def fetch_actor_data(actor)
    base_url = "http://data.moviebuff.com/"
    begin
    		value = HTTParty.get "#{base_url}+#{actor}"
        actor_details = JSON.load(HTTParty.get base_url + actor)
    rescue
        actor_details = {}
    end

end

def match_cast(first_actor_movies_list, second_actor_movies_list,first_name, second_name)
  matched_cast_list = []
  first_actor_movies_list['movies'].each do |first_movie|
    first_movie_cast_list = fetch_details(first_movie['url'])
    if !first_movie_cast_list.empty?
      first_cast_list = first_movie_cast_list['cast']
      second_actor_movie_list['movies'].each do |second_movie|
        second_movie_cast_list = fetch_details(second_movie['url'])
        if !second_movie_cast_list.empty?
          second_cast_list = second_movie_cast_list['cast']
          first_degree_names = first_cast_list.collect{|f| f["name"]}
          second_degree_names =  second_cast_list.collect{|s| s["name"]}
          matched_cast_list = first_degree_names & second_degree_names
          if matched_cast_list.count > 0
            ap"Matched in both #{matched_cast_list}"
            first_seperation = first_cast_list.select {|c| c if c["name"] == matched_cast_list[0]}
            ap"Degrees of Separation: 2"

            ap"1.Movie : #{first_movie_hash['name']}"
            ap"#{first_movie['role']} : #{first_actor_movie_list['name']}"
            ap"#{first_seperation[0]['role']} : #{first_seperation[0]['name']}"

            ap"2.Movie : #{second_movie_hash['name']}"
            ap"#{first_seperation[0]['role']} : #{first_seperation[0]['name']}"
            ap"#{second_movie['role']} : #{second_actor_movie_list['name']}"

            exit
          end
        end
      end
    end
  end
  ap"No Degree of Separation"
end

ap"Enter two names"

actors = gets.chomp.split(' ')
first_actor =  actors[0]
second_actor =  actors[1]

output = set_actors(first_actor,second_actor)
if output.nil?
  ap "No Data found or Unable to access the data"
end



