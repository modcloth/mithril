class ServerRunner
  attr_reader :server_binary, :addr, :port, :amqp_uri, :extra_args, :logfile
  attr_reader :startup_sleep, :start_time, :server_pid

  def initialize(options = {})
    @start_time = options[:start] || Time.now.utc
    start = start_time.strftime('%Y%m%d%H%M%S')
    @server_binary = "#{ENV['GOPATH'].split(/:/).first}/bin/mithril-server"
    @addr = options[:port] ? ":#{options[:port]}" : ENV['ADDR']
    @port = (options[:port] || ENV['ADDR'] || '9494').to_s.gsub(/:/, '').to_i
    @logfile = File.expand_path(
      "../../.artifacts/mithril-server-#{start}-#{port}.log",
      __FILE__
    )
    @amqp_uri = options[:amqp_uri] || 'amqp://guest:guest@localhost:5672'
    @extra_args = options[:extra_args] || ''
    @startup_sleep = Float(
      options[:startup_sleep] || ENV['MITHRIL_STARTUP_SLEEP'] || 0.5
    )

    if !File.exist?(server_binary)
      raise "Can't locate `mithril-server` binary! " <<
            "(it's not here: #{server_binary.inspect})"
    end
  end

  def start
    announce! "Starting mithril server with address #{addr.inspect}, " <<
    "amqp uri #{amqp_uri.inspect}"
    @server_pid = Process.spawn(
      "#{server_binary} s #{extra_args} -b #{addr} -a #{amqp_uri} " <<
      ">> #{logfile} 2>&1"
    ) + 1 # Unknown why we need to add one here
    sleep @startup_sleep
    @server_pid
  end

  def stop
    announce! "Stopping mithril server with address #{addr} " <<
    "(shell PID=#{server_pid})"

    Process.kill(:TERM, server_pid)
  end

  def dump_log
    announce! "Dumping #{logfile}"
    File.read(logfile).split($/).each do |line|
      announce! "--> #{line}"
    end
  end
end
