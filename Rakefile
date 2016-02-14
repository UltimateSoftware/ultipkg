require 'fileutils'

APP_NAME = 'ultipkg'

task :default => [:test, :run]

desc "Run tests"
task :test do |t|
  puts "Running tests...\n"
  unless system %Q{go test ./...}
    exit(1)
  end
end

desc "Build for release"
task :release do
  puts "Building #{APP_NAME}-#{version} for release ..."

  [
    {os: "linux", arch: "amd64"},
    {os: "linux", arch: "386"},
    {os: "windows", arch: "amd64"},
    {os: "windows", arch: "386"},
    {os: "darwin", arch: "amd64"},
  ].each do |t|
    puts "\tBuilding for: #{t}"

    unless system %Q{GOOS=#{t[:os]} GOARCH=#{t[:arch]} go build -o "bin/#{APP_NAME}-#{version}-#{t[:os]}-#{t[:arch]}" -ldflags "#{ldflags}"}
      exit(1)
    end
  end
end

desc "Build to run locally"
task :build => :env do
  cancel_clean('Halting build...')
  puts "Building #{APP_NAME}-#{version} ..."
  unless system %Q{go build -ldflags "#{ldflags}" -o "bin/#{APP_NAME}"}
    exit(1)
  end
end

desc "Clean packages created in bin"
task :clean do
  system "rm -rf ./bin/"
end

desc "Build and Run"
task :run => [:env, :build] do
  cancel_clean('Shutting down ...')
  puts 'Starting application ...'
  run("./bin/#{APP_NAME}")
end

task :env do
  output = []

  {
    'format' => 'LOGGING_FORMAT',
    'level' => 'LOGGING_LEVEL',
    'addr' => 'ADDR',
    'addr_tls' => 'ADDR_TLS',
    'domain' => 'DOMAIN',
    'version' => 'VERSION',
    'certificate' => 'SSL_CERTIFICATE',
    'key' => 'SSL_PRIVATEKEY',
  }.each do |flag, key|
    unless ENV[flag].nil?
      ENV[key] = ENV[flag]
      output << "#{key} = #{ENV[flag]}"
    end
  end

  unless output.empty?
    puts 'Overriding default environment:'
    output.each {|l| puts "\t#{l}" }
  end
end

def version
  ENV['VERSION'] || `git describe --tags --dirty --always`.chomp
end

def run(cmd, &blk)
  subprocess(cmd) do |stdout, stderr, thread|
    if stdout
      if blk.nil?
        puts stdout
      else
        blk.call(stdout)
      end
    end
    puts stderr if stderr
  end
end

def subprocess(cmd, &block)
  # see: http://stackoverflow.com/a/1162850/83386
  Open3.popen3(cmd) do |stdin, stdout, stderr, thread|
    # read each stream from a new thread
    { :out => stdout, :err => stderr }.each do |key, stream|
      Thread.new do
        until (line = stream.gets).nil? do
          # yield the block depending on the stream
          if key == :out
            yield line, nil, thread if block_given?
          else
            yield nil, line, thread if block_given?
          end
        end
      end
    end

    thread.join # don't exit until the external process is done
  end
end

def cancel_clean(msg)
   # Trap ctrl-c and display a nice message instead of Rake barfing
  # a stack trace into our terminal.
  Signal.trap('SIGINT') do
    puts "\n#{msg}"
    exit(0)
  end
end

def ldflags
  if go15?
    "-X main.Version=#{version}"
  else
    "-X main.Version '#{version}'"
  end
end

def go15?
  !!(`go version` =~ /go1.5/)
end
