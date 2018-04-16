require 'httparty'
require 'hashie'
module DegreeOfSeparation
	include HTTParty
	base_uri 'http://data.moviebuff.com'

	def find_degree_of_sep
		puts "Enter Data"
		actors = gets.split(' ')
		return "Please give two actor names" unless actors[1] && actors[2]#ensure two actor names are given as input
		return  "Degrees of Separation: 0" if actors[1] == actors[2]
		@path = []
		if get_movies(actors[1]).size < get_movies(actors[2]).size
			actor2 = actors[2] 
			resp = find_degree(actors[1],actors[2])
		else
			actor2 = actors[1] 
			resp = find_degree(actors[2],actors[1])
		end
	 	(resp.is_a? TrueClass) ? get_op(actor2) : "Actors are not related"
	end

	def find_degree(actor1,actor2)
		@all_movies = []
		@all_actors = []
		actor_set = [actor1]
		@all_actors << actor1
		5.times{|degree|
			actor_set1 = actor_set.uniq
			actor_set = []
			actor_set1.each{|actor_1|
				result = degrees(actor_1, actor2,degree)
				if (result.is_a? TrueClass)
				 	find_degree(actor1,actor_1) if degree != 0
				 	return true
				else
				 	actor_set.concat(result)
				end
			}
		}
	end

	def degrees(actor1,actor2,degree)
		movies = get_movies(actor1)
		@all_movies.concat(movies)
		res = []
		movies.each{|movie|
			@path[degree] = {actor1 => movie}
			result = get_actors(movie)
			next if result.none?
			return true if result.include? actor2
			res << result
			@all_actors.concat(result)
		}
		res.flatten.uniq
	end

	def get_op(actor2)
		path = @path.reduce({}, :merge)
		path[actor2] = nil
		puts "Degrees of Separation: #{@path.size} \n"
		path.values.compact.reverse.each_with_index{|movie, index|
			puts "\nMovie: #{movie}"
			data = get_data(movie)
			(data["cast"] + data["crew"]).each{|crew|
				if crew["url"] == path.keys.reverse[index]
					puts "#{crew["role"]} : #{crew["name"]}"
					break
				end
			}
			(data["cast"] + data["crew"]).each{|crew|
				if crew["url"] == path.keys.reverse[index+1]
					puts "#{crew["role"]} : #{crew["name"]}"
					break
				end
			}
		}
		nil
	end

	def get_movies(actor)
		result = get_data(actor) || {"movies"=>[]}
		result['movies'].map{|movie| movie['url'] unless !@all_movies.nil? && @all_movies.include?(movie['url'])}.compact
	end

	def get_actors(movie)
		result = get_data(movie) || {"cast"=>[],"crew"=>[]}
		(result["cast"] + result["crew"]).map{|actor| actor['url'] unless !@all_movies.nil? && @all_actors.include?(actor['url'])}.compact
	end

	def get_data(inp)
		begin
			res = DegreeOfSeparation.get('/' + inp.to_s)
			res = nil if res['Error']
		rescue => e
			p e.message
		end
		res 
	end
end