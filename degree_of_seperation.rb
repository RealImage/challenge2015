require 'open-uri'
require 'json'
require 'open_uri_redirections'

#to run this install open-uri json open_uri_redirections gem on your local
#then run command ruby degree_of_seperation.rb

def find_degree
  puts "Input two actors name with small letters and separated by -. example:- amitabh-bachchan robert-de-niro\n"
  puts "Enter First Name\n"
  first_name = gets().chomp
  puts "Enter Last Name\n"
  last_name = gets().chomp
  first_name_records = actors_data(first_name)
  last_name_records = actors_data(last_name)
  puts "loading ."
  first_costars = find_shortest_path(first_name_records["movies"],last_name)
  last_costars = find_shortest_path(last_name_records["movies"],first_name)

  common_costars = first_costars.keys & last_costars.keys
  costars_movie = first_costars["#{common_costars.first}"]
  movie_name = movies_data(costars_movie)
  if common_costars
    puts "\nDegree of separation: 2"
    puts "Actor #{first_name_records["name"]}"
    puts "Movie: #{movie_name["name"]}\n\n"

    puts "Movie: #{movie_name["name"]}"
    puts "Actor: #{common_costars.first}\n\n"


    puts "Actor #{last_name_records["name"]}"
    puts "Movie: #{movie_name["name"]}\n"
  else
    puts "Degree of separation: 3"
  end
end

def find_shortest_path(movies,actor_name,costars_info={})
  movies.each do |movie|
    movie_record = movies_data(movie["url"])
    if movie_record["cast"]
      movie_record["cast"].each do |r|
        putc "."

        if r["url"] == actor_name
          puts "\nDegree of separation: 1"
          puts "Movie: #{movie["name"]}"
          puts "#{r["role"]}: #{r["name"]}"
          exit
        end
        costars_info.merge!("#{r["url"]}" => "#{movie["url"]}")
      end
    end
  end
  costars_info
end  

def movies_data(url)
  begin
    base_url = "http://data.moviebuff.com/"
    url_data = URI.open(base_url+url,  :allow_redirections => :safe)
      fetched_data = JSON.load(url_data)
  rescue
      fetched_data = {}
  end
end

def actors_data(url)
  begin
    base_url = "http://data.moviebuff.com/"
    url_data = URI.open(base_url+url,  :allow_redirections => :safe)
      fetched_data = JSON.load(url_data)
  rescue
      fetched_data = {}
  end
end 

find_degree 
