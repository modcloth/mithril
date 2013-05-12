class ServerRunner
  attr_reader :server_binary, :addr, :port, :extra_args, :logfile
  attr_reader :startup_sleep, :start_time, :server_pid, :pidfile

  def initialize(options = {})
    @start_time = options[:start] || Time.now.utc
    start = start_time.strftime('%Y%m%d%H%M%S')
    @server_binary = "#{ENV['GOPATH'].split(/:/).first}/bin/mithril-server"
    @addr = options[:port] ? ":#{options[:port]}" : ENV['ADDR']
    @port = (options[:port] || ENV['ADDR'] || '9494').to_s.gsub(/:/, '').to_i
    @logfile = File.expand_path(
      "../../log/mithril-server-#{start}-#{port}.log",
      __FILE__
    )
    @pidfile = (options[:pidfile] || "mithril-server-#{@port}.pid")
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
    announce! "Starting mithril server with address #{addr}"
    @server_pid = Process.spawn(
      "#{server_binary} -a #{addr} -p #{pidfile} #{extra_args}" <<
        ">> #{logfile} 2>&1"
    )
    sleep @startup_sleep
    @server_pid
  end

  def stop
    real_pid = Integer(File.read(pidfile).chomp) rescue nil
    if server_pid && real_pid
      announce! "Stopping mithril server with address #{addr} " <<
                "(shell PID=#{server_pid}, server PID=#{real_pid})"

      [real_pid, server_pid].each do |pid|
        Process.kill(:TERM, pid) rescue nil
      end
    end
  end

  def dump_log
    announce! "Dumping #{logfile}"
    File.read(logfile).split($/).each do |line|
      announce! "--> #{line}"
    end
  end
end
