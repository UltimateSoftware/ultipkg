require 'fileutils'

APP_NAME = 'ultipkg'
TARGETS = [
  {os: "linux", arch: "amd64"},
  {os: "linux", arch: "386"},
  {os: "windows", arch: "amd64"},
  {os: "windows", arch: "386"},
  {os: "darwin", arch: "amd64"},
]

$process_list = []
at_exit { terminate_all }

task :default => [:test, :run]

desc "Clean packages created in bin"
task :clean do
  system "rm -rf ./bin/"
end

desc "Run go generate"
task :generate do
  step "Generating for #{APP_NAME}-#{version}" do
    go_generate
  end
end

desc "Run tests"
task :test => :generate do |t|
  puts "Running tests...\n"
  exit(go_test)
end

desc "Build to run locally"
task :build => :generate do
  step "Building #{APP_NAME}-#{version}" do
    go_build
  end
end

desc "Build and Run"
task :run => :build do
  puts 'Starting application...'
  go_run
end

desc "Watch and re-run the application"
task :watch, [:paths] => :build do |t, args|

  puts 'Starting application...'
  pid = go_run!

  watch "**/*.go", "*.go", args.paths do
    terminate(pid)

    step "Regenerating" do
      go_generate
    end

    puts "Retesting..."
    go_test(SHORT)

    step "Rebuilding" do
      go_build
    end

    pid = go_run!
  end
end

desc "Build for release"
task :release do
  puts "Building #{APP_NAME}-#{version} for release..."

  TARGETS.each do |t|
    step "\tBuilding for: #{t}" do
      exit(1) unless system %Q{GOOS=#{t[:os]} GOARCH=#{t[:arch]} go build -o "bin/#{APP_NAME}-#{version}-#{t[:os]}-#{t[:arch]}" -ldflags "#{ldflags}"}
    end
  end
end

def version
  ENV['VERSION'] || `git describe --tags --dirty --always`.chomp
end

def run(cmd)
  pid = Process.spawn(cmd)
  $process_list << pid
  pid
end

def wait(pid)
  Process.wait(pid)
  $process_list.delete_if {|pid| pid == pid}
  return $?.exitstatus
end

def run_and_wait(cmd)
  wait(run(cmd))
end

def terminate(pid, sig = "SIGINT")
  begin
    Process.kill(sig, pid)
  rescue Errno::ESRCH
    # noop
  end
end

def terminate_all(sig = "SIGINT")
  $process_list.each { |pid| terminate(pid, sig) }
end

def cancel_clean(msg)
   # Trap ctrl-c and display a nice message instead of Rake barfing
  # a stack trace into our terminal.
  Signal.trap('SIGINT') do
    puts "\n#{msg}"
    exit(0)
  end
end

def go_version
  `go version`.strip.split(' ')[2].split('')[2..4].join.to_f
end

SHORT = true
def go_test(short = false)
  flags = []
  flags << "-short" if short

  run_and_wait("go test #{flags.join(' ')} $(go list ./... | grep -v /vendor/)")
end

def go_generate
  cancel_clean('Halting generate...')
  exit(1) unless system %Q{go generate ./...}
end

def go_build
  cancel_clean('Halting build...')
  exit(1) unless system %Q{go build -ldflags "#{ldflags}" -o "bin/#{APP_NAME}"}
end

def go_run
  wait(go_run!)
end

def go_run!
  cancel_clean('Shutting down...')
  run("./bin/#{APP_NAME}")
end

def watch(*paths)
  with_commands :fswatch do
    cancel_clean('Stopping watcher...')
    while
      system %Q{fswatch -r -1 #{paths.join(' ')}}
      yield
    end
  end
end

def with_commands(*tools)
  tools.each do |t|
    brew_install(t)
  end

  yield if block_given?
end

def installed?(tool)
  `which #{tool} > /dev/null; echo $?`.chomp == "0"
end

def brew_install(tool)
  cancel_clean('Stopping brew install...')
  return if installed?(tool)
  step "Installing missing tool '#{tool}' with brew" do
    exit(1) unless `brew install #{tool}`
  end
end

def step(msg)
  print "#{msg}... "
  if block_given?
    yield
    print "done."
  end
  print "\n"
end

def ldflags
  if go_version >= 1.5
    "-X main.Version=#{version}"
  else
    "-X main.Version '#{version}'"
  end
end


