require 'test_helper'

class LinesControllerTest < ActionController::TestCase
  setup do
    @line = lines(:p1s1ss1l2)
    @sub_section = @line.sub_section
    @section = @sub_section.section
    @play = @section.play
  end

  # test "should get index" do
  #   get :index
  #   assert_response :success
  #   assert_not_nil assigns(:lines)
  # end
  # 
  # test "should get new" do
  #   get :new
  #   assert_response :success
  # end
  # 
  # test "should create line" do
  #   assert_difference('Line.count') do
  #     post :create, line: { character_id: @line.character_id, ordering: @line.ordering, sub_section_id: @line.sub_section_id, text: @line.text }
  #   end
  # 
  #   assert_redirected_to line_path(assigns(:line))
  # end
  # 
  # test "should show line" do
  #   get :show, id: @line
  #   assert_response :success
  # end
  # 
  # test "should get edit" do
  #   get :edit, id: @line
  #   assert_response :success
  # end
  # 
  # test "should update line" do
  #   put :update, id: @line, line: { character_id: @line.character_id, ordering: @line.ordering, sub_section_id: @line.sub_section_id, text: @line.text }
  #   assert_redirected_to line_path(assigns(:line))
  # end
  # 
  test "should destroy line" do
    assert_difference('Line.count', -1) do
      delete :destroy, id: @line, :format => 'json'
    end
 
    assert_response :success
  end
end
