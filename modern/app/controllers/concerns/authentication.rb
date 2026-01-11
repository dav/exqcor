module Authentication
  extend ActiveSupport::Concern

  included do
    before_action :require_authentication
    helper_method :authenticated?, :current_user
  end

  class_methods do
    def allow_unauthenticated_access(**options)
      skip_before_action :require_authentication, **options
    end
  end

  private
    def authenticated?
      resume_session
    end

    def require_authentication
      resume_session || request_authentication
    end

    def resume_session
      return true if resume_api_token

      Current.session ||= find_session_by_cookie
      Current.user ||= Current.session&.user
    end

    def find_session_by_cookie
      Session.find_by(id: cookies.signed[:session_id]) if cookies.signed[:session_id]
    end

    def request_authentication
      if request.format.json?
        render json: { error: "Unauthorized" }, status: :unauthorized
      else
        session[:return_to_after_authenticating] = request.url
        redirect_to new_session_path
      end
    end

    def after_authentication_url
      session.delete(:return_to_after_authenticating) || root_url
    end

    def start_new_session_for(user)
      user.sessions.create!(user_agent: request.user_agent, ip_address: request.remote_ip).tap do |session|
        Current.session = session
        Current.user = session.user
        cookies.signed.permanent[:session_id] = { value: session.id, httponly: true, same_site: :lax }
      end
    end

    def terminate_session
      Current.session.destroy
      cookies.delete(:session_id)
    end

    def resume_api_token
      return false unless request.format.json?

      token = bearer_token
      return false if token.blank?

      user_token = UserToken.find_by_token(token)
      return false if user_token.nil? || user_token.expired?

      if should_rotate_token?(user_token)
        raw_token, new_token = user_token.rotate!
        response.set_header("X-Auth-Token", raw_token)
        response.set_header("X-Auth-Token-Expires-At", new_token.expires_at.iso8601)
        Current.api_token = new_token
        Current.user = new_token.user
        true
      else
        user_token.touch_last_used!
        Current.api_token = user_token
        Current.user = user_token.user
        true
      end
    end

    def bearer_token
      auth = request.authorization.to_s
      return nil unless auth.start_with?("Bearer ")

      auth.split(" ", 2).last
    end

    def current_user
      Current.user
    end

    def should_rotate_token?(user_token)
      user_token.expires_at.present? && user_token.expires_at <= 7.days.from_now
    end
end
