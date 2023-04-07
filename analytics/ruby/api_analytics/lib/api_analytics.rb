# frozen_string_literal: true

require 'uri'
require 'net/http'
require 'json'

module Analytics
  class Middleware
    def initialize(app, api_key)
      @app = app
      @api_key = api_key
      @requests = Array.new
      @last_posted = Time.now
    end

    def call(env)
      start = Time.now
      status, headers, response = @app.call(env)

      request_data = {
        hostname: env['HTTP_HOST'],
        ip_address: env['REMOTE_ADDR'],
        path: env['REQUEST_PATH'],
        user_agent: env['HTTP_USER_AGENT'],
        method: env['REQUEST_METHOD'],
        status: status,
        response_time: (Time.now - start).to_f.round,
        created_at: Time.now.utc.iso8601
      }

      log_request(request_data)

      [status, headers, response]
    end
    
    private

    def post_requests(api_key, requests, framework)
      payload = {
        api_key: api_key,
        requests: requests,
        framework: framework
      }
      uri = URI('https://www.apianalytics-server.com/api/log-request')
      res = Net::HTTP.post(uri, payload.to_json)
    end
    
    def log_request(request_data)
      if @api_key.empty?
        return
      end
      now = Time.now
      @requests.push(request_data)
      if (now - @last_posted) > 5
        requests = @requests.dup
        Thread.new {
          post_requests(@api_key, requests, @framework)
        }
        @requests = Array.new
        @last_posted = now
    end
    end
  end

  private_constant :Middleware
  
  class Rails < Middleware
    def initialize(app, api_key)
      super(app, api_key)
      @framework = "Rails"
    end
  end

  class Sinatra < Middleware
    def initialize(app, api_key)
      super(app, api_key)
      @framework = "Sinatra"
    end
  end
end

