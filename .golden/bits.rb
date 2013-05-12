require 'net/http'

GOLD = "\033\[33;1m"
RESET = "\033\[0m"
GREEN = "\033\[32m"
RED = "\033\[31m"
BRIGHT_GREEN = "\033\[32;1m"
BRIGHT_RED = "\033\[31;1m"

def announce!(something)
  $stderr.puts "#{GOLD}golden#{RESET}: #{GREEN}#{something}#{RESET}"
end

class MiniTestReporter
  def puts(*args)
    args.each { |arg| announce! arg }
  end

  alias print puts
end

def run_psql(command, options = {})
  if !options[:db]
    raise 'No :db option given!'
  end
  exe = %Q{psql -d #{options[:db]} -t -c "#{command}"}
  if options[:user]
    exe = %Q{psql -U #{options[:user]} -d #{options[:db]} -t -c "#{command}"}
  end
  announce! "psql executing: #{exe}"
  output = `#{exe}`.chomp
  return [output, $?]
end

def post_requests(options = {})
  path = "/#{options[:exchange] || 'foop'}/#{options[:routing_key] || 'fwap'}"
  port = options[:port] || 8371
  request = Net::HTTP::Post.new(path)
  request.content_type = case options[:type]
                         when :json
                           'application/json'
                         when :xml
                           'application/xml'
                         else
                           'application/octet-stream'
                         end

  Integer(options[:count] || 1).times do |n|
    response = Net::HTTP.start('localhost', port) do |http|
      request.body = case options[:type]
                     when :json
                       %Q/{"flume":"sandbag #{rand}"}/
                     when :xml
                       %Q(<flume sandbag="#{rand}"></flume>)
                     else
                       %Q(\x99f\x81l\x78u\x93m\x33e\x90)
                     end
      http.request(request)
    end

    if response.code != '204'
      raise "Failed POST: #{response.inspect}"
    end

    announce! "POSTed #{request.body.inspect} to #{request.path.inspect}"
  end
end
