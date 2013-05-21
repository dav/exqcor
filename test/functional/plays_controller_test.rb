require 'test_helper'

class PlaysControllerTest < ActionController::TestCase
  setup do
    @play = plays(:one)
  end

  # test "should get index" do
  #    get :index
  #    assert_response :success
  #    assert_not_nil assigns(:plays)
  #  end
  # 
  #  test "should get new" do
  #    get :new
  #    assert_response :success
  #  end
  # 
  #  test "should create play" do
  #    assert_difference('Play.count') do
  #      post :create, play: { description: @play.description, title: @play.title }
  #    end
  # 
  #    assert_redirected_to play_path(assigns(:play))
  #  end
  # 
  #  test "should show play" do
  #    get :show, id: @play
  #    assert_response :success
  #  end
  # 
  #  test "should get edit" do
  #    get :edit, id: @play
  #    assert_response :success
  #  end
  # 
  #  test "should update play" do
  #    put :update, id: @play, play: { description: @play.description, title: @play.title }
  #    assert_redirected_to play_path(assigns(:play))
  #  end
  # 
  #  test "should destroy play" do
  #    assert_difference('Play.count', -1) do
  #      delete :destroy, id: @play
  #    end
  # 
  #    assert_redirected_to plays_path
  #  end
end
