class ApplicationController < ActionController::Base
  include Authentication
  # Only allow modern browsers supporting webp images, web push, badges, import maps, CSS nesting, and CSS :has.
  allow_browser versions: :modern

  protect_from_forgery with: :exception, unless: -> { request.format.json? }
  protect_from_forgery with: :null_session, if: -> { request.format.json? }

  private

  def render_errors(resource)
    render json: { errors: resource.errors.full_messages }, status: :unprocessable_entity
  end

  def require_admin
    return if current_user&.admin?

    if request.format.json?
      render json: { error: "Forbidden" }, status: :forbidden
    else
      redirect_to root_path, alert: "Admin access required."
    end
  end
end
