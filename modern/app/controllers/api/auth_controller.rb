class Api::AuthController < ApplicationController
  allow_unauthenticated_access only: %i[sign_up sign_in sign_out]

  def sign_up
    user = User.new(user_params)

    if user.save
      token = UserToken.issue!(user, name: params[:device_name])
      expires_at = user.user_tokens.order(created_at: :desc).pick(:expires_at)&.iso8601
      render json: { token: token, user_id: user.id, expires_at: expires_at }, status: :created
    else
      render_errors(user)
    end
  end

  def sign_in
    if (user = User.authenticate_by(params.permit(:email_address, :password)))
      token = UserToken.issue!(user, name: params[:device_name])
      expires_at = user.user_tokens.order(created_at: :desc).pick(:expires_at)&.iso8601
      render json: { token: token, user_id: user.id, expires_at: expires_at }, status: :ok
    else
      render json: { error: "Invalid email or password" }, status: :unauthorized
    end
  end

  def sign_out
    token = bearer_token
    user_token = UserToken.find_by_token(token)
    user_token&.destroy
    head :no_content
  end

  private

  def user_params
    params.require(:user).permit(:email_address, :password, :password_confirmation)
  end
end
