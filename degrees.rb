require 'open-uri'
require 'json'
require 'open_uri_redirections'

def assign_values(first,second)

    first_degree_data = fetch_details(first)
    second_degree_data = fetch_details(second)
    if first_degree_data.empty? || second_degree_data.empty?
      return
    end

    first_degree_names = first_degree_data["movies"].collect{|f| f["name"]}
    second_degree_names =  second_degree_data["movies"].collect{|s| s["name"]}

    matched_name = first_degree_names & second_degree_names

    if matched_name.length != 0
        puts "Matched :#{matched_name}\n"
        puts "\nDegrees of Separation: 1\n\n"

        first_seperation = first_degree_data["movies"].select {|c| c if c["name"] == matched_name[0]}
        second_seperation = second_degree_data["movies"].select{|s| s if s["name"] == matched_name[0] }

        puts "1.Movie : #{first_seperation[0]["name"]}\n"
        puts "#{first_seperation[0]["role"]} : #{first_degree_data["name"]}\n"
        puts "#{second_seperation[0]["role"]} : #{second_degree_data["name"]}\n\n"

        puts "2.Movie : #{second_seperation[0]["name"]}\n"
        puts "#{second_seperation[0]["role"]} : #{second_degree_data["name"]}\n"
        puts "#{first_seperation[0]["role"]} : #{first_degree_data["name"]}\n\n"
        exit
    end

    match_first_second_order_cast(first_degree_data, second_degree_data, first,second)

end

def fetch_details(url)
    base_url = "http://data.moviebuff.com/"
    begin
        fetched_data = JSON.load(open(base_url+url,  :allow_redirections => :safe))
    rescue
        fetched_data = {}
    end

end

def match_first_second_order_cast(first_movie_list, second_movie_list,first_name, second_name)
  matched_cast_list = Array.new
  print "\nLoading."
  first_movie_list['movies'].each do |first_movie_hash| #iterate first movies
    print "."
    #puts "first order movie :#{first_movie_hash['name']}"
    first_movie_cast_list = fetch_details(first_movie_hash['url'])
    if !first_movie_cast_list.empty?
      first_cast_list = first_movie_cast_list['cast']
      second_movie_list['movies'].each do |second_movie_hash| #iterate movies
        print "."
        #puts "second order movie :#{second_movie_hash['name']}"
        second_movie_cast_list = fetch_details(second_movie_hash['url'])
        if !second_movie_cast_list.empty?
          second_cast_list = second_movie_cast_list['cast']
          first_degree_names = first_cast_list.collect{|f| f["name"]}
          second_degree_names =  second_cast_list.collect{|s| s["name"]}
          matched_cast_list = first_degree_names & second_degree_names
          if matched_cast_list.count > 0
            puts "Matched in both #{matched_cast_list}\n"
            first_seperation = first_cast_list.select {|c| c if c["name"] == matched_cast_list[0]}
            #second_seperation = second_cast_list.select{|s| s if s["name"] == second_name }
            #puts "First data degree #{first_seperation}\n"
            #puts "Second data degree #{second_seperation}\n"
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
end

puts "Enter the 2 names with space"
degrees_array = gets.chomp.split(' ')
first_degree =  degrees_array[0]
second_degree =  degrees_array[1]
puts "First Name : #{first_degree}"
puts "Second Name : #{second_degree}"

output = assign_values(first_degree,second_degree)
if output.nil?
  puts "\nData not found for the given names\n"
end
