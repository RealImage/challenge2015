require 'httparty'
require 'pry'

$filtered_artists = []

def log message
  puts message if $allow_puts
end

def get_movies_for(artist_name)
  fetch_details(artist_name).dig("movies")
end

def get_role(cast_and_crew, name)
  cast_and_crew.find { |cast| cast["url"] == name }
end

def fetch_details(key)
  log("Fetching details for = #{key}")

  return_value = {}
  begin
    response = HTTParty.get("http://data.moviebuff.com/#{key}")
    return_value = response.parsed_response
  rescue Exception => e
    puts e.message
    throw e # terminate the program, for now
  end
  return_value
end

def get_movie_crew(movie_url)
  movie_details = fetch_details(movie_url)
  return [ movie_details["cast"], movie_details["crew"] ].flatten.compact
end

# first_name.movies[i].casts includes second_name
def get_first_degree_relation(first_name, second_name)
  result = []
  movies_list = get_movies_for(first_name)
  begin
    movies_list.each do |movie|
      og_cast_and_crew = get_movie_crew(movie['url'])
      cast_and_crew = og_cast_and_crew - $filtered_artists
      cast_object = cast_and_crew.find { |cast| cast['url'] == second_name }
      if cast_object
        result = [{ movie: movie["name"], actor_1: get_role(og_cast_and_crew, first_name), actor_2: cast_object }]
        throw RuntimeError.new
      end
      $filtered_artists.push(cast_and_crew.map { |cast| cast['url'] })
      $filtered_artists.flatten!
    end
  rescue => e
    # do nothing
  end
  $filtered_artists.push(first_name)
  return result
end

def get_n_degree_relation(first_name, second_name)
  result = []
  if($level > 6)
    log("************* Level Exceeded 6 *************")
    return result
  end

  log("\n\nStarting Search for Artist: #{first_name} cast's cast's \n\n")

  # first_name.movies[i].cast[j].movies.cast[k] includes second_name
  first_movies = get_movies_for(first_name)
  begin
    first_movies.each do |movie|
      og_cast_and_crew = get_movie_crew(movie['url'])
      cast_and_crew = og_cast_and_crew - $filtered_artists
      inner_result = []
      cast_and_crew.each do |cast|
        inner_result = get_first_degree_relation(cast['url'], second_name)
        unless inner_result.empty?
          result = [{ movie: movie["name"], actor_1: get_role(og_cast_and_crew, first_name), actor_2: cast }, inner_result].flatten
          throw RuntimeError.new
        end
        $filtered_artists.push(cast['url'])
      end
    end
  rescue => e
    puts e.message
  end

  log("\n\nStarting Nested Search for Artist: #{first_name}\n\n")

  if result.empty?
    begin
      first_movies.each do |movie|
        cast_and_crew = get_movie_crew(movie['url'])
        inner_result = []
        cast_and_crew.each do |cast|
          $level += 1
          inner_result = get_n_degree_relation(cast['url'], second_name)
          unless inner_result.empty?
            result = [{ movie: movie["name"], actor_1: get_role(cast_and_crew, first_name), actor_2: cast }, inner_result].flatten
            throw RuntimeError.new
          end
          $level -= 1
        end
      end
    rescue => e
      puts e.message
    end
  end
  return result
end

def get_input
  print "Please enter 2 (unique) artist names, separated by a space: "
  input = gets.chomp
  first, second = input.split(/ +/)
  return first, second
end

def print_result(result)
  if result.empty?
    puts "\n\nNo degree of separation found!"
    return
  end

  puts "\n\nDegree of Separation: #{result.count}\n\n"
  result.each_with_index do |element, index|
    puts "#{index+1}. Movie: #{element.dig(:movie)}"
    puts "#{element.dig(:actor_1, 'role')} : #{element.dig(:actor_1, 'name')}"
    puts "#{element.dig(:actor_2, 'role')} : #{element.dig(:actor_2, 'name')}\n\n"
  end

end

# main
$allow_puts = ARGV[0]
ARGV.clear
$filtered_artists = []
$level = 1
first, second = get_input
result = get_first_degree_relation(first, second)
if result.empty?
  result = get_n_degree_relation(first, second)
end
print_result(result)