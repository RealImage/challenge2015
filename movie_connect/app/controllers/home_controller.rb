class HomeController < ApplicationController
  require 'net/http'
  
  #Here i have used BFS techniq and crdeate a hash with all the display details if it matches
  #We can cache the search results to make it fast later searches.
   
  def index
    source = params[:a]
    @destination = params[:b]
    init
    success, path = try_with_degree(6)
    path["degrees_of_seperation"] = path.count
    render json: {result: "SUCCESS", :result => path, :found => success}.to_json, :status => 200
  end
  
  def try_with_degree(max_degree = nil)
    max_degree.times do
      flag, path = movies_iteration
      if flag == true
        return flag, path
      end
      person_iteration
    end
    return false, nil
  end
  
  def init
    @movies_array = Queue.new
    @persons_array = Queue.new
    @person_visited = {}
    @movies_visited = {}
    @search_count = 0
    details = get_list_from_moviebuff(params[:a])
    path = {"type" => "person", "url" => details["url"], "name" => details["name"]}
    details["movies"].each do |movie|
      @movies_array << {"path" => [path], "movie" => movie}
    end
  end
  
  
  def get_list_from_moviebuff(search = nil)
    url = URI.parse('http://data.moviebuff.com/'+search)
    req = Net::HTTP::Get.new(url.to_s)
    res = Net::HTTP.start(url.host, url.port) {|http|
      http.request(req)
    }
    if res.message == "OK"
      return JSON.parse(res.body)
    else
      return {}
    end
  end
  
  def movies_iteration
    while !@movies_array.empty?
      movie_obj = @movies_array.pop
      unless @movies_visited.key?([movie_obj["name"]])
        @movies_visited[movie_obj["name"]] = 1
        movie_details = get_list_from_moviebuff(movie_obj["movie"]["url"])
        if movie_details.key?("cast")
          movie_cast = movie_details["cast"]
          path = movie_obj["path"] + [{"type" => "movie", "url" => movie_obj["movie"]["url"], "name" => movie_obj["movie"]["name"]}]
          movie_cast.each do |person|
            if person["url"] == @destination
              return true, {"path" => path, "person" => person}
            else
              unless @person_visited.key?(person['name'])
                @persons_array << {"path" => path, "person" => person}
              end
            end
          end
        end
      end
    end
    return false, nil
  end
  
  def person_iteration
    while !@persons_array.empty?
      person_obj = @persons_array.pop
      unless @person_visited.key?([person_obj["name"]])
        @person_visited[person_obj["name"]] = 1
        person_details = get_list_from_moviebuff(person_obj["person"]["url"])
        if person_details.key?("movies")
          person_movies = person_details["movies"]
          path = person_obj["path"] + [{"type" => "person", "url" => person_obj["person"]["url"], "name" => person_obj["person"]["name"]}]
      
          person_movies.each do |movie|
            unless @movies_visited.key?(movie["name"])
              @movies_array << {"path" => path, "movie" => movie}
            end
          end
        end
      end
    end
  end
end
