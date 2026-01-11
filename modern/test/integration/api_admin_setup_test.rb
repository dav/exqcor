require "test_helper"

class ApiAdminSetupTest < ActionDispatch::IntegrationTest
  setup do
    user = User.create!(email_address: "writer@example.com", password: "secret123", password_confirmation: "secret123")
    @token = UserToken.issue!(user, name: "test")
  end

  test "api script create ensures vosd" do
    post "/api/scripts",
         params: { script: { title: "API Script" } },
         headers: { "Authorization" => "Bearer #{@token}" },
         as: :json

    assert_response :created

    script = Script.find(JSON.parse(response.body)["id"])
    assert script.VOSD
    assert_equal "vosd", script.VOSD.role
  end

  test "api section create ensures vosd character_section" do
    script = Script.create!(title: "API Script")

    post "/api/scripts/#{script.id}/sections",
         params: { section: { name: "Scene 1" } },
         headers: { "Authorization" => "Bearer #{@token}" },
         as: :json

    assert_response :created

    section = Section.find(JSON.parse(response.body)["id"])
    assert_equal script.id, section.script_id
    assert_equal 1, section.character_sections.count
    assert_equal "vosd", section.characters.first.role
  end
end
