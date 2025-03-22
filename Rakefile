task :default => ['run:server']

LISTEN_ADDR = ENV['LISTEN_ADDR'] || ':8000'

namespace :run do

  desc "run server (default: #{LISTEN_ADDR})"
  task :server do
    system %{ GOLANG_ENV=development go run . }
    status = $?&.exitstatus || 1
  rescue Interrupt
    status = 0
  ensure
    exit status
  end

  desc 'run orbstack infra'
  task :infra do
    system %{ docker compose -f stacks/local/docker-compose.yml up --build }
    status = $?&.exitstatus || 1
  rescue Interrupt
    status = 0
  ensure
    exit status
  end
  
end

namespace :down do
  desc 'down orbstack infra'
  task :infra do
    system %{ docker compose -f stacks/local/docker-compose.yml down }
    status = $?&.exitstatus || 1
  rescue Interrupt
    status = 0
  ensure
    exit status
  end
  
end


task :command_exists, [:command] do |_, args|
  abort "#{args.command} doesn't exists" if `command -v #{args.command} > /dev/null 2>&1 && echo $?`.chomp.empty?
end
task :is_repo_clean do
  abort 'please commit your changes first!' unless `git status -s | wc -l`.strip.to_i.zero?
end
task :has_bump_my_version do
  Rake::Task['command_exists'].invoke('bump-my-version')
end


AVAILABLE_REVISIONS = %w[major minor patch].freeze
task :bump, [:revision] => [:has_bump_my_version] do |_, args|
  args.with_defaults(revision: 'patch')
  unless AVAILABLE_REVISIONS.include?(args.revision)
    abort "Please provide valid revision: #{AVAILABLE_REVISIONS.join(',')}"
  end

  system %{ bump-my-version bump #{args.revision} }
  exit $?.exitstatus
end

desc "release new version #{AVAILABLE_REVISIONS.join(',')}, default: patch"
task :release, [:revision] => [:is_repo_clean] do |_, args|
  args.with_defaults(revision: 'patch')
  Rake::Task['bump'].invoke(args.revision)
end

DOCKER_IMAGE_NAME = "bilus.org:latest"

namespace :docker do
  desc "build docker image locally"
  task :build do
    system %{
      GOOS="linux"
      GOARCH=$(go env GOARCH)
      docker build \
        --build-arg="GOOS=${GOOS}" \
        --build-arg="GOARCH=${GOARCH}" \
        --build-arg="BUILD_SHA=$(git describe --tags 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo 'not_available')" \
        --build-arg="BUILD_DATE=$(date)" \
        -t #{DOCKER_IMAGE_NAME} .
    }
    exit $?.exitstatus
  end

  desc "run docker image locally"
  task :run do
    port = LISTEN_ADDR.split(':').last
    system %{
      docker run -p "#{port}:#{port}" #{DOCKER_IMAGE_NAME}
    }
    status = $?&.exitstatus || 1
  rescue Interrupt
    status = 0
  ensure
    exit status
  end
end
