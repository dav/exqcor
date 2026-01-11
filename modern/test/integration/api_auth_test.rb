require "test_helper"

class ApiAuthTest < ActionDispatch::IntegrationTest
  test "sign up returns token" do
    post "/api/sign_up", params: { user: { email_address: "api@example.com", password: "secret123", password_confirmation: "secret123" } }, as: :json

    assert_response :created
    payload = JSON.parse(response.body)
    assert payload["token"].present?
    assert payload["user_id"].present?
    assert payload["expires_at"].present?
  end

  test "sign in returns token" do
    user = User.create!(email_address: "login@example.com", password: "secret123", password_confirmation: "secret123")

    post "/api/sign_in", params: { email_address: user.email_address, password: "secret123" }, as: :json

    assert_response :ok
    payload = JSON.parse(response.body)
    assert payload["token"].present?
    assert payload["expires_at"].present?
  end

  test "api requests accept bearer token" do
    user = User.create!(email_address: "token@example.com", password: "secret123", password_confirmation: "secret123")
    token = UserToken.issue!(user, name: "test")

    get "/api/scripts", headers: { "Authorization" => "Bearer #{token}" }, as: :json

    assert_response :ok
  end
end
