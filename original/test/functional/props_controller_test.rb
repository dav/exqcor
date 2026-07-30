require 'test_helper'

class PropsControllerTest < ActionController::TestCase
  setup do
    @prop = props(:one)
  end

  # test "should get index" do
  #   get :index
  #   assert_response :success
  #   assert_not_nil assigns(:props)
  # end
  # 
  # test "should get new" do
  #   get :new
  #   assert_response :success
  # end
  # 
  # test "should create prop" do
  #   assert_difference('Prop.count') do
  #     post :create, prop: { description: @prop.description, name: @prop.name, play_id: @prop.play_id }
  #   end
  # 
  #   assert_redirected_to prop_path(assigns(:prop))
  # end
  # 
  # test "should show prop" do
  #   get :show, id: @prop
  #   assert_response :success
  # end
  # 
  # test "should get edit" do
  #   get :edit, id: @prop
  #   assert_response :success
  # end
  # 
  # test "should update prop" do
  #   put :update, id: @prop, prop: { description: @prop.description, name: @prop.name, play_id: @prop.play_id }
  #   assert_redirected_to prop_path(assigns(:prop))
  # end
  # 
  # test "should destroy prop" do
  #   assert_difference('Prop.count', -1) do
  #     delete :destroy, id: @prop
  #   end
  # 
  #   assert_redirected_to props_path
  # end
end
