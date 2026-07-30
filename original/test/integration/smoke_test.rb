require 'test_helper'

class SmokeTest < ActionDispatch::IntegrationTest
  test 'root page loads' do
    get '/'
    assert_response :success
  end
end
